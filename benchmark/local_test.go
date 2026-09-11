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
func TestSemanticNeverOverridesStrict(t *testing.T) {
	f := jsonFixture("x", "p", `{"sum":5}`)
	s := "```json\n{\"sum\":5}\n```"
	if !f.check(semanticCandidate(s)) || f.protocol(s) {
		t.Fatal("scores conflated")
	}
	if f.protocol(`{"sum":5,"sum":4}`) {
		t.Fatal("duplicate keys")
	}
	if !f.protocol(`{"sum":4}`) || f.check(`{"sum":4}`) {
		t.Fatal("protocol/meaning conflated")
	}
}
