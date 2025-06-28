# 初回 DXT 作成

## 🎯 目標

このガイドでは、シンプルな「Hello World」DXT 拡張機能を作成し、Claude Desktop にインストールして動作確認を行います。

### 📋 作業の流れ

1. **プロジェクト初期化**: `dxt init` で manifest.json を作成
2. **基本構造作成**: 必要なディレクトリとファイルを手動作成
3. **詳細実装**: MCP サーバーのツール機能を実装
4. **テストと検証**: ローカルでの動作確認
5. **パッケージング**: .dxt ファイルの作成
6. **インストール**: Claude Desktop へのインストールと動作確認

> **注意**: `dxt init` は manifest.json のみを作成します。server/ ディレクトリやその他のファイルは手動で作成する必要があります。

## 📝 ステップ 1: プロジェクトの初期化

### DXT CLI を使用した初期化

```bash
# 新しいディレクトリを作成
mkdir hello-world-dxt
cd hello-world-dxt

# DXTプロジェクトを初期化（manifest.jsonのみ作成）
dxt init
```

対話形式で以下の情報を入力します：

```
? Extension name: hello-world
? Display name: Hello World Extension
? Description: 初回作成用のシンプルなHello World拡張機能
? Author name: Your Name
? Author email: your.email@example.com
? Version: 1.0.0
? License: MIT
? Extension type: node
```

### 追加ディレクトリとファイルの作成

`dxt init` は manifest.json のみを作成するため、残りのファイルを手動で作成します：

```bash
# サーバーディレクトリを作成
mkdir server

# READMEファイルを作成（オプション）
touch README.md
```

### 基本ファイルの作成

#### server/package.json の作成

```bash
cat > server/package.json << 'EOF'
{
  "name": "hello-world-mcp-server",
  "version": "1.0.0",
  "type": "module",
  "main": "index.js",
  "dependencies": {
    "@modelcontextprotocol/sdk": "^1.0.0"
  }
}
EOF
```

#### server/index.js の基本構造作成

```bash
cat > server/index.js << 'EOF'
import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";

// TODO: 詳細な実装はステップ3で行います
const server = new Server({
  name: "Hello World Server",
  version: "1.0.0"
}, {
  capabilities: { tools: {} }
});

console.log("Hello World MCP server template created");
EOF
```

#### 依存関係のインストール

```bash
# serverディレクトリに移動して依存関係をインストール
cd server
npm install
cd ..
```

これで基本的なプロジェクト構造が完成します。詳細な実装は次のステップで行います。

## 📁 ステップ 2: プロジェクト構造の確認

初期化後、以下の構造が作成されます：

```
hello-world-dxt/
├── manifest.json          # 拡張機能のメタデータ
├── server/                # MCPサーバーのコード
│   ├── index.js          # メインサーバーファイル
│   └── package.json      # Node.js依存関係
└── README.md             # 生成された説明書
```

### manifest.json の確認

生成された`manifest.json`を確認しましょう：

```json
{
  "dxt_version": "0.1",
  "name": "hello-world",
  "display_name": "Hello World Extension",
  "version": "1.0.0",
  "description": "初回作成用のシンプルなHello World拡張機能",
  "author": {
    "name": "Your Name",
    "email": "your.email@example.com"
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

⚠️ `mise` を使っていたりして `node` コマンドパスが異なる場合： `server.mcp_confg.env` で `PATH` を通すとうまくいく

```json
{
  "dxt_version": "0.1",
  "name": "hello-world",
  "display_name": "Hello World Extension",
  "version": "1.0.0",
  "description": "初回作成用のシンプルなHello World拡張機能",
  "author": {
    "name": "Your Name",
    "email": "your.email@example.com"
  },
  "license": "MIT",
  "server": {
    "type": "node",
    "entry_point": "server/index.js",
    "mcp_config": {
      "command": "node",
      "args": ["${__dirname}/server/index.js"],
      "env": {
        "PATH": "/path-to-home/.local/share/mise/installs/node/22.17.0/bin:${PATH}"
      }
    }
  }
}
```

## 🔧 ステップ 3: MCP サーバーの実装

### package.json の設定

`server/package.json`:

```json
{
  "name": "hello-world-mcp-server",
  "version": "1.0.0",
  "type": "module",
  "main": "index.js",
  "dependencies": {
    "@modelcontextprotocol/sdk": "^1.0.0"
  }
}
```

### メインサーバーファイルの実装

`server/index.js`:

```javascript
import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  ListToolsRequestSchema,
  CallToolRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";

// MCPサーバーの作成
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

// 利用可能なツール一覧を返すハンドラー
server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
        name: "say_hello",
        description: "指定された名前に対して挨拶を返します",
        inputSchema: {
          type: "object",
          properties: {
            name: {
              type: "string",
              description: "挨拶する相手の名前",
            },
            language: {
              type: "string",
              description: "挨拶の言語 (japanese, english)",
              enum: ["japanese", "english"],
              default: "japanese",
            },
          },
          required: ["name"],
        },
      },
      {
        name: "get_time",
        description: "現在の時刻を返します",
        inputSchema: {
          type: "object",
          properties: {
            format: {
              type: "string",
              description: "時刻のフォーマット (12h, 24h)",
              enum: ["12h", "24h"],
              default: "24h",
            },
          },
        },
      },
    ],
  };
});

// ツール実行のハンドラー
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;

  switch (name) {
    case "say_hello":
      const language = args.language || "japanese";
      const greeting =
        language === "english"
          ? `Hello, ${args.name}!`
          : `こんにちは、${args.name}さん！`;

      return {
        content: [
          {
            type: "text",
            text: greeting,
          },
        ],
      };

    case "get_time":
      const now = new Date();
      const format = args.format || "24h";

      let timeString;
      if (format === "12h") {
        timeString = now.toLocaleTimeString("ja-JP", {
          hour12: true,
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
        });
      } else {
        timeString = now.toLocaleTimeString("ja-JP", {
          hour12: false,
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
        });
      }

      return {
        content: [
          {
            type: "text",
            text: `現在の時刻: ${timeString}`,
          },
        ],
      };

    default:
      throw new Error(`Unknown tool: ${name}`);
  }
});

// エラーハンドリング
server.onerror = (error) => {
  console.error("[MCP Error]", error);
};

// サーバーの起動
async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.log("Hello World MCP server running on stdio");
}

main().catch((error) => {
  console.error("Failed to start server:", error);
  process.exit(1);
});
```

## 📦 ステップ 4: 依存関係のインストール

```bash
# serverディレクトリに移動
cd server

# 依存関係をインストール
npm install

# プロジェクトルートに戻る
cd ..
```

## 🧪 ステップ 5: ローカルテスト

### 基本動作テスト

```bash
# サーバーディレクトリでテスト実行
cd server
node index.js
```

正常に起動すると以下のメッセージが表示されます：

```
Hello World MCP server running on stdio
```

`Ctrl+C`で停止してください。

### マニフェストの検証

```bash
# プロジェクトルートで実行
dxt validate

# 期待される出力
✓ manifest.json is valid
✓ All required files exist
✓ Server configuration is correct
```

## 📦 ステップ 6: DXT ファイルの作成

```bash
# プロジェクトルートで実行
dxt pack

# 成功すると以下のファイルが生成される
# hello-world-1.0.0.dxt
```

### パッケージ内容の確認

```bash
# .dxtファイルはZIPアーカイブなので確認可能
unzip -l hello-world-1.0.0.dxt
```

期待される出力:

```
Archive:  hello-world-1.0.0.dxt
  Length      Date    Time    Name
---------  ---------- -----   ----
     1234  01-01-2025 12:00   manifest.json
     2345  01-01-2025 12:00   server/index.js
     3456  01-01-2025 12:00   server/package.json
     4567  01-01-2025 12:00   server/node_modules/...
---------                     -------
```

## 🚀 ステップ 7: Claude Desktop へのインストール

### 1. Claude Desktop を開く

Claude Desktop アプリケーションを起動します。

### 2. 設定画面を開く

- macOS: `Cmd + ,` または メニューから "Claude Desktop" > "Settings"
- Windows: `Ctrl + ,` または メニューから "File" > "Settings"

### 3. Extensions セクションへ移動

設定画面の「Extensions」タブをクリックします。

### 4. 拡張機能のインストール

1. "Install Extension" ボタンをクリック
2. 作成した `hello-world-1.0.0.dxt` ファイルを選択
3. インストール確認ダイアログで "Install" をクリック

### 5. インストール確認

Extensions 画面に「Hello World Extension」が表示され、ステータスが「Active」になることを確認します。

## ✅ ステップ 8: 動作確認

### 1. 新しいチャットを開始

Claude Desktop で新しいチャットセッションを開始します。

### 2. ツールの動作確認

以下のメッセージを送信してツールが正常に動作するか確認します：

```
私の名前は太郎です。挨拶してください。
```

期待される応答:

```
こんにちは、太郎さん！
```

### 3. 追加テスト

```
現在の時刻を教えてください
```

期待される応答:

```
現在の時刻: 14:30:25
```

### 4. 英語での挨拶テスト

```
Please say hello to Alice in English.
```

期待される応答:

```
Hello, Alice!
```

## 🔧 トラブルシューティング

### よくある問題と解決策

#### 1. 拡張機能が Extensions 画面に表示されない

- Claude Desktop を再起動
- .dxt ファイルが正しい場所に保存されているか確認
- マニフェストファイルの構文エラーをチェック

```bash
dxt validate
```

#### 2. ツールが実行されない

- サーバーのログを確認（Claude Desktop の開発者ツールで確認可能）
- MCP SDK のバージョンが最新かチェック

```bash
cd server
npm update @modelcontextprotocol/sdk
```

#### 3. パッケージングエラー

- 必要なファイルが全て存在するか確認
- `node_modules`が正しくインストールされているか確認

```bash
cd server
rm -rf node_modules package-lock.json
npm install
cd ..
dxt pack
```

## 🎉 次のステップ

おめでとうございます！初めての DXT 拡張機能が完成しました。

次は以下のチュートリアルに進んでください：

- [ファイルシステム連携](./03-filesystem-example.md) - ローカルファイルとの連携
- [データベース連携](./04-database-example.md) - データベースとの連携
- [セキュリティとエンタープライズ機能](./05-security-enterprise.md) - 本格運用向けの設定

## 📚 参考資料

- [MCP SDK API Reference](https://modelcontextprotocol.io/docs/api/)
- [DXT Manifest Specification](https://github.com/anthropics/dxt/blob/main/MANIFEST.md)
- [Claude Desktop Extensions Guide](https://docs.anthropic.com/claude/desktop-extensions)
