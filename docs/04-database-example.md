# 実践例：データベース連携

## 🎯 目標

このガイドでは、データベースと連携する高度な DXT 拡張機能を作成します。SQLite と PostgreSQL の両方に対応し、セキュアな CRUD 操作を実装します。

## 💡 学習内容

- データベース接続の管理
- セキュアな SQL クエリの実装
- CRUD 操作の完全実装
- 設定可能なデータベース選択
- エラーハンドリングとロギング
- パフォーマンス最適化

## 📁 プロジェクト構造

```
database-connector/
├── manifest.json
├── server/
│   ├── package.json
│   ├── index.js
│   ├── database/
│   │   ├── sqlite.js
│   │   ├── postgresql.js
│   │   └── base.js
│   └── utils/
│       ├── validation.js
│       └── logger.js
└── assets/
    └── schema.sql
```

## 🔧 manifest.json の設定

```json
{
  "dxt_version": "0.1",
  "name": "database-connector",
  "version": "2.0.0",
  "description": "セキュアなデータベース連携拡張機能",
  "author": {
    "name": "DXT Tutorial",
    "email": "tutorial@example.com"
  },
  "server": {
    "type": "node",
    "entry_point": "server/index.js",
    "mcp_config": {
      "command": "node",
      "args": ["${__dirname}/server/index.js"]
    }
  },
  "user_settings": [
    {
      "name": "database_type",
      "type": "choice",
      "choices": ["sqlite", "postgresql"],
      "default": "sqlite",
      "description": "使用するデータベースの種類"
    },
    {
      "name": "database_path",
      "type": "string",
      "default": "./data/app.db",
      "description": "SQLite データベースファイルのパス"
    },
    {
      "name": "postgresql_url",
      "type": "string",
      "secure": true,
      "description": "PostgreSQL 接続 URL (postgresql://user:pass@host:port/db)"
    },
    {
      "name": "max_results",
      "type": "number",
      "default": 100,
      "description": "クエリ結果の最大件数"
    },
    {
      "name": "enable_logging",
      "type": "boolean",
      "default": true,
      "description": "詳細ログを有効にする"
    }
  ],
  "permissions": ["filesystem", "network"]
}
```

## 📦 package.json の設定

`server/package.json`:

```json
{
  "name": "database-connector-server",
  "version": "2.0.0",
  "type": "module",
  "dependencies": {
    "@modelcontextprotocol/sdk": "^1.0.0",
    "sqlite3": "^5.1.6",
    "pg": "^8.11.3",
    "joi": "^17.11.0"
  },
  "scripts": {
    "start": "node index.js"
  }
}
```

## 🗄️ データベース基底クラス

`server/database/base.js`:

```javascript
export class DatabaseBase {
  constructor(config) {
    this.config = config;
    this.connection = null;
  }

  async connect() {
    throw new Error("connect() must be implemented");
  }

  async disconnect() {
    throw new Error("disconnect() must be implemented");
  }

  async execute(query, params = []) {
    throw new Error("execute() must be implemented");
  }

  async select(query, params = []) {
    throw new Error("select() must be implemented");
  }

  escapeIdentifier(identifier) {
    // SQL識別子をエスケープ
    return identifier.replace(/[^a-zA-Z0-9_]/g, "");
  }

  validateTableName(tableName) {
    if (!/^[a-zA-Z][a-zA-Z0-9_]*$/.test(tableName)) {
      throw new Error("Invalid table name");
    }
    return tableName;
  }
}
```

## 🗃️ SQLite 実装

`server/database/sqlite.js`:

```javascript
import sqlite3 from "sqlite3";
import { DatabaseBase } from "./base.js";
import fs from "fs/promises";
import path from "path";

export class SQLiteDatabase extends DatabaseBase {
  constructor(config) {
    super(config);
    this.dbPath = config.database_path;
  }

  async connect() {
    try {
      // ディレクトリが存在しない場合は作成
      const dir = path.dirname(this.dbPath);
      await fs.mkdir(dir, { recursive: true });

      return new Promise((resolve, reject) => {
        this.connection = new sqlite3.Database(this.dbPath, (err) => {
          if (err) reject(err);
          else resolve();
        });
      });
    } catch (error) {
      throw new Error(`SQLite connection failed: ${error.message}`);
    }
  }

  async disconnect() {
    if (this.connection) {
      return new Promise((resolve) => {
        this.connection.close(resolve);
      });
    }
  }

  async execute(query, params = []) {
    return new Promise((resolve, reject) => {
      this.connection.run(query, params, function (err) {
        if (err) reject(err);
        else
          resolve({
            changes: this.changes,
            lastID: this.lastID,
          });
      });
    });
  }

  async select(query, params = []) {
    return new Promise((resolve, reject) => {
      this.connection.all(query, params, (err, rows) => {
        if (err) reject(err);
        else resolve(rows);
      });
    });
  }

  async getTableInfo(tableName) {
    const validTableName = this.validateTableName(tableName);
    const query = `PRAGMA table_info(${validTableName})`;
    return await this.select(query);
  }

  async getTables() {
    const query = `
      SELECT name FROM sqlite_master 
      WHERE type='table' AND name NOT LIKE 'sqlite_%'
      ORDER BY name
    `;
    return await this.select(query);
  }
}
```

## 🐘 PostgreSQL 実装

`server/database/postgresql.js`:

```javascript
import pg from "pg";
import { DatabaseBase } from "./base.js";

const { Client } = pg;

export class PostgreSQLDatabase extends DatabaseBase {
  constructor(config) {
    super(config);
    this.connectionString = config.postgresql_url;
  }

  async connect() {
    try {
      this.connection = new Client({
        connectionString: this.connectionString,
        ssl: this.connectionString.includes("sslmode=require")
          ? { rejectUnauthorized: false }
          : false,
      });
      await this.connection.connect();
    } catch (error) {
      throw new Error(`PostgreSQL connection failed: ${error.message}`);
    }
  }

  async disconnect() {
    if (this.connection) {
      await this.connection.end();
    }
  }

  async execute(query, params = []) {
    try {
      const result = await this.connection.query(query, params);
      return {
        changes: result.rowCount,
        rows: result.rows,
      };
    } catch (error) {
      throw new Error(`Query execution failed: ${error.message}`);
    }
  }

  async select(query, params = []) {
    try {
      const result = await this.connection.query(query, params);
      return result.rows;
    } catch (error) {
      throw new Error(`Query execution failed: ${error.message}`);
    }
  }

  async getTableInfo(tableName) {
    const validTableName = this.validateTableName(tableName);
    const query = `
      SELECT 
        column_name, 
        data_type, 
        is_nullable,
        column_default
      FROM information_schema.columns 
      WHERE table_name = $1
      ORDER BY ordinal_position
    `;
    return await this.select(query, [validTableName]);
  }

  async getTables() {
    const query = `
      SELECT table_name as name
      FROM information_schema.tables 
      WHERE table_schema = 'public'
      ORDER BY table_name
    `;
    return await this.select(query);
  }
}
```

## 🔍 バリデーション機能

`server/utils/validation.js`:

```javascript
import Joi from "joi";

export class Validator {
  static schemas = {
    tableName: Joi.string()
      .pattern(/^[a-zA-Z][a-zA-Z0-9_]*$/)
      .required(),
    columnName: Joi.string()
      .pattern(/^[a-zA-Z][a-zA-Z0-9_]*$/)
      .required(),
    limit: Joi.number().integer().min(1).max(1000).default(100),
    offset: Joi.number().integer().min(0).default(0),
    whereClause: Joi.string().max(500).optional(),
    orderBy: Joi.string()
      .pattern(/^[a-zA-Z][a-zA-Z0-9_]*( (ASC|DESC))?$/)
      .optional(),
  };

  static validateTableName(tableName) {
    const { error, value } = this.schemas.tableName.validate(tableName);
    if (error) throw new Error(`Invalid table name: ${error.message}`);
    return value;
  }

  static validateQueryParams(params) {
    const schema = Joi.object({
      table: this.schemas.tableName,
      columns: Joi.array().items(this.schemas.columnName).optional(),
      where: this.schemas.whereClause,
      orderBy: this.schemas.orderBy,
      limit: this.schemas.limit,
      offset: this.schemas.offset,
    });

    const { error, value } = schema.validate(params);
    if (error) throw new Error(`Invalid parameters: ${error.message}`);
    return value;
  }

  static validateInsertData(data) {
    if (!data || typeof data !== "object") {
      throw new Error("Insert data must be an object");
    }

    const validatedData = {};
    for (const [key, value] of Object.entries(data)) {
      // カラム名の検証
      const { error } = this.schemas.columnName.validate(key);
      if (error) throw new Error(`Invalid column name: ${key}`);

      validatedData[key] = value;
    }

    return validatedData;
  }
}
```

## 📝 ロガー機能

`server/utils/logger.js`:

```javascript
export class Logger {
  constructor(enableLogging = true) {
    this.enableLogging = enableLogging;
  }

  log(level, message, meta = {}) {
    if (!this.enableLogging) return;

    const timestamp = new Date().toISOString();
    const logEntry = {
      timestamp,
      level,
      message,
      ...meta,
    };

    console.log(JSON.stringify(logEntry));
  }

  info(message, meta = {}) {
    this.log("info", message, meta);
  }

  error(message, meta = {}) {
    this.log("error", message, meta);
  }

  warn(message, meta = {}) {
    this.log("warn", message, meta);
  }

  debug(message, meta = {}) {
    this.log("debug", message, meta);
  }
}
```

## 🖥️ メインサーバー実装

`server/index.js`:

```javascript
import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  ListToolsRequestSchema,
  CallToolRequestSchema,
} from "@modelcontextprotocol/sdk/types.js";

import { SQLiteDatabase } from "./database/sqlite.js";
import { PostgreSQLDatabase } from "./database/postgresql.js";
import { Validator } from "./utils/validation.js";
import { Logger } from "./utils/logger.js";

class DatabaseConnector {
  constructor() {
    this.db = null;
    this.logger = new Logger(process.env.ENABLE_LOGGING === "true");
    this.config = this.loadConfig();
    this.server = this.createServer();
  }

  loadConfig() {
    return {
      database_type: process.env.DATABASE_TYPE || "sqlite",
      database_path: process.env.DATABASE_PATH || "./data/app.db",
      postgresql_url: process.env.POSTGRESQL_URL,
      max_results: parseInt(process.env.MAX_RESULTS) || 100,
      enable_logging: process.env.ENABLE_LOGGING === "true",
    };
  }

  async initializeDatabase() {
    try {
      if (this.config.database_type === "postgresql") {
        if (!this.config.postgresql_url) {
          throw new Error("PostgreSQL URL is required");
        }
        this.db = new PostgreSQLDatabase(this.config);
      } else {
        this.db = new SQLiteDatabase(this.config);
      }

      await this.db.connect();
      this.logger.info("Database connected", {
        type: this.config.database_type,
      });
    } catch (error) {
      this.logger.error("Database initialization failed", {
        error: error.message,
      });
      throw error;
    }
  }

  createServer() {
    const server = new Server(
      {
        name: "Database Connector",
        version: "2.0.0",
      },
      {
        capabilities: {
          tools: {},
        },
      }
    );

    this.setupToolHandlers(server);
    return server;
  }

  setupToolHandlers(server) {
    // ツール一覧の定義
    server.setRequestHandler(ListToolsRequestSchema, async () => ({
      tools: [
        {
          name: "db_list_tables",
          description: "データベース内のテーブル一覧を取得",
          inputSchema: {
            type: "object",
            properties: {},
          },
        },
        {
          name: "db_describe_table",
          description: "テーブルの構造を取得",
          inputSchema: {
            type: "object",
            properties: {
              table: {
                type: "string",
                description: "テーブル名",
              },
            },
            required: ["table"],
          },
        },
        {
          name: "db_select",
          description: "SELECT クエリの実行",
          inputSchema: {
            type: "object",
            properties: {
              table: {
                type: "string",
                description: "テーブル名",
              },
              columns: {
                type: "array",
                items: { type: "string" },
                description: "取得するカラム（省略時は全カラム）",
              },
              where: {
                type: "string",
                description: "WHERE 条件",
              },
              orderBy: {
                type: "string",
                description: "ORDER BY 指定",
              },
              limit: {
                type: "number",
                description: "取得件数制限",
                minimum: 1,
                maximum: 1000,
              },
              offset: {
                type: "number",
                description: "オフセット",
                minimum: 0,
              },
            },
            required: ["table"],
          },
        },
        {
          name: "db_insert",
          description: "データの挿入",
          inputSchema: {
            type: "object",
            properties: {
              table: {
                type: "string",
                description: "テーブル名",
              },
              data: {
                type: "object",
                description: "挿入するデータ（キー：カラム名、値：データ）",
              },
            },
            required: ["table", "data"],
          },
        },
        {
          name: "db_update",
          description: "データの更新",
          inputSchema: {
            type: "object",
            properties: {
              table: {
                type: "string",
                description: "テーブル名",
              },
              data: {
                type: "object",
                description: "更新するデータ",
              },
              where: {
                type: "string",
                description: "WHERE 条件（必須）",
              },
            },
            required: ["table", "data", "where"],
          },
        },
        {
          name: "db_delete",
          description: "データの削除",
          inputSchema: {
            type: "object",
            properties: {
              table: {
                type: "string",
                description: "テーブル名",
              },
              where: {
                type: "string",
                description: "WHERE 条件（必須）",
              },
            },
            required: ["table", "where"],
          },
        },
        {
          name: "db_execute_sql",
          description: "任意のSQLクエリの実行（上級者向け）",
          inputSchema: {
            type: "object",
            properties: {
              sql: {
                type: "string",
                description: "実行するSQLクエリ",
              },
              params: {
                type: "array",
                description: "SQLパラメータ",
              },
            },
            required: ["sql"],
          },
        },
      ],
    }));

    // ツールの実装
    server.setRequestHandler(CallToolRequestSchema, async (request) => {
      const { name, arguments: args } = request.params;

      try {
        let result;

        switch (name) {
          case "db_list_tables":
            result = await this.handleListTables();
            break;

          case "db_describe_table":
            result = await this.handleDescribeTable(args);
            break;

          case "db_select":
            result = await this.handleSelect(args);
            break;

          case "db_insert":
            result = await this.handleInsert(args);
            break;

          case "db_update":
            result = await this.handleUpdate(args);
            break;

          case "db_delete":
            result = await this.handleDelete(args);
            break;

          case "db_execute_sql":
            result = await this.handleExecuteSQL(args);
            break;

          default:
            throw new Error(`Unknown tool: ${name}`);
        }

        return {
          content: [
            {
              type: "text",
              text: result,
            },
          ],
        };
      } catch (error) {
        this.logger.error(`Tool execution failed: ${name}`, {
          error: error.message,
          args,
        });

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
  }

  // ハンドラーメソッドの実装
  async handleListTables() {
    const tables = await this.db.getTables();
    this.logger.info("Listed tables", { count: tables.length });

    return JSON.stringify(
      {
        success: true,
        tables: tables,
        count: tables.length,
      },
      null,
      2
    );
  }

  async handleDescribeTable(args) {
    const tableName = Validator.validateTableName(args.table);
    const tableInfo = await this.db.getTableInfo(tableName);

    this.logger.info("Described table", { table: tableName });

    return JSON.stringify(
      {
        success: true,
        table: tableName,
        columns: tableInfo,
      },
      null,
      2
    );
  }

  async handleSelect(args) {
    const validated = Validator.validateQueryParams(args);
    const { table, columns, where, orderBy, limit, offset } = validated;

    // クエリ構築
    let query = `SELECT ${columns ? columns.join(", ") : "*"} FROM ${table}`;
    const params = [];
    let paramIndex = 1;

    if (where) {
      query += ` WHERE ${where}`;
    }

    if (orderBy) {
      query += ` ORDER BY ${orderBy}`;
    }

    const actualLimit = Math.min(
      limit || this.config.max_results,
      this.config.max_results
    );

    if (this.config.database_type === "postgresql") {
      query += ` LIMIT $${paramIndex++}`;
      params.push(actualLimit);

      if (offset) {
        query += ` OFFSET $${paramIndex}`;
        params.push(offset);
      }
    } else {
      query += ` LIMIT ${actualLimit}`;
      if (offset) {
        query += ` OFFSET ${offset}`;
      }
    }

    const results = await this.db.select(query, params);

    this.logger.info("Executed SELECT", {
      table,
      rows: results.length,
      query: query.replace(/\s+/g, " "),
    });

    return JSON.stringify(
      {
        success: true,
        table,
        query: query.replace(/\s+/g, " "),
        rowCount: results.length,
        data: results,
      },
      null,
      2
    );
  }

  async handleInsert(args) {
    const tableName = Validator.validateTableName(args.table);
    const data = Validator.validateInsertData(args.data);

    const columns = Object.keys(data);
    const values = Object.values(data);

    let query, params;

    if (this.config.database_type === "postgresql") {
      const placeholders = values.map((_, i) => `$${i + 1}`).join(", ");
      query = `INSERT INTO ${tableName} (${columns.join(
        ", "
      )}) VALUES (${placeholders})`;
      params = values;
    } else {
      const placeholders = values.map(() => "?").join(", ");
      query = `INSERT INTO ${tableName} (${columns.join(
        ", "
      )}) VALUES (${placeholders})`;
      params = values;
    }

    const result = await this.db.execute(query, params);

    this.logger.info("Executed INSERT", {
      table: tableName,
      changes: result.changes,
      lastID: result.lastID,
    });

    return JSON.stringify(
      {
        success: true,
        table: tableName,
        changes: result.changes,
        lastID: result.lastID,
        insertedData: data,
      },
      null,
      2
    );
  }

  async handleUpdate(args) {
    const tableName = Validator.validateTableName(args.table);
    const data = Validator.validateInsertData(args.data);
    const where = args.where;

    if (!where || where.trim() === "") {
      throw new Error("WHERE clause is required for UPDATE operations");
    }

    const setParts = Object.keys(data).map((col, i) => {
      if (this.config.database_type === "postgresql") {
        return `${col} = $${i + 1}`;
      } else {
        return `${col} = ?`;
      }
    });

    let query = `UPDATE ${tableName} SET ${setParts.join(", ")} WHERE ${where}`;
    const params = Object.values(data);

    const result = await this.db.execute(query, params);

    this.logger.info("Executed UPDATE", {
      table: tableName,
      changes: result.changes,
    });

    return JSON.stringify(
      {
        success: true,
        table: tableName,
        changes: result.changes,
        updatedData: data,
        whereClause: where,
      },
      null,
      2
    );
  }

  async handleDelete(args) {
    const tableName = Validator.validateTableName(args.table);
    const where = args.where;

    if (!where || where.trim() === "") {
      throw new Error("WHERE clause is required for DELETE operations");
    }

    const query = `DELETE FROM ${tableName} WHERE ${where}`;
    const result = await this.db.execute(query);

    this.logger.info("Executed DELETE", {
      table: tableName,
      changes: result.changes,
    });

    return JSON.stringify(
      {
        success: true,
        table: tableName,
        changes: result.changes,
        whereClause: where,
      },
      null,
      2
    );
  }

  async handleExecuteSQL(args) {
    const { sql, params = [] } = args;

    // 危険なSQL操作の検出
    const dangerousPatterns = [
      /DROP\s+(?:TABLE|DATABASE|SCHEMA)/i,
      /ALTER\s+TABLE.*DROP/i,
      /TRUNCATE/i,
      /DELETE\s+FROM\s+\w+\s*(?:WHERE\s+1=1|WHERE\s+TRUE|;)/i,
    ];

    for (const pattern of dangerousPatterns) {
      if (pattern.test(sql)) {
        throw new Error("Potentially dangerous SQL operation detected");
      }
    }

    let result;
    if (sql.trim().toUpperCase().startsWith("SELECT")) {
      result = await this.db.select(sql, params);
      this.logger.info("Executed custom SELECT", {
        rows: result.length,
      });

      return JSON.stringify(
        {
          success: true,
          type: "SELECT",
          rowCount: result.length,
          data: result,
        },
        null,
        2
      );
    } else {
      result = await this.db.execute(sql, params);
      this.logger.info("Executed custom SQL", {
        changes: result.changes,
      });

      return JSON.stringify(
        {
          success: true,
          type: "MODIFY",
          changes: result.changes,
          lastID: result.lastID,
        },
        null,
        2
      );
    }
  }

  async start() {
    try {
      await this.initializeDatabase();

      const transport = new StdioServerTransport();
      await this.server.connect(transport);

      this.logger.info("Database connector server started");
    } catch (error) {
      this.logger.error("Failed to start server", { error: error.message });
      process.exit(1);
    }
  }

  async stop() {
    if (this.db) {
      await this.db.disconnect();
      this.logger.info("Database disconnected");
    }
  }
}

// サーバーの起動
const connector = new DatabaseConnector();

// 終了シグナルの処理
process.on("SIGINT", async () => {
  console.log("\nShutting down gracefully...");
  await connector.stop();
  process.exit(0);
});

process.on("SIGTERM", async () => {
  console.log("\nShutting down gracefully...");
  await connector.stop();
  process.exit(0);
});

// サーバー開始
connector.start().catch((error) => {
  console.error("Failed to start server:", error);
  process.exit(1);
});
```

## 🚀 実装とテスト

### 1. プロジェクトの作成

```bash
# DXTプロジェクトの初期化
dxt init database-connector
cd database-connector

# 必要なディレクトリの作成
mkdir -p server/database server/utils assets
```

### 2. 依存関係のインストール

```bash
cd server
npm install @modelcontextprotocol/sdk sqlite3 pg joi
```

### 3. サンプルデータの準備

`assets/schema.sql`:

```sql
-- ユーザーテーブル
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- プロジェクトテーブル
CREATE TABLE projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    user_id INTEGER,
    status TEXT DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- サンプルデータ
INSERT INTO users (name, email) VALUES
('Alice Johnson', 'alice@example.com'),
('Bob Smith', 'bob@example.com'),
('Carol Davis', 'carol@example.com');

INSERT INTO projects (name, description, user_id, status) VALUES
('Website Redesign', 'Complete redesign of company website', 1, 'active'),
('Mobile App', 'iOS and Android app development', 2, 'active'),
('Data Migration', 'Legacy system data migration', 1, 'completed'),
('API Integration', 'Third-party API integration', 3, 'active');
```

### 4. ローカルテスト

```bash
# SQLiteで基本テスト
cd server
DATABASE_TYPE=sqlite DATABASE_PATH=./test.db ENABLE_LOGGING=true node index.js

# 別のターミナルでテスト
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | node index.js
```

### 5. DXT ファイルの作成

```bash
# プロジェクトルートに戻る
cd ..

# DXTファイルを作成
dxt pack

# 作成されたファイルを確認
ls -la *.dxt
```

## 🧪 使用例

### Claude Desktop での使用

1. **テーブル一覧の取得**:

   ```
   利用可能なテーブルを教えてください
   ```

2. **データの検索**:

   ```
   usersテーブルから全てのユーザー情報を取得してください
   ```

3. **データの挿入**:

   ```
   usersテーブルに新しいユーザー「David Wilson (david@example.com)」を追加してください
   ```

4. **データの更新**:

   ```
   ユーザーID 1 のメールアドレスを "alice.new@example.com" に更新してください
   ```

5. **複雑なクエリ**:
   ```
   アクティブなプロジェクトを持つユーザーの一覧を取得してください
   ```

### 期待される結果

```json
{
  "success": true,
  "table": "users",
  "rowCount": 3,
  "data": [
    {
      "id": 1,
      "name": "Alice Johnson",
      "email": "alice@example.com",
      "created_at": "2024-01-15 10:30:00"
    },
    {
      "id": 2,
      "name": "Bob Smith",
      "email": "bob@example.com",
      "created_at": "2024-01-15 11:15:00"
    }
  ]
}
```

## 🔧 設定例

### SQLite 設定（開発用）

```json
{
  "database_type": "sqlite",
  "database_path": "./data/myapp.db",
  "max_results": 50,
  "enable_logging": true
}
```

### PostgreSQL 設定（本番用）

```json
{
  "database_type": "postgresql",
  "postgresql_url": "postgresql://user:pass@localhost:5432/mydb",
  "max_results": 100,
  "enable_logging": false
}
```

## 🚨 トラブルシューティング

### よくある問題

#### 1. データベース接続エラー

```
エラー: PostgreSQL connection failed: connect ECONNREFUSED
```

**解決策**:

- PostgreSQL サーバーが起動しているか確認
- 接続 URL が正しいか確認
- ファイアウォール設定を確認

#### 2. テーブルが見つからない

```
エラー: relation "users" does not exist
```

**解決策**:

- テーブル名の大文字小文字を確認
- スキーマが正しく適用されているか確認
- データベースが正しく選択されているか確認

#### 3. 権限エラー

```
エラー: permission denied for table users
```

**解決策**:

- データベースユーザーの権限を確認
- GRANT 文で適切な権限を付与

### デバッグ方法

```bash
# 詳細ログを有効にしてテスト
ENABLE_LOGGING=true DEBUG=* node server/index.js

# SQLクエリの確認
console.log("Executing query:", query, "with params:", params);
```

## 🎯 実践的な活用例

### 1. 営業データ分析

```
今月の売上上位10件の案件を取得し、担当者ごとに集計してください
```

### 2. 顧客サポート

```
過去30日間のサポートチケットを重要度別に分類し、未解決のものを一覧表示してください
```

### 3. 在庫管理

```
在庫が10個未満の商品を抽出し、発注が必要な商品リストを作成してください
```

## 🔒 セキュリティベストプラクティス

### 1. SQL インジェクション対策

```javascript
// ✅ 良い例：パラメータ化クエリ
const query = "SELECT * FROM users WHERE id = $1";
const result = await db.select(query, [userId]);

// ❌ 悪い例：文字列結合
const query = `SELECT * FROM users WHERE id = ${userId}`; // 危険！
```

### 2. アクセス制御

```javascript
// 重要なテーブルへのアクセス制限
const restrictedTables = ["admin_users", "payment_info", "secrets"];
if (restrictedTables.includes(tableName)) {
  throw new Error("Access denied to restricted table");
}
```

### 3. ログのサニタイゼーション

```javascript
// 機密情報をログから除外
const sanitizeQuery = (query) => {
  return query.replace(/password\s*=\s*['"][^'"]*['"]/gi, "password=***");
};
```

## 📈 パフォーマンス最適化

### 1. 接続プーリング（本番環境）

```javascript
import { Pool } from "pg";

const pool = new Pool({
  connectionString: config.postgresql_url,
  max: 20,
  idleTimeoutMillis: 30000,
  connectionTimeoutMillis: 2000,
});
```

### 2. クエリ最適化

```javascript
// インデックスを活用したクエリ
const optimizedQuery = `
  SELECT u.name, COUNT(p.id) as project_count
  FROM users u
  LEFT JOIN projects p ON u.id = p.user_id
  WHERE u.created_at >= $1
  GROUP BY u.id, u.name
  ORDER BY project_count DESC
  LIMIT $2
`;
```

## 🎉 まとめ

このデータベース連携 DXT 拡張機能により、以下が実現できました：

### ✅ 達成項目

- **マルチ DB 対応**: SQLite と PostgreSQL の両方をサポート
- **セキュア実装**: SQL インジェクション対策、入力検証
- **CRUD 操作**: 完全なデータベース操作機能
- **エラーハンドリング**: 堅牢なエラー処理とログ出力
- **設定可能**: ユーザー設定による柔軟な環境構築

### 🚀 次のステップ

この基礎をもとに、以下のような拡張も可能です：

- **高度な分析機能**: 集計、グループ化、ピボット
- **データ可視化**: チャート生成、レポート作成
- **リアルタイム同期**: WebSocket によるリアルタイム更新
- **バックアップ機能**: 自動バックアップとリストア

データベース連携をマスターしたら、次は [セキュリティとエンタープライズ機能](./05-security-enterprise.md) に進んで、より高度なセキュリティ機能を学習しましょう。
