package queries

import (
	"context"
	"log"

	"github.com/EmptyInsid/db_gui/internal/db"
	"github.com/EmptyInsid/db_gui/internal/models"
)

// выбрать все статьи
func GetAllArticles(ctx context.Context) ([]models.Article, error) {
	rows, err := db.DB.Query(ctx, "SELECT * FROM articles")
	if err != nil {
		log.Printf("Error while get articles: %e", err)
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.ID, &article.Name); err != nil {
			log.Printf("Error while get articles: %e", err)
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// Вывести наименования всех статей, в рамках которых не проводилось операций за заданный период времени.
func GetUnusedArticles(ctx context.Context, startData, finishData string) ([]models.Article, error) {

	query := `
    SELECT name FROM articles 
    WHERE id NOT IN (SELECT operations.article_id FROM operations 
    WHERE $1 <= create_date AND create_date < $2)`

	rows, err := db.DB.Query(ctx, query, startData, finishData)
	if err != nil {
		log.Printf("Error while get unused articles: %e", err)
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		if err := rows.Scan(&article.ID, &article.Name); err != nil {
			log.Printf("Error while get articles: %e", err)
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil

}

// добавить новую статью
func AddArticle(ctx context.Context, name string) error {
	_, err := db.DB.Exec(ctx, "INSERT INTO articles(name) VALUES ($1)", name)
	return err
}

// Удалить статью и операции, выполненные в ее рамках
func DeleteArticleAndRecalculateBalances(ctx context.Context, articleName string) error {
	tx, err := db.DB.Begin(ctx)
	if err != nil {
		log.Printf("Error starting transaction: %e", err)
		return err
	}
	defer tx.Rollback(ctx)

	// Delete the article
	deleteArticleQuery := `DELETE FROM articles WHERE name = $1;`
	_, err = tx.Exec(ctx, deleteArticleQuery, articleName)
	if err != nil {
		log.Printf("Error deleting article %e", err)
		return err
	}

	// Recalculate balances
	updateBalancesQuery := `
	UPDATE balance
	SET debit = COALESCE((SELECT SUM(o.debit) FROM operations o WHERE o.balance_id = balance.id), 0),
	    credit = COALESCE((SELECT SUM(o.credit) FROM operations o WHERE o.balance_id = balance.id), 0),
	    amount = COALESCE((SELECT SUM(o.debit - o.credit) FROM operations o WHERE o.balance_id = balance.id), 0);
	`
	_, err = tx.Exec(ctx, updateBalancesQuery)
	if err != nil {
		log.Printf("Error updating balances %e", err)
		return err
	}

	return tx.Commit(ctx)
}
