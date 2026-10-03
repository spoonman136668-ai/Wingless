package main

import(
	"encoding/json"
	"flag"
	"os"

	"github.com/spoonman136668-ai/Wingless/unitary"
)

func readContextAttribution(path string)[]byte{
	data,err:=os.ReadFile(path)
	if err!=nil{panic(err)}
	return data
}

func main(){
	codePath:=flag.String("code","","verified code source path")
	structuredPath:=flag.String("structured","","verified structured-data source path")
	prosePath:=flag.String("technical-prose","","verified technical-prose source path")
	flag.Parse()
	if *codePath==""||*structuredPath==""||*prosePath==""{panic("SOURCE_PATH_REQUIRED")}
	result:=unitary.RunWlmLmExternalContextStateAttributionR1(
		readContextAttribution(*codePath),
		readContextAttribution(*structuredPath),
		readContextAttribution(*prosePath),
	)
	if err:=json.NewEncoder(os.Stdout).Encode(result);err!=nil{panic(err)}
}
