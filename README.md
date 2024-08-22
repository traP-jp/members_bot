# members_bot

Organizationメンバー管理用のtraQ bot

## デプロイ

GoとMySQL(MariaDB)

### GitHub App

Organizationのメンバーの管理に、GitHub Appを使う。

[GitHub App インストールとしての認証](https://docs.github.com/ja/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app-installation#using-an-installation-access-token-to-authenticate-as-an-app-installation)

GitHub Appを作り、Organizationにinstallする。OrganizationのmembersのRead/Write権限を持たせておく。また、Private Keyをダウンロードしておく。

GitHub AppのInstallation IDが必要になるが、GitHubのUIからは確認できない。リポジトリルートの`installation_id.sh`を実行すると取得できる。

```sh
./installation_id.sh {GitHubAppのClient ID} {Private Keyのパス} {Org名}
```

(jqを使っているが、無い場合は最後の行をコメントアウトして`id`を見ればよい。)

### 環境変数

- `ACCEPT_STAMP_ID` 承認用スタンプのUUID
- `ACCEPT_STAMP_THRESHOLD` 何個スタンプがついたら承認とするか
- `ADMIN_GROUP_ID` adminのtraQ Group UUID
- `ADMIN_GROUP_NAME` adminのtraQ Group名
- `BOT_CHANNEL_ID` botが投稿するチャンネル
- `GITHUB_APP_ID` GitHub AppのID
- `GITHUB_APP_INSTALLATION_ID` GitHub AppのInstallation ID
- `GITHUB_APP_PRIVATE_KEY` GitHub Appの秘密鍵。改行を`\n`に置き変えたもの。
- `GITHUB_ORG_NAME` GitHubのオーガニゼーション名
<!-- - `GITHUB_TOKEN` GitHubのトークン -->
- `INACTIVE_STAMP_ID` 操作を終えたメッセージに押すスタンプのUUID
- `REJECT_STAMP_ID` 却下用スタンプのUUID
- `REJECT_STAMP_THRESHOLD` 何個スタンプがついたら却下とするか
- `TRAQ_BOT_TOKEN` traQのBot token
- `NS_MARIADB_DATABASE`, `MYSQL_DATABASE` (default: `members_bot`) DBのデータベース名。NS_の方が優先される。
- `NS_MARIADB_HOSTNAME`, `MYSQL_HOSTNAME` (default: `db`) DBのホスト名。NS_の方が優先される。
- `NS_MARIADB_PASSWORD`, `MYSQL_PASSWORD` (default `pass`) DBのパスワード。NS_の方が優先される。
- `NS_MARIADB_PORT`, `MYSQL_PORT` (default `3306`) DBのポート番号。NS_の方が優先される。
- `NS_MARIADB_USER`, `MYSQL_USER` (default: `root`) DBのユーザー。NS_の方が優先される。

## 開発

- Go
- Docker

リポジトリルートに上の環境変数を書いた`.env`を置く。[`.env.sample`](./.env.sample)を参考に。

docker compose watchを使ってホットリロードにしている。起動時は

```sh
docker compose watch
```

ログを見たいときは

```sh
docker compose logs
```

### テスト

`handler`パッケージと`repository/impl`パッケージで、ユニットテストを書いている。

- `handler`パッケージでは、serviceとrepositoryのmockとして、[matryer/moq](https://github.com/matryer/moq) を使っている。mockはGit管理に含めていないので、初めてテストを実行するときは`go generate ./...`でmockを生成する。interface定義を変えたときもmock生成が必要である。
- `repository/impl`パッケージでは、[testcontainers/testcontainers-go](https://github.com/testcontainers/testcontainers-go)でDockerコンテナを使ったDB操作のテストを書いている。テストを実行する際はDockerが必要である。

`service/impl/github_test.go`[service/impl/github_test.go]もあるが、このテストはGitHubのトークンが必要なので、デフォルトでは実行されない。`go test`の引数に`-tags github_env`を含めると実行される。
