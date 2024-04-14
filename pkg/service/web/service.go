package web

import (
	"embed"
)

//go:embed web/template/*.htm
var templates embed.FS

//go:embed web/static/*.css
var static embed.FS
