# インストール問題のトラブルシューティング

SilentCastのインストールで発生する一般的な問題と解決方法について説明します。

## 🚀 クイック解決法

### 最も一般的な問題
1. **CGOビルドエラー** → スタブモードを使用
2. **依存関係エラー** → 必要なライブラリをインストール
3. **権限エラー** → 管理者権限で実行
4. **パスエラー** → 絶対パスを使用

## 🛠️ CGOビルド問題

### 問題: CGO関連のビルドエラー
```
# github.com/robotn/gohook
fatal error: X11/keysym.h: No such file or directory
```

### 解決方法

#### 1. スタブモードの使用（推奨）
```bash
# スタブバージョンをビルド（ホットキー機能なし）
make build-stub

# 実行
./build/silentcast --no-tray
```

#### 2. 依存関係のインストール

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install -y \
    gcc \
    libc6-dev \
    libx11-dev \
    libxext-dev \
    libxrandr-dev \
    libxss-dev \
    libgconf-2-4 \
    libxss1 \
    libappindicator3-dev
```

**CentOS/RHEL/Fedora:**
```bash
sudo yum install -y \
    gcc \
    glibc-devel \
    libX11-devel \
    libXext-devel \
    libXrandr-devel \
    libXss-devel
```

**macOS:**
```bash
# Xcode Command Line Tools
xcode-select --install

# または Homebrew で
brew install gcc
```

**Windows:**
```bash
# Chocolatey を使用
choco install mingw

# または MSYS2 を使用
pacman -S mingw-w64-x86_64-gcc
```

#### 3. CGO環境変数の設定
```bash
export CGO_ENABLED=1
export CC=gcc

# Windowsの場合
set CGO_ENABLED=1
set CC=gcc
```

## 📦 依存関係問題

### Go モジュール問題

#### 問題: `go.sum` の不整合
```bash
go mod tidy
go mod download
```

#### 問題: プライベートモジュールアクセス
```bash
export GOPRIVATE=github.com/your-org/*
export GONOPROXY=github.com/your-org/*
export GONOSUMDB=github.com/your-org/*
```

### システム依存関係

#### Linux デスクトップ環境
```bash
# GNOME
sudo apt install -y libayatana-appindicator3-dev

# KDE
sudo apt install -y libappindicator3-dev

# Notification support
sudo apt install -y libnotify-bin zenity
```

#### macOS アプリケーション権限
```bash
# Accessibility permissions required for hotkeys
# System Preferences > Security & Privacy > Privacy > Accessibility
```

#### Windows システム要件
- Windows 10 以降推奨
- PowerShell 5.1 以降
- .NET Framework 4.8 以降（通知機能用）

## 🔐 権限問題

### Linux 権限

#### ファイル権限
```bash
# 実行権限の付与
chmod +x silentcast

# 設定ディレクトリの権限
mkdir -p ~/.config/silentcast
chmod 755 ~/.config/silentcast
```

#### X11/Wayland 権限
```bash
# X11セッション確認
echo $DISPLAY

# Waylandセッション確認
echo $WAYLAND_DISPLAY

# 必要に応じてXvfbを使用
Xvfb :99 -screen 0 1024x768x24 &
export DISPLAY=:99
```

### macOS 権限

#### アクセシビリティ権限
1. システム環境設定を開く
2. セキュリティとプライバシー > プライバシー
3. アクセシビリティを選択
4. SilentCastを追加し、チェックを入れる

#### Gatekeeper 問題
```bash
# 署名されていないアプリケーションの実行許可
sudo spctl --master-disable

# 特定のアプリケーションの許可
sudo spctl --add /path/to/silentcast
```

### Windows 権限

#### UAC問題
```powershell
# 管理者として実行
# PowerShell を管理者として開く
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

#### Windows Defender除外
```powershell
# Windows Defenderの除外設定
Add-MpPreference -ExclusionPath "C:\path\to\silentcast"
```

## 📍 パス問題

### 環境変数の設定

#### Go環境
```bash
# Go パスの確認
go env GOPATH
go env GOROOT

# パスの設定
export PATH=$PATH:$(go env GOPATH)/bin
```

#### システムパス
```bash
# Linux/macOS
export PATH=$PATH:/usr/local/bin:/opt/bin

# Windows
set PATH=%PATH%;C:\Go\bin;C:\tools
```

### 設定ファイルパス

#### デフォルト設定パス
```bash
# Linux
~/.config/silentcast/spellbook.yml

# macOS
~/Library/Application Support/silentcast/spellbook.yml

# Windows
%APPDATA%\silentcast\spellbook.yml
```

#### カスタム設定パス
```bash
# コマンドライン指定
./silentcast --config /custom/path/to/config

# 環境変数
export SILENTCAST_CONFIG_DIR=/custom/path
```

## 🔄 再インストール手順

### 完全なクリーンアップ
```bash
# ビルド成果物の削除
make clean

# Go モジュールキャッシュのクリア
go clean -modcache

# 設定ファイルの削除
rm -rf ~/.config/silentcast

# Windows の場合
rmdir /s %APPDATA%\silentcast
```

### 新規インストール
```bash
# リポジトリのクローン
git clone https://github.com/SphereStacking/SilentCast.git
cd SilentCast

# 依存関係のインストール
cd app
go mod download

# スタブモードでのビルド
make build-stub

# 動作確認
./build/silentcast --version
```

## 🧪 インストール検証

### 基本動作テスト
```bash
# バージョン確認
./silentcast --version

# 設定検証
./silentcast --validate-config

# デバッグモード起動（10秒後に停止）
timeout 10s ./silentcast --debug --no-tray
```

### システム要件確認
```bash
# Go バージョン
go version

# システム情報
uname -a

# 利用可能メモリ
free -h  # Linux
vm_stat  # macOS
systeminfo | find "Available Physical Memory"  # Windows
```

## 📋 インストール診断チェックリスト

- [ ] Go 1.19+ がインストールされている
- [ ] 必要なシステム依存関係がインストールされている
- [ ] CGO が有効（フル機能の場合）
- [ ] 適切なファイル権限が設定されている
- [ ] 環境変数（PATH、GOPATH）が正しく設定されている
- [ ] ファイアウォール/セキュリティソフトが干渉していない
- [ ] スタブモードで基本動作が確認できている

## 🆘 それでも解決しない場合

### 情報収集
```bash
# システム情報の収集
./silentcast --system-info > system-info.txt

# ビルドログの保存
make build-stub 2>&1 | tee build.log

# 環境変数の確認
env | grep -E "(GO|PATH|CGO)" > env.txt
```

### サポート要請時の情報
1. OS とバージョン
2. Go バージョン
3. エラーメッセージの全文
4. インストール方法（ソースビルド/バイナリ）
5. 試行した解決方法

## 🔗 関連リソース

- [システム要件](../guide/installation.md#system-requirements)
- [ビルド手順](../api/building.md)
- [設定ガイド](../guide/configuration.md)
- [開発環境構築](../development/setup.md)