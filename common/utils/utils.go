package utils

import "github.com/google/uuid"

// 获取uuid
func GenUuid() string {
	uuid := uuid.New()
	key := uuid.String()
	return key
}
