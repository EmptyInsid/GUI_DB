package gui

import "time"

func isEndOfMonth(t time.Time) bool {
	// Проверяем, совпадает ли текущий день с последним днем месяца
	return t.Day() == t.AddDate(0, 0, 1).Day()-1
}

func getStartOfMonth(t time.Time) time.Time {
	// Устанавливаем день месяца равным 1
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}
