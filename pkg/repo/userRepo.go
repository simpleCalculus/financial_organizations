package repo

import (
	"crypto/hmac"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"errors"
	"financial_organizations/app/domain"
	"fmt"
	"log"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

type UserInfo struct {
	Id   int    `json:"id"`
	Hmac string `json:"hmac"`
}

var (
	ErrClientNotFound = errors.New("user not found")
)

// Authenticate ...
func (r *UserRepo) Authenticate(userId int, hmac string) error {
	sqlStatement := `select * from users where id = $1 and hmac_sha1 = $2`
	if !rowExists(r.db, sqlStatement, userId, hmac) {
		return ErrClientNotFound
	}
	return nil
}

// GetUserByLogPass ...
func (r *UserRepo) GetUserByLogPass(login, pass string) (*UserInfo, error) {
	h := GetSignature(login, pass)

	sqlStatement := `select id from users where hmac_sha1 = $1`
	if !rowExists(r.db, sqlStatement, h) {
		return nil, ErrClientNotFound
	}

	var clientId int
	row := r.db.QueryRow(sqlStatement, h)

	err := row.Err()
	if err != nil {
		return nil, err
	}

	err = row.Scan(&clientId)
	if err != nil {
		return nil, err
	}

	info := UserInfo{
		Id:   clientId,
		Hmac: h,
	}

	return &info, nil
}

// GetBalance ...
func (r *UserRepo) GetBalance(info UserInfo) (domain.Money, error) {
	sqlStatement := `select balance from users where id = $1 and hmac_sha1 = $2`

	var userBalance domain.Money

	row := r.db.QueryRow(sqlStatement, info.Id, info.Hmac)
	err := row.Scan(&userBalance)
	if err != nil {
		return domain.Money(0), err
	}

	return userBalance, nil
}

// GetIdentified ...
func (r *UserRepo) GetIdentified(info UserInfo) (bool, error) {
	sqlStatement := `select identified from users where id = $1 and hmac_sha1 = $2`

	var identified string

	row := r.db.QueryRow(sqlStatement, info.Id, info.Hmac)
	err := row.Scan(&identified)
	if err != nil {
		return false, err
	}

	if identified == "Y" {
		return true, nil
	} else {
		return false, nil
	}
}

// UpdateBalance ...
func (r *UserRepo) UpdateBalance(info UserInfo, newBalance domain.Money) error {
	sqlStatement := `update users set balance = $1 where id = $2 and hmac_sha1 = $3`
	_, err := r.db.Exec(sqlStatement, newBalance, info.Id, info.Hmac)
	if err != nil {
		return err
	}
	return nil
}

// AddToTransactHistory ...
func (r *UserRepo) AddToTransactHistory(clientId int, replenishment domain.Money) error {
	sqlStatement := `insert into transactions(client_id, amount, date) values($1, $2, now())`
	_, err := r.db.Exec(sqlStatement, clientId, replenishment)
	return err
}

// GetCountAndSum ...
func (r *UserRepo) GetCountAndSum(info UserInfo) (domain.Money, domain.Money, error) {
	sqlStatement := `select count(*), sum(amount) from transactions
		where client_id = $1 and EXTRACT(MONTH FROM now()) = EXTRACT(MONTH FROM date) 
		and EXTRACT(YEAR FROM now()) = EXTRACT(YEAR FROM date);
	`

	var cnt domain.Money
	var sum domain.Money

	rows, _ := r.db.Query(sqlStatement, info.Id)
	if rows.Next() {
		err := rows.Scan(&cnt, &sum)
		if err != nil && rows.Next() {
			return 0, 0, err
		}
	}

	return cnt, sum, nil
}

// GetSignature ...
func GetSignature(key string, value string) string {
	keyForSign := []byte(key)
	h := hmac.New(sha1.New, keyForSign)
	h.Write([]byte(value))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// rowExists ...
func rowExists(db *sql.DB, query string, args ...interface{}) bool {
	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("error checking if row exists '%s' %v", args, err)
	}
	return exists
}

// CreateTables ...
func (r *UserRepo) CreateTables() {
	createUsers := `
	create table if not exists users
		(
			id serial primary key,
			hmac_sha1 varchar(100) not null,
			balance int,
			identified char(1) NOT NULL,
	    	unique (hmac_sha1)
		);
`
	_, err := r.db.Exec(createUsers)
	if err != nil {
		panic(err)
	}

	createTransactions := `
	create table if not exists transactions
	(
		id serial primary key,
		client_id integer not null,
		amount int,
		date timestamp NOT NULL
	);
`
	_, err = r.db.Exec(createTransactions)
	if err != nil {
		panic(err)
	}
}

// AddUsers ...
func (r *UserRepo) AddUsers() {
	hmac1 := GetSignature("test", "123")
	hmac2 := GetSignature("login", "qwerty")

	addUser := `
	insert into users(hmac_sha1, balance, identified) 
	values($1, $2, $3)
	ON conflict do nothing
`

	if _, err := r.db.Exec(addUser, hmac1, 0, "Y"); err != nil {
		panic(err)
	}

	if _, err := r.db.Exec(addUser, hmac2, 0, "N"); err != nil {
		panic(err)
	}
}
