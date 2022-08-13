package main

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
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
	e.GET("/foxes", getFoxes)
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

	fmt.Println("> [generateFoxes] Foxed Loaded:", foxes)
}

func getFoxes(c echo.Context) error {
	fmt.Println("> [getFoxes]...")
	fmt.Println("> [getFoxes] foxes:", foxes)
	return c.JSON(200, foxes)
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
