# パフォーマンス問題のトラブルシューティング

SilentCastのパフォーマンス問題の診断と最適化について説明します。

## 🚀 クイック診断

### パフォーマンス問題の特定
```bash
# リソース使用状況の確認
top -p $(pgrep silentcast)   # Linux
Activity Monitor             # macOS
Task Manager                 # Windows

# プロファイリング有効化
./silentcast --enable-profiling --profile-port 6060
```

### 一般的なパフォーマンス問題
1. **高CPU使用率** → ホットキー監視、ファイル監視
2. **高メモリ使用量** → メモリリーク、バッファサイズ
3. **遅い応答時間** → アクション実行、I/O待機
4. **バッテリー消費** → ポーリング間隔、最適化設定

## 📊 パフォーマンス監視

### 内蔵メトリクス

#### システム情報の取得
```bash
# SilentCast システム情報
./silentcast --system-info

# 出力例:
# Version: 0.1.0-alpha.8
# Go Version: go1.19
# Platform: linux/amd64
# Memory Usage: 15.2 MB
# Goroutines: 12
# GC Cycles: 3
```

#### リアルタイム監視
```bash
# デバッグモードで詳細ログ
./silentcast --debug --performance-monitor

# プロファイリングサーバー
./silentcast --enable-profiling
# http://localhost:6060/debug/pprof/ でアクセス
```

### 外部監視ツール

#### システムレベル監視
```bash
# CPU とメモリ使用量
ps aux | grep silentcast

# ファイルディスクリプタ
lsof -p $(pgrep silentcast)

# ネットワーク接続
netstat -p | grep silentcast
```

#### Go 固有の監視
```bash
# メモリプロファイル
go tool pprof http://localhost:6060/debug/pprof/heap

# CPUプロファイル
go tool pprof http://localhost:6060/debug/pprof/profile

# ゴルーチン分析
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

## 🔧 パフォーマンス最適化

### 設定による最適化

#### パフォーマンス設定の例
```yaml
# spellbook.yml
performance:
  enable_optimization: true
  buffer_size: 2048        # バッファサイズを増加
  gc_percent: 75           # GC頻度を調整
  max_idle_time: 10m       # アイドル時間を延長
  enable_profiling: false  # 本番では無効化

daemon:
  log_level: warn          # ログレベルを下げる
  config_watch: true       # 必要に応じて無効化

logger:
  level: warn              # 詳細ログを無効化
  max_size: 10
  max_backups: 3
  compress: true

hotkeys:
  timeout: 1000ms          # タイムアウトを調整

notification:
  sound: false             # 不要な音声を無効化
  max_output_length: 512   # 出力サイズを制限
```

### CPU使用率の最適化

#### ホットキー監視の最適化
```yaml
hotkeys:
  timeout: 2000ms          # タイムアウトを長くしてCPU負荷軽減
  
# 不要なキー組み合わせを削除
spells:
  # 使用頻度の高いもののみ残す
  e: editor
  t: terminal
```

#### ファイル監視の最適化
```yaml
daemon:
  config_watch: false      # 頻繁に変更しない場合は無効化
```

### メモリ使用量の最適化

#### ガベージコレクション調整
```yaml
performance:
  gc_percent: 50           # メモリ使用量を抑制（CPU負荷増）
  # または
  gc_percent: 200          # CPU負荷を抑制（メモリ使用量増）
```

#### バッファサイズ調整
```yaml
performance:
  buffer_size: 512         # メモリ制約がある場合
  max_idle_time: 2m        # リソースを早期解放
```

#### ログローテーション
```yaml
logger:
  max_size: 5              # ログファイルサイズを制限
  max_backups: 2           # バックアップファイル数を制限
  max_age: 3               # ログ保持期間を短縮
```

## 🐛 メモリリーク診断

### ゴルーチンリーク検出

#### ゴルーチン数の監視
```bash
# ゴルーチン数の定期確認
while true; do
  curl -s http://localhost:6060/debug/pprof/goroutine?debug=1 | grep "^goroutine profile:"
  sleep 10
done
```

#### ゴルーチンリークの調査
```bash
# ゴルーチンプロファイル取得
go tool pprof http://localhost:6060/debug/pprof/goroutine

# プロファイル内で top コマンド実行
(pprof) top
(pprof) list main.main
(pprof) web
```

### メモリリーク検出

#### ヒープ分析
```bash
# ヒープスナップショット取得
go tool pprof http://localhost:6060/debug/pprof/heap

# メモリ使用量分析
(pprof) top
(pprof) list
(pprof) tree
```

#### メモリ使用量のトレース
```bash
# アロケーション分析
go tool pprof http://localhost:6060/debug/pprof/allocs

# 使用中メモリ分析  
go tool pprof http://localhost:6060/debug/pprof/heap
```

## ⚡ アクション実行の最適化

### 並行実行の制限

```yaml
# 同時実行数を制限してリソース使用量を抑制
grimoire:
  heavy-task:
    type: script
    command: heavy-process.sh
    timeout: 30              # タイムアウトを設定
    keep_open: false         # ターミナルを閉じる
```

### キャッシュの活用

```yaml
# パス解決結果のキャッシュ（内部で自動実行）
grimoire:
  editor:
    type: app
    command: code            # パス解決結果がキャッシュされる
```

### 出力バッファリング

```yaml
# 出力サイズを制限
notification:
  max_output_length: 256     # デフォルト1024を削減

grimoire:
  log-command:
    type: script
    command: some-verbose-command
    show_output: false       # 不要な出力を無効化
```

## 🔋 バッテリー寿命の最適化

### ポーリング間隔の調整

```yaml
# ファイル監視の間隔調整（内部設定）
daemon:
  config_watch: false        # バッテリー寿命を優先する場合

# 必要最小限のホットキーのみ使用
spells:
  e: editor                  # 頻繁に使用するもののみ
```

### 不要な機能の無効化

```yaml
# 通知音を無効化
notification:
  sound: false

# システムトレイを無効化
daemon:
  tray: false

# 自動更新を無効化
updater:
  enabled: false
```

## 📈 ベンチマークとテスト

### パフォーマンステスト実行

```bash
# ベンチマークスイート実行
cd app
make benchmark

# 特定コンポーネントのベンチマーク
make benchmark-action
make benchmark-config
make benchmark-hotkey
```

### カスタムベンチマーク

```bash
# アクション実行速度測定
time echo "test" | ./silentcast --stdin-action editor

# 設定読み込み速度測定  
time ./silentcast --validate-config
```

### 負荷テスト

```bash
# 連続実行テスト
for i in {1..100}; do
  ./silentcast --dry-run --config test.yml
done

# 並行実行テスト
seq 1 10 | xargs -P 10 -I {} ./silentcast --validate-config
```

## 🔍 プロファイリング

### CPU プロファイリング

```bash
# CPUプロファイル取得（30秒間）
go tool pprof -seconds 30 http://localhost:6060/debug/pprof/profile

# フレームグラフ生成
go tool pprof -http=:8080 profile.pb.gz
```

### メモリプロファイリング

```bash
# ヒーププロファイル
go tool pprof http://localhost:6060/debug/pprof/heap

# アロケーションプロファイル
go tool pprof http://localhost:6060/debug/pprof/allocs
```

### トレース分析

```bash
# 実行トレース取得
curl http://localhost:6060/debug/pprof/trace?seconds=10 > trace.out

# トレース分析
go tool trace trace.out
```

## 📋 パフォーマンス診断チェックリスト

### システムリソース
- [ ] CPU使用率が常に高くない（< 5%アイドル時）
- [ ] メモリ使用量が適切（< 50MB通常時）  
- [ ] ファイルディスクリプタ数が適正
- [ ] ネットワーク接続数が最小限

### 設定最適化
- [ ] 不要なログ出力を無効化
- [ ] パフォーマンス設定が適用されている
- [ ] 使用頻度の低いアクションを削除
- [ ] タイムアウト値が適切

### メモリ管理
- [ ] ゴルーチン数が安定している
- [ ] メモリリークが発生していない
- [ ] GCが適切な頻度で実行されている
- [ ] バッファプールが効果的に使用されている

## 🆘 深刻なパフォーマンス問題

### 緊急対処法

```bash
# プロセス優先度を下げる
renice +10 $(pgrep silentcast)

# CPUコア数を制限
taskset -c 0 ./silentcast

# メモリ制限
ulimit -v 104857600  # 100MBに制限
```

### プロファイリングデータ収集

```bash
# 包括的なデバッグ情報収集
mkdir perf-debug
cd perf-debug

# CPUプロファイル
go tool pprof -seconds 30 -output cpu.prof http://localhost:6060/debug/pprof/profile

# ヒーププロファイル  
go tool pprof -output heap.prof http://localhost:6060/debug/pprof/heap

# ゴルーチン情報
curl http://localhost:6060/debug/pprof/goroutine?debug=1 > goroutines.txt

# システム情報
./silentcast --system-info > system.txt
```

## 🔗 関連リソース

- [パフォーマンス最適化ガイド](../guide/performance-optimization.md)
- [開発ドキュメント](../development/)
- [システム要件](../guide/installation.md#system-requirements)
- [プロファイリング設定例](../../examples/config/performance_example.yml)