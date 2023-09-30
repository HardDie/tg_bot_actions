package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	"github.com/HardDie/tg_bot_actions/internal/models"
)

func main() {
	urlFlag := flag.String("url", "", "Where to find a list of criminal law statutes")
	outFileFlag := flag.String("out_file", "criminals.json", "Output json file for criminal law records")
	flag.Parse()

	if *urlFlag == "" || *outFileFlag == "" {
		flag.Usage()
		return
	}

	h, err := downloadPageAndParse(*urlFlag)
	if err != nil {
		log.Fatal(err)
	}

	links := extractLinks(h)
	log.Println("Total count of links <a> on page:", len(links))

	fullBaseURL, err := url.Parse(*urlFlag)
	if err != nil {
		log.Fatal("invalid url:", err)
	}
	baseURL := url.URL{
		Scheme: fullBaseURL.Scheme,
		Host:   fullBaseURL.Host,
	}

	var criminals []models.Criminal
	for _, l := range links {
		value := extractTextFromLink(l)
		if !strings.HasPrefix(value, "Статья") {
			continue
		}
		criminal := textToCriminal(value)
		if criminal == nil {
			continue
		}

		// Set valid link for criminal law
		criminal.Link = extractAttrFromLink(l, "href")
		if criminal.Link != "" {
			criminal.Link = baseURL.String() + criminal.Link
		}

		// All records below 105 are not felonies, skip them
		digitalNumber, err := strconv.ParseFloat(criminal.Number, 64)
		if err != nil {
			log.Println("Error parse criminal number:", err.Error())
			continue
		}
		if digitalNumber < 105 {
			continue
		}

		// If the record is no longer valid, skip it
		if strings.Contains(criminal.Description, "Утратила силу") {
			continue
		}

		criminals = append(criminals, *criminal)
	}
	log.Println("Total count valid criminal laws:", len(criminals))

	err = writeToFile(*outFileFlag, criminals)
	if err != nil {
		log.Fatal("error write data to file", err.Error())
	}
	log.Printf("All criminal law records were saved into %s file\n", *outFileFlag)
}

func downloadPageAndParse(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error get page: %w", err)
	}
	defer resp.Body.Close()
	h, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parse html document: %w", err)
	}
	return h, nil
}
func extractLinks(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, extractLinks(c)...)
	}
	return ret
}
func extractTextFromLink(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += extractTextFromLink(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}
func extractAttrFromLink(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if strings.ToLower(attr.Key) == key {
			return attr.Val
		}
	}
	return ""
}
func textToCriminal(text string) *models.Criminal {
	parts := strings.Split(text, " ")
	if len(parts) < 3 {
		return nil
	}
	return &models.Criminal{
		Number:      strings.TrimSuffix(parts[1], "."),
		Description: strings.Join(parts[2:], " "),
	}
}
func writeToFile(filename string, data any) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating output file %q: %w", filename, err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("error closing file %s: %v\n", filename, err.Error())
		}
	}()

	err = json.NewEncoder(file).Encode(data)
	if err != nil {
		return fmt.Errorf("error encoding data into file %q: %w", filename, err)
	}
	return nil
}
