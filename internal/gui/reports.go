package gui

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/EmptyInsid/db_gui/internal/database"
	"github.com/jung-kurt/gofpdf"
)

func LoadArticles(articlesContainer *fyne.Container) []string {
	// Сбор всех статей из контейнера
	articles := []string{}
	for _, obj := range articlesContainer.Objects {
		if entry, ok := obj.(*widget.Select); ok {
			articles = append(articles, strings.TrimSpace(entry.Selected))
			log.Print(entry.Selected)
		}
	}
	return articles
}
func MadeArticlesButton(db database.Service) (*fyne.Container, *widget.Button, *widget.Button, error) {
	// Поля для ввода статей
	articlesList, err := db.GetAllArticles(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}
	var arts []string
	for _, art := range articlesList {
		arts = append(arts, art.Name)
	}
	articlesContainer := container.NewVBox()
	initialArticle := widget.NewSelect(arts, func(value string) {
		log.Println("Select set to", value)
	})
	articlesContainer.Add(initialArticle)

	// Кнопка для добавления новых полей
	addArticleButton := widget.NewButton("Добавить статью", func() {
		initialArticle := widget.NewSelect(arts, func(value string) {
			log.Println("Select set to", value)
		})
		articlesContainer.Add(initialArticle)
		articlesContainer.Refresh()
	})

	// Кнопка для удаления новых полей
	delArticleButton := widget.NewButton("Удалить статью", func() {
		count := len(articlesContainer.Objects) - 1
		if count == 0 {
			return
		}
		articlesContainer.Remove(articlesContainer.Objects[count])
		articlesContainer.Refresh()
	})

	return articlesContainer, addArticleButton, delArticleButton, nil
}
func MadeTitle(titleText string) *canvas.Text {
	title := canvas.NewText(titleText, color.RGBA{R: 135, G: 206, B: 250, A: 255})
	title.TextStyle = fyne.TextStyle{Bold: true}
	return title
}

func MadeDateFields() (*widget.Entry, *widget.Entry) {
	startDate := widget.NewEntry()
	endDate := widget.NewEntry()
	startDate.SetPlaceHolder("2024-11-01")
	endDate.SetPlaceHolder("2024-11-30")

	return startDate, endDate
}

func MainReportFirst(w fyne.Window, db database.Service) (*fyne.Container, error) {

	title := MadeTitle("Отчёт 1")
	articlesContainer, addArticleButton, delArticleButton, err := MadeArticlesButton(db)
	if err != nil {
		return nil, err
	}
	startDate, endDate := MadeDateFields()

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
		delArticleButton,
	)

	// Контейнер для таблицы
	tableContainer := container.NewStack()

	// Функция обновления таблицы
	updateTable := func() {
		// Сбор всех статей из контейнера
		articles := LoadArticles(articlesContainer)

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
	previewButton.OnTapped = updateTable

	// Кнопка "Сохранить"
	saveButton := widget.NewButton("Сохранить", func() {
		dialog.ShowFileSave(

			func(uc fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if uc == nil {
					return // Пользователь отменил выбор
				}

				defer uc.Close()

				articles := LoadArticles(articlesContainer)
				data, err := db.GetIncomeExpenseDynamics(context.Background(), articles, startDate.Text, endDate.Text)
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

	toolbar := container.NewHBox(previewButton, saveButton)
	rightPane := container.NewVBox(inputContainer, toolbar)
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
	title := MadeTitle("Отчёт 2")
	articlesContainer, addArticleButton, delArticleButton, err := MadeArticlesButton(db)
	if err != nil {
		return nil, err
	}
	startDate, endDate := MadeDateFields()

	flow := widget.NewSelect([]string{"credit", "debit", "profit"}, func(value string) {
		log.Printf("Set flow: %s\n", value)
	})

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
		delArticleButton,
	)
	tableContainer := container.NewStack()

	// Функция обновления таблицы
	updateTable := func() {
		// Сбор всех статей из контейнера
		articles := LoadArticles(articlesContainer)

		// Генерация новой таблицы на основе введённых данных
		newTable, err := FinancialPercentagesTable(db, articles, flow.Selected, startDate.Text, endDate.Text)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Очистка старого содержимого контейнера и добавление новой таблицы
		tableContainer.Objects = []fyne.CanvasObject{newTable}
		tableContainer.Refresh() // Обновление отображения
	}
	previewButton := widget.NewButton("Превью", nil)
	previewButton.OnTapped = updateTable

	saveButton := widget.NewButton("Сохранить", func() {
		dialog.ShowFileSave(

			func(uc fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if uc == nil {
					return // Пользователь отменил выбор
				}

				defer uc.Close()

				articles := LoadArticles(articlesContainer)
				ctx := context.Background()
				data, err := db.GetFinancialPercentages(ctx, articles, flow.Selected, startDate.Text, endDate.Text)
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
				err = SaveToPDFSecond(data, filename)
				if err != nil {
					dialog.ShowError(err, w)
				} else {
					dialog.ShowInformation("Успех", "PDF успешно сохранён!", w)
				}
			}, w)
	})

	toolbar := container.NewHBox(previewButton, saveButton)
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
	title := MadeTitle("Отчёт 2")
	startDate, endDate := MadeDateFields()

	inputContainer := container.NewVBox(
		title,
		widget.NewLabel("Введите параметры:"),
		widget.NewLabel("Начальная дата:"),
		startDate,
		widget.NewLabel("Конечная дата:"),
		endDate,
	)

	previewButton := widget.NewButton("Превью", nil)
	tableContainer := container.NewStack()

	updateTable := func() {
		// Генерация новой таблицы на основе введённых данных
		newTable, err := TotalProfitDateTable(db, startDate.Text, endDate.Text)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		tableContainer.Objects = []fyne.CanvasObject{newTable}
		tableContainer.Refresh()
	}

	previewButton.OnTapped = updateTable
	saveButton := widget.NewButton("Сохранить", func() {
		dialog.ShowFileSave(

			func(uc fyne.URIWriteCloser, err error) {

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

	toolbar := container.NewHBox(previewButton, saveButton)
	rightPane := container.NewVBox(inputContainer, toolbar)
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
	pdf.SetFont("Arial", "", 12)

	colWidths := []float64{50, 50, 50}
	headers := []string{"Name", "Debit", "Credit"}
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(10)

	// Таблица
	pdf.SetFont("Arial", "", 10)
	for i, row := range data {
		pdf.CellFormat(colWidths[i], 10, row.Date.Format("2006-01-02"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[i], 10, fmt.Sprintf("%.2f", row.TotalDebit), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[i], 10, fmt.Sprintf("%.2f", row.TotalCredit), "1", 0, "C", false, 0, "")
		pdf.Ln(10)
	}

	return pdf.OutputFileAndClose(filename)
}

func SaveToPDFSecond(data []database.FinancialPercentage, filename string) error {
	pdf := gofpdf.New("P", "mm", "A4", "") // Создаём новый PDF
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	colWidths := []float64{40, 30, 30, 30, 30}
	headers := []string{"Name", "Total_Debit", "Total_Credit", "Total_Profit", "Total_Proc"}
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(10)

	// Таблица
	pdf.SetFont("Arial", "", 10)
	for _, row := range data {
		pdf.CellFormat(40, 10, row.ArticleName, "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", row.TotalDebit), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", row.TotalCredit), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", row.TotalProfit), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", row.TotalProc), "1", 0, "C", false, 0, "")
		pdf.Ln(10)
	}

	return pdf.OutputFileAndClose(filename)
}

func SaveToPDFThird(data []database.DateProfit, filename string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	colWidths := []float64{50, 50}
	headers := []string{"Date", "Profit"}
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(10)

	// Таблица
	pdf.SetFont("Arial", "", 10)
	for i, row := range data {
		pdf.CellFormat(colWidths[i], 10, row.Date.Format("2006-01-02"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidths[i], 10, fmt.Sprintf("%.2f", row.TotalProfit), "1", 0, "C", false, 0, "")
		pdf.Ln(10)
	}

	return pdf.OutputFileAndClose(filename)
}
