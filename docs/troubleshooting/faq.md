# よくある質問 (FAQ)

SilentCastの使用中によく遭遇する質問と解決策をまとめました。

## 🚀 インストールと初期設定

### Q: インストール後、ホットキーが反応しません
**A:** 以下を確認してください：

1. **スタブビルドではないか確認**
   ```bash
   ./silentcast --version
   # "stub" または "nogohook" が表示されていないことを確認
   ```

2. **権限設定（macOS）**
   - システム環境設定 > セキュリティとプライバシー > アクセシビリティ
   - SilentCast を追加してチェックを入れる

3. **管理者権限（Windows）**
   ```powershell
   # 管理者として実行
   Start-Process ./silentcast -Verb RunAs
   ```

### Q: 設定ファイルが見つからないエラーが出ます
**A:** spellbook.yml ファイルを作成してください：

```yaml
# 最小設定例
spells:
  e: editor
  
grimoire:
  editor:
    type: app
    command: "code"
    description: "VS Code を開く"
```

### Q: CGO 関連のビルドエラーが発生します
**A:** スタブビルドを使用してください：

```bash
# CGO 不要のスタブビルド
make build-stub

# または直接ビルド
go build -tags "nogohook notray" cmd/silentcast/main.go
```

## ⌨️ ホットキーとキー操作

### Q: Alt+Space が他のアプリと競合します
**A:** プレフィックスキーを変更してください：

```yaml
hotkeys:
  prefix: "ctrl+alt+space"  # または "cmd+space", "super+space"
  timeout: 1000
```

**一般的な代替案:**
- Windows: `ctrl+alt+grave` (Ctrl+Alt+`)
- macOS: `cmd+option+space`
- Linux: `super+alt+space`

### Q: キーシーケンスが認識されません
**A:** タイムアウト設定を調整してください：

```yaml
hotkeys:
  prefix: "alt+space"
  timeout: 2000           # プレフィックス後の待機時間
  sequence_timeout: 3000  # シーケンス全体のタイムアウト
```

**キー入力のコツ:**
1. プレフィックスキーを確実に離す
2. 次のキーまで少し待つ（500ms程度）
3. シーケンス内のキーは素早く入力

### Q: 特定のキーが認識されません
**A:** キーマッピングを確認してください：

```bash
# デバッグモードでキー監視
./silentcast --debug --monitor-keys

# キーマッピング確認
./silentcast --list-keymaps
```

## 📱 アプリケーション起動

### Q: アプリケーションが起動しません
**A:** パス設定を確認してください：

```yaml
grimoire:
  editor:
    type: app
    # フルパス指定（推奨）
    command: "/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code"
    # または環境変数使用
    command: "code"
    env:
      PATH: "/usr/local/bin:$PATH"
```

### Q: Windows でアプリのパスがわかりません
**A:** PowerShell で確認できます：

```powershell
# アプリケーションの場所を検索
Get-Command code
Get-ChildItem "C:\Program Files" -Recurse -Name "*.exe" | Select-String "code"

# レジストリからアプリケーション情報取得
Get-ItemProperty HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\* | Select-Object DisplayName, InstallLocation
```

### Q: macOS でアプリケーションパスを見つけるには？
**A:** ターミナルで確認できます：

```bash
# アプリケーションバンドル確認
ls -la /Applications/ | grep -i "visual studio"

# 実行可能ファイルのパス確認
find /Applications -name "*code*" -type f -executable

# which コマンド使用
which code
```

## 📜 スクリプト実行

### Q: スクリプトが「permission denied」で失敗します
**A:** 実行権限を確認してください：

```bash
# 権限確認
ls -la script.sh

# 実行権限付与
chmod +x script.sh

# または設定で指定
```

```yaml
grimoire:
  script-action:
    type: script
    script_file: "/path/to/script.sh"
    executable: true  # 自動的に実行権限付与
```

### Q: スクリプト内でコマンドが見つからないエラー
**A:** 環境変数を設定してください：

```yaml
grimoire:
  git-status:
    type: script
    command: "git status"
    env:
      PATH: "/usr/local/bin:/usr/bin:/bin"
    shell: "/bin/bash"
```

### Q: Windows でPowerShell スクリプトを実行したい
**A:** PowerShell を明示的に指定してください：

```yaml
grimoire:
  powershell-script:
    type: script
    command: "Get-Process | Where-Object {$_.ProcessName -eq 'notepad'}"
    shell: "powershell.exe"
    # または
    command: "powershell.exe -Command \"Get-Process\""
```

## 🌐 URL とブラウザ

### Q: URL が既定のブラウザで開きません
**A:** ブラウザを明示的に指定してください：

```yaml
grimoire:
  website:
    type: url
    command: "https://github.com"
    browser: "google-chrome"  # または firefox, safari など
```

### Q: 特定のブラウザプロファイルで開きたい
**A:** ブラウザ引数を使用してください：

```yaml
grimoire:
  github-work:
    type: url
    command: "https://github.com"
    browser: "google-chrome"
    args: ["--profile-directory=Work"]
    
  firefox-private:
    type: url
    command: "https://example.com"  
    browser: "firefox"
    args: ["-private-window"]
```

## 🔧 設定と構成

### Q: 設定ファイルの変更が反映されません
**A:** 設定監視が有効か確認してください：

```yaml
daemon:
  config_watch: true  # 設定ファイル監視を有効化
```

または手動で再読み込み：
```bash
./silentcast --reload-config
```

### Q: プラットフォーム固有の設定はどう書きますか？
**A:** プラットフォーム別セクションを使用してください：

```yaml
grimoire:
  editor:
    type: app
    description: "テキストエディタ"
    windows:
      command: "notepad.exe"
    macos:
      command: "open"
      args: ["-a", "TextEdit"]
    linux:
      command: "gedit"
```

### Q: 環境変数を設定で使いたい
**A:** `${変数名}` の形式で使用できます：

```yaml
grimoire:
  home-folder:
    type: app
    command: "open"
    args: ["${HOME}"]
    
  project:
    type: app
    command: "code"
    args: ["${PROJECT_DIR}/src"]
    env:
      PROJECT_DIR: "/home/user/projects"
```

## 🔐 権限とセキュリティ

### Q: macOS で「アクセシビリティ権限が必要」と表示されます
**A:** システム環境設定で権限を付与してください：

1. システム環境設定 > セキュリティとプライバシー > プライバシー
2. アクセシビリティを選択
3. 🔒 をクリックして管理者パスワード入力
4. + ボタンで SilentCast を追加
5. チェックボックスにチェック
6. アプリケーションを再起動

### Q: Windows で「管理者権限が必要」と表示されます
**A:** 管理者として実行するか、UAC を調整してください：

```powershell
# 管理者として実行
Start-Process ./silentcast -Verb RunAs

# またはタスクスケジューラで昇格タスク作成
schtasks /create /tn "SilentCast" /tr "C:\path\to\silentcast.exe" /rl HIGHEST /f
```

### Q: Linux でホットキーが登録できません
**A:** 権限とデスクトップ環境を確認してください：

```bash
# input グループに追加
sudo usermod -a -G input $(whoami)

# X11 権限設定
xhost +local:

# デスクトップ環境固有設定（GNOME の例）
gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings "['/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom0/']"
```

## 📊 パフォーマンスと問題

### Q: SilentCast が CPU を大量に使用します
**A:** ポーリング間隔を調整してください：

```yaml
performance:
  hotkey_poll_interval: 100  # ミリ秒（デフォルト: 50）
  
# 不要なスペルを削除
spells:
  # 必要最小限のスペルのみ残す
  e: editor
  t: terminal
```

### Q: アプリケーションの応答が遅い
**A:** デバッグモードで原因を特定してください：

```bash
# パフォーマンス監視
./silentcast --debug --performance-monitor

# 特定のアクションをテスト
./silentcast --debug --once --spell slow-action
```

### Q: システムトレイアイコンが表示されません
**A:** トレイ設定を確認してください：

```yaml
daemon:
  tray: true
  tray_icon: "/path/to/icon.png"  # カスタムアイコン
```

または tray なしで実行：
```bash
./silentcast --no-tray
```

## 🛠️ トラブルシューティング

### Q: ログファイルはどこにありますか？
**A:** プラットフォーム別の場所：

- **Linux/macOS**: `~/.local/share/silentcast/silentcast.log`
- **Windows**: `%APPDATA%\SilentCast\silentcast.log`
- **カスタム**: `--log-file` オプションで指定

### Q: 設定の検証方法は？
**A:** 検証コマンドを使用してください：

```bash
# 設定ファイル検証
./silentcast --validate-config

# 設定内容確認
./silentcast --dump-config

# 特定スペル確認
./silentcast --show-spell editor
```

### Q: 問題を報告するときに必要な情報は？
**A:** 以下の情報を収集してください：

```bash
# システム情報
./silentcast --system-info

# 診断情報
./silentcast --diagnose

# ログファイル
cp silentcast.log debug-log.txt

# 設定ファイル
cp spellbook.yml debug-config.yml
```

## 🔄 アップデートと移行

### Q: 新しいバージョンへの更新方法は？
**A:** ビルド済みバイナリを置き換えるか、再ビルドしてください：

```bash
# 新しいバージョンのビルド
git pull origin main
make build

# 設定ファイルはそのまま使用可能（通常）
```

### Q: 設定ファイルの互換性は？
**A:** マイナーバージョンアップは通常互換性があります：

```bash
# 設定互換性確認
./silentcast --check-config-compatibility

# 設定移行（必要に応じて）
./silentcast --migrate-config
```

## 💡 使用のコツ

### Q: よく使う設定パターンは？
**A:** 以下のパターンが便利です：

```yaml
# VS Code スタイルのキーバインド
spells:
  "g,s": git-status
  "g,p": git-pull
  "g,c": git-commit
  "d,r": docker-restart
  
# アプリケーション起動
spells:
  e: editor
  t: terminal
  b: browser
  f: file-manager

# システム操作
spells:
  "s,l": system-lock
  "s,s": system-sleep
  "s,r": system-restart
```

### Q: デバッグ時の便利な設定は？
**A:** デバッグ用設定を作成してください：

```yaml
# debug-spellbook.yml
logger:
  level: debug
  console: true
  
spells:
  test: test-action
  
grimoire:
  test-action:
    type: script
    command: echo "Test: $(date)"
    show_output: true
```

## 🔗 その他のリソース

- [詳細なトラブルシューティングガイド](../troubleshooting/)
- [設定リファレンス](../guide/configuration.md)
- [プラットフォーム固有の問題](platform-specific.md)
- [デバッグガイド](debugging.md)
- [サポート情報](support.md)

---

**この FAQ で解決しない場合は、[デバッグガイド](debugging.md) を参照するか、[サポート](support.md) にお問い合わせください。**