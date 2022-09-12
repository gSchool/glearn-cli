package cmd

import (
	"github.com/gSchool/glearn-cli/api/learn"
	"testing"
)

func Test_parseHostedGitRegex_many(t *testing.T) {
	type urlToRepo struct {
		url    string
		result learn.RepoPieces
	}
	var origins = []urlToRepo{
		{
			url:    "git@gitlab.com:gsei19/nineteen-week/week-02-full-stack-by-feature.git",
			result: learn.RepoPieces{Origin: "gitlab.com", Org: "gsei19", RepoName: "nineteen-week/week-02-full-stack-by-feature"},
		},
		{
			url:    "https://gitlab.com/gsei19/week-05.git",
			result: learn.RepoPieces{Origin: "gitlab.com", Org: "gsei19", RepoName: "nineteen-week/week-02-full-stack-by-feature"},
		},
		{
			url:    "https://gitlab.com/gsei19/week-05.git",
			result: learn.RepoPieces{Origin: "gitlab.com", Org: "gsei19", RepoName: "week-05.git"},
		},
		{
			url:    "git@github.com:golang/from/go.git",
			result: learn.RepoPieces{Origin: "github.com", Org: "golang", RepoName: "from/go.git"},
		},
		{
			url:    "ssh://user@host.xz/path/to/repo.git/",
			result: learn.RepoPieces{Origin: "host.xz", Org: "path", RepoName: "to/repo"},
		},
		{
			url:    "ssh://user@host-website.xz:1234/path-path/to-to/repo-repo.git/",
			result: learn.RepoPieces{Origin: "host-website.xz", Org: "path-path", RepoName: "to-to/repo-repo"},
		},
		{
			url:    "ssh://host.xz:7634/path/to/repo-name-nam.git/",
			result: learn.RepoPieces{Origin: "host.xz", Org: "path", RepoName: "to/repo-name-nam.git"},
		},
		{
			url:    "ssh://host.xz/path/to/repo-name.git/",
			result: learn.RepoPieces{Origin: "host.xz", Org: "path", RepoName: "to/repo-name"},
		},
		{
			url:    "ssh://user@host.xz/path/to/repo-name.git/",
			result: learn.RepoPieces{Origin: "host.xz", Org: "path", RepoName: "to/repo-name"},
		},
		{
			url:    "ssh://host.xz/path/to/repo-name.git/",
			result: learn.RepoPieces{Origin: "host.xz", Org: "path", RepoName: "to/repo-name"},
		},
		{
			url:    "user@host.xz:/path/to/repo.git/",
			result: learn.RepoPieces{Origin: "host.xz", Org: "path", RepoName: "to/repo"},
		},
		{
			url:    "host.xz:/path/to/repo-name.git/",
			result: learn.RepoPieces{Origin: "host.xz", Org: "path", RepoName: "to/repo-name"},
		},
		{
			url:    "user@host.xz:path/to/repo.git",
			result: learn.RepoPieces{Origin: "host.xz", Org: "path", RepoName: "to/repo"},
		},
		{
			url:    "host.xz:path/to/repo.git",
			result: learn.RepoPieces{Origin: "host.xz", Org: "path", RepoName: "to/repo"},
		},
		{
			url:    "rsync://host.xz/path/to/repo.git/",
			result: learn.RepoPieces{Origin: "host.xz", Org: "path", RepoName: "to/repo"},
		},
		{
			url:    "git://host.xz/path/to/repo.git/",
			result: learn.RepoPieces{Origin: "host.xz", Org: "path", RepoName: "to/repo"},
		},
		{
			url:    "http://host.xz/path/to/repo.git/",
			result: learn.RepoPieces{Origin: "host.xz", Org: "path", RepoName: "to/repo"},
		},
		{
			url:    "https://host.xz/path/to/repo.git/",
			result: learn.RepoPieces{Origin: "host.xz", Org: "path", RepoName: "to/repo"},
		},
		{
			url:    "https://oauth2:personaltoken@gitlab.com/galvanize-labs/curriculum-development/agile.git",
			result: learn.RepoPieces{Origin: "gitlab.com", Org: "galvanize-labs", RepoName: "curriculum-development/agile"},
		},
		{
			url:    "https://hello.dolly@gitlab.com/galvanize-labs/curriculum-development/spring/spring-web.git",
			result: learn.RepoPieces{Origin: "gitlab.com", Org: "galvanize-labs", RepoName: "curriculum-development/spring/spring-web"},
		},
		{
			url:    "ssh://git@gitlab.com/galvanize-labs/curriculum-development/application-architectures/12-factor-app/00-introduction.git",
			result: learn.RepoPieces{Origin: "gitlab.com", Org: "galvanize-labs", RepoName: "curriculum-development/application-architectures/12-factor-app/00-introduction"},
		},
	}

	for _, expected := range origins {
		var test, _ = parseHostedGitRegex(expected.url)
		if test.Org != expected.result.Org {
			t.Errorf("Org parse failed, given:  %s\n expect: %s\n result: %s\n", expected.url, expected.result.Org, test.Org)
			return
		}
		if test.Origin != expected.result.Origin {
			t.Errorf("Origin parse failed, given:  %s\n expect: %s\n result: %s\n", expected.url, expected.result.Origin, test.Origin)
			return
		}
		if test.RepoName != expected.result.RepoName {
			t.Errorf("RepoName parse failed, given:  %s\n expect: %s\n result: %s\n", expected.url, expected.result.RepoName, test.RepoName)
			return
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
