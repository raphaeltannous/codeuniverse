package templates

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"time"
)

//go:embed html/verify_email.html
var verifyEmailFS embed.FS

var VerifyEmailTmpl *template.Template

type VerifyEmailTmplData struct {
	Username  string
	Email     string
	VerifyURL string
	LinkTTL   string
	Year      string
	AppName   string
}

func NewVerifyEmailTmplData(username, email, verifyURL, linkTTL string) *VerifyEmailTmplData {
	return &VerifyEmailTmplData{
		Username:  username,
		Email:     email,
		VerifyURL: verifyURL,
		LinkTTL:   linkTTL,
		Year:      fmt.Sprintf("%d", time.Now().Year()),
		AppName:   "CodeUniverse",
	}
}

func init() {
	slog.Info("Initializing verify_email.html as template.")

	VerifyEmailTmpl = template.Must(
		template.ParseFS(verifyEmailFS, "html/verify_email.html"),
	)

	slog.Info("Initializing verify_email.html as template. [OK]")
}
