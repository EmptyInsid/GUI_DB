package database

import (
	"context"
	"time"

	"github.com/EmptyInsid/db_gui/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	CloseDB()

	AddArticle(ctx context.Context, name string) error
	AddOperation(ctx context.Context, articleName string, debit float64, credit float64, date string) error

	CreateBalanceIfProfitable(ctx context.Context, startDate, endDate string, minProfit float64) error

	DeleteArticleAndRecalculateBalances(ctx context.Context, articleName string) error
	DeleteMostUnprofitableBalance(ctx context.Context) error

	GetAllArticles(ctx context.Context) ([]models.Article, error)
	GetAllBalances(ctx context.Context) ([]models.Balance, error)
	GetAllOperations(ctx context.Context) ([]models.Operation, error)

	GetProfitByDate(ctx context.Context, startDate, endDate string) (float64, error)
	GetTotalCreditByArticleAndPeriod(ctx context.Context, articleName string, startDate, finishDate string) (float64, error)

	GetBalanceCountByArticleName(ctx context.Context, articleName string) (int, error)

	GetUnusedArticles(ctx context.Context, startData, finishData string) ([]models.Article, error)
	GetArticlesWithOperations(ctx context.Context) ([]ArticleWithOperations, error)
	GetViewUnaccountedOpertions(ctx context.Context) ([]ArticleTotalMoney, error)
	GetViewCountBalanceOper(ctx context.Context) ([]BalanceOperations, error)

	GetStoreProcLastBalanceOp(ctx context.Context) error
	GetStoreProcArticleMaxExpens(ctx context.Context, balance int, article string) error

	UpdateArticle(ctx context.Context, oldName, newName string) error
	IncreaseExpensesForArticle(ctx context.Context, articleName string, increaseAmount float64) error

	AuthUser(ctx context.Context, username, password string) (string, string, error)
	RegistrUserDB(ctx context.Context, username, password, role string) error
}

type Database struct {
	pool *pgxpool.Pool
}

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
