package student

// StudentRecord is a struct representing a student record.

type StudentRecord struct {
	ID      string
	Name    string
	Address string
	Phone   string
	Status  string
	School  string
}

// Course represents a course that a student is enrolled in.
type Course struct {
	ID          string
	Name        string
	Instructor  string
	Semester    string
}
