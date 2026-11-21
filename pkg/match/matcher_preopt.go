package match

import (
	"context"
	"fmt"
	"math"
)

// PreOptimize runs the LiveOptimizer + DiffOptimizer pair on the best
// candidate available (sweep seed or cohort best). It replaces the old
// NN-seeded refinement step; it also populates State.LiveOpt, which the
// GA uses for permutation candidate generation.
type PreOptimize struct{}

func (s *PreOptimize) Name() string { return "PreOptimize" }

func (s *PreOptimize) Run(ctx context.Context, in MatcherInput) MatcherResult {
	result := MatcherResult{}
	cfg := in.Config
	target := in.Target
	features := in.Features

	if len(in.State.SeedCandidates) == 0 && in.Cohort.Best == EmptyCandidate {
		return result
	}
	rawTargetRMS := features.TargetRMS
	if rawTargetRMS <= 0 {
		rawTargetRMS = CalcRMS(target)
	}

	var initialPatch *Candidate
	if len(in.State.SeedCandidates) > 0 {
		initialPatch = in.State.SeedCandidates[0]
	}
	if in.Cohort.Best != EmptyCandidate && in.Cohort.Best.Fitness < math.MaxFloat64 &&
		(initialPatch == nil || in.Cohort.Best.Fitness < initialPatch.Fitness) {
		initialPatch = &in.Cohort.Best
	}
	if initialPatch == nil {
		return result
	}

	seedDist := initialPatch.Fitness
	optThreshold := math.MaxFloat64
	if in.Cohort.Best != EmptyCandidate && in.Cohort.Best.Fitness < math.MaxFloat64 {
		optThreshold = in.Cohort.Best.Fitness * (1 + cfg.OptMargin)
	}
	if seedDist > optThreshold {
		if cfg.Verbose {
			fmt.Printf("Skipping pre-optimization: dist=%.6f > threshold=%.6f\n", seedDist, optThreshold)
		}
		return result
	}

	liveOpt := NewLiveOptimizer(ctx, target, initialPatch, 1000, cfg.Verbose, cfg.MIDINote, cfg.RenderDur, rawTargetRMS).WithTargetMel(in.TargetMel)
	in.State.LiveOpt = liveOpt
	refined := liveOpt.Run()
	diffOpt := NewDiffOptimizer(ctx, target, refined, 5, cfg.Verbose, cfg.MIDINote, cfg.RenderDur, rawTargetRMS)
	diffRefined := diffOpt.RunAlternating(2)
	if diffOpt.bestDist < liveOpt.bestDist {
		refined = diffRefined
	}
	refinedDist := liveOpt.bestDist
	if diffOpt.bestDist < refinedDist {
		refinedDist = diffOpt.bestDist
	}

	in.State.LiveOptBest = refined
	result.Cohort.Best = *refined
	result.Cohort.Best.Fitness = refinedDist
	result.Cohort.Best.Scores[ScoreTimbre] = refinedDist
	result.Cohort.Timbral[0] = *refined
	result.Cohort.Timbral[0].Fitness = refinedDist
	result.Cohort.Timbral[0].Scores[ScoreTimbre] = refinedDist
	result.Cohort.TimbralN = 1

	if refinedDist < seedDist && len(in.State.SeedCandidates) > 0 {
		in.State.SeedCandidates[0] = refined
	}
	if refinedDist <= 1e-4 {
		logElapsed(cfg, "Exact match found during pre-optimization (dist=%.6f). Returning early!\n", refinedDist)
		result.EarlyReturn = true
		return result
	}
	return result
}
