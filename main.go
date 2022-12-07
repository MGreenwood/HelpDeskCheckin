package main

import (
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func main() {
	// load student data
	studentList = make(map[string]Student)
	LoadStudentList()
	database.Init("\\\\hs-hyperhost-01\\tech3\\HighSchool\\Kiosk\\HelpDeskCheckInDatabase.db")

	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Serve HTML templates
	router.LoadHTMLGlob("./templates/*")

	// Serve frontend static files
	router.Use(static.Serve("/static", static.LocalFile("./static", true)))
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	// setup public routes
	router.GET("/", IndexHandler)
	router.POST("/checkin", CheckinHandler)
	router.GET("/checkin/confirm-visit/*option", ConfirmationPage)

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
			"Options":     database.Options,
			"StudentId":   id,
		},
	)
}

func ConfirmationPage(c *gin.Context) {
	// option - index in that option - id #
	params := strings.Split(c.Param("option"), "/")
	option := params[1]
	ndx, err := strconv.Atoi(params[2])
	if err != nil {
		panic(err)
	}
	id := params[3]

	database.CheckIn(database.Query{Id:id, Topic:option, Option:ndx})

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
