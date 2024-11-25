package utils

import (
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

	gui.LoginMenu(w, db)

	// запуск приложения
	w.Show()
	a.Run()
}
