# Changelog

## [0.29.0](https://github.com/googleapis/gapic-generator-go/compare/v0.28.0...v0.29.0) (2022-05-02)


### Features

* **gengapic:** regapic lro support ([#942](https://github.com/googleapis/gapic-generator-go/issues/942)) ([adb77e9](https://github.com/googleapis/gapic-generator-go/commit/adb77e9477e21709d7a0b2fc08350ec9d5129c11))

## [0.28.0](https://github.com/googleapis/gapic-generator-go/compare/v0.27.1...v0.28.0) (2022-04-13)


### Features

* **gengapic:** regapic implicit/explicit header injection ([#953](https://github.com/googleapis/gapic-generator-go/issues/953)) ([0dec6eb](https://github.com/googleapis/gapic-generator-go/commit/0dec6ebd79d779e223869eaf3a30674e77523fe2))


### Bug Fixes

* **gengapic:** regapic handle url.Parse error ([#951](https://github.com/googleapis/gapic-generator-go/issues/951)) ([098262f](https://github.com/googleapis/gapic-generator-go/commit/098262f98dd00eb109a8439e668fcda6084eee13))
* **gengapic:** REGAPIC support gax.CallSettings.Path override ([#957](https://github.com/googleapis/gapic-generator-go/issues/957)) ([ae6a2e1](https://github.com/googleapis/gapic-generator-go/commit/ae6a2e13f849d4b59d93a1b9d8d188d7ea008407))

### [0.27.1](https://github.com/googleapis/gapic-generator-go/compare/v0.27.0...v0.27.1) (2022-03-30)


### Bug Fixes

* **bazel:** add repo_mapping to all go_repository deps ([#947](https://github.com/googleapis/gapic-generator-go/issues/947)) ([7484ab5](https://github.com/googleapis/gapic-generator-go/commit/7484ab54f607900891130571b708b7e4aab77fd0))

## [0.27.0](https://github.com/googleapis/gapic-generator-go/compare/v0.26.0...v0.27.0) (2022-03-28)


### Features

* **gengapic:** regapic server-streaming ([#933](https://github.com/googleapis/gapic-generator-go/issues/933)) ([a6c0c81](https://github.com/googleapis/gapic-generator-go/commit/a6c0c818b67975c14ebab052c7bf37d45548baff))

## [0.26.0](https://github.com/googleapis/gapic-generator-go/compare/v0.25.1...v0.26.0) (2022-02-23)


### Features

* add dynamic routing header generation ([#887](https://github.com/googleapis/gapic-generator-go/issues/887)) ([e2520c7](https://github.com/googleapis/gapic-generator-go/commit/e2520c7f61228de7939fe8da673eda2b427539f4))
* make versionClient a var ([#912](https://github.com/googleapis/gapic-generator-go/issues/912)) ([7fe9fa5](https://github.com/googleapis/gapic-generator-go/commit/7fe9fa530f81a50b03971f7b3e44782d607106db))

### [0.25.1](https://github.com/googleapis/gapic-generator-go/compare/v0.25.0...v0.25.1) (2022-02-16)


### Bug Fixes

* **gengapic:** switch regapic use of xerrors to fmt for wrapping ([#902](https://github.com/googleapis/gapic-generator-go/issues/902)) ([3b5db4d](https://github.com/googleapis/gapic-generator-go/commit/3b5db4d28709b60be9eed520f007919a7ff21852))

## [0.25.0](https://github.com/googleapis/gapic-generator-go/compare/v0.24.0...v0.25.0) (2022-02-01)


### Features

* **gengapic:** diregapic add operation Wait helper ([#873](https://github.com/googleapis/gapic-generator-go/issues/873)) ([338e6e9](https://github.com/googleapis/gapic-generator-go/commit/338e6e922ec9ba1166e967bf1e091907c8a15b1e))
* **gengapic:** diregapic lro examples use wrapper ([#880](https://github.com/googleapis/gapic-generator-go/issues/880)) ([ccaa033](https://github.com/googleapis/gapic-generator-go/commit/ccaa033656680609d5a6c9527d1c8b6e7caa5594))
* **gengapic:** diregapic lro polling request params ([#876](https://github.com/googleapis/gapic-generator-go/issues/876)) ([fd0a07b](https://github.com/googleapis/gapic-generator-go/commit/fd0a07b92af5df8713b72c2b3fcdbbe04a91a6f5))
* **gengapic:** implement diregapic lro foundation + polling ([#866](https://github.com/googleapis/gapic-generator-go/issues/866)) ([99f2273](https://github.com/googleapis/gapic-generator-go/commit/99f2273c0bdc9ee994e4683bd46074b4d07a416f))


### Bug Fixes

* **bazel:** update rules_go, bazel deps, use bazelisk in ci ([#872](https://github.com/googleapis/gapic-generator-go/issues/872)) ([80a2ab1](https://github.com/googleapis/gapic-generator-go/commit/80a2ab18c32f04f6b4f6f22c58a1722f316522c6))
