package store

import (
	"sync"
	"time"
)

type Cache struct{
	Key 	string     `json:"key"`
	Value	string	   `json:"value"`
	TTL   	time.Time  `json:"ttl"`
}

type Store struct{
	mu sync.RWMutex
	data map[string]Cache
}

func New() *Store{
	s:= &Store{data: make(map[string]Cache)}
	go s.evictExpired()
	return s
}

func (s *Store) GetAll() []Cache{
	s.mu.RLock()
	defer s.mu.RUnlock()
	var list []Cache
	for _,v:=range s.data {
		list=append(list,v)
	}
	return list
}

func (s *Store) Get(key string) (Cache,bool){
	s.mu.RLock()
	defer s.mu.RUnlock()
	val,ok:=s.data[key]
	return val,ok
}

func (s *Store) Delete(key string){
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data,key)
}

func(s *Store) evictExpired(){
	ticker:=time.NewTicker(5*time.Second)
	for range ticker.C{
		s.mu.Lock()
		for key,item:= range s.data{
			if !item.TTL.IsZero() && time.Now().After(item.TTL){
				delete(s.data,key)
			}
		}
		s.mu.Unlock()
	}
}