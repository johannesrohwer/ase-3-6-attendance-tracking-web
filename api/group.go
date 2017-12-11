package api

type Group struct {
	ID           string `json:"id"`
	Place        string `json:"place"`
	Time         string `json:"time"`
	InstructorID string `json:"instructor_id"`
}
