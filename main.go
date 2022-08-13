package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Foxes struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

var foxes []Foxes

func main() {
	fmt.Println(">\n\n --------------------------------------------------")
	fmt.Println(">[main]...")

	// Load foxes from db
	generateFoxes()

	// Echo instance start
	e := echo.New()

	//ALLOW CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{echo.GET},
	}))
	e.GET("/foxes", getFoxes)
	e.GET("/foxes/code/:code", getFoxesByCode)
	e.GET("/foxes/description/:name", getFoxesByDescription)
	// e.POST("/foxes", createFoxes)

	// Start server on port 8080
	e.Logger.Fatal(e.Start(":8080"))

}

func generateFoxes() {
	fmt.Println("> [generateFoxes]...")

	// Load foxes from db
	loadFoxes()

	// foxes = append(foxes, Foxes{
	// 	Code:        404,
	// 	Description: "Not Found",
	// 	Url:         "assets/images/404.png",
	// })

	fmt.Println("> [generateFoxes] Foxed Loaded!")
}

func getFoxes(c echo.Context) error {
	fmt.Println("> [getFoxes]...")
	// fmt.Println("> [getFoxes] foxes:", foxes)
	return c.JSON(200, foxes)
}

func getFoxesByDescription(c echo.Context) error {
	fmt.Println("> [getFoxesByDescription]...")
	description := c.Param("name")
	fmt.Println("> [getFoxesByDescription] Received: ", description)

	temp_foxes := []Foxes{}
	for _, fox := range foxes {
		if strings.Contains(strings.ToUpper(fox.Description), strings.ToUpper(description)) {
			temp_foxes = append(temp_foxes, fox)
		}
	}
	if len(temp_foxes) > 0 {
		fmt.Println("> [getFoxesByDescription] OK")
		return c.JSON(http.StatusOK, temp_foxes)
	}

	return c.String(http.StatusNotFound, "Fox not found")
}

func getFoxesByCode(c echo.Context) error {
	fmt.Println("> [getFoxesByCode]...")
	codeString := c.Param("code")
	fmt.Println("> [getFoxesByCode] Received: ", codeString)
	code, err := strconv.Atoi(codeString)
	if err != nil {
		fmt.Println("> [getFoxesByCode]strconv.Atoi err:", err)
		return c.String(http.StatusBadRequest, "Invalid fox code")
	}

	fmt.Println("> [getFoxesByCode] Fox Code:", code)

	temp_foxes := []Foxes{}
	for _, fox := range foxes {
		if fox.Code == code {
			temp_foxes = append(temp_foxes, fox)
		}
	}
	if len(temp_foxes) > 0 {
		fmt.Println("> [getFoxesByCode] OK")
		return c.JSON(http.StatusOK, temp_foxes)
	}

	return c.String(http.StatusNotFound, "Fox not found")
}

func createFoxes(c echo.Context) error {
	fmt.Println("> [createFoxes]...")
	fox := new(Foxes)
	if err := c.Bind(fox); err != nil {
		return err
	}
	foxes = append(foxes, *fox)
	saveFox(*fox)
	return c.JSON(200, fox)
}

func loadFoxes() error {
	fmt.Println("> [loadFoxes]...")
	db, err := sql.Open("sqlite3", "foxes.db")
	if err != nil {
		fmt.Println("> [loadFoxes]db err:", err)
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM foxes")
	if err != nil {
		fmt.Println("> [loadFoxes]db rows err:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var code int
		var description string
		var url string
		if err := rows.Scan(&code, &description, &url); err != nil {
			fmt.Println("> [loadFoxes]db rows.Scan err:", err)
			return err
		}
		foxes = append(foxes, Foxes{
			Code:        code,
			Description: description,
			Url:         url,
		})
	}
	if err := rows.Err(); err != nil {
		fmt.Println("> [loadFoxes]db rows.Err err:", err)
		return err
	}
	fmt.Println("> [loadFoxes] OK")
	return nil
}

func saveFox(fox Foxes) error {
	fmt.Println("> [saveFox]...")
	db, err := sql.Open("sqlite3", "foxes.db")
	if err != nil {
		fmt.Println("> [saveFox]db err:", err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO foxes (code, description, url) VALUES($1, $2, $3)")
	if err != nil {
		fmt.Println("> [saveFox]db stmt err:", err)
		return err
	}

	_, err = stmt.Exec(fox.Code, fox.Description, fox.Url)
	if err != nil {
		fmt.Println("> [saveFox]db stmt.Exec err:", err)
		return err
	}
	fmt.Println("> [saveFox] OK")
	return nil
}
