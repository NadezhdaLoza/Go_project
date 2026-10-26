package domain

type Order struct {
    OrderHeader    OrderHeader    `json:"order_header"`
    OrderBodyItems []OrderItem    `json:"order_items,omitempty"`
    Performers     []Performer    `json:"performers,omitempty"`
}