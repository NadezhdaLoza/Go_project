package utils

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/nadezhda/documents/services/order/dto"
)

func GetOrderById(id int) (*dto.OrderByIdResponse, error) {
    url := fmt.Sprintf("%s/orders/%d", GetOrdersAddress(), id)
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var order dto.OrderByIdResponse
    if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
        return nil, err
    }
    
    return &order, nil
}