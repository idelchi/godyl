package match

// Assets is a slice of Asset, representing a collection of downloadable files.
type Assets []Asset

// FromNames creates a collection of assets from the provided names.
// It initializes each asset with the given name.
func (as Assets) FromNames(names ...string) Assets {
	assets := make(Assets, len(names))

	for i, name := range names {
		assets[i] = Asset{Name: name}
	}

	return assets
}

// Select filters and sorts the assets based on the provided requirements.
// It returns the best matching assets, sorted by score, if any are qualified.
func (as Assets) Select(req Requirements) Results {
	results := as.Match(req)

	if results.HasErrors() {
		return results
	}

	if !results.HasQualified() {
		return results
	}

	return results.Best().Sorted()
}

// Match evaluates all assets against the provided requirements and returns a list of results.
// Each result contains the asset's name, its matching score, and whether it qualifies.
func (as Assets) Match(req Requirements) Results {
	var results Results

	for _, a := range as {
		score, qualified, err := a.Match(req)
		results = append(results, Result{Asset: a, Score: score, Qualified: qualified, Error: err})
	}

	return results
}
