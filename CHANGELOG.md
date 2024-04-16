# Changelog

All notable changes to this project will be documented in this file. See [convention-change-log](https://github.com/convention-change/convention-change-log) for commit guidelines.

## [1.2.0](https://github.com/woodpecker-kit/woodpecker-npm/compare/1.1.0...v1.2.0) (2024-04-17)

### ✨ Features

* change to new docker build pipeline ([443f5e6b](https://github.com/woodpecker-kit/woodpecker-npm/commit/443f5e6b280e5b40da92302484f7094130f9389c))

### 👷‍ Build System

* update `docker-bake.hcl` at `image-all` ([0fd98143](https://github.com/woodpecker-kit/woodpecker-npm/commit/0fd98143bf4ec9d6401f3254ec273c1e09168f05))

* add docker-bake.hcl ([08ed02d8](https://github.com/woodpecker-kit/woodpecker-npm/commit/08ed02d8978e3f13ec7ebfd162e6aa9497a6fc7b))

## [1.1.0](https://github.com/woodpecker-kit/woodpecker-npm/compare/1.0.1...v1.1.0) (2024-04-07)

### 🐛 Bug Fixes

* skip version check at tag not empty and open `npm-force-tag` ([35286876](https://github.com/woodpecker-kit/woodpecker-npm/commit/352868761426f53e4ba1a0f5400c9e8bb2cb0d6f))

### ✨ Features

* flag `npm-force-tag` check package version by semver and force publish ([c854472a](https://github.com/woodpecker-kit/woodpecker-npm/commit/c854472a902c89c336558e52a6b5d10552ec4815)), fe [#3](https://github.com/woodpecker-kit/woodpecker-npm/issues/3)

* flag `npm-dry-run` to open dry run mode, will not publish to NPM registry ([174a43e4](https://github.com/woodpecker-kit/woodpecker-npm/commit/174a43e4d9dd0e9c04a42477bab2774fe5fada11))

* change .npmrc path, default .npmrc file will write in `npm-folder` ([128057bc](https://github.com/woodpecker-kit/woodpecker-npm/commit/128057bcda9759e14d472d5857c17bd0446e0258))

### 📝 Documentation

* add usage of `npm-tag` ([dfca9ba2](https://github.com/woodpecker-kit/woodpecker-npm/commit/dfca9ba2fe4eb7e87b5b82ea2ce1188250043996)), fe [#1](https://github.com/woodpecker-kit/woodpecker-npm/issues/1)

* update usage of doc/docs.md ([101fafe6](https://github.com/woodpecker-kit/woodpecker-npm/commit/101fafe6b7340ddd3c5e50fc8ed636d22aa71bcf))

## [1.0.1](https://github.com/woodpecker-kit/woodpecker-npm/compare/1.0.0...v1.0.1) (2024-04-07)

### 🐛 Bug Fixes

* fix whoami check by custom `npm-registry` ([71a87dce](https://github.com/woodpecker-kit/woodpecker-npm/commit/71a87dcedf043a7829f776d6a69fdf332f88c3f1))

## 1.0.0 (2024-04-06)

### ✨ Features

* let whoami check by npm-registry and update woodpecker-tools v1.19.0 ([2a9d8512](https://github.com/woodpecker-kit/woodpecker-npm/commit/2a9d851285638d8a76188adbccd3a6ab67dafa67))

* add kubernetes runner patch by /run/drone/env or file by env `kubernetes runner patch` ([40d5bbe6](https://github.com/woodpecker-kit/woodpecker-npm/commit/40d5bbe6649a9c7533fae574a0c39acdc7b2fffd))

* add basic of npm publish ([97105a7b](https://github.com/woodpecker-kit/woodpecker-npm/commit/97105a7bdda6eaf19eb01f430a30bbd7fa6a6d13))

### 📝 Documentation

* doc for docker publish ([7de57dba](https://github.com/woodpecker-kit/woodpecker-npm/commit/7de57dba3fb6edfc5088e13bd1acfa69b28da43f))
