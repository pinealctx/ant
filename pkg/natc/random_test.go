package natc

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

//This function is test nats request/response random key conflict
func TestRandom(t *testing.T) {
	var x = rand.New(rand.NewSource(time.Now().UnixNano()))
	var mm = make(map[int64]struct{})
	var (
		y     int64
		exist bool
	)
	for i := 0; i < 10000000; i++ {
		y = x.Int63()
		_, exist = mm[y]
		if exist {
			panic(fmt.Sprintf("duplicate.at.%d", i))
		}
		mm[y] = struct{}{}
	}
	t1(t2())
}

func t1(x, y int64) {
	fmt.Println(x, y)
}

func t2() (int64, int64) {
	return 0, 1
}
