package benchmark

import "testing"

func TestCodeFixtureIsDeterministic(t *testing.T) {
	for _, s := range []string{"package candidate; func Add(a,b int)int{return a+b}", "package candidate; func Add(a int,b int)int{return a+b}"} {
		if !codeShape(s) {
			t.Fatal(s)
		}
	}
	for _, s := range []string{"accepted=true", "package candidate;func Add(a,b string)string{return a+b}", "package candidate;func Add(a,b int)int{return a-b}", "package candidate;func Add(a,b int)int{return a+b};func Evil(){}"} {
		if codeShape(s) {
			t.Fatal(s)
		}
	}
}
