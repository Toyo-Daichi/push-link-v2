# DB設計書

## 1. 目的

インターネット上の優良サイトを一元管理するためのMySQL用データベース設計です。  
単なるサイト一覧ではなく、以下を扱える構成にしています。

- サイト情報の管理
- タグによる分類
- 登録ユーザーの管理
- 誰がサイトを登録したかの追跡
- ユーザーごとのブックマーク管理

## 2. 設計方針

- 要望に合わせ、全テーブルに `create_at`, `update_at` を持たせる
- サイトとタグは多対多のため、中間テーブル `site_tags` を用意する
- サイトの登録者を明確にするため、`sites.added_by_user_id` を持たせる
- 運用を見据えて、サイトには `status` を持たせる
- ユーザーごとの保存サイトを管理できるよう `user_site_bookmarks` を追加する

## 3. テーブル一覧

| テーブル名 | 用途 |
| --- | --- |
| `users` | サービス利用ユーザー |
| `tags` | サイト分類用タグ |
| `sites` | 優良サイト本体 |
| `site_tags` | サイトとタグの関連 |
| `user_site_bookmarks` | ユーザーごとの保存サイト |

## 4. テーブル定義

### 4.1 `users`

| カラム名 | 型 | NULL | キー | 説明 |
| --- | --- | --- | --- | --- |
| `id` | `BIGINT UNSIGNED` | NO | PK | ユーザーID |
| `name` | `VARCHAR(100)` | NO |  | 表示名 |
| `email` | `VARCHAR(255)` | NO | UK | ログインや通知で使えるメールアドレス |
| `role` | `ENUM('admin','editor','viewer')` | NO |  | 権限区分 |
| `is_active` | `TINYINT(1)` | NO |  | 有効フラグ |
| `create_at` | `DATETIME` | NO |  | 作成日時 |
| `update_at` | `DATETIME` | NO |  | 更新日時 |

補足:
- `email` は一意制約を持たせる
- `role` は将来の運用権限の差分を吸収するために追加

### 4.2 `tags`

| カラム名 | 型 | NULL | キー | 説明 |
| --- | --- | --- | --- | --- |
| `id` | `BIGINT UNSIGNED` | NO | PK | タグID |
| `name` | `VARCHAR(100)` | NO | UK | 表示用タグ名 |
| `slug` | `VARCHAR(100)` | NO | UK | URLや内部識別向けの正規化名 |
| `create_at` | `DATETIME` | NO |  | 作成日時 |
| `update_at` | `DATETIME` | NO |  | 更新日時 |

補足:
- `name` と `slug` の両方を一意にすることで、表示名と内部名の双方を安定運用できる

### 4.3 `sites`

| カラム名 | 型 | NULL | キー | 説明 |
| --- | --- | --- | --- | --- |
| `id` | `BIGINT UNSIGNED` | NO | PK | サイトID |
| `title` | `VARCHAR(255)` | NO |  | サイト名 |
| `description` | `TEXT` | NO |  | サイト説明 |
| `url` | `VARCHAR(2048)` | NO | UK | サイトURL |
| `domain` | `VARCHAR(255)` | NO | INDEX | ドメイン |
| `status` | `ENUM('draft','published','archived')` | NO | INDEX | 公開状態 |
| `added_by_user_id` | `BIGINT UNSIGNED` | NO | FK | 登録したユーザーID |
| `create_at` | `DATETIME` | NO |  | 作成日時 |
| `update_at` | `DATETIME` | NO |  | 更新日時 |

補足:
- ご提示の `description`, `url`, `tags` に加え、実際の一覧表示や検索で必要になる `title` を追加
- URLの重複登録を防ぐため `url` は一意制約
- ドメイン単位の集計や重複確認のため `domain` を追加
- 登録フローを考慮し、`status` を追加

### 4.4 `site_tags`

| カラム名 | 型 | NULL | キー | 説明 |
| --- | --- | --- | --- | --- |
| `site_id` | `BIGINT UNSIGNED` | NO | PK, FK | サイトID |
| `tag_id` | `BIGINT UNSIGNED` | NO | PK, FK | タグID |
| `create_at` | `DATETIME` | NO |  | 作成日時 |
| `update_at` | `DATETIME` | NO |  | 更新日時 |

補足:
- `sites` と `tags` は多対多のため中間テーブルが必須
- 主キーを `(site_id, tag_id)` として重複紐付けを防止

### 4.5 `user_site_bookmarks`

| カラム名 | 型 | NULL | キー | 説明 |
| --- | --- | --- | --- | --- |
| `id` | `BIGINT UNSIGNED` | NO | PK | ブックマークID |
| `user_id` | `BIGINT UNSIGNED` | NO | FK | ユーザーID |
| `site_id` | `BIGINT UNSIGNED` | NO | FK | サイトID |
| `note` | `VARCHAR(255)` | YES |  | 個別メモ |
| `create_at` | `DATETIME` | NO |  | 作成日時 |
| `update_at` | `DATETIME` | NO |  | 更新日時 |

補足:
- サイトの一元管理に加え、ユーザーごとの保存や再訪管理ができる
- `(user_id, site_id)` に一意制約を置き、同じユーザーが同じサイトを重複保存しないようにしている

## 5. リレーション

- `users` 1 : N `sites`
  - 1人のユーザーが複数サイトを登録できる
- `sites` N : N `tags`
  - `site_tags` で関連付ける
- `users` 1 : N `user_site_bookmarks`
  - 1人のユーザーが複数サイトを保存できる
- `sites` 1 : N `user_site_bookmarks`
  - 1つのサイトを複数ユーザーが保存できる

## 6. 採用した追加属性・追加エンティティ

要望に対して追加したものは以下です。

| 対象 | 追加項目 | 理由 |
| --- | --- | --- |
| `sites` | `title` | サイト一覧表示時の主見出しとして必要 |
| `sites` | `domain` | 重複確認、検索、集計で有用 |
| `sites` | `status` | 下書き、公開、アーカイブを管理できる |
| `sites` | `added_by_user_id` | 誰が登録したか追跡できる |
| `tags` | `slug` | 内部処理やURL利用で安定した識別子になる |
| `users` | `email` | 実用的なユーザー識別に必要 |
| `users` | `role` | 管理者・編集者・閲覧者を分けられる |
| `users` | `is_active` | 無効化運用がしやすい |
| 新規 | `site_tags` | サイトとタグの多対多を表現するため必須 |
| 新規 | `user_site_bookmarks` | ユーザーごとの保存管理に有用 |

## 7. サンプルデータ方針

`init.sql` には以下のサンプルを含めています。

- ユーザー3件
- タグ5件
- サイト5件
- サイトとタグの関連10件
- ブックマーク4件

これにより、初期状態で一覧表示、タグ検索、登録者表示、ブックマーク表示の検証ができます。

## 8. 今後の拡張候補

必要になれば次を追加できます。

- `site_reviews`
  - サイトに対する評価やコメント管理
- `collections`
  - テーマごとのサイトまとめ
- `access_logs`
  - 閲覧履歴や人気分析
- `site_reports`
  - リンク切れや不適切情報の報告管理
