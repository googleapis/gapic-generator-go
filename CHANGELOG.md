# Changelog

## [0.33.5](https://github.com/googleapis/gapic-generator-go/compare/v0.33.4...v0.33.5) (2022-12-07)


### Bug Fixes

* **gengapic:** Extraneous protojson import ([#1196](https://github.com/googleapis/gapic-generator-go/issues/1196)) ([023fed2](https://github.com/googleapis/gapic-generator-go/commit/023fed2baa2898fec0e0658963ae549eb2077971))
* **gengapic:** Only gen REST client for RESTable services ([#1199](https://github.com/googleapis/gapic-generator-go/issues/1199)) ([0f063f1](https://github.com/googleapis/gapic-generator-go/commit/0f063f13c3453bf0832ae44d847e69915781ee2f))

## [0.33.4](https://github.com/googleapis/gapic-generator-go/compare/v0.33.3...v0.33.4) (2022-11-08)


### Bug Fixes

* **gengapic:** Document client/bidi streaming unsupported in REST ([#1181](https://github.com/googleapis/gapic-generator-go/issues/1181)) ([f9a9191](https://github.com/googleapis/gapic-generator-go/commit/f9a9191e1116df5127fa31ddd3655594c0071602))

## [0.33.3](https://github.com/googleapis/gapic-generator-go/compare/v0.33.2...v0.33.3) (2022-10-27)


### Bug Fixes

* **bazel:** Manually add compute/metadata dependency ([#1172](https://github.com/googleapis/gapic-generator-go/issues/1172)) ([516cb00](https://github.com/googleapis/gapic-generator-go/commit/516cb00707c396f6ae168814ecf3ff21de38c8ea))

## [0.33.2](https://github.com/googleapis/gapic-generator-go/compare/v0.33.1...v0.33.2) (2022-10-21)


### Bug Fixes

* **gengapic:** Separate repeated prim qp values ([#1161](https://github.com/googleapis/gapic-generator-go/issues/1161)) ([f2edb34](https://github.com/googleapis/gapic-generator-go/commit/f2edb34ab8d8b581e1f9d3b73681973497893982))

## [0.33.1](https://github.com/googleapis/gapic-generator-go/compare/v0.33.0...v0.33.1) (2022-09-13)


### Bug Fixes

* Handle wkt in paged requests ([#1132](https://github.com/googleapis/gapic-generator-go/issues/1132)) ([0c23c8d](https://github.com/googleapis/gapic-generator-go/commit/0c23c8d52b0d8eadc56966fd8d70a7c090375b10))

## [0.33.0](https://github.com/googleapis/gapic-generator-go/compare/v0.32.1...v0.33.0) (2022-09-06)


### Features

* **gengapic:** Add in-line snippet comment to example.go ([#1120](https://github.com/googleapis/gapic-generator-go/issues/1120)) ([88372c5](https://github.com/googleapis/gapic-generator-go/commit/88372c545bf8ac3a51fc26f6986e8c850d89da06))


### Bug Fixes

* Regapic support for proto wkt in query params ([#1124](https://github.com/googleapis/gapic-generator-go/issues/1124)) ([f000c98](https://github.com/googleapis/gapic-generator-go/commit/f000c98586db39e3ab7597b0b1a995e81ea14f6b))

## [0.32.1](https://github.com/googleapis/gapic-generator-go/compare/v0.32.0...v0.32.1) (2022-08-23)


### Bug Fixes

* update deprecation warning on Connection methods ([#1111](https://github.com/googleapis/gapic-generator-go/issues/1111)) ([f9a2c53](https://github.com/googleapis/gapic-generator-go/commit/f9a2c53cb954d6d4f27bdcbf1b643d3ae10622af)), closes [#1110](https://github.com/googleapis/gapic-generator-go/issues/1110)

## [0.32.0](https://github.com/googleapis/gapic-generator-go/compare/v0.31.2...v0.32.0) (2022-08-16)


### Features

* **gengapic:** rest-numeric-enums option enables enum-encoding sys param ([#1022](https://github.com/googleapis/gapic-generator-go/issues/1022)) ([6bbbf6f](https://github.com/googleapis/gapic-generator-go/commit/6bbbf6f7a37bc29861df9864926570c5046c6874))


### Bug Fixes

* **gengapic:** fix linkParser regexp to support multi-line link tags ([#1098](https://github.com/googleapis/gapic-generator-go/issues/1098)) ([863675e](https://github.com/googleapis/gapic-generator-go/commit/863675e499c933b35c14217cb85786d6c91086a2)), closes [#1097](https://github.com/googleapis/gapic-generator-go/issues/1097)
* **gengapic:** fix service-specific constructor name in doc_file.go ([#1099](https://github.com/googleapis/gapic-generator-go/issues/1099)) ([4f80726](https://github.com/googleapis/gapic-generator-go/commit/4f80726aa9f0f0357d5ebcffc7a8670657d35a3d)), closes [#1077](https://github.com/googleapis/gapic-generator-go/issues/1077)

## [0.31.2](https://github.com/googleapis/gapic-generator-go/compare/v0.31.1...v0.31.2) (2022-07-18)


### Bug Fixes

* **gengapic:** regapic GetOperation path fallback logic ([#1072](https://github.com/googleapis/gapic-generator-go/issues/1072)) ([71ff189](https://github.com/googleapis/gapic-generator-go/commit/71ff189c2ec0bc4fd24944faa59a67fe5a388cc0))

## [0.31.1](https://github.com/googleapis/gapic-generator-go/compare/v0.31.0...v0.31.1) (2022-07-14)


### Bug Fixes

* **gengapic:** fix unused imports ([#1071](https://github.com/googleapis/gapic-generator-go/issues/1071)) ([cabfdf3](https://github.com/googleapis/gapic-generator-go/commit/cabfdf323669599512b88ac195ee9b126fbc34d3))
* **gengapic:** regapic fix missing return statement ([#1054](https://github.com/googleapis/gapic-generator-go/issues/1054)) ([7d08d1b](https://github.com/googleapis/gapic-generator-go/commit/7d08d1bddb9cd15fc6cd30f60fc6f352b1c00eeb))

## [0.31.0](https://github.com/googleapis/gapic-generator-go/compare/v0.30.0...v0.31.0) (2022-06-14)


### Features

* **genpapic:** support protobuf-go go_package mapping option ([#1029](https://github.com/googleapis/gapic-generator-go/issues/1029)) ([f40c830](https://github.com/googleapis/gapic-generator-go/commit/f40c8300f39be4584ff195f1cbc0488bebd563ae))


### Bug Fixes

* change go_gapic_library rule  transport argument type from array to string ([#1038](https://github.com/googleapis/gapic-generator-go/issues/1038)) ([a0ee493](https://github.com/googleapis/gapic-generator-go/commit/a0ee493d046cc32a61f742c38befdbd0bdf10547))

## [0.30.0](https://github.com/googleapis/gapic-generator-go/compare/v0.29.2...v0.30.0) (2022-06-01)


### Features

* **gengapic:** REGAPIC retry support ([#993](https://github.com/googleapis/gapic-generator-go/issues/993)) ([4021354](https://github.com/googleapis/gapic-generator-go/commit/4021354b355672c15c937040af9708a41067408d))


### Bug Fixes

* **bazel:** add missing regapic dep to go_library ([#1016](https://github.com/googleapis/gapic-generator-go/issues/1016)) ([ba270da](https://github.com/googleapis/gapic-generator-go/commit/ba270da6d0e8753b71e549d9a4c172a70aa8d7df))

### [0.29.2](https://github.com/googleapis/gapic-generator-go/compare/v0.29.1...v0.29.2) (2022-05-18)


### Bug Fixes

* **gengapic:** REGAPIC encode enums as numbers ([#994](https://github.com/googleapis/gapic-generator-go/issues/994)) ([1a04703](https://github.com/googleapis/gapic-generator-go/commit/1a04703d600be9a1a79f6998876b80e46d5b023b))

### [0.29.1](https://github.com/googleapis/gapic-generator-go/compare/v0.29.0...v0.29.1) (2022-05-09)


### Bug Fixes

* **gengapic:** REGAPIC fix path param parsing to remove duped query param ([#981](https://github.com/googleapis/gapic-generator-go/issues/981)) ([d326973](https://github.com/googleapis/gapic-generator-go/commit/d326973a847c21dce4cb14fca45c8ec4e501e785))

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
