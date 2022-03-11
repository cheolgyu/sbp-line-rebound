package main

import (
	_ "github.com/cheolgyu/sbm-base/db"
	_ "github.com/cheolgyu/sbm-base/env"
	"github.com/cheolgyu/sbm-base/logging"
	"github.com/cheolgyu/sbp-line-rebound/src/dao"
	"github.com/cheolgyu/sbp-line-rebound/src/handler"
)

func main() {
	defer logging.ElapsedTime()()
	project_run()
}
func project_run() {

	handler.ReBoundHandler()
	dao.Update_info()
}
