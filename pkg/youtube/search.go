package youtube

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/sirupsen/logrus"
	"strings"
)

func GetVideosFromSearch(query string) []*QueryResult {

	fmt.Println("Querying for... "+query)
/*
	client := &http.Client{}

	resp, err := client.Get(fmt.Sprintf("https://www.youtube.com/results?search_query=%s", query))
	if err != nil {
		logrus.Error("Issue querying for: "+query)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
*/

	fmt.Println("Search query: "+fmt.Sprintf("https://www.youtube.com/results?search_query=%s", strings.Replace(strings.TrimSpace(query), " ", "+", -1)))

	htmlstring, err := soup.Get(fmt.Sprintf("https://www.youtube.com/results?search_query=%s", strings.Replace(strings.TrimSpace(query), " ", "+", -1)))
	if err != nil {
		logrus.Println("There was an issue making the request....")
	}

	//fmt.Println(htmlstring)

	return ParseQueryResults(htmlstring)
}

type QueryResult struct {
	Url string
	Title string
}

func ParseQueryResults(html string) []*QueryResult {

	doc := soup.HTMLParse(html)

	fmt.Println("Text: "+html)

	links := doc.Find("div", "id", "results").FindAll("div", "class", "yt-lockup-content")

	/*
	fmt.Println(links.NodeValue)


	for _, attribute := range links.Attrs() {
		fmt.Println(attribute)
	}
	*/

	var results []*QueryResult

	for _, block := range links {

		link := block.Find("a")
		fmt.Println(link.Text(), "| Link :", link.Attrs()["href"], " | Title: ", link.Attrs()["title"])
		results = append(results, &QueryResult{link.Attrs()["href"], link.Attrs()["title"]})
	}

	return results
}
