package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	clearData()
	encryptData()
	return router
}

func TestGetCourses(t *testing.T) {
	mockResponse := `["CSSE403","CSSE374","CSSE304"]`
	router := setupRouter()
	router.GET("/courses", getCourses)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/courses", nil)
	router.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	response := strings.ReplaceAll(string(responseData), " ", "")
	response = strings.ReplaceAll(response, "\n", "")

	assert.Equal(t, mockResponse, response)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostCourses(t *testing.T) {
	mockResponse := `"CSSE132"`
	router := setupRouter()
	router.POST("/courses", postCourses)
	router.GET("/courses", getCourses)

	course := "CSSE132"
	val, _ := json.Marshal(course)
	req, _ := http.NewRequest("POST", "/courses", bytes.NewBuffer(val))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	response := strings.ReplaceAll(string(responseData), " ", "")
	response = strings.ReplaceAll(response, "\n", "")

	assert.Equal(t, mockResponse, response)
	assert.Equal(t, http.StatusCreated, w.Code)

	mockResponse = `["CSSE403","CSSE374","CSSE304","CSSE132"]`

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/courses", nil)
	router.ServeHTTP(w, req)
	responseData, _ = ioutil.ReadAll(w.Body)
	response = strings.ReplaceAll(string(responseData), " ", "")
	response = strings.ReplaceAll(response, "\n", "")

	assert.Equal(t, mockResponse, response)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetStudent(t *testing.T) {
	mockResponse := `[{"ID":"1","Name":"JohnDoe","Year":"Junior","Courses":{"CSSE403":100}},{"ID":"2","Name":"JaneDoe","Year":"Sophomore","Courses":{}},{"ID":"3","Name":"DaveSmith","Year":"Senior","Courses":{}}]`
	router := setupRouter()
	router.GET("/students", getStudents)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/students", nil)
	router.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	response := strings.ReplaceAll(string(responseData), " ", "")
	response = strings.ReplaceAll(response, "\n", "")

	assert.Equal(t, mockResponse, response)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostStudent(t *testing.T) {
	mockResponse := `{"ID":"4","Name":"Joey","Year":"Freshman","Courses":{}}`
	router := setupRouter()
	router.POST("/students", postStudents)
	router.GET("/students", getStudents)

	student := student{Name: "Joey", Year: "Freshman"}
	val, _ := json.Marshal(student)
	req, _ := http.NewRequest("POST", "/students", bytes.NewBuffer(val))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	response := strings.ReplaceAll(string(responseData), " ", "")
	response = strings.ReplaceAll(response, "\n", "")

	assert.Equal(t, mockResponse, response)
	assert.Equal(t, http.StatusCreated, w.Code)

	mockResponse = `[{"ID":"1","Name":"JohnDoe","Year":"Junior","Courses":{"CSSE403":100}},{"ID":"2","Name":"JaneDoe","Year":"Sophomore","Courses":{}},{"ID":"3","Name":"DaveSmith","Year":"Senior","Courses":{}},{"ID":"4","Name":"Joey","Year":"Freshman","Courses":{}}]`

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/students", nil)
	router.ServeHTTP(w, req)
	responseData, _ = ioutil.ReadAll(w.Body)
	response = strings.ReplaceAll(string(responseData), " ", "")
	response = strings.ReplaceAll(response, "\n", "")

	assert.Equal(t, mockResponse, response)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostGradeByID(t *testing.T) {
	mockResponse := `{"ID":"3","Name":"DaveSmith","Year":"Senior","Courses":{"CSSE403":40}}`
	router := setupRouter()
	router.POST("/student", postGradeToStudentbyID)
	grades := make(map[string]float64)
	grades["CSSE403"] = 40
	student := student{ID: "3"}
	student.Courses = grades
	val, _ := json.Marshal(student)
	req, _ := http.NewRequest("POST", "/student", bytes.NewBuffer(val))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	response := strings.ReplaceAll(string(responseData), " ", "")
	response = strings.ReplaceAll(response, "\n", "")

	assert.Equal(t, mockResponse, response)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetStudentByID(t *testing.T) {
	mockResponse := `{"ID":"3","Name":"DaveSmith","Year":"Senior","Courses":{}}`
	router := setupRouter()
	router.GET("/student/:id", getStudentByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/student/3", nil)
	router.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	response := strings.ReplaceAll(string(responseData), " ", "")
	response = strings.ReplaceAll(response, "\n", "")

	assert.Equal(t, mockResponse, response)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetStudentGradeByID(t *testing.T) {
	mockResponse := `40`

	students[2].Courses["CSSE403"] = 40

	router := setupRouter()
	router.GET("/student/:id/course/:course", getStudentsGradeById)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/student/3/course/CSSE403", nil)
	router.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	response := strings.ReplaceAll(string(responseData), " ", "")
	response = strings.ReplaceAll(response, "\n", "")

	assert.Equal(t, mockResponse, response)
	assert.Equal(t, http.StatusOK, w.Code)
}
