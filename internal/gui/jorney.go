package gui

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"time"

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
		dialog.ShowError(ErrGetBalance, w)
	}

	editor := AccordionJorney(w, db, table, role)

	return GridViewer(db, table, editor, role), nil
}

func AccordionJorney(w fyne.Window, db database.Service, table *widget.Table, role string) *widget.Accordion {
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
	article := MadeSelectArticle(w, db)

	fieldProfit := widget.NewLabel("Доход: 0.00")

	btnArtOp := widget.NewButton("Доход за период", func() {
		if article.Selected == "" {
			dialog.ShowError(ErrEmptyArt, w)
			return
		}

		if err := CompareDate(startDate.Text, endDate.Text); err != nil {
			dialog.ShowError(ErrEndLessStart, w)
			return
		}

		profit, err := db.GetProfitByDate(ctx, article.Selected, startDate.Text, endDate.Text)
		if err != nil {
			dialog.ShowError(ErrGetProfit, w)
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
	article := MadeSelectArticle(w, db)

	fieldCredit := widget.NewLabel("Расход: 0.00")
	btnArtOp := widget.NewButton("Расход за период", func() {

		if article.Selected == "" {
			dialog.ShowError(ErrEmptyArt, w)
			return
		}

		if err := CompareDate(startDate.Text, endDate.Text); err != nil {
			dialog.ShowError(ErrEndLessStart, w)
			return
		}

		credit, err := db.GetTotalCreditByArticleAndPeriod(ctx, article.Selected, startDate.Text, endDate.Text)
		if err != nil {
			dialog.ShowError(ErrGetTotalCredit, w)
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

	article := MadeSelectArticle(w, db)

	fieldBalances := widget.NewLabel("Количество балансов: 0")

	btnArtOp := widget.NewButton("Количество балансов", func() {
		if article.Selected == "" {
			dialog.ShowError(ErrEmptyArt, w)
			return
		}

		balanceCount, err := db.GetBalanceCountByArticleName(ctx, article.Selected)
		if err != nil {
			dialog.ShowError(ErrGetBalanceCount, w)
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

	_, endDate := MadeDateFields()

	minProf := widget.NewEntry()
	minProf.SetPlaceHolder("100.0")
	periodCont := container.NewStack(container.NewAdaptiveGrid(
		1,
		widget.NewLabel("Конец периода:"),
		endDate,
		widget.NewLabel("Минимальный профит:"),
		minProf,
	))

	btnArtOp := widget.NewButton("Создать баланс", func() {

		date, err := time.Parse("2006-01-02", endDate.Text)
		if err != nil {
			dialog.ShowError(ErrParseDate, w)
			return
		}

		// Проверяем, является ли дата концом месяца
		if !isEndOfMonth(date) {
			dialog.ShowError(ErrCreateBalanceDate, w)
			return
		}

		floatMinProf, err := strconv.ParseFloat(minProf.Text, 64)
		if err != nil {
			dialog.ShowError(ErrParseDebit, w)
			return
		}
		// Получаем начало месяца
		startOfMonth := getStartOfMonth(date)

		err = db.CreateBalanceIfProfitable(ctx, startOfMonth.Format("2006-01-02"), endDate.Text, floatMinProf)
		if err != nil {
			if errors.Is(err, database.ErrLessThenMin) {
				dialog.ShowError(ErrMinBalanceProfit, w)
				return
			}
			dialog.ShowError(ErrCreateBalance, w)
			return
		} else {
			dialog.ShowInformation("Создать баланс", "Новый баланс создан успешно!", w)
		}
		if err := UpdateBalanceTable(db, table); err != nil {
			dialog.ShowError(ErrUpdBalance, w)
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
					dialog.ShowError(ErrDelBalance, w)
					return
				}
				if err := UpdateBalanceTable(db, table); err != nil {
					dialog.ShowError(ErrUpdBalance, w)
					return
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
					dialog.ShowError(ErrDelMinBalance, w)
				}
				if err := UpdateBalanceTable(db, table); err != nil {
					dialog.ShowError(ErrUpdBalance, w)
				}
			},
			w)

	})

	return container.NewVBox(btn)
}
