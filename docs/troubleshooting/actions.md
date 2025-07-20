# アクション実行問題のトラブルシューティング

SilentCastのアクション実行に関する問題の診断と解決方法について説明します。

## 🚀 クイック診断

### アクションが実行されない場合

```bash
# 1. アクション設定の検証
./silentcast --validate-config

# 2. 手動でアクション実行テスト
./silentcast --once --spell <spell-name>

# 3. デバッグモードで詳細確認
./silentcast --debug --no-tray
```

### 一般的なアクション問題
1. **アプリケーション起動失敗** → パス解決、権限、依存関係
2. **スクリプト実行エラー** → シェル環境、権限、構文エラー  
3. **URL が開かない** → ブラウザ設定、プロトコルハンドラー
4. **昇格アクション失敗** → 権限、UAC、sudo 設定

## 📱 アプリケーション起動問題

### 問題: アプリケーションが起動しない

#### 症状
```
[ERROR] Failed to launch application: code
[WARN] Application not found in PATH
[ERROR] Exec format error
```

#### 原因と解決方法

**1. アプリケーションパス解決失敗**
```yaml
# 問題のある設定
grimoire:
  editor:
    type: app
    command: code  # PATHにない場合失敗
    
# 解決策1: フルパス指定
grimoire:
  editor:
    type: app
    command: "/usr/local/bin/code"
    
# 解決策2: プラットフォーム固有設定
grimoire:
  editor:
    type: app
    command: "code"
    windows:
      command: "C:\\Users\\user\\AppData\\Local\\Programs\\Microsoft VS Code\\bin\\code.cmd"
    macos:
      command: "/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code"
    linux:
      command: "/usr/bin/code"
```

**2. 環境変数の問題**
```yaml
# 環境変数を含む起動設定
grimoire:
  editor:
    type: app
    command: "code"
    env:
      PATH: "/usr/local/bin:$PATH"
      EDITOR: "code"
    working_dir: "${HOME}"
```

**3. 引数の問題**
```yaml
# 正しい引数指定
grimoire:
  editor:
    type: app
    command: "code"
    args: ["--new-window", "${HOME}/project"]
    
  browser:
    type: app
    command: "google-chrome"
    args: ["--incognito", "--new-window"]
```

### アプリケーション起動デバッグ

#### デバッグ設定
```yaml
# デバッグ用詳細ログ
logger:
  level: debug
  
# アプリ起動デバッグ
grimoire:
  test-app:
    type: app
    command: "echo"
    args: ["Application", "launched", "successfully"]
    show_output: true
    description: "アプリ起動テスト"
```

#### 手動検証手順
```bash
# 1. コマンド直接実行テスト
code --new-window

# 2. フルパスでテスト
/usr/local/bin/code --new-window

# 3. 環境変数確認
echo $PATH
which code

# 4. アクセス権限確認
ls -la $(which code)
```

## 📜 スクリプト実行問題

### 問題: スクリプトが実行されない

#### 症状
```
[ERROR] Script execution failed: permission denied
[WARN] Shell not found: /bin/bash
[ERROR] Command not found in script
```

#### 解決方法

**1. 実行権限の設定**
```bash
# スクリプトファイルの権限確認
ls -la /path/to/script.sh

# 実行権限付与
chmod +x /path/to/script.sh
```

**2. シェル指定**
```yaml
grimoire:
  git-status:
    type: script
    command: "git status"
    shell: "/bin/bash"     # 明示的なシェル指定
    
  python-script:
    type: script
    command: "python3 /path/to/script.py"
    shell: "/bin/bash"
    env:
      PYTHONPATH: "/usr/local/lib/python3.9"
```

**3. スクリプトファイル実行**
```yaml
grimoire:
  backup:
    type: script
    script_file: "/home/user/scripts/backup.sh"
    working_dir: "/home/user"
    timeout: 300  # 5分タイムアウト
```

### スクリプト実行環境設定

#### 環境変数設定
```yaml
# グローバル環境変数
environment:
  PATH: "/usr/local/bin:/usr/bin:/bin"
  HOME: "${USER_HOME}"
  LANG: "ja_JP.UTF-8"

grimoire:
  dev-setup:
    type: script
    command: |
      export NODE_ENV=development
      npm run dev
    env:
      NODE_VERSION: "18"
      npm_config_cache: "/tmp/npm-cache"
```

#### 複雑なスクリプト例
```yaml
grimoire:
  docker-deploy:
    type: script
    command: |
      #!/bin/bash
      set -e
      
      echo "Building Docker image..."
      docker build -t myapp:latest .
      
      echo "Stopping existing container..."
      docker stop myapp || true
      docker rm myapp || true
      
      echo "Starting new container..."
      docker run -d --name myapp -p 8080:8080 myapp:latest
      
      echo "Deployment completed!"
    working_dir: "/home/user/project"
    show_output: true
    timeout: 600
```

## 🌐 URL アクション問題

### 問題: URL が開かない

#### 症状
```
[ERROR] Failed to open URL: https://example.com
[WARN] No default browser configured
[ERROR] Protocol handler not found
```

#### 解決方法

**1. ブラウザ設定**
```yaml
# デフォルトブラウザ設定
browser:
  default: "google-chrome"
  fallback: ["firefox", "safari", "edge"]

grimoire:
  website:
    type: url
    command: "https://github.com"
    browser: "google-chrome"
    args: ["--incognito"]
```

**2. プラットフォーム固有設定**
```yaml
grimoire:
  open-url:
    type: url
    command: "https://example.com"
    windows:
      browser: "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
    macos:
      browser: "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
    linux:
      browser: "/usr/bin/google-chrome"
```

**3. カスタムプロトコルハンドラー**
```yaml
grimoire:
  vscode-project:
    type: url
    command: "vscode://file/home/user/project"
    description: "VS Code でプロジェクトを開く"
    
  slack-channel:
    type: url
    command: "slack://channel?team=T123&id=C456"
    description: "Slack チャンネルを開く"
```

### URL デバッグ

#### 手動テスト
```bash
# プラットフォーム別 URL 開きテスト

# Linux
xdg-open https://example.com

# macOS  
open https://example.com

# Windows (PowerShell)
Start-Process https://example.com
```

## ⬆️ 昇格アクション問題

### 問題: 昇格アクションが失敗する

#### 症状
```
[ERROR] Elevation failed: access denied
[WARN] sudo password prompt timeout
[ERROR] UAC dialog cancelled by user
```

#### 解決方法

**1. 昇格設定の確認**
```yaml
# 昇格アクション設定
grimoire:
  system-update:
    type: elevated
    command: "apt update && apt upgrade -y"
    description: "システムアップデート"
    require_confirmation: true
    timeout: 1800  # 30分
    
elevation:
  method: "auto"  # auto, sudo, pkexec, uac
  timeout: 300    # 5分
  prompt_user: true
```

**2. プラットフォーム別昇格設定**
```yaml
# Windows UAC
elevation:
  windows:
    method: "uac"
    runas_user: ""  # 現在のユーザーで昇格
    
# macOS sudo
elevation:
  macos:
    method: "sudo"
    sudo_prompt: true
    password_cache: 300  # 5分間キャッシュ
    
# Linux PolicyKit
elevation:
  linux:
    method: "pkexec"
    fallback: "sudo"
    display_name: "SilentCast"
```

**3. 昇格確認プロンプト**
```yaml
grimoire:
  restart-service:
    type: elevated
    command: "systemctl restart nginx"
    confirmation:
      title: "サービス再起動"
      message: "Nginx サービスを再起動しますか？"
      buttons: ["実行", "キャンセル"]
```

## 🔧 アクション設定パターン

### エラーハンドリング付きアクション
```yaml
grimoire:
  robust-backup:
    type: script
    command: |
      #!/bin/bash
      set -e
      
      BACKUP_DIR="/backup/$(date +%Y%m%d)"
      SOURCE_DIR="/home/user/important"
      
      # 事前チェック
      if [ ! -d "$SOURCE_DIR" ]; then
        echo "ERROR: Source directory not found: $SOURCE_DIR"
        exit 1
      fi
      
      # バックアップディレクトリ作成
      mkdir -p "$BACKUP_DIR"
      
      # バックアップ実行
      if rsync -av "$SOURCE_DIR/" "$BACKUP_DIR/"; then
        echo "SUCCESS: Backup completed to $BACKUP_DIR"
        notify-send "Backup Completed" "Files backed up successfully"
      else
        echo "ERROR: Backup failed"
        notify-send "Backup Failed" "Check logs for details"
        exit 1
      fi
    on_error: "notify"
    show_output: true
    timeout: 3600
```

### 条件付きアクション
```yaml
grimoire:
  conditional-deploy:
    type: script
    command: |
      if [ "$(git status --porcelain)" ]; then
        echo "ERROR: Working directory not clean"
        exit 1
      fi
      
      if [ "$(git rev-parse --abbrev-ref HEAD)" != "main" ]; then
        echo "ERROR: Not on main branch"
        exit 1
      fi
      
      echo "Deploying..."
      git push origin main
      ssh server "cd /app && git pull && systemctl restart app"
    description: "Git clean check & deploy"
```

### インタラクティブアクション
```yaml
grimoire:
  interactive-git:
    type: script
    command: |
      echo "Git 操作メニュー:"
      echo "1) Status"
      echo "2) Add all"  
      echo "3) Commit"
      echo "4) Push"
      read -p "選択 (1-4): " choice
      
      case $choice in
        1) git status ;;
        2) git add . ;;
        3) read -p "Commit message: " msg; git commit -m "$msg" ;;
        4) git push ;;
        *) echo "無効な選択" ;;
      esac
    interactive: true
    show_output: true
```

## 📊 アクション監視とログ

### 詳細ログ設定
```yaml
logger:
  level: debug
  file: "silentcast.log"
  action_log: "actions.log"  # アクション専用ログ
  
# アクション実行ログフォーマット
action_logging:
  format: "[{timestamp}] {spell} -> {action} ({duration}ms) [{status}]"
  include_output: true
  max_output_length: 1000
```

### パフォーマンス監視
```yaml
monitoring:
  action_timeout: 300      # デフォルトタイムアウト
  resource_monitoring: true
  cpu_threshold: 80        # CPU使用率警告
  memory_threshold: 1024   # メモリ使用量警告(MB)
```

## 📋 アクション問題診断チェックリスト

### 基本確認
- [ ] アクション設定の構文エラーがない
- [ ] 実行権限が適切に設定されている
- [ ] パスとファイルが存在する
- [ ] 環境変数が正しく設定されている

### アプリケーション起動
- [ ] アプリケーションが PATH にある
- [ ] フルパス指定が正確
- [ ] 引数とオプションが正しい
- [ ] 作業ディレクトリが適切

### スクリプト実行
- [ ] シェルが正しく指定されている
- [ ] スクリプトに実行権限がある
- [ ] 依存コマンドが利用可能
- [ ] タイムアウト設定が適切

### URL アクション
- [ ] URL フォーマットが正しい
- [ ] ブラウザが設定されている
- [ ] プロトコルハンドラーが登録済み

### 昇格アクション
- [ ] 昇格方式が適切に設定されている
- [ ] 権限確認プロセスが機能する
- [ ] タイムアウト設定が十分

## 🆘 それでも解決しない場合

### デバッグ手順
```bash
# 1. 設定ダンプ
./silentcast --dump-config > config-debug.yml

# 2. アクション単体テスト
./silentcast --once --spell test-action

# 3. 詳細ログ確認
./silentcast --debug --log-file debug.log --no-tray

# 4. システム監視
top -p $(pgrep silentcast)
```

### 代替実装
```yaml
# フォールバック付きアクション
grimoire:
  robust-editor:
    type: app
    command: "code"
    fallback:
      - command: "vim"
        type: app
      - command: "nano"
        type: app
    description: "エディタ（フォールバック付き）"
```

## 🔗 関連リソース

- [ホットキー問題](hotkeys.md)
- [権限設定](permissions.md)
- [設定ガイド](../guide/configuration.md)
- [デバッグガイド](debugging.md)