package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/vottundev/vottun-qubic-bridge-go/config"
	"github.com/vottundev/vottun-qubic-bridge-go/utils/log"
)

var (
	ctx              context.Context
	cancel           context.CancelFunc
	redisClient      *redis.Client
	qubicRedisClient *redis.Client
)

func RedisClient() *redis.Client {
	return redisClient
}

func Start(isBridge bool) error {

	var err error

	ctx, cancel = context.WithCancel(context.Background())
	if redisClient, err = newRedisClient(config.Config.Cache.Connections["internal"]); err != nil {
		log.Errorf("%+v", err)
		return err
	}
	if isBridge {
		log.Infof("Internal Redis Client started as bridge. ")
	} else {
		log.Infof("Internal Redis Client started as dispatcher. ")

	}
	if qubicRedisClient, err = newRedisClient(config.Config.Cache.Connections["qubic"]); err != nil {
		log.Errorf("%+v", err)
		return err
	}

	log.Infoln("Qubic Redis Client started")

	if !isBridge {
		go subscribeToQubicEvens(ctx)

		log.Infoln("Qubic redis pubsub subscribed")
	}
	return nil
}

func newRedisClient(info config.CacheInfo) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", info.Host, info.Port),
		Password: "",            // no password set
		DB:       info.Database, // use default DB
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Errorf("REDIS IS NOT AVAILABLE: %+v", err)
		return nil, err
	}

	return redisClient, nil
}

func FlushDB() {
	redisClient.FlushDB(ctx)
}
func Set(key string, value interface{}, expiration time.Duration) error {
	v, _ := json.Marshal(value)
	return redisClient.Set(
		ctx,
		key,
		v,
		expiration,
	).Err()
}

func SetNoExpirationKey(key string, value interface{}) error {

	return Set(key, value, 0)
}
func Get(key string, value interface{}) error {
	v, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return err
	} else if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(v), &value)
	return err

}

func Delete(key string) error {
	result := redisClient.Del(ctx, key)
	return result.Err()
}

func subscribeToQubicEvens(ctx context.Context) {
	if qubicRedisClient != nil {
		// Subscribe to the specified channel
		sub := qubicRedisClient.Subscribe(ctx, config.Config.Cache.QubicEventsChannel)

		// Ensure the subscription is closed when done
		defer sub.Close()

		// Get the channel to receive messages
		channel := sub.Channel()

		// Listen for messages in a loop until context is done
		for {
			select {
			case <-ctx.Done():
				// If the context is done, exit the loop
				log.Infoln("Context done, stopping subscription.")
				return
			case msg := <-channel:
				// Handle incoming messages
				log.Tracef("Received message: %s from channel: %s\n", msg.Payload, msg.Channel)
			}
		}
	}
}

func StopRedisClients() {
	cancel()
	redisClient.Close()
	qubicRedisClient.Close()
}
