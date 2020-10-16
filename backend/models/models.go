package models

type Campaign struct {
    Id       int64 `json:"id"`
	Name     string `json:"name"`
	Status   string  `json:"status"`
	Type 	 string  `json:"type"`
	Budget   int64   `json:"budget"`
	Created_at string `json:"created_at"`
}