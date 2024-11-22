package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func MainWindow(w fyne.Window, role string) {
	MainMenu(w)
	TabMenu(w)
}

func MainMenu(w fyne.Window) {

	reportFirst := fyne.NewMenuItem("Отчёт 1", func() {
		dialog.NewInformation("Создать отчёт 1", "По этой кнопке будет создаваться отчёт первого типа", w)
	})
	reportSecond := fyne.NewMenuItem("Отчёт 2", func() {
		dialog.NewInformation("Создать отчёт 2", "По этой кнопке будет создаваться отчёт второго типа", w)
	})
	reportMenu := fyne.NewMenu("Отчёт", reportFirst, reportSecond)
	w.SetMainMenu(fyne.NewMainMenu(reportMenu))
}

func TabMenu(w fyne.Window) *container.AppTabs {

	btnArticle := widget.NewButton("Статьи", func() {
		dialog.ShowInformation("Справочник статей", "По этой кнопке будет открываться справочник по статьим расходов", w)
	})
	btnOperation := widget.NewButton("Операции", func() {
		dialog.ShowInformation("Справочник операций", "По этой кнопке будет открываться справочник по операциям", w)
	})
	btnBalance := widget.NewButton("Балансы", func() {
		dialog.ShowInformation("Справочник балансов", "По этой кнопке будет открываться справочник по балансам", w)
	})

	dirButtons := container.NewVBox(btnArticle, btnOperation, btnBalance)

	jorney := container.NewTabItem("Журнал", widget.NewLabel("Главный раздел журнала"))
	directory := container.NewTabItem("Справочник", dirButtons)

	tab := container.NewAppTabs(jorney, directory)

	w.SetContent(tab)
	return tab
}
