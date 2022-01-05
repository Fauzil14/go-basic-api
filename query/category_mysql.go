package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"quiz-3/config"
	"quiz-3/models"
	"time"
)

const (
	category_table = "categories"
)

func connectDB() *sql.DB {
	db, err := config.MySQL()
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}

func GetAllCategory(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category

	db := connectDB()
	queryTxt := fmt.Sprintf("SELECT * FROM %v Order By created_at DESC", category_table)
	rowQuery, err := db.QueryContext(ctx, queryTxt)

	if err != nil {
		log.Fatal(err)
	}

	for rowQuery.Next() {
		var category models.Category
		var createdAt, updatedAt string
		// rowQueryScan harus berurutan sesuai susunana kolom di database
		if err = rowQuery.Scan(&category.ID,
			&category.Name,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		category.CreatedAt, err = time.Parse(layoutDateTime, createdAt)

		if err != nil {
			log.Fatal(err)
		}

		category.UpdatedAt, err = time.Parse(layoutDateTime, updatedAt)

		if err != nil {
			log.Fatal(err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func InsertCategory(ctx context.Context, category models.Category) error {
	db := connectDB()

	//query
	queryText := fmt.Sprintf("INSERT INTO %v (name, created_at, updated_at) VALUES ('%v', NOW(), NOW())", category_table, category.Name)
	_, err := db.ExecContext(ctx, queryText)

	if err != nil {
		return err
	}

	return nil
}

func UpdateCategory(ctx context.Context, category models.Category, idCategory string) error {
	db := connectDB()

	//query
	queryText := fmt.Sprintf("UPDATE %v SET name='%v', updated_at=NOW() WHERE id=%v", category_table, category.Name, idCategory)
	_, err := db.ExecContext(ctx, queryText)

	if err != nil {
		return err
	}

	return nil
}

func DeleteCategory(ctx context.Context, idCategory string) error {
	db := connectDB()

	//query
	queryText := fmt.Sprintf("DELETE FROM %v WHERE id=%v", category_table, idCategory)
	s, err := db.ExecContext(ctx, queryText)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	check, err := s.RowsAffected()

	if check == 0 {
		return errors.New("id tidak ada")
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
