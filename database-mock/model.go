package main

type GetModel struct {
	ID string `json:"id"`
}

type CreateModel struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Amount int64  `json:"amount"`
}
