---
name: woodpecker-feishu-group-robot
description: woodpecker plugin template
author: woodpecker-kit
tags: [ node, npm ]
containerImage: sinlov/woodpecker-npm
containerImageUrl: https://hub.docker.com/r/sinlov/woodpecker-npm
url: https://github.com/woodpecker-kit/woodpecker-npm
icon: https://codeberg.org/woodpecker-plugins/node-pm/media/branch/main/nodejs-logo-hexagon.png
---

woodpecker-npm

## before use

- see [package.json define](https://docs.npmjs.com/files/package.json)
- take `username` and `password` or `token` to publish npm package

must set `package.json`
by [npm docs package-json.publishconfig](https://docs.npmjs.com/files/package.json#publishconfig)
args [registry](https://docs.npmjs.com/cli/v10/using-npm/config#registry)

```json
{
  "publishConfig": {
    "registry": "https://registry.npmjs.org/"
  }
}
```

## Settings

| Name           | Required | Default value | Description                                                                                        |
|----------------|----------|---------------|----------------------------------------------------------------------------------------------------|
| `debug`        | **no**   | *false*       | open debug log or open by env `PLUGIN_DEBUG`                                                       |
| `npm-registry` | **no**   | *none*        | NPM registry settings if empty will use https://registry.npmjs.org/                                |
| `npm-username` | **yes**  | *none*        | NPM username                                                                                       |
| `npm-password` | **yes**  | *none*        | NPM password                                                                                       |
| `npm-token`    | **yes**  | *none*        | NPM token to use when publishing packages. if token is set, username and password will be ignored. |
| `npm-email`    | **yes**  | *none*        | NPM email                                                                                          |
| `npm-tag`      | **no**   | *latest*      | NPM publish tag will cover package.json settings                                                   |
| `npm-folder`   | **no**   | *none*        | folder containing package.json, empty will use workspace                                           |
| `npm-access`   | **no**   | *none*        | NPM scoped package access                                                                          |

**custom settings**

| Name                           | Required | Default value | Description                                                         |
|--------------------------------|----------|---------------|---------------------------------------------------------------------|
| `npm-fail-on-version-conflict` | **no**   | *false*       | fail NPM publish if version already exists in NPM registry          |
| `npm-skip-verify-ssl`          | **no**   | *false*       | disables ssl verification when communicating with the NPM registry. |
| `npm-skip-whoami`              | **no**   | *false*       | Skip npm whoami check                                               |

**Hide Settings:**

| Name                                        | Required | Default value                    | Description                                                                      |
|---------------------------------------------|----------|----------------------------------|----------------------------------------------------------------------------------|
| `timeout_second`                            | **no**   | *10*                             | command timeout setting by second                                                |
| `woodpecker-kit-steps-transfer-file-path`   | **no**   | `.woodpecker_kit.steps.transfer` | Steps transfer file path, default by `wd_steps_transfer.DefaultKitStepsFileName` |
| `woodpecker-kit-steps-transfer-disable-out` | **no**   | *false*                          | Steps transfer write disable out                                                 |

## Example

- workflow with backend `docker`

[![docker hub version semver](https://img.shields.io/docker/v/sinlov/woodpecker-npm?sort=semver)](https://hub.docker.com/r/sinlov/woodpecker-npm/tags?page=1&ordering=last_updated)
[![docker hub image size](https://img.shields.io/docker/image-size/sinlov/woodpecker-npm)](https://hub.docker.com/r/sinlov/woodpecker-npm)
[![docker hub image pulls](https://img.shields.io/docker/pulls/sinlov/woodpecker-npm)](https://hub.docker.com/r/sinlov/woodpecker-npm/tags?page=1&ordering=last_updated)

```yml
labels:
  backend: docker
steps:
  woodpecker-npm:
    image: sinlov/woodpecker-npm:latest
    pull: false
    settings:
      # debug: true
      ## registry settings if empty will use https://registry.npmjs.org/
      # npm-registry: https://verdaccio.foo.com
      npm-username: # NPM username
        from_secret: npm_publish_username
      npm-password: # NPM password
        from_secret: npm_publish_password
      npm-email: # NPM email
        from_secret: npm_publish_email
      npm-tag: latest # NPM publish tag will cover package.json settings
      # npm-folder: . # folder containing package.json, empty will use workspace
      ## NPM scoped package access
      # npm-access: foo
```

- workflow with backend `local`, must install at local and effective at evn `PATH`

[![GitHub latest SemVer tag)](https://img.shields.io/github/v/tag/woodpecker-kit/woodpecker-npm)](https://github.com/woodpecker-kit/woodpecker-npm/tags)
[![GitHub release)](https://img.shields.io/github/v/release/woodpecker-kit/woodpecker-npm)](https://github.com/woodpecker-kit/woodpecker-npm/releases)

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
      ## registry settings if empty will use https://registry.npmjs.org/
      # npm-registry: https://verdaccio.foo.com
      npm-username: # NPM username
        from_secret: npm_publish_username
      npm-password: # NPM password
        from_secret: npm_publish_password
      npm-email: # NPM email
        from_secret: npm_publish_email
      npm-tag: latest # NPM publish tag will cover package.json settings
      # npm-folder: . # folder containing package.json, empty will use workspace
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
      ## registry settings if empty will use https://registry.npmjs.org/
      # npm-registry: https://verdaccio.foo.com
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

