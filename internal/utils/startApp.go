package utils

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/EmptyInsid/db_gui/internal/database"
	"github.com/EmptyInsid/db_gui/internal/gui"
)

func StartApp(db database.Service) {

	// инициализация приложения
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())

	w := a.NewWindow("Домашний бюджет")
	w.CenterOnScreen()

	ic, err := fyne.LoadResourceFromPath("..\\..\\src\\icon.png")
	if err != nil {
		log.Print(err)
	}
	w.SetIcon(ic)

	gui.LoginMenu(a, w, db)

	// запуск приложения
	w.Show()
	a.Run()
}
