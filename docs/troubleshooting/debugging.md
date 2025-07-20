# デバッグガイド

SilentCastの問題診断と詳細なデバッグ手法について説明します。

## 🔍 基本デバッグ手順

### デバッグモードの有効化

```bash
# デバッグモードで実行
./silentcast --debug --no-tray

# 詳細ログファイル出力
./silentcast --debug --log-file debug.log --no-tray

# 特定コンポーネントのデバッグ
./silentcast --debug --log-level trace --component hotkey
```

### ログレベルの設定

```yaml
# spellbook.yml - ログ設定
logger:
  level: debug        # trace, debug, info, warn, error
  file: "silentcast.log"
  console: true
  component_levels:
    hotkey: trace
    action: debug
    config: info
```

### 基本診断コマンド

```bash
# システム情報収集
./silentcast --system-info

# 設定検証
./silentcast --validate-config

# 権限確認
./silentcast --check-permissions

# コンポーネント状態確認
./silentcast --status
```

## 📊 ログ分析

### ログファイルの場所

```bash
# デフォルトログ場所
# Linux/macOS
~/.local/share/silentcast/silentcast.log

# Windows
%APPDATA%\SilentCast\silentcast.log

# カスタムログファイル
./silentcast --log-file /path/to/custom.log
```

### ログ形式の理解

```
# 標準ログ形式
[2025-01-20 14:30:15.123] [INFO ] [main    ] Application starting...
[2025-01-20 14:30:15.150] [DEBUG] [config  ] Loading configuration from: spellbook.yml
[2025-01-20 14:30:15.200] [DEBUG] [hotkey  ] Registering prefix hotkey: alt+space
[2025-01-20 14:30:15.250] [INFO ] [hotkey  ] Hotkey manager initialized
[2025-01-20 14:30:15.300] [ERROR] [action  ] Failed to execute action: editor

# フィールド説明:
# [タイムスタンプ] [レベル] [コンポーネント] メッセージ
```

### 重要なログメッセージ

#### 正常起動時のログ
```
[INFO ] Application starting...
[DEBUG] Configuration loaded successfully
[DEBUG] Hotkey manager initialized  
[INFO ] Ready to receive hotkeys
[INFO ] Application ready
```

#### エラー指標
```
[ERROR] Configuration load failed
[ERROR] Permission denied
[ERROR] Hotkey registration failed
[FATAL] Critical component initialization failed
```

### ログ分析スクリプト

```bash
#!/bin/bash
# ログ分析スクリプト

analyze_logs() {
    local log_file="${1:-silentcast.log}"
    
    if [ ! -f "$log_file" ]; then
        echo "❌ Log file not found: $log_file"
        return 1
    fi
    
    echo "📊 SilentCast ログ分析: $log_file"
    echo "=================================="
    
    # 基本統計
    echo "📈 ログ統計:"
    echo "  Total lines: $(wc -l < "$log_file")"
    echo "  ERROR lines: $(grep -c "\[ERROR\]" "$log_file")"
    echo "  WARN lines:  $(grep -c "\[WARN\]" "$log_file")"
    echo "  INFO lines:  $(grep -c "\[INFO\]" "$log_file")"
    echo "  DEBUG lines: $(grep -c "\[DEBUG\]" "$log_file")"
    
    # 最新エラー
    echo ""
    echo "🚨 最新エラー (最大5件):"
    grep "\[ERROR\]" "$log_file" | tail -5
    
    # ホットキー関連問題
    echo ""
    echo "⌨️ ホットキー関連ログ:"
    grep -i "hotkey" "$log_file" | tail -3
    
    # アクション実行ログ
    echo ""
    echo "⚡ アクション実行ログ:"
    grep -i "action\|execute" "$log_file" | tail -3
    
    # 起動・終了ログ
    echo ""
    echo "🚀 起動・終了ログ:"
    grep -E "(starting|ready|stopping|shutdown)" "$log_file" | tail -3
}

# 実行
analyze_logs "$1"
```

## 🛠️ コンポーネント別デバッグ

### ホットキーシステムデバッグ

```bash
# ホットキー詳細デバッグ
./silentcast --debug --component hotkey --test-hotkey

# ホットキー登録状態確認
./silentcast --list-hotkeys

# キー入力監視
./silentcast --monitor-keys
```

#### ホットキーデバッグログ例
```
[DEBUG] [hotkey] Initializing hotkey manager
[DEBUG] [hotkey] Registering prefix: alt+space (keycode: 65, modifiers: 8)
[DEBUG] [hotkey] Hook installed successfully
[TRACE] [hotkey] Key event: key=32 mod=8 state=down
[TRACE] [hotkey] Prefix detected: alt+space
[DEBUG] [hotkey] Waiting for sequence (timeout: 1000ms)
[TRACE] [hotkey] Sequence key: e (keycode: 101)
[DEBUG] [hotkey] Sequence complete: e
[DEBUG] [hotkey] Triggering action: editor
```

### 設定システムデバッグ

```bash
# 設定読み込みデバッグ
./silentcast --debug --component config --validate-config

# 設定カスケード確認
./silentcast --show-config-cascade

# 設定監視デバッグ  
./silentcast --debug --component watcher
```

#### 設定デバッグ設定
```yaml
# デバッグ用設定
logger:
  level: trace
  component_levels:
    config: trace
    watcher: debug
    
debug:
  config_validation: true
  show_cascade: true
  dump_resolved_config: true
```

### アクション実行デバッグ

```bash
# アクション実行デバッグ
./silentcast --debug --component action --once --spell test

# アクション タイムアウト延長
./silentcast --debug --action-timeout 300

# 環境変数デバッグ
./silentcast --debug --show-env
```

#### アクション実行デバッグログ
```
[DEBUG] [action] Executing action: editor
[TRACE] [action] Action type: app
[TRACE] [action] Command: code
[TRACE] [action] Args: [--new-window]
[TRACE] [action] Working dir: /home/user
[TRACE] [action] Environment: PATH=/usr/bin:...
[DEBUG] [action] Process started: PID 12345
[DEBUG] [action] Process completed: exit code 0 (duration: 1.2s)
```

## 🧪 手動テスト手順

### ステップバイステップテスト

#### 1. 設定テスト
```bash
# 設定ファイル構文確認
./silentcast --validate-config

# 設定値確認
./silentcast --dump-config

# 特定スペル確認
./silentcast --show-spell editor
```

#### 2. ホットキーテスト
```bash
# ホットキー登録確認
./silentcast --check-hotkeys

# 手動ホットキーテスト
./silentcast --test-hotkey alt+space

# シーケンステスト
./silentcast --test-sequence "alt+space,e"
```

#### 3. アクションテスト
```bash
# 個別アクションテスト
./silentcast --once --spell editor

# アクション詳細確認
./silentcast --dry-run --spell editor

# アクション環境確認
./silentcast --show-action-env editor
```

### テスト用最小設定

```yaml
# test-spellbook.yml
daemon:
  auto_start: false
  log_level: debug
  
logger:
  level: debug
  console: true
  
hotkeys:
  prefix: "alt+space"
  timeout: 1000
  
spells:
  test: test-action
  
grimoire:
  test-action:
    type: script
    command: echo "Test successful: $(date)"
    show_output: true
    description: "デバッグテスト用アクション"
```

## 🔬 高度なデバッグ技法

### システムコール監視

#### Linux (strace)
```bash
# システムコール監視
strace -f -o trace.log ./silentcast --debug --no-tray

# ファイルアクセス監視
strace -e trace=file ./silentcast --debug --no-tray

# ネットワークアクセス監視  
strace -e trace=network ./silentcast --debug --no-tray
```

#### macOS (dtruss)
```bash
# システムコール監視
sudo dtruss -f ./silentcast --debug --no-tray

# ファイルアクセス監視
sudo fs_usage -w -f filesys ./silentcast
```

### メモリとパフォーマンス分析

#### メモリリーク検出
```bash
# Valgrind (Linux)
valgrind --leak-check=full ./silentcast --debug --no-tray

# Instruments (macOS)
xcrun instruments -t Leaks -D trace.trace ./silentcast
```

#### パフォーマンス監視
```bash
# CPU 使用量監視
top -p $(pgrep silentcast)

# メモリ使用量監視
ps aux | grep silentcast

# リアルタイム監視
watch -n 1 'ps aux | grep silentcast'
```

### ネットワーク監視

```bash
# ネットワーク接続確認
netstat -an | grep silentcast
ss -tuln | grep silentcast

# パケットキャプチャ
sudo tcpdump -i any host localhost and port 8080
```

## 🔧 デバッグツールとユーティリティ

### 内蔵デバッグツール

```bash
# システム診断
./silentcast --diagnose

# 詳細システム情報
./silentcast --system-info --verbose

# 設定診断
./silentcast --config-doctor

# パフォーマンス情報
./silentcast --performance-info
```

### カスタムデバッグビルド

```bash
# デバッグビルド作成
make build-debug

# 詳細デバッグビルド
CGO_ENABLED=1 go build -tags "debug trace" -gcflags "-N -l" cmd/silentcast/main.go
```

### デバッグ環境変数

```bash
# デバッグレベル設定
export SILENTCAST_DEBUG=1
export SILENTCAST_LOG_LEVEL=trace

# コンポーネント別デバッグ
export SILENTCAST_DEBUG_HOTKEY=1
export SILENTCAST_DEBUG_ACTION=1
export SILENTCAST_DEBUG_CONFIG=1

# 詳細ログ
export SILENTCAST_VERBOSE=1
```

## 📋 デバッグ手順チェックリスト

### 問題発生時の初期対応
- [ ] デバッグモードで実行
- [ ] ログファイルを確認
- [ ] 設定ファイルを検証
- [ ] 権限状態を確認
- [ ] システム情報を収集

### 詳細調査
- [ ] コンポーネント別ログを分析
- [ ] 手動テストを実行
- [ ] 環境変数を確認
- [ ] 依存関係を検証
- [ ] プラットフォーム固有設定を確認

### 高度な分析
- [ ] システムコール監視
- [ ] パフォーマンス分析
- [ ] メモリ使用量確認
- [ ] ネットワーク監視
- [ ] 外部ツールとの統合

## 🆘 問題報告の準備

### 情報収集スクリプト
```bash
#!/bin/bash
# 問題報告用情報収集

collect_debug_info() {
    local output_dir="silentcast_debug_$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$output_dir"
    
    echo "🔍 SilentCast デバッグ情報収集中..."
    
    # システム情報
    echo "=== System Information ===" > "$output_dir/system_info.txt"
    uname -a >> "$output_dir/system_info.txt"
    echo >> "$output_dir/system_info.txt"
    
    # SilentCast バージョン
    ./silentcast --version > "$output_dir/version.txt" 2>&1
    
    # 設定ファイル
    cp spellbook*.yml "$output_dir/" 2>/dev/null || true
    
    # ログファイル
    cp silentcast*.log "$output_dir/" 2>/dev/null || true
    
    # システム診断
    ./silentcast --diagnose > "$output_dir/diagnosis.txt" 2>&1
    
    # 設定検証
    ./silentcast --validate-config > "$output_dir/config_validation.txt" 2>&1
    
    # 権限確認
    ./silentcast --check-permissions > "$output_dir/permissions.txt" 2>&1
    
    # アーカイブ作成
    tar -czf "${output_dir}.tar.gz" "$output_dir"
    rm -rf "$output_dir"
    
    echo "✅ デバッグ情報収集完了: ${output_dir}.tar.gz"
    echo "このファイルを問題報告に添付してください。"
}

collect_debug_info
```

### 問題報告テンプレート

```markdown
## SilentCast 問題報告

### 環境情報
- OS: [Windows/macOS/Linux + バージョン]
- SilentCast バージョン: [バージョン番号]
- ビルドタイプ: [通常/スタブ]

### 問題の説明
[問題の詳細な説明]

### 再現手順
1. [手順1]
2. [手順2]
3. [手順3]

### 期待される動作
[期待していた結果]

### 実際の動作
[実際に起こった結果]

### ログ出力
```
[関連するログ出力をここに貼り付け]
```

### 設定ファイル
```yaml
[spellbook.yml の関連部分]
```

### 追加情報
[その他の関連情報]
```

## 🔗 関連リソース

- [ホットキー問題](hotkeys.md)
- [権限設定](permissions.md)
- [アクション実行問題](actions.md)
- [プラットフォーム固有問題](platform-specific.md)
- [FAQ](faq.md)