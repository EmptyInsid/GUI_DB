package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/EmptyInsid/db_gui/internal/database"
)

func MainJorney(w fyne.Window, db database.Service) (*fyne.Container, error) {
	editor, err := AccordionJorney(w)
	if err != nil {
		dialog.ShowError(err, w)
	}

	table, err := BalanceTable(db)
	if err != nil {
		dialog.ShowError(err, w)
	}

	return GridViewer(db, table, editor), nil
}

func AccordionJorney(w fyne.Window) (*widget.Accordion, error) {
	btnsFilters, err := FiltersButtons(w)
	if err != nil {
		return nil, err
	}
	btnsSums, err := SummariesButtons(w)
	if err != nil {
		return nil, err
	}
	btnsEdit, err := EditButtons(w)
	if err != nil {
		return nil, err
	}

	editor := widget.NewAccordion(
		widget.NewAccordionItem("Фильтры", btnsFilters),
		widget.NewAccordionItem("Сводки", btnsSums),
		widget.NewAccordionItem("Редактировать", btnsEdit),
	)
	return editor, nil
}

func FiltersButtons(w fyne.Window) (*fyne.Container, error) {
	btnProfit := widget.NewButton("Доход за период", func() {
		dialog.ShowInformation("Редактировать журнал", "Здесь можно будет редактирвоать журнал", w)
	})

	btnCredit := widget.NewButton("Расход за период", func() {
		dialog.ShowInformation("Редактировать журнал", "Здесь можно будет редактирвоать журнал", w)
	})

	btnBalances := widget.NewButton("Количество балансов", func() {
		dialog.ShowInformation("Редактировать журнал", "Здесь можно будет редактирвоать журнал", w)
	})

	return container.NewVBox(btnProfit, btnCredit, btnBalances), nil
}

func SummariesButtons(w fyne.Window) (*fyne.Container, error) {
	btnArtOp := widget.NewButton("Статьи с операциями", func() {
		dialog.ShowInformation("Редактировать журнал", "Здесь можно будет редактирвоать журнал", w)
	})

	btnUnOp := widget.NewButton("Неучтённые операции", func() {
		dialog.ShowInformation("Редактировать журнал", "Здесь можно будет редактирвоать журнал", w)
	})

	btnOpArt := widget.NewButton("Операции по статье", func() {
		dialog.ShowInformation("Редактировать журнал", "Здесь можно будет редактирвоать журнал", w)
	})

	return container.NewVBox(btnArtOp, btnUnOp, btnOpArt), nil
}

func EditButtons(w fyne.Window) (*fyne.Container, error) {
	btnAdd := widget.NewButton("добавить операцию", func() {
		dialog.ShowInformation("Редактировать журнал", "Здесь можно будет редактирвоать журнал", w)
	})

	btnDel := widget.NewButton("Удалить операцию", func() {
		dialog.ShowInformation("Редактировать журнал", "Здесь можно будет редактирвоать журнал", w)
	})

	btnEdit := widget.NewButton("Редактировать операцию", func() {
		dialog.ShowInformation("Редактировать журнал", "Здесь можно будет редактирвоать журнал", w)
	})

	return container.NewVBox(btnAdd, btnDel, btnEdit), nil
}
