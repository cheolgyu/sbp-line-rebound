package main

import (
	_ "github.com/cheolgyu/stock-write-common/db"
	_ "github.com/cheolgyu/stock-write-common/env"
	"github.com/cheolgyu/stock-write-common/logging"
	"github.com/cheolgyu/stock-write-project-rebound/src/dao"
	"github.com/cheolgyu/stock-write-project-rebound/src/handler"
)

func main() {
	defer logging.ElapsedTime()()
	project_run()
}
func project_run() {

	handler.ReBoundHandler()
	dao.Update_info()
}
