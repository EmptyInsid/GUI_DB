package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/EmptyInsid/db_gui/internal/database"
)

func MainDir(w fyne.Window, db database.Service) (*fyne.Container, error) {
	dirContent, err := TabsDir(w, db)
	if err != nil {
		return nil, err
	}
	return container.NewStack(dirContent), nil
}

func ArticleViewer(w fyne.Window, db database.Service) (*fyne.Container, error) {
	editor, err := AccordionDir(w)
	if err != nil {
		return nil, err
	}

	table, err := ArticleTable(db)
	if err != nil {
		return nil, err
	}

	return GridViewer(db, table, editor), nil
}

func OperationsViewer(w fyne.Window, db database.Service) (*fyne.Container, error) {
	editor, err := AccordionDir(w)
	if err != nil {
		return nil, err
	}

	table, err := OperationsTable(db)
	if err != nil {
		return nil, err
	}

	return GridViewer(db, table, editor), nil
}

func TabsDir(w fyne.Window, db database.Service) (*container.AppTabs, error) {

	articleContent, err := ArticleViewer(w, db)
	if err != nil {
		return nil, err
	}

	operContent, err := OperationsViewer(w, db)
	if err != nil {
		return nil, err
	}

	article := container.NewTabItem("Статьи", articleContent)
	operations := container.NewTabItem("Операции", operContent)

	tab := container.NewAppTabs(article, operations)
	tab.SetTabLocation(container.TabLocationTop)
	return tab, nil
}

func AccordionDir(w fyne.Window) (*widget.Accordion, error) {
	btnsAdd := widget.NewButton("Добавить", func() {
		dialog.ShowInformation("Редактировать журнал", "Здесь можно будет редактирвоать журнал", w)
	})
	btnsDel := widget.NewButton("Удалить", func() {
		dialog.ShowInformation("Редактировать журнал", "Здесь можно будет редактирвоать журнал", w)
	})
	btnsEdit := widget.NewButton("Редактировать", func() {
		dialog.ShowInformation("Редактировать журнал", "Здесь можно будет редактирвоать журнал", w)
	})

	editor := widget.NewAccordion(
		widget.NewAccordionItem("Добавить", btnsAdd),
		widget.NewAccordionItem("Удалить", btnsDel),
		widget.NewAccordionItem("Редактировать", btnsEdit),
	)
	return editor, nil
}
