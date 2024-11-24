package gui

import (
	"context"
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/EmptyInsid/db_gui/internal/database"
)

func MainJorney(w fyne.Window, db database.Service) (*fyne.Container, error) {
	table, err := BalanceTable(db)
	if err != nil {
		dialog.ShowError(err, w)
	}

	editor := AccordionJorney(w, db, table)

	return GridViewer(db, table, editor), nil
}

func AccordionJorney(w fyne.Window, db database.Service, table *widget.Table) *widget.Accordion {
	accFilters := FiltersButtons(w, db, table)
	accSums := SummariesAccord(w, db)
	accEdit := EditAccord(w, db, table)

	editor := widget.NewAccordion(
		widget.NewAccordionItem("Фильтры", accFilters),
		widget.NewAccordionItem("Сводки", accSums),
		widget.NewAccordionItem("Редактировать", accEdit),
	)
	return editor
}

func FiltersButtons(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	btnProfit := widget.NewButton("Баланс по дате", func() {
		dialog.ShowInformation("Баланс по дате формирования", "Здесь возможно будет баланс по дате", w)
	})

	return container.NewVBox(btnProfit)
}

// РАЗДЕЛ СВОДОК
func SummariesAccord(w fyne.Window, db database.Service) *fyne.Container {
	winProfit := WinGetProfit(w, db)
	winCredit := WinGetCredit(w, db)
	winBalance := WinBalanceCount(w, db)

	return container.NewVBox(winProfit, canvas.NewLine(color.Black), winCredit, canvas.NewLine(color.Black), winBalance)
}
func WinGetProfit(w fyne.Window, db database.Service) *fyne.Container {
	ctx := context.Background()

	startDate := widget.NewEntry()
	endDate := widget.NewEntry()
	startDate.SetPlaceHolder("2024-07-01")
	endDate.SetPlaceHolder("2024-07-21")
	periodCont := container.NewStack(container.NewAdaptiveGrid(2, startDate, endDate))

	fieldProfit := widget.NewLabel("0.00")

	btnArtOp := widget.NewButton("Доход за период", func() {

		profit, err := db.GetProfitByDate(ctx, startDate.Text, endDate.Text)
		if err != nil {
			dialog.ShowError(err, w)
		}

		fieldProfit.SetText(fmt.Sprintf("%2f", profit))

	})

	return container.NewVBox(periodCont, fieldProfit, btnArtOp)
}
func WinGetCredit(w fyne.Window, db database.Service) *fyne.Container {
	ctx := context.Background()

	startDate := widget.NewEntry()
	endDate := widget.NewEntry()
	article := widget.NewEntry()
	startDate.SetPlaceHolder("2024-07-01")
	endDate.SetPlaceHolder("2024-07-21")
	article.SetPlaceHolder("статья")
	periodCont := container.NewStack(container.NewAdaptiveGrid(2, startDate, endDate, article))

	fieldCredit := widget.NewLabel("0.00")

	btnArtOp := widget.NewButton("Расход за период", func() {

		credit, err := db.GetTotalCreditByArticleAndPeriod(ctx, article.Text, startDate.Text, endDate.Text)
		if err != nil {
			dialog.ShowError(err, w)
		}

		fieldCredit.SetText(fmt.Sprintf("%2f", credit))

	})

	return container.NewVBox(periodCont, fieldCredit, btnArtOp)
}
func WinBalanceCount(w fyne.Window, db database.Service) *fyne.Container {
	ctx := context.Background()

	article := widget.NewEntry()
	article.SetPlaceHolder("статья")

	fieldBalances := widget.NewLabel("0")

	btnArtOp := widget.NewButton("Количество балансов", func() {

		balanceCount, err := db.GetBalanceCountByArticleName(ctx, article.Text)
		if err != nil {
			dialog.ShowError(err, w)
		}

		fieldBalances.SetText(fmt.Sprint(balanceCount))

	})

	return container.NewVBox(article, fieldBalances, btnArtOp)
}

// РАЗДЕЛ РЕДАКТИРОВАНИЯ
func EditAccord(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	winCreateBalance := WinCreateNewBalance(w, db, table)
	winDelUnprof := WinDeleteUnprofitBalance(w, db, table)
	winDelBalance := WinDelBalance(w, db, table)

	return container.NewVBox(winCreateBalance, canvas.NewLine(color.White), winDelBalance, canvas.NewLine(color.White), winDelUnprof)
}
func WinCreateNewBalance(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	ctx := context.Background()

	startDate := widget.NewEntry()
	endDate := widget.NewEntry()
	minProf := widget.NewEntry()
	startDate.SetPlaceHolder("2024-07-01")
	endDate.SetPlaceHolder("2024-07-21")
	minProf.SetPlaceHolder("100.0")
	periodCont := container.NewStack(container.NewAdaptiveGrid(2, startDate, endDate, minProf))

	btnArtOp := widget.NewButton("Создать баланс", func() {

		floatMinProf, err := strconv.ParseFloat(minProf.Text, 64)
		if err != nil {
			dialog.ShowError(err, w)
		}

		err = db.CreateBalanceIfProfitable(ctx, startDate.Text, endDate.Text, floatMinProf)
		if err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("Создать баланс", "Новый баланс создан успешно!", w)
		}
		if err := UpdateBalanceTable(db, table); err != nil {
			dialog.ShowError(err, w)
		}
	})

	return container.NewVBox(periodCont, btnArtOp)
}
func WinDelBalance(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	ctx := context.Background()

	date := widget.NewEntry()
	date.SetPlaceHolder("дата баланса")

	btn := widget.NewButton("Удалить баланс", func() {

		dialog.ShowConfirm(
			"Удалить баланс",
			"Вы уверены, что хотите удалить баланс?",
			func(bool) {
				err := db.DeleteBalance(ctx, date.Text)
				if err != nil {
					dialog.ShowError(err, w)
				}
				if err := UpdateBalanceTable(db, table); err != nil {
					dialog.ShowError(err, w)
				}
			},
			w)

	})

	return container.NewVBox(date, btn)
}
func WinDeleteUnprofitBalance(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	ctx := context.Background()

	btn := widget.NewButton("Удалить неприбыльный баланс", func() {

		dialog.ShowConfirm(
			"Удалить неприбыльный баланс",
			"Вы уверены, что хотите удалить самый неприбыльный баланс?",
			func(bool) {
				err := db.DeleteMostUnprofitableBalance(ctx)
				if err != nil {
					dialog.ShowError(err, w)
				}
				if err := UpdateBalanceTable(db, table); err != nil {
					dialog.ShowError(err, w)
				}
			},
			w)

	})

	return container.NewVBox(btn)
}
