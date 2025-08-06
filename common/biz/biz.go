package biz

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// 由两个uid拼接双方的对话组id
// 拼接id小号id在前
func GetGroupId(uid1 int64, uid2 int64) string {
	var group_id string
	if uid1 < uid2 {
		group_id = fmt.Sprintf("%d_%d", uid1, uid2)
	} else {
		group_id = fmt.Sprintf("%d_%d", uid2, uid1)
	}
	return group_id
}

// 从groupId和其中一个uid, 得到组内另一人的uid
func GetFriendIdFromGroupId(groupId string, uid int64) (int64, error) {
	// 将字符串按照"_"符号进行分割到arr切片中
	arr := strings.Split(groupId, "_")
	// 将自己的uid转成字符串方便之后进行判断
	uid1_str := strconv.Itoa(int(uid))
	var uid2_str string
	// 检测出另外一个人的uid
	if arr[0] == uid1_str {
		uid2_str = arr[1]
	} else {
		uid2_str = arr[0]
	}
	// 将另一个人的uid转成数字返回
	uid2, err := strconv.ParseInt(uid2_str, 10, 64)
	if err != nil {
		return 0, err
	}
	return uid2, nil
}

// 随机字符串
// 获取随机字符串长度为length
func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

// 获取uuid
// 生成一个uuid
// 格式为123e4567-e89b-12d3-a456-426614174000
func GetUuid() string {
	uuid := uuid.New()
	return uuid.String()
}
