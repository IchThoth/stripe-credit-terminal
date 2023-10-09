package main

import (
	"net/http"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplates(w, r, "terminal", nil); err != nil {
		app.errorLog.Println(err)
	}
}
