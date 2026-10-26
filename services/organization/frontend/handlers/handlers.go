package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nadezhda/documents/services/base"
	"github.com/nadezhda/documents/services/organization/domain"
	"github.com/nadezhda/documents/services/organization/dto"
	"github.com/nadezhda/documents/services/organization/frontend/utils"

	accountsDomain "github.com/nadezhda/documents/services/account/domain"
	accountsDto "github.com/nadezhda/documents/services/account/dto"
)

func OrganizationsHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/organizations", utils.GetOrganizationsAddress())
	resp, _ := http.Get(url)
	var data dto.GetOrganizationsResponse
	json.NewDecoder(resp.Body).Decode(&data)
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusInternalServerError)
		return
	}
	tmpl, _ := template.ParseFiles("../static/html/organizations.html")
	tmpl.Execute(w, data.Organizations)
}

func DeleteOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	url := fmt.Sprintf("%s/organizations", utils.GetOrganizationsAddress())
	body, _ := json.Marshal(map[string]interface{}{"id": id})
	req, _ := http.NewRequest("DELETE", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	var deleteResponse dto.DeleteOrganizationResponse
	json.NewDecoder(resp.Body).Decode(&deleteResponse)
	if deleteResponse.Err != "" {
		http.Error(w, deleteResponse.Err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/organization/organizations", http.StatusSeeOther)
}

func CreateOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	accountsUrl := fmt.Sprintf("%s/accounts", utils.GetAccountsAddress())
	resp, _ := http.Get(accountsUrl)
	var data accountsDto.GetAccountsResponse
	json.NewDecoder(resp.Body).Decode(&data)
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusInternalServerError)
	}
	tmpl, _ := template.ParseFiles("../static/html/create-organization.html")
	tmpl.Execute(w, data.Accounts)
}

func CreateOrganizationPostHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	address := r.FormValue("address")
	accountIdStr := r.FormValue("account_id")
	accountId, accountErr := strconv.Atoi(accountIdStr)
	chief := r.FormValue("chief")
	financialChief := r.FormValue("financial_chief")
	if name == "" ||
		address == "" ||
		accountErr != nil ||
		chief == "" ||
		financialChief == "" {
		http.Error(w, base.FIELDS_VALIDATION_ERROR, http.StatusInternalServerError)
		return
	}
	url := fmt.Sprintf("%s/organizations", utils.GetOrganizationsAddress())
	request := dto.CreateOrganizationRequest{
		Name:           name,
		Address:        address,
		AccountId:      accountId,
		Chief:          chief,
		FinancialChief: financialChief,
	}
	body, _ := json.Marshal(request)
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(body))
	var data dto.CreateOrganizationResponse
	json.NewDecoder(resp.Body).Decode(&data)
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, "/organization/organizations", http.StatusSeeOther)
	}
}

func OrganizationDetailHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	organizationUrl := fmt.Sprintf("%s/organizations/%d", utils.GetOrganizationsAddress(), id)
	resp, _ := http.Get(organizationUrl)
	var organizationResp dto.OrganizationByIdResponse
	json.NewDecoder(resp.Body).Decode(&organizationResp)
	if organizationResp.Err != "" {
		http.Error(w, organizationResp.Err, http.StatusInternalServerError)
		return
	}

	accountsUrl := fmt.Sprintf("%s/accounts/%d", utils.GetAccountsAddress(), organizationResp.Organization.AccountId)
	resp, _ = http.Get(accountsUrl)
	var accountResp accountsDto.AccountByIdResponse
	json.NewDecoder(resp.Body).Decode(&accountResp)
	if accountResp.Err != "" {
		http.Error(w, accountResp.Err, http.StatusInternalServerError)
		return
	}

	type organizationDetail struct {
		Name               string
		Address            string
		BankName           string
		BankIdentityNumber string
		Chief              string
		FinancialChief     string
	}
	tmpl, _ := template.ParseFiles("../static/html/organization-details.html")
	data := organizationDetail{
		Name:               organizationResp.Organization.Name,
		Address:            organizationResp.Organization.Address,
		BankName:           accountResp.Account.BankName,
		BankIdentityNumber: accountResp.Account.BankIdentityNumber,
		Chief:              organizationResp.Organization.Chief,
		FinancialChief:     organizationResp.Organization.FinancialChief,
	}
	tmpl.Execute(w, data)
}

func UpdateOrganizationGetHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	url := fmt.Sprintf("%s/organizations/%d", utils.GetOrganizationsAddress(), id)
	orgResp, _ := http.Get(url)
	var organization dto.OrganizationByIdResponse
	json.NewDecoder(orgResp.Body).Decode(&organization)

	accountsUrl := fmt.Sprintf("%s/accounts", utils.GetAccountsAddress())
	accResp, _ := http.Get(accountsUrl)
	var accounts accountsDto.GetAccountsResponse
	json.NewDecoder(accResp.Body).Decode(&accounts)
	if accounts.Err != "" {
		http.Error(w, accounts.Err, http.StatusInternalServerError)
	}

	type templateData struct {
		Organization domain.Organization
		Accounts     []accountsDomain.Account
	}
	data := templateData{
		Organization: *organization.Organization,
		Accounts:     accounts.Accounts,
	}

	tmpl, _ := template.ParseFiles("../static/html/update-organization.html")
	tmpl.Execute(w, data)
}

func UpdateOrganizationPostHandler(w http.ResponseWriter, r *http.Request) {
	id, idErr := strconv.Atoi(r.FormValue("id"))
	name := r.FormValue("name")
	address := r.FormValue("address")
	accountIdStr := r.FormValue("account_id")
	accountId, accountErr := strconv.Atoi(accountIdStr)
	chief := r.FormValue("chief")
	financialChief := r.FormValue("financial_chief")
	if name == "" ||
		address == "" ||
		accountErr != nil ||
		chief == "" ||
		financialChief == "" ||
		idErr != nil {

		http.Error(w, base.FIELDS_VALIDATION_ERROR, http.StatusInternalServerError)
		return
	}
	url := fmt.Sprintf("%s/organizations", utils.GetOrganizationsAddress())
	organization := domain.Organization{
		Id:             id,
		Name:           name,
		Address:        address,
		AccountId:      accountId,
		Chief:          chief,
		FinancialChief: financialChief,
	}
	request := dto.UpdateOrganizationRequest{
		Organization: organization,
	}
	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	var data dto.UpdateOrganizationResponse
	json.NewDecoder(resp.Body).Decode(&data)
	if data.Err != "" {
		http.Error(w, data.Err, http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, "/organization/organizations", http.StatusSeeOther)
	}
}