package pubsub

import (
	"sync"
)

type Hub struct{
	mu sync.RWMutex
	topics map[string][]chan string
}

func NewHub() *Hub{
	return &Hub{topics: make(map[string][]chan string)}
}

func (h *Hub) Subscribe(topic string)chan string{
	h.mu.Lock()
	defer h.mu.Unlock()
	ch:=make(chan string,10)
	h.topics[topic]=append(h.topics[topic], ch)
	return ch
}

func (h *Hub) Publish(topic string,message string){
	h.mu.RLock()
	defer h.mu.RUnlock()
	if subs,ok:=h.topics[topic];ok{
		for _,ch:=range subs{
			select{
			case ch <- message:
			default:
				// Handle the case where the channel is full
			}
		}
	}
}

func (h *Hub) Unsubscribe(topic string, targetCh chan string){
	h.mu.Lock()
	defer h.mu.Unlock()
	if subs,ok:=h.topics[topic];ok{
		for i,ch:=range subs{
			if ch == targetCh{
				h.topics[topic]=append(subs[:i],subs[i+1:]...)
				close(ch)
				return
		}
	}
}
}