package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	models "../models"
	services "../services"
)

func HandleModeration(w http.ResponseWriter, r *http.Request, params url.Values) {
	t, err := template.ParseFiles("../internal/templates/moderation.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "moderation", nil)
}

func HandleModerationSave(w http.ResponseWriter, r *http.Request, params url.Values) {
	var category models.Category
	var err error

	category.Id = services.GenerateId()
	category.Name = r.FormValue("category")
	fmt.Println("Add ?")
	err = models.AddCategory(category, models.Db)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Println(models.CategoryById(category.Id))
	http.Redirect(w, r, "/", 302)
}
