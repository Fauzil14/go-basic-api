package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"quiz-3/models"
	"time"
)

const (
	book_table     = "books"
	layoutDateTime = "2006-01-02 15:04:05"
)

func GetAllBook(ctx context.Context) ([]models.Book, error) {
	var books []models.Book

	db := connectDB()
	queryTxt := fmt.Sprintf("SELECT * FROM %v Order By created_at DESC", book_table)
	rowQuery, err := db.QueryContext(ctx, queryTxt)

	if err != nil {
		log.Fatal(err)
	}

	for rowQuery.Next() {
		var book models.Book
		var createdAt, updatedAt string
		// rowQueryScan harus berurutan sesuai susunana kolom di database
		if err = rowQuery.Scan(&book.ID,
			&book.Title,
			&book.Description,
			&book.ImageUrl,
			&book.ReleaseYear,
			&book.Price,
			&book.TotalPage,
			&book.Thickness,
			&createdAt,
			&updatedAt,
			&book.CategoryId,
		); err != nil {
			return nil, err
		}

		book.CreatedAt, err = time.Parse(layoutDateTime, createdAt)

		if err != nil {
			log.Fatal(err)
		}

		book.UpdatedAt, err = time.Parse(layoutDateTime, updatedAt)

		if err != nil {
			log.Fatal(err)
		}

		books = append(books, book)
	}

	return books, nil
}

func InsertBook(ctx context.Context, book models.Book) error {
	db := connectDB()

	//query
	queryText := fmt.Sprintf("INSERT INTO %v (title, description, image_url, release_year, price, total_page, thickness, created_at, updated_at, category_id) VALUES ('%v', '%v', '%v', %v, '%v', %v, '%v', NOW(), NOW(), %v)", book_table, book.Title, book.Description, book.ImageUrl, book.ReleaseYear, book.Price, book.TotalPage, book.Thickness, book.CategoryId)
	_, err := db.ExecContext(ctx, queryText)

	if err != nil {
		return err
	}

	return nil
}

func UpdateBook(ctx context.Context, book models.Book, idBook string) error {
	db := connectDB()

	//query
	queryText := fmt.Sprintf("UPDATE %v SET title='%v', description='%v', image_url='%v',release_year=%v, price='%v',total_page=%v,thickness='%v',category_id=%v,updated_at=NOW() WHERE id=%v", book_table, book.Title, book.Description, book.ImageUrl, book.ReleaseYear, book.Price, book.TotalPage, book.Thickness, book.CategoryId, idBook)
	_, err := db.ExecContext(ctx, queryText)

	if err != nil {
		return err
	}

	return nil
}

func DeleteBook(ctx context.Context, idBook string) error {
	db := connectDB()

	//query
	queryText := fmt.Sprintf("DELETE FROM %v WHERE id=%v", book_table, idBook)
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
