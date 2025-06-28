# DXT クイックリファレンス

## 🚀 基本コマンド

### Node.js 関連

```bash
# MCP SDKインストール
npm install @modelcontextprotocol/sdk

# 依存関係インストール
npm install

# ローカルテスト実行
node server/index.js
```

### DXT CLI

```bash
# プロジェクト初期化
dxt init [project-name]

# マニフェスト検証
dxt validate

# DXTファイル作成
dxt pack

# ヘルプ表示
dxt --help

# バージョン確認
dxt --version
```

## 📝 Manifest.json テンプレート

### 基本テンプレート

```json
{
  "dxt_version": "0.1",
  "name": "extension-name",
  "display_name": "拡張機能の表示名",
  "version": "1.0.0",
  "description": "拡張機能の説明",
  "author": {
    "name": "作成者名",
    "email": "author@example.com"
  },
  "license": "MIT",
  "server": {
    "type": "node",
    "entry_point": "server/index.js",
    "mcp_config": {
      "command": "node",
      "args": ["${__dirname}/server/index.js"]
    }
  }
}
```

### ユーザー設定付きテンプレート

```json
{
  "dxt_version": "0.1",
  "name": "extension-with-config",
  "display_name": "設定可能な拡張機能",
  "version": "1.0.0",
  "description": "ユーザー設定を持つ拡張機能",
  "author": {
    "name": "作成者名"
  },
  "server": {
    "type": "node",
    "entry_point": "server/index.js",
    "mcp_config": {
      "command": "node",
      "args": ["${__dirname}/server/index.js"],
      "env": {
        "API_KEY": "${user_config.api_key}",
        "DEBUG_MODE": "${user_config.debug_mode}"
      }
    }
  },
  "user_config": {
    "api_key": {
      "type": "string",
      "title": "API Key",
      "description": "外部サービスのAPIキー",
      "sensitive": true,
      "required": true
    },
    "debug_mode": {
      "type": "string",
      "title": "Debug Mode",
      "description": "デバッグモードの有効化",
      "enum": ["true", "false"],
      "default": "false"
    }
  }
}
```

## 🔧 MCP サーバー テンプレート

### 基本サーバー

```javascript
import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  ListToolsRequestSchema,
  CallToolRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";

const server = new Server(
  {
    name: "Extension Server",
    version: "1.0.0",
  },
  {
    capabilities: {
      tools: {},
    },
  }
);

// ツール一覧の定義
server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
        name: "tool_name",
        description: "ツールの説明",
        inputSchema: {
          type: "object",
          properties: {
            parameter: {
              type: "string",
              description: "パラメータの説明",
            },
          },
          required: ["parameter"],
        },
      },
    ],
  };
});

// ツール実行の処理
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;

  switch (name) {
    case "tool_name":
      return {
        content: [
          {
            type: "text",
            text: `処理結果: ${args.parameter}`,
          },
        ],
      };
    default:
      throw new Error(`Unknown tool: ${name}`);
  }
});

// サーバー起動
async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
}

main().catch(console.error);
```

## 🔑 よく使う設定パターン

### セキュアな設定

```json
{
  "user_config": {
    "api_key": {
      "type": "string",
      "title": "API Key",
      "description": "APIキー",
      "sensitive": true,
      "required": true,
      "validation": {
        "pattern": "^[a-zA-Z0-9]{32}$",
        "message": "32文字の英数字で入力してください"
      }
    }
  }
}
```

### ディレクトリパス設定

```json
{
  "user_config": {
    "workspace_path": {
      "type": "string",
      "title": "Workspace Path",
      "description": "作業用ディレクトリのパス",
      "default": "/Users/USERNAME/Documents"
    }
  }
}
```

### 選択肢設定

```json
{
  "user_config": {
    "log_level": {
      "type": "string",
      "title": "Log Level",
      "description": "ログレベル",
      "enum": ["debug", "info", "warn", "error"],
      "default": "info"
    }
  }
}
```

## 🛠 ツール定義パターン

### ファイル操作ツール

```javascript
{
  name: 'read_file',
  description: 'ファイルを読み取ります',
  inputSchema: {
    type: 'object',
    properties: {
      file_path: {
        type: 'string',
        description: 'ファイルパス'
      },
      encoding: {
        type: 'string',
        enum: ['utf8', 'ascii', 'base64'],
        default: 'utf8'
      }
    },
    required: ['file_path']
  }
}
```

### データベースクエリツール

```javascript
{
  name: 'execute_query',
  description: 'SQLクエリを実行します',
  inputSchema: {
    type: 'object',
    properties: {
      query: {
        type: 'string',
        description: 'SQLクエリ'
      },
      limit: {
        type: 'number',
        description: '結果の最大件数',
        minimum: 1,
        maximum: 100,
        default: 10
      }
    },
    required: ['query']
  }
}
```

### Web API 呼び出しツール

```javascript
{
  name: 'call_api',
  description: 'Web APIを呼び出します',
  inputSchema: {
    type: 'object',
    properties: {
      endpoint: {
        type: 'string',
        description: 'APIエンドポイント'
      },
      method: {
        type: 'string',
        enum: ['GET', 'POST', 'PUT', 'DELETE'],
        default: 'GET'
      },
      body: {
        type: 'object',
        description: 'リクエストボディ'
      }
    },
    required: ['endpoint']
  }
}
```

## 🔍 デバッグとトラブルシューティング

### 一般的なエラーと解決策

#### インストールエラー

```bash
# Node.jsバージョン確認
node --version  # v16以上が必要

# DXT CLIの再インストール
npm uninstall -g @anthropic-ai/dxt
npm install -g @anthropic-ai/dxt
```

#### パッケージングエラー

```bash
# マニフェスト検証
dxt validate

# node_modulesの再構築
cd server
rm -rf node_modules package-lock.json
npm install
cd ..
dxt pack
```

#### 拡張機能が表示されない

1. Claude Desktop の再起動
2. 設定 > Extensions で確認
3. .dxt ファイルの再インストール

### ログ確認方法

#### macOS

```bash
# Claude Desktopのログ
tail -f ~/Library/Logs/Claude\ Desktop/main.log

# システムログ
log stream --predicate 'subsystem == "com.anthropic.claude.desktop"'
```

#### Windows

```powershell
# イベントビューアーでClaude Desktopのログを確認
Get-WinEvent -LogName Application | Where-Object {$_.ProviderName -eq "Claude Desktop"}
```

## 📚 よく使うコードスニペット

### エラーハンドリング

```javascript
try {
  // 処理
  const result = await someOperation();
  return {
    content: [
      {
        type: "text",
        text: `成功: ${result}`,
      },
    ],
  };
} catch (error) {
  return {
    content: [
      {
        type: "text",
        text: `エラーが発生しました: ${error.message}`,
      },
    ],
    isError: true,
  };
}
```

### 設定値の取得

```javascript
// 環境変数から設定値を取得
const apiKey = process.env.API_KEY;
const debugMode = process.env.DEBUG_MODE === "true";
const workspacePath = process.env.WORKSPACE_PATH || "/default/path";
```

### ファイル操作

```javascript
import fs from "fs/promises";
import path from "path";

// ファイル読み取り
const content = await fs.readFile(filePath, "utf8");

// ファイル書き込み
await fs.writeFile(filePath, content, "utf8");

// ディレクトリ作成
await fs.mkdir(dirPath, { recursive: true });

// ファイル存在確認
try {
  await fs.access(filePath);
  // ファイルが存在する
} catch {
  // ファイルが存在しない
}
```

## 🌐 リソースリンク

### 公式ドキュメント

- [DXT GitHub Repository](https://github.com/anthropics/dxt)
- [MCP Specification](https://modelcontextprotocol.io/)
- [Anthropic Blog Post](https://www.anthropic.com/engineering/desktop-extensions)

### 開発ツール

- [VS Code DXT Extension](https://marketplace.visualstudio.com/items?itemName=anthropic.dxt)
- [Claude Desktop](https://claude.ai/download)
- [Node.js](https://nodejs.org/)

### コミュニティ

- [DXT Discord](https://discord.gg/dxt-community)
- [GitHub Discussions](https://github.com/anthropics/dxt/discussions)
- [Stack Overflow Tag](https://stackoverflow.com/questions/tagged/dxt)

---

このクイックリファレンスを参考に、効率的な DXT 開発を行ってください！
