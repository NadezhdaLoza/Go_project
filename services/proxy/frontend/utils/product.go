package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	productDto "github.com/nadezhda/documents/services/product/dto"
)

func GetProducts() (*productDto.GetProductsResponse, error) {
	url := fmt.Sprintf("%s/products", GetProductsAddress())
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var products productDto.GetProductsResponse
	json.NewDecoder(resp.Body).Decode(&products)
	
	return &products, nil
}

func GetProductById(id int) (*productDto.ProductByIdResponse, error) {
	url := fmt.Sprintf("%s/products/%d", GetProductsAddress(), id)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var product productDto.ProductByIdResponse
	json.NewDecoder(resp.Body).Decode(&product)
	
	return &product, nil
}