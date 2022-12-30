# Juice!

Give me the juice! - Simple and quick tool to help extract production licenses from the npm package.


## Usage

```bash
#copy project
git clone git@github.com:egel/juice.git

# build version for your OS
make build

# copy generated 'bin/juice-cli/juice' program to location with 'package.json' and 'package-lock.json' files
# and execute
./bin/juice-cli/juice
```


## FAQ

> **cmd fatal: pipe: too many open files**
>
> Program use the concurrency to get all data asap to you.
> By having this error you probably reach terminal session's limit for open files. You can increase it by: `ulimit -u 4096`


## License

MIT License
