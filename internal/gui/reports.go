package gui

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/EmptyInsid/db_gui/internal/database"
	"github.com/signintech/gopdf"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func MainReportFirst(w fyne.Window, db database.Service) (*fyne.Container, error) {

	title := MadeTitle("Отчёт 1 Динамика изменения расходов и доходов.")
	articlesContainer, addArticleButton, delArticleButton, err := MadeArticlesButton(db)
	if err != nil {
		return nil, err
	}
	startDate, endDate := MadeDateFields()

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

	tableContainer := container.NewStack()
	updateTable := func() {
		articles := LoadArticles(articlesContainer)

		if err := CompareDate(startDate.Text, endDate.Text); err != nil {
			dialog.ShowError(ErrEndLessStart, w)
			return
		}

		newTable, err := IncomeExpenseDynamicsTable(db, articles, startDate.Text, endDate.Text)
		if err != nil {
			dialog.ShowError(ErrIncomeExpence, w)
			return
		}
		tableContainer.Objects = []fyne.CanvasObject{newTable}
		tableContainer.Refresh() // Обновление отображения
	}

	previewButton := widget.NewButton("Превью", nil)
	previewButton.OnTapped = updateTable

	saveButton := widget.NewButton("Сохранить", func() {
		dialog.ShowFileSave(

			func(uc fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(ErrSaveFile, w)
					return
				}
				if uc == nil {
					return // Пользователь отменил выбор
				}

				defer uc.Close()

				if err := CompareDate(startDate.Text, endDate.Text); err != nil {
					dialog.ShowError(ErrEndLessStart, w)
					return
				}

				articles := LoadArticles(articlesContainer)
				data, err := db.GetIncomeExpenseDynamics(context.Background(), articles, startDate.Text, endDate.Text)
				if err != nil {
					dialog.ShowError(ErrIncomeExpence, w)
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
					log.Printf("Error while save first to pdf %v", err)
					dialog.ShowError(ErrSaveFile, w)
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
	title := MadeTitle("Отчёт 2 Процентное соотношение финансовых потоков по статьям.")
	articlesContainer, addArticleButton, delArticleButton, err := MadeArticlesButton(db)
	if err != nil {
		return nil, err
	}
	startDate, endDate := MadeDateFields()

	flow := widget.NewSelect([]string{"расход", "доход", "прибыль"}, func(value string) {
		log.Printf("Set flow: %s\n", value)
	})

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
		if err := CompareDate(startDate.Text, endDate.Text); err != nil {
			dialog.ShowError(ErrEndLessStart, w)
			return
		}

		articles := LoadArticles(articlesContainer)

		newTable, err := FinancialPercentagesTable(db, articles, TranslateFlow(flow.Selected), startDate.Text, endDate.Text)
		if err != nil {
			dialog.ShowError(ErrFinPercTable, w)
			return
		}

		tableContainer.Objects = []fyne.CanvasObject{newTable}
		tableContainer.Refresh() // Обновление отображения
	}
	previewButton := widget.NewButton("Превью", nil)
	previewButton.OnTapped = updateTable

	saveButton := widget.NewButton("Сохранить", func() {
		dialog.ShowFileSave(

			func(uc fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(ErrSaveFile, w)
					return
				}
				if uc == nil {
					return // Пользователь отменил выбор
				}

				defer uc.Close()

				if err := CompareDate(startDate.Text, endDate.Text); err != nil {
					dialog.ShowError(ErrEndLessStart, w)
					return
				}

				articles := LoadArticles(articlesContainer)
				ctx := context.Background()
				data, err := db.GetFinancialPercentages(ctx, articles, TranslateFlow(flow.Selected), startDate.Text, endDate.Text)
				if err != nil {
					dialog.ShowError(ErrFinPercTable, w)
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
					log.Printf("jkasdfkahfg %v", err)
					dialog.ShowError(ErrSaveFile, w)
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
	title := MadeTitle("Отчёт 3 Чистая прибыль бюджета от времени.")
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
		if err := CompareDate(startDate.Text, endDate.Text); err != nil {
			dialog.ShowError(ErrEndLessStart, w)
			return
		}

		newTable, err := TotalProfitDateTable(db, startDate.Text, endDate.Text)
		if err != nil {
			dialog.ShowError(ErrTotalProfTable, w)
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
					dialog.ShowError(ErrSaveFile, w)
					return
				}
				if uc == nil {
					return // Пользователь отменил выбор
				}

				defer uc.Close()

				if err := CompareDate(startDate.Text, endDate.Text); err != nil {
					dialog.ShowError(ErrEndLessStart, w)
					return
				}

				ctx := context.Background()
				data, err := db.GetTotalProfitDate(ctx, startDate.Text, endDate.Text)
				if err != nil {
					dialog.ShowError(ErrTotalProfTable, w)
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
					log.Printf("jkasdfkahfg %v", err)
					dialog.ShowError(ErrSaveFile, w)
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

	pdf, err := createPdf()
	if err != nil {
		return err
	}

	headers := []string{"Дата", "Общий расход", "Общий доход"}
	tableStartY := 10.0
	marginLeft := 10.0

	// Create a new table layout
	table := pdf.NewTableLayout(marginLeft, tableStartY, 25, 5)
	for _, row := range headers {
		table.AddColumn(row, 100, "left")
	}
	// Add rows to the table
	for _, row := range data {
		if row.Date.Format("2006-01-02") != "" {
			table.AddRow([]string{
				row.Date.Format("2006-01-02"),
				fmt.Sprintf("%.2f", row.TotalCredit),
				fmt.Sprintf("%.2f", row.TotalDebit),
			})
		}
		log.Println(row.Date.Format("2006-01-02"))

	}

	table.SetTableStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:    true,
			Left:   true,
			Bottom: true,
			Right:  true,
			Width:  1.0,
		},
		FillColor: gopdf.RGBColor{R: 255, G: 0, B: 255},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		FontSize:  10,
	})

	// Set the style for table header
	table.SetHeaderStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:      true,
			Left:     true,
			Bottom:   true,
			Right:    true,
			Width:    2.0,
			RGBColor: gopdf.RGBColor{R: 100, G: 150, B: 255},
		},
		FillColor: gopdf.RGBColor{R: 255, G: 200, B: 200},
		TextColor: gopdf.RGBColor{R: 255, G: 100, B: 100},
		Font:      "Arial",
		FontSize:  12,
	})

	table.SetCellStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:      true,
			Left:     true,
			Bottom:   true,
			Right:    true,
			Width:    0.5,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		Font:      "Arial",
		FontSize:  10,
	})

	// Draw the table
	table.DrawTable()

	// Сохраняем график как изображение
	plotFile := "chart.png"
	err = createFirstPlot(data, plotFile)
	if err != nil {
		return fmt.Errorf("ошибка создания графика: %v", err)
	}

	// Добавляем новую страницу в PDF
	pdf.AddPage()

	// Добавляем график на новую страницу
	pdf.Image(plotFile, 20, 50, &gopdf.Rect{W: 400, H: 400}) // Размещение графика на второй странице

	// Сохраняем PDF в файл
	err = pdf.WritePdf(filename)
	if err != nil {
		return fmt.Errorf("ошибка сохранения PDF: %v", err)
	}

	os.Remove(plotFile)

	return nil
}

func SaveToPDFSecond(data []database.FinancialPercentage, filename string) error {

	// Создаём новый PDF
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4}) // Размер страницы A4
	pdf.AddPage()

	err := pdf.AddTTFFont("Arial", "C:/Windows/Fonts/arial.ttf")
	if err != nil {
		return fmt.Errorf("ошибка добавления шрифта: %v", err)
	}

	// Устанавливаем шрифт и его размер
	err = pdf.SetFont("Arial", "", 12)
	if err != nil {
		return fmt.Errorf("ошибка установки шрифта: %v", err)
	}

	headers := []string{"Статья", "Общий расход", "Общий доход", "Прибыль", "Процент"}
	// Set the starting Y position for the table
	tableStartY := 10.0
	// Set the left margin for the table
	marginLeft := 10.0

	// Create a new table layout
	table := pdf.NewTableLayout(marginLeft, tableStartY, 25, 5)

	// Add columns to the table
	for _, row := range headers {
		table.AddColumn(row, 100, "left")
	}

	// Add rows to the table
	for _, row := range data {
		if row.ArticleName != "" {
			table.AddRow([]string{
				row.ArticleName,
				fmt.Sprintf("%.2f", row.TotalCredit),
				fmt.Sprintf("%.2f", row.TotalDebit),
				fmt.Sprintf("%.2f", row.TotalProfit),
				fmt.Sprintf("%.2f", row.TotalProc),
			})
		}
	}

	// Set the style for table cells
	table.SetTableStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:    true,
			Left:   true,
			Bottom: true,
			Right:  true,
			Width:  1.0,
		},
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		FontSize:  10,
	})

	// Set the style for table header
	table.SetHeaderStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:      true,
			Left:     true,
			Bottom:   true,
			Right:    true,
			Width:    2.0,
			RGBColor: gopdf.RGBColor{R: 100, G: 150, B: 255},
		},
		FillColor: gopdf.RGBColor{R: 255, G: 200, B: 200},
		TextColor: gopdf.RGBColor{R: 255, G: 100, B: 100},
		Font:      "font2",
		FontSize:  12,
	})

	table.SetCellStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:      true,
			Left:     true,
			Bottom:   true,
			Right:    true,
			Width:    0.5,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		Font:      "font1",
		FontSize:  10,
	})

	// Draw the table
	table.DrawTable()

	// Сохраняем PDF в файл
	err = pdf.WritePdf(filename)
	if err != nil {
		return fmt.Errorf("ошибка сохранения PDF: %v", err)
	}

	return nil
}

func SaveToPDFThird(data []database.DateProfit, filename string) error {

	pdf, err := createPdf()
	if err != nil {
		return err
	}

	headers := []string{"Дата", "Прибыль"}
	tableStartY := 10.0
	marginLeft := 10.0

	// Create a new table layout
	table := pdf.NewTableLayout(marginLeft, tableStartY, 25, 5)
	for _, row := range headers {
		table.AddColumn(row, 100, "left")
	}

	for _, row := range data {
		if row.Date.Format("2006-01-02") != "" {
			table.AddRow([]string{
				row.Date.Format("2006-01-02"),
				fmt.Sprintf("%.2f", row.TotalProfit)})
		}
	}

	table.SetTableStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:    true,
			Left:   true,
			Bottom: true,
			Right:  true,
			Width:  1.0,
		},
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		FontSize:  10,
	})

	// Set the style for table header
	table.SetHeaderStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:      true,
			Left:     true,
			Bottom:   true,
			Right:    true,
			Width:    2.0,
			RGBColor: gopdf.RGBColor{R: 100, G: 150, B: 255},
		},
		FillColor: gopdf.RGBColor{R: 255, G: 200, B: 200},
		TextColor: gopdf.RGBColor{R: 255, G: 100, B: 100},
		Font:      "Arial",
		FontSize:  12,
	})

	table.SetCellStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:      true,
			Left:     true,
			Bottom:   true,
			Right:    true,
			Width:    0.5,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		Font:      "Arial",
		FontSize:  10,
	})

	// Draw the table
	table.DrawTable()

	// Сохраняем график как изображение
	plotFile := "chart.png"
	err = createThirdPlot(data, plotFile)
	if err != nil {
		return fmt.Errorf("ошибка создания графика: %v", err)
	}

	// Добавляем новую страницу в PDF
	pdf.AddPage()

	// Добавляем график на новую страницу
	pdf.Image(plotFile, 20, 50, &gopdf.Rect{W: 400, H: 400}) // Размещение графика на второй странице

	// Сохраняем PDF в файл
	err = pdf.WritePdf(filename)
	if err != nil {
		return fmt.Errorf("ошибка сохранения PDF: %v", err)
	}

	os.Remove(plotFile)

	return nil
}

func createFirstPlot(data []database.DateTotalMoney, filename string) error {
	// Создаем новый график
	p := plot.New()

	// Заголовок графика
	p.Title.Text = "Credit and Debit Over Time"
	p.X.Label.Text = "Date"
	p.Y.Label.Text = "Amount"

	// Создаем графики для кредитов и дебетов
	creditPoints := make(plotter.XYs, len(data))
	debitPoints := make(plotter.XYs, len(data))

	// Сохраняем даты для подписей
	dates := make([]string, len(data))

	for i, row := range data {
		if row.Date.Format("2006-01-02") != "" {
			creditPoints[i].X = float64(i)
			creditPoints[i].Y = row.TotalCredit
			debitPoints[i].X = float64(i)
			debitPoints[i].Y = row.TotalDebit
		}

		dates[i] = row.Date.Format("2006-01-02") // Форматируем дату для подписи
	}

	// Добавляем линии для кредитов и дебетов
	creditLine, err := plotter.NewLine(creditPoints)
	if err != nil {
		return fmt.Errorf("ошибка создания линии кредита: %v", err)
	}
	creditLine.Color = color.RGBA{R: 255, G: 0, B: 0} // Красный для кредита

	debitLine, err := plotter.NewLine(debitPoints)
	if err != nil {
		return fmt.Errorf("ошибка создания линии дебета: %v", err)
	}
	debitLine.Color = color.RGBA{R: 0, G: 0, B: 255} // Синий для дебета

	// Добавляем линии в график
	p.Add(creditLine, debitLine)

	// Настраиваем подписи оси X (даты)
	p.X.Tick.Marker = plot.ConstantTicks(ticksForDates(dates))

	// Настраиваем оси
	p.Y.Tick.Marker = plot.DefaultTicks{}

	// Сохраняем график как PNG
	err = p.Save(400, 400, filename)
	if err != nil {
		return fmt.Errorf("ошибка сохранения графика: %v", err)
	}

	return nil
}

func createThirdPlot(data []database.DateProfit, filename string) error {
	// Создаем новый график
	p := plot.New()

	// Заголовок графика
	p.Title.Text = "Profit Over Time"
	p.X.Label.Text = "Date"
	p.Y.Label.Text = "Amount"

	// Создаем графики для кредитов и дебетов
	creditPoints := make(plotter.XYs, len(data))

	// Сохраняем даты для подписей
	dates := make([]string, len(data))

	for i, row := range data {
		if row.Date.Format("2006-01-02") != "" {
			creditPoints[i].X = float64(i)
			creditPoints[i].Y = row.TotalProfit
		}

		dates[i] = row.Date.Format("2006-01-02") // Форматируем дату для подписи
	}

	// Добавляем линии для кредитов и дебетов
	creditLine, err := plotter.NewLine(creditPoints)
	if err != nil {
		return fmt.Errorf("ошибка создания линии кредита: %v", err)
	}
	creditLine.Color = color.RGBA{R: 255, G: 255, B: 0} // Красный для кредита

	// Добавляем линии в график
	p.Add(creditLine)

	// Настраиваем подписи оси X (даты)
	p.X.Tick.Marker = plot.ConstantTicks(ticksForDates(dates))

	// Настраиваем оси
	p.Y.Tick.Marker = plot.DefaultTicks{}

	// Сохраняем график как PNG
	err = p.Save(400, 400, filename)
	if err != nil {
		return fmt.Errorf("ошибка сохранения графика: %v", err)
	}

	return nil
}

func ticksForDates(dates []string) []plot.Tick {
	var ticks []plot.Tick
	for i, date := range dates {
		ticks = append(ticks, plot.Tick{
			Value: float64(i), // Координата точки на оси X
			Label: date,       // Подпись (дата)
		})
	}
	return ticks
}

func createPdf() (*gopdf.GoPdf, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4}) // Размер страницы A4
	pdf.AddPage()

	err := pdf.AddTTFFont("Arial", "C:/Windows/Fonts/arial.ttf")
	if err != nil {
		return nil, fmt.Errorf("ошибка добавления шрифта: %v", err)
	}

	err = pdf.SetFont("Arial", "", 12)
	if err != nil {
		return nil, fmt.Errorf("ошибка установки шрифта: %v", err)
	}

	return pdf, nil
}
