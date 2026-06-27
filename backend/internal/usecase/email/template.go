package email

import "text/template"

type WelcomeTemplateData struct {
	UserName string
}

var welcomeSubject = "Welcome to Todo App!"

var welcomeBodyTmpl = template.Must(template.New("welcome").Parse(
	`{{.UserName}} さん、ご登録ありがとうございます。

Todo App へようこそ！
このアプリでタスクを管理して、生産性を高めましょう。

添付の QR コードからアプリにアクセスできます。

今後ともよろしくお願いいたします。
Todo App チーム`,
))
