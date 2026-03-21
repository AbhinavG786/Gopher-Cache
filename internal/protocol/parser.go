package parser

import (
	"strings"
	"time"

	store "github.com/AbhinavG786/Gopher-Cache.git/internal/engine"
	"github.com/AbhinavG786/Gopher-Cache.git/internal/pubsub"
)

type StoreInterface interface{
	Get(key string) (store.Cache,bool)
	Set(cache store.Cache)
	Delete(key string)
}

func Process(s StoreInterface,h *pubsub.Hub,input string) (string, chan string){
	parts:=strings.Fields(strings.TrimSpace(input))
	if len(parts)==0{
		return "Empty command", nil
	}
	command:=strings.ToUpper(parts[0])
	switch command{
	case "GET":
		if len(parts)<2{
			return "Usage: GET <key>", nil
		}
		val,ok:=s.Get(parts[1])
		if !ok{
			return "Key not found", nil
		}
		return val.Value, nil
	case "SET":
		var key, value, ttlPart string

		firstQuote := strings.Index(input, `"`)
		lastQuote := strings.LastIndex(input, `"`)

		if firstQuote != -1 && lastQuote != -1 && firstQuote != lastQuote {

			prefix := strings.Fields(input[:firstQuote])
			if len(prefix) < 2 {
				return "Usage: SET <key> \"<value>\" [TTL]", nil
			}
			key = prefix[1]

			value = input[firstQuote+1 : lastQuote]

			suffix := strings.Fields(input[lastQuote+1:])
			if len(suffix) > 0 {
				ttlPart = suffix[0]
			}
		} else {

			if len(parts) < 3 {
				return "Usage: SET <key> <value> [TTL in seconds]", nil
			}
			key = parts[1]
			value = parts[2]
			if len(parts) >= 4 {
				ttlPart = parts[3]
			}
		}

		cache := store.Cache{Key: key, Value: value}
		if ttlPart != "" {
			ttl, err := time.ParseDuration(ttlPart + "s")
			if err != nil {
				return "Invalid TTL format", nil
			}
			cache.TTL = time.Now().Add(ttl)
		}

		s.Set(cache)
		return "Key set", nil
	case "DEL":
		if len(parts)<2{
			return "Usage: DEL <key>", nil
		}
		s.Delete(parts[1])
		return "Key deleted", nil
	case "SUBSCRIBE":
		if len(parts)<2{
			return "Usage: SUBSCRIBE <topic>", nil
		}
		topic:=parts[1]
		userChan:=h.Subscribe(topic)
		return "Subscribed", userChan
	case "PUBLISH":
		if len(parts)<3{
			return "Usage: PUBLISH <topic> <message>", nil
		}
		topic:=parts[1]
		message:=strings.Join(parts[2:]," ")
		h.Publish(topic,message)
		return "Message published",nil
	case "UNSUBSCRIBE":
		if len(parts)<2{
			return "Usage: UNSUBSCRIBE <topic>", nil
		}
		topic:=parts[1]
		userChan:=make(chan string)
		h.Unsubscribe(topic,userChan)
		return "Unsubscribed from " + topic, nil
	case "QUIT":
		return "Goodbye!", nil
	default:
		return "Invalid command", nil
	}
}