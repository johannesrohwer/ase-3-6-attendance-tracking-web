package api

import (
	"encoding/json"
	"net/http"
	"time"

	"io/ioutil"

	"errors"

	"strconv"

	"sort"

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
	router.HandleFunc("/api/login", createCredentials)

	router.HandleFunc("/api/signup", createStudent)
	router.HandleFunc("/api/students", createStudent).Methods("POST")
	router.Handle("/api/students", authMiddleware(readAllStudents, "student", "instructor")).Methods("GET")
	router.Handle("/api/students/{id}", authMiddleware(readStudent, "student", "instructor"))

	router.HandleFunc("/api/instructors", createInstructor).Methods("POST") // TODO: protect this route
	router.Handle("/api/instructors", authMiddleware(readAllInstructors, "student", "instructor")).Methods("GET")
	router.Handle("/api/instructors/{id}", authMiddleware(readInstructor, "student", "instructor"))

	router.HandleFunc("/api/groups", createGroup).Methods("POST")
	router.HandleFunc("/api/groups", readAllGroups).Methods("GET")
	router.Handle("/api/groups/{id}", authMiddleware(readGroup, "student", "instructor"))

	router.Handle("/api/attendances", authMiddleware(readAllAttendances, "student", "instructor")).Methods("GET")
	router.Handle("/api/attendances/new/{student_id}/{presented}", authMiddleware(getAttendanceToken, "student")).Methods("GET")
	router.Handle("/api/attendances/register", authMiddleware(registerAttendanceToken, "instructor")).Methods("POST")
	router.Handle("/api/attendances/for/{student_id}", authMiddleware(readAttendancesForStudent, "student", "instructor")).Methods("GET")

	router.HandleFunc("/api/week/current", readCurrentWeek)
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

func readCurrentWeek(w http.ResponseWriter, r *http.Request) {
	currentWeekObject := map[string]interface{}{
		"current_week": getCurrentWeek(),
		"base_week":    getBaseWeek()}

	sendResponse(w, currentWeekObject, http.StatusOK)
}

// API Login

func createCredentials(w http.ResponseWriter, r *http.Request) {
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
	ID, ok := request["id"].(string)
	if !ok {
		sendErrorResponse(w, errors.New("id empty."), http.StatusBadRequest)
		return
	}

	password, ok := request["password"].(string)
	if !ok {
		sendErrorResponse(w, errors.New("password empty."), http.StatusBadRequest)
		return
	}

	// Load user from datastore
	// Try student first
	if student, err := getStudent(ctx, ID); err == nil {
		if verifyPassword(password, student.Password) {
			permissions := []string{"student"}
			credentials := NewCredentials(ID, permissions)
			expiryTime := time.Now().Add(3 * time.Hour)
			token, err := createJWTToken(jwt.MapClaims{"credentials": credentials, "exp": expiryTime})
			if err != nil {
				sendErrorResponse(w, errors.New("JWT creation failed."), http.StatusInternalServerError)
			}

			response := map[string]interface{}{"token": token}
			sendResponse(w, response, http.StatusOK)
			return
		}

		sendErrorResponse(w, errors.New("Invalid credentials."), http.StatusForbidden)
		return
	}

	// If student failed, try to load an instructor
	if instructor, err := getInstructor(ctx, ID); err == nil {
		if verifyPassword(password, instructor.Password) {
			permissions := []string{"instructor"}
			credentials := NewCredentials(ID, permissions)
			expiryTime := time.Now().Add(3 * time.Hour)
			token, err := createJWTToken(jwt.MapClaims{"credentials": credentials, "exp": expiryTime})
			if err != nil {
				sendErrorResponse(w, errors.New("JWT creation failed."), http.StatusInternalServerError)
			}

			response := map[string]interface{}{"token": token}
			sendResponse(w, response, http.StatusOK)
			return
		}

		sendErrorResponse(w, errors.New("Invalid credentials."), http.StatusForbidden)
		return
	}

	// Neither student nor instructor login was successful, throw an error.
	sendErrorResponse(w, errors.New("Invalid credentials."), http.StatusForbidden)
}

// API Student

func createStudent(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// Extract password
	var request map[string]interface{}
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Parse request into map
	json.Unmarshal(rawBody, &request)
	ID, ok := request["id"].(string)
	if !ok {
		sendErrorResponse(w, errors.New("id empty."), http.StatusBadRequest)
		return
	}

	name, ok := request["name"].(string)
	if !ok {
		sendErrorResponse(w, errors.New("name empty."), http.StatusBadRequest)
		return
	}

	groupID, ok := request["group_id"].(string)
	if !ok {
		sendErrorResponse(w, errors.New("group_id empty."), http.StatusBadRequest)
		return
	}

	password, ok := request["password"].(string)
	if !ok {
		sendErrorResponse(w, errors.New("password empty."), http.StatusBadRequest)
		return
	}

	student, err := NewStudent(ID, name, groupID, password)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if _, err := putStudent(ctx, *student); err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Create autorization token
	permissions := []string{"student"}
	credentials := NewCredentials(ID, permissions)
	expiryTime := time.Now().Add(3 * time.Hour)
	token, err := createJWTToken(jwt.MapClaims{"credentials": credentials, "exp": expiryTime})
	if err != nil {
		sendErrorResponse(w, errors.New("JWT creation failed."), http.StatusInternalServerError)
	}

	response := struct {
		Student
		Token string `json:"token"`
	}{Student: *student, Token: *token}
	sendResponse(w, response, http.StatusCreated)
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

// API Instructor

func createInstructor(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// Extract password
	var request map[string]interface{}
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Parse request into map
	json.Unmarshal(rawBody, &request)
	ID, ok := request["id"].(string)
	if !ok {
		sendErrorResponse(w, errors.New("id empty."), http.StatusBadRequest)
		return
	}

	name, ok := request["name"].(string)
	if !ok {
		sendErrorResponse(w, errors.New("name empty."), http.StatusBadRequest)
		return
	}

	password, ok := request["password"].(string)
	if !ok {
		sendErrorResponse(w, errors.New("password empty."), http.StatusBadRequest)
		return
	}

	instructor, err := NewInstructor(ID, name, password)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if _, err := putInstructor(ctx, *instructor); err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Create authorization token
	permissions := []string{"instructor"}
	credentials := NewCredentials(ID, permissions)
	expiryTime := time.Now().Add(3 * time.Hour)
	token, err := createJWTToken(jwt.MapClaims{"credentials": credentials, "exp": expiryTime})
	if err != nil {
		sendErrorResponse(w, errors.New("JWT creation failed."), http.StatusInternalServerError)
	}

	response := struct {
		Instructor
		Token string `json:"token"`
	}{Instructor: *instructor, Token: *token}
	sendResponse(w, response, http.StatusCreated)
}

func readAllInstructors(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	instructors, err := getInstructors(ctx)
	if err != nil {
		sendResponse(w, emptyArray(), http.StatusOK)
		return
	}

	sendResponse(w, instructors, http.StatusOK)
}

func readInstructor(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	id := vars["id"]
	instructor, err := getInstructor(ctx, id)
	if err != nil {
		sendResponse(w, emptyObject(), http.StatusOK)
		return
	}

	sendResponse(w, instructor, http.StatusOK)
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
	currentWeek := strconv.Itoa(getCurrentWeek())

	student, err := getStudent(ctx, studentID)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	attendance := NewAttendance(studentID, student.GroupID, currentWeek, presented)

	JWTObject := map[string]interface{}{}

	// Prepare claims for JWT
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

		newAttendance := Attendance{ID: tID, WeekID: tWeekID, GroupID: tGroupID, StudentID: tStudentID, Presented: tPresented}
		if existingAttendance, err := putAttendance(ctx, newAttendance); err != nil {
			sendErrorResponse(w, err, http.StatusInternalServerError)
		} else if existingAttendance != nil {
			sendResponse(w, existingAttendance, http.StatusOK)
		} else {
			sendResponse(w, newAttendance, http.StatusCreated)
		}

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
		sendResponse(w, emptyArray(), http.StatusInternalServerError)
		return
	}

	// Check if empty attendances (absences) should be added to the response
	showMissingAttendances := r.URL.Query().Get("missing_attendances")
	if showMissingAttendances == "true" {
		// Create list of attendances that should be there
		currentWeek := getCurrentWeek()
		missingWeekIDs := map[string]string{}
		for i := 1; i <= currentWeek; i++ {
			currentWeekID := strconv.Itoa(i)
			missingWeekIDs[currentWeekID] = currentWeekID
		}

		// Tick off attendances that are actually there
		for _, attendance := range *attendances {
			currentWeekID := attendance.WeekID
			delete(missingWeekIDs, currentWeekID)
		}

		// Add empty attendances for the missing ones
		paddedAttendances := *attendances
		for _, missingWeekID := range missingWeekIDs {
			emptyAttendance := Attendance{ID: "-1", WeekID: missingWeekID}
			paddedAttendances = append(paddedAttendances, emptyAttendance)
		}

		// Make sure that []attendances are sorted
		sort.Sort(ByWeek(paddedAttendances))
		attendances = &paddedAttendances
	}

	if len(*attendances) == 0 {
		sendResponse(w, emptyArray(), http.StatusOK)
	}

	sendResponse(w, attendances, http.StatusOK)
}
