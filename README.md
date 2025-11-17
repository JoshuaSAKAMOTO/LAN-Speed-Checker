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

```bash
git clone <repository-url>
cd lan-speed-tester
go mod download
```

## 使用方法

### サーバーの起動

別のターミナルウィンドウでサーバーを起動します：

```bash
go run cmd/server/main.go
```

サーバーはデフォルトでポート8080で起動します。以下のエンドポイントが利用可能です：

- `http://localhost:8080/` - ヘルスチェック
- `http://localhost:8080/download` - ダウンロード速度測定
- `http://localhost:8080/upload` - アップロード速度測定

### クライアントの実行

サーバーが起動している状態で、クライアントを実行します：

```bash
go run cmd/client/main.go
```

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

## ビルド

実行可能ファイルをビルドするには：

```bash
# サーバー
go build -o server cmd/server/main.go

# クライアント
go build -o client cmd/client/main.go
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
