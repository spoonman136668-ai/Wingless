package unitary

import "testing"

func TestUP46UsesFrozenUP45States(t *testing.T) {
	if len(up45SelectedStep24Offsets) != 6 ||
		len(up45ExactFusionStep29Offsets) != 6 {
		t.Fatalf("UP45 frozen states lost")
	}
	_, selectedCapacity, err :=
		exactOffsetGroups(up45SelectedStep24Offsets)
	if err != nil {
		t.Fatal(err)
	}
	_, fusedCapacity, err :=
		exactOffsetGroups(up45ExactFusionStep29Offsets)
	if err != nil {
		t.Fatal(err)
	}
	if selectedCapacity != 6 {
		t.Fatalf("selected capacity=%d want=6", selectedCapacity)
	}
	if fusedCapacity != 8 {
		t.Fatalf("fused capacity=%d want=8", fusedCapacity)
	}
	if up45ExactFusionStep29Offsets[0] !=
		up45ExactFusionStep29Offsets[3] {
		t.Fatal("UP45 exact fused pair [0,3] lost")
	}
}

func TestUP46DoesNotRunOptimizerOrChangeLambda(t *testing.T) {
	if proximalFusionLambda != 0.00004 {
		t.Fatalf("proximal lambda=%g want=0.00004", proximalFusionLambda)
	}
}
