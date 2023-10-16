package models

import "time"

type ServiceCheck struct {
	Name     string    `json:"name"`
	Url      string    `json:"url"`
	CreateAt time.Time `json:"createAt"`
	Status   string    `json:"status"`
}

type ApiWithName struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
