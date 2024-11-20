package queries

import (
	"context"
	"log"

	"github.com/EmptyInsid/db_gui/internal/db"
	"github.com/EmptyInsid/db_gui/internal/models"
)

func GetAllBalances(ctx context.Context) ([]models.Balance, error) {
	rows, err := db.DB.Query(ctx, "SELECT * FROM balance")
	if err != nil {
		log.Printf("Error while get balances: %e", err)
		return nil, err
	}
	defer rows.Close()

	var balances []models.Balance
	for rows.Next() {
		var balance models.Balance
		if err := rows.Scan(
			&balance.ID,
			&balance.Debit,
			&balance.Credit,
			&balance.Date,
		); err != nil {
			log.Printf("Error while get balances: %e", err)
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
	if err := db.DB.QueryRow(ctx, query, articleName).Scan(&balanceCount); err != nil {
		log.Printf("Error fetching balance count: %e", err)
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
	if err := db.DB.QueryRow(ctx, query, articleName, startDate, finishDate).Scan(&profit); err != nil {
		log.Printf("Error fetching total credit: %e", err)
		return 0, err
	}
	return profit, nil
}

func CreateBalanceIfProfitable(ctx context.Context, startDate, endDate string, minProfit float64) error {
	tx, err := db.DB.Begin(ctx)
	if err != nil {
		log.Printf("Error failed to begin transaction: %e", err)
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
		log.Printf("Error failed to calculate debit/credit: %e", err)
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
		log.Printf("Error failed to insert balance: %e", err)
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
		log.Printf("Error failed to update operations: %e", err)
		return err
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		log.Printf("Error failed to commit transaction: %e", err)
		return err
	}

	return nil
}

// Удалить в рамках транзакции самый убыточный баланс и операции
func DeleteMostUnprofitableBalance(ctx context.Context) error {
	tx, err := db.DB.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %e", err)
		return err
	}
	defer tx.Rollback(ctx)

	// Delete the balance with the minimum amount
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
		log.Printf("Error deleting the most unprofitable balance: %e", err)
		return err
	}

	return tx.Commit(ctx)
}
