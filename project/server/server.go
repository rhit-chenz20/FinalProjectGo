package main

import (
	"encoding/hex"
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

// Encrypted courses
var coursesEncrypted = []string{}

// Encryted
var studentsEncrypted = []student{}

var id_counter = 4

func encryptData() {
	for i := 0; i < len(courses); i++ {
		encodeCourse(courses[i])

	}
	for i := 0; i < len(students); i++ {
		encodeStudent(students[i])
	}

}
func encodeCourse(course string) {

	// block, err := newCipherBlock(0)
	// if err != nil {
	// 	return "", err
	// }

	// ciphertext := make([]byte, aes.BlockSize+len(course))
	// iv := ciphertext[:aes.BlockSize]
	// if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	// 	return "", err
	// }

	// stream := cipher.NewCFBEncrypter(block, iv)
	// stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	// return fmt.Sprintf("%x", ciphertext), nil

	// converted := []byte(course)

	// dst := hex.EncodeToString(converted)
	// coursesEncrypted = append(coursesEncrypted, dst)
	// decoded, _ := hex.DecodeString(dst)
	// // fmt.Printf("%s\n", decoded)
}
func encodeStudent(stu student) {
	convertedS := []byte(stu.ID)
	convertedN := []byte(stu.Name)
	convertedY := []byte(stu.Year)
	dstS := make([]byte, hex.EncodedLen(len(convertedS)))
	dstN := make([]byte, hex.EncodedLen(len(convertedN)))
	dstY := make([]byte, hex.EncodedLen(len(convertedY)))
	hex.Encode(dstS, convertedS)
	hex.Encode(dstN, convertedN)
	hex.Encode(dstY, convertedY)
	studentsEncrypted = append(studentsEncrypted, student{ID: string(dstS), Name: string(dstN), Year: string(dstY)})
}

// func decryptStudent() string {

// }

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
	var student_to_add student
	if err := c.BindJSON(&student_to_add); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "cannot convert to JSON"})
		return
	}
	for _, s := range students {
		if s.ID == student_to_add.ID {
			if len(s.Courses) == 0 {
				s.Courses = make(map[string]float64, len(student_to_add.Courses))
			}
			for course, grade := range student_to_add.Courses {
				s.Courses[course] = grade
			}
			c.IndentedJSON(http.StatusOK, s)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student not found"})
}

// Get a Student's grade in a specified course
func getStudentsGradeById(c *gin.Context) {
	id := c.Param("id")
	course := c.Param("course")

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
	encryptData()
	router.GET("/students", getStudents)
	router.POST("/students", postStudents)

	router.GET("/courses", getCourses)
	router.POST("/courses", postCourses)

	router.GET("/student/:id", getStudentByID)

	router.POST("/student", postGradeToStudentbyID)
	router.GET("/student/:id/course/:course", getStudentsGradeById)

	go router.Run("localhost:8080")
	go router.Run("localhost:8081")
	router.Run("localhost:8082")
}