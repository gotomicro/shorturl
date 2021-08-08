package job

import (
	"fmt"

	"github.com/gotomicro/ego/task/ejob"
	"shorturl/pkg/invoker"
	"shorturl/pkg/mysql"
)

func RunInstall(ctx ejob.Context) error {
	models := []interface{}{
		&mysql.Shorturl{},
	}
	err := invoker.DB.WithContext(ctx.Ctx).AutoMigrate(models...)
	if err != nil {
		return err
	}
	fmt.Println("create table ok")
	return nil
}
