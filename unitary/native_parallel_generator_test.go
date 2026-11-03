package unitary

import (
	"bytes"
	"encoding/json"
	"testing"
)

func nativeGeneratorExamples() []NativeTextExample {
	return []NativeTextExample{{"aaaa", "aA"}, {"zzzz", "bB"}}
}
func TestNativeGeneratorVariableLengthUTF8(t *testing.T) {
	xs := []NativeTextExample{{"aaaa", "é"}, {"zzzz", "A"}}
	m, err := TrainNativeParallelGenerator(xs)
	if err != nil {
		t.Fatal(err)
	}
	if m.length.Kind != "native" {
		t.Fatal("length field did not train")
	}
	raw, _ := m.Checkpoint()
	r, err := RestoreNativeParallelGenerator(raw)
	if err != nil {
		t.Fatal(err)
	}
	for _, x := range xs {
		p, err := r.Generate(x.Input)
		if err != nil || !p.ValidOutput || p.Response == nil || *p.Response != x.Response || p.PredictedLength != len(x.Response) {
			t.Fatalf("variable UTF8 generation failed %+v %v", p, err)
		}
	}
}
func TestNativeGeneratorFieldsFitAndCompleteState(t *testing.T) {
	m, err := TrainNativeParallelGenerator(nativeGeneratorExamples())
	if err != nil {
		t.Fatal(err)
	}
	raw, err := m.Checkpoint()
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(raw, []byte(`"responses"`)) || bytes.Contains(raw, []byte(`"aA"`)) {
		t.Fatal("retained whole response bank")
	}
	for _, x := range nativeGeneratorExamples() {
		p, err := m.Generate(x.Input)
		if err != nil || p.Response == nil || *p.Response != x.Response || !p.SeenTrainingResponse || p.WholeResponseBank || p.LanguageProof {
			t.Fatalf("engineering pair generation failed: %+v %v", p, err)
		}
	}
	after, _ := m.Checkpoint()
	if !bytes.Equal(raw, after) {
		t.Fatal("generation mutated state")
	}
	r, err := RestoreNativeParallelGenerator(raw)
	if err != nil {
		t.Fatal(err)
	}
	again, _ := r.Checkpoint()
	if !bytes.Equal(raw, again) {
		t.Fatal("complete field snapshot changed")
	}
	a, _ := m.Generate("aaaa")
	b, _ := r.Generate("aaaa")
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	if !bytes.Equal(ja, jb) {
		t.Fatal("restored fields changed")
	}
}
func TestNativeGeneratorMechanicalCompositionNotTransfer(t *testing.T) {
	m, err := TrainNativeParallelGenerator(nativeGeneratorExamples())
	if err != nil {
		t.Fatal(err)
	}
	m.fields[0] = nativeGeneratorField{Kind: "constant", Value: 0}
	p, err := m.Generate("aaaa")
	if err != nil || p.Response == nil || *p.Response != "`A" || p.SeenTrainingResponse || p.LanguageProof {
		t.Fatalf("mechanical field intervention failed %+v %v", p, err)
	}
	// This forced intervention demonstrates assembly only, not learned novelty.
	for bit := 0; bit < 8; bit++ {
		m.fields[bit] = nativeGeneratorField{Kind: "constant", Value: 1}
	}
	bad, err := m.Generate("aaaa")
	if err != nil || bad.ValidOutput || bad.Response != nil || bad.RawHex != "ff41" {
		t.Fatalf("invalid native candidate not preserved %+v %v", bad, err)
	}
}
func TestNativeGeneratorBoundsAndCorruption(t *testing.T) {
	for _, xs := range [][]NativeTextExample{nil, {{"x", "same"}, {"y", "same"}}, {{"", "a"}, {"y", "b"}}} {
		if _, err := TrainNativeParallelGenerator(xs); err == nil {
			t.Fatal("invalid training accepted")
		}
	}
	m, _ := TrainNativeParallelGenerator(nativeGeneratorExamples())
	raw, _ := m.Checkpoint()
	var e nativeGeneratorEnvelope
	json.Unmarshal(raw, &e)
	e.Payload.SelectorBlob = "foreign"
	p, _ := json.Marshal(e.Payload)
	e.SHA256 = generatorDigest(p)
	bad, _ := json.Marshal(e)
	if _, err := RestoreNativeParallelGenerator(bad); err == nil {
		t.Fatal("foreign source accepted")
	}
	if _, err := RestoreNativeParallelGenerator(append(raw, []byte("{}")...)); err == nil {
		t.Fatal("trailing JSON accepted")
	}
	before, _ := m.Checkpoint()
	if _, err := m.Generate(""); err == nil {
		t.Fatal("empty input accepted")
	}
	after, _ := m.Checkpoint()
	if !bytes.Equal(before, after) {
		t.Fatal("invalid input mutated state")
	}
}
