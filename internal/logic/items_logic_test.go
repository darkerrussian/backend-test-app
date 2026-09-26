package logic

import (
	"backend-test-app/internal/entity"
	"context"
	"errors"
	"testing"
)

// mockClient — ручной мок entity.SkinportClient.
type mockClient struct {
	tradableItems    []entity.SkinportRawItem
	nonTradableItems []entity.SkinportRawItem
	err              error
}

func (m *mockClient) GetItems(_ context.Context, _ entity.SkinportParams, tradable bool) ([]entity.SkinportRawItem, error) {
	if m.err != nil {
		return nil, m.err
	}
	if tradable {
		return m.tradableItems, nil
	}
	return m.nonTradableItems, nil
}

// mockCache — ручной мок entity.ItemsCache, всегда промах, если не задан hit.
type mockCache struct {
	hit    []entity.Item
	hasHit bool
	setErr error
}

func (m *mockCache) Get(_ string) ([]entity.Item, bool, error) {
	return m.hit, m.hasHit, nil
}

func (m *mockCache) Set(_ string, items []entity.Item) error {
	if m.setErr != nil {
		return m.setErr
	}
	m.hit = items
	return nil
}

func float64Ptr(v float64) *float64 { return &v }

func TestService_GetItems_MergesPricesByName(t *testing.T) {
	client := &mockClient{
		tradableItems: []entity.SkinportRawItem{
			{MarketHashName: "AK-47 | Redline", Currency: "EUR", MinPrice: float64Ptr(10.5)},
		},
		nonTradableItems: []entity.SkinportRawItem{
			{MarketHashName: "AK-47 | Redline", Currency: "EUR", MinPrice: float64Ptr(8.0)},
			{MarketHashName: "AWP | Asiimov", Currency: "EUR", MinPrice: float64Ptr(50.0)},
		},
	}
	cache := &mockCache{}
	service := NewItemsLogic(client, cache)

	got, err := service.GetItems(context.Background(), entity.SkinportParams{AppId: 730, Currency: "EUR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}

	var redline *entity.Item
	for i := range got {
		if got[i].MarketHashName == "AK-47 | Redline" {
			redline = &got[i]
		}
	}
	if redline == nil {
		t.Fatal("AK-47 | Redline not found in result")
	}
	if redline.MinPriceTradable == nil || *redline.MinPriceTradable != 10.5 {
		t.Errorf("expected tradable price 10.5, got %v", redline.MinPriceTradable)
	}
	if redline.MinPriceNonTradable == nil || *redline.MinPriceNonTradable != 8.0 {
		t.Errorf("expected non-tradable price 8.0, got %v", redline.MinPriceNonTradable)
	}
}

func TestService_GetItems_ReturnsCachedResultWithoutCallingClient(t *testing.T) {
	cachedItems := []entity.Item{{MarketHashName: "cached item"}}
	client := &mockClient{err: errors.New("client should not be called")}
	cache := &mockCache{hit: cachedItems, hasHit: true}
	service := NewItemsLogic(client, cache)

	got, err := service.GetItems(context.Background(), entity.SkinportParams{AppId: 730, Currency: "EUR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].MarketHashName != "cached item" {
		t.Errorf("expected cached result, got %+v", got)
	}
}

func TestService_GetItems_ReturnsErrorWhenClientFails(t *testing.T) {
	client := &mockClient{err: errors.New("skinport unavailable")}
	cache := &mockCache{}
	service := NewItemsLogic(client, cache)

	_, err := service.GetItems(context.Background(), entity.SkinportParams{AppId: 730, Currency: "EUR"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
