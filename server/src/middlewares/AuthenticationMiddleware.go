package middlewares

import(
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

func RefreshUserCache(id int, r *redis.Client) {
	err := r.SAdd(strconv.Itoa(id), "live").Err()
	if err != nil {
		panic(err)
	}

	err = r.Expire(strconv.Itoa(id), time.Second*60).Err()
	if err != nil {
		panic(err)
	}
}