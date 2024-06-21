package store

type StoreInterface interface {
	Set(key, value string, secondsTtl int) error
	Get(key string, defaultVal interface{}) (interface{}, error)
	Delete(key string) error
	Has(key string) (bool, error)
}
