package ice
import("os";"path/filepath";"testing")
func TestSourceWins(t *testing.T){r:=t.TempDir();p:=filepath.Join(r,"a.go");os.WriteFile(p,[]byte("package a\nfunc Own() {}\n"),0600);if e:=Save(r);e!=nil{t.Fatal(e)};q,e:=Query(r,"Own");if e!=nil||len(q)!=1{t.Fatal(q,e)};os.WriteFile(p,[]byte("package a\nfunc Changed() {}\n"),0600);if _,e=Query(r,"Own");e==nil{t.Fatal("stale index permitted")};if e=Save(r);e!=nil{t.Fatal(e)};q,e=Query(r,"Changed");if e!=nil||len(q)!=1{t.Fatal(q,e)}}
