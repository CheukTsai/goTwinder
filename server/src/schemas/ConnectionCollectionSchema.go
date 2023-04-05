package schemas

import(
	"goTwinder/src/tools"
	"database/sql"
	"github.com/go-redis/redis"
)

type ConnectionCollection struct {
	RMQChannelPool *tools.ChannelPool
	MySqlDatabase *sql.DB
	RedisClient *redis.Client
}