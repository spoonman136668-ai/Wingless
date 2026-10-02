package unitary

import "testing"

func TestUP47FrozenTrajectoryAuthority(t *testing.T) {
	if len(up45FrozenDenseTrajectory) != 49 {
		t.Fatalf(
			"trajectory count=%d want=49",
			len(up45FrozenDenseTrajectory),
		)
	}
	for index, state := range up45FrozenDenseTrajectory {
		if state.Step != index {
			t.Fatalf(
				"state[%d].step=%d",
				index, state.Step,
			)
		}
		if len(state.Offsets) != compositeChannels {
			t.Fatalf(
				"step=%d offset count=%d",
				state.Step, len(state.Offsets),
			)
		}
	}

	want0 := []float64{
		-0.07319250547113999, -0.043915503282683996,
		-0.014638501094227999, 0.014638501094227999,
		0.043915503282683996, 0.07319250547113999,
	}
	want24 := []float64{
		-0.05841298929708432, -0.03887820934987345,
		0.0051658188506906524, -0.03702321187042985,
		0.07759973733795059, 0.05154885432874637,
	}
	want29 := []float64{
		-0.045664698831175514, -0.04501676812740249,
		0.007218659931181362, -0.045664698831175514,
		0.07894943527059196, 0.05017807058798021,
	}
	want48 := []float64{
		-0.04050481968956773, -0.03772334675340676,
		0.03174652633686573, -0.0666483907461437,
		0.04995983458779618, 0.06317019626445626,
	}
	for _, pair := range []struct {
		step int
		want []float64
	}{
		{0, want0},
		{24, want24},
		{29, want29},
		{48, want48},
	} {
		got := up45FrozenDenseTrajectory[pair.step].Offsets
		for i := range pair.want {
			if got[i] != pair.want[i] {
				t.Fatalf(
					"step=%d offset[%d]=%.17g want=%.17g",
					pair.step, i, got[i], pair.want[i],
				)
			}
		}
	}

	_, cap24, err := exactOffsetGroups(
		up45FrozenDenseTrajectory[24].Offsets,
	)
	if err != nil {
		t.Fatal(err)
	}
	_, cap29, err := exactOffsetGroups(
		up45FrozenDenseTrajectory[29].Offsets,
	)
	if err != nil {
		t.Fatal(err)
	}
	if cap24 != 6 {
		t.Fatalf("step24 capacity=%d want=6", cap24)
	}
	if cap29 != 8 {
		t.Fatalf("step29 capacity=%d want=8", cap29)
	}
	if up45FrozenDenseTrajectory[29].Offsets[0] !=
		up45FrozenDenseTrajectory[29].Offsets[3] {
		t.Fatal("step29 exact pair [0,3] lost")
	}
}
