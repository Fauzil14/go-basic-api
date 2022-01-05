package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"quiz-3/models"
	"quiz-3/query"
	"quiz-3/utils"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func validateBook(book models.Book) string {
	var errstrings []string
	err1 := fmt.Errorf("json:Image URL yang anda masukkan tidak valid")
	err2 := fmt.Errorf("json:Tahun rilis tidak boleh kecil dari 1980 atau besar dari 2021")

	if !isValidUrl(book.ImageUrl) && (book.ReleaseYear < 1980 || book.ReleaseYear > 2021) {
		errstrings = append(errstrings, err1.Error(), err2.Error())
		return strings.Join(errstrings, "\n")
	}
	if !isValidUrl(book.ImageUrl) {
		return err1.Error()
	}
	if book.ReleaseYear < 1980 || book.ReleaseYear > 2021 {
		return err2.Error()
	}
	return ""
}

func GetBook(rw http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	// contex -> untuk oper2 data
	ctx, cancel := context.WithCancel(context.Background())
	// if context return cancel then the process got cancel
	defer cancel()

	books, err := query.GetAllBook(ctx)

	if err != nil {
		fmt.Println(err)
	}

	utils.ResponseJSON(rw, books, http.StatusOK)
}

func PostBook(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var dataToInput models.Book

	if r.Method == "POST" {
		if r.Header.Get("Content-Type") == "application/json" {
			// // decode json from request object with .NewDecoder
			decodeJson := json.NewDecoder(r.Body)
			decodeJson.DisallowUnknownFields()

			// anonymous struct to enforce mandatory fields, and other sanitation checks
			// json tag is not case sensitive
			allowedInput := struct {
				Title       *string `json:"title"`
				Description *string `json:"description"`
				ImageUrl    *string `json:"image_url"`
				ReleaseYear *int    `json:"release_year"`
				Price       *string `json:"price"`
				TotalPage   *int    `json:"total_page"`
				CategoryId  *int    `json:"category_id"`
			}{}

			err := decodeJson.Decode(&allowedInput)
			if err != nil {
				// bad JSON or unrecognized json field
				http.Error(rw, err.Error(), http.StatusBadRequest)
				return
			}
			dataToInput.Title = *allowedInput.Title
			dataToInput.Description = *allowedInput.Description
			dataToInput.ImageUrl = *allowedInput.ImageUrl
			dataToInput.ReleaseYear = *allowedInput.ReleaseYear
			dataToInput.Price = *allowedInput.Price
			dataToInput.TotalPage = *allowedInput.TotalPage
			dataToInput.CategoryId = *allowedInput.CategoryId
		} else {
			dataToInput.Title = r.FormValue("title")
			dataToInput.Description = r.FormValue("description")
			dataToInput.ImageUrl = r.FormValue("image_url")
			getReleaseYear, _ := strconv.Atoi(r.FormValue("release_year"))
			dataToInput.Price = r.FormValue("price")
			getTotalPage, _ := strconv.Atoi(r.FormValue("total_page"))
			getCategoryId, _ := strconv.Atoi(r.FormValue("category_id"))
			dataToInput.ReleaseYear = getReleaseYear
			dataToInput.TotalPage = getTotalPage
			dataToInput.CategoryId = getCategoryId
		}

		err := validateBook(dataToInput)
		if err != "" {
			http.Error(rw, err, http.StatusUnprocessableEntity)
			return
		}

		if err := query.InsertBook(ctx, dataToInput); err != nil {
			utils.ResponseJSON(rw, err, http.StatusInternalServerError)
			return
		}

		res := map[string]string{
			"Status": "Successfully",
		}

		utils.ResponseJSON(rw, res, http.StatusOK)
		return
	}

	utils.ResponseJSON(rw, "Method Not Allowed", http.StatusBadRequest)

}

func UpdateBook(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var dataToInput models.Book

	// pengecekan inputannya json atau bukan
	if r.Header.Get("Content-Type") == "application/json" {
		// // decode json from request object with .NewDecoder
		decodeJson := json.NewDecoder(r.Body)
		decodeJson.DisallowUnknownFields()

		// anonymous struct to enforce mandatory fields, and other sanitation checks
		// json tag is not case sensitive
		allowedInput := struct {
			Title       *string `json:"title"`
			Description *string `json:"description"`
			ImageUrl    *string `json:"image_url"`
			ReleaseYear *int    `json:"release_year"`
			Price       *string `json:"price"`
			TotalPage   *int    `json:"total_page"`
			CategoryId  *int    `json:"category_id"`
		}{}

		err := decodeJson.Decode(&allowedInput)
		if err != nil {
			// bad JSON or unrecognized json field
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		dataToInput.Title = *allowedInput.Title
		dataToInput.Description = *allowedInput.Description
		dataToInput.ImageUrl = *allowedInput.ImageUrl
		dataToInput.ReleaseYear = *allowedInput.ReleaseYear
		dataToInput.Price = *allowedInput.Price
		dataToInput.TotalPage = *allowedInput.TotalPage
		dataToInput.CategoryId = *allowedInput.CategoryId
	} else {
		dataToInput.Title = r.FormValue("title")
		dataToInput.Description = r.FormValue("description")
		dataToInput.ImageUrl = r.FormValue("image_url")
		getReleaseYear, _ := strconv.Atoi(r.FormValue("release_year"))
		dataToInput.Price = r.FormValue("price")
		getTotalPage, _ := strconv.Atoi(r.FormValue("total_page"))
		getCategoryId, _ := strconv.Atoi(r.FormValue("category_id"))
		dataToInput.ReleaseYear = getReleaseYear
		dataToInput.TotalPage = getTotalPage
		dataToInput.CategoryId = getCategoryId
	}

	validateBook(dataToInput)
	var idBook = ps.ByName("id")

	if err := query.UpdateBook(ctx, dataToInput, idBook); err != nil {
		utils.ResponseJSON(rw, err, http.StatusInternalServerError)
		return
	}

	res := map[string]string{
		"Status": "Successfully",
	}

	utils.ResponseJSON(rw, res, http.StatusOK)
}

func DeleteBook(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var idBook = ps.ByName("id")

	if err := query.DeleteBook(ctx, idBook); err != nil {
		errDelete := map[string]string{
			"error": fmt.Sprintf("%v", err),
		}
		utils.ResponseJSON(rw, errDelete, http.StatusInternalServerError)
		return
	}

	res := map[string]string{
		"Status": "Successfully",
	}

	utils.ResponseJSON(rw, res, http.StatusOK)
}
