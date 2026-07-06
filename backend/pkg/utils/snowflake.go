package utils

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

const defaultSnowflakeNodeID int64 = 1

var (
	snowflakeNode     *snowflake.Node
	snowflakeNodeOnce sync.Once
)

// GetSnowflakeNode returns the singleton snowflake node used by the application.
//
// 节点 ID 当前固定为 1，用于保持原有单节点部署下的 ID 生成行为。初始化失败代表节点配置不合法，
// 属于启动期不可恢复错误，因此会触发 panic。
func GetSnowflakeNode() *snowflake.Node {
	snowflakeNodeOnce.Do(func() {
		node, err := snowflake.NewNode(defaultSnowflakeNodeID)
		if err != nil {
			panic(err)
		}
		snowflakeNode = node
	})
	return snowflakeNode
}

// GenerateID generates a new Snowflake ID as uint64.
//
// 返回值用于业务模型主键。底层开源库生成的是正数 int64 Snowflake ID，转换为 uint64 后保持数值不变。
func GenerateID() uint64 {
	return uint64(GetSnowflakeNode().Generate().Int64())
}
