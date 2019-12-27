package main

type UserData struct {
	Name          string `json:"name"`
	Code          string `json:"code"`
	Avatar        []int  `json:"avatar"`
	Join          string `json:"join"`
	Language      string `json:"language"`
	CreatePrivate bool   `json:"createPrivate"`
}
