package templates

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"time"
)

//go:embed html/2fa_code.html
var twoFAFS embed.FS

var TwoFATmpl *template.Template

type TwoFATmplData struct {
	Username string
	Code     string
	CodeTTL  string
	Year     string
	AppName  string
}

func NewTwoFATmplData(username, code, codeTTL string) *TwoFATmplData {
	return &TwoFATmplData{
		Username: username,
		Code:     code,
		CodeTTL:  codeTTL,
		Year:     fmt.Sprintf("%d", time.Now().Year()),
		AppName:  "CodeUniverse",
	}
}

func init() {
	slog.Info("Initializing 2fa_code.html as template.")

	TwoFATmpl = template.Must(
		template.ParseFS(twoFAFS, "html/2fa_code.html"),
	)

	slog.Info("Initializing 2fa_code.html as template. [OK]")
}
