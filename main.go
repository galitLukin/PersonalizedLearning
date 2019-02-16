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
	"github.com/dghubble/oauth1"
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
//var dbUsers = map[string]user{}       // user ID, user -> TODO: should be singular
//var dbSessions = map[string]session{} // session ID, session
var dbSessionsCleaned time.Time
var qd QuestionData
var uid string
var purl string
var an string

const sessionLength int = 30

func init() {
	db, _ = sql.Open("mysql", "arieg419:Nyknicks4191991!@tcp(mydbinstance.cmsj8sgg5big.us-east-2.rds.amazonaws.com:3306)/test02?charset=utf8")
	tpl = template.Must(template.ParseGlob("./templates/*"))
	dbSessionsCleaned = time.Now()
	//uid = "2"
	//purl = "https://nba.com"
	//an = "Reading Test Scores"
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
	purl = req.FormValue("lis_outcome_service_url")

	config := oauth1.Config{
    ConsumerKey:    "oandg_key",
    ConsumerSecret: "oandg_secret",
    //CallbackURL:    "http://localhost/getstarted",
    //Endpoint:       twitter.AuthorizeEndpoint,
	}
	requestToken, requestSecret, err := config.RequestToken()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(requestToken,requestSecret)

	qd.Score = dbInitFetchUser(db, uid, an)


	u := user{
		UserName:       "arieg419@gmail.com",
		Password:       "Beatles",
		First:          "Omer",
		Last:           "Goldberg",
		Uid:            uid,
		AssignmentName: an,
		PostUrl:        purl,
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
		fmt.Println("Quiz is done ...")
		qd.Score.Grade = dbAssignmentDone(db, qd)
		fmt.Println("Users Grade Is: ", qd.Score.Grade)
		return qd.Score.Grade
	}
	return 0.0
}

func quiz(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		if err := req.ParseForm(); err != nil {
			fmt.Println("Failed to parse form...")
			return
		}
		fmt.Println("Continue quiz...")
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
		fmt.Println("Initial question...")
		qd.User.Username = uid
		qd.Question.Assignment = an
		qd.PrevLocation = dbGetUserPrevLocation(db, qd)
		fmt.Println(qd.Question.Assignment)
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
		PostUrl:        purl,
	}

	qpd := QuizPageData{
		UserData:               u,
		QuestionData:           qd,
		PageType:               "quiz",
		HTMLContentText:        template.HTML(qd.Question.Text),
		HTMLContentExplanation: template.HTML(qd.Question.Explanation),
	}

	tpl.ExecuteTemplate(w, "layout", qpd)
	//send post request to edX with the value qd.Score.Grade
}

func logPostBody(req *http.Request) {
	if err := req.ParseForm(); err != nil {
		fmt.Println("Failed to parse form...")
		return
	}

	fmt.Println("Edx - CustomComponentDisplayName: " + req.FormValue("custom_component_display_name"))
	fmt.Println("Edx - LTI Version: " + req.FormValue("lti_version"))
	fmt.Println("Edx - Oauth nonce: " + req.FormValue("oauth_nonce"))
	fmt.Println("Edx - Resource link id: " + req.FormValue("resource_link_id"))
	fmt.Println("Edx - Context ID: " + req.FormValue("context_id"))
	fmt.Println("Edx -  Oauth signature method: " + req.FormValue("oauth_signature_method"))
	fmt.Println("Edx -  Oauth timestamp: " + req.FormValue("oauth_timestamp"))
	fmt.Println("Edx -  Oauth version: " + req.FormValue("oauth_version"))
	fmt.Println("Edx -  Oauth signature: " + req.FormValue("oauth_signature"))
	fmt.Println("Edx - Context Title: " + req.FormValue("context_title"))
	fmt.Println("Edx - LTI Message Type: " + req.FormValue("lti_message_type"))
	fmt.Println("Edx - Launch presentation return url: " + req.FormValue("launch_presentation_return_url"))
	fmt.Println("Edx - Context label: " + req.FormValue("context_label"))
	fmt.Println("Edx - User ID: " + req.FormValue("user_id"))
	fmt.Println("Edx - Roles: " + req.FormValue("roles"))
	fmt.Println("Edx - Custom component due date: " + req.FormValue("custom_component_due_date"))
	fmt.Println("Edx - Oauth consumer key: " + req.FormValue("oauth_consumer_key"))
	fmt.Println("Edx - LIS result sourcedid: " + req.FormValue("lis_result_sourcedid"))
	fmt.Println("Edx - Launch Presentation locale: " + req.FormValue("launch_presentation_locale"))
	fmt.Println("Edx - LISOutcomeService URL: " + req.FormValue("lis_outcome_service_url"))
	fmt.Println("Edx -  Custom component grace period: " + req.FormValue("custom_component_graceperiod"))
	fmt.Println("Edx - Oauth callback: " + req.FormValue("oauth_callback"))
}
