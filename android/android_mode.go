package android

import (
	"fmt"
	"log"
	"os"
	"time"

	"gioui.org/app"
	_ "gioui.org/app/permission/networkstate"
	"gioui.org/io/system"
	"github.com/inkeliz/gowebview"
	"github.com/wailovet/gofunc"
	"github.com/wailovet/nuwa"
)

var logger *log.Logger

func Logger() *log.Logger {
	if logger != nil {
		return logger
	}
	logFile, err := os.OpenFile("/sdcard/gom.log", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return nil
	}
	logger = log.New(logFile, "<PS>", log.Lshortfile|log.Ldate|log.Ltime)
	return logger
}

func Run(title string, hes ...*nuwa.HttpEngine) {
	Logger().Println("Run:", title)
	port := nuwa.Helper().GetFreePort()

	var he *nuwa.HttpEngine
	if len(hes) > 0 {
		he = hes[0]
	} else {
		he = nuwa.Http()
	}
	he.InstanceConfig.Port = fmt.Sprint(port)
	gofunc.New(func() {
		for {
			he.Run()
			time.Sleep(time.Second * 5)
		}
	})

	for {
		_, _, errs := nuwa.HttpClient().Get(fmt.Sprint("http://127.0.0.1:", port)).End()
		if len(errs) == 0 {
			break
		}
		time.Sleep(time.Second)
	}

	// gofunc.New(func() {
	// 	startWebview(title, port)
	// })
	gofunc.New(func() {
		w := app.NewWindow()
		err := loop(w, title, port)
		if err != nil {
			panic(err)
		}
	})
	app.Main()
}

var config *gowebview.Config
var webview gowebview.WebView

func startWebview(view uintptr, title string, port int) {
	var config = &gowebview.Config{URL: fmt.Sprint("http://127.0.0.1:", port), WindowConfig: &gowebview.WindowConfig{Title: title, VM: app.JavaVM()}}

	config.WindowConfig.Window = view // Here, sets the GioView. ;)
	gofunc.New(func() {
		if webview != nil {
			webview.Destroy()
		}
		webview, _ = gowebview.New(config)
		Logger().Println("webview:start")
		webview.Run()
		Logger().Println("webview:end")
	})
}

func loop(w *app.Window, title string, port int) error {

	for {
		e := <-w.Events()
		Logger().Println("Events:", e, nuwa.Helper().JsonEncode(e))

		switch e := e.(type) {
		case app.ViewEvent: // ViewEvent is the main event type.
			if e.View > 0 {
				startWebview(e.View, title, port)
			}

		case system.DestroyEvent:
			return e.Err
		}
	}

}
