package match

// AttachNearestRawBase finds the closest codebook entry and copies its RawBase to the candidate.
func AttachNearestRawBase(c *Candidate, codebook *Codebook) *Candidate {
	if codebook == nil || len(codebook.Entries) == 0 {
		return c
	}

	bestDist := 999999.0
	var bestEntry *CodebookEntry

	for i := range codebook.Entries {
		entry := &codebook.Entries[i]
		dist := 0.0

		if entry.Candidate.Algorithm != c.Algorithm {
			dist += 100.0
		}

		for op := 0; op < 6; op++ {
			d := float64(entry.Candidate.Ops[op].FreqCoarse) - float64(c.Ops[op].FreqCoarse)
			dist += d * d
		}

		if dist < bestDist {
			bestDist = dist
			bestEntry = entry
		}
	}

	if bestEntry != nil && bestEntry.Candidate.RawBase != NoRawBase {
		c2 := c.Copy()
		c2.RawBase = bestEntry.Candidate.RawBase // value copy
		c2.CopyEnvelopeFromRawBase()
		return c2
	}

	return c
}

// AttachNearestRawBaseTimbre finds the closest timbre codebook entry and copies its RawBase.
func AttachNearestRawBaseTimbre(c *Candidate, tc *TimbreCodebook) *Candidate {
	if tc == nil || len(tc.Entries) == 0 {
		return c
	}

	bestDist := 999999.0
	var bestEntry *TimbreCodebookEntry

	for i := range tc.Entries {
		entry := &tc.Entries[i]
		dist := 0.0

		if entry.Candidate.Algorithm != c.Algorithm {
			dist += 100.0
		}

		for op := 0; op < 6; op++ {
			d := float64(entry.Candidate.Ops[op].FreqCoarse) - float64(c.Ops[op].FreqCoarse)
			dist += d * d
		}

		if dist < bestDist {
			bestDist = dist
			bestEntry = &tc.Entries[i]
		}
	}

	if bestEntry != nil && bestEntry.Candidate.RawBase != NoRawBase {
		c2 := c.Copy()
		c2.RawBase = bestEntry.Candidate.RawBase // value copy
		c2.CopyEnvelopeFromRawBase()
		return c2
	}

	return c
}
