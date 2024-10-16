// Package github provides functionality to interact with GitHub repositories,
// releases, and assets via the GitHub API. It includes utilities for filtering
// and retrieving release information, assets, and repository metadata.
//
// The main types and functions in this package include:
//
//   - Asset: Represents a GitHub release asset with name, URL, and content type.
//   - Assets: A collection of Asset objects with filtering methods.
//   - Release: Represents a GitHub release containing a tag, name, and assets.
//   - Repository: Represents a GitHub repository, with methods for retrieving
//     releases, assets, and repository details.
//   - NewClient: Creates a new GitHub API client.
//   - NewRepository: Creates a new Repository instance for accessing repository data.
//
// This package uses the Google GitHub client (github.com/google/go-github/v64) for
// making API calls and provides additional utilities to handle GitHub data.
package github
