package main

type comic struct {
	Num        int    `json:"num"`
	SafeTitle  string `json:"safeTitle"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
}

type latestComic struct {
	Num int `json:"num"`
}
