package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const xkcdURL = "https://xkcd.com/%d/info.0.json"

func GetComic(comicNum int) (Comic, error) {
	resp, err := http.Get(fmt.Sprintf(xkcdURL, comicNum))
	if err != nil {
		return Comic{}, fmt.Errorf("failed to GET comic %d: %s", comicNum, err)
	}
	defer resp.Body.Close()

	var comicRetrieved Comic
	if err := json.NewDecoder(resp.Body).Decode(&comicRetrieved); err != nil {
		return Comic{}, fmt.Errorf("failed to GET comic %d: %s", comicNum, err)
	}

	return comicRetrieved, nil
}

func GetLatestComicNumber() (int, error) {
	resp, err := http.Get("https://xkcd.com/info.0.json")
	if err != nil {
		return 0, fmt.Errorf("failed to perform GET request: %s", err)
	}
	defer resp.Body.Close()

	var latest LatestComic
	if err := json.NewDecoder(resp.Body).Decode(&latest); err != nil {
		return 0, fmt.Errorf("failed to GET the latest comic number %d: %s", latest.Num, err)
	}

	fmt.Printf("Latest Comic Number is: %d", latest.Num)

	return latest.Num, nil
}

func RetrieveComicsSequentially(start, end int) ([]Comic, time.Duration) {
	var comics []Comic
	startTime := time.Now()
	for i := start; i <= end; i++ {
		comic, err := GetComic(i)
		if err != nil {
			fmt.Printf("Error retrieving sequentially comic %d: %v\n", i, err)
			continue
		}
		// fmt.Printf("comic number %d is %v\n", i, comic.Title)
		comics = append(comics, comic)
	}
	duration := time.Since(startTime)
	return comics, duration
}

func RetrieveComicsConcurrently(start, end int) ([]Comic, time.Duration) {
	var comics []Comic
	var wg sync.WaitGroup
	startTime := time.Now()

	for i := start; i <= end; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			comic, err := GetComic(num)
			if err != nil {
				fmt.Printf("Error retrieving concurrently comic %d: %v\n", num, err)
				return
			}
			// fmt.Printf("comic number %d is %v\n", i, comic.Title)
			comics = append(comics, comic)
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)
	return comics, duration
}
