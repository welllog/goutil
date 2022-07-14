package structtag

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"strconv"
)

type FileStructs struct {
	fileset *token.FileSet
	file    *ast.File
}

func ParseStruct(src []byte) (*FileStructs, error) {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, "", src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("ParseFile error: %w", err)
	}

	return &FileStructs{
		fileset: fileset,
		file:    file,
	}, nil
}

func (f *FileStructs) HandleStruct(fn func(*string, *Tags)) {
	ast.Inspect(f.file, func(x ast.Node) bool {
		s, ok := x.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range s.Fields.List {
			if len(field.Names) == 0 {
				continue
			}

			tags := &Tags{}
			if field.Tag != nil {
				tagValue, _ := strconv.Unquote(field.Tag.Value)
				pTags, err := Parse(tagValue)
				if err == nil {
					tags = pTags
				}
			} else {
				field.Tag = &ast.BasicLit{
					ValuePos: field.End() + 1,
					Kind:     token.STRING,
				}
			}
			fn(&field.Names[0].Name, tags)
			field.Tag.Value = "`" + tags.String() + "`"
		}

		return true
	})
}

func (f *FileStructs) Save(dst io.Writer) error {
	return format.Node(dst, f.fileset, f.file)
}
