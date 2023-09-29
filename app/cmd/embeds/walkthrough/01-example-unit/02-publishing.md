# Publishing

In order to use curriculum content, it must be published in Learn under a 'block'. Each block represents one remote repository, and each block can have several released versions of curriculum.

In this way curriculum development of blocks mirrors working on production software. Changes can be made in the repository until a new 'deployment' of curriculum is ready to be created from a single commit. Publishing 'deploys' a new release of curriculum content for that block at that commit SHA.

## Requirements

In order to publish curriculum, the project must be in a remote git based VCS, and Learn must have permission to access the repository.

Learn works with GitHub, Gitlab SaaS, and self hosted GitLab instances. Private repositories can be used; when doing so for GitHub, the `github-forge-production` user must have read access. For GitLab, the `galvanize-learn-production` user must have read access.

## Push to the Remote

In the VCS of your choice, create a repository/project ([GitHub](https://docs.github.com/en/migrations/importing-source-code/using-the-command-line-to-import-source-code/adding-locally-hosted-code-to-github), [GitLab](https://docs.gitlab.com/ee/user/project/)), then commit the all contents of the walkthrough:
```
git add -A
git commit -m "testing learn publish"
git push origin main
```

## Publish the curriculum

From with the project simply run
```
learn publish
```
from the command line.

The `learn` CLI tool will ensure that an `autoconfig.yaml` file exists (unless a `config.yaml` is found) in a commit on the remote, then it will publish the curriculum contents.

The block will now be discoverable by users with proper access from the [searchable blocks index page](https://learn-2.galvanize.com/blocks).
