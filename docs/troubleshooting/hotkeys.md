# ホットキー問題のトラブルシューティング

SilentCastのホットキー機能で発生する問題の診断と解決方法について説明します。

## 🚀 クイック診断

### ホットキーが反応しない場合

```bash
# 1. アプリケーションの状態確認
./silentcast --debug --no-tray

# 2. ホットキー機能のテスト
./silentcast --test-hotkey

# 3. 設定の検証
./silentcast --validate-config
```

### 一般的なホットキー問題
1. **ホットキーが登録されない** → 権限不足またはプラットフォーム制限
2. **他のアプリとの競合** → キーの組み合わせ変更
3. **シーケンスが認識されない** → タイムアウト設定確認
4. **全く反応しない** → スタブモードで実行されている可能性

## 🔧 ホットキー登録問題

### 問題: ホットキーが登録できない

#### Windows
```powershell
# 管理者権限で実行
# PowerShell を管理者として開く
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# SilentCast を管理者として実行
Start-Process ./silentcast -Verb RunAs
```

#### macOS
```bash
# アクセシビリティ権限の確認
# システム環境設定 > セキュリティとプライバシー > プライバシー > アクセシビリティ
# SilentCast を追加し、チェックを入れる

# ターミナルにもアクセシビリティ権限が必要な場合があります
```

#### Linux
```bash
# X11 セッションの確認
echo $DISPLAY

# 必要な権限の確認
xhost +local:

# デスクトップ環境固有の設定
# GNOME
gsettings list-schemas | grep hotkey

# KDE
kwriteconfig5 --file kglobalshortcutsrc --group silentcast --key _k_friendly_name SilentCast
```

### 解決方法: ホットキー権限の設定

#### macOS アクセシビリティ権限
1. **システム環境設定** を開く
2. **セキュリティとプライバシー** → **プライバシー**
3. **アクセシビリティ** を選択
4. 🔒 をクリックして管理者パスワード入力
5. **+** ボタンで SilentCast を追加
6. SilentCast にチェックを入れる
7. アプリケーションを再起動

#### Windows UAC 設定
```powershell
# UAC レベルの確認
Get-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" -Name "ConsentPromptBehaviorAdmin"

# SilentCast を例外に追加（管理者として）
New-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" -Name "SilentCastBypass" -Value 1 -PropertyType DWORD
```

## ⌨️ キー競合の解決

### 問題: 他のアプリケーションとのキー競合

#### 競合の確認
```bash
# 現在のキー設定確認
./silentcast --list-spells

# デバッグモードで競合検出
./silentcast --debug --test-hotkey
```

#### 競合解決策

**1. キーの組み合わせ変更**
```yaml
# spellbook.yml
hotkeys:
  prefix: "ctrl+alt+space"  # デフォルトから変更
  timeout: 1000

spells:
  e: editor
  t: terminal
```

**2. より具体的なキー組み合わせ**
```yaml
spells:
  "ctrl+shift+e": editor
  "alt+f1": terminal
  "win+space": launcher
```

**3. プラットフォーム固有設定**
```yaml
# spellbook.windows.yml
hotkeys:
  prefix: "win+space"

# spellbook.darwin.yml  
hotkeys:
  prefix: "cmd+space"

# spellbook.linux.yml
hotkeys:
  prefix: "super+space"
```

### よくある競合アプリケーション

#### Windows
- **Windows キー + Space**: 入力方式切り替え
- **Ctrl + Shift + Esc**: タスクマネージャー
- **Alt + Tab**: アプリケーション切り替え

**回避策:**
```yaml
hotkeys:
  prefix: "ctrl+alt+grave"  # Ctrl+Alt+`
```

#### macOS
- **Cmd + Space**: Spotlight
- **Ctrl + Space**: 入力ソース切り替え
- **Cmd + Tab**: アプリケーション切り替え

**回避策:**
```yaml
hotkeys:
  prefix: "cmd+option+space"  # Cmd+Option+Space
```

#### Linux
- **Super + Space**: アクティビティ概要（GNOME）
- **Alt + F2**: 実行ダイアログ
- **Ctrl + Alt + T**: ターミナル

**回避策:**
```yaml
hotkeys:
  prefix: "super+alt+space"  # Super+Alt+Space
```

## ⏱️ タイムアウトとシーケンス問題

### 問題: キーシーケンスが認識されない

#### タイムアウト設定の調整
```yaml
hotkeys:
  prefix: "alt+space"
  timeout: 2000           # プレフィックス後の待機時間 (ms)
  sequence_timeout: 3000  # シーケンス全体のタイムアウト (ms)
```

#### シーケンス入力のコツ
1. **プレフィックスキーを確実に離す**
2. **次のキーまで少し待つ**（500ms程度）
3. **シーケンス内のキーは素早く入力**

#### 複雑なシーケンスの例
```yaml
spells:
  # 短いシーケンス（推奨）
  "g,s": git-status
  "g,p": git-pull
  
  # 長いシーケンス（注意が必要）
  "d,o,c,k,e,r": docker-status
  
grimoire:
  git-status:
    type: script
    command: git status
    description: "Git リポジトリの状態確認"
```

### デバッグ用設定
```yaml
# デバッグ情報を増やす
daemon:
  log_level: debug

logger:
  level: debug
  file: "silentcast.log"

# ホットキーテスト専用設定
spells:
  test: test-action
  
grimoire:
  test-action:
    type: script
    command: echo "ホットキーテスト成功"
    show_output: true
    notify: true
```

## 🔍 ホットキー診断

### 基本診断手順

#### 1. スタブモード確認
```bash
# スタブモードで実行されているか確認
./silentcast --version
# 出力に "stub" または "nogohook" が含まれている場合、ホットキー機能は無効
```

#### 2. ホットキーマネージャーの状態確認
```bash
# デバッグモードで詳細ログ確認
./silentcast --debug --no-tray

# ログで以下を確認:
# [DEBUG] Hotkey manager initialized
# [DEBUG] Registered hotkey: alt+space
# [DEBUG] Hotkey detected: alt+space
```

#### 3. 手動テスト
```bash
# 設定ファイルのテスト
./silentcast --dry-run

# 特定のスペル実行テスト
./silentcast --once --spell e

# ホットキー機能のテスト
./silentcast --test-hotkey
```

### 詳細ログ分析

#### 正常なホットキー登録ログ
```
[INFO] SilentCast starting...
[DEBUG] Configuration loaded: spellbook.yml
[DEBUG] Hotkey manager initializing...
[DEBUG] Registered prefix key: alt+space
[DEBUG] Registered 5 spells
[INFO] Ready to receive hotkeys
```

#### 問題のあるログ例
```
[ERROR] Failed to register hotkey: alt+space (already in use)
[WARN] Hotkey manager initialization failed, falling back to stub mode
[ERROR] Permission denied: accessibility features required
```

## 🛠️ 高度なトラブルシューティング

### システムレベル診断

#### Windows レジストリ確認
```powershell
# ホットキー登録の確認
Get-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced" -Name "StartMenuInit"

# グローバルホットキーの確認
Get-WinEvent -LogName Application | Where-Object {$_.ProviderName -eq "SilentCast"}
```

#### macOS システム診断
```bash
# アクセシビリティデータベースの確認
sudo sqlite3 /Library/Application\ Support/com.apple.TCC/TCC.db "SELECT * FROM access WHERE service='kTCCServiceAccessibility';"

# ホットキー登録の確認
sudo fs_usage -w -f filesys silentcast
```

#### Linux イベント監視
```bash
# キーイベントの監視
sudo evtest

# X11 イベントの確認
xev | grep KeyPress

# デバッグ用環境変数
export SILENTCAST_DEBUG_HOTKEYS=1
./silentcast --debug
```

### パフォーマンス問題

#### 高CPU使用率の原因
```bash
# ホットキー監視の CPU 使用量確認
top -p $(pgrep silentcast)

# イベントループの確認
sudo strace -p $(pgrep silentcast) | grep poll
```

#### 解決策
```yaml
# ポーリング間隔の調整
performance:
  hotkey_poll_interval: 100  # ミリ秒
  
# 不要なホットキーの削除
spells:
  # 必要最小限のスペルのみ
  e: editor
  t: terminal
```

## 📋 ホットキー診断チェックリスト

### 基本確認
- [ ] アプリケーションがスタブモードで実行されていない
- [ ] 適切な権限（アクセシビリティ等）が設定されている
- [ ] 他のアプリケーションとのキー競合がない
- [ ] 設定ファイルの構文エラーがない

### プラットフォーム固有
- [ ] **Windows**: 管理者権限または UAC 例外設定
- [ ] **macOS**: アクセシビリティ権限とターミナル許可
- [ ] **Linux**: X11/Wayland 環境とデスクトップ権限

### 設定確認
- [ ] タイムアウト値が適切（1000-3000ms）
- [ ] キーシーケンスが短い（2-3キー以下推奨）
- [ ] プレフィックスキーが他と競合していない

### デバッグ確認
- [ ] デバッグログでホットキー登録確認
- [ ] 手動テストでスペル実行確認
- [ ] システムイベント監視で入力検出確認

## 🆘 それでも解決しない場合

### 情報収集
```bash
# システム情報の収集
./silentcast --system-info > system-info.txt

# デバッグログの詳細記録
./silentcast --debug --log-file debug.log --no-tray

# 設定ダンプ
./silentcast --dump-config > config-dump.txt
```

### 代替手段
1. **手動実行モード**
   ```bash
   ./silentcast --once --spell <spell-name>
   ```

2. **スクリプトによる代替**
   ```bash
   # シェルエイリアスまたは関数として設定
   alias sce='./silentcast --once --spell e'
   ```

3. **他のホットキーツールとの連携**
   ```bash
   # AutoHotkey (Windows), Hammerspoon (macOS), xbindkeys (Linux) 等
   ```

## 🔗 関連リソース

- [インストール問題](installation.md)
- [権限設定](permissions.md)
- [設定ガイド](../guide/configuration.md)
- [プラットフォーム固有問題](platform-specific.md)