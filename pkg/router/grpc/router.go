package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gotomicro/ego/server/egrpc"
	"gorm.io/gorm"
	"shorturl/pkg/invoker"
	"shorturl/pkg/mysql"
	"shorturl/proto"
)

func ServeGRPC() *egrpc.Component {
	srv := egrpc.Load("server.grpc").Build()
	proto.RegisterShorturlServer(srv.Server, &Shorturl{})
	return srv
}

type Shorturl struct {
	proto.UnimplementedShorturlServer
}

func (*Shorturl) Create(ctx context.Context, req *proto.CreateRequest) (*proto.CreateReply, error) {
	var url mysql.Shorturl
	err := invoker.DB.Where("origin_url = ?", req.OriginURL).Find(&url).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("short url create failed, err: %w", err)
	}
	if url.ID > 0 {
		return nil, fmt.Errorf("short url exist")
	}

	newURL := mysql.Shorturl{
		OriginUrl: req.OriginURL,
		Code:      "",
		CallCnt:   0,
		Ctime:     time.Now().Unix(),
		Utime:     time.Now().Unix(),
	}
	err = invoker.DB.Create(&newURL).Error
	if err != nil {
		return nil, fmt.Errorf("create short url failed, err: %w", err)
	}
	code := base10ToBase58(newURL.ID)
	err = invoker.DB.Model(mysql.Shorturl{}).Where("id = ?", newURL.ID).Update("code", code).Error
	if err != nil {
		return nil, fmt.Errorf("update short url failed, err: %w", err)
	}

	return &proto.CreateReply{
		Code: code,
	}, nil
}

// 去掉 L、l、o、O
var elements = "0123456789abcdefghjkmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ"

// base10ToBase58 十进制转换成58进制
func base10ToBase58(n int64) string {
	var str string
	for n != 0 {
		str += string(elements[n%58])
		n /= 58
	}
	return str
}
