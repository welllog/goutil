package base

import "testing"

func TestBatchUpdateSqlBuilder_Build(t *testing.T) {
	builder := BatchUpdateSqlBuilder{TableName: "users"}
	builder.SetWhenEqual("name", "test", "id", 1)
	builder.SetWhenEqual("name", "test2", "id", 2)
	builder.SetWhen("name", "when id = ? then ?", "test3", 4)
	builder.Where("status = ?", 1)
	t.Log(builder.Build())
}
