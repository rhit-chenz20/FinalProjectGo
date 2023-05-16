package main

import (
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

// /key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")

type student struct {
	ID      string
	Name    string
	Year    string
	Courses map[string]float64
}

// starter data for courses
var courses = []string{
	"CSSE403",
	"CSSE374",
	"CSSE304",
}

// starter data for students
var students = []student{
	{ID: "1", Name: "John Doe", Year: "Junior"},
	{ID: "2", Name: "Jane Doe", Year: "Sophomore"},
	{ID: "3", Name: "Dave Smith", Year: "Senior"},
}

//Encrypted courses

//Encryted

var id_counter = 4

func getCourses(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, courses)
}

func postCourses(c *gin.Context) {
	var newCourse string
	err := c.BindJSON(&newCourse)
	if err != nil {
		return
	}
	courses = append(courses, newCourse)
	c.IndentedJSON(http.StatusCreated, newCourse)
}

func getStudents(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, students)
}

func postStudents(c *gin.Context) {
	var newStudent student

	// Call BindJSON to bind the received JSON to
	// newStudent.
	if err := c.BindJSON(&newStudent); err != nil {
		return
	}
	newStudent.ID = strconv.Itoa(id_counter)
	id_counter += 1
	// Add the new student to the slice.
	students = append(students, newStudent)
	c.IndentedJSON(http.StatusCreated, newStudent)
}

func getStudentByID(c *gin.Context) {
	id := c.Param("id")

	for _, s := range students {
		if s.ID == id {
			c.IndentedJSON(http.StatusOK, s)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student not found"})
}

// Overwrite a Student's grade in a specified course
func postGradeToStudentbyID(c *gin.Context) {
	id := c.Param("id")
	course := c.Param("coursename")
	grade := c.Param("grade")

	for _, s := range students {
		if s.ID == id {
			gradeFloat, err := strconv.ParseFloat(grade, 32)
			if err == nil {
				s.Courses[course] = gradeFloat
				c.IndentedJSON(http.StatusOK, gradeFloat)
			} else {
				c.IndentedJSON(http.StatusNotFound, gin.H{"message": "grade not found"})
			}
		}
	}
}

// Get a Student's grade in a specified course
func getStudentsGradeById(c *gin.Context) {
	id := c.Param("id")
	course := c.Param("coursename")

	for _, s := range students {
		if s.ID == id {
			for name, grade := range s.Courses {
				if name == course {
					c.IndentedJSON(http.StatusOK, grade)
				}
			}
		}
	}
}

func main() {
	router := gin.Default()
	router.GET("/students", getStudents)
	router.POST("/students", postStudents)

	router.GET("/courses", getCourses)
	router.POST("/courses", postCourses)

	router.GET("/student/:id", getStudentByID)

	router.POST("/student/:id/course/:coursename/grade/:grade", postGradeToStudentbyID)
	router.GET("/student/:id/course/:coursename", getStudentsGradeById)

	router.Run("localhost:8080")
}
