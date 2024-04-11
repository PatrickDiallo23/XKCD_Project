package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"xkcd/comics/model"
)

const xkcdURL = "https://xkcd.com/%d/info.0.json"

func GetComic(comicNum int) (*model.Comic, error) {
    resp, err := http.Get(fmt.Sprintf(xkcdURL, comicNum))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var comic model.Comic
    if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
        return nil, err
    }

    return &comic, nil
}

func GetLatestComicNumber() (int, error) {
    resp, err := http.Get("https://xkcd.com/info.0.json")
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    var latest model.LatestComic
    if err := json.NewDecoder(resp.Body).Decode(&latest); err != nil {
        return 0, err
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
        comics = append(comics, *comic)
    }
    duration := time.Since(startTime)
    return comics, duration
}

func RetrieveComicsConcurrently(start, end int) ([]model.Comic, time.Duration) {
    var comics []model.Comic
    var wg sync.WaitGroup
    startTime := time.Now()
    // We could also use wg outside the for scope because we know the number of comics
    //todo: check if this is more efficient
    // wg.Add(getLatestComicNumber())
    // defer wg.Done()

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
            comics = append(comics, *comic)
        }(i)
    }

    wg.Wait()
    duration := time.Since(startTime)
    return comics, duration
}