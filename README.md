[![ci](https://github.com/woodpecker-kit/woodpecker-npm/workflows/ci/badge.svg)](https://github.com/woodpecker-kit/woodpecker-npm/actions/workflows/ci.yml)

[![go mod version](https://img.shields.io/github/go-mod/go-version/woodpecker-kit/woodpecker-npm?label=go.mod)](https://github.com/woodpecker-kit/woodpecker-npm)
[![GoDoc](https://godoc.org/github.com/woodpecker-kit/woodpecker-npm?status.png)](https://godoc.org/github.com/woodpecker-kit/woodpecker-npm)
[![goreportcard](https://goreportcard.com/badge/github.com/woodpecker-kit/woodpecker-npm)](https://goreportcard.com/report/github.com/woodpecker-kit/woodpecker-npm)

[![GitHub license](https://img.shields.io/github/license/woodpecker-kit/woodpecker-npm)](https://github.com/woodpecker-kit/woodpecker-npm)
[![codecov](https://codecov.io/gh/woodpecker-kit/woodpecker-npm/branch/main/graph/badge.svg)](https://codecov.io/gh/woodpecker-kit/woodpecker-npm)
[![GitHub latest SemVer tag)](https://img.shields.io/github/v/tag/woodpecker-kit/woodpecker-npm)](https://github.com/woodpecker-kit/woodpecker-npm/tags)
[![GitHub release)](https://img.shields.io/github/v/release/woodpecker-kit/woodpecker-npm)](https://github.com/woodpecker-kit/woodpecker-npm/releases)

## for what

- this project used to woodpecker plugin

## Contributing

[![Contributor Covenant](https://img.shields.io/badge/contributor%20covenant-v1.4-ff69b4.svg)](.github/CONTRIBUTING_DOC/CODE_OF_CONDUCT.md)
[![GitHub contributors](https://img.shields.io/github/contributors/woodpecker-kit/woodpecker-npm)](https://github.com/woodpecker-kit/woodpecker-npm/graphs/contributors)

We welcome community contributions to this project.

Please read [Contributor Guide](.github/CONTRIBUTING_DOC/CONTRIBUTING.md) for more information on how to get started.

请阅读有关 [贡献者指南](.github/CONTRIBUTING_DOC/zh-CN/CONTRIBUTING.md) 以获取更多如何入门的信息

## Features

- [x] publish npm package by npm cli, so must install npm cli or under nodejs env
  - [x] by default docker image `node:20.11.1-alpine` for env of nodejs
  - [ ] if you use `local` backend, must install `npm` and `node` at local
- [x] support `npm-token` or `npm-username` and `npm-password` to publish
- [x] support `npm-tag` to publish, as `latest`
- [x] support `npm-access` to publish scoped package
- [x] support `npm-folder` to publish, which must containing `package.json`
- [x] can skip `npm whoami` check by open `npm-skip-whoami`
- [x] can skip `npm ssl` verify by open `npm-skip-verify-ssl`
- [x] can fail on version conflict by open `npm-fail-on-version-conflict`
- [ ] more perfect test case coverage
- [ ] more perfect benchmark case

## usage

- use this template, replace list below and add usage
    - `github.com/woodpecker-kit/woodpecker-npm` to your package name
    - `woodpecker-kit` to your owner name
    - `woodpecker-npm` to your project name

- use github action for this workflow push to docker hub, must add at github secrets
    - `DOCKERHUB_OWNER` user of docker hub
    - `DOCKERHUB_REPO_NAME` repo name of docker hub
    - `DOCKERHUB_TOKEN` token of docker hub user

- if use `wd_steps_transfer` just add `.woodpecker_kit.steps.transfer` at git ignore

### workflow usage

- workflow with backend `docker`

```yml
labels:
  backend: docker
steps:
  woodpecker-npm:
    image: sinlov/woodpecker-npm:latest
    pull: false
    settings:
      # debug: true
      ## registry settings if not will use https://registry.npmjs.org/
      # registry:
      ## NPM username
      npm-username:
        from_secret: npm_publish_username
      ## NPM password
      npm-password:
        from_secret: npm_publish_password
      ## NPM email
      npm-email:
        from_secret: npm_publish_email
      ## folder containing package.json, empty will use workspace
      # npm-folder:
      ## NPM scoped package access
      # npm-access: foo
```

- workflow with backend `local`, must install at local and effective at evn `PATH`
- install at ${GOPATH}/bin, latest

```bash
go install -a github.com/woodpecker-kit/woodpecker-npm/cmd/woodpecker-npm@latest
```

- install at ${GOPATH}/bin, v1.0.0

```bash
go install -v github.com/woodpecker-kit/woodpecker-npm/cmd/woodpecker-npm@v1.0.0
```

```yml
labels:
  backend: local
steps:
  woodpecker-npm:
    image: woodpecker-npm
    settings:
      # debug: true
      ## registry settings if not will use https://registry.npmjs.org/
      # registry:
      ## NPM username
      npm-username:
        from_secret: npm_publish_username
      ## NPM password
      npm-password:
        from_secret: npm_publish_password
      ## NPM email
      npm-email:
        from_secret: npm_publish_email
      ## folder containing package.json, empty will use workspace
      # npm-folder:
      ## NPM scoped package access
      # npm-access: foo
```

### settings.debug

- if open `settings.debug` will try file browser use `override` for debug.

### full config

```yaml
labels:
  backend: docker
steps:
  woodpecker-npm:
    image: sinlov/woodpecker-npm:latest
    pull: false
    settings:
      # debug: true
      ## registry settings if not will use https://registry.npmjs.org/
      # registry:
      ## NPM username
      npm-username:
        from_secret: npm_publish_username
      ## NPM password
      npm-password:
        from_secret: npm_publish_password
      ## NPM token to use when publishing packages. if token is set, username and password will be ignored.
      npm-token:
        from_secret: npm_publish_token
      ## NPM email
      npm-email:
        from_secret: npm_publish_email
      ## NPM tag to use when publishing packages. this will cover package.json version field.
      npm-tag: latest
      ## NPM scoped package access
      npm-access: foo
      ## folder containing package.json, empty will use workspace
      # npm-folder:
      ## fail NPM publish if version already exists in NPM registry
      npm-fail-on-version-conflict: true
      ## disables ssl verification when communicating with the NPM registry.
      npm-skip-verify-ssl: true
      ## Skip npm whoami check
      npm-skip-whoami: true
```

---

- want dev this project, see [doc](doc/README.md)