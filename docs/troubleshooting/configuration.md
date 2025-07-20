# 設定問題のトラブルシューティング

SilentCastの設定に関する問題と解決方法について説明します。

## 🚀 クイック診断

### 設定検証コマンド
```bash
# 設定ファイルの検証
./silentcast --validate-config

# 設定内容の表示
./silentcast --dump-config

# デバッグモードで設定読み込み確認
./silentcast --debug --no-tray
```

### 最も一般的な設定問題
1. **YAML構文エラー** → インデントやクォート確認
2. **ファイルが見つからない** → パス確認
3. **設定が反映されない** → ファイル監視とキャッシュ
4. **アクション設定エラー** → 必須フィールド確認

## 📝 YAML構文問題

### インデント問題

#### ❌ 間違った例
```yaml
spells:
e: editor          # インデントが足りない
  t:terminal       # スペースとタブが混在
grimoire:
    editor:        # インデント過多
  type: app        # インデント不一致
```

#### ✅ 正しい例
```yaml
spells:
  e: editor        # 2スペースインデント
  t: terminal      # 一貫したインデント

grimoire:
  editor:          # 2スペースインデント
    type: app      # 4スペースインデント（2段階）
    command: code
```

### クォート問題

#### ❌ 間違った例
```yaml
spells:
  g,s: git-status  # コンマを含む場合はクォートが必要

grimoire:
  git-status:
    command: git status && echo "完了"  # 特殊文字の処理
```

#### ✅ 正しい例
```yaml
spells:
  "g,s": git-status  # コンマを含む場合

grimoire:
  git-status:
    command: 'git status && echo "完了"'  # シングルクォートで囲む
    # または
    command: |
      git status
      echo "完了"
```

### 特殊文字とエスケープ

```yaml
spells:
  "ctrl+shift+p": command-palette  # 特殊キー組み合わせ
  "\"": quote-action               # クォート文字

grimoire:
  quote-action:
    command: 'echo "引用符: \""'     # エスケープ
    
  multiline:
    command: |
      echo "複数行の"
      echo "コマンド例"
    description: >
      長い説明文は
      このように書ける
```

## 📁 ファイルパスとディレクトリ

### 設定ファイルの場所

#### デフォルトパス
```bash
# Linux
~/.config/silentcast/spellbook.yml

# macOS  
~/Library/Application Support/silentcast/spellbook.yml

# Windows
%APPDATA%\silentcast\spellbook.yml
```

#### カスタムパス
```bash
# コマンドライン指定
./silentcast --config /path/to/custom/spellbook.yml

# 環境変数
export SILENTCAST_CONFIG_DIR=/custom/config/dir
```

### カスケード設定

#### ファイル読み込み順序
1. `spellbook.yml` (基本設定)
2. `spellbook.linux.yml` (Linux固有)
3. `spellbook.darwin.yml` (macOS固有)  
4. `spellbook.windows.yml` (Windows固有)

#### 例: プラットフォーム固有設定
```yaml
# spellbook.yml (共通設定)
spells:
  e: editor
  t: terminal

grimoire:
  editor:
    type: app
    command: code  # デフォルト

---
# spellbook.darwin.yml (macOS専用)
grimoire:
  editor:
    command: /Applications/Visual Studio Code.app/Contents/Resources/app/bin/code
  terminal:
    command: open -a Terminal

---
# spellbook.windows.yml (Windows専用)
grimoire:
  editor:
    command: code.exe
  terminal:
    command: cmd.exe
```

## 🔄 設定の再読み込み

### ファイル監視問題

#### 監視が動作しない場合
```bash
# ファイル監視の確認
./silentcast --debug --no-tray

# ログ出力例:
# [DEBUG] File watcher started: /path/to/spellbook.yml
# [DEBUG] Configuration reloaded: file changed
```

#### 手動再読み込み
```bash
# プロセスにSIGHUPシグナル送信
kill -HUP $(pgrep silentcast)

# または設定ファイルの touch
touch spellbook.yml
```

### キャッシュ問題

#### 設定キャッシュのクリア
```bash
# 設定ファイルを再保存
cp spellbook.yml spellbook.yml.bak
mv spellbook.yml.bak spellbook.yml

# または
echo "" >> spellbook.yml
```

## ⚙️ アクション設定の検証

### 必須フィールドの確認

#### アプリケーションアクション
```yaml
grimoire:
  editor:
    type: app        # 必須
    command: code    # 必須
    description: "VS Code editor"  # 推奨
```

#### スクリプトアクション
```yaml
grimoire:
  git-status:
    type: script     # 必須
    command: git status  # 必須
    show_output: true    # オプション
    working_dir: /project  # オプション
```

#### URLアクション
```yaml
grimoire:
  github:
    type: url        # 必須
    command: https://github.com  # 必須（URL形式）
```

### 設定検証エラーの解読

#### エラーメッセージの例
```
Error: spell 'e' references unknown action 'editr'
Available actions: [editor, terminal, browser]
Suggestion: check grimoire section for typos
```

#### 解決方法
1. スペルミスの確認
2. アクション名の一致確認
3. grimoireセクションの存在確認

## 🔧 高度な設定問題

### 環境変数の展開

```yaml
grimoire:
  home:
    type: script
    command: echo $HOME        # 展開される
    working_dir: ${PROJECT_DIR}  # 展開される
    env:
      CUSTOM_VAR: ${USER}_custom  # 展開される
```

### 条件付き設定

```yaml
# 開発環境用設定
spells:
  d: dev-server

grimoire:
  dev-server:
    type: script
    command: |
      if [ "$NODE_ENV" = "development" ]; then
        npm run dev
      else
        echo "Development environment not set"
      fi
```

### パフォーマンス設定

```yaml
# パフォーマンス最適化設定
performance:
  enable_optimization: true
  buffer_size: 2048
  gc_percent: 75
  max_idle_time: 10m

daemon:
  config_watch: true   # ファイル監視有効
  log_level: warn      # ログレベル調整

logger:
  level: warn
  max_size: 10
```

## 🧪 設定テストと検証

### 段階的な設定テスト

#### 1. 最小設定でテスト
```yaml
# minimal.yml
spells:
  t: test

grimoire:
  test:
    type: script
    command: echo "test"
```

```bash
./silentcast --config minimal.yml --debug --no-tray
```

#### 2. 一つずつアクションを追加
```yaml
spells:
  t: test
  e: editor  # 追加

grimoire:
  test:
    type: script
    command: echo "test"
  editor:     # 追加
    type: app
    command: echo  # 安全なコマンド
```

### 設定の自動テスト

```bash
# 設定ファイルの構文チェック
yamllint spellbook.yml

# SilentCast独自の検証
./silentcast --validate-config

# ドライランモード（実際には実行しない）
./silentcast --dry-run --config test.yml
```

## 📋 設定問題診断チェックリスト

### YAML構文
- [ ] インデントが一貫している（スペースのみ）
- [ ] クォート文字が正しく使用されている
- [ ] 特殊文字がエスケープされている
- [ ] コロンの後にスペースがある

### ファイルとパス
- [ ] 設定ファイルが存在し、読み取り可能
- [ ] パスが正しい（絶対パスまたは相対パス）
- [ ] カスケード設定が正しい順序で読み込まれている

### アクション設定
- [ ] 全てのspellに対応するgrimoireエントリが存在
- [ ] 必須フィールド（type, command）が設定されている
- [ ] コマンドパスが有効
- [ ] 権限が適切に設定されている

### システム設定
- [ ] ファイル監視が動作している
- [ ] 環境変数が正しく展開されている
- [ ] プラットフォーム固有設定が適用されている

## 🆘 それでも解決しない場合

### デバッグ情報の収集
```bash
# 詳細なデバッグ情報
./silentcast --debug --dump-config > debug-config.txt

# 設定検証の詳細
./silentcast --validate-config --verbose > validation.txt

# ファイル権限の確認
ls -la spellbook.yml
ls -la ~/.config/silentcast/
```

### よくある原因
1. **タブ文字の混入**: エディタでタブを可視化
2. **ファイルエンコーディング**: UTF-8で保存
3. **改行コード**: Unix形式（LF）を使用
4. **ファイル権限**: 読み取り権限を確認
5. **ディスク容量**: 十分な空き容量を確認

## 🔗 関連リソース

- [設定ガイド](../guide/configuration.md)
- [設定ファイル構造](../config/file-structure.md)
- [設定例](../../examples/config/)
- [FAQ](../guide/faq.md)