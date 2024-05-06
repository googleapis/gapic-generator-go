# Changelog

## [0.42.0](https://github.com/googleapis/gapic-generator-go/compare/v0.41.3...v0.42.0) (2024-05-06)


### Features

* **gengapic:** Add x-goog-api-version header ([#1498](https://github.com/googleapis/gapic-generator-go/issues/1498)) ([9e22b2c](https://github.com/googleapis/gapic-generator-go/commit/9e22b2c073129eb8ff0746db44d3c5abe16a0751))


### Bug Fixes

* **gencli:** Don't generate grpc.WithInsecure ([#1503](https://github.com/googleapis/gapic-generator-go/issues/1503)) ([27f02ea](https://github.com/googleapis/gapic-generator-go/commit/27f02ea30caa2199c1e2c5570781d2f25aaab40a))

## [0.41.3](https://github.com/googleapis/gapic-generator-go/compare/v0.41.2...v0.41.3) (2024-04-11)


### Bug Fixes

* **bazel:** Upgrade rules_go & gazelle to drop go_googleapis shading ([#1486](https://github.com/googleapis/gapic-generator-go/issues/1486)) ([7fdb7aa](https://github.com/googleapis/gapic-generator-go/commit/7fdb7aa68038bac0f65b8d6bb90b098b88fbb81a))
* Refactor codebase to protobuf-go v2 ([#1489](https://github.com/googleapis/gapic-generator-go/issues/1489)) ([e84b5ef](https://github.com/googleapis/gapic-generator-go/commit/e84b5ef2669849eeb30d8925bee3b15524bcc521))


### Performance Improvements

* Use errors.New when possible ([#1490](https://github.com/googleapis/gapic-generator-go/issues/1490)) ([c1fd44c](https://github.com/googleapis/gapic-generator-go/commit/c1fd44ceb367d9e4cf43864b34b1cc39c4138a61))

## [0.41.2](https://github.com/googleapis/gapic-generator-go/compare/v0.41.1...v0.41.2) (2024-04-06)


### Bug Fixes

* **bazel:** Add missing otel deps & resolve patches ([#1482](https://github.com/googleapis/gapic-generator-go/issues/1482)) ([1f9fa36](https://github.com/googleapis/gapic-generator-go/commit/1f9fa3611328e57d6a5a9b2c8db4e4017e0c92c5))

## [0.41.1](https://github.com/googleapis/gapic-generator-go/compare/v0.41.0...v0.41.1) (2024-03-25)


### Bug Fixes

* **internal/gengapic:** Add iter response access example ([#1468](https://github.com/googleapis/gapic-generator-go/issues/1468)) ([02e6c65](https://github.com/googleapis/gapic-generator-go/commit/02e6c65e2df34e5bcb1c85460a2a422db971d843))

## [0.41.0](https://github.com/googleapis/gapic-generator-go/compare/v0.40.0...v0.41.0) (2024-02-27)


### Features

* **gengapic:** Add support for AutoPopulatedFields UUID4 ([#1460](https://github.com/googleapis/gapic-generator-go/issues/1460)) ([2f3b7b9](https://github.com/googleapis/gapic-generator-go/commit/2f3b7b99c41b1a4083af215fb75658fcf3c4d30e))

## [0.40.0](https://github.com/googleapis/gapic-generator-go/compare/v0.39.4...v0.40.0) (2024-01-23)


### Features

* **gengapic:** Add universe domain support ([#1452](https://github.com/googleapis/gapic-generator-go/issues/1452)) ([c72b650](https://github.com/googleapis/gapic-generator-go/commit/c72b650af246b5db4d30a6730c943cde103936a7))

## [0.39.4](https://github.com/googleapis/gapic-generator-go/compare/v0.39.3...v0.39.4) (2023-11-08)


### Bug Fixes

* **gencli:** Title case dot delimited segments ([#1442](https://github.com/googleapis/gapic-generator-go/issues/1442)) ([a1d372c](https://github.com/googleapis/gapic-generator-go/commit/a1d372c88fc5a26cd61df0efa57611cefb23ac22))

## [0.39.3](https://github.com/googleapis/gapic-generator-go/compare/v0.39.2...v0.39.3) (2023-11-07)


### Bug Fixes

* **gengapic:** Rename file aux to auxiliary ([#1440](https://github.com/googleapis/gapic-generator-go/issues/1440)) ([9c0fbd5](https://github.com/googleapis/gapic-generator-go/commit/9c0fbd58ac0d2d278973877ee690230cf135fc38))

## [0.39.2](https://github.com/googleapis/gapic-generator-go/compare/v0.39.1...v0.39.2) (2023-11-02)


### Bug Fixes

* **gengapic:** Fix generator reset ordering ([#1436](https://github.com/googleapis/gapic-generator-go/issues/1436)) ([8d64d03](https://github.com/googleapis/gapic-generator-go/commit/8d64d03a37f22ce17cdbe4906c74cc874e2ee56b))

## [0.39.1](https://github.com/googleapis/gapic-generator-go/compare/v0.39.0...v0.39.1) (2023-11-01)


### Bug Fixes

* **gengapic:** Ensure input type imported ([#1434](https://github.com/googleapis/gapic-generator-go/issues/1434)) ([f6be536](https://github.com/googleapis/gapic-generator-go/commit/f6be536fff68c698a52f370d99b73f6dd133b272))

## [0.39.0](https://github.com/googleapis/gapic-generator-go/compare/v0.38.2...v0.39.0) (2023-11-01)


### Features

* **internal/gengapic:** Move operations & iterators to aux.go ([#1428](https://github.com/googleapis/gapic-generator-go/issues/1428)) ([e8ad272](https://github.com/googleapis/gapic-generator-go/commit/e8ad27239eca2cb09623812e1bb0bf88f5f7a5c1))

## [0.38.2](https://github.com/googleapis/gapic-generator-go/compare/v0.38.1...v0.38.2) (2023-09-22)


### Bug Fixes

* **internal/gengapic:** Add workaround to a delete lro ([#1398](https://github.com/googleapis/gapic-generator-go/issues/1398)) ([096e74d](https://github.com/googleapis/gapic-generator-go/commit/096e74d5f822716e1c6eaa22b6326e003e1cdf28))

## [0.38.1](https://github.com/googleapis/gapic-generator-go/compare/v0.38.0...v0.38.1) (2023-09-21)


### Bug Fixes

* **gengapic:** Support deprecated as a release-level option value ([#1390](https://github.com/googleapis/gapic-generator-go/issues/1390)) ([0b0f823](https://github.com/googleapis/gapic-generator-go/commit/0b0f823ce3e2be05570e572eff5c2f98252cc8a5))
* **internal/gengapic:** Add workaround for operation collision ([#1397](https://github.com/googleapis/gapic-generator-go/issues/1397)) ([edb3b8f](https://github.com/googleapis/gapic-generator-go/commit/edb3b8fb66bd6e8e57ac345f51030175807afc1d))

## [0.38.0](https://github.com/googleapis/gapic-generator-go/compare/v0.37.2...v0.38.0) (2023-08-07)


### Features

* **gengapic:** Use gax.BuildHeaders and gax.InsertMetadataIntoOutgoingContext ([#1368](https://github.com/googleapis/gapic-generator-go/issues/1368)) ([6f782f9](https://github.com/googleapis/gapic-generator-go/commit/6f782f96a29a27b6c7ca5d21a433533270679bcc)), closes [#1300](https://github.com/googleapis/gapic-generator-go/issues/1300) [#1301](https://github.com/googleapis/gapic-generator-go/issues/1301)


### Bug Fixes

* Update client docs to ref base docs more ([#1375](https://github.com/googleapis/gapic-generator-go/issues/1375)) ([b78472c](https://github.com/googleapis/gapic-generator-go/commit/b78472cc517c0f37582ee69fec21fdf992aca92b))

## [0.37.2](https://github.com/googleapis/gapic-generator-go/compare/v0.37.1...v0.37.2) (2023-06-20)


### Bug Fixes

* **gengapic:** Remove encoded quotes in query params ([#1364](https://github.com/googleapis/gapic-generator-go/issues/1364)) ([5d62c34](https://github.com/googleapis/gapic-generator-go/commit/5d62c344a7a7aac5b49979a800430b144062900b)), closes [#1363](https://github.com/googleapis/gapic-generator-go/issues/1363)

## [0.37.1](https://github.com/googleapis/gapic-generator-go/compare/v0.37.0...v0.37.1) (2023-06-13)


### Bug Fixes

* **gengapic:** Remove unknown enum error helper ([#1358](https://github.com/googleapis/gapic-generator-go/issues/1358)) ([34af96c](https://github.com/googleapis/gapic-generator-go/commit/34af96cccbe3fe5f700133fbe5b6f4595a6996fc))
* **gengapic:** Use gax GoVersion ([#1359](https://github.com/googleapis/gapic-generator-go/issues/1359)) ([9116eca](https://github.com/googleapis/gapic-generator-go/commit/9116eca768029a58f5d353795748b6351c0fc9eb))
* Refactor usage of deprecated io/ioutil to io ([#1336](https://github.com/googleapis/gapic-generator-go/issues/1336)) ([455a421](https://github.com/googleapis/gapic-generator-go/commit/455a421978e82d79155fbe120af2725f5cc8b9da))

## [0.37.0](https://github.com/googleapis/gapic-generator-go/compare/v0.36.0...v0.37.0) (2023-05-17)


### Features

* **gengapic:** Use WithTimeout for default logical timeout ([#1267](https://github.com/googleapis/gapic-generator-go/issues/1267)) ([7d1418f](https://github.com/googleapis/gapic-generator-go/commit/7d1418faaa4d7bc0ef76b661ab1be33c85be8a7f)), closes [#1206](https://github.com/googleapis/gapic-generator-go/issues/1206)

## [0.36.0](https://github.com/googleapis/gapic-generator-go/compare/v0.35.7...v0.36.0) (2023-05-09)


### Features

* **gengapic:** Raise Operation errors from diregapic ([#1323](https://github.com/googleapis/gapic-generator-go/issues/1323)) ([66d43c6](https://github.com/googleapis/gapic-generator-go/commit/66d43c661329f26440aae237e5f5bf1716489e68)), closes [#1320](https://github.com/googleapis/gapic-generator-go/issues/1320)

## [0.35.7](https://github.com/googleapis/gapic-generator-go/compare/v0.35.6...v0.35.7) (2023-04-21)


### Bug Fixes

* **internal/gengapic:** Write snippet output to cloud.google.com/go ([#1313](https://github.com/googleapis/gapic-generator-go/issues/1313)) ([dfc5ce2](https://github.com/googleapis/gapic-generator-go/commit/dfc5ce2336a0f1a8b10732e32164540ba5686883))

## [0.35.6](https://github.com/googleapis/gapic-generator-go/compare/v0.35.5...v0.35.6) (2023-04-18)


### Bug Fixes

* **deps:** Revert middleware version, drop s2a-go ([#1308](https://github.com/googleapis/gapic-generator-go/issues/1308)) ([4609f47](https://github.com/googleapis/gapic-generator-go/commit/4609f47aae6f175fde04db51b0e564056ac135e4))
* **internal/gengapic:** Fix mixin file filter ([#1310](https://github.com/googleapis/gapic-generator-go/issues/1310)) ([79a7a34](https://github.com/googleapis/gapic-generator-go/commit/79a7a3443dca516746f9998f12749b9cefdac61c))

## [0.35.5](https://github.com/googleapis/gapic-generator-go/compare/v0.35.4...v0.35.5) (2023-04-07)


### Bug Fixes

* Add time import for compute operations ([#1292](https://github.com/googleapis/gapic-generator-go/issues/1292)) ([a18ff1e](https://github.com/googleapis/gapic-generator-go/commit/a18ff1eeccec94f94685f1641fe3b130a0d4d834))

## [0.35.4](https://github.com/googleapis/gapic-generator-go/compare/v0.35.3...v0.35.4) (2023-04-06)


### Bug Fixes

* Add wkt desc, bazel snippet copy ([#1288](https://github.com/googleapis/gapic-generator-go/issues/1288)) ([d7965aa](https://github.com/googleapis/gapic-generator-go/commit/d7965aa74fdcffa7b13d77e6897c58ed2259e7a8))
* **gengapic:** Check error from exampleMethodBody ([#1291](https://github.com/googleapis/gapic-generator-go/issues/1291)) ([068217c](https://github.com/googleapis/gapic-generator-go/commit/068217c4ddf88ce7abebbf1c6f86e7464ee87598))

## [0.35.3](https://github.com/googleapis/gapic-generator-go/compare/v0.35.2...v0.35.3) (2023-04-05)


### Bug Fixes

* Restore omit-snippets for bazel ([#1278](https://github.com/googleapis/gapic-generator-go/issues/1278)) ([9810cf4](https://github.com/googleapis/gapic-generator-go/commit/9810cf44cc6348cc6b1034310ebf512e7ff7bc40))

## [0.35.2](https://github.com/googleapis/gapic-generator-go/compare/v0.35.1...v0.35.2) (2023-03-10)


### Bug Fixes

* Explicitly override the import paths for iam and longrunning ([#1257](https://github.com/googleapis/gapic-generator-go/issues/1257)) ([132fb43](https://github.com/googleapis/gapic-generator-go/commit/132fb43b6034fa77790964c4a8a0a5ddb214a451))

## [0.35.1](https://github.com/googleapis/gapic-generator-go/compare/v0.35.0...v0.35.1) (2023-03-09)


### Bug Fixes

* Omit-snippets for bazel ([#1254](https://github.com/googleapis/gapic-generator-go/issues/1254)) ([76efded](https://github.com/googleapis/gapic-generator-go/commit/76efded07898baeb9d2773c4bb7cb6694f6c1adc))

## [0.35.0](https://github.com/googleapis/gapic-generator-go/compare/v0.34.0...v0.35.0) (2023-03-08)


### Features

* **gengapic:** Add snippets ([#1220](https://github.com/googleapis/gapic-generator-go/issues/1220)) ([98f7a13](https://github.com/googleapis/gapic-generator-go/commit/98f7a13fb501df907462041eb39ae93cac472ebd))
* Update generator for new iam and longrunning locations ([#1247](https://github.com/googleapis/gapic-generator-go/issues/1247)) ([2584c29](https://github.com/googleapis/gapic-generator-go/commit/2584c29a5a545d25622962fae983e2b16640df9d))

## [0.34.0](https://github.com/googleapis/gapic-generator-go/compare/v0.33.7...v0.34.0) (2023-03-01)


### Features

* Migrate to lro & iam submodules ([#1240](https://github.com/googleapis/gapic-generator-go/issues/1240)) ([76944b9](https://github.com/googleapis/gapic-generator-go/commit/76944b9f52ff6bacf5f68b267fbbf210d091c823))

## [0.33.7](https://github.com/googleapis/gapic-generator-go/compare/v0.33.6...v0.33.7) (2023-01-19)


### Bug Fixes

* **gengapic:** Fix panic when a file not conain any service ([#1214](https://github.com/googleapis/gapic-generator-go/issues/1214)) ([cd0c02f](https://github.com/googleapis/gapic-generator-go/commit/cd0c02f43a041b62c2cc9dfafe5aff8ccabe3485)), closes [#1213](https://github.com/googleapis/gapic-generator-go/issues/1213)

## [0.33.6](https://github.com/googleapis/gapic-generator-go/compare/v0.33.5...v0.33.6) (2023-01-09)


### Bug Fixes

* **gengapic:** Inject gRPC server stream call opts ([#1202](https://github.com/googleapis/gapic-generator-go/issues/1202)) ([1b93213](https://github.com/googleapis/gapic-generator-go/commit/1b93213c0eb80c6de85ff6de6009b31c473ad3ab))
* **gengapic:** Move top-level package doc links to top ([#1175](https://github.com/googleapis/gapic-generator-go/issues/1175)) ([8cf6194](https://github.com/googleapis/gapic-generator-go/commit/8cf619464bcb316adc7c1d93bab51425981d9f55)), closes [#1140](https://github.com/googleapis/gapic-generator-go/issues/1140)

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
