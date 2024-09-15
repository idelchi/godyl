// Package gitlab provides functionality to interact with GitLab repositories,
// releases, and assets via the GitLab API. It includes utilities for filtering
// and retrieving release information, assets, and repository metadata.
//
// The main types and functions in this package include:
//
//   - Asset: Represents a GitLab release asset with name, URL, and content type.
//   - Assets: A collection of Asset objects with filtering methods.
//   - Release: Represents a GitLab release containing a tag, name, and assets.
//   - Repository: Represents a GitLab repository, with methods for retrieving
//     releases, assets, and repository details.
//   - NewClient: Creates a new GitLab API client.
//   - NewRepository: Creates a new Repository instance for accessing repository data.
//
// This package uses the official GitLab client (gitlab.com/gitlab-org/api/client-go) for
// making API calls and provides additional utilities to handle GitLab data.
package gitlab
