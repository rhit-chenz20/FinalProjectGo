package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	requestType int
	context     *gin.Context
}

type student struct {
	ID      string
	Name    string
	Year    string
	Courses map[string]float64
}

var id_counter = 4

var students = []student{
	{ID: "1", Name: "John Doe", Year: "Junior"},
	{ID: "2", Name: "Jane Doe", Year: "Sophomore"},
	{ID: "3", Name: "Dave Smith", Year: "Senior"},
}

// starter data for courses
var courses = []string{
	"CSSE403",
	"CSSE374",
	"CSSE304",
}

const CHANNEL_SIZE = 250
const NUM_WORKERS = 5

var cin = make(chan RequestBody, CHANNEL_SIZE)

// generatlized handler
func addRequestToQueue(requestType int, c *gin.Context) {
	request := RequestBody{
		requestType: requestType,
		context:     c,
	}
	cin <- request
}

// GET REQUESTS
// ---------------------------------------------------------------
func getStudentsRequest(c *gin.Context) { addRequestToQueue(1, c) }

func getCoursesRequest(c *gin.Context) { addRequestToQueue(2, c) }

func getStudentByIDRequest(c *gin.Context) { addRequestToQueue(3, c) }

func getStudentCourseGradeByIdRequest(c *gin.Context) { addRequestToQueue(4, c) }

//---------------------------------------------------------------

// POST REQUESTS
// ---------------------------------------------------------------
func postStudentsRequest(c *gin.Context) { addRequestToQueue(5, c) }

func postCoursesRequest(c *gin.Context) { addRequestToQueue(6, c) }

func postGradeToStudentByIDRequest(c *gin.Context) { addRequestToQueue(7, c) }

//---------------------------------------------------------------

// handle requests
func handleRequests() {
	for {
		request := <-cin
		switch request.requestType {
		case 1: // GET /students
			request.context.IndentedJSON(http.StatusOK, students)
		case 2: // GET /courses
			request.context.IndentedJSON(http.StatusOK, courses)
		case 3: // GET /student/:id
			getStudentByID(request)
		case 4: // GET /student/:id/courses/:course
			getStudentCourseGradeByID(request)
		case 5: // POST /students
			addStudent(request)
		case 6: // POST /courses
			addCourse(request)
		case 7: // POST /student/:id/course/:course
			setStudentGradeByID(request)
		}
	}
}

func getStudentByID(request RequestBody) {
	c := request.context
	id := c.Param("id")

	for _, s := range students {
		if s.ID == id {
			c.IndentedJSON(http.StatusOK, s)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student not found"})
}

func getStudentCourseGradeByID(request RequestBody) {
	c := request.context
	id := c.Param("id")
	course := c.Param("course")

	for _, s := range students {
		if s.ID == id {
			for name, grade := range s.Courses {
				if name == course {
					c.IndentedJSON(http.StatusOK, grade)
					return
				}
			}
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "course not found"})
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student not found"})
}

func addStudent(request RequestBody) {
	c := request.context
	var newStudent student

	err := c.BindJSON(&newStudent)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not add student"})
		return
	}

	newStudent.ID = strconv.Itoa(id_counter)
	id_counter++

	students = append(students, newStudent)
	c.IndentedJSON(http.StatusCreated, newStudent)
}

func addCourse(request RequestBody) {
	c := request.context
	var newCourse string
	err := c.BindJSON(&newCourse)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "cannot convert to JSON"})
		return
	}
	for i := 0; i < len(courses); i++ {
		if courses[i] == newCourse {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Course already exists"})
			return
		}
	}

	courses = append(courses, newCourse)
	c.IndentedJSON(http.StatusCreated, newCourse)
}

func setStudentGradeByID(request RequestBody) {
	c := request.context
	var student_to_modify student
	err := c.BindJSON(&student_to_modify)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "cannot convert student to JSON"})
		return
	}

	var course_to_set string
	err = c.BindJSON(&course_to_set)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "cannot convert course to JSON"})
		return
	}

	var grade_for_course float64
	err = c.BindJSON(&grade_for_course)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "cannot convert grade to JSON"})
		return
	}

	for _, s := range students {
		if s.ID == student_to_modify.ID {
			if len(s.Courses) == 0 {
				s.Courses = make(map[string]float64)
			}
			s.Courses[course_to_set] = grade_for_course
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given course not found"})
		}
	}
	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "student with given ID not found"})
}

// Goal: every time a request is captured in the server, we want to spawn a goroutine that will handle that request
func main() {
	router := gin.Default()

	for i := 0; i < 5; i++ {
		go handleRequests()
	}

	router.GET("/students", getStudentsRequest)
	router.GET("/courses", getCoursesRequest)
	router.GET("/student/:id", getStudentByIDRequest)
	router.GET("/student/:id/course/:course", getStudentCourseGradeByIdRequest)

	router.POST("/students", postStudentsRequest)
	router.POST("/courses", postCoursesRequest)
	router.POST("/student/", postGradeToStudentByIDRequest)

	go router.Run("localhost:8083")
	go router.Run("localhost:8084")
	router.Run("localhost:8085")
}
