package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"xkcdcomics/model"
)

const xkcdURL = "https://xkcd.com/%d/info.0.json"

func GetComic(comicNum int) (model.Comic, error) {
	resp, err := http.Get(fmt.Sprintf(xkcdURL, comicNum))
	if err != nil {
		return model.Comic{}, fmt.Errorf("failed to GET comic %d: %s", comicNum, err)
	}
	defer resp.Body.Close()

	var comicRetrieved model.Comic
	if err := json.NewDecoder(resp.Body).Decode(&comicRetrieved); err != nil {
		return model.Comic{}, fmt.Errorf("failed to GET comic %d: %s", comicNum, err)
	}

	return comicRetrieved, nil
}

func GetLatestComicNumber() (int, error) {
	resp, err := http.Get("https://xkcd.com/info.0.json")
	if err != nil {
		return 0, fmt.Errorf("failed to perform GET request: %s", err)
	}
	defer resp.Body.Close()

	var latest model.LatestComic
	if err := json.NewDecoder(resp.Body).Decode(&latest); err != nil {
		return 0, fmt.Errorf("failed to unmarshall %d: %s", latest.Num, err)
	}

	fmt.Printf("Latest Comic Number is: %d", latest.Num)

	return latest.Num, nil
}

func RetrieveComicsSequentially(start, end int) ([]model.Comic, time.Duration) {
	var comics []model.Comic
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

// #1 attempt - not very efficient
func RetrieveComicsConcurrently(start, end int) ([]model.Comic, time.Duration) {
	var comics []model.Comic
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

// #2 attempt using proposed solution with a channel and a single writer go routine
func RetrieveComicsConcurrently2(start, end int) ([]model.Comic, time.Duration) {
	var comics []model.Comic
	comicCh := make(chan model.Comic, end)
	var getComicsWG sync.WaitGroup
	var writeWG sync.WaitGroup
	startTime := time.Now()

	// writer go routine
	writeWG.Add(1)
	go func() {
		defer writeWG.Done()
		for comic := range comicCh {
			comics = append(comics, comic)
		}
	}()

	for i := start; i <= end; i++ {
		getComicsWG.Add(1)
		go func(num int) {
			defer getComicsWG.Done()
			comic, err := GetComic(num)
			if err != nil {
				fmt.Printf("Error retrieving concurrently comic %d: %v\n", num, err)
				return
			}
			// fmt.Printf("comic number %d is %v\n", i, comic.Title)
			comicCh <- comic
		}(i)
	}

	// first wait to finish retrieving all comics
	getComicsWG.Wait()
	// then close the channel to signal to the writer go routine
	close(comicCh)
	// wait for the writer go routine to finish
	writeWG.Wait()

	duration := time.Since(startTime)
	return comics, duration
}

// #3 attempt using Mutex
func RetrieveComicsConcurrently3(start, end int) ([]model.Comic, time.Duration) {
	var comics []model.Comic
	var wg sync.WaitGroup
	var mutex sync.Mutex
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

			// Lock the mutex before appending to the comics slice and unlock it after appending
			mutex.Lock()
			comics = append(comics, comic)
			mutex.Unlock()
		}(i)
	}

	wg.Wait() // Wait for all goroutines to finish
	duration := time.Since(startTime)
	return comics, duration
}
