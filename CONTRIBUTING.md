# How to Contribute

We'd love to accept your patches and contributions to this project. There are
just a few small guidelines you need to follow.

## General guidance

Development should be done on a fork. CI is configured to run tests against fork PRs, but you should
run the tests locally prior to opening a PR.

To only test `internal/gengapic`:

    make test-gapic

When changing the generated output, update the `want` files for tests:

    make golden

Test everything and run Showcase integration tests:

    make test

When opening a PR, please prefix the title with the main component the change effects,
separated from the message with a colon. This chould be one of `gapic`, `gencli`, `samples`,
`bazel`, or `chore`. For example:

    gapic: strip reference links in comments

    gencli: ignore output_only enums
    
    bazel: copy go_gapic_assembly_pkg macro impl

This allows the release automation to categorize the commits in the release notes.
If omitted, the PR title will be changed by a maintainer prior to submission.

## Bazel BUILD files

All of the normal Go tooling is sufficient to develop this project, the Makefile utilizes them.
However, this project supports Bazel build environments as well. As such, the Bazel build files need
to be kept up-to-date.

When new files are added, or existing ones are deleted or moved, the appropriate BUILD files must
be updated. This can, and should, be done automatically via Gazelle e.g. `bazel run //:gazelle`.

If there are changes to the `go.mod`, an equivalent change must be made on the `repositories.bzl`
macro that defines the Go repository dependencies. This can be done manually or by executing the
Make target: `make update-bazel-repos`.

## Releases

Releases are made on GitHub. CircleCI is configured to build and push tagged images upon release.
Tags always begin with a `v` and follow semver.

    git tag v1.2.3 && git push upstream --tags

A GitHub release will be made automatically with release notes generated based on the commit
messages of commits since the previous tag.

## Contributor License Agreement

Contributions to this project must be accompanied by a Contributor License
Agreement. You (or your employer) retain the copyright to your contribution;
this simply gives us permission to use and redistribute your contributions as
part of the project. Head over to <https://cla.developers.google.com/> to see
your current agreements on file or to sign a new one.

You generally only need to submit a CLA once, so if you've already submitted one
(even if it was for a different project), you probably don't need to do it
again.

## Code reviews

All submissions, including submissions by project members, require review. We
use GitHub pull requests for this purpose. Consult
[GitHub Help](https://help.github.com/articles/about-pull-requests/) for more
information on using pull requests.

## Community Guidelines

This project follows [Google's Open Source Community
Guidelines](https://opensource.google.com/conduct/).
