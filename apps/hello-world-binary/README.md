# Hello World Go Binary DXT Extension

Go 言語で実装された Hello World DXT 拡張機能です。

## 機能

- **go_say_hello**: 指定された名前に対して挨拶を返します（Go 版）
- **go_get_time**: 現在の時刻を返します（Go 版）

## 特徴

✅ **完全自己完結**: Node.js 環境不要  
✅ **高速起動**: バイナリ実行で高速  
✅ **クロスプラットフォーム**: macOS/Windows/Linux 対応  
✅ **軽量**: 6.6MB 単一バイナリ

## Node.js 版との違い

| 項目         | Node.js 版              | Go Binary 版                  |
| ------------ | ----------------------- | ----------------------------- |
| **ツール名** | `say_hello`, `get_time` | `go_say_hello`, `go_get_time` |
| **依存関係** | Node.js 必要            | 完全自己完結                  |
| **起動速度** | 普通                    | 高速                          |

## 使用方法

Claude Desktop にインストール後、以下のように使用できます：

```
go_say_hello を使って私の名前は太郎です。挨拶してください。
```

```
go_get_time を使って現在の時刻を教えてください
```

## ローカルテスト

```bash
cd server
./hello-world-binary
```

## パッケージング

```bash
dxt pack
```

## ⚠️ バイナリの実行権限がないエラーが発生する可能性あり

実際、筆者がローカル環境で試してみたところ、 `EACCESS` エラーが発生した。
`chmod +x /Library/Application\ Support/Claude/Claude\ Extensions/local.dxt.<name>.hello-world-go/server/hello-world-binary` のようなコマンドで実行権限を与えてあげれば実行可能になったが、配布する上での障壁になるのでどうにかならないか・・・😢

## 技術仕様

- **言語**: Go
- **SDK**: [公式 MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- **通信**: stdin/stdout (MCP Protocol)
- **バイナリサイズ**: 6.6MB
- **DXT サイズ**: 3.6MB
