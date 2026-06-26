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

var menuApp *App

func main() {
	menuApp = NewApp()

	appMenu := menu.NewMenu()

	igniteMenu := appMenu.AddSubmenu("Ignite")
	igniteMenu.AddText("About Ignite", nil, func(_ *menu.CallbackData) {
		if menuApp != nil {
			runtime.EventsEmit(menuApp.ctx, "menu-about")
		}
	})
	igniteMenu.AddSeparator()
	igniteMenu.AddText("Settings...", keys.CmdOrCtrl(","), func(_ *menu.CallbackData) {
		if menuApp != nil {
			runtime.EventsEmit(menuApp.ctx, "menu-settings")
		}
	})
	igniteMenu.AddSeparator()
	igniteMenu.Append(menu.AppMenu())
	appMenu.Append(menu.EditMenu())

	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("New Project", keys.CmdOrCtrl("n"), func(_ *menu.CallbackData) {
		if menuApp != nil {
			runtime.EventsEmit(menuApp.ctx, "menu-new-project")
		}
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("Export Chat", keys.CmdOrCtrl("e"), func(_ *menu.CallbackData) {
		if menuApp != nil {
			runtime.EventsEmit(menuApp.ctx, "menu-export")
		}
	})

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
				Message: "Provisioning with a heartbeat\n\n© 2026 Morten Johansen",
				Icon:    nil,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
