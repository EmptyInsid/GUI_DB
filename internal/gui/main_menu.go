package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/EmptyInsid/db_gui/internal/database"
)

func MainWindow(myApp fyne.App, w fyne.Window, db database.Service, role string) {
	MainMenu(myApp, w, db, role)
	emptyArea := container.NewStack()
	w.Resize(fyne.NewSize(1000, 500))
	w.CenterOnScreen()
	w.SetContent(container.NewBorder(nil, nil, nil, nil, emptyArea))
}

func MainMenu(myApp fyne.App, w fyne.Window, db database.Service, role string) {

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

	infoContent := fyne.NewMenuItem("Информация", func() {
		aboutWindow := createAboutWindow(myApp)
		aboutWindow.Show()
	})
	infoMenu := fyne.NewMenu("О приложении", infoContent)

	exit := fyne.NewMenuItem("Выход", func() {
		dialog.ShowConfirm("Выход", "Вы уверены, что хотите выйти из приложения?",
			func(bool) {
				w.SetMainMenu(nil)
				LoginMenu(myApp, w, db)
			},
			w,
		)
	})
	exitMenu := fyne.NewMenu("Выход", exit)

	w.SetMainMenu(fyne.NewMainMenu(jorneyMenu, dirMenu, reportMenu, infoMenu, exitMenu))
}

func createAboutWindow(app fyne.App) fyne.Window {
	aboutWindow := app.NewWindow("О приложении")

	instruction := `
	Это приложение предназначено для ведение домашнего бюджета.
	Инструкция по использованию:
	
	1. Общий концепт - приложение предоставляет интерфейс для просмотра и редактирования операций 
	с информацией о статье, доходе, расходе и дате. Статьи выбираются из списка ранее добавленных.
	В конце месяца доступна функция формирования баланса за месяц с информацией о расходах, доходах
	и прибыли за месяц.
	2. Предоставлено три вкладки - Журнал, Справочник, Отчёты.
	3. В разделе Журнал предоставлен следующий интерфейс:
	  3.1. Просмотр сформированных балансов.
	  3.2. Просмотр сводных данных о доходах и расходах.
	  3.3. Формирование и расформирование балансов [admin]
	4. В разделе Справочник предоставлен следующий интерфейс:
	  4.1. Вкладка статей с возможностью добавить, редактировать, удалить статью
	  4.2. Вкладка операций с возможностью добавить, редактировать, удалить операцию
	5. В разделе отчёты предоставлен следующий интерфейс:
	  5.1. Выбор типа отчёта из возможных
	  5.2. Введение данных для формирования по ним отчёта
	  5.3. Сохранение документа сформированного отчёта
	
	Обратите внимание: 
	- Все данные сохраняются автоматически.
	- Для получения возможностей редактирования нелбходимо иметь роль администратора
	- Сформированный баланс характеризует закрытый период, то есть нельзя редактировать
	операции, который входят в закрытый период.

	Для вопросов и поддержки обратитесь к разработчику.`

	label := `
	Приложение: Автоматизация домашнего бюджета
	Версия: 1.0.0
	Автор: жив
	`

	textLabel := widget.NewLabel(instruction)
	textLabel.Wrapping = fyne.TextWrapWord // Перенос текста по словам
	// Оборачиваем текст в скролл
	scroll := container.NewScroll(textLabel)

	aboutContent := container.NewVSplit(
		widget.NewLabel(label),
		scroll,
	)
	aboutContent.SetOffset(0.2)

	aboutWindow.Resize(fyne.NewSize(400, 300))
	aboutWindow.SetContent(aboutContent)
	return aboutWindow
}
