module github.com/Byfengfeng/wl_net

go 1.19

require (
	git.byfengfeng.cn/byaf/ringbuf v0.0.1
	github.com/bwmarrin/snowflake v0.3.0
	github.com/panjf2000/ants/v2 v2.7.4
	go.uber.org/zap v1.24.0
	golang.org/x/net v0.10.0
)

require (
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
)

replace git.byfengfeng.cn/byaf/ringbuf => ../ring_buf