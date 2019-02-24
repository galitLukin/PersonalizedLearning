package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

type PageData struct {
	UserData user
	PageType string
}

type QuizPageData struct {
	UserData               user
	QuestionData           QuestionData
	PageType               string
	HTMLContentText        template.HTML
	HTMLContentExplanation template.HTML
}

type user struct {
	UserName       string
	Password       string
	First          string
	Last           string
	Role           string
	Uid            string
	PostUrl        string
	AssignmentName string
}

type session struct {
	un           string
	first        string
	last         string
	lastActivity time.Time
}

type EdxPOSTBody struct {
	CustomerComponentDisplayName int64  `json:"custom_component_display_name"`
	LTIVersion                   string `json:"lti_version"`
	OauthNonce                   string `json:"oauth_nonce"`
	ResourceLinkId               string `json:"resource_link_id"`
	ContextId                    string `json:"context_id"`
	OauthSignatureMethod         string `json:"oauth_signature_method"`
	OauthTimestamp               string `json:"oauth_timestamp"`
	OauthVersion                 string `json:"oauth_version"`
	OauthSignature               string `json:"oauth_signature"`
	ContextTitle                 string `json:"context_title"`
	LTIMessageType               string `json:"lti_message_type"`
	UserID                       string `json:"user_id"`
	OauthConsumerKey             string `json:"oauth_consumer_key"`
	LISOutcomeServiceURL         string `json:"lis_outcome_service_url"`
}

var db *sql.DB
var err error
var tpl *template.Template

var dbSessionsCleaned time.Time
var qd QuestionData
var uid string
var an string

const sessionLength int = 30

func init() {
	db, _ = sql.Open("mysql", "arieg419:Nyknicks4191991!@tcp(mydbinstance.cmsj8sgg5big.us-east-2.rds.amazonaws.com:3306)/test02?charset=utf8")
	tpl = template.Must(template.ParseGlob("./templates/*"))
	dbSessionsCleaned = time.Now()
	//uid = "1"
	//an = "Detecting Flu Epedemics"
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/getstarted", getStarted)
	http.HandleFunc("/quiz", quiz)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/instance", instance)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":80", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/getstarted", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "OK")
}

func instance(w http.ResponseWriter, req *http.Request) {
	fmt.Println("instace route hit")
	resp, err := http.Get("http://3.16.157.40/latest/meta-data/instance-id")
	if err != nil {
		io.WriteString(w, "Couldn't fetch meta data "+err.Error())
		fmt.Println(err)
		return
	}

	bs := make([]byte, resp.ContentLength)
	resp.Body.Read(bs)
	resp.Body.Close()
	io.WriteString(w, string(bs))
	fmt.Println("instance end" + string(bs))
}

func getStarted(w http.ResponseWriter, req *http.Request) {
	d, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(d))

	logPostBody(req)
	uid = req.FormValue("user_id")
	an = req.FormValue("custom_component_display_name")

	qd.Score = dbInitFetchUser(db, uid, an)

	u := user{
		UserName:       "arieg419@gmail.com",
		Password:       "Beatles",
		First:          "Omer",
		Last:           "Goldberg",
		Uid:            uid,
		AssignmentName: an,
	}
	qpd := QuizPageData{
		UserData:     u,
		QuestionData: qd,
		PageType:     "getstarted",
	}
	tpl.ExecuteTemplate(w, "layout", qpd)
}

func finishAssignment(db *sql.DB, qd QuestionData) float32 {
	if qd.QuestionInstance.Status == "Done" {
		fmt.Println("Quiz is done ...",uid,an)
		qd.Score.Grade = dbAssignmentDone(db, qd)
		fmt.Println("Users Grade Is: ", qd.Score.Grade,uid,an)

		return float32(int(qd.Score.Grade * 100))
	}
	return 0.0
}

func quiz(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		if err := req.ParseForm(); err != nil {
			fmt.Println("Failed to parse form...",uid,an)
			return
		}
		fmt.Println("Continue quiz...",uid,an)
		for key, values := range req.PostForm {
			if key == SelectedAnswers {
				qd.QuestionInstance.Answer = values
				qd.PrevLocation.IsFirst = false
				dbInsertResponse(db, qd)
				if qd.QuestionInstance.Status == "Correct" || qd.QuestionInstance.Status == "IncorrectNoAttempts" {
					qd.Score = dbFetchUserInScores(db, qd)
				}
				qd = getNextQuizState(qd)
				dbUpdateFinishedQuestion(db, qd)
				qd.Score.Grade = finishAssignment(db, qd)
			}
		}
	} else {
		fmt.Println("Initial question...",uid,an)
		qd.User.Username = uid
		qd.Question.Assignment = an
		qd.PrevLocation = dbGetUserPrevLocation(db, qd)
		qd = getNextQuizState(qd)
		qd.Score.Grade = finishAssignment(db, qd)
	}

	u := user{
		UserName:       "arieg419@gmail.com",
		Password:       "Beatles",
		First:          "Omer",
		Last:           "Goldberg",
		Uid:            uid,
		AssignmentName: an,
	}

	qpd := QuizPageData{
		UserData:               u,
		QuestionData:           qd,
		PageType:               "quiz",
		HTMLContentText:        template.HTML(qd.Question.Text),
		HTMLContentExplanation: template.HTML(qd.Question.Explanation),
	}

	tpl.ExecuteTemplate(w, "layout", qpd)
}

func logPostBody(req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Println("Failed to parse form...")
		return
	}
	fmt.Println("Edx - CustomComponentDisplayName: " + req.FormValue("custom_component_display_name"))
	fmt.Println("Edx - User ID: " + req.FormValue("user_id"))
	fmt.Println("Edx - Roles: " + req.FormValue("roles"))
}
