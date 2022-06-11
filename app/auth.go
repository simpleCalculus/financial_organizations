package app

import (
	"financial_organizations/pkg/resp"
	"net/http"
)

// Authentication ...
func (a *App) Authentication(w http.ResponseWriter, _ *http.Request) {
	resp.ResponseByCode(w, resp.None, http.StatusOK)
}
