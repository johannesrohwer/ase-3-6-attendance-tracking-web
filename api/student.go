package api

type Student struct {
	ID       string `json:"id"`
	GroupID  string `json:"group_id"`
	Name     string `json:"name"`
	Password []byte `json:"-"`
}

func NewStudent(ID string, name string, groupID string, password string) (*Student, error) {
	student := &Student{ID: ID, GroupID: groupID, Name: name}
	_, err := student.setPassword(password)
	if err != nil {
		return nil, err
	}

	return student, nil
}

func (s *Student) setPassword(pw string) ([]byte, error) {
	passwordHash, err := generatePasswordHash(pw)
	s.Password = passwordHash
	return passwordHash, err
}
