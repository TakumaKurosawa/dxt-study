# 実践例：ファイルシステム連携

## 🎯 目標

このガイドでは、ローカルファイルシステムと連携する DXT 拡張機能を作成します。ファイルの読み取り、書き込み、検索、統計情報の取得などの機能を実装します。

## 📋 機能概要

作成する拡張機能の機能：

- 📁 ディレクトリ内容の一覧表示
- 📄 ファイル内容の読み取り
- ✏️ ファイルの作成・編集
- 🔍 ファイル内容の検索
- 📊 ファイル・ディレクトリの統計情報
- 🗑️ ファイル・ディレクトリの削除

## 🚀 ステップ 1: プロジェクトの初期化

```bash
# 新しいプロジェクトディレクトリを作成
mkdir filesystem-manager-dxt
cd filesystem-manager-dxt

# DXTプロジェクトを初期化
dxt init
```

対話形式の入力：

```
? Extension name: filesystem-manager
? Display name: File System Manager
? Description: ローカルファイルシステムとの包括的な連携を提供
? Author name: Your Name
? Author email: your.email@example.com
? Version: 1.0.0
? License: MIT
? Extension type: node
```

## 📝 ステップ 2: manifest.json の設定

セキュリティ設定とユーザー設定を追加します：

```json
{
  "dxt_version": "0.1",
  "name": "filesystem-manager",
  "display_name": "File System Manager",
  "version": "1.0.0",
  "description": "ローカルファイルシステムとの包括的な連携を提供",
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
        "ALLOWED_DIRECTORIES": "${user_config.allowed_directories}",
        "MAX_FILE_SIZE": "${user_config.max_file_size}",
        "SAFE_MODE": "${user_config.safe_mode}"
      }
    }
  },
  "user_config": {
    "allowed_directories": {
      "type": "string",
      "title": "許可ディレクトリ",
      "description": "アクセスを許可するディレクトリのパス（カンマ区切り）",
      "default": "/Users/YOUR_USERNAME/Documents,/Users/YOUR_USERNAME/Desktop"
    },
    "max_file_size": {
      "type": "string",
      "title": "最大ファイルサイズ",
      "description": "読み取り可能な最大ファイルサイズ（MB）",
      "default": "10"
    },
    "safe_mode": {
      "type": "string",
      "title": "セーフモード",
      "description": "削除操作を無効にする（true/false）",
      "default": "true"
    }
  },
  "tools": [
    {
      "name": "list_directory",
      "description": "ディレクトリの内容を一覧表示"
    },
    {
      "name": "read_file",
      "description": "ファイルの内容を読み取り"
    },
    {
      "name": "write_file",
      "description": "ファイルの作成または編集"
    },
    {
      "name": "search_files",
      "description": "ファイル内容やファイル名を検索"
    },
    {
      "name": "get_file_stats",
      "description": "ファイル・ディレクトリの統計情報を取得"
    },
    {
      "name": "delete_file",
      "description": "ファイルまたはディレクトリを削除"
    }
  ]
}
```

## 🔧 ステップ 3: package.json の設定

`server/package.json`:

```json
{
  "name": "filesystem-manager-mcp-server",
  "version": "1.0.0",
  "type": "module",
  "main": "index.js",
  "dependencies": {
    "@modelcontextprotocol/sdk": "^1.0.0"
  }
}
```

## 💻 ステップ 4: サーバー実装

### メインサーバーファイル

`server/index.js`:

```javascript
import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  ListToolsRequestSchema,
  CallToolRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";
import fs from "fs/promises";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// 設定の読み込み
const ALLOWED_DIRECTORIES = process.env.ALLOWED_DIRECTORIES?.split(",") || [];
const MAX_FILE_SIZE = parseInt(process.env.MAX_FILE_SIZE || "10") * 1024 * 1024; // MB to bytes
const SAFE_MODE = process.env.SAFE_MODE === "true";

// MCPサーバーの作成
const server = new Server(
  {
    name: "File System Manager Server",
    version: "1.0.0",
  },
  {
    capabilities: {
      tools: {},
    },
  }
);

// セキュリティ: パスの検証
function isPathAllowed(targetPath) {
  const normalizedPath = path.resolve(targetPath);
  return ALLOWED_DIRECTORIES.some((allowedDir) =>
    normalizedPath.startsWith(path.resolve(allowedDir))
  );
}

// ファイルサイズの検証
async function isFileSizeAllowed(filePath) {
  try {
    const stats = await fs.stat(filePath);
    return stats.size <= MAX_FILE_SIZE;
  } catch {
    return false;
  }
}

// ツール一覧の定義
server.setRequestHandler(ListToolsRequestSchema, async () => {
  return {
    tools: [
      {
        name: "list_directory",
        description: "ディレクトリの内容を一覧表示します",
        inputSchema: {
          type: "object",
          properties: {
            directory_path: {
              type: "string",
              description: "リストするディレクトリのパス",
            },
            show_hidden: {
              type: "boolean",
              description: "隠しファイルを表示するか",
              default: false,
            },
            detailed: {
              type: "boolean",
              description: "詳細情報を表示するか",
              default: false,
            },
          },
          required: ["directory_path"],
        },
      },
      {
        name: "read_file",
        description: "ファイルの内容を読み取ります",
        inputSchema: {
          type: "object",
          properties: {
            file_path: {
              type: "string",
              description: "読み取るファイルのパス",
            },
            encoding: {
              type: "string",
              description: "ファイルのエンコーディング",
              enum: ["utf8", "ascii", "base64"],
              default: "utf8",
            },
          },
          required: ["file_path"],
        },
      },
      {
        name: "write_file",
        description: "ファイルの作成または編集を行います",
        inputSchema: {
          type: "object",
          properties: {
            file_path: {
              type: "string",
              description: "書き込むファイルのパス",
            },
            content: {
              type: "string",
              description: "ファイルに書き込む内容",
            },
            encoding: {
              type: "string",
              description: "ファイルのエンコーディング",
              enum: ["utf8", "ascii"],
              default: "utf8",
            },
            create_directories: {
              type: "boolean",
              description: "必要に応じてディレクトリを作成するか",
              default: true,
            },
          },
          required: ["file_path", "content"],
        },
      },
      {
        name: "search_files",
        description: "ファイル内容やファイル名を検索します",
        inputSchema: {
          type: "object",
          properties: {
            search_directory: {
              type: "string",
              description: "検索するディレクトリのパス",
            },
            query: {
              type: "string",
              description: "検索クエリ",
            },
            search_type: {
              type: "string",
              description: "検索タイプ",
              enum: ["filename", "content", "both"],
              default: "both",
            },
            file_extensions: {
              type: "string",
              description: "検索対象のファイル拡張子（カンマ区切り）",
              default: ".txt,.md,.js,.py,.json",
            },
            case_sensitive: {
              type: "boolean",
              description: "大文字小文字を区別するか",
              default: false,
            },
          },
          required: ["search_directory", "query"],
        },
      },
      {
        name: "get_file_stats",
        description: "ファイル・ディレクトリの統計情報を取得します",
        inputSchema: {
          type: "object",
          properties: {
            target_path: {
              type: "string",
              description: "統計情報を取得するパス",
            },
          },
          required: ["target_path"],
        },
      },
      {
        name: "delete_file",
        description: "ファイルまたはディレクトリを削除します",
        inputSchema: {
          type: "object",
          properties: {
            target_path: {
              type: "string",
              description: "削除するファイル・ディレクトリのパス",
            },
            recursive: {
              type: "boolean",
              description: "ディレクトリを再帰的に削除するか",
              default: false,
            },
          },
          required: ["target_path"],
        },
      },
    ],
  };
});

// ツール実行のハンドラー
server.setRequestHandler(CallToolRequestSchema, async (request) => {
  const { name, arguments: args } = request.params;

  try {
    switch (name) {
      case "list_directory":
        return await handleListDirectory(args);
      case "read_file":
        return await handleReadFile(args);
      case "write_file":
        return await handleWriteFile(args);
      case "search_files":
        return await handleSearchFiles(args);
      case "get_file_stats":
        return await handleGetFileStats(args);
      case "delete_file":
        return await handleDeleteFile(args);
      default:
        throw new Error(`Unknown tool: ${name}`);
    }
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
});

// ディレクトリ一覧の処理
async function handleListDirectory(args) {
  const { directory_path, show_hidden = false, detailed = false } = args;

  if (!isPathAllowed(directory_path)) {
    throw new Error("指定されたディレクトリへのアクセスが許可されていません");
  }

  const items = await fs.readdir(directory_path);
  const filteredItems = show_hidden
    ? items
    : items.filter((item) => !item.startsWith("."));

  if (!detailed) {
    return {
      content: [
        {
          type: "text",
          text: `ディレクトリ: ${directory_path}\n\n${filteredItems.join(
            "\n"
          )}`,
        },
      ],
    };
  }

  // 詳細情報を取得
  const detailedItems = await Promise.all(
    filteredItems.map(async (item) => {
      const itemPath = path.join(directory_path, item);
      try {
        const stats = await fs.stat(itemPath);
        return {
          name: item,
          type: stats.isDirectory() ? "directory" : "file",
          size: stats.size,
          modified: stats.mtime.toISOString(),
          permissions: stats.mode.toString(8),
        };
      } catch {
        return { name: item, type: "unknown", error: "アクセスできません" };
      }
    })
  );

  const formatSize = (bytes) => {
    const units = ["B", "KB", "MB", "GB"];
    let size = bytes;
    let unitIndex = 0;
    while (size >= 1024 && unitIndex < units.length - 1) {
      size /= 1024;
      unitIndex++;
    }
    return `${size.toFixed(1)}${units[unitIndex]}`;
  };

  const output = detailedItems
    .map((item) => {
      if (item.error) {
        return `${item.name} (${item.error})`;
      }
      const sizeStr =
        item.type === "directory" ? "<DIR>" : formatSize(item.size);
      return `${item.type === "directory" ? "📁" : "📄"} ${item.name.padEnd(
        30
      )} ${sizeStr.padStart(10)} ${item.modified.substring(0, 19)}`;
    })
    .join("\n");

  return {
    content: [
      {
        type: "text",
        text: `ディレクトリ: ${directory_path}\n\n${output}`,
      },
    ],
  };
}

// ファイル読み取りの処理
async function handleReadFile(args) {
  const { file_path, encoding = "utf8" } = args;

  if (!isPathAllowed(file_path)) {
    throw new Error("指定されたファイルへのアクセスが許可されていません");
  }

  if (!(await isFileSizeAllowed(file_path))) {
    throw new Error(
      `ファイルサイズが制限（${MAX_FILE_SIZE / 1024 / 1024}MB）を超えています`
    );
  }

  const content = await fs.readFile(file_path, encoding);

  return {
    content: [
      {
        type: "text",
        text: `ファイル: ${file_path}\n\n${content}`,
      },
    ],
  };
}

// ファイル書き込みの処理
async function handleWriteFile(args) {
  const {
    file_path,
    content,
    encoding = "utf8",
    create_directories = true,
  } = args;

  if (!isPathAllowed(file_path)) {
    throw new Error("指定されたファイルへのアクセスが許可されていません");
  }

  if (create_directories) {
    const dir = path.dirname(file_path);
    await fs.mkdir(dir, { recursive: true });
  }

  await fs.writeFile(file_path, content, encoding);

  return {
    content: [
      {
        type: "text",
        text: `ファイル "${file_path}" に正常に書き込みました`,
      },
    ],
  };
}

// ファイル検索の処理
async function handleSearchFiles(args) {
  const {
    search_directory,
    query,
    search_type = "both",
    file_extensions = ".txt,.md,.js,.py,.json",
    case_sensitive = false,
  } = args;

  if (!isPathAllowed(search_directory)) {
    throw new Error("指定されたディレクトリへのアクセスが許可されていません");
  }

  const extensions = file_extensions.split(",").map((ext) => ext.trim());
  const searchQuery = case_sensitive ? query : query.toLowerCase();
  const results = [];

  async function searchInDirectory(dir) {
    const items = await fs.readdir(dir);

    for (const item of items) {
      const itemPath = path.join(dir, item);

      try {
        const stats = await fs.stat(itemPath);

        if (stats.isDirectory()) {
          await searchInDirectory(itemPath);
        } else {
          const ext = path.extname(item).toLowerCase();
          if (!extensions.includes(ext)) continue;

          // ファイル名検索
          if (search_type === "filename" || search_type === "both") {
            const filename = case_sensitive ? item : item.toLowerCase();
            if (filename.includes(searchQuery)) {
              results.push({
                type: "filename",
                path: itemPath,
                match: item,
              });
            }
          }

          // ファイル内容検索
          if (search_type === "content" || search_type === "both") {
            if (await isFileSizeAllowed(itemPath)) {
              const content = await fs.readFile(itemPath, "utf8");
              const searchContent = case_sensitive
                ? content
                : content.toLowerCase();

              if (searchContent.includes(searchQuery)) {
                // マッチした行を取得
                const lines = content.split("\n");
                const matchedLines = lines
                  .map((line, index) => ({ line, number: index + 1 }))
                  .filter(({ line }) => {
                    const checkLine = case_sensitive
                      ? line
                      : line.toLowerCase();
                    return checkLine.includes(searchQuery);
                  })
                  .slice(0, 3); // 最初の3行のみ

                results.push({
                  type: "content",
                  path: itemPath,
                  matches: matchedLines,
                });
              }
            }
          }
        }
      } catch {
        // アクセスできないファイルはスキップ
      }
    }
  }

  await searchInDirectory(search_directory);

  if (results.length === 0) {
    return {
      content: [
        {
          type: "text",
          text: `"${query}" に一致する結果が見つかりませんでした`,
        },
      ],
    };
  }

  const output = results
    .slice(0, 20) // 最大20件
    .map((result) => {
      if (result.type === "filename") {
        return `📄 ${result.path}\n   ファイル名にマッチ: "${result.match}"`;
      } else {
        const matchLines = result.matches
          .map((m) => `   ${m.number}: ${m.line.trim()}`)
          .join("\n");
        return `📄 ${result.path}\n${matchLines}`;
      }
    })
    .join("\n\n");

  return {
    content: [
      {
        type: "text",
        text: `検索結果 (${results.length}件中最大20件表示):\n\n${output}`,
      },
    ],
  };
}

// ファイル統計情報の処理
async function handleGetFileStats(args) {
  const { target_path } = args;

  if (!isPathAllowed(target_path)) {
    throw new Error("指定されたパスへのアクセスが許可されていません");
  }

  const stats = await fs.stat(target_path);

  const formatSize = (bytes) => {
    const units = ["B", "KB", "MB", "GB"];
    let size = bytes;
    let unitIndex = 0;
    while (size >= 1024 && unitIndex < units.length - 1) {
      size /= 1024;
      unitIndex++;
    }
    return `${size.toFixed(2)} ${units[unitIndex]}`;
  };

  const info = {
    path: target_path,
    type: stats.isDirectory() ? "ディレクトリ" : "ファイル",
    size: formatSize(stats.size),
    created: stats.birthtime.toLocaleString("ja-JP"),
    modified: stats.mtime.toLocaleString("ja-JP"),
    accessed: stats.atime.toLocaleString("ja-JP"),
    permissions: stats.mode.toString(8),
  };

  const output = Object.entries(info)
    .map(([key, value]) => `${key}: ${value}`)
    .join("\n");

  return {
    content: [
      {
        type: "text",
        text: output,
      },
    ],
  };
}

// ファイル削除の処理
async function handleDeleteFile(args) {
  if (SAFE_MODE) {
    throw new Error("セーフモードが有効なため、削除操作は無効化されています");
  }

  const { target_path, recursive = false } = args;

  if (!isPathAllowed(target_path)) {
    throw new Error("指定されたパスへのアクセスが許可されていません");
  }

  const stats = await fs.stat(target_path);

  if (stats.isDirectory()) {
    if (recursive) {
      await fs.rm(target_path, { recursive: true });
      return {
        content: [
          {
            type: "text",
            text: `ディレクトリ "${target_path}" を再帰的に削除しました`,
          },
        ],
      };
    } else {
      await fs.rmdir(target_path);
      return {
        content: [
          {
            type: "text",
            text: `ディレクトリ "${target_path}" を削除しました`,
          },
        ],
      };
    }
  } else {
    await fs.unlink(target_path);
    return {
      content: [
        {
          type: "text",
          text: `ファイル "${target_path}" を削除しました`,
        },
      ],
    };
  }
}

// エラーハンドリング
server.onerror = (error) => {
  console.error("[MCP Error]", error);
};

// サーバーの起動
async function main() {
  console.log("File System Manager MCP Server");
  console.log("許可されたディレクトリ:", ALLOWED_DIRECTORIES);
  console.log("最大ファイルサイズ:", `${MAX_FILE_SIZE / 1024 / 1024}MB`);
  console.log("セーフモード:", SAFE_MODE ? "有効" : "無効");

  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.log("Server running on stdio");
}

main().catch((error) => {
  console.error("Failed to start server:", error);
  process.exit(1);
});
```

## 📦 ステップ 5: ビルドとテスト

### 依存関係のインストール

```bash
cd server
npm install
cd ..
```

### ローカルテスト

```bash
cd server
node index.js
```

### DXT ファイルの作成

```bash
dxt validate
dxt pack
```

## 🔧 ステップ 6: セキュリティ設定

### インストール時のユーザー設定

DXT をインストールする際、以下の設定が求められます：

1. **許可ディレクトリ**: アクセスを許可するディレクトリのパス
2. **最大ファイルサイズ**: 読み取り可能な最大ファイルサイズ
3. **セーフモード**: 削除操作の有効/無効

### 推奨設定例

```
許可ディレクトリ: /Users/USERNAME/Documents,/Users/USERNAME/Desktop,/Users/USERNAME/Projects
最大ファイルサイズ: 5
セーフモード: true
```

## 🧪 ステップ 7: 動作確認

### 1. ディレクトリ一覧表示

```
Documentsフォルダの内容を詳細情報付きで表示してください
```

### 2. ファイル読み取り

```
Desktop/sample.txtファイルの内容を読み取ってください
```

### 3. ファイル作成

```
Desktop/test.txtファイルに「Hello, DXT!」と書き込んでください
```

### 4. ファイル検索

```
Documentsフォルダ内で「JavaScript」という単語を含むファイルを検索してください
```

### 5. ファイル統計

```
Documentsフォルダの統計情報を取得してください
```

## 🔍 トラブルシューティング

### パーミッションエラー

```bash
# macOSでのフルディスクアクセス許可
# システム環境設定 > セキュリティとプライバシー > プライバシー > フルディスクアクセス
# Claude Desktopを追加
```

### ファイルが読み取れない

- ファイルサイズが制限を超えていないか確認
- ファイルのエンコーディングが正しいか確認
- 許可されたディレクトリ内にあるか確認

## 🎯 応用例

この拡張機能を基にして、以下のような機能追加が可能です：

- 📄 **文書解析**: PDF、Word 文書の内容解析
- 🖼️ **画像処理**: 画像メタデータの取得、リサイズ
- 🗜️ **アーカイブ操作**: ZIP、TAR.GZ ファイルの展開
- 🔄 **同期機能**: 複数ディレクトリ間のファイル同期

次は [データベース連携](./04-database-example.md) に進んで、さらに高度な機能を学習しましょう。
