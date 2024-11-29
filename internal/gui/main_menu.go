package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"github.com/EmptyInsid/db_gui/internal/database"
)

func MainWindow(w fyne.Window, db database.Service, role string) {
	MainMenu(w, db, role)
	emptyArea := container.NewStack()
	w.Resize(fyne.NewSize(1000, 500))
	w.CenterOnScreen()
	w.SetContent(container.NewBorder(nil, nil, nil, nil, emptyArea))
}

func MainMenu(w fyne.Window, db database.Service, role string) {

	reportFirst := fyne.NewMenuItem("Отчёт 1", func() {
		cont, err := MainReportFirst(w, db)
		if err != nil {
			dialog.ShowError(ErrReport, w)
		}
		w.SetContent(cont)
	})
	reportSecond := fyne.NewMenuItem("Отчёт 2", func() {
		cont, err := MainReportSecond(w, db)
		if err != nil {
			dialog.ShowError(ErrReport, w)
		}
		w.SetContent(cont)
	})
	reportThird := fyne.NewMenuItem("Отчёт 3", func() {
		cont, err := MainReportThird(w, db)
		if err != nil {
			dialog.ShowError(ErrReport, w)
		}
		w.SetContent(cont)
	})

	reportMenu := fyne.NewMenu("Отчёт", reportFirst, reportSecond, reportThird)

	jorney := fyne.NewMenuItem("Балансы", func() {
		jorneyContent, err := MainJorney(w, db, role)
		if err != nil {
			dialog.ShowError(ErrShowJorney, w)
		}
		w.SetContent(jorneyContent)
	})
	jorneyMenu := fyne.NewMenu("Журнал", jorney)

	dir := fyne.NewMenuItem("Справочник", func() {
		dirContent, err := MainDir(w, db, role)
		if err != nil {
			dialog.ShowError(ErrShowDir, w)
		}
		w.SetContent(dirContent)
	})
	dirMenu := fyne.NewMenu("Справочник", dir)

	w.SetMainMenu(fyne.NewMainMenu(jorneyMenu, dirMenu, reportMenu))
}
