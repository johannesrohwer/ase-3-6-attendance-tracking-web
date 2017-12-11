package api

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"time"

	"fmt"

	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// This method is the entry point for app engine.
func init() {
	router := mux.NewRouter()

	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	router.HandleFunc("/", renderIndex)
	router.HandleFunc("/login", renderIndex)
	router.HandleFunc("/signup", renderSignUp)
	router.HandleFunc("/dashboard", renderDashboard)
	router.HandleFunc("/createGroup", renderCreateGroup)
	router.HandleFunc("/test", test)

	router.HandleFunc("/api/signup", createStudent)
	router.HandleFunc("/api/students", readAllStudents).Methods("GET")
	router.HandleFunc("/api/students", createStudent).Methods("POST")
	router.HandleFunc("/api/students/{id}", readStudent)
	router.HandleFunc("/api/groups", readAllGroups).Methods("GET")
	router.HandleFunc("/api/groups", createGroup).Methods("POST")
	router.HandleFunc("/api/groups/{id}", readGroup)
	router.HandleFunc("/api/attendance/{student_id}/{presented}", getAttendanceToken).Methods("GET")
	router.HandleFunc("/api/attendance", registerAttendanceToken).Methods("POST")
	router.HandleFunc("/api/version", readVersion)
	http.Handle("/", router)
}

// Template rendering

func test(w http.ResponseWriter, r *http.Request) {

	t, err := loadTemplate("test.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := map[string]string{"name": "world"}
	t.Execute(w, data)
}

func loadTemplate(templateName string) (*template.Template, error) {

	// Assemble template path
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	templatePath := filepath.Join(dir, "/template/"+templateName)
	return template.New(templateName).Delims("[[", "]]").ParseFiles(templatePath)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(version)
}

// API Student

func createStudent(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var student Student
	if err := decoder.Decode(&student); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := appengine.NewContext(r)
	key := studentKeyFromString(ctx, student.ID)
	_, err := datastore.Put(ctx, key, &student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

func readAllStudents(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	q := datastore.NewQuery("Student")
	var students []Student
	if _, err := q.GetAll(ctx, &students); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if students != nil {
		json.NewEncoder(w).Encode(students)
	} else {
		emptyArrayResponse(w)
	}
}

func readStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := appengine.NewContext(r)
	var student []Student
	q := datastore.NewQuery("Student").Filter("ID =", id)
	if _, err := q.GetAll(ctx, &student); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if student != nil {
		json.NewEncoder(w).Encode(student[0])
	} else {
		emptyObjectResponse(w)
	}
}

// API Group

func createGroup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var group Group
	if err := decoder.Decode(&group); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := appengine.NewContext(r)
	key := groupKeyFromString(ctx, group.ID)

	_, err := datastore.Put(ctx, key, &group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

func readAllGroups(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	q := datastore.NewQuery("Group")

	var groups []Group
	if _, err := q.GetAll(ctx, &groups); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if groups != nil {
		json.NewEncoder(w).Encode(groups)
	} else {
		emptyArrayResponse(w)
	}
}

func readGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := appengine.NewContext(r)
	var group []Group
	q := datastore.NewQuery("Group").Filter("ID =", id)
	if _, err := q.GetAll(ctx, &group); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if group != nil {
		json.NewEncoder(w).Encode(group[0])
	} else {
		emptyObjectResponse(w)
	}
}

// API Attendance

func getAttendanceToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["student_id"]
	presented := vars["presented"] == "true"

	// Fetch Student object
	ctx := appengine.NewContext(r)
	var student []Student
	q := datastore.NewQuery("Student").Filter("ID =", studentID)
	if _, err := q.GetAll(ctx, &student); err != nil || student == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentWeek := "0" // TODO: replace placeholder week
	attendance := NewAttendance(studentID, student[0].GroupID, currentWeek, presented)

	// Encode attendance object to JWT
	JWTObject := map[string]interface{}{}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        time.Date(2019, 10, 10, 12, 0, 0, 0, time.UTC).Unix(), // TODO: insert actual expiry
		"attendance": attendance,
	})

	tokenString, err := token.SignedString([]byte{123})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	JWTObject["token"] = tokenString

	// Send encoded JWT back to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JWTObject)
}

func registerAttendanceToken(w http.ResponseWriter, r *http.Request) {
	// Extract token from POST request
	var request map[string]interface{}
	rawBody, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(rawBody, &request)
	tokenString := request["token"].(string) // TODO: error handler

	// Validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// TODO: read real secret
		return []byte{123}, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// TODO: validate and check for duplicates
		// TODO: insert into datastore
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(claims["attendance"])
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Database helpers

func studentKeyFromString(c context.Context, key string) *datastore.Key {
	return datastore.NewKey(c, "Student", key, 0, nil)
}
func groupKeyFromString(c context.Context, key string) *datastore.Key {
	return datastore.NewKey(c, "Group", key, 0, nil)
}

// Utils

func emptyObjectResponse(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func emptyArrayResponse(w http.ResponseWriter) {
	json.NewEncoder(w).Encode([]interface{}{})
}
