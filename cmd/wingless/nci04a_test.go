package main
import("bytes";"testing";"strings")
func TestNativeNCI04AOfflineCLI(t *testing.T){
 var out bytes.Buffer
 e:=runNCI04A(strings.NewReader(`{"op":"observe","src":[1,0,0,0],"dst":[0,1,0,0],"action":0,"outcome":1}`+"\n"+`{"op":"predict","src":[1,0,0,0],"dst":[0,1,0,0],"action":0}`+"\n"+`{"op":"snapshot"}`+"\n"),&out)
 if e!=nil{t.Fatal(e)}
 if !strings.Contains(out.String(),"probability")||!strings.Contains(out.String(),"entities"){t.Fatal("incomplete result",out.String())}
 var bad bytes.Buffer
 if e=runNCI04A(strings.NewReader(`{"op":"execute"}`+"\n"),&bad);e==nil{t.Fatal("unauthorized action permitted")}
}
