package templates

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"time"
)

//go:embed html/reset_password.html
var ResetPasswordFS embed.FS

var ResetPasswordTmpl *template.Template

type ResetPasswordTmplData struct {
	Username    string
	ResetURL    string
	ResetURLTTL string
	Year        string
	AppName     string
}

func NewResetPasswordTmplData(username, resetURL, resetURLTTL string) *ResetPasswordTmplData {
	return &ResetPasswordTmplData{
		Username:    username,
		ResetURL:    resetURL,
		ResetURLTTL: resetURLTTL,
		Year:        fmt.Sprintf("%d", time.Now().Year()),
		AppName:     "CodeUniverse",
	}
}

func init() {
	slog.Info("Initializing reset_password.html as template.")

	ResetPasswordTmpl = template.Must(
		template.ParseFS(ResetPasswordFS, "html/reset_password.html"),
	)

	slog.Info("Initializing reset_password.html as template. [OK]")
}
