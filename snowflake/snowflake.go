package snowflake

import (
	"github.com/bwmarrin/snowflake"
)

func SnowflakeId(i int64) (id snowflake.ID) {
	node, err := snowflake.NewNode(i)
	if err != nil {
		return
	}

	id = node.Generate()
	return
}
