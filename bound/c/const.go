package c

import (
	"os"
)

const SQL_DIR_DAILY = "data/sql-daily/"
const SQL_FILE_DAILY_REBOUND = SQL_DIR_DAILY
const DOTSQL_NAME_REBOUND = "insert-hist.rebound-table"

const FILE_FLAG_APPEND = os.O_RDWR | os.O_CREATE | os.O_APPEND
const FILE_FLAG_TRUNC = os.O_RDWR | os.O_CREATE | os.O_TRUNC

const INFO_NAME_UPDATED = "rebound_updated"
