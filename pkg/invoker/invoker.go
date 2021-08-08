package invoker

import (
	"github.com/gotomicro/ego-component/egorm"
	"github.com/gotomicro/ego-component/eredis"
)

var (
	DB    *egorm.Component
	Redis *eredis.Component
)

func Init() error {
	Redis = eredis.Load("redis").Build()
	DB = egorm.Load("mysql").Build()
	return nil
}
