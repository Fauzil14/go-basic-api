package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"quiz-3/models"
	"quiz-3/query"
	"quiz-3/utils"

	"github.com/julienschmidt/httprouter"
)

func GetCategory(rw http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	// contex -> untuk oper2 data
	ctx, cancel := context.WithCancel(context.Background())
	// if context return cancel then the process got cancel
	defer cancel()

	categories, err := query.GetAllCategory(ctx)

	if err != nil {
		fmt.Println(err)
	}

	utils.ResponseJSON(rw, categories, http.StatusOK)
}

func PostCategory(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// pengecekan inputannya json atau bukan
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(rw, "Guanakan Content-Type application/json", http.StatusBadRequest)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var category models.Category

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.ResponseJSON(rw, err, http.StatusBadRequest)
		return
	}

	if err := query.InsertCategory(ctx, category); err != nil {
		utils.ResponseJSON(rw, err, http.StatusInternalServerError)
		return
	}

	res := map[string]string{
		"Status": "Successfully",
	}

	utils.ResponseJSON(rw, res, http.StatusOK)
}

func UpdateCategory(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// pengecekan inputannya json atau bukan
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(rw, "Guanakan Content-Type application/json", http.StatusBadRequest)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var category models.Category

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.ResponseJSON(rw, err, http.StatusBadRequest)
		return
	}

	var idCategory = ps.ByName("id")

	if err := query.UpdateCategory(ctx, category, idCategory); err != nil {
		utils.ResponseJSON(rw, err, http.StatusInternalServerError)
		return
	}

	res := map[string]string{
		"Status": "Successfully",
	}

	utils.ResponseJSON(rw, res, http.StatusOK)
}

func DeleteCategory(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var idCategory = ps.ByName("id")

	if err := query.DeleteCategory(ctx, idCategory); err != nil {
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
