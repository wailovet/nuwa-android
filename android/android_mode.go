package android

import (
	"fmt"
	"image/color"
	"time"

	"gioui.org/app"
	_ "gioui.org/app/permission/networkstate"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
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
	for {
		_, _, errs := nuwa.HttpClient().Get(fmt.Sprint("http://127.0.0.1:", port)).End()
		if len(errs) == 0 {
			break
		}
		time.Sleep(time.Second)
	}
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

	th := material.NewTheme(gofont.Collection())
	var ops op.Ops

	var (
		config  = &gowebview.Config{URL: fmt.Sprint("http://127.0.0.1:", port), WindowConfig: &gowebview.WindowConfig{Title: title, VM: app.JavaVM()}}
		webview gowebview.WebView
	)

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case app.ViewEvent:
			//----------------------------------
			config.WindowConfig.Window = e.View // Here, sets the GioView. ;)
			//----------------------
			if webview == nil {
				go func() {
					var err error
					webview, err = gowebview.New(config)
					if err != nil {
						panic(err)
					}
					defer webview.Destroy()
					webview.Run()
				}()
			} else {
				webview.SetVisibility(gowebview.VisibilityMaximized)
			}
		case *system.CommandEvent:
			// You can minimize when click on "Back", without destroy.
			// It can be used when the webview is open in response to a button-click.
			if e.Type == system.CommandBack {
				e.Cancel = true
				webview.SetVisibility(gowebview.VisibilityMinimized)
			}
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			l := material.H1(th, "Hello, Gio")
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			l.Color = maroon
			l.Alignment = text.Middle
			l.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}

	// var wv gowebview.WebView
	// for {
	// 	e := <-w.Events()
	// 	switch e := e.(type) {
	// 	case app.ViewEvent:
	// 		if wv == nil {
	// 			go func() {
	// 				var err error
	// 				wv, err = gowebview.New(&gowebview.Config{URL: fmt.Sprint("http://127.0.0.1:", port), WindowConfig: &gowebview.WindowConfig{Title: title, Window: e.View, VM: app.JavaVM()}})
	// 				if err != nil {
	// 					panic(err)
	// 				}
	// 				defer wv.Destroy()
	// 				wv.Run()
	// 			}()
	// 		}
	// 	}
	// }
}
