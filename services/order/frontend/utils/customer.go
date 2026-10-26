package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	customerDto "github.com/nadezhda/documents/services/customer/dto"
)

func GetCustomers() (customerDto.GetCustomersResponse, error) {
	url := fmt.Sprintf("%s/customers", GetCustomersAddress())
	resp, err := http.Get(url)
	if err != nil {
		return customerDto.GetCustomersResponse{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var customers customerDto.GetCustomersResponse
	json.Unmarshal(body, &customers)
	return customers, nil
}

func GetCustomerById(id int) (customerDto.CustomerByIdResponse, error) {
	url := fmt.Sprintf("%s/customers/%d", GetCustomersAddress(), id)
	resp, err := http.Get(url)
	if err != nil {
		return customerDto.CustomerByIdResponse{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var customer customerDto.CustomerByIdResponse
	json.Unmarshal(body, &customer)
	return customer, nil
}