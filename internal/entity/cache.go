package entity

// Возвращаем error на случай если будем использовать хранилища по типу Redis
type ItemsCache interface {
	Get(key string) ([]Item, bool, error)
	Set(key string, items []Item) error
}
