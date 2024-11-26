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

func MainJorney(w fyne.Window, db database.Service, role string) (*container.Split, error) {
	table, err := BalanceTable(db)
	if err != nil {
		dialog.ShowError(err, w)
	}

	editor := AccordionJorney(w, db, table, role)

	return GridViewer(db, table, editor, role), nil
}

func AccordionJorney(w fyne.Window, db database.Service, table *widget.Table, role string) *widget.Accordion {
	//accFilters := FiltersButtons(w, db, table)
	accSums := SummariesAccord(w, db)
	accEdit := EditAccord(w, db, table)

	sumAccItem := widget.NewAccordionItem("Сводки", accSums)
	editAccItem := widget.NewAccordionItem("Редактировать", accEdit)

	editor := widget.NewAccordion(
		sumAccItem,
	)
	if role == "admin" {
		editor.Append(editAccItem)
	}

	return editor
}

func FiltersButtons(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	btnProfit := widget.NewButton("Баланс по дате", func() {
		dialog.ShowInformation("Баланс по дате формирования", "Здесь возможно будет баланс по дате", w)
	})

	return container.NewVBox(canvas.NewLine(color.White), btnProfit)
}

// РАЗДЕЛ СВОДОК
func SummariesAccord(w fyne.Window, db database.Service) *fyne.Container {
	winProfit := WinGetProfit(w, db)
	winCredit := WinGetCredit(w, db)
	winBalance := WinBalanceCount(w, db)

	return container.NewVBox(canvas.NewLine(color.White), winProfit, canvas.NewLine(color.White), winCredit, canvas.NewLine(color.White), winBalance)
}
func WinGetProfit(w fyne.Window, db database.Service) *fyne.Container {
	ctx := context.Background()

	startDate, endDate := MadeDateFields()
	article := widget.NewEntry()
	article.SetPlaceHolder("продукты")
	fieldProfit := widget.NewLabel("Доход: 0.00")

	btnArtOp := widget.NewButton("Доход за период", func() {
		if article.Text == "" {
			dialog.ShowError(fmt.Errorf("Укажите статью!"), w)
			return
		}

		profit, err := db.GetProfitByDate(ctx, article.Text, startDate.Text, endDate.Text)
		if err != nil {
			dialog.ShowError(err, w)
		}

		fieldProfit.SetText(fmt.Sprintf("Доход: %.2f", profit))

	})

	periodCont := container.NewStack(container.NewAdaptiveGrid(
		2,
		widget.NewLabel("Начало периода:"),
		startDate,
		widget.NewLabel("Конец периода:"),
		endDate,
		widget.NewLabel("Статья:"),
		article,
		fieldProfit,
	))

	return container.NewVBox(periodCont, btnArtOp)
}
func WinGetCredit(w fyne.Window, db database.Service) *fyne.Container {
	ctx := context.Background()

	startDate, endDate := MadeDateFields()
	article := widget.NewEntry()
	article.SetPlaceHolder("продукты")

	fieldCredit := widget.NewLabel("Расход: 0.00")

	btnArtOp := widget.NewButton("Расход за период", func() {

		if article.Text == "" {
			dialog.ShowError(fmt.Errorf("Укажите статью!"), w)
			return
		}

		credit, err := db.GetTotalCreditByArticleAndPeriod(ctx, article.Text, startDate.Text, endDate.Text)
		if err != nil {
			dialog.ShowError(err, w)
		}

		fieldCredit.SetText(fmt.Sprintf("Расход: %.2f", credit))

	})

	periodCont := container.NewStack(container.NewAdaptiveGrid(
		2,
		widget.NewLabel("Начало периода:"),
		startDate,
		widget.NewLabel("Конец периода:"),
		endDate,
		widget.NewLabel("Статья:"),
		article,
		fieldCredit,
	))

	return container.NewVBox(periodCont, btnArtOp)
}
func WinBalanceCount(w fyne.Window, db database.Service) *fyne.Container {
	ctx := context.Background()

	article := widget.NewEntry()
	article.SetPlaceHolder("продукты")

	fieldBalances := widget.NewLabel("Количество балансов: 0")

	btnArtOp := widget.NewButton("Количество балансов", func() {
		if article.Text == "" {
			dialog.ShowError(fmt.Errorf("Укажите статью!"), w)
			return
		}

		balanceCount, err := db.GetBalanceCountByArticleName(ctx, article.Text)
		if err != nil {
			dialog.ShowError(err, w)
		}

		fieldBalances.SetText(fmt.Sprintf("Количество балансов: %d", balanceCount))

	})

	commCont := container.NewAdaptiveGrid(
		2, widget.NewLabel("Статья:"), article, fieldBalances)

	return container.NewVBox(commCont, btnArtOp)
}

// РАЗДЕЛ РЕДАКТИРОВАНИЯ
func EditAccord(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	winCreateBalance := WinCreateNewBalance(w, db, table)
	winDelUnprof := WinDeleteUnprofitBalance(w, db, table)
	winDelBalance := WinDelBalance(w, db, table)

	return container.NewVBox(canvas.NewLine(color.White), winCreateBalance, canvas.NewLine(color.White), winDelBalance, canvas.NewLine(color.White), winDelUnprof)
}
func WinCreateNewBalance(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	ctx := context.Background()

	startDate, endDate := MadeDateFields()

	minProf := widget.NewEntry()
	minProf.SetPlaceHolder("100.0")
	periodCont := container.NewStack(container.NewAdaptiveGrid(
		2,
		widget.NewLabel("Начало периода:"),
		startDate,
		widget.NewLabel("Конец периода:"),
		endDate,
		widget.NewLabel("Минимальный профит:"),
		minProf,
	))

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

	periodCont := container.NewAdaptiveGrid(2, widget.NewLabel("Дата баланса:"), date)

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

	return container.NewVBox(periodCont, btn)
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
