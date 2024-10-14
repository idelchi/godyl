package match

type Assets []Asset

func (as Assets) FromNames(names ...string) Assets {
	assets := make(Assets, len(names))

	for i, name := range names {
		assets[i] = Asset{Name: name}
	}

	return assets
}

func (as Assets) Select(req Requirements) Results {
	results := as.Match(req)

	if !results.HasQualified() {
		return results
	}

	return results.Best().Sorted()
}

// Results evaluates all assets against requirements and returns Results.
func (as Assets) Match(req Requirements) Results {
	var results Results
	for _, a := range as {
		score, qualified := a.Match(req)
		results = append(results, Result{Name: a.Name, Score: score, Qualified: qualified})
	}

	return results
}
