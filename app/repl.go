package app

import (
	"encoding/json"
	"financial_organizations/app/domain"
	"financial_organizations/pkg/repo"
	"financial_organizations/pkg/resp"
	"log"
	"net/http"
)

// ReplBody ...
type ReplBody struct {
	Amount domain.Money `json:"amount"`
}

// WalletReplenishment ...
func (a *App) WalletReplenishment(w http.ResponseWriter, r *http.Request) {
	var data ReplBody
	userInfo, ok := r.Context().Value("user-info").(repo.UserInfo)
	if !ok {
		log.Print("bad value in context")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		resp.ResponseByCode(w, resp.BadRequest, http.StatusBadRequest)
		return
	}

	oldBalance, err := a.userRepo.GetBalance(userInfo)
	if err != nil {
		resp.ResponseByCode(w, resp.ServerError, http.StatusInternalServerError)
		return
	}

	isIdentified, err := a.userRepo.GetIdentified(userInfo)
	if err != nil {
		resp.ResponseByCode(w, resp.ServerError, http.StatusInternalServerError)
		return
	}

	newBalance := oldBalance + data.Amount

	switch {
	case newBalance <= 10_000:
	case newBalance > 10_000 && !isIdentified:
		resp.ResponseByCode(w, resp.UserUnidentifiedError, http.StatusConflict)
		return
	case newBalance > 10_000 && newBalance <= 100_000:
	default:
		resp.ResponseByCode(w, resp.UserIdentifiedError, http.StatusConflict)
		return
	}

	err = a.userRepo.UpdateBalance(userInfo, newBalance)
	if err != nil {
		resp.ResponseByCode(w, resp.ServerError, http.StatusInternalServerError)
	}

	err = a.userRepo.AddToTransactHistory(userInfo.Id, data.Amount)
	if err != nil {
		log.Print(err)
	}

	resp.ResponseByCode(w, resp.None, http.StatusOK)
}
