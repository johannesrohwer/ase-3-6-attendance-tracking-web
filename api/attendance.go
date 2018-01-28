package api

import "github.com/google/uuid"

type Attendance struct {
	ID        string `json:"id"`
	WeekID    string `json:"week_id"`
	Presented bool   `json:"presented"`
	StudentID string `json:"student_id,omitempty"`
	GroupID   string `json:"group_id,omitempty"`
}

func NewAttendance(studentID string, groupID string, weekID string, presented bool) *Attendance {
	attendanceID := uuid.New().String()
	return &Attendance{ID: attendanceID, StudentID: studentID, GroupID: groupID, WeekID: weekID, Presented: presented}
}

// ByWeek implements the sort interface for []Attendance
type ByWeek []Attendance

func (a ByWeek) Len() int {
	return len(a)
}

func (a ByWeek) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByWeek) Less(i, j int) bool {
	return a[i].WeekID < a[j].WeekID
}
