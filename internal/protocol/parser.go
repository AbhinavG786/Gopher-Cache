package parser

import (
	"strings"
	"github.com/AbhinavG786/Gopher-Cache.git/internal/engine"
	"time"
)

type StoreInterface interface{
	Get(key string) (engine.Cache,bool)
	Set(cache engine.Cache)
	Delete(key string)
}

func Process(s StoreInterface,input string) string{
	parts:=strings.Fields(strings.TrimSpace(input))
	if len(parts)==0{
		return "Empty command"
	}
	command:=strings.ToUpper(parts[0])
	switch command{
	case "GET":
		if len(parts)<2{
			return "Usage: GET <key>"
		}
		val,ok:=s.Get(parts[1])
		if !ok{
			return "Key not found"
		}
		return val.Value
	case "SET":
		var key, value, ttlPart string

		firstQuote := strings.Index(input, `"`)
		lastQuote := strings.LastIndex(input, `"`)

		if firstQuote != -1 && lastQuote != -1 && firstQuote != lastQuote {

			prefix := strings.Fields(input[:firstQuote])
			if len(prefix) < 2 {
				return "Usage: SET <key> \"<value>\" [TTL]"
			}
			key = prefix[1]

			value = input[firstQuote+1 : lastQuote]

			suffix := strings.Fields(input[lastQuote+1:])
			if len(suffix) > 0 {
				ttlPart = suffix[0]
			}
		} else {

			if len(parts) < 3 {
				return "Usage: SET <key> <value> [TTL in seconds]"
			}
			key = parts[1]
			value = parts[2]
			if len(parts) >= 4 {
				ttlPart = parts[3]
			}
		}

		cache := engine.Cache{Key: key, Value: value}
		if ttlPart != "" {
			ttl, err := time.ParseDuration(ttlPart + "s")
			if err != nil {
				return "Invalid TTL format"
			}
			cache.TTL = time.Now().Add(ttl)
		}

		s.Set(cache)
		return "Key set"
	case "DEL":
		if len(parts)<2{
			return "Usage: DEL <key>"
		}
		s.Delete(parts[1])
		return "Key deleted"
	case "QUIT":
		return "Goodbye!"
	default:
		return "Invalid command"
	}
}