package gui

import (
	"context"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/EmptyInsid/db_gui/internal/auth"
	"github.com/EmptyInsid/db_gui/internal/database"
)

// первое окно входа
func LoginMenu(w fyne.Window, db database.Service) {
	btnLogin := widget.NewButton("Войти", func() {
		log.Println("User clicked 'Войти'")
		AuthForm(w, db)
	})

	// btnRegistr := widget.NewButton("Зарегистрироваться", func() {
	// 	log.Println("User clicked 'Зарегистрироваться'")
	// 	//RegistrationForm(w, db)
	// })

	w.Resize(fyne.NewSize(250, 250))
	w.SetContent(container.NewCenter(btnLogin))
}

// форма входа
func AuthForm(w fyne.Window, db database.Service) {
	login := widget.NewEntry()
	password := widget.NewPasswordEntry()

	form := widget.NewForm(
		widget.NewFormItem("Логин", login),
		widget.NewFormItem("Пароль", password),
	)

	form.SubmitText = "Вход"
	form.OnSubmit = func() {
		log.Printf("User %s is trying to login", login.Text)

		ctx := context.Background()

		_, role, err := auth.AuthenticateUser(db, ctx, login.Text, password.Text)
		if err != nil {
			log.Printf("Failed to fetch user names: %v", err)
			dialog.ShowError(ErrAuth, w)
			return
		}

		MainWindow(w, db, role)
	}

	form.CancelText = "Отмена"
	form.OnCancel = func() {
		log.Println("User clicked 'Отмена'")
		w.Close()
	}

	w.SetContent(container.NewCenter(form))
}
