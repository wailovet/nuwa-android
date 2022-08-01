package main

import (
	"github.com/wailovet/nuwa"
	"github.com/wailovet/nuwa-android/android"
)

func main() {
	nuwa.Http().HandleFunc("/", func(ctx nuwa.HttpContext) {
		ctx.DisplayByString("Hello world")
	})

	android.Run("HelloWorld")
}
