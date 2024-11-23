package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/EmptyInsid/db_gui/internal/database"
)

// table + toolbar
func GridViewer(db database.Service, table *widget.Table, toolbar *widget.Accordion) *fyne.Container {
	tableContainer := container.NewStack(table)
	toolBarContainer := container.NewVBox(toolbar)

	return container.NewBorder(nil, nil, nil, toolBarContainer, tableContainer)
}
