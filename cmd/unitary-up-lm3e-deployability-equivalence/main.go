package main

import("encoding/json";"os";"github.com/spoonman136668-ai/Wingless/unitary")

func main(){r,_:=unitary.RunUPLM3E();e:=json.NewEncoder(os.Stdout);e.SetIndent("","  ");_=e.Encode(r)}
