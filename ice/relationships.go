package ice

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type Relationship struct {
	Kind string `json:"kind"`
	From string `json:"from"`
	To   string `json:"to"`
	File string `json:"file"`
	Line int    `json:"line"`
}
type TypeReport struct {
	GOOS          string         `json:"goos"`
	GOARCH        string         `json:"goarch"`
	Relationships []Relationship `json:"relationships"`
	Diagnostics   []string       `json:"diagnostics"`
	Authority     string         `json:"authority"`
}
type typeLoader struct {
	root     string
	set      *token.FileSet
	groups   map[string][]*ast.File
	packages map[string]*types.Package
	busy     map[string]bool
	standard types.Importer
	report   TypeReport
}

func (l *typeLoader) Import(path string) (*types.Package, error) {
	if p := l.packages[path]; p != nil {
		return p, nil
	}
	files, local := l.groups[path]
	if !local {
		return l.standard.Import(path)
	}
	if l.busy[path] {
		return nil, fmt.Errorf("import cycle")
	}
	l.busy[path] = true
	defer delete(l.busy, path)
	info := &types.Info{Uses: map[*ast.Ident]types.Object{}, Defs: map[*ast.Ident]types.Object{}}
	cfg := types.Config{Importer: l, GoVersion: "go1.25"}
	pkg, e := cfg.Check(path, l.set, files, info)
	if e != nil {
		return nil, e
	}
	l.packages[path] = pkg
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			var id *ast.Ident
			switch x := call.Fun.(type) {
			case *ast.Ident:
				id = x
			case *ast.SelectorExpr:
				id = x.Sel
			}
			if id == nil {
				return true
			}
			target, ok := info.Uses[id].(*types.Func)
			if !ok {
				return true
			}
			position := l.set.Position(call.Pos())
			file, _ := filepath.Rel(l.root, position.Filename)
			l.report.Relationships = append(l.report.Relationships, Relationship{Kind: "resolved_call_target", From: path, To: target.FullName(), File: filepath.ToSlash(file), Line: position.Line})
			return true
		})
	}
	names := pkg.Scope().Names()
	for _, name := range names {
		obj, ok := pkg.Scope().Lookup(name).(*types.TypeName)
		if !ok {
			continue
		}
		iface, ok := obj.Type().Underlying().(*types.Interface)
		if !ok {
			continue
		}
		iface.Complete()
		for _, other := range names {
			candidate, ok := pkg.Scope().Lookup(other).(*types.TypeName)
			if !ok {
				continue
			}
			if _, ok := candidate.Type().Underlying().(*types.Interface); ok {
				continue
			}
			if types.Implements(candidate.Type(), iface) || types.Implements(types.NewPointer(candidate.Type()), iface) {
				pos := l.set.Position(candidate.Pos())
				file, _ := filepath.Rel(l.root, pos.Filename)
				l.report.Relationships = append(l.report.Relationships, Relationship{"same_package_implementation", path + "." + other, path + "." + name, filepath.ToSlash(file), pos.Line})
			}
		}
	}
	return pkg, nil
}

// Relationships checks active-platform non-test Go source without executing a compiler or model.
// Only successfully type-checked packages produce relationships; diagnostics preserve incomplete coverage.
func Relationships(root string) (TypeReport, error) {
	if e := Validate(root); e != nil {
		return TypeReport{}, e
	}
	idx, e := Load(root)
	if e != nil {
		return TypeReport{}, e
	}
	set := token.NewFileSet()
	l := &typeLoader{root: root, set: set, groups: map[string][]*ast.File{}, packages: map[string]*types.Package{}, busy: map[string]bool{}, standard: importer.ForCompiler(set, "source", nil), report: TypeReport{GOOS: runtime.GOOS, GOARCH: runtime.GOARCH, Authority: "derived; active platform; no runtime call-graph guarantee"}}
	module := "github.com/spoonman136668-ai/Wingless"
	for _, f := range idx.Files {
		if f.Package == "" || strings.HasSuffix(f.Path, "_test.go") {
			continue
		}
		path := filepath.Join(root, filepath.FromSlash(f.Path))
		match, e := build.Default.MatchFile(filepath.Dir(path), filepath.Base(path))
		if e != nil {
			return l.report, e
		}
		if !match {
			continue
		}
		tree, e := parser.ParseFile(set, path, nil, 0)
		if e != nil {
			return l.report, e
		}
		pkg := module + "/" + filepath.ToSlash(filepath.Dir(f.Path))
		l.groups[pkg] = append(l.groups[pkg], tree)
	}
	keys := []string{}
	for k := range l.groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if _, e = l.Import(k); e != nil {
			l.report.Diagnostics = append(l.report.Diagnostics, k+": "+e.Error())
		}
	}
	sort.Slice(l.report.Relationships, func(i, j int) bool {
		a, b := l.report.Relationships[i], l.report.Relationships[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		return a.To < b.To
	})
	return l.report, nil
}
