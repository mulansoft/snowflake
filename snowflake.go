package snowflake

import (
	"sync"
	"time"
	"errors"
)

const (
	nodeBits  uint8 = 10
	stepBits  uint8 = 12
	nodeMax   int64 = -1 ^ (-1 << nodeBits)
	stepMax   int64 = -1 ^ (-1 << stepBits)
	timeShift uint8 = nodeBits + stepBits
	nodeShift uint8 = stepBits
)

// 起始时间戳 (毫秒数显示)
var Epoch int64 = 1530892800000 // timestamp 2018-07-07:0:0:0 GMT +8

type Node struct {
	mu sync.Mutex	// 保证并发安全
	timestamp int64
	node	  int64
	step	  int64
}

func (n *Node) Generate() int64 {
	n.mu.Lock() // 保证并发安全, 加锁
	defer n.mu.Unlock() // 解锁

	// 获取当前时间的时间戳 (毫秒数显示)
	now := time.Now().UnixNano() / 1e6

	if n.timestamp == now {
		// step 步进 1 
		n.step ++
		// 当前 step 用完
		if n.step > stepMax {
			// 等待本毫秒结束
			for now <= n.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 本毫秒内 step 用完
		n.step = 0
	}

	n.timestamp = now
	return (now - Epoch) << timeShift | (n.node << nodeShift) | (n.step)
}

func NewNode(node int64) (*Node, error) {
	if node < 0 || node > nodeMax {
		return nil, errors.New("node must be between 0 and 1023")
	}

	return &Node{
		timestamp: 0,
		node:      node,
		step:	   0,
	}, nil
}