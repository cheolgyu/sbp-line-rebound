module github.com/cheolgyu/sbp-line-rebound

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1 // indirect
	github.com/cheolgyu/sbm-base v0.0.0
	github.com/cheolgyu/sbm-struct v0.0.0
	github.com/gchaincl/dotsql v1.0.0
	github.com/swithek/dotsqlx v1.0.0
)

replace (
	github.com/cheolgyu/sbm-base v0.0.0 => ../sbm-base
	github.com/cheolgyu/sbm-struct v0.0.0 => ../sbm-struct
)
