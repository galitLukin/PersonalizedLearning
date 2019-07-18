package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
	"html/template"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

type QuizDisplay struct {
	HTMLContentText        template.HTML
	HTMLContentExplanation template.HTML
	Options                []string
	Correctness            bool
	Answer                 string
	CorrectAnswer          []string
	AttemptsOverall        int
	Weight                 int
	QUIndex                int
}

type user struct {
	Uid            string
	AssignmentName string
}

type session struct {
	un           string
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

var dbSessions = map[string]session{}
var dbUserState = map[string]QuestionData{}
var dbSessionsCleaned time.Time

const sessionLength int = 1800

var uid string
var an string

func init() {
	db, _ = sql.Open("mysql", "arieg419:Nyknicks4191991!@tcp(mydbinstance.cmsj8sgg5big.us-east-2.rds.amazonaws.com:3306)/test02?charset=utf8")
	tpl = template.Must(template.ParseGlob("./templates/*"))
	dbSessionsCleaned = time.Now()
	uid = "2"
	an = "Climate Change"
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/getstarted", getStarted)
	http.HandleFunc("/quiz", quiz)
	http.HandleFunc("/pastQuestions", pastQuestions)
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

	//logPostBody(req)
	//uid := req.FormValue("user_id")
	//an := req.FormValue("custom_component_display_name")

	var qd QuestionData
	user_assignment := an + "+" + uid
	qd.User = user{
		Uid:            uid,
		AssignmentName: an,
	}
	qd.Score = dbInitFetchUser(db, uid, an)

	if !alreadyLoggedIn(w, req) {
		// create session
		fmt.Println("Creating session")
		sID, _ := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		c.MaxAge = sessionLength
		http.SetCookie(w, c)
		dbSessions[c.Value] = session{un: user_assignment, lastActivity: time.Now()}
		dbUserState[user_assignment] = qd
	} else {
		if getOldState(w, req) != an {
			fmt.Println("switiching assignments")
			c, err := req.Cookie("session")
			if err != nil {
				return
			}
			dbSessions[c.Value] = session{un: user_assignment, lastActivity: time.Now()}
			dbUserState[user_assignment] = qd
		}
	}

	qpd := struct {
		QuestionData QuestionData
		PageType     string
	}{
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

		return float32(int(qd.Score.Grade * 100))
	}
	return 0.0
}

func quiz(w http.ResponseWriter, req *http.Request) {
	myqd := getUserAsmt(w, req)
	user_assignment := myqd.User.AssignmentName + "+" + myqd.User.Uid
	fmt.Println(user_assignment)
	var newqd QuestionData
	if req.Method == http.MethodPost {
		if err := req.ParseForm(); err != nil {
			fmt.Println("Failed to parse form...", user_assignment)
			return
		}
		fmt.Println("Continue quiz...", user_assignment)
		for key, values := range req.PostForm {
			if key == SelectedAnswers {
				if len(values[0]) > 0 && !valueinOptions(strings.Split(values[0], ",")[0], myqd.Question.Options) {
					//user did back and got to
					fmt.Println("User went back...", user_assignment)
					break
				}
				myqd.QuestionInstance.Answer = values
				myqd.PrevLocation.IsFirst = false
				dbInsertResponse(db, myqd)
				if myqd.QuestionInstance.Status == "Correct" || myqd.QuestionInstance.Status == "IncorrectNoAttempts" {
					myqd.Score = dbFetchUserInScores(db, myqd)
				}
				newqd = getNextQuizState(myqd)
				dbUpdateFinishedQuestion(db, newqd)
				newqd.Score.Grade = finishAssignment(db, newqd)
				dbUserState[user_assignment] = newqd
			}
		}
	} else {
		fmt.Println("Initial question...", user_assignment)
		myqd.PrevLocation = dbGetUserPrevLocation(db, myqd)
		newqd = getNextQuizState(myqd)
		newqd.Score.Grade = finishAssignment(db, newqd)
		dbUserState[user_assignment] = newqd
	}

	qpd := struct {
		QuestionData           QuestionData
		PageType               string
		HTMLContentText        template.HTML
		HTMLContentExplanation template.HTML
	}{
		QuestionData:           dbUserState[user_assignment],
		PageType:               "quiz",
		HTMLContentText:        template.HTML(dbUserState[user_assignment].Question.Text),
		HTMLContentExplanation: template.HTML(dbUserState[user_assignment].Question.Explanation),
	}

	tpl.ExecuteTemplate(w, "layout", qpd)
}

func pastQuestions(w http.ResponseWriter, req *http.Request) {
	myqd := getUserAsmt(w, req)
	questions := dbFetchUserInResponses(db, myqd)
	var quizpd []QuizDisplay
	quizpd = nil
	if questions != nil {
		pastQs := getAllPastQuestions(myqd, questions)
		var qd QuizDisplay
		var q Question
		for i, pq := range pastQs {
			q = pq.Question
			if pq.Correctness || q.AttemptsOverall == pq.Attempts {
				qd = QuizDisplay{
					HTMLContentText:        template.HTML(q.Text),
					HTMLContentExplanation: template.HTML(q.Explanation),
					Options:                q.Options,
					Correctness:            pq.Correctness,
					Answer:                 pq.Answer,
					CorrectAnswer:          q.CorrectAnswer,
					AttemptsOverall:        q.AttemptsOverall,
					Weight:                 q.Weight,
					QUIndex:                i + 1,
				}
				quizpd = append(quizpd, qd)
			}
		}
	}

	u := user{
		Uid:            uid,
		AssignmentName: an,
	}

	qpd := struct {
		UserData  user
		Questions []QuizDisplay
		PageType  string
	}{
		UserData:  u,
		Questions: quizpd,
		PageType:  "pastQuestions",
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

func valueinOptions(a string, list []string) bool {
	for _, b := range list {
		if strings.Split(b, ",")[0] == a {
			return true
		}
	}
	return false
}
