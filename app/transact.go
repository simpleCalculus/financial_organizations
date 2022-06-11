package app

import (
	"financial_organizations/pkg/repo"
	"financial_organizations/pkg/resp"
	"log"
	"net/http"
)

func (a *App) Transactions(w http.ResponseWriter, r *http.Request) {
	userInfo, ok := r.Context().Value("user-info").(repo.UserInfo)
	if !ok {
		log.Print("bad value in context")
		return
	}

	cnt, sum, err := a.userRepo.GetCountAndSum(userInfo)
	if err != nil {
		resp.ResponseByCode(w, resp.ServerError, http.StatusInternalServerError)
		return
	}

	resp.SendCountAndAmount(w, cnt, sum)
}
