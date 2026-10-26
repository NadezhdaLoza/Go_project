package domain

type OrderItem struct {
    ID          int     `json:"id"`
    OrderID     int     `json:"order_id"`
    WorkID      int     `json:"work_id"`
    Description string  `json:"description"`
    Unit        string  `json:"unit"`
    Quantity    float64 `json:"quantity"`
    Rate        float64 `json:"rate"`
    Amount      float64 `json:"amount"`
}