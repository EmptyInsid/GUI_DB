package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EmptyInsid/db_gui/internal/database"
)

func LoadDb(config *Config) (database.Service, error) {

	// Создаем экземпляр структуры Database
	db := &database.Database{}

	// Инициализация подключения
	db.Init(buildConnectionString(config))

	// Проверка подключения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Pool().Ping(ctx); err != nil {
		log.Fatalf("Error connection with bd: %v", err)
		return nil, err
	}

	// if err := testQuery(ctx); err != nil {
	// 	db.CloseDB()
	// 	log.Fatal(err)
	// }

	return db, nil
}

func testQuery(ctx context.Context, db *database.Database) error {
	if err := db.AddArticle(ctx, "психолог"); err != nil {
		return err
	}
	fmt.Println("Succsess add article")
	if err := db.AddOperation(ctx, "кофейня", 200, 0, "2024-09-21"); err != nil {
		return err
	}
	fmt.Println("Succsess add operation")
	if err := db.CreateBalanceIfProfitable(ctx, "2024-11-01", "2024-11-30", 0); err != nil {
		return err
	}
	fmt.Println("Succsess create balance")
	if err := db.DeleteArticleAndRecalculateBalances(ctx, "транспорт"); err != nil {
		return err
	}
	fmt.Println("Succsess delete article")
	if err := db.DeleteMostUnprofitableBalance(ctx); err != nil {
		return err
	}
	fmt.Println("Succsess delete balance")

	articles, err := db.GetAllArticles(ctx)
	if err != nil {
		return err
	}
	fmt.Println("====Article====")
	for _, article := range articles {
		fmt.Printf("Article name: %s\n", article.Name)
	}

	balances, err := db.GetAllBalances(ctx)
	if err != nil {
		return err
	}
	fmt.Println("====Balance====")
	for _, balance := range balances {
		fmt.Printf(
			"Balance id: %d\tBalance date: %s\tBalance credit: %f\tBalance debit: %f\n",
			balance.ID, balance.Date, balance.Credit, balance.Debit)
	}

	operations, err := db.GetAllOperations(ctx)
	if err != nil {
		return err
	}
	fmt.Println("====Operations====")
	for _, operation := range operations {
		bid := 0
		if operation.BalanceID != nil {
			bid = *operation.BalanceID
		}
		fmt.Printf(
			"Operation id: %d\tOperation date: %s\tOperation credit: %f\tOperation debit: %f\tBalanceID: %d",
			operation.ID, operation.Date, operation.Credit, operation.Debit, bid)
	}

	artwithops, err := db.GetArticlesWithOperations(ctx)
	if err != nil {
		return err
	}
	fmt.Println("====Artwithops====")
	for _, artwithop := range artwithops {
		fmt.Printf(
			"ArticleName: %s\tCredit: %f\tDebit: %f\tOperationID: %d\tCreateDate: %s\n",
			artwithop.ArticleName, artwithop.Credit, artwithop.Debit, artwithop.OperationID, artwithop.CreateDate)
	}

	balanceCount, err := db.GetBalanceCountByArticleName(ctx, "стипендия")
	if err != nil {
		return err
	}
	fmt.Printf("====Balance count for стипендия: %d\n", balanceCount)

	profit, err := db.GetProfitByDate(ctx, "2024-09-01", "2024-09-30")
	if err != nil {
		return err
	}
	fmt.Printf("====Profit By Date 2024-09-01, 2024-09-30: %f\n", profit)

	var art string
	if err = db.GetStoreProcArticleMaxExpens(ctx, 2, art); err != nil {
		return err
	}
	fmt.Printf("====ArticleMaxExpens: %s\n", art) //questions

	if err = db.GetStoreProcLastBalanceOp(ctx); err != nil {
		return err
	}
	fmt.Println("Sucsess GetStoreProcLastBalanceOp")

	credit, err := db.GetTotalCreditByArticleAndPeriod(ctx, "развлечения", "2024-09-01", "2024-09-30")
	if err != nil {
		return err
	}
	fmt.Printf("====TotalCreditByArticleAndPeriod: %f\n", credit)

	unArticles, err := db.GetUnusedArticles(ctx, "2024-09-01", "2024-09-30")
	if err != nil {
		return err
	}
	fmt.Println("====unArticle====")
	for _, unArticle := range unArticles {
		fmt.Printf("unArticle name: %s\n", unArticle.Name)
	}

	balanceOpers, err := db.GetViewCountBalanceOper(ctx)
	if err != nil {
		return err
	}
	fmt.Println("====CountBalanceOper====")
	for _, balanceOper := range balanceOpers {
		fmt.Printf(
			"balanceOper id: %d\tbalanceOper count: %d\tbalanceOper date: %s\n",
			balanceOper.BalanceId, balanceOper.OperationCount, balanceOper.BalanceDate)
	}

	unacOpers, err := db.GetViewUnaccountedOpertions(ctx)
	if err != nil {
		return err
	}
	fmt.Println("====UnaccountedOpertions====")
	for _, unacOper := range unacOpers {
		fmt.Printf(
			"unacOper name: %s\tunacOper totCredit: %f\tunacOper totDebit: %f\n",
			unacOper.ArticleName, unacOper.TotalCredit, unacOper.TotlaDebit)
	}

	// if err = db.IncreaseExpensesForArticle(ctx, "продукты", 100); err != nil { //это не будет работать, скипай
	// 	return err
	// }
	fmt.Println("Succsess increase article")
	if err = db.UpdateArticle(ctx, "профессия", "техника"); err != nil {
		return err
	}
	fmt.Println("Succsess update article")

	return nil
}

func buildConnectionString(config *Config) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
		config.DBSSLMode,
	)
}
