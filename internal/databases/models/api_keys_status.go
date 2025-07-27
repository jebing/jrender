package models

//go:generate enumer -type=ApiKeysStatus -json -output=api_keys_status_enumer.go

type ApiKeysStatus int8

const (
	ApiKeysStatusActive ApiKeysStatus = iota + 1
	ApiKeysStatusRevoked
)
