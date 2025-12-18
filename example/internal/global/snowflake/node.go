package snowflake

import "github.com/syralon/snowflake"

var instance *snowflake.Generator

func Setup(nodeID int) (err error) {
	instance, err = snowflake.New(nodeID)
	return err
}

func Next() int64 {
	return instance.Next()
}
