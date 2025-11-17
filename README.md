# LAN Speed Tester

Go言語で実装したネットワーク速度測定ツールです。ダウンロード速度とアップロード速度を測定できます。

## 特徴

- **ダウンロード速度測定**: サーバーからデータをダウンロードして速度を計測
- **アップロード速度測定**: サーバーにデータをアップロードして速度を計測
- **並列処理**: 複数の接続を使用して精度を向上
- **視覚的なCLI**: プログレスインジケーターとボックス描画で見やすい出力
- **クライアント/サーバー分離**: 柔軟なアーキテクチャ

## 必要要件

- Go 1.20以上

## インストール

### 1. リポジトリをクローン

```bash
git clone https://github.com/JoshuaSAKAMOTO/LAN-Speed-Checker.git
cd LAN-Speed-Checker
```

または、既にクローン済みの場合：

```bash
cd ~/lan-speed-tester
```

### 2. 依存関係の確認

このプロジェクトは**Go標準ライブラリのみ**を使用しているため、追加のパッケージインストールは不要です。

```bash
go mod download  # 念のため実行（何もダウンロードされません）
```

### 3. Goのバージョン確認

Go 1.20以上が必要です：

```bash
go version
# 例: go version go1.25.4 darwin/arm64
```

## 使用方法

### ステップバイステップガイド

#### 1. ターミナルを2つ開く

このツールはサーバーとクライアントを**別々のターミナルウィンドウ**で実行する必要があります。

- **ターミナル1**: サーバー用（常に起動したまま）
- **ターミナル2**: クライアント用（速度測定を実行）

#### 2. サーバーを起動する（ターミナル1）

```bash
cd ~/lan-speed-tester
go run cmd/server/main.go
```

**成功すると以下のように表示されます：**
```
2025/11/17 14:28:38 Starting LAN Speed Tester Server on :8080
2025/11/17 14:28:38 Endpoints:
2025/11/17 14:28:38   - Health Check: http://localhost:8080/
2025/11/17 14:28:38   - Download Test: http://localhost:8080/download
2025/11/17 14:28:38   - Upload Test: http://localhost:8080/upload
```

**⚠️ このターミナルはそのまま開いたままにしておいてください！**

#### 3. クライアントを実行する（ターミナル2）

**新しいターミナルウィンドウを開いて**以下を実行：

```bash
cd ~/lan-speed-tester
go run cmd/client/main.go
```

**実行結果の例：**
```
╔═══════════════════════════════════╗
║   LAN Speed Tester Client        ║
╚═══════════════════════════════════╝
Server: http://localhost:8080

Connecting to server... ✓

┌─────────────────────────────────┐
│  Single Connection Test         │
└─────────────────────────────────┘
  Downloading test data... ✓
  Download: 20248.89 Mbps

  Uploading test data... ✓
  Upload:   14052.33 Mbps

┌─────────────────────────────────┐
│  Parallel Test (4 connections) │
└─────────────────────────────────┘
  Downloading with 4 parallel connections... ✓
  Download: 17440.95 Mbps

  Uploading with 4 parallel connections... ✓
  Upload:   25647.16 Mbps

╔═══════════════════════════════════╗
║      Test Results Summary        ║
╠═══════════════════════════════════╣
║ Single Connection:               ║
║   Download:  20248.89 Mbps      ║
║   Upload:    14052.33 Mbps      ║
║                                  ║
║ Parallel (4 connections):       ║
║   Download:  17440.95 Mbps      ║
║   Upload:    25647.16 Mbps      ║
╚═══════════════════════════════════╝
```

#### 4. 測定時間について

測定には**約30秒〜1分**程度かかります。以下のデータ転送が行われるためです：

- ダウンロードテスト（シングル）: 100MB
- アップロードテスト（シングル）: 50MB
- ダウンロードテスト（並列4接続）: 100MB × 4 = 400MB
- アップロードテスト（並列4接続）: 50MB × 4 = 200MB

**合計: 約750MB**

進行状況はアニメーション付きスピナーで表示されます。

### 利用可能なエンドポイント

サーバーが起動すると、以下のエンドポイントが利用可能になります：

- `http://localhost:8080/` - ヘルスチェック
- `http://localhost:8080/download` - ダウンロード速度測定
- `http://localhost:8080/upload` - アップロード速度測定

### サーバーの停止方法

ターミナル1（サーバー）で **Ctrl+C** を押すと、サーバーが停止します。

## プロジェクト構造

```
lan-speed-tester/
├── cmd/
│   ├── server/
│   │   └── main.go    # サーバー実装
│   └── client/
│       └── main.go    # クライアント実装
├── go.mod             # Go モジュールファイル
└── README.md          # このファイル
```

## 実装の詳細

### サーバー側

- HTTPサーバーとして実装
- ランダムデータを生成してクライアントに送信（ダウンロードテスト）
- クライアントからのデータを受信（アップロードテスト）

### クライアント側

- 単一接続と並列接続（デフォルト4接続）の両方で測定
- ゴルーチンを使用した並列処理
- プログレスインジケーターによる視覚的フィードバック
- 結果をMbps単位で表示

## 設定

### 並列接続数の変更

`cmd/client/main.go`の`parallelConnections`定数を変更してください：

```go
const (
    defaultServerURL   = "http://localhost:8080"
    parallelConnections = 4 // ここを変更
)
```

### サーバーポートの変更

`cmd/server/main.go`の`defaultPort`定数を変更してください：

```go
const (
    defaultPort = "8080" // ここを変更
)
```

## 技術スタック

- **言語**: Go 1.25.4
- **標準ライブラリのみ使用**:
  - `net/http` - HTTPサーバー/クライアント
  - `sync` - 並列処理の同期
  - `time` - 時間計測
  - `crypto/rand` - ランダムデータ生成

## トラブルシューティング

### `command not found: go`

**問題**: `go`コマンドが見つからない

**解決方法**:

Goのインストールパスを確認して、パスを通してください：

```bash
# Goがインストールされているか確認
/usr/local/go/bin/go version

# パスを一時的に通す
export PATH=$PATH:/usr/local/go/bin

# 永続的にパスを通す（推奨）
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
source ~/.zshrc
```

### `connection refused`

**問題**: `dial tcp [::1]:8080: connect: connection refused`

**原因**: サーバーが起動していない

**解決方法**:

1. サーバーが起動しているか確認：
   ```bash
   lsof -i :8080
   ```

2. サーバーを起動してから、クライアントを実行：
   ```bash
   # ターミナル1
   cd ~/lan-speed-tester
   go run cmd/server/main.go

   # ターミナル2（新しいターミナル）
   cd ~/lan-speed-tester
   go run cmd/client/main.go
   ```

### クライアントがハングする

**問題**: クライアントの実行が途中で止まる

**原因**: レスポンスボディの読み取りが完了していない

**解決方法**:

最新版のコードを使用していることを確認してください。古いバージョンではHTTPレスポンスボディの読み取りに関するバグがありました。

```bash
git pull origin feature/lan-speed-tester
```

### `stat cmd/server/main.go: no such file or directory`

**問題**: ファイルが見つからない

**原因**: プロジェクトディレクトリにいない

**解決方法**:

プロジェクトディレクトリに移動してから実行：

```bash
cd ~/lan-speed-tester
go run cmd/server/main.go
```

### Xcode license エラー

**問題**: `You have not agreed to the Xcode license agreements`

**解決方法**:

```bash
sudo xcodebuild -license
# スペースキーで最後まで読み、"agree"と入力
```

## ビルド

実行可能ファイルをビルドするには：

```bash
# サーバー
go build -o server cmd/server/main.go

# クライアント
go build -o client cmd/client/main.go
```

ビルド後は、以下のように実行できます：

```bash
# サーバー
./server

# クライアント（別のターミナルで）
./client
```

## 今後の拡張案

- [ ] コマンドライン引数でサーバーURLやポートを指定
- [ ] 測定結果のログファイル出力
- [ ] JSON形式での結果出力
- [ ] カスタムデータサイズの指定
- [ ] グラフィカルな結果表示
- [ ] Ping測定機能の追加

## ライセンス

MIT License

## 参考

このプロジェクトは以下のチュートリアルを参考に作成されました：
- [Go言語で作る！回線速度測定ツール：LAN Speed Tester](https://zenn.dev/haruki1009/books/7d3b6cfaac560d)
