package database

import (
	"context"
	"time"

	"github.com/EmptyInsid/db_gui/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	CloseDB()

	AddArticle(ctx context.Context, name string) error                                                      //справочник статей +
	AddOperation(ctx context.Context, articleName string, debit float64, credit float64, date string) error //справочник операций +

	CreateBalanceIfProfitable(ctx context.Context, startDate, endDate string, minProfit float64) error //журнал +

	DeleteArticle(ctx context.Context, articleName string) error //справочник статей +
	DeleteOperation(ctx context.Context, id int) error           //справочник операций +
	DeleteMostUnprofitableBalance(ctx context.Context) error     //журнал +

	GetAllArticles(ctx context.Context) ([]models.Article, error)     //справочник статей +
	GetAllBalances(ctx context.Context) ([]models.Balance, error)     //журнал +
	GetAllOperations(ctx context.Context) ([]models.Operation, error) //

	GetProfitByDate(ctx context.Context, startDate, endDate string) (float64, error)                                         //журнал +
	GetTotalCreditByArticleAndPeriod(ctx context.Context, articleName string, startDate, finishDate string) (float64, error) //журнал +
	GetBalanceCountByArticleName(ctx context.Context, articleName string) (int, error)                                       //журнал +

	GetUnusedArticles(ctx context.Context, startData, finishData string) ([]models.Article, error) //справочник статей
	GetArticlesWithOperations(ctx context.Context) ([]ArticleWithOperations, error)                //справочник операций +
	GetViewUnaccountedOpertions(ctx context.Context) ([]ArticleTotalMoney, error)                  //
	GetViewCountBalanceOper(ctx context.Context) ([]BalanceOperations, error)                      //

	GetStoreProcLastBalanceOp(ctx context.Context) error                                 //
	GetStoreProcArticleMaxExpens(ctx context.Context, balance int, article string) error //

	UpdateArticle(ctx context.Context, oldName, newName string) error                                                  //справочник статей +
	UpdateOpertions(ctx context.Context, id int, articleName string, debit float64, credit float64, date string) error //справочник операций +
	IncreaseExpensesForArticle(ctx context.Context, articleName string, increaseAmount float64) error                  //справочник операций +

	AuthUser(ctx context.Context, username, password string) (string, string, error) //вход +
	RegistrUserDB(ctx context.Context, username, password, role string) error        //регистрация -
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
	TotalDebit  float64
	TotalCredit float64
}

type BalanceOperations struct {
	BalanceId      int
	BalanceDate    time.Time
	OperationCount int
}
