package db

import "time"

type ArticleWithOperations struct {
	ArticleID   int
	ArticleName string
	OperationID int
	Debit       float64
	Credit      float64
	CreateDate  time.Time // NULL, если операции нет
}

type ArticleTotalMoney struct {
	ArticleName string
	TotlaDebit  float64
	TotalCredit float64
}

type BalanceOperations struct {
	BalanceId      int
	BalanceDate    time.Time
	OperationCount int
}
