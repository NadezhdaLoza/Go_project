package domain

type Work struct {
    ID   int     `json:"id"`
    Name string  `json:"name"`
    Rate float64 `json:"rate"`
    Unit string  `json:"unit"`
}