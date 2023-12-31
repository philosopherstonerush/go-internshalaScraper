package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
	_ "github.com/lib/pq"
)

const URL = "https://internshala.com/internships/work-from-home-%s-internships-in-chennai/"

const connStr = "user=postgres password=test dbname=postgres sslmode=disable"

type offer struct {
	id      string
	company string
	posted  string
	stipend string
	link    string
}

var offers = []offer{}
var skills_internshala = []string{"python", "django", "flutter", "flutter-development", "c-programming", "sql", "mysql", "bash", "java", "hibernate-java", "rust", "javascript", "javascript-development", "data-analytics", "data-science", "database-building", "embedded-systems", "arduino", "machine-learning", "artificial-intelligence-ai"}

func main() {

	defer addtoPostGres()

	// new colly instance with async
	c := colly.NewCollector(
		colly.Async(true),
	)

	// rate limit
	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	url := []string{}
	for _, value := range skills_internshala {
		temp := fmt.Sprintf(URL, value) // convert the base url to an url with skill inserted
		url = append(url, temp)
	}

	// Callback method extract that gets called whenever an element div with attribute internshipid is found
	c.OnHTML("div[internshipid]", extract)

	// Visit all the pages with the skills inserted
	for _, link := range url {
		c.Visit(link)
	}

	c.Wait() // wait for all the threads to finish executing
}

func extract(r *colly.HTMLElement) {
	temp := offer{}

	//select the value of attribute "internshipid" to uniquely identify the job details
	temp.id = r.Attr("internshipid")

	//select the child text with the class mentioned. These are go queries
	temp.company = strings.ReplaceAll(r.ChildText(".link_display_like_text"), "'", "''")
	moneyGiveth := extractStipendNumber(r.ChildText(".stipend"))
	if moneyGiveth != "" {
		temp.stipend = moneyGiveth
	} else {
		temp.stipend = "0"
	}

	temp.posted = time.Now().Format("01-02-2006")

	// Form the whole internshala url
	temp.link = "www.internshala.com" + r.ChildAttr("a.view_detail_button", "href")
	offers = append(offers, temp)
}

func tranform2D(d []offer) [][]string {
	results := [][]string{}
	results = append(results, []string{"id", "company", "stipend", "link"}) // append the column headers
	for _, value := range d {
		results = append(results, []string{value.id, value.company, value.posted, value.stipend, value.link})
	}
	return results
}

func removeDup(d []offer) []offer {
	processed := make(map[string]bool) // ids already scraped
	result := make([]offer, 0)         // unique structs
	for _, val := range d {
		if _, ok := processed[val.id]; ok {
			continue
		}
		result = append(result, val)
		processed[val.id] = true
	}
	return result
}

func extractStipendNumber(num string) string {
	re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
	return re.FindString(num)
}

func outputInCSV() {
	// remove duplicates
	offers = removeDup(offers)

	// Create a csv file to store the extracted information
	f, err := os.Create("file.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// convert to [][]string or list of list of string from struct. Marshalling the struct
	offersString := tranform2D(offers)

	// Initialise csv writer and then write all of them to the file. CSV requries [][]string to write.
	w := csv.NewWriter(f)
	w.WriteAll(offersString)
	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}

func addtoPostGres() {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	dbApi := Dbase{
		db: db,
	}
	count := 0
	dbApi.addAll(offers)
	fmt.Printf("%v", count)
}
