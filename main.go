package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
	"html/template"
	"net/http"
	"time"
)

type PageData struct {
	UserData user
	PageType string
}

type QuizPageData struct {
	UserData     user
	QuestionData QuestionData
	PageType     string
}

type user struct {
	UserName string
	Password string
	First    string
	Last     string
	Role     string
}

type session struct {
	un           string
	first        string
	last         string
	lastActivity time.Time
}

var db *sql.DB
var err error
var tpl *template.Template
var dbUsers = map[string]user{}       // user ID, user -> TODO: should be singular
var dbSessions = map[string]session{} // session ID, session
var dbSessionsCleaned time.Time
var qd QuestionData

const sessionLength int = 30

func init() {
	db, _ = sql.Open("mysql", "arieg419:Nyknicks4191991!@tcp(mydbinstance.cmsj8sgg5big.us-east-2.rds.amazonaws.com:3306)/test02?charset=utf8")
	tpl = template.Must(template.ParseGlob("./templates/*"))
	dbSessionsCleaned = time.Now()
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/quiz", quiz)
	http.HandleFunc("/getUsers", getUsers)
	http.HandleFunc("/deleteUser", deleteUser)
	http.HandleFunc("/home", home)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":80", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	// u := dbGetUsers(db)
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/home", http.StatusSeeOther)
		return
	}
	// fmt.Fprintln(w, "ALL DB USERS\n", u)
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func quiz(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		if err := req.ParseForm(); err != nil {
			fmt.Println("Failed to parse form...")
			return
		}
		for key, values := range req.PostForm {
			if key == SelectedAnswers {
				qd.QuestionInstance.Answer = values
				qd = getNextQuizState(qd)
			}
		}
	} else {
		qd.QuestionInstance.Answer = nil
		qd = getNextQuizState(qd)
	}

	u := user{
		UserName: "arieg419@gmail.com",
		Password: "Beatles",
		First:    "Omer",
		Last:     "Goldberg",
	}
	qpd := QuizPageData{
		UserData:     u,
		QuestionData: qd,
		PageType:     "quiz",
	}
	tpl.ExecuteTemplate(w, "layout", qpd)
}

func home(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	showSessions()
	pd := PageData{
		UserData: u,
		PageType: "home",
	}
	tpl.ExecuteTemplate(w, "layout", pd)
}

func signup(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	var u user
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")

		u = user{
			UserName: un,
			Password: p,
			First:    f,
			Last:     l,
		}
		dbCreateUser(db, u)
		// create session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		c.MaxAge = sessionLength
		http.SetCookie(w, c)
		dbSessions[c.Value] = session{un: u.UserName, lastActivity: time.Now(), first: u.First, last: u.Last}
		dbUsers[u.UserName] = user{UserName: u.UserName, First: u.First, Last: u.Last, Password: u.Password}
		// redirect
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	pd := PageData{
		UserData: u,
		PageType: "signup",
	}
	tpl.ExecuteTemplate(w, "layout", pd)
}

func login(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	var u user
	// process form submission
	if req.Method == http.MethodPost {
		// fetch user
		email := req.FormValue("email")
		u = dbGetUser(db, email)

		// create session
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		c.MaxAge = sessionLength
		http.SetCookie(w, c)
		dbSessions[c.Value] = session{un: email, lastActivity: time.Now(), first: u.First, last: u.Last}
		dbUsers[email] = user{UserName: email, First: u.First, Last: u.Last}

		// go to home page
		http.Redirect(w, req, "/home", http.StatusSeeOther)
		return
	}
	pd := PageData{
		UserData: user{},
		PageType: "login",
	}
	e := tpl.ExecuteTemplate(w, "layout", pd)
	if e != nil {
		fmt.Println(e)
	}
}

func logout(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	c, _ := req.Cookie("session")

	// delete the session
	delete(dbSessions, c.Value)
	cleanActiveUsers()

	// remove the cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	// clean up dbSessions after certain time passed
	if time.Now().Sub(dbSessionsCleaned) > (time.Second * 30) {
		go cleanSessions()
	}

	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func checkUserProvidedAnswer(w http.ResponseWriter, req *http.Request) {

}

func getUsers(w http.ResponseWriter, req *http.Request) {
	res := dbGetUsers(db)
	fmt.Fprintln(w, res)
}

func deleteUser(w http.ResponseWriter, req *http.Request) {
	res := dbDeleteUser(db)
	fmt.Fprintln(w, res)
}
