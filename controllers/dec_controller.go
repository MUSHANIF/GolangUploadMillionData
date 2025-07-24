package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golangReact/database"
	"golangReact/models"
)
func ptrInt(v int) *int {
	return &v
} 

func parseDate(value string) (string, error) {
	t, err := time.Parse("02/01/2006", value)
	if err != nil {
		return "", err
	}
	return t.Format("2006-01-02"), nil
}

func ImportDeclareFromCSV(c *gin.Context) {
	start := time.Now()

	// Ambil file dari form-data
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	csvFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file"})
		return
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	// Set delimiter to semicolon pake koma blog
	reader.Comma = ';'

	// Skip header
	if _, err := reader.Read(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read CSV header"})
		return
	}

	// Ambil No Resi yang sudah ada di database
	var existingResis []string
	if err := database.DB.Model(&models.DeclareExcel{}).Where("status = ?", 1).Pluck("no_resi", &existingResis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Buat map untuk cek No Resi yang sudah ada
	resiMap := make(map[string]struct{})
	for _, resi := range existingResis {
		resiMap[resi] = struct{}{}
	}

	// Konstanta bisnis
	const (
		limitAnyone        = 2000000000
		minPremium         = 300.0
		rateInsuredMaster  = 0.2
		defaultStatus      = 1
	)

	var declares []models.DeclareExcel
	rowNumber := 2

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error reading CSV row %d: %v", rowNumber, err)})
			return
		}
		//digunakan untuk melewati (skip) baris yang tidak memiliki cukup data (kurang dari 10 kolom/field)
		if len(record) < 10 {
			fmt.Printf("Skipping row %d: too few fields (%d)\n", rowNumber, len(record))
			rowNumber++
			continue
		}

		// Parsing tanggal
		conoteDate, err := parseDate(record[1])
		if err != nil {
			fmt.Printf("Skipping row %d: invalid date format: %v\n", rowNumber, record[1])
			rowNumber++
			continue
		}

		noResi := record[5]
		status := defaultStatus
		reason := ""
		if _, exists := resiMap[noResi]; exists {
			reason = "No Resi already exists"
			status = 4
		}

		sumInsuredRaw := strings.ReplaceAll(record[8], ",", "")
		sumInsuredFloat, err := strconv.ParseFloat(sumInsuredRaw, 64)
		if err != nil {
			reason = "Invalid sum insured format"
			status = 4
		}

		if sumInsuredFloat > float64(limitAnyone) {
			reason = "Sum insured exceeds limit"
			status = 4
		}

		// Perhitungan premi
		premiumCalc := sumInsuredFloat * (rateInsuredMaster / 100)
		validatedPremium := math.Max(minPremium, premiumCalc)

		// Perhitungan rate
		var rate float64
		if sumInsuredFloat > 0 {
			rate = validatedPremium * 100 / sumInsuredFloat
		}
		
		// Buat objek untuk dimasukkan ke database
		recordData := models.DeclareExcel{
			ProductCode: record[0],
			ConoteDate:  conoteDate,
			CustNo:      record[2],
			CustName:    record[3],
			NoStikb:     record[4],
			NoResi:      noResi,
			Origin:      record[6],
			Destination: record[7],
			SumInsured:  record[8],
			Premium:     fmt.Sprintf("%.2f", validatedPremium),
			Rate:        fmt.Sprintf("%.4f", rate),    
			Status:      status,
			Reason:     reason,
			UpdatedBy: ptrInt(1),
		}

		declares = append(declares, recordData)
		rowNumber++
	}

	// Simpan batch
	if err := database.DB.CreateInBatches(declares, 1000).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	duration := time.Since(start)
	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("Imported %d rows successfully", len(declares)),
		"duration": duration.String(),
	})
}
