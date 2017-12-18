package api

type Instructor struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password []byte `json:"-"`
}

func NewInstructor(ID string, name string, password string) (*Instructor, error) {
	instructor := &Instructor{ID: ID, Name: name}
	_, err := instructor.setPassword(password)
	if err != nil {
		return nil, err
	}

	return instructor, nil
}

func (ins *Instructor) setPassword(pw string) ([]byte, error) {
	passwordHash, err := generatePasswordHash(pw)
	ins.Password = passwordHash
	return passwordHash, err
}
