package main

import (
	"embed"

	_ "gioui.org/app/permission/storage"
	"github.com/wailovet/nuwa"
	"github.com/wailovet/nuwa-android/android"
)

//go:embed static
var static embed.FS

func main() {
	nuwa.Http().Static(static, "/static")

	android.Run("HelloWorld")
}
