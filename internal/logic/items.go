package logic

import (
	"backend-test-app/internal/entity"
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type itemsLogic struct {
	client entity.SkinportClient
	cache  entity.ItemsCache
}

func NewItemsLogic(client entity.SkinportClient, cacheItems entity.ItemsCache) entity.ItemsLogic {
	return &itemsLogic{
		client: client,
		cache:  cacheItems,
	}
}

func (i *itemsLogic) GetItems(ctx context.Context, params entity.SkinportParams) ([]entity.Item, error) {
	key := fmt.Sprintf("%d:%s", params.AppId, params.Currency)

	cached, ok, err := i.cache.Get(key)
	if err != nil {
		return nil, fmt.Errorf("cant get items from cache: %w", err)
	}
	if ok {
		return cached, nil
	}

	var tradable, nonTradable []entity.SkinportRawItem

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		tradable, err = i.client.GetItems(gCtx, params, true)
		if err != nil {
			return fmt.Errorf("cant get tradable items: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		nonTradable, err = i.client.GetItems(gCtx, params, false)
		if err != nil {
			return fmt.Errorf("cant get not tradable items: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	merged := mergeItems(tradable, nonTradable)

	if err = i.cache.Set(key, merged); err != nil {
		return nil, fmt.Errorf("cant set items in cache: %w", err)
	}
	return merged, nil
}

func mergeItems(tradable, nonTradable []entity.SkinportRawItem) []entity.Item {
	byName := make(map[string]*entity.Item)

	for _, raw := range tradable {
		item := getOrCreate(byName, raw)
		item.MinPriceTradable = raw.MinPrice
	}

	for _, raw := range nonTradable {
		item := getOrCreate(byName, raw)
		item.MinPriceNonTradable = raw.MinPrice
	}

	result := make([]entity.Item, 0, len(byName))
	for _, item := range byName {
		result = append(result, *item)
	}
	return result
}

func getOrCreate(m map[string]*entity.Item, raw entity.SkinportRawItem) *entity.Item {
	if item, ok := m[raw.MarketHashName]; ok {
		return item
	}
	item := &entity.Item{
		MarketHashName: raw.MarketHashName,
		Currency:       raw.Currency,
	}
	m[raw.MarketHashName] = item
	return item
}
