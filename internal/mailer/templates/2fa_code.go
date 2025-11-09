package templates

import (
	"embed"
	"html/template"
	"log/slog"
)

//go:embed html/2fa_code.html
var twoFAFS embed.FS

var TwoFATmpl *template.Template

type TwoFATmplData struct {
	Name  string
	Email string
	Code  string
}

func init() {
	slog.Info("Initializing 2fa_code.html as template.")

	TwoFATmpl = template.Must(
		template.ParseFS(twoFAFS, "html/2fa_code.html"),
	)

	slog.Info("Initializing 2fa_code.html as template. [OK]")
}
