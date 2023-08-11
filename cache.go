package main

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var codeCache *cache.Cache

func createCache() {
	codeCache = cache.New(cacheTTS*time.Minute, 2*cacheTTS*time.Minute)
}

func insertCache(code string, response Response) {
	codeCache.Set(code, &response, cache.DefaultExpiration)
}

func findCache(code string) (Response, bool) {
	cachedItem, found := codeCache.Get(code)
	if found {
		if response, ok := cachedItem.(*Response); ok {
			return *response, true
		}
	}
	return Response{}, false
}
