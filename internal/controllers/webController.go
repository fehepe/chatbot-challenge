package controllers

import (
	"net/http"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("../../web/*.html"))
}

func Index(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}
