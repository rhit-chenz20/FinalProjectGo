package main

import (
	"net/http"
	"strconv"
	"sync"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	//"errors"
	//"fmt"
	"io"
	"github.com/gin-gonic/gin"
	"bufio"
    "fmt"
    "os"
)

// /key, _ := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")

var m sync.Mutex


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
	{ID: "1", Name: "John Doe", Year: "Junior", Courses: make(map[string]float64)},
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
		coursesEncrypted= append(coursesEncrypted,EncryptString(courses[i])) 

 	}
 	for i := 0; i < len(students); i++ {
		id:=EncryptString(students[i].ID)
		name:=EncryptString(students[i].Name)
		year:=EncryptString(students[i].Year)
		courses:=students[i].Courses
		
		
		var student = student{ID: id, Name: name, Year: year, Courses: courses}
		studentsEncrypted = append(studentsEncrypted, student)
 	}

 }

func getCourses(c *gin.Context) {
	m.Lock()
	tempcourses:= []string{}
	for i := 0; i < len(coursesEncrypted); i++ { 		
		tempcourses= append(tempcourses,DecryptString(coursesEncrypted[i])) 
 	}
	c.IndentedJSON(http.StatusOK, tempcourses)
	m.Unlock()
}

func postCourses(c *gin.Context) {
	m.Lock()
	var newCourse string
	err := c.BindJSON(&newCourse)
	if err != nil {
		return
	}
	coursesEncrypted = append(coursesEncrypted, EncryptString(newCourse))
	c.IndentedJSON(http.StatusCreated,  newCourse)
	m.Unlock()
}

func getStudents(c *gin.Context) {
	m.Lock()
	tempstudent:= []student{}
	for i := 0; i < len(studentsEncrypted); i++ {
		id:=DecryptString(studentsEncrypted[i].ID)
		name:=DecryptString(studentsEncrypted[i].Name)
		year:=DecryptString(studentsEncrypted[i].Year)
		courses:=studentsEncrypted[i].Courses
		
		
		var student = student{ID: id, Name: name, Year: year, Courses: courses}
		tempstudent = append(tempstudent, student)
 	}	

	c.IndentedJSON(http.StatusOK, tempstudent)
	m.Unlock()
}

func postStudents(c *gin.Context) {
	m.Lock()
	var newStudent student

	// Call BindJSON to bind the received JSON to
	// newStudent.
	if err := c.BindJSON(&newStudent); err != nil {
		return
	}
	newStudent.ID = EncryptString(strconv.Itoa(id_counter))
	if newStudent.Courses == nil {
		newStudent.Courses = make(map[string]float64)
	}
	id_counter += 1
	// Add the new student to the slice.
	studentsEncrypted = append(studentsEncrypted, newStudent)
	c.IndentedJSON(http.StatusCreated, newStudent)
	m.Unlock()
}

func getStudentByID(c *gin.Context) {
	m.Lock()
	id := EncryptString(c.Param("id"))

	for _, s := range studentsEncrypted {
		if s.ID == id {
			c.IndentedJSON(http.StatusOK, s)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student not found"})
	m.Unlock()


}

// Overwrite a Student's grade in a specified course
func postGradeToStudentbyID(c *gin.Context) {
	m.Lock()
	var student_to_add student
	if err := c.BindJSON(&student_to_add); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "cannot convert to JSON"})
		return
	}
	for _, s := range studentsEncrypted {
		if s.ID == EncryptString(student_to_add.ID) {
			for course, grade := range student_to_add.Courses {
				for _, c := range courses {
					if c == course {
						s.Courses[course] = grade
					}
				}
			}
			c.IndentedJSON(http.StatusCreated, s)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student not found"})
	m.Unlock()
}
func sendToFile(){

    f, err:= os.Create("data.txt")
	if err != nil {
        panic(err)
    }
	w := bufio.NewWriter(f)
	w.WriteString("Courses:\n")
	for i := 0; i < len(coursesEncrypted); i++ {
		 w.WriteString(coursesEncrypted[i]+"\n")

	}
	w.WriteString("Students:\n")
	for i := 0; i < len(studentsEncrypted); i++ {
		w.WriteString("id: "+studentsEncrypted[i].ID+"\n")
		w.WriteString("name: "+studentsEncrypted[i].Name+"\n")
		w.WriteString("year: "+studentsEncrypted[i].Year+"\n")
	//	w.WriteString("courses: "studentsEncrypted[i].Courses+"\n")
		w.WriteString(" \n")

   }

    w.Flush()
}
func sendToFileRequest(c *gin.Context){
	sendToFile()
	c.IndentedJSON(http.StatusOK, "ok")
}
// Get a Student's grade in a specified course
func getStudentsGradeById(c *gin.Context) {
	m.Lock()
	id := c.Param("id")
	course := c.Param("course")

	for _, s := range studentsEncrypted {
		if s.ID == id {
			for name, grade := range s.Courses {
				if name == course {
					c.IndentedJSON(http.StatusOK, grade)
					return
				}
			}
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "student/course not found"})
	m.Unlock()
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
