package http

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/core/util/xtime"
	"github.com/gotomicro/ego/server/egin"
	"golang.org/x/sync/singleflight"
	"shorturl/pkg/invoker"
	"shorturl/pkg/mysql"
)

var (
	g           = &singleflight.Group{}
	ErrNotFound = fmt.Errorf("code not exist")
)

func ServeHTTP() *egin.Component {
	r := egin.Load("server.http").Build()
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Short URL")
	})
	r.GET("/:code", func(ctx *gin.Context) {
		code := ctx.Param("code")
		originURL, err := getRedirectURLWrapper(ctx.Request.Context(), g, code)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				ctx.JSON(200, "code码不存在")
				return
			}
			elog.Error("系统错误", elog.FieldErr(err), elog.FieldValueAny(code))
			ctx.JSON(500, "系统错误")
			return
		}
		ctx.Redirect(302, originURL)
	})
	return r
}

func getRedirectURLWrapper(ctx context.Context, sg *singleflight.Group, code string) (string, error) {
	v, err, _ := sg.Do(code, func() (interface{}, error) {
		return getRedirectURL(ctx, code)
	})
	return v.(string), err
}

func getRedirectURL(ctx context.Context, code string) (string, error) {
	originURL, err := invoker.Redis.Get(ctx, econf.GetString("shorturl.redisPrefix")+":"+code)
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("find redis failed, err: %w", err)
	}
	// 说明这个时候，是可以获取到数据
	if !errors.Is(err, redis.Nil) {
		return originURL, nil
	}

	// 以下是redis nil的情况
	var short mysql.Shorturl
	err = invoker.DB.WithContext(ctx).Where("code = ?", code).Find(&short).Error
	// gorm的 not found比较坑，懒得判断了
	if err != nil {
		return "", fmt.Errorf("find db failed, err: %w", err)
	}
	if short.ID == 0 {
		return "", ErrNotFound
	}
	err = invoker.Redis.SetEX(ctx, econf.GetString("shorturl.redisPrefix")+":"+code, short.OriginUrl, xtime.Duration("1h"))
	if err != nil {
		// 记录一下日志。这个错误不影响服务
		elog.Error("设置redis", elog.FieldErr(err), elog.FieldValueAny(code))
		return originURL, nil
	}
	return short.OriginUrl, nil
}
