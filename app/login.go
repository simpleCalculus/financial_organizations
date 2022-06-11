package app

import (
	"encoding/json"
	"errors"
	"financial_organizations/app/domain"
	"financial_organizations/pkg/repo"
	"financial_organizations/pkg/resp"
	"net/http"
)

// Login ...
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	var enter domain.UserAuth

	err := json.NewDecoder(r.Body).Decode(&enter)
	if err != nil {
		resp.ResponseByCode(w, resp.BadRequest, http.StatusBadRequest)
		return
	}

	userInf, err := a.userRepo.GetUserByLogPass(enter.Login, enter.Password)
	if err != nil {
		if errors.Is(err, repo.ErrClientNotFound) {
			resp.ResponseByCode(w, resp.UserNotFound, http.StatusNotFound)
			return
		}

		resp.ResponseByCode(w, resp.ServerError, http.StatusInternalServerError)
		return
	}

	resp.SuccessEntered(w, userInf)
}
