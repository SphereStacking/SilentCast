# SilentCast 🤫⚡

[![CI](https://github.com/SphereStacking/silentcast/actions/workflows/ci.yml/badge.svg)](https://github.com/SphereStacking/silentcast/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/SphereStacking/silentcast)](https://github.com/SphereStacking/silentcast/releases)
[![License](https://img.shields.io/github/license/SphereStacking/silentcast)](LICENSE)

[English](README.md) | [日本語](README.ja.md)

SilentCast は、シンプルなキーボードショートカットでタスクを実行できる、サイレントなホットキー駆動タスクランナーです。キーボードから手を離すことなくワークフローを効率化したい開発者に最適です。

## ✨ 機能

- 🎯 **グローバルホットキー** - システム全体で動作、アプリケーションをフォーカスする必要なし
- 🏃 **高速実行** - アプリケーション起動とスクリプト実行が瞬時に
- 📝 **シーケンシャルキー** - VS Code スタイルのキーコンビネーション（例：`g,s` で git status）
- 🌐 **クロスプラットフォーム** - Windows と macOS で動作
- 🎨 **直感的** - spells と actions によるシンプルな設定
- 🧪 **軽量** - 最小限の CPU とメモリ使用量（約15MB）
- 🔄 **自動リロード** - 設定変更が自動的に適用される
- 📋 **システムトレイ** - 邪魔にならないシステムトレイ統合
- 📊 **構造化ログ** - ローテーション対応の包括的なログ

## 📦 インストール

### バイナリをダウンロード

[リリースページ](https://github.com/SphereStacking/silentcast/releases)から最新版をダウンロードしてください。

#### macOS
```bash
# Intel
curl -L https://github.com/SphereStacking/silentcast/releases/latest/download/silentcast-darwin-amd64.tar.gz | tar xz
sudo mv silentcast/silentcast /usr/local/bin/

# Apple Silicon
curl -L https://github.com/SphereStacking/silentcast/releases/latest/download/silentcast-darwin-arm64.tar.gz | tar xz
sudo mv silentcast/silentcast /usr/local/bin/
```

#### Windows

リリースページから ZIP ファイルをダウンロードし、PATH の通ったディレクトリに展開してください。

### ソースからビルド

```bash
# リポジトリをクローン
git clone https://github.com/SphereStacking/silentcast.git
cd silentcast

# ビルドしてインストール
make install

# または現在のプラットフォーム用にビルドのみ
make build
```

## 🔮 設定

SilentCast は YAML 設定ファイルを使用します：

- **Spells（呪文）** - キーボードショートカット
- **Grimoire（魔導書）** - 実行するコマンド/アプリケーション
- **Command（コマンド）** - 実行するコマンドまたはパス

### 基本設定

`spellbook.yml` ファイルを作成：

```yaml
# SilentCast 設定
daemon:
  auto_start: false
  log_level: info
  config_watch: true

logger:
  level: info
  file: ""                 # 空 = コンソールのみ
  max_size: 10             # MB
  max_backups: 3
  max_age: 7               # 日数
  compress: false

hotkeys:
  prefix: "alt+space"      # 魔法のキー
  timeout: 1000            # プレフィックス後の待機時間（ミリ秒）
  sequence_timeout: 2000   # 完全なシーケンスのタイムアウト（ミリ秒）

spells:
  # 単一キーの呪文
  e: "editor"
  t: "terminal"
  b: "browser"
  
  # マルチキーシーケンス（VS Code スタイル）
  "g,s": "git_status"
  "g,p": "git_pull"
  "g,c": "git_commit"

grimoire:
  editor:
    type: app
    command: "code"    # VS Code
    description: "VS Code を開く"
  
  terminal:
    type: app
    command: "wt"      # Windows Terminal
    description: "ターミナルを開く"
  
  git_status:
    type: script
    command: "git status"
    description: "git status を表示"
```

### プラットフォーム別設定

SilentCast はプラットフォーム固有のオーバーライドをサポートします：

- `spellbook.yml` - 基本設定（最初に読み込まれる）
- `spellbook.mac.yml` - macOS オーバーライド
- `spellbook.windows.yml` - Windows オーバーライド

`spellbook.mac.yml` の例：
```yaml
grimoire:
  terminal:
    type: app
    command: "Terminal"
  
  browser:
    type: app
    command: "Safari"
```

### 設定例

`examples/config/` ディレクトリに設定例があります：
- `spellbook.yml` - よく使うショートカットを含む完全な例
- `spellbook.windows.yml` - Windows 固有のオーバーライド
- `spellbook.mac.yml` - macOS 固有のオーバーライド

詳細な設定ガイドは [CONFIG.md](CONFIG.md) を参照してください。

## 🎮 使い方

### SilentCast の起動

```bash
# デフォルト設定で実行
silentcast

# システムトレイなしで実行
silentcast --no-tray

# カスタム設定場所で実行
SILENTCAST_CONFIG=/path/to/config silentcast
```

### 呪文の詠唱

1. プレフィックスキーを押す（デフォルト：`Alt+Space`）
2. 呪文を詠唱：
   - **単一キー**：`e` を押してエディタを開く
   - **シーケンス**：`g` を押してから `s` で git status
   - **長いシーケンス**：`d`、`o`、`c` と順に押してドキュメントを表示

### システムトレイ

システムトレイサポートで実行時：
- **Show Hotkeys** - 設定されたすべてのショートカットを表示
- **Reload Config** - 設定を手動でリロード
- **About** - バージョン情報を表示
- **Quit** - SilentCast を終了

## 💻 プラットフォームサポート

### Windows
- ✅ 完全なグローバルホットキーサポート
- ✅ システムトレイ統合
- ✅ 特別な権限不要

### macOS  
- ✅ 完全なグローバルホットキーサポート
- ⚠️ アクセシビリティ権限が必要
- ✅ システムトレイ統合
- 📝 初回実行時：システム環境設定 > セキュリティとプライバシー > アクセシビリティで権限を付与

## 🔧 開発

### 前提条件

- Go 1.21 以降
- Make（オプションですが推奨）

## 📁 プロジェクト構造

```
SilentCast/
├── app/              # アプリケーションソースコード
│   ├── cmd/          # メインエントリポイント
│   ├── internal/     # 内部パッケージ
│   ├── pkg/          # 公開パッケージ
│   └── Makefile      # ビルド設定
├── docs/             # VitePress ドキュメント
│   ├── guide/        # ユーザーガイド
│   ├── config/       # 設定リファレンス
│   └── api/          # 開発者ドキュメント
├── examples/         # 設定例
└── README.ja.md      # このファイル
```

### クイックスタート

```bash
# リポジトリをクローン
git clone https://github.com/SphereStacking/silentcast.git
cd silentcast

# 開発用にビルド
make build-stub

# 実行
./app/build/silentcast --no-tray
```

### 完全なドキュメント

```bash
# ドキュメントサーバーを起動
make docs-dev
# http://localhost:5173 を開く
```

詳細なビルド手順は [docs/api/build.md](docs/api/build.md) を参照してください。

### 利用可能なコマンド

```bash
# アプリケーション
make build         # アプリケーションをビルド
make build-stub    # C 依存関係なしでビルド
make test          # テストを実行
make clean         # ビルド成果物をクリーン

# ドキュメント
make docs-dev      # ドキュメント開発サーバーを起動
make docs-build    # ドキュメントをビルド

# セットアップ
make setup         # 開発環境をセットアップ
```

### テスト

```bash
# すべてのテストを実行
make test

# カバレッジ付きでテストを実行
make test-coverage

# 特定のパッケージのテストを実行
go test -v ./internal/config/
```

### コントリビューション

1. リポジトリをフォーク
2. フィーチャーブランチを作成（`git checkout -b feature/amazing-spell`）
3. 変更をコミット（`git commit -m 'Add amazing spell'`）
4. ブランチにプッシュ（`git push origin feature/amazing-spell`）
5. プルリクエストを作成

## 📄 ライセンス

MIT ライセンス - 詳細は [LICENSE](LICENSE) ファイルを参照してください。

## 🙏 謝辞

- [gohook](https://github.com/robotn/gohook) - グローバルホットキーサポート
- [systray](https://github.com/getlantern/systray) - システムトレイ統合
- [fsnotify](https://github.com/fsnotify/fsnotify) - ファイルシステム通知
- [lumberjack](https://github.com/natefinch/lumberjack) - ログローテーション

---

<p align="center">
  キーボードショートカットを ❤️ する開発者たちが 🪄 で作りました
</p>