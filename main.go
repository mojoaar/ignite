package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

const version = "0.1.0"

var menuApp *App

func main() {
	menuApp = NewApp()

	appMenu := menu.NewMenu()
	appMenu.Append(menu.AppMenu())

	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("New Project", keys.CmdOrCtrl("n"), func(_ *menu.CallbackData) {
		if menuApp != nil {
			runtime.EventsEmit(menuApp.ctx, "menu-new-project")
		}
	})
	fileMenu.AddText("Export Chat", keys.CmdOrCtrl("e"), func(_ *menu.CallbackData) {
		if menuApp != nil {
			runtime.EventsEmit(menuApp.ctx, "menu-export")
		}
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("Settings...", keys.CmdOrCtrl(","), func(_ *menu.CallbackData) {
		if menuApp != nil {
			runtime.EventsEmit(menuApp.ctx, "menu-settings")
		}
	})

	appMenu.Append(menu.EditMenu())

	err := wails.Run(&options.App{
		Title:  "Ignite",
		Width:  1024,
		Height: 768,

		MinWidth:  800,
		MinHeight: 600,

		AssetServer: &assetserver.Options{Assets: assets},
		OnStartup:   menuApp.startup,
		OnShutdown:  menuApp.shutdown,
		Menu:        appMenu,
		Bind:        []interface{}{menuApp},
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "Ignite",
				Message: "Provisioning with a heartbeat\nv" + version + "\n\n© 2026 Morten Johansen",
				Icon:    nil,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
