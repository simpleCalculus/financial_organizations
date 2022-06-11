package app

import (
	"financial_organizations/pkg/repo"
	"financial_organizations/pkg/resp"
	"log"
	"net/http"
)

// WalletBalance ...
func (a *App) WalletBalance(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value("user-info").(repo.UserInfo)
	if !ok {
		log.Print("bad value in context")
		return
	}

	balance, err := a.userRepo.GetBalance(userInfo)
	if err != nil {
		resp.ResponseByCode(w, resp.ServerError, http.StatusInternalServerError)
		return
	}

	resp.SendBalance(w, balance)
}
