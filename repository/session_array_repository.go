package repository

import "errors"

type SessionArrayRepository struct {
	Sessions map[string]uint64
}

func (s SessionArrayRepository) Remove(session string) {
	delete(s.Sessions, session)
}

func (s SessionArrayRepository) GetID(session string) (uint64, error) {
	if val,ok:=s.Sessions[session];ok{
		return uint64(val),nil
	}else {
		return 0,errors.New("No such cookie")
	}
}

func (s SessionArrayRepository) Contains(session string) bool {
	if _,ok:=s.Sessions[session];ok{
		return ok
	}
	return false
}

func (s SessionArrayRepository) Put(session string, id uint64) error {
	s.Sessions[session]=id
	return nil
}

func NewSessionArrayRepository() SessionRepository{
	return SessionArrayRepository{Sessions:make(map[string]uint64,0)}
}


