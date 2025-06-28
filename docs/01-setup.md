# 基本セットアップ

## 📋 前提条件

### 必要な環境

- **Node.js**: v16 以上（推奨: v18 以上）
- **npm**: Node.js に付属
- **Claude Desktop**: 最新版

### 環境確認

```bash
# Node.jsのバージョン確認
node --version

# npmのバージョン確認
npm --version
```

## 🛠 DXT CLI ツールのインストール

### グローバルインストール

```bash
npm install -g @anthropic-ai/dxt
```

### インストール確認

```bash
# DXT CLIのバージョン確認
dxt --version

# 利用可能なコマンド一覧
dxt --help
```

期待される出力:

```
Usage: dxt [command] [options]

Commands:
  init     Create a new DXT extension
  pack     Package extension into .dxt file
  validate Validate manifest.json
  help     Show help information

Options:
  --version  Show version number
  --help     Show help
```

## 🎯 Claude Desktop の設定

### 1. Claude Desktop のダウンロード

[公式サイト](https://claude.ai/download)から最新版をダウンロードしてインストールします。

### 2. Extensions 設定の確認

1. Claude Desktop を起動
2. 設定（Settings）を開く
3. "Extensions" セクションを確認
4. "Install Extension" ボタンが表示されることを確認

### 3. 開発者モードの有効化（オプション）

開発中のテストを容易にするため、以下の設定を有効にできます：

**macOS:**

```bash
defaults write com.anthropic.claude.desktop DeveloperMode -bool true
```

**Windows:**
レジストリエディタで以下を設定：

```
[HKEY_CURRENT_USER\Software\Anthropic\Claude Desktop]
"DeveloperMode"=dword:00000001
```

## 📁 ワークスペースの準備

### プロジェクトディレクトリの作成

```bash
# 作業用ディレクトリを作成
mkdir dxt-extensions
cd dxt-extensions

# 最初の拡張機能用ディレクトリ
mkdir hello-world-extension
cd hello-world-extension
```

### 基本的なディレクトリ構造

```
hello-world-extension/
├── manifest.json          # 拡張機能の設定
├── server/                # MCPサーバーのコード
│   ├── index.js          # メインエントリーポイント
│   └── package.json      # Node.js依存関係
├── assets/               # 追加リソース（オプション）
│   └── icon.png         # 拡張機能のアイコン
└── README.md            # 拡張機能の説明
```

## 🔧 開発環境の最適化

### VS Code の設定（推奨）

以下の拡張機能をインストールすることを推奨します：

```json
{
  "recommendations": [
    "ms-vscode.vscode-json",
    "bradlc.vscode-tailwindcss",
    "ms-vscode.vscode-typescript-next"
  ]
}
```

### デバッグ設定

`.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug MCP Server",
      "type": "node",
      "request": "launch",
      "program": "${workspaceFolder}/server/index.js",
      "console": "integratedTerminal",
      "skipFiles": ["<node_internals>/**"]
    }
  ]
}
```

## 🧪 テスト環境の準備

### MCP SDK のインストール

```bash
cd server
npm init -y
npm install @modelcontextprotocol/sdk
```

### 基本的なテストスクリプト

`server/test.js`:

```javascript
import { spawn } from "child_process";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";

// MCPサーバーをテスト実行
const server = spawn("node", ["index.js"]);

server.stdout.on("data", (data) => {
  console.log(`Server output: ${data}`);
});

server.stderr.on("data", (data) => {
  console.error(`Server error: ${data}`);
});
```

## 🔍 トラブルシューティング

### よくある問題と解決策

#### 1. DXT CLI が見つからない

```bash
# npmのグローバルパスを確認
npm config get prefix

# パスを追加（bash/zshの場合）
echo 'export PATH="$(npm config get prefix)/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

#### 2. Node.js のバージョンが古い

```bash
# nvm（Node Version Manager）を使用
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install 18
nvm use 18
```

#### 3. Claude Desktop で拡張機能が表示されない

- Claude Desktop の再起動
- システムの再起動
- ファイルの権限確認

```bash
# macOSでの権限確認
ls -la *.dxt
chmod 644 your-extension.dxt
```

## ✅ セットアップ確認

以下のコマンドが正常に動作することを確認してください：

```bash
# DXT CLIの動作確認
dxt --version

# Node.jsの動作確認
node -e "console.log('Node.js is working!')"

# NPMパッケージのインストール確認
npm list -g @anthropic-ai/dxt
```

すべて正常に動作したら、次は [初回 DXT 作成](./02-first-extension.md) に進んでください。

## 📚 参考資料

- [Node.js 公式ドキュメント](https://nodejs.org/docs/)
- [npm 公式ドキュメント](https://docs.npmjs.com/)
- [DXT CLI Reference](https://github.com/anthropics/dxt/blob/main/CLI.md)
- [MCP SDK Documentation](https://modelcontextprotocol.io/docs)
