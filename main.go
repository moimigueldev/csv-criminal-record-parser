package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Report struct {
	Date      string `json:"date"`
	Address   string `json:"address"`
	District  string `json:"district"`
	Beat      string `json:"beat"`
	Grid      string `json:"grid"`
	Crimedesc string `json:"crimedesc"`
	UCRCode   string `json:"ucr_ncic_code"`
	Lat       string `json:"latitude"`
	Lon       string `json:"longitude"`
}

type DateData struct {
	Count           int
	Districts       map[string]int
	PopularDistrict string
}

func main() {
	start := time.Now()
	// Code to measure

	reports := CreateReport()
	dates := ParseDates(reports)
	FindAllPopularDistricts(dates)
	duration := time.Since(start)

	fmt.Println("duration", duration)

}

func CreateReport() []Report {
	csvFile, err := os.Open("./files/sacramento-jan-2006.csv")
	// csvFile, err := os.Open("./files/sacramento-test.csv")
	if err != nil {
		log.Fatal("error opeining file", err)
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	var reports []Report

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("error reading file", err)
		}

		reports = append(reports, Report{
			Date:      line[0],
			Address:   line[1],
			District:  line[2],
			Beat:      line[3],
			Grid:      line[4],
			Crimedesc: line[5],
			UCRCode:   line[6],
			Lat:       line[7],
			Lon:       line[8],
		})

	}

	return reports
}

func ParseDates(reports []Report) map[string]DateData {
	var dates = make(map[string]DateData)

	for _, report := range reports {

		dateIndex := strings.LastIndex(report.Date, "/")
		date := report.Date[0 : dateIndex+3]

		val, ok := dates[date]
		if !ok {
			dates[date] = DateData{
				Count:           1,
				Districts:       make(map[string]int),
				PopularDistrict: "",
			}
			dates[date].Districts[report.District] = 1

		} else {
			val.Count += 1

			_, ok := val.Districts[report.District]
			if !ok {
				val.Districts[report.District] = 0
			}
			val.Districts[report.District] += 1
			dates[date] = val
		}

	}

	return dates
}

func (d DateData) FindPopularDistrict() DateData {

	for key, val := range d.Districts {

		if d.PopularDistrict == "" {
			d.PopularDistrict = key
		}
		if d.Districts[d.PopularDistrict] < val {
			d.PopularDistrict = key
		}
	}

	return d
}

func FindAllPopularDistricts(dates map[string]DateData) {

	for key, _ := range dates {
		dates[key] = dates[key].FindPopularDistrict()
	}
}
