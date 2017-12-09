package api

type Attendance struct {
	ID         string `json:"id"`
	Week_ID    string `json:"week_id"`
	Presented  bool   `json:"presented"`
	Student_ID string `json:"student_id"`
	Group_ID   string `json:"group_id"`
}
