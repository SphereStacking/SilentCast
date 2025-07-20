# パーミッション問題のトラブルシューティング

SilentCastの権限関連問題の診断と解決方法について説明します。

## 🚀 クイック診断

### 権限エラーが発生する場合

```bash
# 1. 現在の権限状態確認
./silentcast --check-permissions

# 2. 権限要求プロセス実行
./silentcast --request-permissions

# 3. デバッグモードで詳細確認
./silentcast --debug --no-tray
```

### 一般的な権限問題
1. **アクセシビリティ権限なし** → macOS固有、手動設定必要
2. **管理者権限不足** → Windows UAC、昇格アクション実行時
3. **X11/Wayland権限なし** → Linux環境、ディスプレイアクセス
4. **ファイルシステム権限** → 設定ファイル読み書き、スクリプト実行

## 🍎 macOS パーミッション

### 問題: アクセシビリティ権限が必要

#### 症状
```
[ERROR] Permission denied: accessibility features required
[WARN] Hotkey registration failed: insufficient permissions
[ERROR] Failed to monitor keyboard events
```

#### 解決方法: アクセシビリティ権限設定

**1. システム環境設定から設定**
```bash
# 1. システム環境設定を開く
open "x-apple.systempreferences:com.apple.preference.security?Privacy_Accessibility"

# 2. または手動で:
# システム環境設定 > セキュリティとプライバシー > プライバシー > アクセシビリティ
```

**2. SilentCast を追加**
1. 🔒 をクリックして管理者パスワード入力
2. **+** ボタンで SilentCast バイナリを追加
3. チェックボックスにチェックを入れる
4. アプリケーションを再起動

**3. ターミナルにも権限付与（必要に応じて）**
```bash
# ターミナルからSilentCastを実行する場合
# ターミナルもアクセシビリティリストに追加
```

#### 自動化スクリプト
```bash
#!/bin/bash
# macOS アクセシビリティ権限チェックスクリプト

check_accessibility_permission() {
    # TCC データベースを確認
    local app_bundle="/Applications/SilentCast.app"
    local result=$(sqlite3 /Library/Application\ Support/com.apple.TCC/TCC.db \
        "SELECT allowed FROM access WHERE service='kTCCServiceAccessibility' AND client='$app_bundle';" 2>/dev/null)
    
    if [ "$result" = "1" ]; then
        echo "✅ アクセシビリティ権限: 付与済み"
        return 0
    else
        echo "❌ アクセシビリティ権限: 未付与"
        return 1
    fi
}

request_accessibility_permission() {
    echo "🔧 アクセシビリティ権限の設定手順:"
    echo "1. システム環境設定 > セキュリティとプライバシー > プライバシー"
    echo "2. アクセシビリティを選択"
    echo "3. 🔒をクリックして管理者パスワード入力"
    echo "4. + ボタンでSilentCastを追加"
    echo "5. チェックボックスにチェック"
    echo "6. アプリケーション再起動"
    
    # システム環境設定を直接開く
    open "x-apple.systempreferences:com.apple.preference.security?Privacy_Accessibility"
}

if ! check_accessibility_permission; then
    request_accessibility_permission
fi
```

### 問題: フルディスクアクセス権限

#### 症状
```
[ERROR] Permission denied reading configuration file
[WARN] Cannot access user home directory
```

#### 解決方法
```bash
# フルディスクアクセス権限設定
# システム環境設定 > セキュリティとプライバシー > プライバシー > フルディスクアクセス
# SilentCast とターミナルを追加
```

### macOS 特有の制限

#### System Integrity Protection (SIP)
```bash
# SIP 状態確認
csrutil status

# 保護されたディレクトリ
echo "SIPにより保護される場所:"
echo "- /System/"
echo "- /usr/ (一部)"
echo "- /bin/"
echo "- /sbin/"
echo "- その他システムディレクトリ"
```

#### Gatekeeper 問題
```bash
# 署名されていないアプリの実行許可
sudo spctl --master-disable

# 個別アプリの許可
sudo xattr -rd com.apple.quarantine /path/to/silentcast

# Gatekeeper 状態確認
spctl --status
```

## 🪟 Windows パーミッション

### 問題: UAC (User Account Control) 制限

#### 症状
```
[ERROR] Access denied: administrative privileges required
[WARN] Elevated action execution failed
[ERROR] Registry access denied
```

#### 解決方法: 管理者として実行

**1. 右クリックで管理者実行**
```powershell
# PowerShell から管理者実行
Start-Process ./silentcast -Verb RunAs

# または cmd から
runas /user:Administrator silentcast.exe
```

**2. UAC 設定の調整**
```powershell
# 現在のUAC レベル確認
Get-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" -Name "ConsentPromptBehaviorAdmin"

# UAC レベル変更（注意: セキュリティリスク）
Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" -Name "ConsentPromptBehaviorAdmin" -Value 0
```

**3. UAC 例外設定**
```powershell
# SilentCast を UAC 例外に追加
New-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" -Name "SilentCastBypass" -Value 1 -PropertyType DWORD

# タスクスケジューラで昇格タスク作成
schtasks /create /tn "SilentCast" /tr "C:\path\to\silentcast.exe" /rl HIGHEST /f
```

### 問題: Windows Defender / アンチウイルス干渉

#### 症状
```
[WARN] Hotkey registration blocked by security software
[ERROR] Process creation blocked
```

#### 解決方法
```powershell
# Windows Defender 除外設定
Add-MpPreference -ExclusionPath "C:\path\to\silentcast"
Add-MpPreference -ExclusionProcess "silentcast.exe"

# 除外設定確認
Get-MpPreference | Select-Object -Property ExclusionPath, ExclusionProcess
```

### Windows サービスとして実行

#### サービス登録
```powershell
# サービス作成
sc create SilentCast binpath= "C:\path\to\silentcast.exe --service" start= auto

# サービス開始
sc start SilentCast

# サービス状態確認
sc query SilentCast
```

#### サービス用設定
```yaml
# spellbook.yml - サービス実行用設定
daemon:
  service_mode: true
  auto_start: true
  log_level: info
  
logger:
  file: "C:\\ProgramData\\SilentCast\\silentcast.log"
  
performance:
  background_mode: true
```

## 🐧 Linux パーミッション

### 問題: X11/Wayland 権限不足

#### 症状
```
[ERROR] Cannot connect to X server
[WARN] Wayland compositor access denied
[ERROR] Input device access permission denied
```

#### X11 環境での解決方法
```bash
# DISPLAY 環境変数確認
echo $DISPLAY

# X11 権限設定
xhost +local:

# 特定ユーザーにアクセス許可
xhost +si:localuser:$(whoami)

# X11 認証情報確認
echo $XAUTHORITY
xauth list
```

#### Wayland 環境での解決方法
```bash
# Wayland セッション確認
echo $WAYLAND_DISPLAY

# 必要な環境変数設定
export WAYLAND_DISPLAY=wayland-0
export XDG_RUNTIME_DIR=/run/user/$(id -u)

# Wayland 権限確認
ls -la $XDG_RUNTIME_DIR/wayland-*
```

### 問題: デスクトップ環境権限

#### GNOME
```bash
# GNOME 設定ツール
gsettings list-schemas | grep -i hotkey
gsettings list-schemas | grep -i keyboard

# ホットキー設定確認
gsettings get org.gnome.desktop.wm.keybindings switch-applications

# カスタムショートカット設定
gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings "['/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/']"
```

#### KDE
```bash
# KDE ショートカット設定
kwriteconfig5 --file kglobalshortcutsrc --group silentcast --key _k_friendly_name SilentCast

# ショートカット確認
kreadconfig5 --file kglobalshortcutsrc --group silentcast
```

#### i3/sway
```bash
# i3 設定例
echo "bindsym Mod1+space exec silentcast --once" >> ~/.config/i3/config

# sway 設定例
echo "bindsym Mod1+space exec silentcast --once" >> ~/.config/sway/config
```

### 問題: デバイスファイル権限

#### 入力デバイスアクセス
```bash
# 入力デバイス確認
ls -la /dev/input/

# 必要なグループ確認
groups $(whoami)

# input グループに追加
sudo usermod -a -G input $(whoami)

# udev ルール作成
sudo tee /etc/udev/rules.d/99-silentcast.rules << EOF
SUBSYSTEM=="input", GROUP="input", MODE="0664"
KERNEL=="event*", GROUP="input", MODE="0664"
EOF

# udev ルール再読み込み
sudo udevadm control --reload-rules
sudo udevadm trigger
```

## 🔐 昇格アクション (Elevated Actions)

### 設定例
```yaml
# 昇格が必要なアクション
grimoire:
  system-update:
    type: elevated
    command: apt update && apt upgrade
    description: "システムアップデート"
    require_confirmation: true
    
  service-restart:
    type: elevated  
    command: systemctl restart nginx
    description: "Nginx再起動"
    
  log-cleanup:
    type: elevated
    command: find /var/log -name "*.log" -mtime +30 -delete
    description: "古いログファイル削除"
```

### 昇格方式の設定
```yaml
# プラットフォーム別昇格設定
elevation:
  windows:
    method: "uac"           # UAC プロンプト
    runas_user: ""          # 空の場合は現在ユーザーで昇格
    
  macos:
    method: "sudo"          # sudo 使用
    sudo_prompt: true       # パスワードプロンプト表示
    
  linux:
    method: "pkexec"        # PolicyKit 使用
    fallback: "sudo"        # フォールバック方式
```

## 🛠️ 権限診断ツール

### 包括的権限チェックスクリプト
```bash
#!/bin/bash
# SilentCast 権限診断スクリプト

check_permissions() {
    echo "🔍 SilentCast 権限診断"
    echo "========================"
    
    # OS 検出
    case "$(uname -s)" in
        Darwin)
            check_macos_permissions
            ;;
        Linux)
            check_linux_permissions
            ;;
        CYGWIN*|MINGW*|MSYS*)
            check_windows_permissions
            ;;
        *)
            echo "❓ 不明なOS: $(uname -s)"
            ;;
    esac
}

check_macos_permissions() {
    echo "🍎 macOS 権限チェック"
    
    # アクセシビリティ権限
    if sqlite3 /Library/Application\ Support/com.apple.TCC/TCC.db \
       "SELECT allowed FROM access WHERE service='kTCCServiceAccessibility';" 2>/dev/null | grep -q "1"; then
        echo "✅ アクセシビリティ権限: OK"
    else
        echo "❌ アクセシビリティ権限: 必要"
        echo "   → システム環境設定 > セキュリティとプライバシー > アクセシビリティ"
    fi
}

check_linux_permissions() {
    echo "🐧 Linux 権限チェック"
    
    # X11/Wayland チェック
    if [ -n "$DISPLAY" ]; then
        echo "✅ X11 セッション: $DISPLAY"
    elif [ -n "$WAYLAND_DISPLAY" ]; then
        echo "✅ Wayland セッション: $WAYLAND_DISPLAY"
    else
        echo "❌ ディスプレイセッション: 検出されず"
    fi
    
    # 入力デバイス権限
    if groups | grep -q input; then
        echo "✅ input グループ: 所属済み"
    else
        echo "❌ input グループ: 未所属"
        echo "   → sudo usermod -a -G input \$(whoami)"
    fi
}

check_windows_permissions() {
    echo "🪟 Windows 権限チェック"
    
    # 管理者権限チェック (PowerShell が必要)
    echo "管理者権限の確認が必要です"
    echo "PowerShell で以下を実行:"
    echo '([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)'
}

# 実行
check_permissions
```

## 📋 権限問題診断チェックリスト

### 基本確認
- [ ] OS 固有の権限要件を理解している
- [ ] 現在のユーザー権限を確認済み
- [ ] 必要な権限が適切に設定されている
- [ ] セキュリティソフトウェアの干渉を確認済み

### プラットフォーム固有
- [ ] **macOS**: アクセシビリティ権限とフルディスクアクセス
- [ ] **Windows**: UAC 設定と管理者権限
- [ ] **Linux**: X11/Wayland 権限と input グループ

### 昇格アクション
- [ ] 昇格方式が適切に設定されている
- [ ] 昇格確認プロンプトが機能する
- [ ] セキュリティ警告が適切に表示される

### デバッグ確認
- [ ] 権限エラーがログに記録されている
- [ ] 診断ツールで権限状態確認済み
- [ ] 権限要求プロセスが正常動作する

## 🆘 それでも解決しない場合

### 代替手段
1. **権限なしモード**
   ```bash
   ./silentcast --no-hotkeys --manual-mode
   ```

2. **スクリプト経由実行**
   ```bash
   # 権限付きシェルスクリプト作成
   echo '#!/bin/bash' > run_silentcast.sh
   echo 'exec /path/to/silentcast "$@"' >> run_silentcast.sh
   chmod +x run_silentcast.sh
   ```

3. **コンテナ実行**
   ```bash
   # Docker で分離実行
   docker run --privileged -v /tmp/.X11-unix:/tmp/.X11-unix silentcast
   ```

## 🔗 関連リソース

- [インストール問題](installation.md)
- [ホットキー問題](hotkeys.md)
- [プラットフォーム固有問題](platform-specific.md)
- [デバッグガイド](debugging.md)