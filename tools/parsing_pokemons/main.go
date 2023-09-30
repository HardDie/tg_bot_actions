package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"sort"

	"github.com/go-resty/resty/v2"

	"github.com/HardDie/tg_bot_actions/internal/models"
)

func main() {
	urlFlag := flag.String("url", "", "Where to find a list of pokemons")
	fileFlag := flag.String("file", "", "File with pokemon description")
	outFileFlag := flag.String("out_file", "pokemons.json", "Output json file for pokemons")
	flag.Parse()

	if (*urlFlag == "" && *fileFlag == "") || *outFileFlag == "" {
		flag.Usage()
		return
	}

	var err error
	var data []byte
	if *fileFlag != "" {
		data, err = openAndRead(*fileFlag)
	} else {
		data, err = downloadPage(*urlFlag)
	}
	if err != nil {
		log.Fatal(err)
	}

	var list []models.Pokemon
	err = json.Unmarshal(data, &list)
	if err != nil {
		log.Fatal("error parse data:", err.Error())
	}

	fullBaseURL, err := url.Parse(*urlFlag)
	if err != nil {
		log.Fatal("invalid url:", err)
	}
	baseURL := url.URL{
		Scheme: fullBaseURL.Scheme,
		Host:   fullBaseURL.Host,
	}

	// There are duplicates in the document, replace all previous elements with the latest ones
	pokemons := make(map[int]models.Pokemon)
	for _, l := range list {
		// Set full url for detailed page
		l.DetailPageURL = baseURL.String() + l.DetailPageURL
		pokemons[l.ID] = l
	}

	// Convert map to sorted slice
	pokemonList := make([]models.Pokemon, 0, len(pokemons))
	for _, p := range pokemons {
		pokemonList = append(pokemonList, p)
	}
	sort.Slice(pokemonList, func(i, j int) bool {
		return pokemonList[i].ID < pokemonList[j].ID
	})

	log.Println("Count of pokemons:", len(pokemonList))

	err = writeToFile(*outFileFlag, pokemonList)
	if err != nil {
		log.Fatal("error write data to file", err.Error())
	}
	log.Printf("All pokemons were saved into %s file\n", *outFileFlag)
}

func downloadPage(url string) ([]byte, error) {
	cli := resty.New().
		SetRedirectPolicy(resty.FlexibleRedirectPolicy(20)).
		SetHeaders(map[string]string{
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"Accept-Encoding":           "gzip, deflate, br",
			"Accept-Language":           "en-US,en;q=0.9",
			"Sec-Ch-Ua":                 "\"Chromium\";v=\"115\", \"Not/A)Brand\";v=\"99\"",
			"Sec-Ch-Ua-Mobile":          "?0",
			"Sec-Ch-Ua-Platform":        "\"Linux\"",
			"Sec-Fetch-Dest":            "document",
			"Sec-Fetch-Mode":            "navigate",
			"Sec-Fetch-Site":            "none",
			"Sec-Fetch-User":            "?1",
			"Upgrade-Insecure-Requests": "1",
			"User-Agent":                "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36",
		})

	resp, err := cli.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("error get page: %w", err)
	}

	return resp.Body(), nil
}
func openAndRead(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error open file %s: %w", filename, err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error read data from file %s: %w", filename, err)
	}
	return data, nil
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
