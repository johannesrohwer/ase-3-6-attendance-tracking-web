package api

import (
	"encoding/json"
	"net/http"
	"time"

	"io/ioutil"

	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

// This method is the entry point for app engine.
func init() {
	router := mux.NewRouter()

	// Load /static folder
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	// Register routes for template engine
	router.HandleFunc("/", renderIndex)
	router.HandleFunc("/login", renderIndex)
	router.HandleFunc("/signup", renderSignUp)
	router.HandleFunc("/dashboard", renderDashboard)
	router.HandleFunc("/createGroup", renderCreateGroup)

	// Register json API routes
	router.HandleFunc("/api/signup", createStudent)
	router.HandleFunc("/api/students", readAllStudents).Methods("GET")
	router.HandleFunc("/api/students", createStudent).Methods("POST")
	router.HandleFunc("/api/students/{id}", readStudent)

	router.HandleFunc("/api/groups", readAllGroups).Methods("GET")
	router.HandleFunc("/api/groups", createGroup).Methods("POST")
	router.HandleFunc("/api/groups/{id}", readGroup)

	router.HandleFunc("/api/attendances", readAllAttendances).Methods("GET")
	router.HandleFunc("/api/attendances/new/{student_id}/{presented}", getAttendanceToken).Methods("GET")
	router.HandleFunc("/api/attendances/register", registerAttendanceToken).Methods("POST")
	router.HandleFunc("/api/attendances/for/{student_id}", readAttendancesForStudent).Methods("GET")

	router.HandleFunc("/api/version", readVersion)

	http.Handle("/", router)
}

func renderIndex(w http.ResponseWriter, r *http.Request) {

	t, err := loadTemplate("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := map[string]string{"name": "world"}
	t.Execute(w, data)
}

func renderCreateGroup(w http.ResponseWriter, r *http.Request) {

	t, err := loadTemplate("createGroup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := map[string]string{"name": "world"}
	t.Execute(w, data)
}

func renderSignUp(w http.ResponseWriter, r *http.Request) {

	t, err := loadTemplate("signup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := map[string]string{}

	t.Execute(w, data)
}

func renderDashboard(w http.ResponseWriter, r *http.Request) {

	t, err := loadTemplate("dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := map[string]string{"name": "world"}
	t.Execute(w, data)
}

// API handlers

func readVersion(w http.ResponseWriter, r *http.Request) {
	authors := []string{"Sandra Grujovic", "Paul Schmiedermayer", "Johannes Rohwer"}
	version := NewVersion("v0.1", authors)

	sendResponse(w, version, http.StatusOK)
}

// API Student

func createStudent(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	decoder := json.NewDecoder(r.Body)
	var student Student
	if err := decoder.Decode(&student); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err := putStudent(ctx, student)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	sendResponse(w, student, http.StatusCreated)
}

func readAllStudents(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	students, err := getStudents(ctx)
	if err != nil {
		sendResponse(w, emptyArray(), http.StatusOK)
		return
	}

	sendResponse(w, students, http.StatusOK)
}

func readStudent(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	id := vars["id"]
	student, err := getStudent(ctx, id)
	if err != nil {
		sendResponse(w, emptyObject(), http.StatusOK)
		return
	}

	sendResponse(w, student, http.StatusOK)
}

// API Group

func createGroup(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	decoder := json.NewDecoder(r.Body)
	var group Group
	if err := decoder.Decode(&group); err != nil {
		sendErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	if _, err := putGroup(ctx, group); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	sendResponse(w, group, http.StatusCreated)
}

func readAllGroups(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	groups, err := getGroups(ctx)
	if err != nil {
		sendResponse(w, emptyArray(), http.StatusOK)
		return
	}

	sendResponse(w, groups, http.StatusOK)
}

func readGroup(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	id := vars["id"]
	group, err := getGroup(ctx, id)
	if err != nil {
		sendResponse(w, emptyObject(), http.StatusNotFound)
		return
	}

	sendResponse(w, group, http.StatusOK)
}

// API Attendance

func getAttendanceToken(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	studentID := vars["student_id"]
	presented := vars["presented"] == "true"
	currentWeek := "0" // TODO: replace placeholder week

	student, err := getStudent(ctx, studentID)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	attendance := NewAttendance(studentID, student.GroupID, currentWeek, presented)

	JWTObject := map[string]interface{}{}

	// Prepare claims for JQT
	expiryTime := time.Now().Add(24 * time.Hour) // One day expiration time
	claims := jwt.MapClaims{
		"exp":        expiryTime,
		"attendance": attendance,
	}

	JWTObject["token"], err = createJWTToken(claims)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	sendResponse(w, JWTObject, http.StatusCreated)
}

func registerAttendanceToken(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// Extract token from POST request
	var request map[string]interface{}
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Parse request into map
	json.Unmarshal(rawBody, &request)
	tokenString, ok := request["token"].(string)
	if !ok {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Validate and process token
	if claims, err := validateJWTToken(tokenString); err == nil {
		attendanceToken := claims["attendance"].(map[string]interface{})
		tID := attendanceToken["id"].(string)
		tPresented := attendanceToken["presented"].(bool)
		tWeekID := attendanceToken["week_id"].(string)
		tGroupID := attendanceToken["group_id"].(string)
		tStudentID := attendanceToken["student_id"].(string)
		if err != nil {
			sendErrorResponse(w, err, http.StatusInternalServerError)
		}

		attendance := Attendance{ID: tID, WeekID: tWeekID, GroupID: tGroupID, StudentID: tStudentID, Presented: tPresented}
		if _, err := putAttendance(ctx, attendance); err != nil {
			sendErrorResponse(w, err, http.StatusInternalServerError)
		}

		sendResponse(w, attendance, http.StatusCreated)
	} else {
		sendErrorResponse(w, errors.New("Invalid token."), http.StatusInternalServerError)
	}
}

func readAllAttendances(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	attendances, err := getAttendances(ctx)
	if err != nil {
		sendResponse(w, emptyArray(), http.StatusOK)
		return
	}

	sendResponse(w, attendances, http.StatusOK)
}

func readAttendancesForStudent(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	studentID := vars["student_id"]
	attendances, err := getAttendancesForStudent(ctx, studentID)
	if err != nil {
		sendResponse(w, emptyArray(), http.StatusNotFound)
		return
	}

	sendResponse(w, attendances, http.StatusOK)
}
