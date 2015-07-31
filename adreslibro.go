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
var mgoSession *mgo.Session

/* Data structures */

type User struct {
	Id       bson.ObjectId `bson:"_id"`
	Name     string
	Email    string
	Password string
	Location string
}

/* Helper functions */

func getMgoSession() *mgo.Session {
    if mgoSession == nil {
        var err error
        mgoSession, err = mgo.Dial("mongodb://adreslibro:adreslibro@dogen.mongohq.com:10096/adreslibro_dev")
        if err != nil {
            panic(err)
        }
    }
    return mgoSession.Clone()
}

func getUserSession(r *http.Request) *sessions.Session {
    userSession, _ := store.Get(r, "adreslibro")
    return userSession
}

func getCurrentUser(r *http.Request) *User {
    user := User{}
    userSession := getUserSession(r)
    mgoSession := getMgoSession()
    id := userSession.Values["id"].(string)
    if bson.IsObjectIdHex(id) {
        mId := bson.ObjectIdHex(id)
        c := mgoSession.DB("").C("users")
        _ = c.FindId(mId).One(&user)
    }
    return &user
}

func renderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
    if data == nil {
        data = make(map[string]interface{})
    }
    data["currentUser"] = getCurrentUser(r)
    err := tmplt.ExecuteTemplate(w, tmpl, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

/* Main */

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/ensaluti", login)
	http.HandleFunc("/elsaluti", logout)
	http.HandleFunc("/registri", register)
	http.HandleFunc("/esperantuloj/chrisfosterelli", user)
    log.Println("Starting Server")
	http.ListenAndServe(":8080", nil)
}

/* Routes */

func index(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, r, "index", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        user := User{}
        mgoSession := getMgoSession()
		c := mgoSession.DB("").C("users")
        err := c.Find(bson.M{"email":r.FormValue("email")}).One(&user)
        if (err != nil) {
            log.Println("User not found")
            http.Redirect(w, r, "/ensaluti", http.StatusFound)
            return
        }
        if (user.Password != r.FormValue("password")) {
            log.Println("invalid pass")
            http.Redirect(w, r, "/ensaluti", http.StatusFound)
            return
        }
        userSession := getUserSession(r)
        userSession.Values["id"] = user.Id.Hex()
        userSession.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
    } else {
        renderTemplate(w, r, "login", nil)
    }
}

func logout(w http.ResponseWriter, r *http.Request) {
    userSession := getUserSession(r)
    userSession.Values["id"] = ""
    userSession.Save(r, w)
    http.Redirect(w, r, "/", http.StatusFound)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		newUser := User{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
			Location: r.FormValue("location"),
		}
        newUser.Id = bson.NewObjectId()
        mgoSession = getMgoSession()
		c := mgoSession.DB("").C("users")
		err := c.Insert(&newUser)
		if err != nil {
			panic(err)
		}
        userSession := getUserSession(r)
        userSession.Values["id"] = newUser.Id.Hex()
        userSession.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
        renderTemplate(w, r, "register", nil)
	}
}

func user(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, r, "user", nil)
}
