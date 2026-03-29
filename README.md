# Push Link v2

HTMX をフロント、Go をバックエンドに採用したリンク管理アプリです。  
一覧 UI は Qiita のような情報密度を意識し、バックエンドは Go の標準的なレイヤ構成で整理しています。

## 現在の実装範囲

- HTMX による一覧更新
- Go の `net/http` ベース API
- Strategy パターンによる一覧ロジック切り替え
- 月次登録数の棒グラフ表示
- 共通カラートークンの分離
- メモリ上のサンプルデータによる画面確認

現時点では DB 接続は未実装で、`internal/repository` のメモリリポジトリを使っています。  
MySQL スキーマ案は [docs/db-design.md](/Users/toyo/Terminal/push-link-v2/docs/db-design.md) に整理済みです。

## ディレクトリ構成

- `cmd/server`
  - アプリケーションのエントリーポイント
- `internal/domain`
  - エンティティとレスポンスモデル
- `internal/repository`
  - データ取得層
- `internal/usecase`
  - アプリケーションロジック
- `internal/service`
  - Strategy パターン実装
- `internal/handler`
  - HTTP ハンドラ
- `internal/view`
  - Go テンプレート描画
- `web/templates`
  - HTMX 用 HTML テンプレート
- `web/static`
  - 画面スタイル
- `web/static/common`
  - 共通カラートークン
- `docs`
  - 設計ドキュメント

## 画面と API

- `GET /`
  - 一覧画面
- `GET /ui/sites`
  - HTMX 用部分更新
- `GET /api/v1/sites`
  - JSON API
- `GET /healthz`
  - ヘルスチェック

### API クエリ例

```text
/api/v1/sites?strategy=trending&tag=development&q=docs
```

### 利用可能な Strategy

- `default`
  - 通常一覧
- `published`
  - `published` のみを表示
- `trending`
  - ブックマーク数優先で表示

Go ではクラスではなくインターフェースを使って Strategy を表現しています。

## 起動方法

```bash
go run ./cmd/server
```

デフォルトの起動ポートは `8080` です。  
ブラウザで `http://localhost:8080` を開くと画面を確認できます。

## 検証

```bash
go build ./...
```

## ドキュメント

- [docs/app-architecture.md](/Users/toyo/Terminal/push-link-v2/docs/app-architecture.md)
- [docs/ui-design.md](/Users/toyo/Terminal/push-link-v2/docs/ui-design.md)
- [docs/db-design.md](/Users/toyo/Terminal/push-link-v2/docs/db-design.md)
