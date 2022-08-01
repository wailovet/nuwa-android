package android

import (
	"fmt"
	"time"

	"gioui.org/app"
	"github.com/inkeliz/gowebview"
	"github.com/wailovet/gofunc"
	"github.com/wailovet/nuwa"
)

func Run(title string, hes ...*nuwa.HttpEngine) {
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

	go func() {
		w := app.NewWindow()
		err := loop(w, title, port)
		if err != nil {
			panic(err)
		}
	}()
	app.Main()
}

func loop(w *app.Window, title string, port int) error {

	var wv gowebview.WebView
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case app.ViewEvent:
			if wv == nil {
				go func() {
					var err error
					wv, err = gowebview.New(&gowebview.Config{URL: fmt.Sprint("http://127.0.0.1:", port), WindowConfig: &gowebview.WindowConfig{Title: title, Window: e.View, VM: app.JavaVM()}})
					if err != nil {
						panic(err)
					}
					defer wv.Destroy()
					wv.Run()
				}()
			}
		}
	}
}
