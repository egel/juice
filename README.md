![GitHub License](https://img.shields.io/github/license/egel/juice)
[![Go Reference](https://pkg.go.dev/badge/github.com/egel/juice.svg)](https://pkg.go.dev/github.com/egel/juice)
[![Go Report Card](https://goreportcard.com/badge/github.com/egel/juice)](https://goreportcard.com/report/github.com/egel/juice)

<p align="center"><img alt="The Juice" width="175px" src="./docs/assets/juice.svg"/></p>

# Juice

Quick and easy tool to help extract licensing information on production node packages.

> -- Yeah! Give me the Juice!

![Sample demo](./docs/screenshots/2023-01-05_demo_loop.gif)

This tool is aiming to simplify and reduce your time with gathering informations about all your application's node production dependencies in single summary file, with all vital information about each package, like: _name_, _version_, _license type_, _license text_, _links to NPM/repository/homepage_, _errors_, or _direct license link (experimental)_.

## Install

```bash
go install github.com/egel/juice/cmd/juice-cli@latest
```

> Check if juice has been installed correctly
>
> ```bash
> which juice-cli
> ```

### Upgrade

To upgrade to latest version use

```bash
GONOPROXY=github.com/egel go install github.com/egel/juice/cmd/juice-cli@latest
```

## Usage

```bash
# enter the location with package.json and package-lock.json files, and run:
juice-cli get
```

> If `juice-cli` is failing, possibly your project may have old, incorrect, or outdated dependencies. Make sure you have a working installation of `package.json` and `package-lock.json`, so it can download all required packages into `node_modules`.

## License

MIT License. Logo by [Vincent Le moign](https://iconscout.com/icon/juice-247)
