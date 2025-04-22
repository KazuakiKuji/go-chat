package random

import (
	"math/rand"
	"time"
)

// ローカル乱数生成器
var LocalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
