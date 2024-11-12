package cache

const (
	cache_key_bridges_list string = "CK_BRIDGES_LIST"
	cache_key_bridges_set  string = "CK_BRIDGES_SET"
)

// func AddBridgeToSet(bridge string) error {

// 	result := redisClient.SAdd(ctx, cache_key_bridges_set,bridgeTopic)

// 	if result.Err()!=nil {
// 		log.Errorf("Error adding new bridge with id")
// 		return result.Err()
// 	}
// }
