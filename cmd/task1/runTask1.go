package task1

import (
	"fmt"
	"math/rand"
	"time"
	"xkcdcomics/cmd"
)

func RunTask1(mode *string) {
	latestComicNum, err := cmd.GetLatestComicNumber()
	if err != nil {
		fmt.Printf("Error retrieving latest comic number: %v\n", err)
		return
	}

	switch *mode {
	case "sequential":
		// Retrieve comics sequentially
		seqComics, seqDuration := cmd.RetrieveComicsSequentially(1, latestComicNum)
		fmt.Printf("Number of comics retrieved is: %d\n", len(seqComics))
		fmt.Printf("Time taken to retrieve all comics sequentially: %v\n", seqDuration)
	case "concurrent":
		// Retrieve comics concurrently
		conComics1, conDuration1 := cmd.RetrieveComicsConcurrently(1, latestComicNum)
		fmt.Println("#1 attempt of retriving comics concurrently")
		fmt.Printf("Number of comics retrieved is: %d\n", len(conComics1))
		fmt.Printf("Time taken to retrieve all comics concurrently: %v\n", conDuration1)
		fmt.Println("-----------------------------------")
		fmt.Println("#2 attempt of retriving comics concurrently with channel and a single write go routine")
		conComics2, conDuration2 := cmd.RetrieveComicsConcurrently2(1, latestComicNum)
		fmt.Printf("Number of comics retrieved is: %d\n", len(conComics2))
		fmt.Printf("Time taken to retrieve all comics concurrently: %v\n", conDuration2)
		fmt.Println("-----------------------------------")
		fmt.Println("#3 attempt of retriving comics concurrently with mutex")
		conComics3, conDuration3 := cmd.RetrieveComicsConcurrently3(1, latestComicNum)
		fmt.Printf("Number of comics retrieved is: %d\n", len(conComics3))
		fmt.Printf("Time taken to retrieve all comics concurrently: %v\n", conDuration3)
	default:
		fmt.Println("Invalid mode. Please use 'sequential' or 'concurrent'.")
		return
	}

	// Retrieve information about a random comic
	rand.NewSource(time.Now().UnixNano())
	randomComicNum := rand.Intn(latestComicNum) + 1
	comic, err := cmd.GetComic(randomComicNum)
	if err != nil {
		fmt.Printf("Error retrieving comic %d: %v\n", randomComicNum, err)
		return
	}
	fmt.Printf("Information about random comic %d:\n", randomComicNum)
	fmt.Printf("Title: %s\n", comic.Title)
	fmt.Printf("Image URL: %s\n", comic.Img)
	fmt.Printf("Transcript: %s\n", comic.Transcript)
}
