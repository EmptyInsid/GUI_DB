package gui

import (
	"context"
	"image/color"
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/EmptyInsid/db_gui/internal/database"
)

func isEndOfMonth(date time.Time) bool {
	nextDay := date.AddDate(0, 0, 1)
	return nextDay.Day() == 1
}

func getStartOfMonth(t time.Time) time.Time {
	// Устанавливаем день месяца равным 1
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func MadeSelectArticle(w fyne.Window, db database.Service) *widget.Select {
	articlesList, err := db.GetAllArticles(context.Background())
	if err != nil {
		dialog.ShowError(ErrGetArt, w)
	}
	var arts []string
	for _, art := range articlesList {
		arts = append(arts, art.Name)
	}
	return widget.NewSelect(arts, func(value string) {
		log.Println("Select set to", value)
	})
}

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

func MadeDateFields() (*widget.Entry, *widget.Entry) {
	startDate := widget.NewEntry()
	endDate := widget.NewEntry()
	startDate.SetPlaceHolder("2024-11-01")
	endDate.SetPlaceHolder("2024-11-30")

	return startDate, endDate
}

func MadeTitle(titleText string) *canvas.Text {
	title := canvas.NewText(titleText, color.RGBA{R: 135, G: 206, B: 250, A: 255})
	title.TextStyle = fyne.TextStyle{Bold: true}
	return title
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

func TranslateFlow(flow string) string {
	switch flow {
	case "расход":
		{
			return "credit"
		}
	case "доход":
		{
			return "debit"
		}
	case "прибыль":
		{
			return "profit"
		}
	default:
		{
			return "idk"
		}
	}
}

func CompareDate(start, end string) error {
	const layout = "2006-01-02"

	// Парсим строки в объекты time.Time
	dateStart, err := time.Parse(layout, start)
	if err != nil {
		log.Printf("error parse first date: %v\n", err)
		return err
	}

	dateEnd, err := time.Parse(layout, end)
	if err != nil {
		log.Printf("error parse second date:: %v\n", err)
		return err
	}

	if dateStart.After(dateEnd) {
		return ErrEndLessStart
	}
	return nil
}
