package htmltmp

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Table struct {
	id     string
	border int
	styles []string
	titles []string
	rows   [][]string
}

func NewTable(cap int) *Table {
	if cap <= 0 {
		cap = 30
	}
	id := strconv.Itoa(time.Now().Nanosecond()) + strconv.Itoa(rand.Intn(9999))
	return &Table{id: id, rows: make([][]string, 0, cap+1)}
}

func (t *Table) AddTitle(titles []string) {
	t.titles = titles
}

func (t *Table) AddRow(rows []string) {
	t.rows = append(t.rows, rows)
}

func (t *Table) SetBorder(border int) {
	t.border = border
}

func (t *Table) SetStyle(style string) {
	t.styles = append(t.styles, style)
}

func (t *Table) DefaultStyle() {
	t.SetBorder(1)
	t.SetStyle("{border-collapse: collapse;text-align: center;}")
}

func (t *Table) String() string {
	var sb strings.Builder
	sb.WriteString("<style>\n")
	for _, style := range t.styles {
		sb.WriteString("#tmp-table-")
		sb.WriteString(t.id)
		sb.WriteString(" ")
		sb.WriteString(style)
		sb.WriteString("\n")
	}
	sb.WriteString("</style>\n")

	sb.WriteString(`<table border="`)
	sb.WriteString(strconv.Itoa(t.border))
	sb.WriteString("\" id=\"tmp-table-")
	sb.WriteString(t.id)
	sb.WriteString("\">\n")

	if len(t.titles) > 0 {
		sb.WriteString("<tr>\n")
		for _, title := range t.titles {
			sb.WriteString("<th>")
			sb.WriteString(title)
			sb.WriteString("</th>\n")
		}
		sb.WriteString("</tr>\n")
	}

	for _, row := range t.rows {
		sb.WriteString("<tr>\n")
		for _, td := range row {
			sb.WriteString("<td>")
			sb.WriteString(td)
			sb.WriteString("</td>\n")
		}
		sb.WriteString("</tr>\n")
	}

	sb.WriteString("</table>\n")
	return sb.String()
}
