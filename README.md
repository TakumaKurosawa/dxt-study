# DXT (Desktop Extensions) 完全チュートリアル

## 📚 目次

DXT について学習するためのチュートリアルです。以下の順序で読み進めることをお勧めします：

1. [DXT とは何か](#dxtとは何か)
2. [従来の課題と DXT の解決策](#従来の課題とdxtの解決策)
3. [基本セットアップ](./docs/01-setup.md)
4. [初回 DXT 作成](./docs/02-first-extension.md)
5. [実践例：ファイルシステム連携](./docs/03-filesystem-example.md)
6. [実践例：データベース連携](./docs/04-database-example.md)
7. [セキュリティとエンタープライズ機能](./docs/05-security-enterprise.md)
8. [高度な機能とベストプラクティス](./docs/06-advanced-practices.md)

## DXT とは何か

**Desktop Extensions (DXT)** は、Anthropic 社が 2025 年 6 月 27 日に発表した革新的な技術です。Model Context Protocol (MCP) サーバーをワンクリックでインストールできるパッケージング形式として設計されています。

### 🎯 主な特徴

- **ワンクリックインストール**: 技術的知識不要で簡単インストール
- **自己完結型パッケージ**: 全依存関係を内包
- **セキュア**: API キーは OS キーチェーンで管理
- **自動アップデート**: 最新版を自動で取得
- **オープン仕様**: 他の AI アプリでも利用可能

### 🔄 DXT の仕組み

```
.dxt ファイル (ZIPアーカイブ)
├── manifest.json         # 拡張機能のメタデータ
├── server/               # MCPサーバーのコード
│   ├── index.js
│   └── package.json
├── node_modules/         # 依存関係（バンドル済み）
└── assets/               # その他のリソース
```

## 従来の課題と DXT の解決策

### ❌ 従来の MCP サーバー導入の問題

```bash
# 1. ランタイムのインストール
brew install node  # または python

# 2. MCPサーバーのクローン
git clone https://github.com/example/mcp-server

# 3. 依存関係のインストール
npm install

# 4. 設定ファイルの編集
vi ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

**課題:**

- パスの設定ミスで動作しない
- 依存関係の競合
- アップデート時の再設定
- チーム内での環境差異
- セキュリティリスク（平文での API キー保存）

### ✅ DXT による解決

**Before（従来）:**

```bash
# Node.js環境をインストール
npm install -g @example/mcp-server
# 設定ファイルを手動編集
# Claude Desktopを再起動
# 動作するかを祈る
```

**After（DXT）:**

1. `.dxt` ファイルをダウンロード
2. ダブルクリックで Claude Desktop で開く
3. "Install" ボタンをクリック

これだけです！

## 🚀 クイックスタート

```bash
# CLIツールのインストール
npm install -g @anthropic-ai/dxt

# 新規プロジェクトの作成
dxt init my-extension

# DXTファイルの作成
dxt pack
```

## 📖 サンプル例

### 基本的な manifest.json

```json
{
  "dxt_version": "0.1",
  "name": "my-extension",
  "version": "1.0.0",
  "description": "シンプルなMCP拡張機能",
  "author": {
    "name": "Your Name"
  },
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

### MCP サーバーの実装例

```javascript
import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";

const server = new Server(
  {
    name: "Hello World Server",
    version: "1.0.0",
  },
  {
    capabilities: {
      tools: {},
    },
  }
);

// ツールの定義
server.setRequestHandler(ListToolsRequestSchema, async () => ({
  tools: [
    {
      name: "say_hello",
      description: "挨拶をします",
      inputSchema: {
        type: "object",
        properties: {
          name: {
            type: "string",
            description: "挨拶する相手の名前",
          },
        },
        required: ["name"],
      },
    },
  ],
}));

// ツールの実装
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;

  if (name === "say_hello") {
    return {
      content: [
        {
          type: "text",
          text: `こんにちは、${args.name}さん！`,
        },
      ],
    };
  }

  throw new Error(`Unknown tool: ${name}`);
});

// サーバーの起動
const transport = new StdioServerTransport();
await server.connect(transport);
```

## 🌟 公式サンプルの衝撃例

DXT の可能性を示す公式サンプル：

- **Apple Notes 連携**: 「今日の議事録を Notes に保存」
- **Chrome 完全制御**: 全タブの情報を横断的に収集
- **iMessage 連携**: メッセージ履歴の解析と自動返信

## 📈 エコシステムの拡大

DXT 仕様は完全にオープンソースであり、他の AI アプリケーションでも採用可能です：

1. **標準化**: MCP サーバーの配布フォーマットのデファクトスタンダード
2. **相互運用性**: 一度作成した DXT は複数の AI アプリで利用可能
3. **マーケットプレイス**: 将来的には DXT 専用のストアが登場予定

## 🔗 参考リンク

- [Anthropic 公式ブログ](https://www.anthropic.com/engineering/desktop-extensions)
- [DXT GitHub Repository](https://github.com/anthropics/dxt)
- [MCP 仕様](https://modelcontextprotocol.io/)
- [Claude Desktop](https://claude.ai/download)

## 📝 ライセンス

このチュートリアルは MIT ライセンスの下で公開されています。

---

**今すぐ始める理由:**

- 🚀 先行者利益を得られる
- 🔧 社内ツールの価値を最大化
- 🌍 オープンエコシステムへの貢献
- 💡 新しいビジネスチャンスの創出

次は [基本セットアップ](./docs/01-setup.md) に進んでください。
