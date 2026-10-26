package utils

import (
    "encoding/json"
    "fmt"
    "net/http"
    
    workdto "github.com/nadezhda/documents/services/work/dto"  // Добавляем псевдоним
)

func GetWorks() (*workdto.GetWorksResponse, error) {
    url := fmt.Sprintf("%s/works", GetWorksAddress())
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var works workdto.GetWorksResponse
    if err := json.NewDecoder(resp.Body).Decode(&works); err != nil {
        return nil, err
    }
    
    return &works, nil
}

func GetWorkById(id int) (*workdto.WorkByIdResponse, error) {
    url := fmt.Sprintf("%s/works/%d", GetWorksAddress(), id)
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var work workdto.WorkByIdResponse
    if err := json.NewDecoder(resp.Body).Decode(&work); err != nil {
        return nil, err
    }
    
    return &work, nil
}