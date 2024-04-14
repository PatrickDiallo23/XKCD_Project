package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Hello, World! Welcome to XKCD go bootcamp")

	mode := flag.String("mode", "sequential", "mode of retrieval (sequential or concurrent)")
	flag.Parse()

	latestComicNum, err := GetLatestComicNumber()
	if err != nil {
		fmt.Printf("Error retrieving latest comic number: %v\n", err)
		return
	}

	switch *mode {
	case "sequential":
		// Retrieve comics sequentially
		seqComics, seqDuration := RetrieveComicsSequentially(1, latestComicNum)
		fmt.Printf("Number of comics retrieved is: %d\n", len(seqComics))
		fmt.Printf("Time taken to retrieve all comics sequentially: %v\n", seqDuration)
	case "concurrent":
		// Retrieve comics concurrently
		conComics, conDuration := RetrieveComicsConcurrently(1, latestComicNum)
		fmt.Printf("Number of comics retrieved is: %d\n", len(conComics))
		fmt.Printf("Time taken to retrieve all comics concurrently: %v\n", conDuration)
	default:
		fmt.Println("Invalid mode. Please use 'sequential' or 'concurrent'.")
		return
	}

	// Retrieve information about a random comic
	rand.NewSource(time.Now().UnixNano())
	randomComicNum := rand.Intn(latestComicNum) + 1
	comic, err := GetComic(randomComicNum)
	if err != nil {
		fmt.Printf("Error retrieving comic %d: %v\n", randomComicNum, err)
		return
	}
	fmt.Printf("Information about random comic %d:\n", randomComicNum)
	fmt.Printf("Title: %s\n", comic.Title)
	fmt.Printf("Image URL: %s\n", comic.Img)
	fmt.Printf("Transcript: %s\n", comic.Transcript)
}
