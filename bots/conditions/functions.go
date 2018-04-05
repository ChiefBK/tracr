package conditions

import (
	"tracr-cache"
	"tracr-store"
	"tracr/bots/properties"
)

var ConditionFunctions = make(map[string]ConditionFunction)

type ConditionFunction func(props properties.Props, cache *tracr_cache.CacheClient, store tracr_store.Store) bool

func ReturnTrueFunction(props properties.Props, cache *tracr_cache.CacheClient, store tracr_store.Store) ConditionFunction {
	return TrueFunction
}

func TrueFunction(props properties.Props, cache *tracr_cache.CacheClient, store tracr_store.Store) bool {
	return true
}