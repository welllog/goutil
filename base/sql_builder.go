package base

import (
	"strings"
)

type BatchUpdateSqlBuilder struct {
	TableName string
	updateMap map[string][]caseElm
	condition []whereElm
}

type caseElm struct {
	express string
	args    []interface{}
}

type whereElm struct {
	express string
	args    []interface{}
}

func (b *BatchUpdateSqlBuilder) Where(express string, args ...interface{}) {
	b.condition = append(b.condition, whereElm{
		express: express,
		args:    args,
	})
}

func (b *BatchUpdateSqlBuilder) SetWhen(field, express string, args ...interface{}) {
	if b.updateMap == nil {
		b.updateMap = make(map[string][]caseElm, 10)
	}
	if _, ok := b.updateMap[field]; !ok {
		b.updateMap[field] = make([]caseElm, 0, 100)
	}
	b.updateMap[field] = append(b.updateMap[field], caseElm{
		express: express,
		args:    args,
	})
}

func (b *BatchUpdateSqlBuilder) SetWhenEqual(field string, value interface{}, caseField string, caseValue interface{}) {
	b.SetWhen(field, "when "+caseField+" = ? then ?", caseValue, value)
}

func (b *BatchUpdateSqlBuilder) Build() (string, []interface{}) {
	if b.TableName == "" {
		return "", nil
	}

	updateSQL := "update " + b.TableName + " set "
	args := make([]interface{}, 0, len(b.updateMap)*2)
	for field, cases := range b.updateMap {
		updateSQL += field + " = case "
		for _, caseElm := range cases {
			updateSQL += caseElm.express + " "
			args = append(args, caseElm.args...)
		}
		updateSQL += "end,"
	}

	updateSQL = strings.TrimRight(updateSQL, ",")
	if len(b.condition) > 0 {
		updateSQL += " where "
		for _, whereElm := range b.condition {
			updateSQL += whereElm.express + " and "
			args = append(args, whereElm.args...)
		}
		updateSQL = strings.TrimSuffix(updateSQL, " and ")
	} else {
		// 防止误更全部
		updateSQL += " where 1 = 2"
	}

	return updateSQL, args
}
