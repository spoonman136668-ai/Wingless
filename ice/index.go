// Package ice builds advisory, reproducible syntax indexes. Source remains authoritative.
package ice

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Symbol struct {
	Name     string   `json:"name"`
	Kind     string   `json:"kind"`
	File     string   `json:"file"`
	Start    int      `json:"start"`
	End      int      `json:"end"`
	Receiver string   `json:"receiver,omitempty"`
	Calls    []string `json:"calls,omitempty"`
}
type File struct {
	Path    string   `json:"path"`
	SHA256  string   `json:"sha256"`
	Package string   `json:"package,omitempty"`
	Imports []string `json:"imports,omitempty"`
}
type Implementation struct {
	Interface  string `json:"interface"`
	Expression string `json:"expression"`
	File       string `json:"file"`
	Line       int    `json:"line"`
}
type Index struct {
	Implementations []Implementation `json:"implementations"`
	Version         int              `json:"version"`
	Authority       string           `json:"authority"`
	Files           []File           `json:"files"`
	Symbols         []Symbol         `json:"symbols"`
}

func Build(root string) (Index, error) {
	idx := Index{Version: 1, Authority: "advisory syntax only; source wins"}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, e := filepath.Rel(root, path)
		if e != nil {
			return e
		}
		if d.IsDir() {
			if rel != "." && ((strings.HasPrefix(d.Name(), ".") && d.Name() != ".github") || d.Name() == "bin" || d.Name() == "evidence") {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink rejected: %s", rel)
		}
		ext := filepath.Ext(path)
		if ext != ".go" && ext != ".md" && ext != ".json" && ext != ".ps1" && ext != ".yml" && ext != ".yaml" && filepath.Base(path) != ".gitattributes" && filepath.Base(path) != ".gitignore" && filepath.Base(path) != "go.mod" {
			return nil
		}
		b, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		h := sha256.Sum256(b)
		f := File{Path: filepath.ToSlash(rel), SHA256: hex.EncodeToString(h[:])}
		if ext == ".go" {
			set := token.NewFileSet()
			tree, e := parser.ParseFile(set, path, b, 0)
			if e != nil {
				return e
			}
			f.Package = tree.Name.Name
			for _, i := range tree.Imports {
				f.Imports = append(f.Imports, strings.Trim(i.Path.Value, "\""))
			}
			for _, decl := range tree.Decls {
				switch x := decl.(type) {
				case *ast.FuncDecl:
					s := Symbol{Name: x.Name.Name, Kind: "function", File: f.Path, Start: set.Position(x.Pos()).Line, End: set.Position(x.End()).Line}
					if x.Recv != nil {
						s.Kind = "method"
						t := x.Recv.List[0].Type
						if p, ok := t.(*ast.StarExpr); ok {
							t = p.X
						}
						if n, ok := t.(*ast.Ident); ok {
							s.Receiver = n.Name
						}
					}
					if x.Body != nil {
						ast.Inspect(x.Body, func(n ast.Node) bool {
							if c, ok := n.(*ast.CallExpr); ok {
								switch a := c.Fun.(type) {
								case *ast.Ident:
									s.Calls = append(s.Calls, a.Name)
								case *ast.SelectorExpr:
									if q, ok := a.X.(*ast.Ident); ok {
										s.Calls = append(s.Calls, q.Name+"."+a.Sel.Name)
									}
								}
							}
							return true
						})
					}
					sort.Strings(s.Calls)
					idx.Symbols = append(idx.Symbols, s)
				case *ast.GenDecl:
					for _, spec := range x.Specs {
						if v, ok := spec.(*ast.ValueSpec); ok {
							for _, name := range v.Names {
								if name.Name != "_" {
									idx.Symbols = append(idx.Symbols, Symbol{Name: name.Name, Kind: x.Tok.String(), File: f.Path, Start: set.Position(v.Pos()).Line, End: set.Position(v.End()).Line})
								}
							}
						}

						if v, ok := spec.(*ast.ValueSpec); ok && v.Type != nil {
							for n, name := range v.Names {
								if name.Name == "_" && n < len(v.Values) {
									var a, b bytes.Buffer
									format.Node(&a, set, v.Type)
									format.Node(&b, set, v.Values[n])
									idx.Implementations = append(idx.Implementations, Implementation{a.String(), b.String(), f.Path, set.Position(v.Pos()).Line})
								}
							}
						}

						if t, ok := spec.(*ast.TypeSpec); ok {
							kind := "type"
							if _, ok := t.Type.(*ast.InterfaceType); ok {
								kind = "interface"
							}
							idx.Symbols = append(idx.Symbols, Symbol{Name: t.Name.Name, Kind: kind, File: f.Path, Start: set.Position(t.Pos()).Line, End: set.Position(t.End()).Line})
						}
					}
				}
			}
		}
		idx.Files = append(idx.Files, f)
		return nil
	})
	return idx, err
}
func derived(idx Index) map[string]any {
	deps := map[string][]string{}
	tests := []Symbol{}
	for _, f := range idx.Files {
		if f.Package != "" {
			deps[f.Path] = f.Imports
		}
	}
	for _, s := range idx.Symbols {
		if strings.HasPrefix(s.Name, "Test") && strings.HasSuffix(s.File, "_test.go") {
			tests = append(tests, s)
		}
	}
	return map[string]any{"manifest.json": idx, "architecture-map.json": idx.Files, "symbol-map.json": idx.Symbols, "dependency-map.json": deps, "test-map.json": tests, "implementations.json": idx.Implementations}
}
func Save(root string) error {
	idx, e := Build(root)
	if e != nil {
		return e
	}
	if e = os.MkdirAll(filepath.Join(root, ".ice"), 0755); e != nil {
		return e
	}
	for name, value := range derived(idx) {
		b, e := json.MarshalIndent(value, "", "  ")
		if e != nil {
			return e
		}
		if e = os.WriteFile(filepath.Join(root, ".ice", name), append(b, '\n'), 0644); e != nil {
			return e
		}
	}
	for _, name := range []string{"integration-seams.json", "decisions.json"} {
		b, e := os.ReadFile(filepath.Join(root, "docs", name))
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return e
		}
		if e = os.WriteFile(filepath.Join(root, ".ice", name), b, 0644); e != nil {
			return e
		}
	}
	return nil
}
func Load(root string) (Index, error) {
	var i Index
	b, e := os.ReadFile(filepath.Join(root, ".ice", "manifest.json"))
	if e != nil {
		return i, e
	}
	e = json.Unmarshal(b, &i)
	return i, e
}
func Validate(root string) error {
	old, e := Load(root)
	if e != nil {
		return e
	}
	now, e := Build(root)
	if e != nil {
		return e
	}
	a, _ := json.Marshal(old)
	b, _ := json.Marshal(now)
	if string(a) != string(b) {
		return fmt.Errorf("ICE_STALE: rebuild from source")
	}
	for name, value := range derived(now) {
		want, _ := json.MarshalIndent(value, "", "  ")
		got, e := os.ReadFile(filepath.Join(root, ".ice", name))
		if e != nil || string(got) != string(append(want, '\n')) {
			return fmt.Errorf("ICE_STALE: %s", name)
		}
	}
	for _, name := range []string{"integration-seams.json", "decisions.json"} {
		source, e := os.ReadFile(filepath.Join(root, "docs", name))
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return e
		}
		got, e := os.ReadFile(filepath.Join(root, ".ice", name))
		if e != nil || string(got) != string(source) {
			return fmt.Errorf("ICE_STALE: %s", name)
		}
	}
	return nil
}
func Query(root, q string) ([]Symbol, error) {
	if strings.TrimSpace(q) == "" {
		return nil, fmt.Errorf("empty query")
	}
	if e := Validate(root); e != nil {
		return nil, e
	}
	i, e := Load(root)
	if e != nil {
		return nil, e
	}
	var out []Symbol
	implementationFiles := map[string]bool{}
	for _, impl := range i.Implementations {
		if strings.Contains(strings.ToLower(impl.Interface), strings.ToLower(q)) {
			implementationFiles[impl.File] = true
		}
	}
	for _, s := range i.Symbols {
		if implementationFiles[s.File] || strings.Contains(strings.ToLower(s.Name+" "+s.File+" "+s.Receiver+" "+strings.Join(s.Calls, " ")), strings.ToLower(q)) {
			out = append(out, s)
		}
	}
	return out, nil
}
