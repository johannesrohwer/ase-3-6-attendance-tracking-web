package api

import (
	"context"

	"errors"

	"google.golang.org/appengine/datastore"
)

// Key helpers

func studentKeyFromString(ctx context.Context, key string) *datastore.Key {
	return datastore.NewKey(ctx, "Student", key, 0, nil)
}

func groupKeyFromString(ctx context.Context, key string) *datastore.Key {
	return datastore.NewKey(ctx, "Group", key, 0, nil)
}

func attendanceKeyFromString(ctx context.Context, key string) *datastore.Key {
	return datastore.NewKey(ctx, "Attendance", key, 0, nil)
}

func instructorKeyFromString(ctx context.Context, key string) *datastore.Key {
	return datastore.NewKey(ctx, "Instructor", key, 0, nil)
}

// Convenience access methods

func getStudent(ctx context.Context, ID string) (*Student, error) {
	var student []Student
	q := datastore.NewQuery("Student").Filter("ID =", ID)
	if _, err := q.GetAll(ctx, &student); err != nil {
		return nil, err
	}

	if len(student) == 0 {
		return nil, errors.New("Not found.")
	}

	return &student[0], nil
}

func getStudents(ctx context.Context) (*[]Student, error) {
	q := datastore.NewQuery("Student")
	var students []Student
	if _, err := q.GetAll(ctx, &students); err != nil {
		return nil, err
	}

	return &students, nil
}

func putStudent(ctx context.Context, student Student) (*Student, error) {
	key := studentKeyFromString(ctx, student.ID)
	_, err := datastore.Put(ctx, key, &student)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func getInstructor(ctx context.Context, ID string) (*Instructor, error) {
	var instructor []Instructor
	q := datastore.NewQuery("Instructor").Filter("ID =", ID)
	if _, err := q.GetAll(ctx, &instructor); err != nil {
		return nil, err
	}

	if len(instructor) == 0 {
		return nil, errors.New("Not found.")
	}

	return &instructor[0], nil
}

func getInstructors(ctx context.Context) (*[]Instructor, error) {
	q := datastore.NewQuery("Instructor")
	var instructors []Instructor
	if _, err := q.GetAll(ctx, &instructors); err != nil {
		return nil, err
	}

	return &instructors, nil
}

func putInstructor(ctx context.Context, instructor Instructor) (*Instructor, error) {
	key := instructorKeyFromString(ctx, instructor.ID)
	_, err := datastore.Put(ctx, key, &instructor)
	if err != nil {
		return nil, err
	}

	return &instructor, nil
}

func getGroup(ctx context.Context, ID string) (*Group, error) {
	var group []Group
	q := datastore.NewQuery("Group").Filter("ID =", ID)
	if _, err := q.GetAll(ctx, &group); err != nil {
		return nil, err
	}

	if len(group) == 0 {
		return nil, errors.New("Not found.")
	}

	return &group[0], nil
}

func getGroups(ctx context.Context) (*[]Group, error) {
	q := datastore.NewQuery("Group")
	var groups []Group
	if _, err := q.GetAll(ctx, &groups); err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return nil, errors.New("No groups found.")
	}

	return &groups, nil
}

func putGroup(ctx context.Context, group Group) (*Group, error) {
	key := groupKeyFromString(ctx, group.ID)
	_, err := datastore.Put(ctx, key, &group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func putAttendance(ctx context.Context, attendance Attendance) (*Attendance, error) {
	// Duplicate IDs do not have to be checked since their values are guaranteed to be unique by the JWT signature.

	// Check if an attendance record already exists for that week
	if pastAttendances, err := getAttendancesForStudent(ctx, attendance.StudentID); err == nil {
		for _, a := range *pastAttendances {
			if a.WeekID == attendance.WeekID {
				return &a, nil
			}
		}

	}
	// FIXME dirty fix
	// else {
	//	return nil, nil
	// }

	key := attendanceKeyFromString(ctx, attendance.ID)
	_, err := datastore.Put(ctx, key, &attendance)
	if err != nil {
		return nil, err
	}

	return &attendance, nil
}

func getAttendancesForStudent(ctx context.Context, studentID string) (*[]Attendance, error) {
	var attendances []Attendance
	q := datastore.NewQuery("Attendance").Filter("StudentID =", studentID)
	if _, err := q.GetAll(ctx, &attendances); err != nil {
		return nil, err
	}

	if len(attendances) == 0 {
		return nil, errors.New("Not found.")
	}

	return &attendances, nil
}

func getAttendances(ctx context.Context) (*[]Attendance, error) {
	q := datastore.NewQuery("Attendance")
	var attendance []Attendance
	if _, err := q.GetAll(ctx, &attendance); err != nil {
		return nil, err
	}

	return &attendance, nil
}
