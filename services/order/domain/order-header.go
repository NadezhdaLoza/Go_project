package domain

type OrderHeader struct {
    ID            int     `json:"id"`
    OrderNumber   string  `json:"order_number"`
    OrderDate     string  `json:"order_date"`
    DepartmentID  int     `json:"department_id"`
    ForemanID     int     `json:"foreman_id"`
    CustomerID    int     `json:"customer_id"`
    Status        string  `json:"status"`
    TotalAmount   float64 `json:"total_amount"`
}