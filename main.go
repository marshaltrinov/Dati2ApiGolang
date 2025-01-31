package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

type Location struct {
	PostalCode        string `json:"PostalCode"`
	KelurahanCode     string `json:"KelurahanCode"`
	KelurahanName     string `json:"KelurahanName"`
	KecamatanCode     string `json:"KecamatanCode"`
	KecamatanName     string `json:"KecamatanName"`
	Dati2Code         string `json:"Dati2Code"`
	Dati2Name         string `json:"Dati2Name"`
	IsDati2Flag       string `json:"IsDati2Flag"`
	MainKelurahanCode string `json:"MainKelurahanCode"`
	MainKecamatanCode string `json:"MainKecamatanCode"`
	MainDati2Code     string `json:"MainDati2Code"`
	CityCode          string `json:"CityCode"`
	CityName          string `json:"CityName"`
	ProvinceCode      string `json:"ProvinceCode"`
	ProvinceName      string `json:"ProvinceName"`
}

type InsertRequest struct {
	Dati2Data struct {
		Row []Location `json:"Row"`
	} `json:"Dati2Data"`
}

type Dati2CodeRequest struct {
	Dati2Code string `json:"Dati2Code"`
}

func initDB() *sql.DB {
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	query := `
	CREATE TABLE IF NOT EXISTS locations (
		PostalCode TEXT,
		KelurahanCode TEXT,
		KelurahanName TEXT,
		KecamatanCode TEXT,
		KecamatanName TEXT,
		Dati2Code TEXT,
		Dati2Name TEXT,
		IsDati2Flag TEXT,
		MainKelurahanCode TEXT,
		MainKecamatanCode TEXT,
		MainDati2Code TEXT,
		CityCode TEXT,
		CityName TEXT,
		ProvinceCode TEXT,
		ProvinceName TEXT
	);`
	_, err = db.Exec(query)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return db
}

func insertData(c *gin.Context) {
	var req InsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := initDB()
	defer db.Close()

	for _, row := range req.Dati2Data.Row {
		query := `
		INSERT INTO locations (
			PostalCode, KelurahanCode, KelurahanName, KecamatanCode, KecamatanName,
			Dati2Code, Dati2Name, IsDati2Flag, MainKelurahanCode, MainKecamatanCode,
			MainDati2Code, CityCode, CityName, ProvinceCode, ProvinceName
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		_, err := db.Exec(query,
			row.PostalCode, row.KelurahanCode, row.KelurahanName, row.KecamatanCode, row.KecamatanName,
			row.Dati2Code, row.Dati2Name, row.IsDati2Flag, row.MainKelurahanCode, row.MainKecamatanCode,
			row.MainDati2Code, row.CityCode, row.CityName, row.ProvinceCode, row.ProvinceName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Data inserted successfully"})
}

func getData(c *gin.Context) {
	db := initDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM locations")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var locations []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		location := make(map[string]interface{})
		for i, col := range columns {
			location[col] = values[i]
		}
		locations = append(locations, location)
	}

	c.JSON(http.StatusOK, locations)
}

func getDataByDati2Code(c *gin.Context) {
	var req Dati2CodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dati2Code is required in the request body"})
		return
	}

	db := initDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM locations WHERE Dati2Code = ?", req.Dati2Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var locations []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		location := make(map[string]interface{})
		for i, col := range columns {
			location[col] = values[i]
		}
		locations = append(locations, location)
	}

	c.JSON(http.StatusOK, locations)
}

func main() {
	r := gin.Default()

	r.POST("/insert", insertData)
	r.GET("/data", getData)
	r.POST("/data/by-dati2code", getDataByDati2Code)

	if err := r.Run(":8080"); err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Example: Insert a new record
	_, err = db.Exec("INSERT INTO locations (PostalCode, KelurahanCode) VALUES (?, ?)", "12345", "KEL001")
	if err != nil {
		log.Fatal(err)
	}
}
