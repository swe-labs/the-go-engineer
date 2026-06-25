package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type DocStatus int

const (
	Missing DocStatus = iota
	Present
)

func (d DocStatus) String() string {
	switch d {
	case Missing:
		return "Missing"
	case Present:
		return "Present"
	default:
		return "Unknown"
	}
}

type DocItem struct {
	Name    string
	Kind    string
	Status  DocStatus
	Line    int
	DocText string
}

type DocReport struct {
	FilePath string
	Items    []DocItem
	Covered  int
	Total    int
}

func (r DocReport) Coverage() float64 {
	if r.Total == 0 {
		return 100.0
	}
	return float64(r.Covered) / float64(r.Total) * 100.0
}

func AnalyzeFile(filePath string) (DocReport, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return DocReport{}, fmt.Errorf("parse error: %w", err)
	}

	report := DocReport{FilePath: filePath}

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.TYPE || d.Tok == token.VAR || d.Tok == token.CONST {
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						doc := s.Doc
						if doc == nil {
							doc = d.Doc
						}
						item := DocItem{
							Name:   s.Name.Name,
							Kind:   "type",
							Status: docStatusFromComment(doc),
							Line:   fset.Position(s.Pos()).Line,
						}
						if doc != nil {
							item.DocText = doc.Text()
						}
						report.Items = append(report.Items, item)
						if item.Status == Present {
							report.Covered++
						}
						report.Total++
					case *ast.ValueSpec:
						for _, name := range s.Names {
							if name.IsExported() {
								item := DocItem{
									Name:   name.Name,
									Kind:   varOrConst(d.Tok),
									Status: docStatusFromComment(s.Doc),
									Line:   fset.Position(s.Pos()).Line,
								}
								if s.Doc != nil {
									item.DocText = s.Doc.Text()
								}
								report.Items = append(report.Items, item)
								if item.Status == Present {
									report.Covered++
								}
								report.Total++
							}
						}
					}
				}
			}
		case *ast.FuncDecl:
			if d.Name.IsExported() {
				item := DocItem{
					Name:   d.Name.Name,
					Kind:   funcKind(d),
					Status: docStatusFromComment(d.Doc),
					Line:   fset.Position(d.Pos()).Line,
				}
				if d.Doc != nil {
					item.DocText = d.Doc.Text()
				}
				report.Items = append(report.Items, item)
				if item.Status == Present {
					report.Covered++
				}
				report.Total++
			}
		}
	}

	return report, nil
}

func docStatusFromComment(doc *ast.CommentGroup) DocStatus {
	if doc != nil && len(doc.List) > 0 {
		return Present
	}
	return Missing
}

func varOrConst(tok token.Token) string {
	if tok == token.CONST {
		return "const"
	}
	return "var"
}

func funcKind(d *ast.FuncDecl) string {
	if d.Recv != nil {
		return "method"
	}
	return "function"
}

func (r DocReport) Summary() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Documentation Coverage Report: %s\n", r.FilePath))
	b.WriteString(fmt.Sprintf("Coverage: %.1f%% (%d/%d exported symbols documented)\n\n", r.Coverage(), r.Covered, r.Total))
	for _, item := range r.Items {
		b.WriteString(fmt.Sprintf("  %s %s (line %d): %s\n", item.Kind, item.Name, item.Line, item.Status))
		if item.DocText != "" {
			b.WriteString(fmt.Sprintf("    %s", item.DocText))
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <file.go>")
		fmt.Println("Example: go run . main.go")
		os.Exit(1)
	}

	report, err := AnalyzeFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(report.Summary())
}
