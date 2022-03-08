package cmd

import (
	"testing"
)

func Test_parseHostedGitRegex_many(t *testing.T) {
	var origins = []string{
		"git@gitlab.com:gsei19/nineteen-week/week-02-full-stack-by-feature.git",
		"https://gitlab.com/gsei19/week-05.git",
		"git@github.com:golang/from/go.git",
		"ssh://user@host.xz/path/to/repo.git/",
		"ssh://user@host-website.xz:1234/path-path/to-to/repo-repo.git/",
		"ssh://host.xz:7634/path/to/repo-name-nam.git/",
		"ssh://host.xz/path/to/repo-name.git/",
		"ssh://user@host.xz/path/to/repo-name.git/",
		"ssh://host.xz/path/to/repo-name.git/",
		"user@host.xz:/path/to/repo.git/",
		"host.xz:/path/to/repo-name.git/",
		"user@host.xz:path/to/repo.git",
		"host.xz:path/to/repo.git",
		"rsync://host.xz/path/to/repo.git/",
		"git://host.xz/path/to/repo.git/",
		"http://host.xz/path/to/repo.git/",
		"https://host.xz/path/to/repo.git/",
	}
	for _, origin := range origins {
		var repoPieces, _ = parseHostedGitRegex(origin)
		if repoPieces.Org == "" {
			t.Errorf("Org parse failed")
		}
	}
}

func Test_parseRepoPieces_TopLevelOrg(t *testing.T) {
	repoPieces, err := parseHostedGitRegex("https://gitlab.com/gsei19/week-05.git")

	if err != nil {
		t.Errorf("Threw an error")
	}

	// https://gitlab.com/gsei19/nineteen-week/week-05.git

	if repoPieces.Org != "gsei19" {
		t.Errorf("Incorrect org parse")
	}
}

func Test_parseRepoPieces_TopLevelOrg_Ssh(t *testing.T) {
	repoPieces, err := parseHostedGitRegex("git@github.com:golang/go.git")

	if err != nil {
		t.Errorf("Threw an error")
	}

	// https://gitlab.com/gsei19/nineteen-week/week-05.git

	if repoPieces.Origin != "github.com" {
		t.Errorf("Incorrect Origin parse")
	}

	if repoPieces.RepoName != "go" {
		t.Errorf("Incorrect RepoName parse")
	}

	if repoPieces.Org != "golang" {
		t.Errorf("Incorrect org parse")
	}
}

func Test_parseRepoPieces_NestedOrg(t *testing.T) {
	repoPieces, err := parseHostedGitRegex("https://gitlab.com/gsei19/nineteen-week/week-05.git")

	if err != nil {
		t.Errorf("Threw an error")
	}

	// https://gitlab.com/gsei19/nineteen-week/week-05.git

	if repoPieces.Org != "gsei19/nineteen-week" {
		t.Errorf("Incorrect org parse")
	}
}
