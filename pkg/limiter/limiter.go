package limiter

import (
	"github.com/gin-gonic/gin"
	limiterlib "github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"gohub/pkg/config"
	"gohub/pkg/logger"
	"gohub/pkg/redis"
	"strings"
)

// GetKeyIP 获取 Limitor 的 Key，IP
func GetKeyIP(c *gin.Context) string {
	return c.ClientIP()
}

// GetKeyRouteWithIP Limitor 的 Key，路由+IP，针对单个路由做限流
func GetKeyRouteWithIP(c *gin.Context) string {
	return routeToKeyString(c.FullPath()) + c.ClientIP()
}

// CheckRate 检测请求是否超额
func CheckRate(c *gin.Context, key string, formatted string) (limiterlib.Context, error) {
	var context limiterlib.Context
	rate, err := limiterlib.NewRateFromFormatted(formatted)
	if err != nil {
		logger.LogIf(err)
		return context, err
	}

	store, err := sredis.NewStoreWithOptions(redis.Redis.Client, limiterlib.StoreOptions{
		Prefix: config.GetString("app.name") + ":limiter",
	})
	if err != nil {
		logger.LogIf(err)
		return context, err
	}

	limiterObj := limiterlib.New(store, rate)

	if c.GetBool("limiter-once") {
		return limiterObj.Peek(c, key)
	} else {
		c.Set("limiter-once", true)
		return limiterObj.Get(c, key)
	}
}

func routeToKeyString(routeName string) string {
	routeName = strings.ReplaceAll(routeName, "/", "-")
	routeName = strings.ReplaceAll(routeName, ":", "_")
	return routeName
}
