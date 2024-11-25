package gui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/EmptyInsid/db_gui/internal/database"
)

// table + toolbar
func GridViewer(db database.Service, table *widget.Table, toolbar *widget.Accordion) *container.Split {
	tableContainer := container.NewStack(table)
	toolBarContainer := container.NewVBox(toolbar)
	mainContent := container.NewHSplit(tableContainer, toolBarContainer)
	mainContent.SetOffset(0.7) // Устанавливает пропорцию (70% для таблицы, 30% для правой панели)

	return mainContent
}
