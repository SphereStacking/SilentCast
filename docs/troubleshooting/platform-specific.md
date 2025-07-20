# プラットフォーム固有問題のトラブルシューティング

SilentCastの各プラットフォーム（Windows、macOS、Linux）固有の問題と解決方法について説明します。

## 🪟 Windows 固有問題

### Windows バージョン互換性

#### 問題: Windows 10/11 での動作問題
```
[ERROR] API not supported on this Windows version
[WARN] Legacy Windows compatibility mode required
```

**解決方法:**
```powershell
# Windows バージョン確認
Get-ComputerInfo | Select-Object WindowsProductName, WindowsVersion

# 互換性モード設定
Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\AppCompatFlags\Layers" -Name "C:\path\to\silentcast.exe" -Value "WIN81"
```

#### Windows Defender と SmartScreen
```powershell
# SilentCast を Windows Defender 除外に追加
Add-MpPreference -ExclusionPath "C:\Program Files\SilentCast"
Add-MpPreference -ExclusionProcess "silentcast.exe"

# SmartScreen 警告の回避
Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer" -Name "SmartScreenEnabled" -Value "Off"

# 特定ファイルの信頼設定
Unblock-File -Path "C:\path\to\silentcast.exe"
```

### Windows レジストリ問題

#### レジストリアクセス権限
```powershell
# レジストリキー権限確認
Get-Acl -Path "HKLM:\SOFTWARE\SilentCast" | Format-Table

# 権限付与
$acl = Get-Acl "HKLM:\SOFTWARE\SilentCast"
$permission = "Users","FullControl","Allow"
$rule = New-Object System.Security.AccessControl.RegistryAccessRule $permission
$acl.SetAccessRule($rule)
$acl | Set-Acl -Path "HKLM:\SOFTWARE\SilentCast"
```

#### ホットキー登録問題
```powershell
# グローバルホットキー確認
Get-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced"

# ホットキー競合確認
tasklist /svc | findstr "hotkey"

# Windows ホットキー無効化
Set-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced" -Name "DisableHotkeys" -Value 1
```

### PowerShell 実行ポリシー

#### スクリプト実行問題
```powershell
# 現在の実行ポリシー確認
Get-ExecutionPolicy

# 実行ポリシー変更
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# 特定スクリプトのブロック解除
Unblock-File -Path "C:\path\to\script.ps1"
```

### Windows サービス統合

#### サービス登録とトラブルシューティング
```powershell
# サービス作成
sc create SilentCast binpath= "C:\SilentCast\silentcast.exe --service" start= auto displayname= "SilentCast Hotkey Manager"

# サービス状態確認
Get-Service -Name SilentCast | Format-List

# サービスログ確認
Get-WinEvent -LogName Application | Where-Object {$_.ProviderName -eq "SilentCast"}

# サービス設定修正
sc config SilentCast start= demand
sc config SilentCast depend= ""
```

## 🍎 macOS 固有問題

### macOS セキュリティとプライバシー

#### Gatekeeper 問題
```bash
# アプリの署名確認
codesign -dv --verbose=4 /Applications/SilentCast.app

# Gatekeeper バイパス
sudo xattr -rd com.apple.quarantine /Applications/SilentCast.app

# 個別アプリ許可
sudo spctl --add /Applications/SilentCast.app
sudo spctl --enable --label "SilentCast"
```

#### System Integrity Protection (SIP)
```bash
# SIP 状態確認
csrutil status

# 保護されたプロセスの確認
ps aux | grep -E "(WindowServer|loginwindow)"

# SIP 制限回避（再起動必要）
# 1. Recovery Mode で起動 (Cmd+R)
# 2. Terminal を開く
# 3. csrutil disable (非推奨)
```

### macOS アクセシビリティ問題

#### TCC (Transparency, Consent, and Control) データベース
```bash
# TCC データベース確認
sudo sqlite3 /Library/Application\ Support/com.apple.TCC/TCC.db "SELECT * FROM access WHERE service='kTCCServiceAccessibility';"

# アクセシビリティ権限リセット
sudo tccutil reset Accessibility com.silentcast.app

# 権限要求の強制
sudo tccutil reset All
```

#### AppleScript 実行問題
```bash
# AppleScript 権限確認
osascript -e 'tell application "System Events" to get name of processes'

# AppleScript セキュリティ設定
sudo sqlite3 /Library/Application\ Support/com.apple.TCC/TCC.db "SELECT * FROM access WHERE service='kTCCServiceAppleEvents';"
```

### macOS アプリケーションバンドル

#### .app バンドル構造問題
```bash
# バンドル構造確認
find /Applications/SilentCast.app -type f -name "*.plist"

# Info.plist 確認
plutil -p /Applications/SilentCast.app/Contents/Info.plist

# バンドル修復
touch /Applications/SilentCast.app
```

#### macOS Notarization
```bash
# Notarization 状態確認
spctl -a -v /Applications/SilentCast.app

# Notarization 要求確認
xcrun altool --notarization-history 0 -u developer@example.com
```

### macOS 環境変数問題

#### launchd 環境設定
```bash
# launchd 環境変数確認
launchctl print-cache

# 環境変数設定
launchctl setenv SILENTCAST_HOME /Users/user/.silentcast

# launchd plist 作成
sudo tee /Library/LaunchDaemons/com.silentcast.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.silentcast</string>
    <key>ProgramArguments</key>
    <array>
        <string>/Applications/SilentCast.app/Contents/MacOS/silentcast</string>
        <string>--daemon</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
EOF
```

## 🐧 Linux 固有問題

### ディストリビューション互換性

#### パッケージ管理システム対応
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install libx11-dev libxext-dev libxrandr-dev libxss-dev

# CentOS/RHEL/Fedora
sudo dnf install libX11-devel libXext-devel libXrandr-devel libXScrnSaver-devel

# Arch Linux
sudo pacman -S libx11 libxext libxrandr libxss

# openSUSE
sudo zypper install libX11-devel libXext-devel libXrandr-devel libXss-devel
```

#### 依存関係解決
```bash
# 実行時依存関係確認
ldd /usr/local/bin/silentcast

# 不足ライブラリ確認
LD_DEBUG=libs /usr/local/bin/silentcast 2>&1 | grep "no version"

# ライブラリパス設定
export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH
echo 'export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH' >> ~/.bashrc
```

### デスクトップ環境固有問題

#### GNOME
```bash
# GNOME Shell 拡張競合
gnome-extensions list
gnome-extensions disable example@example.com

# dconf 設定確認
dconf dump /org/gnome/settings-daemon/plugins/media-keys/

# ホットキー設定
gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings "['/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/']"

gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/ name 'SilentCast'
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/ command '/usr/local/bin/silentcast --once'
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/ binding '<Alt>space'
```

#### KDE Plasma
```bash
# KDE ショートカット管理
kwriteconfig5 --file kglobalshortcutsrc --group silentcast --key _k_friendly_name "SilentCast"
kwriteconfig5 --file kglobalshortcutsrc --group silentcast --key "activate" "Alt+Space,none,Activate SilentCast"

# KDE 設定再読み込み
qdbus org.kde.kglobalaccel /kglobalaccel reloadConfiguration
```

#### i3/sway
```bash
# i3 設定
echo "bindsym Mod1+space exec --no-startup-id silentcast --once" >> ~/.config/i3/config
i3-msg reload

# sway 設定
echo "bindsym Mod1+space exec silentcast --once" >> ~/.config/sway/config
swaymsg reload
```

### X11/Wayland 互換性

#### X11 環境
```bash
# X11 サーバー確認
echo $DISPLAY
xrandr --version

# X11 権限設定
xhost +local:
xauth list

# X11 拡張確認
xdpyinfo | grep -i "extension"
```

#### Wayland 環境
```bash
# Wayland セッション確認
echo $WAYLAND_DISPLAY
echo $XDG_SESSION_TYPE

# Wayland プロトコル確認
wayland-scanner --version

# XWayland 互換性
echo $DISPLAY  # XWayland が動作している場合設定される
```

### systemd 統合

#### systemd サービス作成
```bash
# ユーザーサービス
mkdir -p ~/.config/systemd/user
tee ~/.config/systemd/user/silentcast.service << EOF
[Unit]
Description=SilentCast Hotkey Manager
After=graphical-session.target

[Service]
Type=simple
ExecStart=/usr/local/bin/silentcast --daemon
Restart=always
RestartSec=5
Environment=DISPLAY=%i
Environment=WAYLAND_DISPLAY=wayland-0

[Install]
WantedBy=default.target
EOF

# サービス有効化
systemctl --user daemon-reload
systemctl --user enable silentcast.service
systemctl --user start silentcast.service
```

#### systemd 環境変数
```bash
# systemd 環境変数確認
systemctl --user show-environment

# 環境変数設定
systemctl --user set-environment SILENTCAST_HOME=/home/user/.silentcast

# サービスログ確認
journalctl --user -u silentcast.service -f
```

## 🔧 プラットフォーム固有設定パターン

### 条件付き設定
```yaml
# プラットフォーム固有設定
spells:
  e: editor
  
grimoire:
  editor:
    type: app
    description: "テキストエディタ"
    # Windows 設定
    windows:
      command: "C:\\Program Files\\Microsoft VS Code\\Code.exe"
      args: ["--new-window"]
    # macOS 設定  
    macos:
      command: "/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code"
      args: ["--new-window"]
    # Linux 設定
    linux:
      command: "/usr/bin/code"
      args: ["--new-window"]
```

### プラットフォーム検出
```yaml
# システム情報による動的設定
system:
  auto_detect: true
  
# 条件付きスペル
spells:
  w: windows-only
  m: macos-only
  l: linux-only
  
grimoire:
  windows-only:
    type: app
    command: "notepad.exe"
    platforms: ["windows"]
    
  macos-only:
    type: app  
    command: "open"
    args: ["-a", "TextEdit"]
    platforms: ["darwin"]
    
  linux-only:
    type: app
    command: "gedit"
    platforms: ["linux"]
```

## 📊 プラットフォーム診断スクリプト

### 統合診断スクリプト
```bash
#!/bin/bash
# プラットフォーム固有問題診断

diagnose_platform() {
    echo "🔍 プラットフォーム固有診断"
    echo "================================"
    
    # OS とバージョン確認
    case "$(uname -s)" in
        Darwin)
            echo "🍎 macOS $(sw_vers -productVersion)"
            diagnose_macos
            ;;
        Linux)
            echo "🐧 Linux $(uname -r)"
            if [ -f /etc/os-release ]; then
                . /etc/os-release
                echo "   Distribution: $NAME $VERSION"
            fi
            diagnose_linux
            ;;
        CYGWIN*|MINGW*|MSYS*)
            echo "🪟 Windows ($(uname -s))"
            diagnose_windows
            ;;
        *)
            echo "❓ Unknown OS: $(uname -s)"
            ;;
    esac
}

diagnose_macos() {
    echo "--- macOS 固有チェック ---"
    
    # SIP 状態
    echo -n "SIP Status: "
    csrutil status
    
    # アクセシビリティ権限
    echo -n "Accessibility Permission: "
    if sqlite3 /Library/Application\ Support/com.apple.TCC/TCC.db \
       "SELECT allowed FROM access WHERE service='kTCCServiceAccessibility';" 2>/dev/null | grep -q "1"; then
        echo "✅ Granted"
    else
        echo "❌ Required"
    fi
}

diagnose_linux() {
    echo "--- Linux 固有チェック ---"
    
    # デスクトップ環境
    echo "Desktop Environment: ${XDG_CURRENT_DESKTOP:-Unknown}"
    echo "Session Type: ${XDG_SESSION_TYPE:-Unknown}"
    
    # ディスプレイサーバー
    if [ -n "$WAYLAND_DISPLAY" ]; then
        echo "Display Server: Wayland ($WAYLAND_DISPLAY)"
    elif [ -n "$DISPLAY" ]; then
        echo "Display Server: X11 ($DISPLAY)"
    else
        echo "Display Server: None detected"
    fi
    
    # systemd 確認
    if systemctl --version >/dev/null 2>&1; then
        echo "Init System: systemd"
        systemctl --user is-enabled silentcast.service 2>/dev/null || echo "SilentCast service not installed"
    fi
}

diagnose_windows() {
    echo "--- Windows 固有チェック ---"
    echo "PowerShell でより詳細な診断を実行してください:"
    echo "Get-ComputerInfo | Select-Object WindowsProductName, WindowsVersion"
    echo "Get-ExecutionPolicy"
    echo "Get-Service | Where-Object {\\$_.Name -like '*silentcast*'}"
}

# 診断実行
diagnose_platform
```

## 📋 プラットフォーム固有チェックリスト

### Windows
- [ ] Windows Defender 除外設定済み
- [ ] PowerShell 実行ポリシー適切
- [ ] UAC 設定確認済み
- [ ] レジストリアクセス権限適切
- [ ] Windows サービス設定（必要な場合）

### macOS
- [ ] アクセシビリティ権限付与済み
- [ ] Gatekeeper 設定適切
- [ ] SIP 制限理解済み
- [ ] アプリバンドル構造正常
- [ ] launchd 設定（必要な場合）

### Linux
- [ ] デスクトップ環境互換性確認
- [ ] X11/Wayland 権限適切
- [ ] 必要な依存ライブラリインストール済み
- [ ] systemd サービス設定（必要な場合）
- [ ] ディストリビューション固有設定確認

## 🔗 関連リソース

- [権限設定](permissions.md)
- [ホットキー問題](hotkeys.md)
- [インストール問題](installation.md)
- [デバッグガイド](debugging.md)