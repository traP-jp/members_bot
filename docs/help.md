## GitHub Organization 管理Bot

GitHubの {{ .ORG_NAME }} Organizationのメンバーを管理するためのtraQ botです。
`@{{ .BOT_NAME }} <コマンド> [引数(任意)]` のように使います。

### `/invite` (`@{{ .BOT_NAME }} /invite <traQID1> <GitHubID1> ...`)

Organizationへの招待を申請するコマンドです。
Organizationのadminのグループにメンションが飛び、一定数のスタンプがついたら承認・却下されます。
現在は承認は{{ .ACCEPT_STAMP_THRESHOLD }}個、却下は{{ .REJECT_STAMP_THRESHOLD }}個に設定されています。adminに承認されると招待が送られます。

### `/list` (`@{{ .BOT_NAME }} /list`)

現在の申請状態を示します。

### `/help` (`@{{ .BOT_NAME }} /help`)

この文章を表示します。
