package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nadezhda/documents/services/proxy/dto"
)

func GetProxyById(id int) (*dto.ProxyByIdResponse, error) {
	url := fmt.Sprintf("%s/proxy/%d", GetProxiesAddress(), id)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var proxy dto.ProxyByIdResponse
	json.NewDecoder(resp.Body).Decode(&proxy)
	
	return &proxy, nil
}