package main

import (
	"html/template"
    "gopkg.in/mgo.v2"
	"net/http"
    "log"
)

type User struct {
	Name     string
	Email    string
	Password string
	Location string
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/registri", register)
	http.HandleFunc("/esperantuloj/chrisfosterelli", user)
    log.Println("Starting Server")
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	tem := template.Must(template.ParseFiles("template/index.html"))
	tem.Execute(w, r.FormValue("Content"))
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		session, err := mgo.Dial("mongodb://adreslibro:adreslibro@dogen.mongohq.com:10096/adreslibro_dev")
		if err != nil {
			panic(err)
		}
		defer session.Close()
		u := User{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
			Location: r.FormValue("location"),
		}
		c := session.DB("").C("users")
		err = c.Insert(&u)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		tem := template.Must(template.ParseFiles("template/register.html"))
		tem.Execute(w, r.FormValue("Content"))
	}
}

func user(w http.ResponseWriter, r *http.Request) {
	tem := template.Must(template.ParseFiles("template/user.html"))
	tem.Execute(w, r.FormValue("Content"))
}
