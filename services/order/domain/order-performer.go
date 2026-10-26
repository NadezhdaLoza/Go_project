package domain

type Performer struct {
    ID         int     `json:"id"`
    OrderID    int     `json:"order_id"`
    EmployeeID int     `json:"employee_id"`
    KTU        float64 `json:"ktu"`
    Amount     float64 `json:"amount"`
}