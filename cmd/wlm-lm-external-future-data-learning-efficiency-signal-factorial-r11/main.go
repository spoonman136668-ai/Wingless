package main

import (
  "encoding/json"
  "flag"
  "os"

  "github.com/spoonman136668-ai/Wingless/unitary"
)

func read(path string) []byte {
  b,err:=os.ReadFile(path)
  if err!=nil { panic(err) }
  return b
}

func main() {
  code:=flag.String("code","","")
  structured:=flag.String("structured","","")
  prose:=flag.String("technical-prose","","")
  flag.Parse()
  if *code==""||*structured==""||*prose=="" { panic("SOURCE_PATH_REQUIRED") }
  if err:=json.NewEncoder(os.Stdout).Encode(
    unitary.RunWlmLmExternalFutureDataLearningEfficiencySignalFactorialR11(
      read(*code),read(*structured),read(*prose),
    ),
  ); err!=nil { panic(err) }
}
