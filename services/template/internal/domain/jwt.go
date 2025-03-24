package domain

import (
	"os"
	"strconv"
	"time"
)

func TokenExpiry() int64 {
	expireTime, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRE_TIME"))
	if err != nil || expireTime <= 0 {
		expireTime = 24
	}

	return time.Now().Add(time.Duration(expireTime) * time.Hour).Unix()
}
