package main

import (
    "log"
	"net/http"
	"html/template"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("!adr3sl1br0s3cr3t5@"))
var tmplt = template.Must(template.ParseGlob("template/*"))

type User struct {
	Name     string
	Email    string
	Password string
	Location string
}

func getCurrentUser(r *http.Request) *User {
    session, err := mgo.Dial("mongodb://adreslibro:adreslibro@dogen.mongohq.com:10096/adreslibro_dev")
    if err != nil {
        panic(err)
    }
    defer session.Close()
    sessionz, _ := store.Get(r, "adreslibro")
    name := sessionz.Values["user"]
    c := session.DB("").C("users")
    user := User{}
    err = c.Find(bson.M{"name":name}).One(&user)
    if err != nil {
        return &user
    }
    return &user
}

func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
    data["currentUser"] = getCurrentUser(r)
    log.Println(data["currentUser"])
    err := tmplt.ExecuteTemplate(w, tmpl, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/registri", register)
	http.HandleFunc("/esperantuloj/chrisfosterelli", user)
    log.Println("Starting Server")
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, r, "index", make(map[string]interface{}))
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
        sessionz, _ := store.Get(r, "adreslibro")
        sessionz.Values["user"] = u.Name
        sessionz.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
        renderTemplate(w, r, "register", make(map[string]interface{}))
	}
}

func user(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, r, "user", make(map[string]interface{}))
}
