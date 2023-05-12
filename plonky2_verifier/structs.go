package plonky2_verifier

import (
	. "gnark-plonky2-verifier/field"
)

type MerkleCap = []Hash // Length = 2^CircuitConfig.FriConfig.CapHeight

type MerkleProof struct {
	Siblings []Hash // Length = CircuitConfig.FriConfig.DegreeBits + CircuitConfig.FriConfig.RateBits - CircuitConfig.FriConfig.CapHeight
}

type EvalProof struct {
	Elements    []F
	MerkleProof MerkleProof
}

type FriInitialTreeProof struct {
	EvalsProofs []EvalProof
}

type FriQueryStep struct {
	Evals       []QuadraticExtension
	MerkleProof MerkleProof
}

type FriQueryRound struct {
	InitialTreesProof FriInitialTreeProof
	Steps             []FriQueryStep
}

type PolynomialCoeffs struct {
	Coeffs []QuadraticExtension
}

type FriProof struct {
	CommitPhaseMerkleCaps []MerkleCap     // Length = Len(CommonCircuitData.FriParams.ReductionArityBits)
	QueryRoundProofs      []FriQueryRound // Length = CommonCircuitData.FriConfig.FriParams.NumQueryRounds
	FinalPoly             PolynomialCoeffs
	PowWitness            F
}

type OpeningSet struct {
	Constants       []QuadraticExtension // Length = CommonCircuitData.Constants
	PlonkSigmas     []QuadraticExtension // Length = CommonCircuitData.NumRoutedWires
	Wires           []QuadraticExtension // Length = CommonCircuitData.NumWires
	PlonkZs         []QuadraticExtension // Length = CommonCircuitData.NumChallenges
	PlonkZsNext     []QuadraticExtension // Length = CommonCircuitData.NumChallenges
	PartialProducts []QuadraticExtension // Length = CommonCircuitData.NumChallenges * CommonCircuitData.NumPartialProducts
	QuotientPolys   []QuadraticExtension // Length = CommonCircuitData.NumChallenges * CommonCircuitData.QuotientDegreeFactor
}

type Proof struct {
	WiresCap                  MerkleCap
	PlonkZsPartialProductsCap MerkleCap
	QuotientPolysCap          MerkleCap
	Openings                  OpeningSet
	OpeningProof              FriProof
}

type ProofWithPublicInputs struct {
	Proof        Proof
	PublicInputs []F // Length = CommonCircuitData.NumPublicInputs
}

type VerifierOnlyCircuitData struct {
	ConstantSigmasCap MerkleCap
	CircuitDigest     Hash
}

type FriConfig struct {
	RateBits        uint64
	CapHeight       uint64
	ProofOfWorkBits uint64
	NumQueryRounds  uint64
	// TODO: add FriReductionStrategy
}

func (fc *FriConfig) rate() float64 {
	return 1.0 / float64((uint64(1) << fc.RateBits))
}

type FriParams struct {
	Config             FriConfig
	Hiding             bool
	DegreeBits         uint64
	ReductionArityBits []uint64
}

type CircuitConfig struct {
	NumWires                uint64
	NumRoutedWires          uint64
	NumConstants            uint64
	UseBaseArithmeticGate   bool
	SecurityBits            uint64
	NumChallenges           uint64
	ZeroKnowledge           bool
	MaxQuotientDegreeFactor uint64
	FriConfig               FriConfig
}

type CommonCircuitData struct {
	Config               CircuitConfig
	FriParams            FriParams
	Gates                []gate
	SelectorsInfo        SelectorsInfo
	DegreeBits           uint64
	QuotientDegreeFactor uint64
	NumGateConstraints   uint64
	NumConstants         uint64
	NumPublicInputs      uint64
	KIs                  []F
	NumPartialProducts   uint64
}

type ProofChallenges struct {
	PlonkBetas    []F
	PlonkGammas   []F
	PlonkAlphas   []F
	PlonkZeta     QuadraticExtension
	FriChallenges FriChallenges
}

type FriChallenges struct {
	FriAlpha        QuadraticExtension
	FriBetas        []QuadraticExtension
	FriPowResponse  F
	FriQueryIndices []F
}
