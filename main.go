package main

import (
	_ "github.com/cheolgyu/base/db"
	_ "github.com/cheolgyu/base/env"
	"github.com/cheolgyu/base/logging"
	"github.com/cheolgyu/graph/bound/dao"
	"github.com/cheolgyu/graph/bound/handler"
)

func main() {
	defer logging.ElapsedTime()()
	project_run()
}
func project_run() {

	handler.ReBoundHandler()
	dao.Update_info()
}
