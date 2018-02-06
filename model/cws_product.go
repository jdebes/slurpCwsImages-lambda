package model

import "time"

type ProductsResponse struct {
	Items []CwsProduct `json:"items"`
}

type CwsProduct struct {
	ProductID   string    `json:"productId"`
	Identifier  string    `json:"identifier"`
	Name        string    `json:"name"`
	Platform    string    `json:"platform"`
	Quantity    int       `json:"quantity"`
	Images      []Image   `json:"images"`
	Regions     []string  `json:"regions"`
	Languages   []string  `json:"languages"`
	Prices      []Price   `json:"prices"`
	Links       []Link    `json:"links"`
	ReleaseDate time.Time `json:"releaseDate"`
}

type Image struct {
	Image  string `json:"image"`
	Format string `json:"format"`
}

type Price struct {
	Value float64 `json:"value"`
	From  float64 `json:"from"`
	To    float64 `json:"to"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}
