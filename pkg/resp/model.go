package resp

import (
	"encoding/json"
	"financial_organizations/app/domain"
	"financial_organizations/pkg/repo"
	"log"
	"net/http"
	"strconv"
)

// Code ...
type Code int

const (
	None = Code(iota)
	BadRequest
	BadRequestHeader
	BadRequestBody
	UserNotFound
	ServerError
	UserIdentifiedError
	UserUnidentifiedError
)

var (
	MsgByCode = map[Code]string{
		None:                  "None",
		BadRequest:            "Плохой запрос",
		BadRequestHeader:      "Неправильный заголвок",
		BadRequestBody:        "Ошибка чтения тела запроса",
		UserNotFound:          "Пользователь не найден",
		ServerError:           "Ошибка на стороне сервера",
		UserUnidentifiedError: "Максимальный баланс для неиндентифированного клиента 10000",
		UserIdentifiedError:   "Максимальный баланс для индентифированного клиента 100000",
	}
)

// CodeResp ...
type CodeResp struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

// ResponseByCode ...
func ResponseByCode(w http.ResponseWriter, code Code, status int) {
	w.WriteHeader(status)

	msgBytes, err := json.Marshal(CodeResp{
		Code:    code,
		Message: MsgByCode[code],
	})
	if err != nil {
		log.Printf("err on response by code, %+v", err)
		return
	}

	_, err = w.Write(msgBytes)
	if err != nil {
		log.Printf("err on write, %+v", err)
	}
}

// SuccessEntered ...
func SuccessEntered(w http.ResponseWriter, info *repo.UserInfo) {
	w.WriteHeader(http.StatusOK)

	infoBytes, err := json.Marshal(CodeResp{
		Code:    None,
		Message: strconv.Itoa(info.Id) + " " + info.Hmac,
	})
	if err != nil {
		log.Printf("err on response by code, %+v", err)
		return
	}

	_, err = w.Write(infoBytes)
	if err != nil {
		log.Printf("err on write, %+v", err)
	}
}

// SendBalance ...
func SendBalance(w http.ResponseWriter, balance domain.Money) {
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(
		struct {
			Balance domain.Money
		}{
			Balance: balance,
		})
	if err != nil {
		log.Printf("err on encode, %+v", err)
	}
}

// SendCountAndAmount ...
func SendCountAndAmount(w http.ResponseWriter, cnt domain.Money, sum domain.Money) {
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(
		struct {
			Count  domain.Money
			Amount domain.Money
		}{
			Count:  cnt,
			Amount: sum,
		})
	if err != nil {
		log.Printf("err on encode, %+v", err)
	}
}
