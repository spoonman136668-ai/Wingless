package main

import (
 "bytes"
 "errors"
 "io"
 "strings"
 "testing"
)

const inputLimit = 1 << 20
const testPredict = `{"op":"predict","action":0,"prior":0}`

func paddedPredictLine(length int) []byte {
 line := []byte(testPredict)
 if length < len(line)+1 { panic("invalid test line length") }
 return append(append(line, bytes.Repeat([]byte(" "), length-len(line)-1)...), '\n')
}

func fullInputAtLimit() []byte {
 line := paddedPredictLine(inputLimit/64)
 return bytes.Repeat(line,64)
}

func TestCognitionInputByteLimitExactAccepted(t *testing.T) {
 input := fullInputAtLimit()
 if len(input)!=inputLimit {t.Fatalf("size = %d",len(input))}
 var out bytes.Buffer
 if err:=runCognition(bytes.NewReader(input),&out);err!=nil {t.Fatal(err)}
 if got:=bytes.Count(out.Bytes(),[]byte("probability"));got!=64 {t.Fatalf("got %d outputs, want 64",got)}
}

func TestCognitionInputByteLimitPlusOneRejectedWithoutOutput(t *testing.T) {
 input := append(fullInputAtLimit(),byte(' '))
 var out bytes.Buffer
 if err:=runCognition(bytes.NewReader(input),&out);err==nil||!strings.Contains(err.Error(),"exceeds byte limit") {t.Fatalf("err=%v",err)}
 if out.Len()!=0 {t.Fatalf("unexpected output of %d bytes",out.Len())}
}

func TestCognitionInputFullAdditionalEventRejectedWithoutOutput(t *testing.T) {
 input:=append(fullInputAtLimit(), []byte(testPredict+"\n")...)
 var out bytes.Buffer
 if err:=runCognition(bytes.NewReader(input),&out);err==nil||!strings.Contains(err.Error(),"exceeds byte limit") {t.Fatalf("err=%v",err)}
 if out.Len()!=0 {t.Fatalf("unexpected output of %d bytes",out.Len())}
}

func TestCognitionInputScannerTokenLimitStillDenied(t *testing.T) {
 var out bytes.Buffer
 if err:=runCognition(strings.NewReader(strings.Repeat(" ",1<<16)+"\n"),&out);err==nil {t.Fatal("overlong scanner token accepted")}
 if out.Len()!=0 {t.Fatalf("unexpected output of %d bytes",out.Len())}
}

type readerWithFailure struct{ done bool }
func (r *readerWithFailure) Read(b []byte) (int,error) {
 if !r.done {r.done=true;return copy(b,[]byte(testPredict+"\n")),nil}
 return 0,errors.New("injected reader failure")
}
func TestCognitionInputReaderErrorDeniedWithoutOutput(t *testing.T) {
 var out bytes.Buffer
 if err:=runCognition(&readerWithFailure{},&out);err==nil||!strings.Contains(err.Error(),"injected reader failure") {t.Fatalf("err=%v",err)}
 if out.Len()!=0 {t.Fatalf("unexpected output of %d bytes",out.Len())}
}
func TestCognitionInputOrdinaryFlowPreserved(t *testing.T) {
 input := `{"op":"observe","action":1,"prior":0,"outcome":1,"rate":0.1}`+"\n"+`{"op":"forecast_two","prior":0,"action":1,"next_action":0}`+"\n"+`{"op":"snapshot"}`+"\n"
 var out bytes.Buffer
 if err:=runCognition(strings.NewReader(input),&out);err!=nil {t.Fatal(err)}
 if !bytes.Contains(out.Bytes(),[]byte("probability"))||!bytes.Contains(out.Bytes(),[]byte("weights")){t.Fatalf("missing outputs %s",out.String())}
}
var _ io.Reader = (*readerWithFailure)(nil)
