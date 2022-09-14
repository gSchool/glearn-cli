package cmd

import (
	"github.com/gSchool/glearn-cli/api/learn"
	"testing"
)

func Test_parseHostedGit_many(t *testing.T) {
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
			result: learn.RepoPieces{Origin: "gitlab.com", Org: "gsei19", RepoName: "week-05"},
		},
		{
			url:    "https://gitlab.com/gsei19/week-05.git",
			result: learn.RepoPieces{Origin: "gitlab.com", Org: "gsei19", RepoName: "week-05"},
		},
		{
			url:    "git@github.com:golang/from/go.git",
			result: learn.RepoPieces{Origin: "github.com", Org: "golang", RepoName: "from/go"},
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
			result: learn.RepoPieces{Origin: "host.xz", Org: "path", RepoName: "to/repo-name-nam"},
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
		var test, err = parseHostedGit(expected.url)
		if err != nil {
			t.Errorf("given %s returned error: %s\n", expected.url, err)
			return
		}
		if test.Org != expected.result.Org {
			t.Errorf("Org parse failed!\ngiven:  %s\n expect: %s\n result: %s\n", expected.url, expected.result.Org, test.Org)
			return
		}
		if test.Origin != expected.result.Origin {
			t.Errorf("Origin parse failed!\n given:  %s\n expect: %s\n result: %s\n", expected.url, expected.result.Origin, test.Origin)
			return
		}
		if test.RepoName != expected.result.RepoName {
			t.Errorf("RepoName parse failed!\n given:  %s\n expect: %s\n result: %s\n", expected.url, expected.result.RepoName, test.RepoName)
			return
		}
	}
}
