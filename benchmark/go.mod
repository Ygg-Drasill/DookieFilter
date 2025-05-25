module benchmark

go 1.24

replace github.com/Ygg-Drasill/DookieFilter/common => ../common

require (
	github.com/Ygg-Drasill/DookieFilter/common v0.0.0-00010101000000-000000000000
	github.com/pebbe/zmq4 v1.4.0
)

require github.com/google/uuid v1.6.0 // indirect
