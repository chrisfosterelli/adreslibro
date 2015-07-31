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
	Id       bson.ObjectId `bson:"_id"`
	Name     string
	Email    string
	Password string
	Location string
}

//func getUser(email string) {
    //user := User{}
    //query := bson.M{"email":email}
    //err := c.Find(query).One(&user)
    //if err != nil {
        //log.Error("Could not fetch user", err)
    //}
    //return &user
//}

func getCurrentUser(r *http.Request) *User {
    mgoSession, err := mgo.Dial("mongodb://adreslibro:adreslibro@dogen.mongohq.com:10096/adreslibro_dev")
    if err != nil {
        panic(err)
    }
    defer mgoSession.Close()
    userSession, _ := store.Get(r, "adreslibro")
    stringId := userSession.Values["id"].(string)
    log.Println(stringId)
    user := User{}
    if !bson.IsObjectIdHex(stringId) {
        log.Println("Not a valid Object ID Hex")
        return &user
    }
    mgoId := bson.ObjectIdHex(stringId)
    c := mgoSession.DB("").C("users")
    err = c.FindId(mgoId).One(&user)
    log.Println(err)
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
	http.HandleFunc("/ensaluti", login)
	http.HandleFunc("/elsaluti", logout)
	http.HandleFunc("/registri", register)
	http.HandleFunc("/esperantuloj/chrisfosterelli", user)
    log.Println("Starting Server")
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, r, "index", make(map[string]interface{}))
}

func login(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
		session, err := mgo.Dial("mongodb://adreslibro:adreslibro@dogen.mongohq.com:10096/adreslibro_dev")
        if err != nil {
            panic(err)
        }
        defer session.Close()
        user := User{}
		c := session.DB("").C("users")
        err = c.Find(bson.M{"email":r.FormValue("email")}).One(&user)
        if (err != nil) {
            log.Println("not found")
            http.Redirect(w, r, "/ensaluti", http.StatusFound)
            return
        }
        if (user.Password != r.FormValue("password")) {
            log.Println("invalid pass")
            log.Println(r.FormValue("password"))
            log.Println(user)
            log.Println(user.Password)
            http.Redirect(w, r, "/ensaluti", http.StatusFound)
            return
        }
        sessionz, _ := store.Get(r, "adreslibro")
        sessionz.Values["id"] = user.Id.Hex()
        sessionz.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
    } else {
        renderTemplate(w, r, "login", make(map[string]interface{}))
    }
}

func logout(w http.ResponseWriter, r *http.Request) {
    sessionz, _ := store.Get(r, "adreslibro")
    sessionz.Values["id"] = ""
    sessionz.Save(r, w)
    http.Redirect(w, r, "/", http.StatusFound)
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
        u.Id = bson.NewObjectId()
		c := session.DB("").C("users")
		err = c.Insert(&u)
		if err != nil {
			log.Fatal(err)
		}
        log.Println("Result")
        log.Println(u.Id.Hex())
        sessionz, _ := store.Get(r, "adreslibro")
        sessionz.Values["id"] = u.Id.Hex()
        sessionz.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
        renderTemplate(w, r, "register", make(map[string]interface{}))
	}
}

func user(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, r, "user", make(map[string]interface{}))
}
