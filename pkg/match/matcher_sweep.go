package match

import (
	"context"
	"fmt"
	"math"
	"sort"
)

// AlgoSweep deterministically re-topologizes the best candidate produced by
// earlier steps across all 32 DX7 algorithm topologies, evaluates every
// variant, and keeps the top seeds for downstream steps (GA islands,
// formant refinement, top-32 export).
type AlgoSweep struct{}

func (s *AlgoSweep) Name() string { return "AlgoSweep" }

func (s *AlgoSweep) Run(ctx context.Context, in MatcherInput) MatcherResult {
	result := MatcherResult{}
	cfg := in.Config
	target := in.Target
	features := in.Features

	rawTargetRMS := features.TargetRMS
	if rawTargetRMS <= 0 {
		rawTargetRMS = CalcRMS(target)
	}

	// Base parameter set: best candidate from earlier steps, else a
	// deterministic heuristic candidate built from the target features.
	var base *Candidate
	if in.Cohort.Best != EmptyCandidate && in.Cohort.Best.Fitness < math.MaxFloat64 {
		base = in.Cohort.Best.Copy()
	} else {
		base = HeuristicCandidateWithAlgo(features, cfg.NoteFreq, 0)
	}

	type algoScore struct {
		algo uint8
		dist float64
		cand *Candidate
	}
	candidates := make([]*Candidate, 0, 32)
	for algo := uint8(0); algo < 32; algo++ {
		cand := base.Copy()
		cand.Algorithm = algo
		cand.Fitness = math.MaxFloat64
		if in.TimbreCB != nil {
			cand = AttachNearestRawBaseTimbre(cand, in.TimbreCB)
		} else if in.Codebook != nil {
			cand = AttachNearestRawBase(cand, in.Codebook)
		}
		candidates = append(candidates, cand)
	}

	fits, _ := evaluateBatchParallel(ctx, candidates, target, cfg.MIDINote, cfg.RenderDur, rawTargetRMS, in.TargetMel)

	scores := make([]algoScore, 0, 32)
	for i, cand := range candidates {
		cand.Fitness = fits[i]
		cand.Scores[ScoreTimbre] = fits[i]
		scores = append(scores, algoScore{cand.Algorithm, fits[i], cand})
	}
	sort.Slice(scores, func(i, j int) bool { return scores[i].dist < scores[j].dist })

	// Top-4 sweep results become the shared seed list for downstream steps,
	// plus the best one into the Timbral cohort.
	seeds := make([]*Candidate, 0, 4)
	for i := 0; i < 4 && i < len(scores); i++ {
		seeds = append(seeds, scores[i].cand.Copy())
	}
	in.State.SeedCandidates = seeds

	if len(scores) > 0 {
		result.Cohort.Best = *scores[0].cand
		result.Cohort.Best.Fitness = scores[0].dist
		result.Cohort.Timbral[0] = *scores[0].cand
		result.Cohort.Timbral[0].Fitness = scores[0].dist
		result.Cohort.Timbral[0].Scores[ScoreTimbre] = scores[0].dist
		result.Cohort.TimbralN = 1
	}
	if cfg.Verbose {
		fmt.Printf("Algorithm sweep: best algo=%d dist=%.6f\n", scores[0].algo, scores[0].dist)
	}
	if scores[0].dist <= 1e-4 {
		logElapsed(cfg, "Exact sweep match found (dist=%.6f). Returning early!\n", scores[0].dist)
		result.EarlyReturn = true
		return result
	}
	return result
}
