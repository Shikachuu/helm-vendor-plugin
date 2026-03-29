# Changelog

## [1.2.2](https://github.com/Shikachuu/helm-vendor-plugin/compare/v1.2.1...v1.2.2) (2026-03-29)


### 🐛 Bug Fixes

* replace release workflow in gorelease and release-please ([93cd391](https://github.com/Shikachuu/helm-vendor-plugin/commit/93cd391dcf316843c98c1a125e679f1e26b3e459))


### 🔧 Miscellaneous

* **deps:** Bump docker/login-action from 3.7.0 to 4.0.0 ([#36](https://github.com/Shikachuu/helm-vendor-plugin/issues/36)) ([fb4d7dc](https://github.com/Shikachuu/helm-vendor-plugin/commit/fb4d7dc06f3d556809aaba75cc17ed35a91b1494))
* **deps:** Bump github.com/cloudflare/circl from 1.6.1 to 1.6.3 ([#32](https://github.com/Shikachuu/helm-vendor-plugin/issues/32)) ([8391fca](https://github.com/Shikachuu/helm-vendor-plugin/commit/8391fca81244012328e84356c609b11b4fc7cfff))
* **deps:** Bump github/codeql-action from 4.32.3 to 4.32.6 ([#35](https://github.com/Shikachuu/helm-vendor-plugin/issues/35)) ([7c4c337](https://github.com/Shikachuu/helm-vendor-plugin/commit/7c4c337f3602ff0a4c4ac9f90e71da80fa6a426a))
* **deps:** Bump github/codeql-action from 4.32.6 to 4.34.1 ([#42](https://github.com/Shikachuu/helm-vendor-plugin/issues/42)) ([9d182dd](https://github.com/Shikachuu/helm-vendor-plugin/commit/9d182dd23069386eeeda10804512ebaf75940259))
* update mise and it's lock file + go dependencies ([#44](https://github.com/Shikachuu/helm-vendor-plugin/issues/44)) ([caf08cc](https://github.com/Shikachuu/helm-vendor-plugin/commit/caf08cc6c2d48b2770c0f3f94a426add51a0c083))

## [1.2.1](https://github.com/Shikachuu/helm-vendor-plugin/compare/v1.2.0...v1.2.1) (2026-02-19)


### 🐛 Bug Fixes

* publish tags on merged release PRs ([#27](https://github.com/Shikachuu/helm-vendor-plugin/issues/27)) ([6506b2d](https://github.com/Shikachuu/helm-vendor-plugin/commit/6506b2d5c5a32a64ca5f80af5ab0ac86e64abb57))


### 📚 Documentation

* update readme with extract field and new dev tools ([#28](https://github.com/Shikachuu/helm-vendor-plugin/issues/28)) ([0405b54](https://github.com/Shikachuu/helm-vendor-plugin/commit/0405b54f67ddb62211220de67d91199e6932d0de))


### 🔧 Miscellaneous

* **deps:** Bump docker/login-action from 3.6.0 to 3.7.0 ([#23](https://github.com/Shikachuu/helm-vendor-plugin/issues/23)) ([857e2dc](https://github.com/Shikachuu/helm-vendor-plugin/commit/857e2dce736c6a97946c843258921a3a56036cb3))
* **deps:** Bump github/codeql-action from 4.31.9 to 4.32.1 ([#22](https://github.com/Shikachuu/helm-vendor-plugin/issues/22)) ([489ed21](https://github.com/Shikachuu/helm-vendor-plugin/commit/489ed21e544fd24ec35fd43b0b7a7e5c23c4ed0d))
* **deps:** Bump github/codeql-action from 4.32.1 to 4.32.3 ([#26](https://github.com/Shikachuu/helm-vendor-plugin/issues/26)) ([802ca21](https://github.com/Shikachuu/helm-vendor-plugin/commit/802ca216c5b422c9bc409c9724aeb513f3937e37))
* **deps:** Bump jdx/mise-action from 3.5.1 to 3.6.1 ([#19](https://github.com/Shikachuu/helm-vendor-plugin/issues/19)) ([339f19e](https://github.com/Shikachuu/helm-vendor-plugin/commit/339f19e53d7af42c976085ca2b89e8c6e36aa6b1))
* **deps:** Bump the go-dependency-updates group across 1 directory with 2 updates ([#20](https://github.com/Shikachuu/helm-vendor-plugin/issues/20)) ([ba6d4bb](https://github.com/Shikachuu/helm-vendor-plugin/commit/ba6d4bb44e5e66cd5557287602d12b4b787c3138))
* **deps:** Bump the go-dependency-updates group with 2 updates ([#25](https://github.com/Shikachuu/helm-vendor-plugin/issues/25)) ([e64dae4](https://github.com/Shikachuu/helm-vendor-plugin/commit/e64dae444db66daa51702fd8b5973debb2264046))

## [1.2.0](https://github.com/Shikachuu/helm-vendor-plugin/compare/v1.1.0...v1.2.0) (2026-01-10)


### ✨ Features

* implement codeql pipeline for security checks ([#10](https://github.com/Shikachuu/helm-vendor-plugin/issues/10)) ([a4ba2fe](https://github.com/Shikachuu/helm-vendor-plugin/commit/a4ba2fe24e2bcb73c54cd62ecc0e47386d26cfc1))


### 🐛 Bug Fixes

* replace old release please version ([#13](https://github.com/Shikachuu/helm-vendor-plugin/issues/13)) ([da51173](https://github.com/Shikachuu/helm-vendor-plugin/commit/da51173a740bba73d1ddb4ed317ceb48c4c50809))
* update release type and plugin.yaml update ([#15](https://github.com/Shikachuu/helm-vendor-plugin/issues/15)) ([ce45bbe](https://github.com/Shikachuu/helm-vendor-plugin/commit/ce45bbe805fad67475ab79020392a0de1d39ded1))


### 🔧 Miscellaneous

* **deps:** Bump the go-dependency-updates group with 5 updates ([#11](https://github.com/Shikachuu/helm-vendor-plugin/issues/11)) ([1b8cb3a](https://github.com/Shikachuu/helm-vendor-plugin/commit/1b8cb3a9dae3a2a707f61beeb104bda8ef31b007))


### ✅ Tests

* add initial unit and integration tests ([#16](https://github.com/Shikachuu/helm-vendor-plugin/issues/16)) ([34abec8](https://github.com/Shikachuu/helm-vendor-plugin/commit/34abec828c1e4bbbdf0c7b342cfafdaafa7c2656))
