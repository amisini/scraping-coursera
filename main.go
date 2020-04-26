package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Course struct {
	Title       string
	Description string
	Instructors []string
	URL         string
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("coursera.org", "www.coursera.org"),
	)

	detailCollector := c.Clone()

	courses := []Course{}

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
		courseURL := e.Text
		detailCollector.Visit(courseURL)
	})

	detailCollector.OnHTML(`div[id=rendered-content]`, func(e *colly.HTMLElement) {
		log.Println("Course found", e.Request.URL)
		title := e.ChildText(".banner-title")
		if title == "" {
			log.Println("No title")
			return
		}

		instructors := []string{}

		e.ForEach("h3.instructor-name", func(_ int, el *colly.HTMLElement) {
			instructors = append(instructors, el.Text)
		})

		course := Course{
			Title:       title,
			URL:         e.Request.URL.String(),
			Description: e.ChildText("div.content"),
			Instructors: instructors,
		}

		courses = append(courses, course)

	})

	c.Visit("https://www.coursera.org/sitemap~www~courses.xml")

	file, _ := json.MarshalIndent(courses, "", " ")

	_ = ioutil.WriteFile("coursera.json", file, os.ModePerm)
}
