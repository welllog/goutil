package htmltmp

import (
	"fmt"
	"testing"
)

func TestTable_String(t *testing.T) {
	table := NewTable(20)
	table.AddRow([]string{"china", "1000"})
	table.AddTitle([]string{"country", "count"})
	table.DefaultStyle()
	fmt.Println(table.String())
}
