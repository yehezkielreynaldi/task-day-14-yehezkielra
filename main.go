package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"task-day-14/connection"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
)

type Project struct {
	Id           int
	ProjectName  string
	StartDate    time.Time
	EndDate      time.Time
	Duration     string
	Description  string
	Technologies []string
	Image        string
	Tech1        bool
	Tech2        bool
	Tech3        bool
	Tech4        bool
	FormatStart  string
	FormatEnd    string
}

// var dataProject = []Project{
// {
// 	ProjectName: "Project 1",
// 	StartDate:   "2023-05-01",
// 	EndDate:     "2023-06-01",
// 	Duration:    "1 Bulan",
// 	Description: "Ini Project 1",
// 	Tech1:       true,
// 	Tech2:       true,
// 	Tech3:       true,
// 	Tech4:       true,
// },
// {
// 	ProjectName: "Project 2",
// 	// StartDate:   "2023-05-02",
// 	// EndDate:     "2023-06-02",
// 	Duration:    "1 Bulan",
// 	Description: "Ini Project 2",
// 	Tech1:       true,
// 	Tech2:       true,
// 	Tech3:       true,
// 	Tech4:       true,
// },
// {
// 	ProjectName: "Project 3",
// 	StartDate:   "2023-05-03",
// 	EndDate:     "2023-06-03",
// 	Duration:    "1 Bulan",
// 	Description: "Ini Project 3",
// 	Tech1:       true,
// 	Tech2:       true,
// 	Tech3:       true,
// 	Tech4:       true,
// },
// {
// 	ProjectName: "Project 4",
// 	StartDate:   "2023-05-04",
// 	EndDate:     "2023-06-04",
// 	Duration:    "1 Bulan",
// 	Description: "Ini Project 4",
// 	Tech1:       false,
// 	Tech2:       false,
// 	Tech3:       true,
// 	Tech4:       true,
// },
// {
// 	ProjectName: "Project 5",
// 	StartDate:   "2023-05-05",
// 	EndDate:     "2023-06-05",
// 	Duration:    "1 Bulan",
// 	Description: "Ini Project 5",
// 	Tech1:       true,
// 	Tech2:       false,
// 	Tech3:       true,
// 	Tech4:       false,
// },
// {
// 	ProjectName: "Project 6",
// 	StartDate:   "2023-05-06",
// 	EndDate:     "2023-06-06",
// 	Duration:    "1 Bulan",
// 	Description: "Ini Project 6",
// 	Tech1:       true,
// 	Tech2:       false,
// 	Tech3:       true,
// 	Tech4:       true,
// },
// }

func main() {
	connection.DatabaseConnect()

	e := echo.New()

	// e = echo package
	// GET/POST = run the method
	// "/" = endpoint/routing (ex. localhost:5000'/' | ex. dumbways.id'/lms')
	// helloWorld = function that will run if the routes are opened

	// Serve a static files from "public" directory
	e.Static("/public", "public")

	// Routing

	// GET
	e.GET("/", home)
	e.GET("/contact", contact)
	e.GET("/my-project", myProject)
	e.GET("/project-detail/:id", projectDetail)
	e.GET("/testimonials", testimonials)
	e.GET("/update-project/:id", updateMyProject)

	// POST
	e.POST("/add-project", addProject)
	e.POST("/project-delete/:id", deleteProject)
	e.POST("/update-project/:id", updateProject)

	e.Logger.Fatal(e.Start("localhost:5000"))
}

func home(c echo.Context) error {
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, name, start_date, end_date, duration, description, tech1, tech2, tech3, tech4, image FROM tb_proyek ORDER BY id ASC")

	fmt.Println(data)
	var result []Project
	for data.Next() {
		var each = Project{}

		err := data.Scan(&each.Id, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Duration, &each.Description, &each.Tech1, &each.Tech2, &each.Tech3, &each.Tech4, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
		}
		fmt.Println(each)

		each.FormatStart = each.StartDate.Format("2 January 2006")
		each.FormatEnd = each.EndDate.Format("2 January 2006")
		// each.Author = "Abel Dustin"

		result = append(result, each)
	}

	fmt.Println(result)

	projects := map[string]interface{}{
		"Projects": result,
	}

	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil { // null
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), projects)
}

func contact(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"messsage": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func myProject(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/my-project.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func projectDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, duration, description, tech1, tech2, tech3, tech4, image FROM tb_proyek WHERE id=$1", id).Scan(
		&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.Description, &ProjectDetail.Tech1, &ProjectDetail.Tech2, &ProjectDetail.Tech3, &ProjectDetail.Tech4, &ProjectDetail.Image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	ProjectDetail.FormatStart = ProjectDetail.StartDate.Format("2 January 2006")
	ProjectDetail.FormatEnd = ProjectDetail.EndDate.Format("2 January 2006")

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	var tmpl, errTemplate = template.ParseFiles("views/project-detail.html")

	if errTemplate != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

func testimonials(c echo.Context) error {
	var tmpl, err = template.ParseFiles("views/testimonials.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func updateMyProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Project{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, duration, description, tech1, tech2, tech3, tech4, image FROM tb_proyek WHERE id=$1", id).Scan(&ProjectDetail.Id, &ProjectDetail.ProjectName, &ProjectDetail.StartDate, &ProjectDetail.EndDate, &ProjectDetail.Duration, &ProjectDetail.Description, &ProjectDetail.Tech1, &ProjectDetail.Tech2, &ProjectDetail.Tech3, &ProjectDetail.Tech4, &ProjectDetail.Image)

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	var tmpl, errTmplt = template.ParseFiles("views/update-project.html")

	if errTmplt != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return tmpl.Execute(c.Response(), data)
}

func addProject(c echo.Context) error {
	projectName := c.FormValue("projectName")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	duration := hitungDurasi(startDate, endDate)
	description := c.FormValue("descriptionProject")
	tech1 := (c.FormValue("tech1") == "tech1")
	tech2 := (c.FormValue("tech2") == "tech2")
	tech3 := (c.FormValue("tech3") == "tech3")
	tech4 := (c.FormValue("tech4") == "tech4")
	image := c.FormValue("input-image")

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_proyek (name, start_date, end_date, duration, description, tech1, tech2, tech3, tech4, image) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", projectName, startDate, endDate, duration, description, tech1, tech2, tech3, tech4, image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	// fmt.Println(dataProject)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	fmt.Println("Index : ", id)

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_proyek WHERE id=$1", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func updateProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	fmt.Println("Index :", id)

	projectName := c.FormValue("projectName")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	duration := hitungDurasi(startDate, endDate)
	description := c.FormValue("descriptionProject")
	tech1 := (c.FormValue("tech1") == "tech1")
	tech2 := (c.FormValue("tech2") == "tech2")
	tech3 := (c.FormValue("tech3") == "tech3")
	tech4 := (c.FormValue("tech4") == "tech4")

	_, err := connection.Conn.Exec(context.Background(), "UPDATE tb_proyek SET name=$1, start_date=$2, end_date=$3, duration=$4, description=$5, tech1=$6, tech2=$7, tech3=$8, tech4=$9 WHERE id=$10", projectName, startDate, endDate, duration, description, tech1, tech2, tech3, tech4, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func hitungDurasi(startDate, endDate string) string {
	startTime, _ := time.Parse("2006-01-02", startDate)
	endTime, _ := time.Parse("2006-01-02", endDate)

	durationTime := int(endTime.Sub(startTime).Hours())
	durationDays := durationTime / 24
	durationWeeks := durationDays / 7
	durationMonths := durationWeeks / 4
	durationYears := durationMonths / 12

	var duration string

	if durationYears > 1 {
		duration = strconv.Itoa(durationYears) + " Tahun"
	} else if durationYears > 0 {
		duration = strconv.Itoa(durationYears) + " Tahun"
	} else {
		if durationMonths > 1 {
			duration = strconv.Itoa(durationMonths) + " Bulan"
		} else if durationMonths > 0 {
			duration = strconv.Itoa(durationMonths) + " Bulan"
		} else {
			if durationWeeks > 1 {
				duration = strconv.Itoa(durationWeeks) + " Minggu"
			} else if durationWeeks > 0 {
				duration = strconv.Itoa(durationWeeks) + " Minggu"
			} else {
				if durationDays > 1 {
					duration = strconv.Itoa(durationDays) + " Hari"
				} else {
					duration = strconv.Itoa(durationDays) + " Hari"
				}
			}
		}
	}

	return duration
}
