package main

import (
	"fmt"
	"math/rand"
	"time"
	"xkcd/comics/service"
)

func main() {
    fmt.Println("Hello, World! Welcome to XKCD go bootcamp")

    latestComicNum, err := service.GetLatestComicNumber()
    if err != nil {
        fmt.Printf("Error retrieving latest comic number: %v\n", err)
        return
    }

    // Retrieve comics sequentially
    seqComics, seqDuration := service.RetrieveComicsSequentially(1, latestComicNum)
    fmt.Printf("Number of comics retrieved is: %d", len(seqComics))
    fmt.Printf("Time taken to retrieve all comics sequentially: %v\n", seqDuration)

    // Retrieve comics concurrently
    conComics, conDuration := service.RetrieveComicsConcurrently(1, latestComicNum)
    fmt.Printf("Number of comics retrieved is: %d", len(conComics))
    fmt.Printf("Time taken to retrieve all comics concurrently: %v\n", conDuration)

    // Retrieve information about a random comic
    rand.NewSource(time.Now().UnixNano())
    randomComicNum := rand.Intn(latestComicNum) + 1
    comic, err := service.GetComic(randomComicNum)
    if err != nil {
        fmt.Printf("Error retrieving comic %d: %v\n", randomComicNum, err)
        return
    }
    fmt.Printf("Information about random comic %d:\n", randomComicNum)
    fmt.Printf("Title: %s\n", comic.Title)
    fmt.Printf("Image URL: %s\n", comic.Img)
    fmt.Printf("Transcript: %s\n", comic.Transcript)
}