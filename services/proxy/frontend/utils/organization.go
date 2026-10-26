package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	organizationDto "github.com/nadezhda/documents/services/organization/dto"
)

func GetOrganizationById(id int) (*organizationDto.OrganizationByIdResponse, error) {
	url := fmt.Sprintf("%s/organizations/%d", GetOrganizationsAddress(), id)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var organization organizationDto.OrganizationByIdResponse
	json.NewDecoder(resp.Body).Decode(&organization)
	
	return &organization, nil
}

func GetOrganizations() (organizationDto.GetOrganizationsResponse, error) {
	url := fmt.Sprintf("%s/organizations", GetOrganizationsAddress())
	
	resp, err := http.Get(url)
	if err != nil {
		return organizationDto.GetOrganizationsResponse{}, err
	}
	defer resp.Body.Close()
	
	var organizations organizationDto.GetOrganizationsResponse
	json.NewDecoder(resp.Body).Decode(&organizations)
	
	return organizations, nil
}