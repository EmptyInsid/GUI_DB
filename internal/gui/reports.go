package gui

import (
	"context"
	"fmt"
	"image/color"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/EmptyInsid/db_gui/internal/database"
	"github.com/jung-kurt/gofpdf"
)

func MainReportFirst(w fyne.Window, db database.Service) (*fyne.Container, error) {
	// Заголовок
	title := canvas.NewText("Отчёт 1", color.RGBA{R: 135, G: 206, B: 250, A: 255})
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Поля для ввода статей
	articlesContainer := container.NewVBox()
	initialArticle := widget.NewEntry()
	initialArticle.SetPlaceHolder("Введите статью")
	articlesContainer.Add(initialArticle)

	// Кнопка для добавления новых полей
	addArticleButton := widget.NewButton("Добавить статью", func() {
		newArticle := widget.NewEntry()
		newArticle.SetPlaceHolder("Введите статью")
		articlesContainer.Add(newArticle)
		articlesContainer.Refresh()
	})

	// Поля ввода дат
	startDate := widget.NewEntry()
	endDate := widget.NewEntry()
	startDate.SetPlaceHolder("2024-11-01")
	endDate.SetPlaceHolder("2024-11-30")

	// Контейнер для ввода данных
	inputContainer := container.NewVBox(
		title,
		widget.NewLabel("Введите параметры:"),
		widget.NewLabel("Начальная дата:"),
		startDate,
		widget.NewLabel("Конечная дата:"),
		endDate,
		widget.NewLabel("Статьи:"),
		articlesContainer,
		addArticleButton,
	)

	// Контейнер для таблицы
	tableContainer := container.NewStack()

	// Функция обновления таблицы
	updateTable := func() {
		// Сбор всех статей из контейнера
		articles := []string{}
		for _, obj := range articlesContainer.Objects {
			if entry, ok := obj.(*widget.Entry); ok {
				articles = append(articles, entry.Text)
			}
		}

		// Генерация новой таблицы на основе введённых данных
		newTable, err := IncomeExpenseDynamicsTable(db, articles, startDate.Text, endDate.Text)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Очистка старого содержимого контейнера и добавление новой таблицы
		tableContainer.Objects = []fyne.CanvasObject{newTable}
		tableContainer.Refresh() // Обновление отображения
	}

	// Кнопка "Превью"
	previewButton := widget.NewButton("Превью", nil)

	// Привязываем функцию обновления к кнопке "Превью"
	previewButton.OnTapped = updateTable

	// Кнопка "Сохранить"
	saveButton := widget.NewButton("Сохранить", func() {
		dialog.ShowFileSave(

			func(uc fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
				}

				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if uc == nil {
					return // Пользователь отменил выбор
				}

				defer uc.Close()

				articles := []string{}
				for _, obj := range articlesContainer.Objects {
					if entry, ok := obj.(*widget.Entry); ok {
						articles = append(articles, entry.Text)
					}
				}
				ctx := context.Background()
				data, err := db.GetIncomeExpenseDynamics(ctx, articles, startDate.Text, endDate.Text)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}

				// Получаем путь к файлу
				filename := uc.URI().Path()

				// Проверяем расширение и добавляем, если отсутствует
				if filepath.Ext(filename) != ".pdf" {
					filename += ".pdf"
				}

				// Сохраняем в PDF
				err = SaveToPDFFirst(data, filename)
				if err != nil {
					dialog.ShowError(err, w)
				} else {
					dialog.ShowInformation("Успех", "PDF успешно сохранён!", w)
				}
			}, w)
	})

	// Верхний тулбар
	toolbar := container.NewHBox(previewButton, saveButton)

	// Правый элемент: ввод данных
	rightPane := container.NewVBox(inputContainer, toolbar)

	// Центральная область: таблица
	mainContent := container.NewHSplit(tableContainer, rightPane)
	mainContent.SetOffset(0.7) // Устанавливает пропорцию (70% для таблицы, 30% для правой панели)

	// Основной контейнер
	mainContainer := container.NewBorder(
		nil,         // Верхняя часть
		nil,         // Нижняя часть
		nil,         // Левая часть
		nil,         // Правая часть
		mainContent, // Центральный контент
	)

	return mainContainer, nil
}

func MainReportSecond(w fyne.Window, db database.Service) (*fyne.Container, error) {
	// Заголовок
	title := canvas.NewText("Отчёт 2", color.RGBA{R: 135, G: 206, B: 250, A: 255})
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Поля для ввода статей
	articlesContainer := container.NewVBox()
	initialArticle := widget.NewEntry()
	initialArticle.SetPlaceHolder("Введите статью")
	articlesContainer.Add(initialArticle)

	// Кнопка для добавления новых полей
	addArticleButton := widget.NewButton("Добавить статью", func() {
		newArticle := widget.NewEntry()
		newArticle.SetPlaceHolder("Введите статью")
		articlesContainer.Add(newArticle)
		articlesContainer.Refresh()
	})

	// Поля ввода дат
	startDate := widget.NewEntry()
	endDate := widget.NewEntry()
	flow := widget.NewEntry()
	startDate.SetPlaceHolder("2024-11-01")
	endDate.SetPlaceHolder("2024-11-30")
	flow.SetPlaceHolder("поток")

	// Контейнер для ввода данных
	inputContainer := container.NewVBox(
		title,
		widget.NewLabel("Введите параметры:"),
		widget.NewLabel("Начальная дата:"),
		startDate,
		widget.NewLabel("Конечная дата:"),
		endDate,
		widget.NewLabel("Тип потока:"),
		flow,
		widget.NewLabel("Статьи:"),
		articlesContainer,
		addArticleButton,
	)

	// Кнопка "Превью"
	previewButton := widget.NewButton("Превью", nil)

	// Контейнер для таблицы
	tableContainer := container.NewStack()

	// Функция обновления таблицы
	updateTable := func() {
		// Сбор всех статей из контейнера
		articles := []string{}
		for _, obj := range articlesContainer.Objects {
			if entry, ok := obj.(*widget.Entry); ok {
				articles = append(articles, entry.Text)
			}
		}

		// Генерация новой таблицы на основе введённых данных
		newTable, err := FinancialPercentagesTable(db, articles, flow.Text, startDate.Text, endDate.Text)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Очистка старого содержимого контейнера и добавление новой таблицы
		tableContainer.Objects = []fyne.CanvasObject{newTable}
		tableContainer.Refresh() // Обновление отображения
	}

	// Привязываем функцию обновления к кнопке "Превью"
	previewButton.OnTapped = updateTable

	// Кнопка "Сохранить"
	saveButton := widget.NewButton("Сохранить", func() {
		dialog.ShowInformation("Сохранение", "Отчёт успешно сохранён.", w)
	})

	// Верхний тулбар
	toolbar := container.NewHBox(previewButton, saveButton)

	// Правый элемент: ввод данных
	rightPane := container.NewVBox(inputContainer, toolbar)

	// Центральная область: таблица
	mainContent := container.NewHSplit(tableContainer, rightPane)
	mainContent.SetOffset(0.7) // Устанавливает пропорцию (70% для таблицы, 30% для правой панели)

	// Основной контейнер
	mainContainer := container.NewBorder(
		nil,         // Верхняя часть
		nil,         // Нижняя часть
		nil,         // Левая часть
		nil,         // Правая часть
		mainContent, // Центральный контент
	)

	return mainContainer, nil
}
func MainReportThird(w fyne.Window, db database.Service) (*fyne.Container, error) {
	// Заголовок
	title := canvas.NewText("Отчёт 3", color.RGBA{R: 135, G: 206, B: 250, A: 255})
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Поля ввода дат
	startDate := widget.NewEntry()
	endDate := widget.NewEntry()
	startDate.SetPlaceHolder("2024-11-01")
	endDate.SetPlaceHolder("2024-11-30")

	// Контейнер для ввода данных
	inputContainer := container.NewVBox(
		title,
		widget.NewLabel("Введите параметры:"),
		widget.NewLabel("Начальная дата:"),
		startDate,
		widget.NewLabel("Конечная дата:"),
		endDate,
	)

	// Кнопка "Превью"
	previewButton := widget.NewButton("Превью", nil)

	// Контейнер для таблицы
	tableContainer := container.NewStack()

	// Функция обновления таблицы
	updateTable := func() {
		// Генерация новой таблицы на основе введённых данных
		newTable, err := TotalProfitDateTable(db, startDate.Text, endDate.Text)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Очистка старого содержимого контейнера и добавление новой таблицы
		tableContainer.Objects = []fyne.CanvasObject{newTable}
		tableContainer.Refresh() // Обновление отображения
	}

	// Привязываем функцию обновления к кнопке "Превью"
	previewButton.OnTapped = updateTable

	// Кнопка "Сохранить"
	saveButton := widget.NewButton("Сохранить", func() {
		dialog.ShowFileSave(

			func(uc fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
				}

				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if uc == nil {
					return // Пользователь отменил выбор
				}

				defer uc.Close()

				ctx := context.Background()
				data, err := db.GetTotalProfitDate(ctx, startDate.Text, endDate.Text)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}

				// Получаем путь к файлу
				filename := uc.URI().Path()

				// Проверяем расширение и добавляем, если отсутствует
				if filepath.Ext(filename) != ".pdf" {
					filename += ".pdf"
				}

				// Сохраняем в PDF
				err = SaveToPDFThird(data, filename)
				if err != nil {
					dialog.ShowError(err, w)
				} else {
					dialog.ShowInformation("Успех", "PDF успешно сохранён!", w)
				}
			}, w)
	})

	// Верхний тулбар
	toolbar := container.NewHBox(previewButton, saveButton)

	// Правый элемент: ввод данных
	rightPane := container.NewVBox(inputContainer, toolbar)

	// Центральная область: таблица
	mainContent := container.NewHSplit(tableContainer, rightPane)
	mainContent.SetOffset(0.7) // Устанавливает пропорцию (70% для таблицы, 30% для правой панели)

	// Основной контейнер
	mainContainer := container.NewBorder(
		nil,         // Верхняя часть
		nil,         // Нижняя часть
		nil,         // Левая часть
		nil,         // Правая часть
		mainContent, // Центральный контент
	)

	return mainContainer, nil

}

func SaveToPDFFirst(data []database.DateTotalMoney, filename string) error {
	pdf := gofpdf.New("P", "mm", "A4", "") // Создаём новый PDF
	pdf.AddPage()

	// Устанавливаем шрифт
	pdf.SetFont("Arial", "B", 12)

	// Заголовок таблицы
	pdf.Cell(40, 10, "Date")
	pdf.Cell(40, 10, "Debit")
	pdf.Cell(40, 10, "Credit")
	pdf.Ln(10) // Переход на следующую строку

	// Таблица
	pdf.SetFont("Arial", "", 10)
	for _, row := range data {
		pdf.Cell(40, 10, row.Date.Format("2006-01-02"))
		pdf.Cell(40, 10, fmt.Sprintf("%2f", row.TotalDebit))
		pdf.Cell(40, 10, fmt.Sprintf("%2f", row.TotalCredit))
		pdf.Ln(10)
	}

	// Сохраняем PDF
	return pdf.OutputFileAndClose(filename)
}

func SaveToPDFThird(data []database.DateProfit, filename string) error {
	pdf := gofpdf.New("P", "mm", "A4", "") // Создаём новый PDF
	pdf.AddPage()

	// Устанавливаем шрифт
	pdf.SetFont("Arial", "B", 12)

	// Заголовок таблицы
	pdf.Cell(40, 10, "Date")
	pdf.Cell(40, 10, "Profit")
	pdf.Ln(10) // Переход на следующую строку

	// Таблица
	pdf.SetFont("Arial", "", 10)
	for _, row := range data {
		pdf.Cell(40, 10, row.Date.Format("2006-01-02"))
		pdf.Cell(40, 10, fmt.Sprintf("%2f", row.TotalProfit))
		pdf.Ln(10)
	}

	// Сохраняем PDF
	return pdf.OutputFileAndClose(filename)
}
