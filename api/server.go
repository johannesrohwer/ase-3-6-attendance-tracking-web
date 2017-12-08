package api

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// This method is the entry point for app engine.
func init() {
	router := mux.NewRouter()
	router.HandleFunc("/", renderIndex)
	router.HandleFunc("/signup", renderSignUp)
	router.HandleFunc("/dashboard", renderDashboard)
	router.HandleFunc("/createGroup", renderCreateGroup)

	router.HandleFunc("/api/signup", createStudent)
	router.HandleFunc("/api/students", readStudents)
	router.HandleFunc("/api/groups", readAllGroups).Methods("GET")
	router.HandleFunc("/api/groups", createGroup).Methods("POST")
	router.HandleFunc("/api/groups/{number}", readGroup)
	router.HandleFunc("/api/version", readVersion)
	http.Handle("/", router)
}

// Template rendering

func loadTemplate(templateName string) (*template.Template, error) {

	// Assemble template path
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	templatePath := filepath.Join(dir, "/template/"+templateName)
	return template.New(templateName).ParseFiles(templatePath)
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

	// Load groups
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Group").Ancestor(groupKey(c))
	groups := make([]Group, 0, 10)
	if _, err := q.GetAll(c, &groups); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{"groups": groups}
	//data := map[string]string{"name": "world"}

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

func readStudents(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Student").Ancestor(studentKey(c))
	students := make([]Student, 0, 10)
	if _, err := q.GetAll(c, &students); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(students)
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	matrNr := r.Form["matrNr"][0]
	groupNumber := r.Form["groupNumber"][0]
	student := Student{MatrNr: matrNr, GroupNumber: groupNumber}

	context := appengine.NewContext(r)
	key := datastore.NewIncompleteKey(context, "Student", studentKey(context))
	_, err := datastore.Put(context, key, &student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(student)
}

func createGroup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	number := r.Form["number"][0]
	place := r.Form["place"][0]
	time := r.Form["time"][0]
	instructorName := r.Form["instructorName"][0]
	group := Group{Number: number, Place: place, Time: time, InstructorName: instructorName}

	context := appengine.NewContext(r)
	key := datastore.NewIncompleteKey(context, "Group", groupKey(context))
	_, err := datastore.Put(context, key, &group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(group)
}

func readAllGroups(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Group").Ancestor(groupKey(c))
	groups := make([]Group, 0, 10)
	if _, err := q.GetAll(c, &groups); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(groups)
}

func readGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	number := vars["number"]

	context := appengine.NewContext(r)

	var group []Group
	q := datastore.NewQuery("Group").Filter("Number <=", number)
	if _, err := q.GetAll(context, group); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(group)
}

// Database helpers

func studentKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "Student", "default_list", 0, nil)
}

func groupKey(c context.Context) *datastore.Key {
	return datastore.NewKey(c, "Group", "default_list", 0, nil)
}
