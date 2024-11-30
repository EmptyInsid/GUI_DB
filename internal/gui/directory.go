package gui

import (
	"context"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/EmptyInsid/db_gui/internal/database"
)

func MainDir(w fyne.Window, db database.Service, role string) (*fyne.Container, error) {
	dirContent, err := TabsDir(w, db, role)
	if err != nil {
		return nil, err
	}
	return container.NewStack(dirContent), nil
}

func TabsDir(w fyne.Window, db database.Service, role string) (*container.AppTabs, error) {

	articleContent, err := ArticleViewer(w, db, role)
	if err != nil {
		return nil, err
	}

	operContent, err := OperationsViewer(w, db, role)
	if err != nil {
		return nil, err
	}

	article := container.NewTabItem("Статьи", articleContent)
	operations := container.NewTabItem("Операции", operContent)

	tab := container.NewAppTabs(article, operations)
	tab.SetTabLocation(container.TabLocationTop)
	return tab, nil
}

func ArticleViewer(w fyne.Window, db database.Service, role string) (*container.Split, error) {
	table, err := ArticleTable(db)
	if err != nil {
		return nil, err
	}
	editor, err := AccordionDirArticle(w, db, table, role)
	if err != nil {
		return nil, err
	}

	if role != "admin" {
		editor.Hide()
	}

	return GridViewer(db, table, editor, role), nil
}

func OperationsViewer(w fyne.Window, db database.Service, role string) (*container.Split, error) {
	table, err := OperationsTable(db)
	if err != nil {
		return nil, err
	}
	editor, err := AccordionDirOper(w, db, table)
	if err != nil {
		return nil, err
	}

	if role != "admin" {
		editor.Hide()
	}

	return GridViewer(db, table, editor, role), nil
}

// СПИСОК ДЕЙСТВИЙ ДЛЯ СТАТЕЙ
func AccordionDirArticle(w fyne.Window, db database.Service, table *widget.Table, role string) (*widget.Accordion, error) {
	accAdd := AddArticle(w, db, table, role)
	accEdit := EditArticle(w, db, table, role)
	accDel := DelArticle(w, db, table, role)

	editor := widget.NewAccordion(
		widget.NewAccordionItem("Добавить", accAdd),
		widget.NewAccordionItem("Редактировать", accEdit),
		widget.NewAccordionItem("Удалить", accDel),
	)
	return editor, nil
}

// РАЗДЕЛ ДОБАВИТЬ СТАТЬЮ
func AddArticle(w fyne.Window, db database.Service, table *widget.Table, role string) *fyne.Container {
	winAddArticle := WinAddArticle(w, db, table, role)
	return container.NewVBox(canvas.NewLine(color.White), winAddArticle)
}
func WinAddArticle(w fyne.Window, db database.Service, table *widget.Table, role string) *fyne.Container {
	ctx := context.Background()

	article := widget.NewEntry()
	article.SetPlaceHolder("новое название")

	cont := container.NewAdaptiveGrid(2, widget.NewLabel("Статья"), article)

	btn := widget.NewButton("Добавить статью", func() {

		if article.Text == "" {
			dialog.ShowError(ErrEmptyArt, w)
			return
		}

		err := db.AddArticle(ctx, article.Text)
		if err != nil {
			dialog.ShowError(ErrAddArt, w)
			return
		} else {
			dialog.ShowInformation("Добавить статью", "Новая статья успешно добавлена!", w)
		}

		err = UpdateArticleTable(db, table)
		if err != nil {
			dialog.ShowError(ErrUpdArt, w)
			return
		}

		dirContent, err := MainDir(w, db, role)
		if err != nil {
			dialog.ShowError(ErrShowDir, w)
		}
		w.SetContent(dirContent)

	})

	return container.NewVBox(cont, btn)
}

// РАЗДЕЛ РЕДАКТИРОВАТЬ СТАТЬЮ
func EditArticle(w fyne.Window, db database.Service, table *widget.Table, role string) *fyne.Container {
	winEditArticle := WinEditArticle(w, db, table, role)
	return container.NewVBox(canvas.NewLine(color.White), winEditArticle)
}
func WinEditArticle(w fyne.Window, db database.Service, table *widget.Table, role string) *fyne.Container {
	ctx := context.Background()

	oldName := MadeSelectArticle(w, db)
	newName := widget.NewEntry()
	newName.SetPlaceHolder("еда")
	fieldsCont := container.NewStack(container.NewAdaptiveGrid(
		2,
		widget.NewLabel("Старое имя"),
		oldName,
		widget.NewLabel("Новое имя"),
		newName,
	))

	btn := widget.NewButton("Изменить имя", func() {

		if oldName.Selected == "" {
			dialog.ShowError(ErrEmptyArt, w)
			return
		}

		err := db.UpdateArticle(ctx, oldName.Selected, newName.Text)
		if err != nil {
			dialog.ShowError(ErrUpdArtName, w)
			return
		} else {
			dialog.ShowInformation("Изменить имя", "Название статьи успешно изменено!", w)
		}
		err = UpdateArticleTable(db, table)
		if err != nil {
			dialog.ShowError(ErrUpdArt, w)
			return
		}

		dirContent, err := MainDir(w, db, role)
		if err != nil {
			dialog.ShowError(ErrShowDir, w)
		}
		w.SetContent(dirContent)

	})

	return container.NewVBox(fieldsCont, btn)
}

// РАЗДЕЛ УДАЛЕНИЯ СТАТЬИ
func DelArticle(w fyne.Window, db database.Service, table *widget.Table, role string) *fyne.Container {
	winDelArticle := WinDelArticle(w, db, table, role)
	return container.NewVBox(canvas.NewLine(color.White), winDelArticle)
}
func WinDelArticle(w fyne.Window, db database.Service, table *widget.Table, role string) *fyne.Container {
	ctx := context.Background()

	article := MadeSelectArticle(w, db)

	cont := container.NewAdaptiveGrid(2, widget.NewLabel("Статья"), article)

	btn := widget.NewButton("Удалить статью", func() {

		if article.Selected == "" {
			dialog.ShowError(ErrEmptyArt, w)
			return
		}

		err := db.DeleteArticle(ctx, article.Selected)
		if err != nil {
			dialog.ShowError(ErrDelArt, w)
			return
		} else {
			dialog.ShowInformation("Удалить статью", "Сатья успешно удалена!", w)
		}
		err = UpdateArticleTable(db, table)
		if err != nil {
			dialog.ShowError(ErrUpdArt, w)
			return
		}

		dirContent, err := MainDir(w, db, role)
		if err != nil {
			dialog.ShowError(ErrShowDir, w)
		}
		w.SetContent(dirContent)

	})

	return container.NewVBox(cont, btn)
}

// СПИСОК ДЕЙСТВИЙ ДЛЯ ОПЕРАЦИЙ
func AccordionDirOper(w fyne.Window, db database.Service, table *widget.Table) (*widget.Accordion, error) {
	accAdd := AddOperation(w, db, table)
	accEdit := EditOperation(w, db, table)
	accDel := DelOperation(w, db, table)

	editor := widget.NewAccordion(
		widget.NewAccordionItem("Добавить", accAdd),
		widget.NewAccordionItem("Редактировать", accEdit),
		widget.NewAccordionItem("Удалить", accDel),
	)
	return editor, nil
}

// РАЗДЕЛ ДОБАВИТЬ ОПЕРАЦИЮ
func AddOperation(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	winAddOperation := WinAddOperation(w, db, table)
	return container.NewVBox(canvas.NewLine(color.White), winAddOperation)
}
func WinAddOperation(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	ctx := context.Background()

	article := MadeSelectArticle(w, db)
	date := widget.NewEntry()
	debit := widget.NewEntry()
	credit := widget.NewEntry()

	date.SetPlaceHolder("2024-11-03")
	debit.SetPlaceHolder("0")
	credit.SetPlaceHolder("1000")

	cont := container.NewAdaptiveGrid(
		2,
		widget.NewLabel("Статья"), article,
		widget.NewLabel("Дата"), date,
		widget.NewLabel("Доход"), debit,
		widget.NewLabel("расход"), credit,
	)

	btn := widget.NewButton("Добавить операцию", func() {

		floatDebit, err := strconv.ParseFloat(debit.Text, 64)
		if err != nil {
			dialog.ShowError(ErrParseDebit, w)
			return
		}

		floatCredit, err := strconv.ParseFloat(credit.Text, 64)
		if err != nil {
			dialog.ShowError(ErrParseCredit, w)
			return
		}

		if article.Selected == "" {
			dialog.ShowError(ErrEmptyArt, w)
			return
		}

		err = db.AddOperation(ctx, article.Selected, floatDebit, floatCredit, date.Text)
		if err != nil {
			dialog.ShowError(ErrAddOp, w)
			return
		} else {
			dialog.ShowInformation("Добавить операцию", "Новая операция успешно добавлена!", w)
		}

		err = UpdateOperationTable(db, table)
		if err != nil {
			dialog.ShowError(ErrUpdOp, w)
			return
		}

	})

	return container.NewVBox(cont, btn)
}

// РАЗДЕЛ РЕДАКТИРОВАТЬ ОПЕРАЦИЮ
func EditOperation(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	winEditOperation := WinEditOperation(w, db, table)
	winIncOperation := WinIncreaseOperation(w, db, table)
	return container.NewVBox(canvas.NewLine(color.White), winEditOperation, canvas.NewLine(color.White), winIncOperation)
}
func WinEditOperation(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	ctx := context.Background()

	id := widget.NewEntry()
	article := MadeSelectArticle(w, db)
	debit := widget.NewEntry()
	credit := widget.NewEntry()

	id.SetPlaceHolder("66")
	debit.SetPlaceHolder("0")
	credit.SetPlaceHolder("666")

	cont := container.NewAdaptiveGrid(
		2,
		widget.NewLabel("ID операции"), id,
		widget.NewLabel("Статья"), article,
		widget.NewLabel("Доход"), debit,
		widget.NewLabel("расход"), credit,
	)

	btn := widget.NewButton("Изменить операцию", func() {

		intId, err := strconv.ParseInt(id.Text, 0, 0)
		if err != nil {
			dialog.ShowError(ErrParseId, w)
			return
		}

		floatDebit, err := strconv.ParseFloat(debit.Text, 64)
		if err != nil {
			dialog.ShowError(ErrParseDebit, w)
			return
		}

		floatCredit, err := strconv.ParseFloat(credit.Text, 64)
		if err != nil {
			dialog.ShowError(ErrParseCredit, w)
			return
		}

		if article.Selected == "" {
			dialog.ShowError(ErrEmptyArt, w)
			return
		}

		err = db.UpdateOpertions(ctx, int(intId), article.Selected, floatDebit, floatCredit)
		if err != nil {
			dialog.ShowError(ErrUpdOpData, w)
			return
		} else {
			dialog.ShowInformation("Изменить операцию", "Операция успешно изменена!", w)
		}

		err = UpdateOperationTable(db, table)
		if err != nil {
			dialog.ShowError(ErrUpdOp, w)
			return
		}

	})

	return container.NewVBox(cont, btn)
}
func WinIncreaseOperation(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	ctx := context.Background()

	article := MadeSelectArticle(w, db)
	amount := widget.NewEntry()
	amount.SetPlaceHolder("100")

	cont := container.NewAdaptiveGrid(
		2,
		widget.NewLabel("Статья"), article,
		widget.NewLabel("Сумма повышения"), amount,
	)

	btn := widget.NewButton("Повысить расходы по статье", func() {

		floatAmount, err := strconv.ParseFloat(amount.Text, 64)
		if err != nil {
			dialog.ShowError(ErrParseAmount, w)
			return
		}

		if article.Selected == "" {
			dialog.ShowError(ErrEmptyArt, w)
			return
		}

		err = db.IncreaseExpensesForArticle(ctx, article.Selected, floatAmount)
		if err != nil {
			dialog.ShowError(ErrIncArtCredit, w)
			return
		} else {
			dialog.ShowInformation("Повысить расходы по статье", "Расход по статье успешно изменён!", w)
		}
		err = UpdateOperationTable(db, table)
		if err != nil {
			dialog.ShowError(ErrUpdOp, w)
			return
		}

	})

	return container.NewVBox(cont, btn)
}

// РАЗДЕЛ УДАЛЕНИЯ ОПЕРАЦИЮ
func DelOperation(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	winDelOperation := WinDelOperation(w, db, table)
	return container.NewVBox(canvas.NewLine(color.White), winDelOperation)
}
func WinDelOperation(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	ctx := context.Background()

	id := widget.NewEntry()
	id.SetPlaceHolder("id")

	cont := container.NewAdaptiveGrid(
		2,
		widget.NewLabel("ID операции"), id,
	)

	btn := widget.NewButton("Удалить операцию", func() {

		intId, err := strconv.ParseInt(id.Text, 0, 0)
		if err != nil {
			dialog.ShowError(ErrParseId, w)
			return
		}

		err = db.DeleteOperation(ctx, int(intId))
		if err != nil {
			dialog.ShowError(ErrDelOp, w)
			return
		} else {
			dialog.ShowInformation("Удалить операцию", "Операция успешно удалена", w)
		}
		err = UpdateOperationTable(db, table)
		if err != nil {
			dialog.ShowError(ErrUpdOp, w)
			return
		}

	})

	return container.NewVBox(cont, btn)
}
