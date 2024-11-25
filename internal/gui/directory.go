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

func MainDir(w fyne.Window, db database.Service) (*fyne.Container, error) {
	dirContent, err := TabsDir(w, db)
	if err != nil {
		return nil, err
	}
	return container.NewStack(dirContent), nil
}

func ArticleViewer(w fyne.Window, db database.Service) (*container.Split, error) {
	table, err := ArticleTable(db)
	if err != nil {
		return nil, err
	}
	editor, err := AccordionDirArticle(w, db, table)
	if err != nil {
		return nil, err
	}

	return GridViewer(db, table, editor), nil
}

func OperationsViewer(w fyne.Window, db database.Service) (*container.Split, error) {
	table, err := OperationsTable(db)
	if err != nil {
		return nil, err
	}
	editor, err := AccordionDirOper(w, db, table)
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

// СПИСОК ДЕЙСТВИЙ ДЛЯ СТАТЕЙ
func AccordionDirArticle(w fyne.Window, db database.Service, table *widget.Table) (*widget.Accordion, error) {
	accAdd := AddArticle(w, db, table)
	accEdit := EditArticle(w, db, table)
	accDel := DelArticle(w, db, table)

	editor := widget.NewAccordion(
		widget.NewAccordionItem("Добавить", accAdd),
		widget.NewAccordionItem("Редактировать", accEdit),
		widget.NewAccordionItem("Удалить", accDel),
	)
	return editor, nil
}

// РАЗДЕЛ ДОБАВИТЬ СТАТЬЮ
func AddArticle(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	winAddArticle := WinAddArticle(w, db, table)
	return container.NewVBox(canvas.NewLine(color.White), winAddArticle)
}
func WinAddArticle(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	ctx := context.Background()

	article := widget.NewEntry()
	article.SetPlaceHolder("статья")

	btn := widget.NewButton("Добавить статью", func() {

		err := db.AddArticle(ctx, article.Text)
		if err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("Добавить статью", "Новая статья успешно добавлена!", w)
		}

		err = UpdateArticleTable(db, table)
		if err != nil {
			dialog.ShowError(err, w)
		}

	})

	return container.NewVBox(article, btn)
}

// РАЗДЕЛ РЕДАКТИРОВАТЬ СТАТЬЮ
func EditArticle(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	winEditArticle := WinEditArticle(w, db, table)
	return container.NewVBox(canvas.NewLine(color.White), winEditArticle)
}
func WinEditArticle(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	ctx := context.Background()

	oldName := widget.NewEntry()
	newName := widget.NewEntry()
	oldName.SetPlaceHolder("старое имя")
	newName.SetPlaceHolder("новое имя")
	fieldsCont := container.NewStack(container.NewAdaptiveGrid(2, oldName, newName))

	btn := widget.NewButton("Изменить имя", func() {

		err := db.UpdateArticle(ctx, oldName.Text, newName.Text)
		if err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("Изменить имя", "Название статьи успешно изменено!", w)
		}
		err = UpdateArticleTable(db, table)
		if err != nil {
			dialog.ShowError(err, w)
		}

	})

	return container.NewVBox(fieldsCont, btn)
}

// РАЗДЕЛ УДАЛЕНИЯ СТАТЬИ
func DelArticle(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	winDelArticle := WinDelArticle(w, db, table)
	return container.NewVBox(canvas.NewLine(color.White), winDelArticle)
}
func WinDelArticle(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	ctx := context.Background()

	article := widget.NewEntry()
	article.SetPlaceHolder("статья")

	btn := widget.NewButton("Удалить статью", func() {

		err := db.DeleteArticle(ctx, article.Text)
		if err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("Удалить статью", "Сатья успешно удалена!", w)
		}
		err = UpdateArticleTable(db, table)
		if err != nil {
			dialog.ShowError(err, w)
		}

	})

	return container.NewVBox(article, btn)
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

	article := widget.NewEntry()
	date := widget.NewEntry()
	debit := widget.NewEntry()
	credit := widget.NewEntry()

	article.SetPlaceHolder("статья")
	date.SetPlaceHolder("дата")
	debit.SetPlaceHolder("доход")
	credit.SetPlaceHolder("расход")

	сont := container.NewStack(container.NewAdaptiveGrid(2, article, date, debit, credit))

	btn := widget.NewButton("Добавить операцию", func() {

		floatDebit, err := strconv.ParseFloat(debit.Text, 64)
		if err != nil {
			dialog.ShowError(err, w)
		}

		floatCredit, err := strconv.ParseFloat(credit.Text, 64)
		if err != nil {
			dialog.ShowError(err, w)
		}

		err = db.AddOperation(ctx, article.Text, floatDebit, floatCredit, date.Text)
		if err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("Добавить операцию", "Новая операция успешно добавлена!", w)
		}

		err = UpdateOperationTable(db, table)
		if err != nil {
			dialog.ShowError(err, w)
		}

	})

	return container.NewVBox(сont, btn)
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
	article := widget.NewEntry()
	debit := widget.NewEntry()
	credit := widget.NewEntry()

	id.SetPlaceHolder("id")
	article.SetPlaceHolder("статья")
	debit.SetPlaceHolder("доход")
	credit.SetPlaceHolder("расход")

	сont := container.NewStack(container.NewAdaptiveGrid(2, id, article, debit, credit))

	btn := widget.NewButton("Изменить операцию", func() {

		intId, err := strconv.ParseInt(id.Text, 0, 0)
		if err != nil {
			dialog.ShowError(err, w)
		}

		floatDebit, err := strconv.ParseFloat(debit.Text, 64)
		if err != nil {
			dialog.ShowError(err, w)
		}

		floatCredit, err := strconv.ParseFloat(credit.Text, 64)
		if err != nil {
			dialog.ShowError(err, w)
		}

		err = db.UpdateOpertions(ctx, int(intId), article.Text, floatDebit, floatCredit)
		if err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("Изменить операцию", "Операция успешно изменена!", w)
		}

		err = UpdateOperationTable(db, table)
		if err != nil {
			dialog.ShowError(err, w)
		}

	})

	return container.NewVBox(сont, btn)
}
func WinIncreaseOperation(w fyne.Window, db database.Service, table *widget.Table) *fyne.Container {
	ctx := context.Background()

	article := widget.NewEntry()
	amount := widget.NewEntry()
	article.SetPlaceHolder("название статьи")
	amount.SetPlaceHolder("сумма повышения")
	сont := container.NewStack(container.NewAdaptiveGrid(1, article, amount))

	btn := widget.NewButton("Повысить расходы по статье", func() {

		floatAmount, err := strconv.ParseFloat(amount.Text, 64)
		if err != nil {
			dialog.ShowError(err, w)
		}

		err = db.IncreaseExpensesForArticle(ctx, article.Text, floatAmount)
		if err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("Повысить расходы по статье", "Расход по статье успешно изменён!", w)
		}
		err = UpdateOperationTable(db, table)
		if err != nil {
			dialog.ShowError(err, w)
		}

	})

	return container.NewVBox(сont, btn)
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

	btn := widget.NewButton("Удалить статью", func() {

		intId, err := strconv.ParseInt(id.Text, 0, 0)
		if err != nil {
			dialog.ShowError(err, w)
		}

		err = db.DeleteOperation(ctx, int(intId))
		if err != nil {
			dialog.ShowError(err, w)
		} else {
			dialog.ShowInformation("Удалить операцию", "Надо оно тебе?", w)
		}
		err = UpdateOperationTable(db, table)
		if err != nil {
			dialog.ShowError(err, w)
		}

	})

	return container.NewVBox(id, btn)
}
