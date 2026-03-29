# アプリケーション設計

## 1. 目的

Push Link v2 のアプリケーション構成と、Go 側の責務分離を整理するためのドキュメントです。

## 2. 技術方針

- フロントは HTMX を採用する
- バックエンドは Go の `net/http` を採用する
- API は `GET /api/v1/...` 形式でバージョニングする
- 画面描画は Go テンプレートで行う
- 一覧取得ロジックは Strategy パターンで差し替え可能にする

## 3. 全体構造

アプリケーション全体の依存関係は以下の形です。

```text
[cmd/server]
  サーバ起動
  依存関係の組み立て
        |
        v
[handler]
  HTTP リクエスト受付
  クエリ受取
  HTML/JSON 返却
        |
        v
[usecase]
  一覧取得の業務ロジック
  Strategy 選択
  レスポンス整形
        |
        +--------------------+
        |                    |
        v                    v
[repository]           [service]
  データ取得             Strategy 群
  現状はメモリ実装       default/published/trending
        |                    |
        +----------+---------+
                   |
                   v
                [domain]
         Site, Filter, Result など共通モデル
```

ポイントは以下です。

- `handler` は HTTP の入出力だけを扱う
- `usecase` は業務ロジックの中心であり、画面と API の両方から利用される
- `repository` はデータ取得を担当し、保存先の詳細を外へ漏らさない
- `service` は差し替え可能な戦略ロジックを担当する
- `domain` は各レイヤで共有するモデルを持つ

## 4. レイヤ構成

### 4.1 `cmd/server`

- サーバ起動
- 依存関係の組み立て
- ルーティング登録

[main.go](/Users/toyo/Terminal/push-link-v2/cmd/server/main.go) では以下を実施しています。

- テンプレートレンダラ生成
- リポジトリ生成
- Strategy レジストリ生成
- ユースケース生成
- ハンドラ生成
- ルーティング登録

このファイルは DI の起点であり、各レイヤの組み立てだけを行います。

### 4.2 `internal/domain`

- `Site`
- `SiteFilter`
- `SiteListResult`
- `MonthlyRegistration`

ドメインモデルと、画面・API で共有するレスポンスモデルをここに集約しています。

主な役割は以下です。

- `Site`
  - 1 件のサイト情報
- `SiteFilter`
  - 一覧取得時の検索条件
- `SiteListResult`
  - 一覧表示や JSON API の返却モデル
- `MonthlyRegistration`
  - 月次棒グラフ用の集計値

### 4.3 `internal/repository`

- `SiteRepository`
- `MemorySiteRepository`

現状はサンプルデータを返すメモリ実装です。  
将来的に MySQL 実装へ差し替える前提で、リポジトリはインターフェースで抽象化しています。

今の設計では、ユースケースは `SiteRepository` だけを知っていればよく、メモリか SQL かは意識しません。

### 4.4 `internal/service`

- `SiteStrategy`
- `SiteStrategyRegistry`
- `DefaultSiteStrategy`
- `PublishedSiteStrategy`
- `TrendingSiteStrategy`

一覧の並び順や抽出条件の違いを Strategy として分離しています。  
追加戦略は `SiteStrategy` を実装してレジストリへ登録すれば拡張できます。

### 4.5 `internal/usecase`

- 一覧取得
- タグ一覧取得
- 月次登録数集計

リポジトリから取得したデータを、画面と API の両方で使える形に組み立てる責務を持ちます。

`usecase` は以下の流れで処理します。

1. リポジトリからサイト一覧を取得する
2. リポジトリからタグ一覧を取得する
3. Strategy 名を解決する
4. Strategy を適用して一覧を作る
5. 月次登録数を集計する
6. 画面と API が共有できる `SiteListResult` を返す

### 4.6 `internal/handler`

- `GET /`
- `GET /ui/sites`
- `GET /api/v1/sites`
- `GET /healthz`

HTTP リクエストからフィルタ条件を取り出し、ユースケースを呼び出し、HTML または JSON を返します。

画面系と API 系で責務は以下のように分かれます。

- `Index`
  - 初期ページ表示
- `SiteListPartial`
  - HTMX による一覧部分更新
- `SiteListAPI`
  - JSON を返す API
- `Healthz`
  - 死活監視用

### 4.7 `internal/view`

- テンプレート読み込み
- テンプレート関数登録

現状は `formatDate` と `add1` をテンプレート関数として提供しています。

## 5. リクエストフロー

一覧画面の主な流れは以下です。

```text
Browser
  -> GET /
Handler
  -> query を SiteFilter に詰める
Usecase
  -> Repository から sites と tags を取得
  -> StrategyRegistry で戦略を解決
  -> Strategy で絞り込み・並び替え
  -> 月次登録数を集計
Handler
  -> HTML テンプレートへ渡して返却
```

HTMX 更新時は `GET /ui/sites` を使い、画面全体ではなく一覧領域だけを差し替えます。  
JSON API の場合は同じユースケースを通し、返却形式だけを JSON に変えます。

## 6. Strategy パターンの考え方

このプロジェクトでの Strategy パターンは、一般的な OOP のクラス設計ではなく、Go のインターフェース設計で実現しています。

- `default`
  - 基本的なフィルタ処理とタイトル順表示
- `published`
  - 公開済みのみを抽出
- `trending`
  - ブックマーク数を優先して並び替え

`usecase` は戦略名を受け取り、`service` のレジストリから実装を解決します。

Go での表現は以下の考え方です。

- Java や PHP のような継承ベースのクラス設計ではない
- `interface` で振る舞いを定義する
- `struct` で具体実装を持つ
- 必要な依存だけをコンストラクタ関数で注入する

戦略追加の手順は以下です。

1. `SiteStrategy` を満たす新しい型を作る
2. `Name()` と `Apply()` を実装する
3. `NewSiteStrategyRegistry()` に登録する
4. 必要なら UI 側の選択肢を追加する

## 7. API 設計

### 7.1 エンドポイント

- `GET /api/v1/sites`

### 7.2 クエリパラメータ

- `strategy`
- `q`
- `tag`
- `status`

### 7.3 レスポンス内容

- サイト一覧
- 適用中 Strategy
- 利用可能タグ一覧
- 月次登録数

画面 HTML と API JSON の両方で同じユースケース結果を利用することで、表示ロジックの重複を避けています。

## 8. クラス設計に相当する考え方

Go には典型的な意味でのクラスはありません。  
このプロジェクトでは以下の組み合わせで、クラス設計に相当する構造を作っています。

- データは `struct`
- 契約は `interface`
- 初期化は `New...` 関数
- 差し替えはインターフェース越しに行う

対応関係の例は以下です。

- `SiteUsecase`
  - Java/PHP でいうサービスクラスに近い
- `SiteRepository`
  - リポジトリインターフェース
- `MemorySiteRepository`
  - リポジトリ実装クラスに近い役割
- `SiteStrategy`
  - Strategy のインターフェース
- `DefaultSiteStrategy` など
  - Strategy の具体実装
- `SiteHandler`
  - コントローラに近い役割

## 9. 今後の実装候補

- MySQL 接続実装
- `repository` の SQL 実装追加
- `POST /api/v1/sites` など登録 API の追加
- 認証と権限制御
- ページネーション
- 並び替え条件の追加
