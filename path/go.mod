module example.com/path

go 1.15

replace example.com/graph => ../graph

replace example.com/config => ../config

replace example.com/request => ../request

require (
	example.com/graph v0.0.0-00010101000000-000000000000
	example.com/request v0.0.0-00010101000000-000000000000
)
