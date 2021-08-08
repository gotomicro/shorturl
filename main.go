package main

import (
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egovernor"
	"github.com/gotomicro/ego/task/ejob"
	"shorturl/pkg/invoker"
	"shorturl/pkg/job"
	"shorturl/pkg/router/grpc"
	"shorturl/pkg/router/http"
)

func main() {
	err := ego.New().
		Invoker(
			invoker.Init,
		).
		Job(
			ejob.Job("install", job.RunInstall),
		).
		Serve(
			egovernor.Load("server.governor").Build(),
			http.ServeHTTP(),
			grpc.ServeGRPC(),
		).
		Run()
	if err != nil {
		elog.Panic("start up error: " + err.Error())
	}

}
