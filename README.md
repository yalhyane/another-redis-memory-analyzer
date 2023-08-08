# arma - Another Redis Memory Analyzer

`arma` is a command Line Tool for analyzing Redis memory, it is designed to help Redis developers understand how DBs and keys are using memory in Redis instance.
This tool uses the Redis MEMORY USAGE command to retrieve the number of bytes that a key and its value are taking in RAM, and then group them by key pattern and by db.

## Installation

To install `arma`, you can either download the pre-built binaries from the [releases](https://github.com/yalhyane/another-redis-memory-analyzer/releases) page, or build from source by running the following commands:

```bash
git clone https://github.com/yalhyane/another-redis-memory-analyzer
cd another-redis-memory-analyzer
go build -o arma
```


## Usage
Once you have installed arma, you can run it from the command line with the following syntax:
```bash
arma analyze [OPTIONS]
```

Here are the available options:

```lua
Options:
    -a, --ask-password       Let tty asks for password (recommended)
    -D, --db int             Redis DB to analyse ( -1 will analyze all dbs ) (default -1)
    -d, --delimiter string   The delimiter of keys (default ":")
    -o, --format strings     Output format, supported values: , , table, json (default [table])
    -g, --goroutines int     Number of concurrent goroutines to analyze a redis DB (default 50)
    -h, --help               help for analyze
    -H, --host string        Redis server host (default "127.0.0.1")
    -l, --level int          The delimiter level of keys (default 1)
    -s, --min-size string    Minimum size of group of keys to show, if group of keys is less than this size it won't be shown. Human readable size (KB, MB, GB...) (default "1KB")
    -b, --no-progress-bar    No progress bar
    -p, --password string    Redis server password, it's recommended to use -a flag which will ask for your password which will prevent it from appearing in tty history
    -P, --port int           Redis server port (default 6379)
    -S, --scan int           Redis scan range length, number of keys to scan at once (default 500)
```
Here are some examples of how you can use redis-memory-analysis:

```bash
# To analyze memory usage of a Redis instance running on the default port on the local machine:
./arma analyze

# To analyze memory usage of a Redis instance running on a remote server with a custom port and password:
./arma analyze --host example.com --port 6380 --ask-password

# To analyze memory usage of a Redis instance and pattern keys by custom delimiter and a minimum level:
./arma analyze --delimiter "-" --level 2

# To analyze memory usage of a Redis instance on a specific database with larger scan range and only show big keys:
./arma analyze --db 3 --scan 10000 --min-size 1MB
```

## Contributing
If you find any bugs or issues with `arma`, please open a new issue on the [Issues](https://github.com/yalhyane/another-redis-memory-analyzer/issues) page.
Contributions are also welcome! If you'd like to contribute to the project, feel free to open a pull request with your changes.

## Credits
This tool was inspired by the work of [redis-memory-analyzer](https://github.com/hto/redis-memory-analyzer).

## License
yet-another-redis-memory-analyzer is licensed under the [MIT License](https://github.com/yalhyane/another-redis-memory-analyzer/blob/main/LICENSE).


