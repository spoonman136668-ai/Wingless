// Package ice builds advisory, reproducible syntax indexes. Source remains authoritative.
package ice

import (
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "fmt"
 "go/ast"
 "go/parser"
 "go/token"
 "io/fs"
 "os"
 "path/filepath"
 "sort"
 "strings"
)
type Symbol struct { Name string `json:"name"`; Kind string `json:"kind"`; File string `json:"file"`; Start int `json:"start"`; End int `json:"end"`; Receiver string `json:"receiver,omitempty"`; Calls []string `json:"calls,omitempty"` }
type File struct { Path string `json:"path"`; SHA256 string `json:"sha256"`; Package string `json:"package,omitempty"`; Imports []string `json:"imports,omitempty"` }
type Index struct { Version int `json:"version"`; Authority string `json:"authority"`; Files []File `json:"files"`; Symbols []Symbol `json:"symbols"` }
func Build(root string) (Index,error) {
 idx:=Index{Version:1,Authority:"advisory syntax only; source wins"}
 err:=filepath.WalkDir(root,func(path string,d fs.DirEntry,err error)error{
  if err!=nil{return err}; rel,e:=filepath.Rel(root,path);if e!=nil{return e}
  if d.IsDir(){if rel!="." && (strings.HasPrefix(d.Name(),".")||d.Name()=="bin"){return filepath.SkipDir};return nil}
  if d.Type()&os.ModeSymlink!=0{return fmt.Errorf("symlink rejected: %s",rel)}
  ext:=filepath.Ext(path);if ext!=".go"&&ext!=".md"&&ext!=".json"&&filepath.Base(path)!="go.mod"{return nil}
  b,e:=os.ReadFile(path);if e!=nil{return e};h:=sha256.Sum256(b);f:=File{Path:filepath.ToSlash(rel),SHA256:hex.EncodeToString(h[:])}
  if ext==".go" {set:=token.NewFileSet();tree,e:=parser.ParseFile(set,path,b,0);if e!=nil{return e};f.Package=tree.Name.Name
   for _,i:=range tree.Imports{f.Imports=append(f.Imports,strings.Trim(i.Path.Value,"\""))}
   for _,decl:=range tree.Decls { switch x:=decl.(type){case *ast.FuncDecl:
    s:=Symbol{Name:x.Name.Name,Kind:"function",File:f.Path,Start:set.Position(x.Pos()).Line,End:set.Position(x.End()).Line}
    if x.Recv!=nil {s.Kind="method"; t:=x.Recv.List[0].Type;if p,ok:=t.(*ast.StarExpr);ok{t=p.X};if n,ok:=t.(*ast.Ident);ok{s.Receiver=n.Name}}
    ast.Inspect(x.Body,func(n ast.Node)bool{if c,ok:=n.(*ast.CallExpr);ok{switch a:=c.Fun.(type){case *ast.Ident:s.Calls=append(s.Calls,a.Name);case *ast.SelectorExpr:if q,ok:=a.X.(*ast.Ident);ok{s.Calls=append(s.Calls,q.Name+"."+a.Sel.Name)}}};return true});sort.Strings(s.Calls);idx.Symbols=append(idx.Symbols,s)
   case *ast.GenDecl: for _,spec:=range x.Specs{if t,ok:=spec.(*ast.TypeSpec);ok{kind:="type";if _,ok:=t.Type.(*ast.InterfaceType);ok{kind="interface"};idx.Symbols=append(idx.Symbols,Symbol{Name:t.Name.Name,Kind:kind,File:f.Path,Start:set.Position(t.Pos()).Line,End:set.Position(t.End()).Line})}}
   }
   }
  };idx.Files=append(idx.Files,f);return nil
 });return idx,err
}
func Save(root string)error{idx,e:=Build(root);if e!=nil{return e};b,e:=json.MarshalIndent(idx,"","  ");if e!=nil{return e};if e=os.MkdirAll(filepath.Join(root,".ice"),0755);e!=nil{return e};return os.WriteFile(filepath.Join(root,".ice","manifest.json"),append(b,'\n'),0644)}
func Load(root string)(Index,error){var i Index;b,e:=os.ReadFile(filepath.Join(root,".ice","manifest.json"));if e!=nil{return i,e};e=json.Unmarshal(b,&i);return i,e}
func Validate(root string)error{old,e:=Load(root);if e!=nil{return e};now,e:=Build(root);if e!=nil{return e};a,_:=json.Marshal(old);b,_:=json.Marshal(now);if string(a)!=string(b){return fmt.Errorf("ICE_STALE: rebuild from source")};return nil}
func Query(root,q string)([]Symbol,error){if strings.TrimSpace(q)==""{return nil,fmt.Errorf("empty query")};if e:=Validate(root);e!=nil{return nil,e};i,e:=Load(root);if e!=nil{return nil,e};var out []Symbol;for _,s:=range i.Symbols{if strings.Contains(strings.ToLower(s.Name+" "+s.File+" "+s.Receiver+" "+strings.Join(s.Calls," ")),strings.ToLower(q)){out=append(out,s)}};return out,nil}
