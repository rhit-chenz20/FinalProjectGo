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

	// router.POST("/students", postStudents)

	// router.GET("/student/:id", getStudentByID)

	// router.POST("/student", postGradeToStudentbyID)
	// router.GET("/student/:id/course/:course", getStudentsGradeById)
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
	mockResponse := `[{"ID":"1","Name":"JohnDoe","Year":"Junior","Courses":{}},{"ID":"2","Name":"JaneDoe","Year":"Sophomore","Courses":{}},{"ID":"3","Name":"DaveSmith","Year":"Senior","Courses":{}}]`
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

// func TestPostStudent(t *testing.T) {
// 	mockResponse := `"CSSE132"`
// 	router := setupRouter()
// 	router.POST("/courses", postCourses)
// 	router.GET("/courses", getCourses)

// 	course := "CSSE132"
// 	val, _ := json.Marshal(course)
// 	req, _ := http.NewRequest("POST", "/courses", bytes.NewBuffer(val))
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	responseData, _ := ioutil.ReadAll(w.Body)
// 	response := strings.ReplaceAll(string(responseData), " ", "")
// 	response = strings.ReplaceAll(response, "\n", "")

// 	assert.Equal(t, mockResponse, response)
// 	assert.Equal(t, http.StatusCreated, w.Code)

// 	mockResponse = `["CSSE403","CSSE374","CSSE304","CSSE132"]`

// 	w = httptest.NewRecorder()
// 	req, _ = http.NewRequest("GET", "/courses", nil)
// 	router.ServeHTTP(w, req)
// 	responseData, _ = ioutil.ReadAll(w.Body)
// 	response = strings.ReplaceAll(string(responseData), " ", "")
// 	response = strings.ReplaceAll(response, "\n", "")

// 	assert.Equal(t, mockResponse, response)
// 	assert.Equal(t, http.StatusOK, w.Code)
// }
