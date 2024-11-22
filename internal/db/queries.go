package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EmptyInsid/db_gui/internal/models"
)

// выбрать все статьи
func GetAllArticles(ctx context.Context) ([]models.Article, error) {
	rows, err := DB.Query(ctx, "SELECT * FROM articles")
	if err != nil {
		log.Printf("Error while get articles: %v", err)
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.ID, &article.Name); err != nil {
			log.Printf("Error while get articles: %v", err)
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// Вывести наименования всех статей, в рамках которых не проводилось операций за заданный период времени.
func GetUnusedArticles(ctx context.Context, startData, finishData string) ([]models.Article, error) {

	query := `
    SELECT DISTINCT id, name FROM articles 
    WHERE id NOT IN (SELECT DISTINCT operations.article_id FROM operations 
    WHERE $1 <= create_date AND create_date < $2)`

	rows, err := DB.Query(ctx, query, startData, finishData)
	if err != nil {
		log.Printf("Error while get unused articles: %v", err)
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.ID, &article.Name); err != nil {
			log.Printf("Error while get articles: %v", err)
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil

}

// добавить новую статью
func AddArticle(ctx context.Context, name string) error {
	_, err := DB.Exec(ctx, "INSERT INTO articles(name) VALUES ($1)", name)
	return err
}

// В рамках транзакции поменять заданную статью во всех операциях на другую и удалить ее.
func UpdateArticle(ctx context.Context, oldName, newName string) error {
	query := `UPDATE articles SET name = $1 WHERE name = $2`

	commandTag, err := DB.Exec(ctx, query, newName, oldName)
	if err != nil {
		log.Printf("Error failed to update article name: %v", err)
		return err
	}

	if commandTag.RowsAffected() == 0 {
		log.Printf("Error no articles found with name: %s", oldName)
		return fmt.Errorf("error no articles found with name: %s", oldName)
	}

	return nil
}

// Удалить статью и операции, выполненные в ее рамках
func DeleteArticleAndRecalculateBalances(ctx context.Context, articleName string) error {
	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	query := `SELECT id FROM articles WHERE name = $1;`
	var articleId int
	if err := DB.QueryRow(ctx, query, articleName).Scan(&articleId); err != nil {
		log.Printf("Error fetching balance count: %v", err)
		return err
	}

	// Recalculate balances
	updateBalancesQuery := `
	UPDATE balance
	SET debit = COALESCE((SELECT SUM(o.debit) FROM operations o WHERE o.balance_id = balance.id AND o.article_id != $1), 0),
		credit = COALESCE((SELECT SUM(o.credit) FROM operations o WHERE o.balance_id = balance.id AND o.article_id != $1), 0),
		amount = COALESCE((SELECT SUM(o.debit - o.credit) FROM operations o WHERE o.balance_id = balance.id AND o.article_id != $1), 0);
	`
	_, err = tx.Exec(ctx, updateBalancesQuery, articleId)
	if err != nil {
		log.Printf("Error updating balances %v", err)
		return err
	}

	// Delete the article
	deleteArticleQuery := `DELETE FROM articles WHERE name = $1;`
	_, err = tx.Exec(ctx, deleteArticleQuery, articleName)
	if err != nil {
		log.Printf("Error deleting article %v", err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("Error commit transaction: %v\n", err)
		return err
	}

	return nil
}

// получить все балансы
func GetAllBalances(ctx context.Context) ([]models.Balance, error) {
	rows, err := DB.Query(ctx, "SELECT * FROM balance")
	if err != nil {
		log.Printf("Error while get balances: %v", err)
		return nil, err
	}
	defer rows.Close()

	var balances []models.Balance
	for rows.Next() {
		var balance models.Balance
		if err := rows.Scan(
			&balance.ID,
			&balance.Date,
			&balance.Debit,
			&balance.Credit,
			&balance.Amount,
		); err != nil {
			log.Printf("Error while get balances: %v", err)
			return nil, err
		}
		balances = append(balances, balance)
	}

	return balances, nil
}

// Вывести число балансов, в которых учтены операции, принадлежащие статье с заданным наименованием
func GetBalanceCountByArticleName(ctx context.Context, articleName string) (int, error) {
	query := `
	SELECT COUNT(DISTINCT b.id) AS balance_count
	FROM balance b
	JOIN operations o ON b.id = o.balance_id
	JOIN articles a ON o.article_id = a.id
	WHERE a.name = $1;
	`

	var balanceCount int
	if err := DB.QueryRow(ctx, query, articleName).Scan(&balanceCount); err != nil {
		log.Printf("Error fetching balance count: %v", err)
		return 0, err
	}
	return balanceCount, nil
}

// Вывести сумму расходов по заданной статье, агрегируя по балансам за указанный период
func GetTotalCreditByArticleAndPeriod(ctx context.Context, articleName string, startDate, finishDate string) (float64, error) {
	query := `
	SELECT SUM(o.credit) AS total_credit
	FROM operations o
	JOIN articles a ON o.article_id = a.id
	JOIN balance b ON o.balance_id = b.id
	WHERE a.name = $1
	  AND b.create_date BETWEEN $2 AND $3;
	`

	var profit float64
	if err := DB.QueryRow(ctx, query, articleName, startDate, finishDate).Scan(&profit); err != nil {
		log.Printf("Error fetching total credit: %v", err)
		return 0, err
	}
	return profit, nil
}

// Сформировать баланс. Если сумма прибыли меньше некоторой суммы – транзакцию откатить.
func CreateBalanceIfProfitable(ctx context.Context, startDate, endDate string, minProfit float64) error {
	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Printf("Error failed to begin transaction: %v", err)
		return err
	}

	defer tx.Rollback(ctx)

	var totalDebit, totalCredit float64

	// Calculate debit and credit for the given period
	query := `
	SELECT COALESCE(SUM(debit), 0), COALESCE(SUM(credit), 0)
	FROM operations
	WHERE create_date BETWEEN $1 AND $2
	`
	err = tx.QueryRow(ctx, query, startDate, endDate).Scan(&totalDebit, &totalCredit)
	if err != nil {
		log.Printf("Error failed to calculate debit/credit: %v", err)
		return err
	}

	profit := totalDebit - totalCredit
	if profit < minProfit {
		log.Printf("Error profit (%f) is less than the minimum required (%f)", profit, minProfit)
		return err
	}

	// Insert balance
	insertQuery := `
	INSERT INTO balance (create_date, debit, credit, amount)
	VALUES ($1, $2, $3, $4) RETURNING id
	`
	var newBalanceID int
	err = tx.QueryRow(ctx, insertQuery, endDate, totalDebit, totalCredit, profit).Scan(&newBalanceID)
	if err != nil {
		log.Printf("Error failed to insert balance: %v", err)
		return err
	}

	// Update operations with new balance ID
	updateQuery := `
	UPDATE operations
	SET balance_id = $1
	WHERE create_date BETWEEN $2 AND $3
	`
	_, err = tx.Exec(ctx, updateQuery, newBalanceID, startDate, endDate)
	if err != nil {
		log.Printf("Error failed to update operations: %v", err)
		return err
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Error failed to commit transaction: %v", err)
		return err
	}

	return nil
}

// Удалить в рамках транзакции самый убыточный баланс и операции
func DeleteMostUnprofitableBalance(ctx context.Context) error {
	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	deleteBalanceQuery := `
	DELETE FROM balance
	WHERE id = (
	    SELECT id FROM balance
	    WHERE amount = (SELECT MIN(amount) FROM balance)
	    LIMIT 1
	);
	`
	_, err = tx.Exec(ctx, deleteBalanceQuery)
	if err != nil {
		log.Printf("Error deleting the most unprofitable balance: %v", err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("Error commit transaction: %v\n", err)
		return err
	}

	return nil
}

// получить все операции
func GetAllOperations(ctx context.Context) ([]models.Operation, error) {
	rows, err := DB.Query(ctx, "SELECT * FROM operations")
	if err != nil {
		log.Printf("Error while get operations: %v", err)
		return nil, err
	}
	defer rows.Close()

	var operations []models.Operation
	for rows.Next() {
		var operation models.Operation
		if err := rows.Scan(
			&operation.ID,
			&operation.ArticleID,
			&operation.Debit,
			&operation.Credit,
			&operation.Date,
			&operation.BalanceID,
		); err != nil {
			log.Printf("Error while get operations: %v", err)
			return nil, err
		}
		operations = append(operations, operation)
	}

	return operations, nil
}

// Вывести операции и наименования статей, включая статьи, в рамках которых не проводились операции.
func GetArticlesWithOperations(ctx context.Context) ([]ArticleWithOperations, error) {
	query := `
	SELECT 
		articles.id AS article_id,
		articles.name AS article_name,
		operations.id AS operation_id,
		operations.debit,
		operations.credit,
		operations.create_date
	FROM 
		articles
	RIGHT JOIN 
		operations 
	ON 
		articles.id = operations.article_id
	ORDER BY 
		operations.create_date;
	`

	rows, err := DB.Query(ctx, query)
	if err != nil {
		log.Printf("Error fetching articles with operations: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []ArticleWithOperations
	for rows.Next() {
		var record ArticleWithOperations
		err := rows.Scan(
			&record.ArticleID,
			&record.ArticleName,
			&record.OperationID,
			&record.Debit,
			&record.Credit,
			&record.CreateDate,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		results = append(results, record)
	}

	return results, nil
}

// Посчитать прибыль за заданную дату
func GetProfitByDate(ctx context.Context, startDate, endDate string) (float64, error) {
	var totalProfit float64

	query := `
	SELECT COALESCE(SUM(debit - credit), 0)
	FROM operations
	WHERE create_date BETWEEN $1 AND $2
	`

	err := DB.QueryRow(context.Background(), query, startDate, endDate).Scan(&totalProfit)
	if err != nil {
		log.Printf("Error while get operations: %v", err)
		return 0, err
	}

	return totalProfit, nil
}

// Добавить операцию в рамках статьи
func AddOperation(ctx context.Context, articleName string, debit float64, credit float64, date string) error {

	operationDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("invalid date format: %w", err)
	}

	// Рассчитать последний день месяца для операции
	lastDayOfMonth := time.Date(operationDate.Year(), operationDate.Month()+1, 0, 0, 0, 0, 0, operationDate.Location())

	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %v\n", err)
		return err
	}
	defer tx.Rollback(ctx)

	// Найти баланс
	var balanceID *int
	selectBalanceQuery := `
		SELECT id FROM balance
		WHERE create_date = $1
		LIMIT 1`
	err = tx.QueryRow(ctx, selectBalanceQuery, lastDayOfMonth).Scan(&balanceID)
	if err != nil {
		// Если баланса нет, просто пропускаем присвоение balance_id
		log.Printf("No balance found for the operation date: %s. Balance ID will be NULL.", operationDate.Format("2006-01-02"))
		balanceID = nil
	}

	// Вставить операцию
	queryAddOp := `
	INSERT INTO operations(article_id, debit, credit, create_date, balance_id) VALUES
	((SELECT id FROM articles WHERE articles.name = $1), $2, $3, $4, $5)
	`

	if _, err := DB.Exec(ctx, queryAddOp, articleName, debit, credit, date, balanceID); err != nil {
		log.Printf("Error insert operation: %v\n", err)
		return err
	}

	if balanceID != nil {
		// Recalculate balances
		updateBalancesQuery := `
		UPDATE balance
		SET debit = COALESCE((SELECT SUM(o.debit) FROM operations o WHERE o.balance_id = balance.id), 0),
			credit = COALESCE((SELECT SUM(o.credit) FROM operations o WHERE o.balance_id = balance.id), 0),
			amount = COALESCE((SELECT SUM(o.debit - o.credit) FROM operations o WHERE o.balance_id = balance.id), 0)
		WHERE balance.id = $1;
		`
		_, err = tx.Exec(ctx, updateBalancesQuery, balanceID)
		if err != nil {
			log.Printf("Error updating balances %v", err)
			return err
		}

	}

	return nil
}

// Увеличить сумму расхода операций для статьи, заданной по наименованию
func IncreaseExpensesForArticle(ctx context.Context, articleName string, increaseAmount float64) error {
	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}
	defer tx.Rollback(context.Background())

	// Update operations for the given article
	updateOperationsQuery := `
	UPDATE operations o
	SET credit = credit + $1
	FROM articles a
	WHERE o.article_id = a.id
	AND a.name = $2;
	`
	_, err = tx.Exec(ctx, updateOperationsQuery, increaseAmount, articleName)
	if err != nil {
		log.Printf("Error upgrading operations: %v", err)
		return err
	}

	// Recalculate balances
	updateBalancesQuery := `
	UPDATE balance b
	SET debit = (
	        SELECT SUM(o.debit)
	        FROM operations o
	        WHERE o.balance_id = b.id
	    ),
	    credit = (
	        SELECT SUM(o.credit)
	        FROM operations o
	        WHERE o.balance_id = b.id
	    ),
	    amount = debit - credit
	WHERE b.id IN (
	    SELECT DISTINCT o.balance_id
	    FROM operations o
	    JOIN articles a ON o.article_id = a.id
	    WHERE a.name = $2
	);
	`
	_, err = tx.Exec(ctx, updateBalancesQuery, articleName)
	if err != nil {
		log.Printf("Error updating balances: %v", err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("Error commit transaction: %v\n", err)
		return err
	}

	return nil
}

// Создать представление, отображающее все статьи и суммы приход/расход неучтенных операций
func GetViewUnaccountedOpertions(ctx context.Context) ([]ArticleTotalMoney, error) {
	rows, err := DB.Query(ctx, "SELECT * FROM unaccounted_operations")
	if err != nil {
		log.Printf("Error getting unaccounted_operations: %v", err)
		return nil, err
	}

	var articalTotals []ArticleTotalMoney
	for rows.Next() {
		var articalTotal ArticleTotalMoney
		if err = rows.Scan(&articalTotal.ArticleName, &articalTotal.TotlaDebit, &articalTotal.TotalCredit); err != nil {
			log.Printf("Errorgetting unaccounted_operation: %v", err)
			return nil, err
		}
		articalTotals = append(articalTotals, articalTotal)
	}
	return articalTotals, nil
}

// Создать представление, отображающее все балансы и число операций, на основании которых они были сформированы
func GetViewCountBalanceOper(ctx context.Context) ([]BalanceOperations, error) {
	rows, err := DB.Query(ctx, "SELECT * FROM balance_operations_count")
	if err != nil {
		log.Printf("Error getting balance_operations_count: %v", err)
		return nil, err
	}

	var balOps []BalanceOperations
	for rows.Next() {
		var balOp BalanceOperations
		if err = rows.Scan(&balOp.BalanceId, &balOp.BalanceDate, &balOp.OperationCount); err != nil {
			log.Printf("Error getting balance_operation_count: %v", err)
			return nil, err
		}
		balOps = append(balOps, balOp)
	}
	return balOps, nil
}

// Вызвать хранимую процедуру, выводящую все операции последнего баланса и прибыли по каждой.
func GetStoreProcLastBalanceOp(ctx context.Context) error {
	if _, err := DB.Exec(ctx, "CALL get_last_balance_operations()"); err != nil {
		log.Printf("Failed to call get_last_balance_operations procedure: %v", err)
		return err
	}

	return nil
}

// Создать хранимую процедуру, имеющую два параметра «статья1» и «статья2».
// Она должна возвращать балансы, операции по «статье1» в которых составили прибыль большую, чем по «статье2».
// Если в балансе отсутствуют операции по одной из статей – он не рассматривается.
func GetStoreProcBalanceWithProfit(ctx context.Context, articleFirst, articleSecond string) error {
	if _, err := DB.Exec(ctx, "CALL get_balances_with_profit_comparison($1, $2)", articleFirst, articleSecond); err != nil {
		log.Printf("Failed to call get_balances_with_profit_comparison procedure: %v", err)
		return err
	}

	return nil
}

// Создать хранимую процедуру с входным параметром баланс и выходным параметром – статья, операции по которой проведены с наибольшими расходами
func GetStoreProcArticleMaxExpens(ctx context.Context, balance int, article string) error {
	if err := DB.QueryRow(ctx, "CALL get_article_with_max_expenses($1, $2)", balance, &article).Scan(&article); err != nil {
		log.Printf("Failed to call get_balances_with_profit_comparison procedure: %v", err)
		return err
	}

	return nil
}
