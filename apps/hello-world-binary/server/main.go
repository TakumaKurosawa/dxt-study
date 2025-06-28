package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// デバッグ用のヘルパー関数
func getWorkingDir() string {
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return "unknown"
}

func getExecutablePath() string {
	if exe, err := os.Executable(); err == nil {
		return exe
	}
	return "unknown"
}

// say_hello ツール用のパラメータ構造体
type SayHelloParams struct {
	Name     string `json:"name"`
	Language string `json:"language,omitempty"`
}

// get_time ツール用のパラメータ構造体
type GetTimeParams struct {
	Format string `json:"format,omitempty"`
}

// say_hello ツールの実装
func SayHello(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[SayHelloParams]) (*mcp.CallToolResultFor[any], error) {
	name := params.Arguments.Name
	language := params.Arguments.Language
	
	if language == "" {
		language = "japanese"
	}

	var greeting string
	switch language {
	case "english":
		greeting = fmt.Sprintf("Hello, %s!", name)
	default:
		greeting = fmt.Sprintf("こんにちは、%sさん！", name)
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: greeting}},
	}, nil
}

// get_time ツールの実装
func GetTime(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[GetTimeParams]) (*mcp.CallToolResultFor[any], error) {
	format := params.Arguments.Format
	if format == "" {
		format = "default"
	}

	now := time.Now()
	var timeStr string

	switch format {
	case "rfc3339":
		timeStr = now.Format(time.RFC3339)
	case "unix":
		timeStr = fmt.Sprintf("%d", now.Unix())
	default:
		timeStr = now.Format("2006-01-02 15:04:05")
	}

	result := fmt.Sprintf("現在時刻: %s", timeStr)
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: result}},
	}, nil
}

func main() {
	// デバッグ情報をstderrに出力（Claude Desktopのログに表示される）
	fmt.Fprintf(os.Stderr, "[DEBUG] Go MCP server starting...\n")
	fmt.Fprintf(os.Stderr, "[DEBUG] Working directory: %s\n", getWorkingDir())
	fmt.Fprintf(os.Stderr, "[DEBUG] Executable path: %s\n", getExecutablePath())
	
	// MCPサーバーの作成
	server := mcp.NewServer("Hello World Go Server", "1.0.0", nil)
	
	// ツールの追加
	server.AddTools(
		// go_say_hello ツール (重複回避のためプレフィックス追加)
		mcp.NewServerTool("go_say_hello", "Go版: 指定された名前に対して挨拶を返します", SayHello, mcp.Input(
			mcp.Property("name", mcp.Description("挨拶する相手の名前"), mcp.Required(true)),
			mcp.Property("language", mcp.Description("挨拶の言語 (japanese, english)"), mcp.Enum("japanese", "english")),
		)),
		
		// go_get_time ツール (重複回避のためプレフィックス追加)
		mcp.NewServerTool("go_get_time", "Go版: 現在の時刻を返します", GetTime, mcp.Input(
			mcp.Property("format", mcp.Description("時刻のフォーマット (default, rfc3339, unix)"), mcp.Enum("default", "rfc3339", "unix")),
		)),
	)

	log.Printf("Hello World Go MCP server starting...")
	fmt.Fprintf(os.Stderr, "[DEBUG] About to start MCP server...\n")
	
	// deferでサーバー終了時のログ
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Server panicked: %v\n", r)
		}
		fmt.Fprintf(os.Stderr, "[DEBUG] Server shutting down...\n")
	}()

	// サーバーの起動（stdio transport）
	fmt.Fprintf(os.Stderr, "[DEBUG] Starting stdio transport...\n")
	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Server run failed: %v\n", err)
		log.Fatalf("Server error: %v", err)
	}
	fmt.Fprintf(os.Stderr, "[DEBUG] Server run completed normally\n")
}
