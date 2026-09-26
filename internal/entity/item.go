package entity

import "context"

// Item — предмет с агрегированными минимальными ценами Skinport:
// отдельно среди tradable- и non-tradable-лотов. Если один из вариантов
// сейчас отсутствует на рынке, соответствующее поле будет nil, а не 0 —
// это осознанное решение, чтобы отличать "цены нет" от "цена нулевая".
type Item struct {
	MarketHashName      string   `json:"market_hash_name"`
	Currency            string   `json:"currency"`
	MinPriceTradable    *float64 `json:"min_price_tradable"`
	MinPriceNonTradable *float64 `json:"min_price_non_tradable"`
}

type ItemsLogic interface {
	GetItems(ctx context.Context, params SkinportParams) ([]Item, error)
}
