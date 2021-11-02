module github.com/cheolgyu/stock-write-project-rebound

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1 // indirect
	github.com/cheolgyu/stock-write-common v0.0.0
	github.com/cheolgyu/stock-write-model v0.0.0
	github.com/gchaincl/dotsql v1.0.0
	github.com/swithek/dotsqlx v1.0.0
)

replace (
	github.com/cheolgyu/stock-write-common v0.0.0 => ../stock-write-common
	github.com/cheolgyu/stock-write-model v0.0.0 => ../stock-write-model
)
