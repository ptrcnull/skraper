package main

type UserData struct {
	Name          string `json:"name"`
	Code          string `json:"code"`
	Avatar        []int  `json:"avatar"`
	Join          string `json:"join"`
	Language      string `json:"language"`
	CreatePrivate bool   `json:"createPrivate"`
}

type WordsEvent struct {
	Id float64 `json:"int"`
	Words []string `json:"words"`
}

type Chans struct {
	Word chan string
	Key chan string
	NewWord chan string
}