package main

import (
	"encoding/csv"
	"fmt"
	"github.com/xuri/excelize/v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"ptu_gost/auth"
	"regexp"
	"strconv"
)

func main() {
	handler := RegexHandler{}
	initRoutes(&handler)

	err := http.ListenAndServe(":8080", &handler)
	if err != nil {
		log.Fatal(err)
	}
}

func initRoutes(handler *RegexHandler) {
	handler.HandleFunc(regexp.MustCompile(`(?m)^\/$`), index)
	handler.HttpHandler(regexp.MustCompile(`(?m)^\/(?:js|css|images)\/.*\.[A-z]+$`), http.FileServer(http.Dir("public/")))
	handler.HandleFunc(regexp.MustCompile(`(?m)^\/saveperson$`), savePerson)
	handler.HandleFunc(regexp.MustCompile(`(?m)^\/download$`), download)
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("public/mainpage.html")

	if err != nil {
		log.Println(err)
	}

	login := auth.User{}
	login.CheckLogin(r)

	err = tmpl.Execute(w, login)

	if err != nil {
		log.Println("Error template view")
	}
}

func savePerson(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	req.ParseForm()
	scvString := req.FormValue("fio") + "," + req.FormValue("birthday") + "," + req.FormValue("code") + "\n"

	file, err := os.OpenFile("persons.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		_, _ = os.Create("persons.csv")
		file, _ = os.OpenFile("persons.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	}

	file.WriteString(scvString)
	file.Close()
}

func download(res http.ResponseWriter, req *http.Request) {
	createXlsx()
	Openfile, err := os.Open("persons.xlsx") //Open the file to be downloaded later
	file, _ := ioutil.ReadAll(Openfile)
	defer Openfile.Close() //Close after function return

	if err != nil {
		http.Error(res, "File not found.", 404) //return 404 if file is not found
		return
	}

	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	Filename := "persons.xlsx"

	//Set the headers
	res.Header().Set("Content-Description", "File Transfer")
	res.Header().Set("Content-Type", "application/octet-stream")
	res.Header().Set("Content-Disposition", "attachment; filename="+Filename)
	res.Header().Set("Content-Transfer-Encoding", "binary")
	res.Header().Set("Connection", "Keep-Alive")
	res.Header().Set("Expires", "0")
	res.Header().Set("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	res.Header().Set("Content-Length", FileSize)

	res.Write(file)
	os.Remove("persons.xlsx")
}

func createXlsx() {
	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "A1", "ФИО")
	f.SetCellValue("Sheet1", "B1", "Дата рождения")
	f.SetCellValue("Sheet1", "C1", "Код")

	file, err := os.Open("persons.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3
	reader.Comment = '#'

	i := 2
	for {
		record, e := reader.Read()
		if e != nil {
			fmt.Println(e)
			fmt.Println("csv error")
			break
		}

		stringI := strconv.Itoa(i)

		f.SetCellValue("Sheet1", "A"+stringI, record[0])
		f.SetCellValue("Sheet1", "B"+stringI, record[1])
		f.SetCellValue("Sheet1", "C"+stringI, record[2])

		i++
	}
	// Save spreadsheet by the given path.
	if err := f.SaveAs("persons.xlsx"); err != nil {
		fmt.Println(err)
		fmt.Println("save error")
	}
}
