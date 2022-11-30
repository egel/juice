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

## License

MIT License