package templates

import (
	"embed"
	"html/template"
	"log/slog"
)

//go:embed html/reset_password.html
var ResetPasswordFS embed.FS

var ResetPasswordTmpl *template.Template

type ResetPasswordTmplData struct {
	Name  string
	Email string
	Link  string
}

func init() {
	slog.Info("Initializing reset_password.html as template.")

	ResetPasswordTmpl = template.Must(
		template.ParseFS(ResetPasswordFS, "html/reset_password.html"),
	)

	slog.Info("Initializing reset_password.html as template. [OK]")
}
