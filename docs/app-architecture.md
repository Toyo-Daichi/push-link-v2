# アプリケーション設計

## 1. 目的

Push Link v2 のアプリケーション構成と、Go 側の責務分離を整理するためのドキュメントです。

## 2. 技術方針

- フロントは HTMX を採用する
- バックエンドは Go の `net/http` を採用する
- API は `GET /api/v1/...` 形式でバージョニングする
- 画面描画は Go テンプレートで行う
- 一覧取得ロジックは Strategy パターンで差し替え可能にする

## 3. レイヤ構成

### 3.1 `cmd/server`

- サーバ起動
- 依存関係の組み立て
- ルーティング登録

### 3.2 `internal/domain`

- `Site`
- `SiteFilter`
- `SiteListResult`
- `MonthlyRegistration`

ドメインモデルと、画面・API で共有するレスポンスモデルをここに集約しています。

### 3.3 `internal/repository`

- `SiteRepository`
- `MemorySiteRepository`

現状はサンプルデータを返すメモリ実装です。  
将来的に MySQL 実装へ差し替える前提で、リポジトリはインターフェースで抽象化しています。

### 3.4 `internal/service`

- `SiteStrategy`
- `SiteStrategyRegistry`
- `DefaultSiteStrategy`
- `PublishedSiteStrategy`
- `TrendingSiteStrategy`

一覧の並び順や抽出条件の違いを Strategy として分離しています。  
追加戦略は `SiteStrategy` を実装してレジストリへ登録すれば拡張できます。

### 3.5 `internal/usecase`

- 一覧取得
- タグ一覧取得
- 月次登録数集計

リポジトリから取得したデータを、画面と API の両方で使える形に組み立てる責務を持ちます。

### 3.6 `internal/handler`

- `GET /`
- `GET /ui/sites`
- `GET /api/v1/sites`
- `GET /healthz`

HTTP リクエストからフィルタ条件を取り出し、ユースケースを呼び出し、HTML または JSON を返します。

### 3.7 `internal/view`

- テンプレート読み込み
- テンプレート関数登録

現状は `formatDate` と `add1` をテンプレート関数として提供しています。

## 4. Strategy パターンの考え方

このプロジェクトでの Strategy パターンは、一般的な OOP のクラス設計ではなく、Go のインターフェース設計で実現しています。

- `default`
  - 基本的なフィルタ処理とタイトル順表示
- `published`
  - 公開済みのみを抽出
- `trending`
  - ブックマーク数を優先して並び替え

`usecase` は戦略名を受け取り、`service` のレジストリから実装を解決します。

## 5. API 設計

### 5.1 エンドポイント

- `GET /api/v1/sites`

### 5.2 クエリパラメータ

- `strategy`
- `q`
- `tag`
- `status`

### 5.3 レスポンス内容

- サイト一覧
- 適用中 Strategy
- 利用可能タグ一覧
- 月次登録数

## 6. 今後の実装候補

- MySQL 接続実装
- `repository` の SQL 実装追加
- `POST /api/v1/sites` など登録 API の追加
- 認証と権限制御
- ページネーション
- 並び替え条件の追加
