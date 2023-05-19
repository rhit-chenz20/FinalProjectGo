package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"

	// "sync"

	//"errors"
	//"fmt"
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

// /key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")

// var m sync.Mutex

func EncryptString(plaintext string) string {
	block, err := newCipherBlock("0")
	if err != nil {
		return ""
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return ""
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return fmt.Sprintf("%x", ciphertext)
}

func DecryptString(cipherHex string) string {
	block, err := newCipherBlock("0")
	if err != nil {
		return ""
	}

	ciphertext, err := hex.DecodeString(cipherHex)
	if err != nil {
		return ""
	}

	if len(ciphertext) < aes.BlockSize {
		return ""
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext)
}
func newCipherBlock(key string) (cipher.Block, error) {
	hasher := md5.New()
	fmt.Fprint(hasher, key)
	cipherKey := hasher.Sum(nil)
	return aes.NewCipher(cipherKey)
}

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
	{ID: "1", Name: "John Doe", Year: "Junior", Courses: map[string]float64{"CSSE403": 100}},
	{ID: "2", Name: "Jane Doe", Year: "Sophomore", Courses: make(map[string]float64)},
	{ID: "3", Name: "Dave Smith", Year: "Senior", Courses: make(map[string]float64)},
}

// Encrypted courses
var coursesEncrypted = []string{}

// Encryted students
var studentsEncrypted = []student{}

var id_counter = 4

func encryptData() {
	for i := 0; i < len(courses); i++ {
		coursesEncrypted = append(coursesEncrypted, EncryptString(courses[i]))

	}
	for i := 0; i < len(students); i++ {
		id := EncryptString(students[i].ID)
		name := EncryptString(students[i].Name)
		year := EncryptString(students[i].Year)
		courses := make(map[string]float64)
		for name, grade := range students[i].Courses {
			courses[EncryptString(name)] = grade
		}

		var student = student{ID: id, Name: name, Year: year, Courses: courses}
		studentsEncrypted = append(studentsEncrypted, student)
	}

}

func clearData() {
	studentsEncrypted = []student{}
	coursesEncrypted = []string{}
}

func getCourses(c *gin.Context) {
	//m.Lock()
	tempcourses := []string{}
	for i := 0; i < len(coursesEncrypted); i++ {
		tempcourses = append(tempcourses, DecryptString(coursesEncrypted[i]))
	}
	c.IndentedJSON(http.StatusOK, tempcourses)
	//m.Unlock()
}

func postCourses(c *gin.Context) {
	//m.Lock()
	var newCourse string
	err := c.BindJSON(&newCourse)
	if err != nil {
		return
	}
	coursesEncrypted = append(coursesEncrypted, EncryptString(newCourse))
	c.IndentedJSON(http.StatusCreated, newCourse)
	//m.Unlock()
}

func getStudents(c *gin.Context) {
	//m.Lock()
	tempstudent := []student{}
	for i := 0; i < len(studentsEncrypted); i++ {
		id := DecryptString(studentsEncrypted[i].ID)
		name := DecryptString(studentsEncrypted[i].Name)
		year := DecryptString(studentsEncrypted[i].Year)
		//courses:=studentsEncrypted[i].Courses
		courses := make(map[string]float64)
		for name, grade := range studentsEncrypted[i].Courses {
			courses[DecryptString(name)] = grade
		}

		var student = student{ID: id, Name: name, Year: year, Courses: courses}
		tempstudent = append(tempstudent, student)
	}

	c.IndentedJSON(http.StatusOK, tempstudent)
	//m.Unlock()
}

func postStudents(c *gin.Context) {
	//m.Lock()
	var newStudent student
	var encStudent student

	// Call BindJSON to bind the received JSON to
	// newStudent.
	if err := c.BindJSON(&newStudent); err != nil {
		return
	}

	fmt.Println("NEW STUDENT: ", newStudent)
	encStudent.ID = EncryptString(strconv.Itoa(id_counter))
	newStudent.ID = strconv.Itoa(id_counter)
	if newStudent.Name != "" {
		encStudent.Name = EncryptString(newStudent.Name)
	}
	if newStudent.Year != "" {
		encStudent.Year = EncryptString(newStudent.Year)
	}
	if newStudent.Courses == nil {
		encStudent.Courses = make(map[string]float64)
		newStudent.Courses = make(map[string]float64)
	} else {
		encryptCourse := make(map[string]float64)
		nromalCourse := make(map[string]float64)
		for name, grade := range newStudent.Courses {
			encryptCourse[EncryptString(name)] = grade
			nromalCourse[name] = grade
		}
		newStudent.Courses = nromalCourse
		encStudent.Courses = encryptCourse
	}
	id_counter += 1
	// Add the new student to the slice.
	studentsEncrypted = append(studentsEncrypted, encStudent)
	c.IndentedJSON(http.StatusCreated, newStudent)
	//m.Unlock()
}

func getStudentByID(c *gin.Context) {
	//m.Lock()
	id := (c.Param("id"))
	for _, s := range studentsEncrypted {
		fmt.Println("encrypted ID: ", s.ID)
		if DecryptString(s.ID) == id {
			id := DecryptString(s.ID)
			name := DecryptString(s.Name)
			year := DecryptString(s.Year)
			courses := make(map[string]float64)
			for name, grade := range s.Courses {
				courses[DecryptString(name)] = grade
			}
			var student = student{ID: id, Name: name, Year: year, Courses: courses}
			c.IndentedJSON(http.StatusOK, student)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student not found"})
	//m.Unlock()

}

// Overwrite a Student's grade in a specified course
func postGradeToStudentbyID(c *gin.Context) {
	//m.Lock()
	var student_to_add student
	if err := c.BindJSON(&student_to_add); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "cannot convert to JSON"})
		return
	}
	for _, s := range studentsEncrypted {
		if DecryptString(s.ID) == student_to_add.ID {
			for course, grade := range student_to_add.Courses {
				for _, c := range coursesEncrypted {
					if DecryptString(c) == course {
						fmt.Println("UPDATING COURSE" + course)
						s.Courses[c] = grade
					}
				}
			}
			var newStudent student
			newStudent.ID = DecryptString(s.ID)
			newStudent.Name = DecryptString(s.Name)
			newStudent.Year = DecryptString(s.Year)
			newStudent.Courses = make(map[string]float64)
			for c, g := range s.Courses {

				newStudent.Courses[DecryptString(c)] = g

			}
			c.IndentedJSON(http.StatusCreated, newStudent)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student not found"})
	//m.Unlock()
}
func sendToFile() {

	f, err := os.Create("data.txt")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	w.WriteString("Courses:\n")
	for i := 0; i < len(coursesEncrypted); i++ {
		w.WriteString(coursesEncrypted[i] + "\n")

	}
	w.WriteString("Students:\n")
	for i := 0; i < len(studentsEncrypted); i++ {
		w.WriteString("id: " + studentsEncrypted[i].ID + "\n")
		w.WriteString("name: " + studentsEncrypted[i].Name + "\n")
		w.WriteString("year: " + studentsEncrypted[i].Year + "\n")
		for name, grade := range studentsEncrypted[i].Courses {
			s := fmt.Sprintf("%f", grade)
			w.WriteString("Grade for " + name + ": " + s + " \n")
		}
		//	w.WriteString("courses: "studentsEncrypted[i].Courses+"\n")
		w.WriteString(" \n")

	}

	w.Flush()
}
func sendToFileRequest(c *gin.Context) {
	sendToFile()
	c.IndentedJSON(http.StatusOK, "ok")
}

// Get a Student's grade in a specified course
func getStudentsGradeById(c *gin.Context) {
	//m.Lock()
	id := c.Param("id")
	course := c.Param("course")

	for _, s := range studentsEncrypted {
		if DecryptString(s.ID) == id {
			fmt.Println("Looking for course!!")
			for name, grade := range s.Courses {
				if DecryptString(name) == course {
					c.IndentedJSON(http.StatusOK, grade)
					return
				}
			}
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student/course not found"})
	//m.Unlock()
}

func main() {
	router := gin.Default()
	encryptData()
	sendToFile()
	router.GET("/students", getStudents)
	router.POST("/students", postStudents)

	router.GET("/courses", getCourses)
	router.POST("/courses", postCourses)

	router.GET("/student/:id", getStudentByID)

	router.POST("/student", postGradeToStudentbyID)
	router.GET("/student/:id/course/:course", getStudentsGradeById)
	router.GET("/grades", sendToFileRequest)
	router.Run("localhost:8080")

}
