# fmguessr

DX7 FM synthesizer patch matcher. Uses genetic algorithms, formant analysis, and hill-climbing to match target audio WAV files and generate DX7 `.syx` bank files.

**Browser frontend:** [OpenPatches](https://github.com/kompadre/openpatches) (live at [openpatch.es](https://openpatch.es))

## Features

- **Matcher Pipeline:** Modular 9-step pipeline — each step is a `Matcher` interface with early-return support, reconciled into a shared `ResultsCohort` after each step.
- **Formant Candidate Generation:** Extracts spectral envelope peaks (F1/F2/F3) from target, generates DX7 candidates whose FM sidebands match. Zero-render, pure math.
- **Formant-Constrained Refinement:** Optimizes candidates on formant distance while keeping mel distance within 10% tolerance.
- **Algorithm Sweep:** Deterministically re-topologizes the best candidate across all 32 DX7 algorithms, evaluates each, and keeps the top seeds.
- **Preset Codebook:** Dual KNN lookup (fingerprint + timbre) for instant GA seeding.
- **32-Island GA:** Niche evolution across algorithm topologies → global convergence → settle phase. Seeded from codebook, sweep, beam search, and LiveOptimizer results.
- **LiveOptimizer:** Multi-threaded hill-climbing with acceleration trough cacher.
- **DiffOptimizer:** Jacobian steering on temporal timbre features.
- **Beam Search:** Constructive operator-by-operator patch building.
- **CSPOA:** Codebook-seeded progressive operator attachment.
- **Template Preservation:** LFO, pitch envelopes, and scaling from matched presets preserved.

## Build

```bash
make          # Build all tools
make test     # Run all tests
```

Requires Go 1.21+, C++ compiler (MSFA engine).

## Testing

```bash
make test               # Unit + E2E tests
make test-smoke         # Fast e2e smoke tests (no external fixtures)
make test-regression    # Regression tests (-tags regression)
go test -v ./pkg/match/...           # Unit tests only
go test -v -tags regression ./tests/...  # Regression tests only
go test -race ./pkg/match/...        # Race detector
```

## Usage

Flags must come before the WAV path.

### Match a Target

```bash
# Standard GA matching
./match -v target.wav -syx result.syx
```

### Play/Export

```bash
./play result.syx 0 C3 -o match.wav
```

### Index a Codebook

```bash
./index Dexed_01.syx
```

### REST API

```bash
./server -port 8080 -codebook dx7_codebook.bin
```

- `POST /api/probe` — Analyze WAV envelope/pitch/formants (truncates to 2s centered on energy peak, ~5s). Returns `read_ms`, `decode_ms`, `probe_ms`, `peak_f0_hz`.
- `POST /api/match` — Async matching (returns job ID, poll `GET /api/status/:id`)
- `POST /api/play` — Render patch to WAV
- `POST /api/voice` — Extract voice data from .syx
- `POST /api/morph` — Morph two patches with spectral optimization
- `POST /api/export-syx` — Export patches as .syx bank
- `GET /api/health` — Server health check
- `GET /exports/*` — Static files (server-generated only, no user uploads stored)

See [API.md in openpatches](https://github.com/kompadre/openpatches/blob/main/API.md) for full request/response schemas.

## Architecture

```
target.wav
  ├─→ Matcher Pipeline (9 steps, each reconciled into ResultsCohort):
  │     1. CodebookSearch    — dual KNN (fingerprint + timbre)
  │     2. FormantGenerate   — spectral envelope → DX7 candidates
  │     3. AlgoSweep         — deterministic sweep across 32 algorithms
  │     4. PreOptimize       — LiveOptimizer + DiffOptimizer on best seed
  │     5. BeamSearch        — constructive operator-by-operator
  │     6. CSPOAStep         — progressive operator attachment
  │     7. GAStep            — 32-island niche → global convergence
  │     8. HillClimbStep     — LiveOptimizer + DiffOptimizer on GA best
  │     9. FormantRefineStep — formant-constrained refinement
  ├─→ Finalize: rescue → top-32 assembly → carrier calibration → temporal re-eval
  └─→ DX7 .syx export (top-32 candidates from diverse sources)
```

## Result Buckets

The matcher produces 5 result categories:

- **Top 32** — Diverse candidates from island bests + formant + sweep seeds + GA (for .syx export)
- **TimbralTop5** — Best by static timbre distance
- **FullTop5** — Best by temporal distance (DiffOptimizer-refined)
- **FormantTop5** — Best by formant-constrained refinement
- **Favorites** — Candidates appearing in both TimbralTop5 and FullTop5

## Key Types

- `Candidate` — DX7 patch with `RawBase [128]byte`, `Scores [NumScoreTypes]float64`, `Fitness`
- `ResultsCohort` — Best + Timbral[5] + Full[5] + Formant[5] (fixed-size arrays)
- `CandidatePool` — 1024-slot pre-allocated pool (zero-allocation)
- `Matcher` interface — `Name() string` + `Run(ctx, MatcherInput) MatcherResult`

## Audio Tour

`LONGREAD.txt` contains a spoken-word architecture walkthrough suitable for text-to-speech. Covers the full signal flow from WAV input through the 9-step matcher pipeline to DX7 .syx output.

## License

MIT
