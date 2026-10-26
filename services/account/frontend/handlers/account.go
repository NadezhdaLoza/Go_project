package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/nadezhda/documents/services/account/domain"
	"github.com/nadezhda/documents/services/account/dto"
	"github.com/nadezhda/documents/services/account/frontend/utils"
	"github.com/nadezhda/documents/services/base"
)

func AccountsHandler(w http.ResponseWriter, r *http.Request) {
    // Получаем данные из бэкенда
    accountsResp, err := utils.GetAccountsFromBackend()
    if err != nil {
        base.LogError(err)
        http.Error(w, "Failed to fetch accounts from backend", http.StatusInternalServerError)
        return
    }
    
    if accountsResp.Err != "" {
        http.Error(w, accountsResp.Err, http.StatusInternalServerError)
        return
    }

    // Парсим HTML-шаблон
    accountsPage, err := template.ParseFiles("static/html/accounts.html")
    if err != nil {
        // Попробуем альтернативный путь
        accountsPage, err = template.ParseFiles("../static/html/accounts.html")
        if err != nil {
            base.LogError(fmt.Errorf("failed to parse template: %v", err))
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
    }
    
    // Выполняем шаблон с данными
    err = accountsPage.Execute(w, accountsResp.Accounts)
    if err != nil {
        base.LogError(fmt.Errorf("failed to execute template: %v", err))
        http.Error(w, "Internal server error", http.StatusInternalServerError)
    }
}

func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	accountsUrl := fmt.Sprintf("%s/accounts", utils.GetBackendAddress())
	
	body, _ := json.Marshal(map[string]interface{}{"id": id})
	req, _ := http.NewRequest("DELETE", accountsUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	
	var deleteResponse dto.DeleteAccountResponse
	json.NewDecoder(resp.Body).Decode(&deleteResponse)
	if deleteResponse.Err != "" {
		http.Error(w, deleteResponse.Err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/account/accounts", http.StatusSeeOther)
}

func CreateAccountGetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../static/html/create-account.html")
}

func CreateAccountPostHandler(w http.ResponseWriter, r *http.Request) {
	bankName := r.FormValue("bank_name")
	bankIdentityNumber := r.FormValue("bin")
	account := r.FormValue("account")
	if bankName == "" || len(bankIdentityNumber) != 9 {
		http.Error(w, base.FIELDS_VALIDATION_ERROR, http.StatusUnprocessableEntity)
		return
	}
	request := dto.CreateAccountRequest{
		Account:            account,
		BankName:           bankName,
		BankIdentityNumber: bankIdentityNumber,
	}
	accountsUrl := fmt.Sprintf("%s/accounts", utils.GetBackendAddress())
	
	body, _ := json.Marshal(request)
	resp, _ := http.Post(accountsUrl, "application/json", bytes.NewBuffer(body))
	
	var data dto.CreateAccountResponse
	json.NewDecoder(resp.Body).Decode(&data)
	if data.Err == "" {
		http.Redirect(w, r, "/account/accounts", http.StatusSeeOther)
	} else {
		http.Error(w, data.Err, http.StatusInternalServerError)
	}
}

func UpdateAccountGetHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	url := fmt.Sprintf("%s/accounts/%d", utils.GetBackendAddress(), id)
	
	resp, _ := http.Get(url)
	
	var account dto.AccountByIdResponse
	json.NewDecoder(resp.Body).Decode(&account)
	if account.Err != "" {
		http.Error(w, account.Err, http.StatusBadRequest)
		return
	}
	updateAccountPage, _ := template.ParseFiles("../static/html/update-account.html")
	updateAccountPage.Execute(w, account.Account)
}

func UpdateAccountPostHandler(w http.ResponseWriter, r *http.Request) {
	accountId, _ := strconv.Atoi(r.FormValue("id"))
	bankName := r.FormValue("bank_name")
	bankIdentityCode := r.FormValue("bin")
	acc := r.FormValue("account")
	account := domain.Account{
		Id:                 accountId,
		Account:            acc,
		BankName:           bankName,
		BankIdentityNumber: bankIdentityCode,
	}

	accountsUrl := fmt.Sprintf("%s/accounts", utils.GetBackendAddress())
	request := dto.UpdateAccountRequest{
		Account: account,
	}
	

	body, _ := json.Marshal(request)
	req, _ := http.NewRequest("PUT", accountsUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, _ := client.Do(req)

	var updateResponse dto.UpdateAccountResponse
	json.NewDecoder(response.Body).Decode(&updateResponse)  
	if updateResponse.Err != "" {
		http.Error(w, updateResponse.Err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/account/accounts", http.StatusSeeOther)
}