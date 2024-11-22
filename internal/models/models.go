package models

import "time"

// Article представляет статью доходов или расходов
type Article struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Operation представляет операцию (доход/расход)
type Operation struct {
	ID        int       `json:"id"`
	ArticleID int       `json:"article_id"`
	Debit     float64   `json:"debit"`
	Credit    float64   `json:"credit"`
	Date      time.Time `json:"create_date"`
	BalanceID *int      `json:"balance_id"`
}

// Balance представляет баланс за месяц
type Balance struct {
	ID     int       `json:"id"`
	Date   time.Time `json:"create_date"`
	Debit  float64   `json:"debit"`
	Credit float64   `json:"credit"`
	Amount float64   `json:"amount"`
}
