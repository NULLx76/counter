package store

type Repository interface {
	Get(key string) int
	Set(key string, value int)
	Increment(key string)
	Decrement(key string)
}
