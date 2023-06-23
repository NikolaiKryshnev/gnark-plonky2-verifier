package gates

import (
	"regexp"

	"github.com/consensys/gnark/frontend"
	"github.com/succinctlabs/gnark-plonky2-verifier/field"
	"github.com/succinctlabs/gnark-plonky2-verifier/gl"
	"github.com/succinctlabs/gnark-plonky2-verifier/poseidon"
)

var poseidonMdsGateRegex = regexp.MustCompile("PoseidonMdsGate.*")

func deserializePoseidonMdsGate(parameters map[string]string) Gate {
	// Has the format "PoseidonMdsGate(PhantomData<plonky2_field::goldilocks_field::GoldilocksField>)<WIDTH=12>"
	return NewPoseidonMdsGate()
}

type PoseidonMdsGate struct {
}

func NewPoseidonMdsGate() *PoseidonMdsGate {
	return &PoseidonMdsGate{}
}

func (g *PoseidonMdsGate) Id() string {
	return "PoseidonMdsGate"
}

func (g *PoseidonMdsGate) WireInput(i uint64) Range {
	if i >= poseidon.SPONGE_WIDTH {
		panic("Input less than sponge width")
	}
	return Range{i * field.D, (i + 1) * field.D}
}

func (g *PoseidonMdsGate) WireOutput(i uint64) Range {
	if i >= poseidon.SPONGE_WIDTH {
		panic("Input less than sponge width")
	}
	return Range{(poseidon.SPONGE_WIDTH + i) * field.D, (poseidon.SPONGE_WIDTH + i + 1) * field.D}
}

func (g *PoseidonMdsGate) mdsRowShfAlgebra(
	r uint64,
	v [poseidon.SPONGE_WIDTH]gl.QuadraticExtensionAlgebraVariable,
	api frontend.API,
) gl.QuadraticExtensionAlgebraVariable {
	glApi := gl.NewChip(api)
	if r >= poseidon.SPONGE_WIDTH {
		panic("MDS row index out of range")
	}

	res := gl.ZeroExtensionAlgebra()
	for i := uint64(0); i < poseidon.SPONGE_WIDTH; i++ {
		coeff := gl.NewQuadraticExtensionVariable(gl.NewVariable(poseidon.MDS_MATRIX_CIRC[i]), gl.Zero())
		res = glApi.AddExtensionAlgebra(res, glApi.ScalarMulExtensionAlgebra(coeff, v[(i+r)%poseidon.SPONGE_WIDTH]))
	}

	coeff := gl.NewQuadraticExtensionVariable(gl.NewVariable(poseidon.MDS_MATRIX_DIAG[r]), gl.Zero())
	res = glApi.AddExtensionAlgebra(res, glApi.ScalarMulExtensionAlgebra(coeff, v[r]))

	return res
}

func (g *PoseidonMdsGate) mdsLayerAlgebra(
	state [poseidon.SPONGE_WIDTH]gl.QuadraticExtensionAlgebraVariable,
	api frontend.API,
) [poseidon.SPONGE_WIDTH]gl.QuadraticExtensionAlgebraVariable {
	var result [poseidon.SPONGE_WIDTH]gl.QuadraticExtensionAlgebraVariable
	for r := uint64(0); r < poseidon.SPONGE_WIDTH; r++ {
		result[r] = g.mdsRowShfAlgebra(r, state, api)
	}
	return result
}

func (g *PoseidonMdsGate) EvalUnfiltered(
	api frontend.API,
	vars EvaluationVars,
) []gl.QuadraticExtensionVariable {
	glApi := gl.NewChip(api)
	constraints := []gl.QuadraticExtensionVariable{}

	var inputs [poseidon.SPONGE_WIDTH]gl.QuadraticExtensionAlgebraVariable
	for i := uint64(0); i < poseidon.SPONGE_WIDTH; i++ {
		inputs[i] = vars.GetLocalExtAlgebra(g.WireInput(i))
	}

	computed_outputs := g.mdsLayerAlgebra(inputs, api)

	for i := uint64(0); i < poseidon.SPONGE_WIDTH; i++ {
		output := vars.GetLocalExtAlgebra(g.WireOutput(i))
		diff := glApi.SubExtensionAlgebra(output, computed_outputs[i])
		constraints = append(constraints, diff[0], diff[1])
	}

	return constraints
}
