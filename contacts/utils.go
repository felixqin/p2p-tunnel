package contacts

import (
	"math/rand"
	"time"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().Unix()))
}

// RandString 生成随机字符串
func makeRandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		n := random.Intn(62)
		a := n / 26
		r := n % 26
		switch a {
		case 0:
			bytes[i] = byte(r + 'A')
		case 1:
			bytes[i] = byte(r + 'a')
		default:
			bytes[i] = byte(r + '0')
		}
	}

	return string(bytes)
}
