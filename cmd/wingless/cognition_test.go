package main
import("bytes";"errors";"strings";"testing")
func TestOfflineCognition(t *testing.T){
 input:=`{"op":"observe","action":1,"prior":0,"outcome":1,"rate":0.1}`+"\n"+`{"op":"forecast_two","prior":0,"action":1,"next_action":0}`+"\n"+`{"op":"snapshot"}`+"\n"
 var out bytes.Buffer
 if e:=runCognition(strings.NewReader(input),&out);e!=nil{t.Fatal(e)}
 if !strings.Contains(out.String(),"probability")||!strings.Contains(out.String(),"weights"){t.Fatal(out.String())}
 if e:=runCognition(strings.NewReader(`{"op":"execute"}`+"\n"),&out);e==nil{t.Fatal("unauthorized command accepted")}
}

func TestOfflineCognitionOversizeInputHasNoPartialOutput(t *testing.T) {
 line:=append([]byte(`{"op":"snapshot"}`),bytes.Repeat([]byte(" "),60000)...)
 line=append(line,'\n')
 data:=bytes.Repeat(line,18)
 if len(data)<=(1<<20){t.Fatal("fixture does not exceed source byte limit")}
 var out bytes.Buffer
 err:=runCognition(bytes.NewReader(data),&out)
 if err==nil||!strings.Contains(err.Error(),"exceeds 1 MiB"){t.Fatalf("oversized request not rejected: %v",err)}
 if out.Len()!=0{t.Fatalf("rejected request emitted %d bytes",out.Len())}
}

type cognitionFailingReader struct{read bool}
func (r *cognitionFailingReader) Read(p []byte)(int,error){
 if !r.read{r.read=true;return copy(p,[]byte(`{"op":"snapshot"}`+"\n")),nil}
 return 0,errors.New("injected cognition reader failure")
}
func TestOfflineCognitionReaderFailureHasNoPartialOutput(t *testing.T) {
 var out bytes.Buffer
 err:=runCognition(&cognitionFailingReader{},&out)
 if err==nil||!strings.Contains(err.Error(),"injected cognition reader failure"){t.Fatalf("reader failure was not returned: %v",err)}
 if out.Len()!=0{t.Fatalf("failed reader emitted %d bytes",out.Len())}
}
