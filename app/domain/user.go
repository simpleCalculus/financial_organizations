package domain

// UserAuth ...
type UserAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// User ...
type User struct {
	ID         int    `json:"id"`
	HMAC       string `json:"hmac_sha1"`
	Identified bool   `json:"identified"`
	Balance    Money  `json:"balance"`
}

// Money ...
type Money int
