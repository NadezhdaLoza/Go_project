package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nadezhda/documents/services/account/dto"
	"github.com/nadezhda/documents/services/base"
	//"github.com/levigross/grequests"
)

const (
	BACKEND_ENV = "ACCOUNT_BACKEND_ADDRESS"
)

func GetBackendAddress() string {
	address := base.GetEnv(BACKEND_ENV, "http://localhost:7073")
	return address
}

func GetAccountsFromBackend() (*dto.GetAccountsResponse, error) {
    url := fmt.Sprintf("%s/accounts", GetBackendAddress())
    
    resp, err := http.Get(url) 
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var data dto.GetAccountsResponse
    err = json.NewDecoder(resp.Body).Decode(&data)
    if err != nil {
        return nil, err
    }
    return &data, nil
}