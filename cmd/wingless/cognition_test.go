package main
import("bytes";"strings";"testing")
func TestOfflineCognition(t *testing.T){
 input:=`{"op":"observe","action":1,"prior":0,"outcome":1,"rate":0.1}`+"\n"+`{"op":"forecast_two","prior":0,"action":1,"next_action":0}`+"\n"+`{"op":"snapshot"}`+"\n"
 var out bytes.Buffer
 if e:=runCognition(strings.NewReader(input),&out);e!=nil{t.Fatal(e)}
 if !strings.Contains(out.String(),"probability")||!strings.Contains(out.String(),"weights"){t.Fatal(out.String())}
 if e:=runCognition(strings.NewReader(`{"op":"execute"}`+"\n"),&out);e==nil{t.Fatal("unauthorized command accepted")}
}
