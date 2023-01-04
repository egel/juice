<p align="center"><img alt="The Juice" width="175px" src="./docs/assets/juice.svg"/></p>

# Juice

Quick and easy tool to help extract licensing information on production node packages.

>  -- Yeah! Give me the Juice!

---

This tool is aiming to simplify and reduce your time with gathering informations about all your application's node production dependencies in single summary file, with all vital information about each package, like: *name*, *version*, *license type*, *license text*, *links to NPM/repository/homepage*, *errors*, or *direct license link (experimental)*.


## Install

```bash
go install github.com/egel/juice/cmd/juice-cli@latest

# check if installed correctly
which juice-cli
```


## Usage

```bash
# enter the location with package.json and package-lock.json files, and run:
juice-cli
```

> If `juice-cli` is failing, possibly your project may have old, incorrect or outdated dependencies. Make sure you have a working installation of `package.json` and `package-lock.json`, so it can download all required packages into `node_modules`.


## License

MIT License. Logo by [Vincent Le moign](https://iconscout.com/icon/juice-247)
