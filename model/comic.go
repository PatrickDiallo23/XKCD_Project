package model

type Comic struct {
	Num        int    `json:"num"`
	SafeTitle  string `json:"safeTitle"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
}

type LatestComic struct {
	Num int `json:"num"`
}
