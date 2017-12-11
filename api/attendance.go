package api

import "github.com/google/uuid"

type Attendance struct {
	ID        string `json:"id"`
	WeekID    string `json:"week_id"`
	Presented bool   `json:"presented"`
	StudentID string `json:"student_id"`
	GroupID   string `json:"group_id"`
}

func NewAttendance(studentID string, groupID string, weekID string, presented bool) *Attendance {
	attendanceID := uuid.New().String()
	return &Attendance{ID: attendanceID, StudentID: studentID, GroupID: groupID, WeekID: weekID, Presented: presented}
}
