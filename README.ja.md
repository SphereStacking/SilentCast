# SilentCast

<div align="center">
  <img src="https://spherestacking.github.io/SilentCast/logo.svg" alt="SilentCast Logo" width="200" height="200">
  
  <h3>🪄 呪文を唱えて、タスクを実行</h3>
  
  <p>シンプルなキーボードの呪文でタスクを実行する、サイレントなホットキー駆動タスクランナー</p>
</div>

<p align="center">
  <a href="https://github.com/SphereStacking/silentcast/actions/workflows/ci.yml"><img src="https://github.com/SphereStacking/silentcast/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://github.com/SphereStacking/silentcast/releases"><img src="https://img.shields.io/github/v/release/SphereStacking/silentcast" alt="Release"></a>
  <a href="https://goreportcard.com/report/github.com/SphereStacking/silentcast"><img src="https://goreportcard.com/badge/github.com/SphereStacking/silentcast" alt="Go Report Card"></a>
  <a href="LICENSE"><img src="https://img.shields.io/github/license/SphereStacking/silentcast" alt="License"></a>
  <a href="https://pkg.go.dev/github.com/SphereStacking/silentcast"><img src="https://pkg.go.dev/badge/github.com/SphereStacking/silentcast.svg" alt="Go Reference"></a>
</p>

<p align="center">
  <a href="README.md">English</a> | <a href="README.ja.md">日本語</a> | <a href="https://spherestacking.github.io/SilentCast/">ドキュメント</a>
</p>

---

## 🌟 SilentCastとは？

SilentCastは、バックグラウンドで静かに動作し、キーボードの呪文を待ち受けて事前定義されたタスクを実行する、軽量でクロスプラットフォームなアプリケーションです。開発者、システム管理者、パワーユーザーのいずれであっても、SilentCastはシンプルなキーボードショートカットで反復的なタスクを自動化します。

### ✨ 主な機能

- **🎯 グローバルホットキー** - どこでも動作、ウィンドウフォーカス不要
- **⚡ 超高速** - 最小限のリソース使用で瞬時にタスク実行
- **🔮 魔法の用語** - spells（呪文）とgrimoire（魔導書）を使用
- **🎹 VS Codeスタイルのシーケンス** - `g,s`でgit statusのような複数キーの組み合わせ
- **🌍 クロスプラットフォーム** - Windows、macOS、Linuxにネイティブ対応
- **🔄 ライブ設定** - 再起動なしで変更が即座に適用
- **📊 スマート出力** - コマンド結果を通知またはターミナルに表示
- **🔐 管理者権限実行** - 必要に応じて管理者権限でタスクを実行
- **🧪 開発者フレンドリー** - テストとデバッグのための包括的なCLIツール

## 🚀 クイックスタート

### インストール

#### パッケージマネージャーを使用

```bash
# macOS (Homebrew)
brew install spherestacking/tap/silentcast

# Windows (Scoop)
scoop bucket add spherestacking https://github.com/spherestacking/scoop-bucket
scoop install silentcast

# Linux (Snap)
sudo snap install silentcast
```

#### 直接ダウンロード

[リリースページ](https://github.com/SphereStacking/silentcast/releases)から最新のバイナリをダウンロードしてください。

```bash
# macOS (Apple Silicon)の例
curl -L https://github.com/SphereStacking/silentcast/releases/latest/download/silentcast-darwin-arm64.tar.gz | tar xz
sudo mv silentcast /usr/local/bin/
```

### 最初の呪文

1. 設定ファイル `spellbook.yml` を作成:

```yaml
# 基本的なspellbook設定
hotkeys:
  prefix: "alt+space"      # アクティベーションキー

spells:
  e: "editor"              # Alt+Space、その後 E
  t: "terminal"            # Alt+Space、その後 T
  "g,s": "git_status"      # Alt+Space、その後 G、その後 S

grimoire:
  editor:
    type: app
    command: "code"        # VS Codeを開く
    
  terminal:
    type: app
    command: "wt"          # Windows Terminalを開く
    
  git_status:
    type: script
    command: "git status"
    show_output: true      # 結果を通知に表示
```

2. SilentCastを起動:

```bash
silentcast
```

3. 最初の呪文を唱える:
   - `Alt+Space`（プレフィックスキー）を押す
   - `e`を押してエディタを開く
   - または`g`、その後`s`を押してgit statusを確認

## 🎮 使用例

### 基本コマンド

```bash
# SilentCastを起動
silentcast                           # システムトレイで実行
silentcast --no-tray                 # システムトレイなしで実行
silentcast --debug                   # デバッグログを有効化

# 設定管理
silentcast --validate-config         # 設定構文をチェック
silentcast --show-config             # マージされた設定を表示
silentcast --list-spells             # 利用可能なすべての呪文を表示

# テストとデバッグ
silentcast --test-spell --spell=e    # 特定の呪文をテスト
silentcast --dry-run --spell=g,s     # 実行せずにプレビュー
silentcast --test-hotkey             # ホットキー検出をテスト

# 一回限りの実行
silentcast --once --spell=e          # 呪文を実行して終了
```

### 高度な設定

```yaml
grimoire:
  # コマンド出力を表示
  docker_ps:
    type: script
    command: "docker ps"
    show_output: true
    description: "Dockerコンテナ一覧"
    
  # 実行後もターミナルを開いたままにする
  python_shell:
    type: script
    command: "python"
    terminal: true
    keep_open: true
    description: "対話型Pythonシェル"
    
  # 管理者権限で実行
  system_update:
    type: script
    command: "apt update && apt upgrade -y"
    admin: true
    terminal: true
    description: "システムパッケージを更新"
    
  # URLを開く
  github_profile:
    type: url
    command: "https://github.com/{{.Username}}"
    description: "GitHubプロフィールを開く"
    
  # カスタムシェルとタイムアウト
  long_process:
    type: script
    command: "./backup.sh"
    shell: "bash"
    timeout: 300
    show_output: true
    description: "5分のタイムアウトでバックアップを実行"
```

## 📚 ドキュメント

### ユーザーガイド
- [はじめに](https://spherestacking.github.io/SilentCast/guide/getting-started)
- [設定ガイド](https://spherestacking.github.io/SilentCast/guide/configuration)
- [呪文と魔導書](https://spherestacking.github.io/SilentCast/guide/spells)
- [プラットフォーム設定](https://spherestacking.github.io/SilentCast/guide/platforms)

### リファレンス
- [CLIリファレンス](https://spherestacking.github.io/SilentCast/guide/cli-reference)
- [設定スキーマ](https://spherestacking.github.io/SilentCast/config/)
- [トラブルシューティング](https://spherestacking.github.io/SilentCast/troubleshooting/)

### 開発者リソース
- [APIドキュメント](https://pkg.go.dev/github.com/SphereStacking/silentcast)
- [アーキテクチャガイド](https://spherestacking.github.io/SilentCast/api/architecture)
- [コントリビュート](https://spherestacking.github.io/SilentCast/contributing)

## 💻 プラットフォームサポート

| プラットフォーム | グローバルホットキー | システムトレイ | 通知 | 管理者/Sudo | 自動起動 |
|------------------|----------------------|----------------|------|-------------|----------|
| Windows          | ✅ | ✅ | ✅ (ネイティブ) | ✅ | ✅ |
| macOS            | ✅* | ✅ | ✅ (ネイティブ) | ✅ | ✅ |
| Linux            | ✅ | ✅** | ✅ (複数) | ✅ | ✅ |

\* macOSでは初回実行時にアクセシビリティ権限が必要です  
\** Linuxではシステムトレイに`libappindicator3-1`が必要です

## 🔧 開発

### 前提条件
- Go 1.23以降
- Make（オプションですが推奨）
- Cコンパイラ（ホットキーサポート付きの本番ビルド用）

### ソースからのビルド

```bash
# リポジトリをクローン
git clone https://github.com/SphereStacking/silentcast.git
cd silentcast

# 開発環境をセットアップ
make setup

# ビルドオプション
make build-dev      # ホットキーサポートなしの高速ビルド（開発用）
make build          # フル機能の本番ビルド
make build-all      # すべてのプラットフォーム用にビルド

# テストを実行
make test           # ユニットテスト
make test-all       # 統合テストを含むすべてのテスト
make benchmark      # パフォーマンスベンチマーク

# 開発ワークフロー
make lint           # リンターを実行
make fmt            # コードをフォーマット
make docs-dev       # ドキュメントサーバーを起動
```

### プロジェクト構造

```
SilentCast/
├── app/                    # アプリケーションソースコード
│   ├── cmd/                # メインエントリポイント
│   ├── internal/           # 内部パッケージ
│   │   ├── action/         # アクション実行
│   │   ├── config/         # 設定管理
│   │   ├── hotkey/         # ホットキー検出
│   │   └── notify/         # 通知システム
│   └── pkg/                # 公開パッケージ
├── docs/                   # ドキュメント（VitePress）
├── examples/               # 設定例
└── .ticket/                # チケットベースの開発システム
```

## 🤝 コントリビュート

コントリビュートを歓迎します！詳細は[コントリビューションガイド](CONTRIBUTING.md)をご覧ください。

### クイックコントリビューションガイド

1. 既存の[issues](https://github.com/SphereStacking/silentcast/issues)と[tickets](.ticket/README.md)を確認
2. リポジトリをフォーク
3. フィーチャーブランチを作成（`git checkout -b feature/amazing-spell`）
4. コーディング標準に従い、魔法の用語を使用
5. 新機能のテストを作成
6. プルリクエストを送信

### 開発理念

- **魔法の用語**: spells、grimoire、spellbookを一貫して使用
- **テスト駆動開発**: テストを先に書き、その後実装
- **クリーンアーキテクチャ**: 明確な関心の分離
- **ユーザー体験優先**: ユーザーにはシンプル、開発者にはパワフル

## 📊 パフォーマンス

SilentCastは軽量で効率的に設計されています：

- **メモリ使用量**: アイドル時約15MB、アクティブ時約25MB
- **CPU使用率**: アイドル時<0.1%、実行時<1%
- **起動時間**: <100ms
- **ホットキー応答**: <10ms

最適化のヒントは[パフォーマンスガイド](docs/performance/README.md)を参照してください。

## 🔒 セキュリティ

- セルフアップデートチェック（オプション）を除き、ネットワーク接続なし
- すべての設定はローカル
- 管理者/sudo実行には明示的な設定が必要
- テレメトリやデータ収集なし

セキュリティの問題を報告: security@spherestacking.com

## 📄 ライセンス

SilentCastは[MITライセンス](LICENSE)の下でライセンスされたオープンソースソフトウェアです。

## 🙏 謝辞

これらの優れたライブラリで構築されています：
- [gohook](https://github.com/robotn/gohook) - クロスプラットフォームホットキーサポート
- [systray](https://github.com/getlantern/systray) - システムトレイ統合
- [fsnotify](https://github.com/fsnotify/fsnotify) - ファイル監視
- [lumberjack](https://github.com/natefinch/lumberjack) - ログローテーション
- [yaml.v3](https://github.com/go-yaml/yaml) - YAML設定

## 🌟 スター履歴

[![Star History Chart](https://api.star-history.com/svg?repos=SphereStacking/silentcast&type=Date)](https://star-history.com/#SphereStacking/silentcast&Date)

---

<div align="center">
  <p>キーボードの魔法を❤️する開発者たちによって🪄で作られました</p>
  
  <p>
    <a href="https://github.com/SphereStacking/silentcast/issues/new?labels=bug">バグを報告</a>
    ·
    <a href="https://github.com/SphereStacking/silentcast/issues/new?labels=enhancement">機能をリクエスト</a>
    ·
    <a href="https://spherestacking.github.io/SilentCast/">ドキュメントを読む</a>
  </p>
</div>