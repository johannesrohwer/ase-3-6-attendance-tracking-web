package api

import "golang.org/x/crypto/bcrypt"

type Student struct {
	ID       string `json:"id"`
	GroupID  string `json:"group_id"`
	Name     string `json:"name"`
	password []byte
}

func (s *Student) Password() []byte {
	return s.password
}

func (s *Student) setPassword(pw string) ([]byte, error) {
	// Use bcrypt to (automatically) salt and encrypt password string
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	s.password = bytes
	return bytes, err
}
