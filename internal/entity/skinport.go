package entity

import "context"

// SupportedCurrencies — валюты, поддерживаемые Skinport API для параметра currency.
var SupportedCurrencies = map[string]struct{}{
	"AUD": {}, "BRL": {}, "CAD": {}, "CHF": {}, "CNY": {},
	"CZK": {}, "DKK": {}, "EUR": {}, "GBP": {}, "HRK": {},
	"NOK": {}, "PLN": {}, "RUB": {}, "SEK": {}, "TRY": {}, "USD": {},
}

// Строки которые отдает skinport api
type SkinportRawItem struct {
	MarketHashName string   `json:"market_hash_name"`
	Currency       string   `json:"currency"`
	MinPrice       *float64 `json:"min_price"`
	SuggestedPrice *float64 `json:"suggested_price"`
}

type SkinportClient interface {
	GetItems(ctx context.Context, params SkinportParams, tradable bool) ([]SkinportRawItem, error)
}

type SkinportParams struct {
	AppId    int
	Currency string
}

func IsSupportedCurrency(currency string) bool {
	_, ok := SupportedCurrencies[currency]
	return ok
}
