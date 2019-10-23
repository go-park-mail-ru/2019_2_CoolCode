package repository

type SessionRepository interface {
	GetID(session string) (uint64, error)
	Contains(session string) bool
	Put(session string, id uint64) error
	Remove(session string) error
}
