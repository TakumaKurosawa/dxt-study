# セキュリティとエンタープライズ機能

## 🔐 セキュリティ概要

DXT は企業環境での安全な利用を前提として設計されており、複数のセキュリティ層を提供しています。

### セキュリティの基本原則

1. **最小権限の原則**: 必要最小限のアクセス権限のみを付与
2. **深層防御**: 複数のセキュリティ層による保護
3. **透明性**: すべての操作がログ記録され、監査可能
4. **分離**: 各拡張機能は独立したサンドボックス環境で実行

## 🔑 認証情報の安全な管理

### OS キーチェーンの活用

DXT は機密情報を OS のセキュアストレージに保存します：

#### macOS (Keychain)

```json
{
  "user_config": {
    "api_key": {
      "type": "string",
      "title": "API Key",
      "description": "外部サービスのAPIキー",
      "sensitive": true,
      "required": true
    },
    "database_password": {
      "type": "string",
      "title": "Database Password",
      "description": "データベース接続パスワード",
      "sensitive": true,
      "required": true
    }
  }
}
```

#### Windows (Credential Manager)

```json
{
  "user_config": {
    "service_token": {
      "type": "string",
      "title": "Service Token",
      "description": "サービストークン",
      "sensitive": true,
      "validation": {
        "pattern": "^sk-[a-zA-Z0-9]{32}$",
        "message": "有効なトークン形式で入力してください"
      }
    }
  }
}
```

### 環境変数での機密情報注入

```javascript
// サーバーコードでの機密情報の取得
const apiKey = process.env.API_KEY; // ${user_config.api_key}から自動注入
const dbPassword = process.env.DATABASE_PASSWORD;

// 接続設定
const dbConfig = {
  host: process.env.DATABASE_HOST,
  user: process.env.DATABASE_USER,
  password: dbPassword, // セキュアストレージから取得
  database: process.env.DATABASE_NAME,
};
```

## 🏢 エンタープライズ向け管理機能

### Windows Group Policy 設定

#### GPO テンプレートの作成

`DXTExtensions.admx`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<policyDefinitions xmlns:xsd="http://www.w3.org/2001/XMLSchema"
                   xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                   revision="1.0"
                   schemaVersion="1.0"
                   xmlns="http://schemas.microsoft.com/GroupPolicy/2006/07/PolicyDefinitions">

  <policyNamespaces>
    <target prefix="dxtextensions" namespace="Anthropic.Claude.DXTExtensions" />
  </policyNamespaces>

  <resources minRequiredRevision="1.0" />

  <categories>
    <category name="DXTExtensions" displayName="$(string.DXTExtensions)">
      <parentCategory ref="windows:WindowsComponents" />
    </category>
  </categories>

  <policies>
    <!-- 許可された拡張機能リスト -->
    <policy name="AllowedExtensions"
            class="Machine"
            displayName="$(string.AllowedExtensions)"
            explainText="$(string.AllowedExtensions_Explain)"
            key="SOFTWARE\Policies\Anthropic\Claude\Extensions"
            valueName="AllowedExtensions">
      <parentCategory ref="DXTExtensions" />
      <supportedOn ref="windows:SUPPORTED_WindowsVista" />
      <elements>
        <list id="AllowedExtensionsList" valueName="AllowedExtensions" />
      </elements>
    </policy>

    <!-- ブロックされた拡張機能リスト -->
    <policy name="BlockedExtensions"
            class="Machine"
            displayName="$(string.BlockedExtensions)"
            explainText="$(string.BlockedExtensions_Explain)"
            key="SOFTWARE\Policies\Anthropic\Claude\Extensions"
            valueName="BlockedExtensions">
      <parentCategory ref="DXTExtensions" />
      <supportedOn ref="windows:SUPPORTED_WindowsVista" />
      <elements>
        <list id="BlockedExtensionsList" valueName="BlockedExtensions" />
      </elements>
    </policy>

    <!-- 拡張機能ディレクトリの無効化 -->
    <policy name="DisableExtensionDirectory"
            class="Machine"
            displayName="$(string.DisableExtensionDirectory)"
            explainText="$(string.DisableExtensionDirectory_Explain)"
            key="SOFTWARE\Policies\Anthropic\Claude\Extensions"
            valueName="DisableExtensionDirectory">
      <parentCategory ref="DXTExtensions" />
      <supportedOn ref="windows:SUPPORTED_WindowsVista" />
      <enabledValue>
        <decimal value="1" />
      </enabledValue>
      <disabledValue>
        <decimal value="0" />
      </disabledValue>
    </policy>
  </policies>
</policyDefinitions>
```

#### PowerShell での一括設定

```powershell
# 許可された拡張機能の設定
$allowedExtensions = @(
    "com.company.internal-tools",
    "com.company.database-connector",
    "com.company.file-manager"
)

# レジストリキーの作成
$registryPath = "HKLM:\SOFTWARE\Policies\Anthropic\Claude\Extensions"
New-Item -Path $registryPath -Force | Out-Null

# 許可リストの設定
$allowedExtensions | ForEach-Object {
    $index = $allowedExtensions.IndexOf($_)
    Set-ItemProperty -Path $registryPath -Name "AllowedExtensions$index" -Value $_ -Type String
}

# ディレクトリの無効化
Set-ItemProperty -Path $registryPath -Name "DisableExtensionDirectory" -Value 1 -Type DWord

Write-Host "DXT Enterprise settings applied successfully"
```

### macOS MDM 設定

#### Configuration Profile

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadContent</key>
    <array>
        <dict>
            <key>PayloadDisplayName</key>
            <string>Claude Desktop Extensions</string>
            <key>PayloadIdentifier</key>
            <string>com.company.claude.extensions</string>
            <key>PayloadType</key>
            <string>com.anthropic.claude.desktop</string>
            <key>PayloadUUID</key>
            <string>12345678-1234-1234-1234-123456789012</string>
            <key>PayloadVersion</key>
            <integer>1</integer>

            <!-- 許可された拡張機能 -->
            <key>AllowedExtensions</key>
            <array>
                <string>com.company.internal-tools</string>
                <string>com.company.database-connector</string>
                <string>com.company.file-manager</string>
            </array>

            <!-- ブロックされた拡張機能 -->
            <key>BlockedExtensions</key>
            <array>
                <string>com.untrusted.*</string>
            </array>

            <!-- 拡張機能ディレクトリの無効化 -->
            <key>DisableExtensionDirectory</key>
            <true/>

            <!-- プライベート拡張機能ディレクトリ -->
            <key>ExtensionDirectories</key>
            <array>
                <string>https://extensions.company.com/</string>
            </array>
        </dict>
    </array>
    <key>PayloadDisplayName</key>
    <string>Claude Desktop Extensions Policy</string>
    <key>PayloadIdentifier</key>
    <string>com.company.claude.extensions.policy</string>
    <key>PayloadType</key>
    <string>Configuration</string>
    <key>PayloadUUID</key>
    <string>87654321-4321-4321-4321-210987654321</string>
    <key>PayloadVersion</key>
    <integer>1</integer>
</dict>
</plist>
```

#### Jamf Pro 設定

```bash
#!/bin/bash
# Jamf Pro Extension Attributes

# 現在インストールされているDXT拡張機能を取得
extensions_path="/Users/*/Library/Application Support/Claude/Extensions"

if [ -d "$extensions_path" ]; then
    installed_extensions=$(find "$extensions_path" -name "*.dxt" -exec basename {} .dxt \;)
    echo "<result>$installed_extensions</result>"
else
    echo "<result>No extensions installed</result>"
fi
```

## 🛡️ セキュリティベストプラクティス

### 拡張機能開発のセキュリティガイドライン

#### 1. 入力検証

```javascript
// 悪い例：入力検証なし
async function readFile(filePath) {
  return await fs.readFile(filePath, "utf8");
}

// 良い例：適切な入力検証
async function readFile(filePath) {
  // パスの正規化
  const normalizedPath = path.resolve(filePath);

  // 許可されたディレクトリ内かチェック
  if (!isPathAllowed(normalizedPath)) {
    throw new Error("Access denied: Path not allowed");
  }

  // パストラバーサル攻撃の防止
  if (normalizedPath.includes("..")) {
    throw new Error("Access denied: Invalid path");
  }

  // ファイルサイズ制限
  const stats = await fs.stat(normalizedPath);
  if (stats.size > MAX_FILE_SIZE) {
    throw new Error("File too large");
  }

  return await fs.readFile(normalizedPath, "utf8");
}
```

#### 2. SQL インジェクション対策

```javascript
// 悪い例：SQLインジェクション脆弱性
async function searchUsers(query) {
  const sql = `SELECT * FROM users WHERE name LIKE '%${query}%'`;
  return await db.query(sql);
}

// 良い例：パラメータ化クエリ
async function searchUsers(query) {
  // 入力の検証
  if (typeof query !== "string" || query.length > 100) {
    throw new Error("Invalid query parameter");
  }

  // エスケープ処理
  const sanitizedQuery = query.replace(/[%_]/g, "\\$&");

  // パラメータ化クエリ
  const sql = "SELECT id, name, email FROM users WHERE name LIKE ? LIMIT 50";
  return await db.query(sql, [`%${sanitizedQuery}%`]);
}
```

#### 3. 機密情報のログ出力防止

```javascript
// セキュアなログ出力クラス
class SecureLogger {
  static sensitivePatterns = [
    /api[_-]?key[s]?[:=]\s*[\"\']?([a-zA-Z0-9_-]+)[\"\']?/gi,
    /password[:=]\s*[\"\']?([^\s\"\',]+)[\"\']?/gi,
    /token[:=]\s*[\"\']?([a-zA-Z0-9_-]+)[\"\']?/gi,
  ];

  static sanitize(message) {
    let sanitized = message;

    this.sensitivePatterns.forEach((pattern) => {
      sanitized = sanitized.replace(pattern, (match) => {
        return match.replace(/([a-zA-Z0-9_-]{4})[a-zA-Z0-9_-]+/g, "$1****");
      });
    });

    return sanitized;
  }

  static log(level, message, metadata = {}) {
    const sanitizedMessage = this.sanitize(message);
    const sanitizedMetadata = JSON.parse(
      this.sanitize(JSON.stringify(metadata))
    );

    console.log({
      timestamp: new Date().toISOString(),
      level,
      message: sanitizedMessage,
      metadata: sanitizedMetadata,
    });
  }
}

// 使用例
SecureLogger.log("info", "Database connection established", {
  host: dbConfig.host,
  database: dbConfig.database,
  // password は自動的にマスクされる
});
```

### セキュリティ監査とコンプライアンス

#### 拡張機能のセキュリティチェックリスト

```markdown
## DXT 拡張機能セキュリティチェックリスト

### ✅ 基本セキュリティ

- [ ] 入力検証が適切に実装されている
- [ ] パラメータ化クエリを使用している
- [ ] パストラバーサル攻撃を防いでいる
- [ ] ファイルサイズ制限を設けている
- [ ] エラーメッセージに機密情報が含まれていない

### ✅ 認証・認可

- [ ] 機密情報は OS キーチェーンに保存されている
- [ ] アクセス権限が最小限に制限されている
- [ ] セッション管理が適切に実装されている
- [ ] タイムアウト機能が実装されている

### ✅ ログ・監査

- [ ] セキュリティイベントがログ記録されている
- [ ] ログに機密情報が含まれていない
- [ ] 異常なアクセスパターンを検出できる
- [ ] 監査証跡が保持されている

### ✅ 暗号化

- [ ] 通信は暗号化されている（HTTPS/TLS）
- [ ] 保存時暗号化が実装されている
- [ ] 強力な暗号化アルゴリズムを使用している
- [ ] 鍵管理が適切に行われている
```

## 📊 監視とログ機能

### セキュリティイベントの監視

```javascript
// セキュリティイベント監視クラス
class SecurityMonitor {
  constructor() {
    this.events = [];
    this.alertThresholds = {
      failedAttempts: 5,
      timeWindow: 300000, // 5分
    };
  }

  logSecurityEvent(eventType, details) {
    const event = {
      timestamp: Date.now(),
      type: eventType,
      details: details,
      severity: this.getSeverity(eventType),
    };

    this.events.push(event);

    // アラート条件をチェック
    this.checkAlertConditions(eventType);

    // ログ出力
    SecureLogger.log("security", `Security event: ${eventType}`, details);
  }

  getSeverity(eventType) {
    const severityMap = {
      authentication_failure: "medium",
      unauthorized_access: "high",
      suspicious_activity: "high",
      data_access: "low",
      privilege_escalation: "critical",
    };

    return severityMap[eventType] || "medium";
  }

  checkAlertConditions(eventType) {
    const now = Date.now();
    const windowStart = now - this.alertThresholds.timeWindow;

    const recentEvents = this.events.filter(
      (event) => event.timestamp >= windowStart && event.type === eventType
    );

    if (recentEvents.length >= this.alertThresholds.failedAttempts) {
      this.triggerAlert(eventType, recentEvents);
    }
  }

  triggerAlert(eventType, events) {
    const alert = {
      timestamp: Date.now(),
      type: "security_alert",
      eventType: eventType,
      count: events.length,
      timeWindow: this.alertThresholds.timeWindow,
    };

    // アラート通知（例：Webhookまたはメール）
    this.sendAlert(alert);
  }

  async sendAlert(alert) {
    // 実装例：Slackまたはメール通知
    const webhookUrl = process.env.SECURITY_WEBHOOK_URL;
    if (webhookUrl) {
      await fetch(webhookUrl, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          text: `🚨 Security Alert: ${alert.eventType}`,
          details: alert,
        }),
      });
    }
  }
}

// 使用例
const securityMonitor = new SecurityMonitor();

// 認証失敗の記録
securityMonitor.logSecurityEvent("authentication_failure", {
  user: "user@example.com",
  source: "extension_api",
  reason: "invalid_token",
});

// 不正アクセス試行の記録
securityMonitor.logSecurityEvent("unauthorized_access", {
  path: "/sensitive/data",
  method: "GET",
  userAgent: "suspicious-client/1.0",
});
```

## 🏛️ コンプライアンス対応

### GDPR 対応

```javascript
// GDPRコンプライアンス機能
class GDPRCompliance {
  async handleDataSubjectRequest(requestType, userId, details) {
    switch (requestType) {
      case "access":
        return await this.exportUserData(userId);
      case "rectification":
        return await this.updateUserData(userId, details);
      case "erasure":
        return await this.deleteUserData(userId);
      case "portability":
        return await this.exportPortableData(userId);
      default:
        throw new Error("Invalid request type");
    }
  }

  async exportUserData(userId) {
    const userData = {
      profile: await this.getUserProfile(userId),
      activityLog: await this.getUserActivity(userId),
      preferences: await this.getUserPreferences(userId),
    };

    // 個人データの匿名化
    return this.anonymizeData(userData);
  }

  async deleteUserData(userId) {
    // データの完全削除
    await this.deleteUserProfile(userId);
    await this.deleteUserActivity(userId);
    await this.deleteUserPreferences(userId);

    // 削除ログの記録
    SecureLogger.log("gdpr", "User data deleted", { userId });

    return { status: "completed", deletedAt: new Date().toISOString() };
  }
}
```

### SOC 2 対応

```javascript
// SOC 2コンプライアンス機能
class SOC2Compliance {
  constructor() {
    this.auditLog = [];
  }

  logAuditEvent(category, action, details) {
    const auditEvent = {
      timestamp: new Date().toISOString(),
      category: category, // CC6.1, CC6.2, etc.
      action: action,
      details: details,
      user: this.getCurrentUser(),
      source: "dxt_extension",
    };

    this.auditLog.push(auditEvent);

    // 外部監査システムに送信
    this.sendToAuditSystem(auditEvent);
  }

  async generateComplianceReport(startDate, endDate) {
    const events = this.auditLog.filter((event) => {
      const eventDate = new Date(event.timestamp);
      return eventDate >= startDate && eventDate <= endDate;
    });

    return {
      period: { start: startDate, end: endDate },
      totalEvents: events.length,
      categorizedEvents: this.categorizeEvents(events),
      securityIncidents: this.getSecurityIncidents(events),
      accessReports: this.getAccessReports(events),
    };
  }
}
```

## 🔧 実装例：企業向けセキュア DXT

完全な企業向け DXT の例を見てみましょう：

```json
{
  "dxt_version": "0.1",
  "name": "enterprise-secure-connector",
  "display_name": "Enterprise Secure Connector",
  "version": "1.0.0",
  "description": "企業向けセキュアデータ連携ツール",
  "author": {
    "name": "Corporate IT Department",
    "email": "it-security@company.com"
  },
  "license": "Proprietary",
  "security": {
    "required_permissions": ["network", "filesystem_read", "keychain_access"],
    "security_level": "enterprise",
    "compliance": ["SOC2", "GDPR", "HIPAA"]
  },
  "server": {
    "type": "node",
    "entry_point": "server/index.js",
    "mcp_config": {
      "command": "node",
      "args": ["${__dirname}/server/index.js"],
      "env": {
        "ENTERPRISE_API_KEY": "${user_config.enterprise_api_key}",
        "AUDIT_ENDPOINT": "${user_config.audit_endpoint}",
        "ENCRYPTION_KEY": "${user_config.encryption_key}",
        "COMPLIANCE_MODE": "${user_config.compliance_mode}"
      }
    }
  },
  "user_config": {
    "enterprise_api_key": {
      "type": "string",
      "title": "Enterprise API Key",
      "description": "企業システム接続用APIキー",
      "sensitive": true,
      "required": true,
      "validation": {
        "pattern": "^ent_[a-zA-Z0-9]{32}$"
      }
    },
    "audit_endpoint": {
      "type": "string",
      "title": "Audit Endpoint",
      "description": "監査ログ送信先エンドポイント",
      "default": "https://audit.company.com/api/v1/logs"
    },
    "encryption_key": {
      "type": "string",
      "title": "Encryption Key",
      "description": "データ暗号化キー",
      "sensitive": true,
      "required": true
    },
    "compliance_mode": {
      "type": "string",
      "title": "Compliance Mode",
      "description": "コンプライアンスモード",
      "enum": ["SOC2", "GDPR", "HIPAA", "ALL"],
      "default": "ALL"
    }
  }
}
```

このガイドを参考に、企業環境での DXT 拡張機能の安全な開発と運用を行ってください。

次は [高度な機能とベストプラクティス](./06-advanced-practices.md) で、より高度な実装テクニックを学習しましょう。
