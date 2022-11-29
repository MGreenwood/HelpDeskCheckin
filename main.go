package main

import (
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"website/database"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

type Student struct {
	FName string
	LName string
	Grade string
}

var studentList map[string]Student
var options [4]string

func main() {
	// load student data
	studentList = make(map[string]Student)
	LoadStudentList()
	options = [...]string{"Borrow a laptop", "Connect to Wifi", "IT Help", "Broken Device"}
	database.Init("../HelpdeskCheckinDatabase.db")

	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve HTML templates
	router.LoadHTMLGlob("./templates/*")

	// Serve frontend static files
	router.Use(static.Serve("/static", static.LocalFile("./static", true)))
	router.StaticFile("/favicon.ico", ".static/favicon.ico")

	// setup public routes
	router.GET("/", IndexHandler)
	router.POST("/checkin", CheckinHandler)
	router.GET("/checkin/confirm-visit", ConfirmationPage)

	router.Run(":8080")
}

func IndexHandler(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.gohtml",
		gin.H{
			"StudentList": studentList,
		},
	)
}

func CheckinHandler(c *gin.Context) {
	id := c.PostForm("id")

	if _, exists := studentList[id]; !exists {
		IndexHandler(c)
		return
	}

	c.HTML(
		http.StatusOK,
		"checkin.gohtml",
		gin.H{
			"StudentInfo": studentList[id],
			"Options":     options,
			"Id":          id,
		},
	)
}

func ConfirmationPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"confirm-checkin.gohtml",
		"Thank you",
	)
}

func LoadStudentList() {
	file, err := os.Open("./static/AllStudents.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		studentList[record[2]] = Student{record[1], record[0], record[3]}
	}
}
