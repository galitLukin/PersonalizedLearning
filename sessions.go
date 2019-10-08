package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

func getUserAsmt(w http.ResponseWriter, req *http.Request) QuestionData {
	// get cookie
	c, err := req.Cookie("session")
	if err != nil {
		sID, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}

	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)

	// if the user exists already, get user
	var qd QuestionData
	if s, ok := dbSessions[c.Value]; ok {
		s.lastActivity = time.Now().UTC()
		dbSessions[c.Value] = s
		qd = dbUserState[s.un]
	}
	return qd
}

func alreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false // user is not logged in
	}
	s, ok := dbSessions[c.Value]
	if ok {
		s.lastActivity = time.Now().UTC()
		dbSessions[c.Value] = s
	}
	_, ok = dbUserState[s.un]
	// refresh session
	c.MaxAge = sessionLength
	http.SetCookie(w, c)
	cleanSessions()
	showSessions()
	return ok
}

func getOldState(w http.ResponseWriter, req *http.Request) string {
	c, err := req.Cookie("session")
	if err != nil {
		return ""
	}
	s, ok := dbSessions[c.Value]
	if ok {
		s.lastActivity = time.Now().UTC()
		dbSessions[c.Value] = s
	}
	qd, ok := dbUserState[s.un]
	return qd.User.AssignmentName
}

func cleanActiveUsers() {
	for k := range dbUserState {
		delete(dbUserState, k)
	}
	dbSessionsCleaned = time.Now().UTC()
}

func cleanSessions() {
	for k, v := range dbSessions {
		if time.Now().UTC().Sub(v.lastActivity) > (time.Second * 3600) {
			delete(dbSessions, k)
		}
	}
	dbSessionsCleaned = time.Now().UTC()
}

// for demonstration purposes
func showSessions() {
	fmt.Println("********")
	for k, v := range dbSessions {
		fmt.Println(k, v.un)
	}
	fmt.Println("********")
}

// // for demonstration purposes
// func showDBUsers() {
// 	fmt.Println("********")
// 	for k, v := range dbUsers {
// 		fmt.Println(k, v.UserName)
// 	}
// 	fmt.Println("")
// }
