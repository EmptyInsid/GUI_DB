package queries

import (
	"context"
	"log"

	"github.com/EmptyInsid/db_gui/internal/db"
	"github.com/EmptyInsid/db_gui/internal/models"
)

type ArticleWithOperations struct {
	ArticleID   int
	ArticleName string
	OperationID *int // NULL, если операции нет
	Debit       *float64
	Credit      *float64
	CreateDate  *string // NULL, если операции нет
}

func GetAllOperations(ctx context.Context) ([]models.Operation, error) {
	rows, err := db.DB.Query(ctx, "SELECT * FROM operations")
	if err != nil {
		log.Printf("Error while get operations: %e", err)
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
			log.Printf("Error while get operations: %e", err)
			return nil, err
		}
		operations = append(operations, operation)
	}

	return operations, nil
}

// GetArticlesWithOperations fetches articles and their associated operations
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
	LEFT JOIN 
		operations 
	ON 
		articles.id = operations.article_id
	ORDER BY 
		articles.name, operations.create_date;
	`

	rows, err := db.DB.Query(ctx, query)
	if err != nil {
		log.Printf("Error fetching articles with operations: %e", err)
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
			log.Printf("Error scanning row: %e", err)
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

	err := db.DB.QueryRow(context.Background(), query, startDate, endDate).Scan(&totalProfit)
	if err != nil {
		log.Printf("Error while get operations: %e", err)
		return 0, err
	}

	return totalProfit, nil
}

// Добавить операцию в рамках статьи
func AddOperation(ctx context.Context, articleName string, debit float64, credit float64, date string) error {
	query := `
		INSERT INTO operations(article_id, debit, credit, create_date, balance_id) VALUES
		((SELECT id FROM articles WHERE articles.name = $1), $2, $3, $3, NULL)`

	_, err := db.DB.Exec(ctx, query, articleName, debit, credit, date)
	return err
}

// Увеличить сумму расхода операций для статьи, заданной по наименованию
func IncreaseExpensesForArticle(ctx context.Context, articleName string, increaseAmount float64) error {
	tx, err := db.DB.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %e", err)
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
		log.Printf("Error upgrading operations: %e", err)
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
		log.Printf("Error updating balances: %e", err)
		return err
	}

	return tx.Commit(ctx)
}
