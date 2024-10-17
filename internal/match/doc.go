// Package match provides functionality to evaluate and match assets against
// specific platform requirements and name-based hints. It enables scoring
// assets based on their compatibility with target platforms (e.g., OS, architecture)
// and their matching patterns, which can include regex or string comparisons.
//
// The core types in this package include Asset, which represents a downloadable
// file or package, and Requirements, which defines the criteria for matching
// assets. These criteria can include platform specifications and pattern hints.
// The Results type facilitates sorting and evaluating the best-matching assets.

package match
