# SilentCast Troubleshooting Guide

このガイドでは、SilentCastの一般的な問題と解決方法について説明します。問題を迅速に解決するために、まず「クイック診断」を試してから、詳細なトラブルシューティングに進んでください。

## 🚀 クイック診断

### 問題の種類を特定
1. **インストール問題**: SilentCastが起動しない、ビルドエラー
2. **設定問題**: 設定が読み込まれない、YAML エラー
3. **ホットキー問題**: キーが認識されない、他のアプリと競合
4. **アクション実行問題**: アプリが起動しない、スクリプトが実行されない
5. **パフォーマンス問題**: 動作が遅い、高いCPU/メモリ使用率
6. **プラットフォーム固有問題**: OS特有の権限やセキュリティ問題

### 最初に試すこと

#### 1. バージョンとビルド情報の確認
```bash
./silentcast --version
```

#### 2. 設定の検証
```bash
./silentcast --validate-config
```

#### 3. デバッグモードでの実行
```bash
./silentcast --debug --no-tray
```

#### 4. ログの確認
```bash
tail -f silentcast.log
```

## 📋 詳細トラブルシューティング

### [インストール問題](installation.md)
- CGOビルド問題とスタブモード
- 依存関係の問題
- パーミッション問題
- プラットフォーム固有のインストール

### [設定問題](configuration.md)
- YAML 構文エラー
- 設定検証失敗
- ファイル監視問題
- カスケード設定読み込み

### [ホットキー問題](hotkeys.md)
- ホットキー登録失敗
- キー競合の解決
- プラットフォーム固有のキーマッピング
- アクセシビリティ権限

### [アクション実行問題](actions.md)
- アプリケーション起動失敗
- スクリプト実行エラー
- URL開く問題
- パス解決問題

### [パフォーマンス問題](performance.md)
- 高CPU/メモリ使用率
- 遅い応答時間
- リソースリーク
- 最適化設定

### [権限とセキュリティ](permissions.md)
- macOS アクセシビリティ権限
- Windows UAC と昇格
- Linux デスクトップ環境
- セキュリティソフトとの競合

### [プラットフォーム固有](platform-specific.md)
- Windows 固有の問題
- macOS 固有の問題
- Linux 固有の問題

### [デバッグとログ分析](debugging.md)
- ログレベルの設定
- デバッグツールの使用
- エラーメッセージの理解
- 問題の報告方法

## 🔧 診断ツール

### 内蔵診断コマンド

#### 設定検証
```bash
# 設定ファイルの構文と内容を検証
./silentcast --validate-config

# 詳細な設定情報を表示
./silentcast --dump-config
```

#### システム情報
```bash
# システム環境とSilentCast情報を表示
./silentcast --system-info
```

#### 権限チェック
```bash
# 必要な権限の確認
./silentcast --check-permissions
```

### ログレベルの設定

```yaml
# spellbook.yml
logger:
  level: debug  # error, warn, info, debug
  file: silentcast.log
```

### デバッグ出力の例

```bash
# デバッグモードで実行
./silentcast --debug --no-tray

# 出力例:
[DEBUG] Configuration loaded: /path/to/spellbook.yml
[DEBUG] Hotkey registered: alt+space
[DEBUG] Action executor created: type=app
[INFO] SilentCast started successfully
```

## ❓ よくある質問 (FAQ)

### Q: SilentCastが起動しない
**A**: まずスタブモードを試してください:
```bash
make build-stub
./build/silentcast --no-tray
```

### Q: ホットキーが反応しない
**A**: 
1. 他のアプリケーションとの競合を確認
2. 管理者権限で実行を試行（Windowsの場合）
3. アクセシビリティ権限を確認（macOSの場合）

### Q: 設定が反映されない
**A**:
1. YAML構文エラーがないか確認: `--validate-config`
2. ファイル監視が動作しているか確認: `--debug` で確認
3. キャッシュをクリア: 設定ファイルを再保存

### Q: アプリケーションが起動しない
**A**:
1. アプリケーションパスが正しいか確認
2. 実行権限があるか確認
3. PATH環境変数が設定されているか確認

### Q: 動作が遅い
**A**:
1. パフォーマンス設定を調整: [performance-optimization.md](../guide/performance-optimization.md)
2. ログレベルを下げる: `level: warn`
3. 不要な設定監視を無効化

## 🆘 サポートとコミュニティ

### 問題の報告
バグや問題を発見した場合は、以下の情報と共にGitHub Issueを作成してください:

1. **SilentCastバージョン**: `./silentcast --version`
2. **OS情報**: OS名、バージョン
3. **エラーメッセージ**: 完全なエラーログ
4. **再現手順**: 問題を再現する詳細な手順
5. **設定ファイル**: spellbook.yml（機密情報は除く）

### テンプレート
```markdown
## 環境情報
- SilentCast Version: 
- OS: 
- Installation Method: 

## 問題の説明


## 再現手順
1. 
2. 
3. 

## 期待される動作


## 実際の動作


## ログ/エラーメッセージ
```

### コミュニティサポート
- [GitHub Discussions](https://github.com/SphereStacking/SilentCast/discussions)
- [GitHub Issues](https://github.com/SphereStacking/SilentCast/issues)
- [ドキュメント](https://spherestacking.github.io/SilentCast/)

## 🔗 関連リソース

- [インストールガイド](../guide/installation.md)
- [設定ガイド](../guide/configuration.md)
- [パフォーマンス最適化](../guide/performance-optimization.md)
- [開発者ガイド](../development/README.md)
- [FAQ](../guide/faq.md)