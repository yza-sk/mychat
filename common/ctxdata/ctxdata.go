package ctxdata

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

// 从context中解析出uid并转化为int64返回
func GetUidFromCtx(ctx context.Context) int64 {
	var uid int64
	if jsonUid, ok := ctx.Value("uid").(json.Number); ok {
		if int64Uid, err := jsonUid.Int64(); err == nil {
			uid = int64Uid
		} else {
			logx.WithContext(ctx).Errorf("GetUidFromCtx err : %+v", err)
		}
	}
	return uid
}
