package adreslibro

import (
    "html/template"
    "net/http"
)

func init() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/esperantuloj/chrisfosterelli", handler2)
}

func handler(w http.ResponseWriter, r *http.Request) {
    tem := template.Must(template.ParseFiles("template/index.html"))
    tem.Execute(w, r.FormValue("Content"))
}

func handler2(w http.ResponseWriter, r *http.Request) {
    tem := template.Must(template.ParseFiles("template/user.html"))
    tem.Execute(w, r.FormValue("Content"))
}
