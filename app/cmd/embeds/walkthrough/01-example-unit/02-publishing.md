---
Type: Lesson
UID: publishing
---

# Publishing

In order to use curriculum content, it must be published in Learn under a 'block'. Each block represents one remote repository, and each block can have several released versions of curriculum.

In this way curriculum development of blocks mirrors working on production software. Changes can be made in the repository until a new 'deployment' of curriculum is ready to be created from a single commit.

Publishing 'deploys' a new release of curriculum content for that block at that commit SHA. Different branches of the same repository can be published to provide greater flexibility when delivering content.

## Requirements

In order to publish curriculum, the project must be in a remote git based VCS, and Learn must have permission to access the repository.

Learn works with GitHub, Gitlab SaaS, and self hosted GitLab instances. Private repositories can be used; when doing so for GitHub, the `github-forge-production` user must have read access to the repository. For GitLab, the `galvanize-learn-production` user must have read access to the project.

## Push to the Remote

In the VCS of your choice, create a repository/project ([GitHub](https://docs.github.com/en/migrations/importing-source-code/using-the-command-line-to-import-source-code/adding-locally-hosted-code-to-github), [GitLab](https://docs.gitlab.com/ee/user/project/)), then commit all contents of the walkthrough:
```
git add -A
git commit -m "testing learn publish"
git push origin main
```

## Publish the curriculum

From within the project simply run
```
learn publish
```
from the command line.

The `learn` CLI tool will ensure that an `autoconfig.yaml` file exists (unless a `config.yaml` is found) in a commit on the remote, then it will publish the curriculum contents.

Any errors encountered while attempting to parse the directory will be provided in the event that the publish fails.

The block will now be discoverable by users with proper access from the [searchable blocks index page](https://learn-2.galvanize.com/blocks).

## Whats going on with all these UIDs?

A lot of first time curriculum developers wonder why we have to define so many identifiers within the content, "Shouldn't the database handle this?" they ask. Well, because git repositories back the content rendered in Learn, there needs to be a way for Learn to keep track of the same content _across releases_.

Suppose

Content files, units, and as you'll see challenges all require their own identifiers.
