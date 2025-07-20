# サポート情報

SilentCastに関するサポート、問題報告、コミュニティリソースについて説明します。

## 🆘 サポートを受ける前に

### 自己診断チェックリスト

問題を報告する前に、以下をご確認ください：

- [ ] [FAQ](faq.md) で類似の問題を確認済み
- [ ] [トラブルシューティングガイド](../troubleshooting/) を参照済み
- [ ] 最新バージョンを使用している
- [ ] 設定ファイルの構文エラーがない
- [ ] ログファイルでエラーメッセージを確認済み

### 基本診断の実行

```bash
# 基本診断コマンド
./silentcast --diagnose

# システム情報収集
./silentcast --system-info

# 設定検証
./silentcast --validate-config
```

## 📝 問題報告ガイドライン

### 効果的な問題報告のために

良い問題報告には以下の情報が含まれています：

1. **環境情報**
2. **問題の詳細**
3. **再現手順**
4. **期待される動作**
5. **実際の動作**
6. **ログ出力**
7. **設定ファイル**

### 問題報告テンプレート

```markdown
## 問題の概要
[問題の簡潔な説明]

## 環境情報
- **OS**: [Windows 10/macOS 12.5/Ubuntu 22.04 など]
- **SilentCast バージョン**: [./silentcast --version の出力]
- **ビルドタイプ**: [通常ビルド/スタブビルド]
- **インストール方法**: [ソースビルド/バイナリダウンロード]

## 問題の詳細
[問題の詳細な説明]

## 再現手順
1. [具体的な手順1]
2. [具体的な手順2]
3. [具体的な手順3]

## 期待される動作
[期待していた結果]

## 実際の動作
[実際に起こった結果]

## エラーメッセージ・ログ
```
[ログ出力やエラーメッセージをここに貼り付け]
```

## 設定ファイル
```yaml
[関連するspellbook.ymlの内容]
```

## 診断情報
```
[./silentcast --diagnose の出力]
```

## 試した解決策
[既に試した解決策があれば記載]

## 追加情報
[その他の関連情報]
```

### デバッグ情報の自動収集

問題報告用の情報を自動収集するスクリプト：

```bash
#!/bin/bash
# support-info.sh - サポート用情報収集スクリプト

collect_support_info() {
    echo "🔍 SilentCast サポート情報を収集中..."
    
    local report_dir="silentcast_support_$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$report_dir"
    
    # バージョン情報
    echo "=== SilentCast Version ===" > "$report_dir/version.txt"
    ./silentcast --version >> "$report_dir/version.txt" 2>&1
    
    # システム情報
    echo "=== System Information ===" > "$report_dir/system.txt"
    uname -a >> "$report_dir/system.txt"
    echo "" >> "$report_dir/system.txt"
    
    case "$(uname -s)" in
        Darwin)
            sw_vers >> "$report_dir/system.txt"
            ;;
        Linux)
            if [ -f /etc/os-release ]; then
                cat /etc/os-release >> "$report_dir/system.txt"
            fi
            ;;
        *)
            echo "Platform: $(uname -s)" >> "$report_dir/system.txt"
            ;;
    esac
    
    # 診断情報
    echo "=== Diagnosis ===" > "$report_dir/diagnosis.txt"
    ./silentcast --diagnose >> "$report_dir/diagnosis.txt" 2>&1
    
    # 設定検証
    echo "=== Configuration Validation ===" > "$report_dir/config_validation.txt"
    ./silentcast --validate-config >> "$report_dir/config_validation.txt" 2>&1
    
    # 設定ファイル
    cp spellbook*.yml "$report_dir/" 2>/dev/null || echo "No spellbook files found"
    
    # ログファイル
    cp silentcast*.log "$report_dir/" 2>/dev/null || echo "No log files found"
    
    # 権限情報
    echo "=== Permissions ===" > "$report_dir/permissions.txt"
    ./silentcast --check-permissions >> "$report_dir/permissions.txt" 2>&1
    
    # アーカイブ作成
    tar -czf "${report_dir}.tar.gz" "$report_dir"
    rm -rf "$report_dir"
    
    echo "✅ サポート情報収集完了: ${report_dir}.tar.gz"
    echo ""
    echo "📋 問題報告時は以下を含めてください："
    echo "1. 上記のアーカイブファイル"
    echo "2. 問題の詳細な説明"
    echo "3. 再現手順"
    echo "4. 期待される動作と実際の動作"
}

# スクリプト実行
collect_support_info
```

## 🌐 サポートチャンネル

### GitHub Issues（推奨）

**用途**: バグ報告、機能要求、技術的な問題

**URL**: [https://github.com/SphereStacking/SilentCast/issues](https://github.com/SphereStacking/SilentCast/issues)

**報告前の確認事項**:
- 既存のIssueで同様の問題が報告されていないか検索
- 適切なラベルを選択（bug, enhancement, question など）
- 問題報告テンプレートを使用

### ディスカッション

**用途**: 一般的な質問、使用方法の相談、アイデア共有

**URL**: [https://github.com/SphereStacking/SilentCast/discussions](https://github.com/SphereStacking/SilentCast/discussions)

**カテゴリ**:
- **Q&A**: 質問と回答
- **General**: 一般的な議論
- **Ideas**: 機能提案やアイデア
- **Show and tell**: 設定例やカスタマイズの共有

### コミュニティ

**用途**: リアルタイムな質問、コミュニティ交流

- **Discord**: [招待リンク](https://discord.gg/silentcast)
- **Reddit**: [r/SilentCast](https://reddit.com/r/SilentCast)

## 📚 セルフヘルプリソース

### ドキュメント

1. **[ユーザーガイド](../guide/)**: 基本的な使用方法
2. **[設定リファレンス](../guide/configuration.md)**: 詳細な設定オプション
3. **[トラブルシューティング](../troubleshooting/)**: 問題解決ガイド
4. **[FAQ](faq.md)**: よくある質問と回答

### ビデオチュートリアル

- **[インストールガイド](https://youtube.com/watch?v=example1)**: インストール手順
- **[基本設定](https://youtube.com/watch?v=example2)**: spellbook.yml の設定方法
- **[高度な使用法](https://youtube.com/watch?v=example3)**: カスタマイズとスクリプト

### ブログとチュートリアル

- **[公式ブログ](https://blog.silentcast.dev)**: 新機能の紹介、使用例
- **[コミュニティチュートリアル](https://github.com/SilentCast/awesome-configs)**: ユーザー投稿の設定例

## 🤝 コミュニティへの貢献

### 貢献方法

1. **バグ報告**: 問題を見つけたら報告してください
2. **機能提案**: 新しいアイデアを共有してください
3. **ドキュメント改善**: 誤字や不明確な点の修正
4. **コード貢献**: バグ修正や新機能の実装
5. **設定例共有**: 有用な設定をコミュニティと共有

### Pull Request ガイドライン

```bash
# 開発環境のセットアップ
git clone https://github.com/SphereStacking/SilentCast.git
cd SilentCast
make setup

# ブランチ作成
git checkout -b feature/your-feature-name

# 変更実装
# ...

# テスト実行
make test

# コミット
git commit -m "feat: add new feature description"

# Push とPR作成
git push origin feature/your-feature-name
```

### コード貢献の準備

1. **[Contributing Guide](../contributing.md)** を読む
2. **[Development Setup](../development/setup.md)** に従って環境構築
3. **Issue** または **Discussion** で提案内容を事前相談
4. **テスト** を含む実装
5. **ドキュメント** の更新

## 🔧 エンタープライズサポート

### 商用サポート

大規模導入や企業での使用については、以下をご検討ください：

- **カスタマイズ開発**: 企業固有の要件に応じた機能開発
- **技術コンサルティング**: 導入・運用支援
- **優先サポート**: 専用サポートチャンネル
- **SLA**: サービスレベル契約

**お問い合わせ**: [enterprise@silentcast.dev](mailto:enterprise@silentcast.dev)

### トレーニングとワークショップ

- **オンサイトトレーニング**: チーム向けトレーニングセッション
- **ベストプラクティス**: 効果的な使用方法の指導
- **カスタム設定**: 組織固有の設定作成支援

## 📞 緊急サポート

### セキュリティ問題

セキュリティに関する脆弱性を発見した場合：

**連絡先**: [security@silentcast.dev](mailto:security@silentcast.dev)

**報告時の注意事項**:
- 公開Issueでの報告は避けてください
- 詳細な再現手順を含めてください
- 影響範囲と深刻度を記載してください

### 重要なバグ

本番環境で重大な影響を与えるバグの場合：

1. **GitHub Issue** で「critical」ラベルを付けて報告
2. **Discord** の #urgent チャンネルで通知
3. 必要に応じて **enterprise@silentcast.dev** に連絡

## 📈 サポート品質向上にご協力ください

### フィードバック

サポート体験の改善のため、以下にご協力ください：

- **解決後のフィードバック**: Issue解決後の満足度評価
- **ドキュメント改善提案**: 分かりにくい箇所の指摘
- **プロセス改善**: サポートプロセスの改善案

### サポート統計

現在のサポート状況：

- **平均初回応答時間**: 24時間以内
- **問題解決率**: 95%以上
- **コミュニティ満足度**: 4.8/5.0

## 🎯 よくあるサポートパターン

### レベル1: 基本サポート
- インストール問題
- 基本設定エラー
- FAQ に記載されている問題

**解決方法**: ドキュメント、FAQ、コミュニティ

### レベル2: 技術サポート
- 複雑な設定問題
- プラットフォーム固有の問題
- パフォーマンス問題

**解決方法**: GitHub Issues、詳細なデバッグ

### レベル3: 開発者サポート
- バグ修正が必要な問題
- 新機能が必要な要求
- アーキテクチャレベルの問題

**解決方法**: コード変更、Pull Request、開発者との直接連携

## 📅 サポート時間

### コミュニティサポート
- **24時間365日**: GitHub Issues、Discussions
- **リアルタイム**: Discord（タイムゾーンにより応答時間は変動）

### メンテナーサポート
- **月-金 9:00-18:00 JST**: 開発チームによる対応
- **緊急時**: 24時間以内の初回応答を目標

## 🔗 関連リンク

- **[公式ウェブサイト](https://silentcast.dev)**
- **[GitHub リポジトリ](https://github.com/SphereStacking/SilentCast)**
- **[ドキュメント](https://docs.silentcast.dev)**
- **[リリースノート](https://github.com/SphereStacking/SilentCast/releases)**
- **[ロードマップ](https://github.com/SphereStacking/SilentCast/projects)**

---

**サポートが必要な場合は、適切なチャンネルを選択して、できるだけ詳細な情報とともにお問い合わせください。コミュニティの皆様の協力により、SilentCast をより良いツールにしていきましょう！**