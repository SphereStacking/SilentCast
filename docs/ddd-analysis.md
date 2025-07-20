# SilentCast DDD導入検討書

## 概要

このドキュメントでは、SilentCastプロジェクトにDomain-Driven Design (DDD)を導入することの是非を検討します。

## 現在のアーキテクチャ分析

### 現状の構造
```
app/
├── cmd/silentcast/     # エントリーポイント
├── internal/           # 内部パッケージ
│   ├── action/         # アクション実行
│   ├── config/         # 設定管理
│   ├── hotkey/         # ホットキー検知
│   ├── notify/         # 通知
│   ├── permission/     # 権限管理
│   └── tray/          # システムトレイ
└── pkg/               # 公開パッケージ
```

### 現在のアプローチ
- 機能ベースのパッケージ分割
- インターフェースによる抽象化
- プラットフォーム固有実装の分離

## DDD導入の検討

### SilentCastのドメイン分析

#### コアドメイン
1. **Spell Management (呪文管理)**
   - Spell (呪文): ホットキーの組み合わせ
   - Grimoire Entry (魔導書エントリー): 実行されるアクション定義
   - Spellbook (呪文書): 設定ファイル全体

2. **Action Execution (アクション実行)**
   - Action: 実行可能なタスク
   - Executor: アクション実行者
   - Result: 実行結果

3. **Input Detection (入力検知)**
   - Hotkey: キーボード入力
   - Sequence: キーシーケンス
   - Trigger: トリガーイベント

### DDD適用案

```
app/
├── cmd/silentcast/
├── domain/              # ドメイン層
│   ├── spell/          # 呪文ドメイン
│   │   ├── spell.go    # Spell エンティティ
│   │   ├── grimoire.go # Grimoire 値オブジェクト
│   │   └── repository.go
│   ├── action/         # アクションドメイン
│   │   ├── action.go   # Action エンティティ
│   │   ├── executor.go # Executor インターフェース
│   │   └── result.go   # Result 値オブジェクト
│   └── hotkey/         # ホットキードメイン
│       ├── hotkey.go   # Hotkey 値オブジェクト
│       └── detector.go # Detector インターフェース
├── application/        # アプリケーション層
│   ├── spell/         # 呪文管理ユースケース
│   ├── action/        # アクション実行ユースケース
│   └── hotkey/        # ホットキー検知ユースケース
├── infrastructure/     # インフラストラクチャ層
│   ├── config/        # 設定ファイル読み込み
│   ├── executor/      # OS固有のアクション実行
│   ├── hotkey/        # OS固有のキー検知
│   └── notify/        # OS固有の通知
└── presentation/       # プレゼンテーション層
    ├── cli/           # CLIインターフェース
    └── tray/          # システムトレイ
```

## メリット・デメリット分析

### DDDを導入するメリット

1. **ビジネスロジックの明確化**
   - 呪文管理の概念が明確になる
   - アクション実行のルールが集約される

2. **テスタビリティの向上**
   - ドメインロジックが独立してテスト可能
   - インフラから分離された純粋なビジネスロジック

3. **拡張性の向上**
   - 新しいアクションタイプの追加が容易
   - 新しい入力方法（音声など）の追加が容易

4. **チーム開発での利点**
   - ユビキタス言語による共通理解
   - 責任範囲の明確化

### DDDを導入するデメリット

1. **複雑性の増加**
   - 小規模プロジェクトには過剰
   - 学習コストが高い

2. **開発速度の低下**
   - 初期実装に時間がかかる
   - ボイラープレートコードの増加

3. **現在のコードベースとの乖離**
   - 大規模なリファクタリングが必要
   - 移行期間中の複雑性

## 推奨事項

### 現時点では部分的な導入を推奨

SilentCastの現在の規模と複雑さを考慮すると、完全なDDD導入よりも、DDDの良い部分を取り入れた**軽量なアプローチ**が適切です。

#### 推奨する改善

1. **ドメインモデルの導入**
   ```go
   // domain/spell/spell.go
   type Spell struct {
       Key      string
       Sequence []string
   }
   
   // domain/action/action.go
   type Action struct {
       Type    ActionType
       Command string
       Args    []string
   }
   ```

2. **ユースケース層の追加**
   ```go
   // application/execute_spell.go
   type ExecuteSpellUseCase struct {
       spellRepo  SpellRepository
       executor   ActionExecutor
       notifier   Notifier
   }
   ```

3. **リポジトリパターンの採用**
   ```go
   // domain/spell/repository.go
   type SpellRepository interface {
       FindByKey(key string) (*Spell, error)
       LoadFromConfig(path string) error
   }
   ```

### 段階的移行計画

#### Phase 1: ドメインモデルの定義（推奨）
- Spell, Action, Hotkeyなどの基本型を定義
- 既存コードをラップする形で導入
- 影響範囲: 小

#### Phase 2: リポジトリパターンの導入（推奨）
- 設定ロードをリポジトリとして抽象化
- テスタビリティの向上
- 影響範囲: 中

#### Phase 3: ユースケース層の追加（オプション）
- ビジネスロジックの集約
- 依存性の整理
- 影響範囲: 大

#### Phase 4: 完全なDDD構造への移行（非推奨）
- プロジェクトが大規模化した場合のみ検討
- チーム開発が本格化した場合のみ検討

## 結論

### 現時点での推奨アプローチ

1. **T068（パッケージ構造の整理）で以下を実施**：
   - ドメインモデルの定義（Phase 1）
   - リポジトリパターンの部分的導入（Phase 2）
   - 既存構造を大きく変えない範囲での改善

2. **DDDの考え方を部分的に採用**：
   - ユビキタス言語（Spell, Grimoire等）の徹底
   - ドメインロジックとインフラの分離
   - インターフェースによる抽象化の強化

3. **将来的な拡張性を考慮**：
   - 完全なDDD移行への道を残す
   - 段階的な移行が可能な構造

### 実装例

```go
// domain/types.go - 軽量なドメインモデル
package domain

type Spell struct {
    Key      string
    Action   Action
}

type Action struct {
    Type     string
    Command  string
    Args     []string
    Timeout  time.Duration
}

// internal/spell/service.go - 既存構造を活かしたサービス層
package spell

type Service struct {
    config   *config.Loader
    executor action.Executor
}

func (s *Service) ExecuteSpell(key string) error {
    spell, err := s.config.GetSpell(key)
    if err != nil {
        return err
    }
    
    return s.executor.Execute(spell.Action)
}
```

この軽量なアプローチにより、DDDの利点を享受しながら、現在のコードベースを大きく変更することなく、段階的な改善が可能になります。