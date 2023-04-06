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

When opening a PR, please follow [Conventional Commit](https://www.conventionalcommits.org/en/v1.0.0-beta.4/)
style for commit messages, and include the subdirectory as the scope, exluding "internal" if present. For
example, a commit that changes code in the [internal/gengapic](internal/gengapic) directory would be
`fix(gengapic): change foo to bar`.

## Bazel BUILD files

All of the normal Go tooling is sufficient to develop this project, the Makefile utilizes them.
However, this project supports Bazel build environments as well. As such, the Bazel build files need
to be kept up-to-date.

When new files are added, or existing ones are deleted or moved, the appropriate BUILD files must
be updated. This can, and should, be done automatically via Gazelle e.g. `make gazelle`.

_Note: If running on linux, you must comment remove the `''` value of the `-i` flag in this
target's `sed` command._

If there are changes to the `go.mod`, an equivalent change must be made on the `repositories.bzl`
macro that defines the Go repository dependencies. This can be done manually or by executing the
Make target: `make update-bazel-repos`.

_Note: If running on linux, you must comment remove the `''` value of the `-i` flag in this
target's `sed` command._

## Releases

Releases are managed by [Release Please](https://github.com/googleapis/release-please). Any commit
merged starting with `fix` or `feat` will trigger a Release Please pull request. Merge that to create
a release. Generator binaries will be added by CI after the release is created.

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
