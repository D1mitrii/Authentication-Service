package models

type Token struct {
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}
