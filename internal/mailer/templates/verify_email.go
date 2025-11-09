package templates

import (
	"embed"
	"html/template"
	"log/slog"
)

//go:embed html/verify_email.html
var verifyEmailFS embed.FS

var VerifyEmailTmpl *template.Template

type VerifyEmailTmplData struct {
	Name  string
	Email string
	Link  string
}

func init() {
	slog.Info("Initializing verify_email.html as template.")

	VerifyEmailTmpl = template.Must(
		template.ParseFS(verifyEmailFS, "html/verify_email.html"),
	)

	slog.Info("Initializing verify_email.html as template. [OK]")
}
