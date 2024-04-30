package models

type JWTPair struct {
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}
