package gui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/EmptyInsid/db_gui/internal/database"
)

func OperationsTable(db database.Service) (*widget.Table, error) {
	ctx := context.Background()

	data, err := db.GetArticlesWithOperations(ctx)
	if err != nil {
		return nil, err
	}

	header := []string{"Номер", "Id", "Статья", "Доход", "Расход", "Дата"}

	table := widget.NewTable(
		func() (int, int) {
			return len(data) + 1, len(header)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("very very wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			lable := o.(*widget.Label)
			col, row := i.Col, i.Row

			if row == 0 {
				//lable.Alignment = fyne.TextAlignCenter
				//lable.TextStyle = fyne.TextStyle{Bold: true}
				lable.SetText(header[col])
			} else {
				switch col {
				case 0:
					lable.SetText(fmt.Sprint(row))
				case 1:
					lable.SetText(fmt.Sprint(data[row-1].OperationID))
				case 2:
					lable.SetText(fmt.Sprint(data[row-1].ArticleName))
				case 3:
					lable.SetText(fmt.Sprint(data[row-1].Debit))
				case 4:
					lable.SetText(fmt.Sprint(data[row-1].Credit))
				case 5:
					lable.SetText(fmt.Sprint(data[row-1].CreateDate.Format("2006-01-02")))
				default:
					lable.SetText("-")
				}

			}
		})

	table.SetColumnWidth(0, widget.NewLabel("Number").MinSize().Width)
	table.SetColumnWidth(1, widget.NewLabel("Number").MinSize().Width)
	table.SetColumnWidth(2, widget.NewLabel("very very wide content").MinSize().Width)
	table.SetColumnWidth(3, widget.NewLabel("10000000").MinSize().Width)
	table.SetColumnWidth(4, widget.NewLabel("10000000").MinSize().Width)
	table.SetColumnWidth(5, widget.NewLabel("2024-11-01").MinSize().Width)

	return table, nil
}

func BalanceTable(db database.Service) (*widget.Table, error) {
	ctx := context.Background()

	data, err := db.GetAllBalances(ctx)
	if err != nil {
		return nil, err
	}

	header := []string{"Номер", "Дата", "Доход", "Расход", "Итог"}

	table := widget.NewTable(
		func() (int, int) {
			return len(data) + 1, len(header)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("very very wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			lable := o.(*widget.Label)
			col, row := i.Col, i.Row

			if row == 0 {
				lable.SetText(header[col])
			} else {
				switch col {
				case 0:
					lable.SetText(fmt.Sprint(row))
				case 1:
					lable.SetText(fmt.Sprint(data[row-1].Date.Format("2006-01-02")))
				case 2:
					lable.SetText(fmt.Sprint(data[row-1].Debit))
				case 3:
					lable.SetText(fmt.Sprint(data[row-1].Credit))
				case 4:
					lable.SetText(fmt.Sprint(data[row-1].Amount))
				default:
					lable.SetText("-")
				}

			}
		})

	table.SetColumnWidth(0, widget.NewLabel("Number").MinSize().Width)
	table.SetColumnWidth(1, widget.NewLabel("2024-11-01").MinSize().Width)
	table.SetColumnWidth(2, widget.NewLabel("10000000").MinSize().Width)
	table.SetColumnWidth(3, widget.NewLabel("10000000").MinSize().Width)
	table.SetColumnWidth(4, widget.NewLabel("10000000").MinSize().Width)

	return table, nil
}

func ArticleTable(db database.Service) (*widget.Table, error) {
	ctx := context.Background()

	data, err := db.GetAllArticles(ctx)
	if err != nil {
		return nil, err
	}

	header := []string{"Номер", "Статья"}

	table := widget.NewTable(
		func() (int, int) {
			return len(data) + 1, len(header)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("very very wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			lable := o.(*widget.Label)
			col, row := i.Col, i.Row

			if row == 0 {
				lable.SetText(header[col])
			} else {
				switch col {
				case 0:
					lable.SetText(fmt.Sprint(row))
				case 1:
					lable.SetText(fmt.Sprint(data[row-1].Name))
				default:
					lable.SetText("-")
				}

			}
		})

	table.SetColumnWidth(0, widget.NewLabel("Number").MinSize().Width)
	table.SetColumnWidth(1, widget.NewLabel("very very wide content").MinSize().Width)

	return table, nil
}

func UnaccountedOpertionsMoneyTable(db database.Service) (*widget.Table, error) {
	ctx := context.Background()

	data, err := db.GetViewUnaccountedOpertions(ctx)
	if err != nil {
		return nil, err
	}

	header := []string{"Номер", "Статья", "Общий доход", "Общий расход"}

	table := widget.NewTable(
		func() (int, int) {
			return len(data) + 1, len(header)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("very very wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			lable := o.(*widget.Label)
			col, row := i.Col, i.Row

			if row == 0 {
				lable.SetText(header[col])
			} else {
				switch col {
				case 0:
					lable.SetText(fmt.Sprint(row))
				case 1:
					lable.SetText(fmt.Sprint(data[row-1].ArticleName))
				case 2:
					lable.SetText(fmt.Sprint(data[row-1].TotalDebit))
				case 3:
					lable.SetText(fmt.Sprint(data[row-1].TotalCredit))
				default:
					lable.SetText("-")
				}

			}
		})

	table.SetColumnWidth(0, widget.NewLabel("Number").MinSize().Width)
	table.SetColumnWidth(1, widget.NewLabel("very very wide content").MinSize().Width)

	return table, nil
}

func UnusedArticlesTable(db database.Service, startData string, finishData string) (*widget.Table, error) {
	ctx := context.Background()

	data, err := db.GetUnusedArticles(ctx, startData, finishData)
	if err != nil {
		return nil, err
	}

	header := []string{"Номер", "Статья"}

	table := widget.NewTable(
		func() (int, int) {
			return len(data) + 1, len(header)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("very very wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			lable := o.(*widget.Label)
			col, row := i.Col, i.Row

			if row == 0 {
				lable.SetText(header[col])
			} else {
				switch col {
				case 0:
					lable.SetText(fmt.Sprint(row))
				case 1:
					lable.SetText(fmt.Sprint(data[row-1].Name))
				default:
					lable.SetText("-")
				}

			}
		})

	table.SetColumnWidth(0, widget.NewLabel("Number").MinSize().Width)
	table.SetColumnWidth(1, widget.NewLabel("very very wide content").MinSize().Width)

	return table, nil
}

func CountBalanceOperTable(db database.Service, startData string, finishData string) (*widget.Table, error) {
	ctx := context.Background()

	data, err := db.GetViewCountBalanceOper(ctx)
	if err != nil {
		return nil, err
	}

	header := []string{"Номер", "Дата", "Кол-во операций"}

	table := widget.NewTable(
		func() (int, int) {
			return len(data) + 1, len(header)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("very very wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			lable := o.(*widget.Label)
			col, row := i.Col, i.Row

			if row == 0 {
				lable.SetText(header[col])
			} else {
				switch col {
				case 0:
					lable.SetText(fmt.Sprint(row))
				case 1:
					lable.SetText(fmt.Sprint(data[row-1].BalanceDate))
				case 2:
					lable.SetText(fmt.Sprint(data[row-1].OperationCount))
				default:
					lable.SetText("-")
				}

			}
		})

	table.SetColumnWidth(0, widget.NewLabel("Number").MinSize().Width)
	table.SetColumnWidth(1, widget.NewLabel("2024-11-01").MinSize().Width)
	table.SetColumnWidth(2, widget.NewLabel("100").MinSize().Width)

	return table, nil
}
