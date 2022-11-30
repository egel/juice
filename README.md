# Juice!

Give me the juice! - Simple and quick tool to help extract production licenses from the npm package.


## Usage

```bash
#copy project
git clone git@github.com:egel/juice.git

# build version for your OS
go build -o bin/juice-cli/ple cmd/juice-cli/main.go

# copy generated 'bin/juice-cli/ple' program to location with 'package.json' and 'package-lock.json' files
# and execute
./juice
```


## FAQ

> **cmd fatal: pipe: too many open files**
>
> Program use the concurrency to get all data asap to you.
> By having this error you probably reach terminal session's limit for open files. You can increase it by: `ulimit -u 4096`


## License

MIT License