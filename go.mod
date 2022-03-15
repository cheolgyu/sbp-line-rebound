module github.com/cheolgyu/graph

go 1.16

require (
	github.com/cheolgyu/tb v0.0.0
	github.com/cheolgyu/base v0.0.0
	github.com/cheolgyu/model v0.0.0
	github.com/gchaincl/dotsql v1.0.0
	github.com/swithek/dotsqlx v1.0.0
)

replace (
	github.com/cheolgyu/tb v0.0.0 => ../tb
	github.com/cheolgyu/base v0.0.0 => ../base
	github.com/cheolgyu/model v0.0.0 => ../model
)
