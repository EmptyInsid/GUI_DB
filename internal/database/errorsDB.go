package database

import "errors"

var (
	ErrEmptyRow    = errors.New("Return empty row")
	ErrLessThenMin = errors.New("Balance profit less then minimum")

	ErrGetProfit = errors.New("Error while getting profit")
	ErrGetCredit = errors.New("Error while getting credit")
)
