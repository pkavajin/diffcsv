# diffcsv

diffcsv is a small utility to diff two csv files and get the additions, deletions and equal rows as separate csv files.

## Installation

Download the latest binary from the [releases page](https://github.com/pkavajin/diffcsv/releases).

## Usage

```
Usage of ./diffcsv-darwin-amd64-v0.1.0:
  -a string
    	first file
  -b string
    	second file
  -header
    	file contains header (default true)
  -out-added string
    	file where the rows should be written which are in b but not in a (aka: added) (default "added.csv")
  -out-both string
    	file where the rows should be written which are in both b and a (aka: equal (default "both.csv")
  -out-deleted string
    	file where the rows should be written which are in a but not in b (aka: deleted (default "deleted.csv")
```

## Example

```
./diffcsv-darwin-amd64-v0.1.0 -a old.csv -b new.csv
```

Creates three files:
* added.csv - Contains the rows which are in b but not in a
* both.csv - Contains the rows which are in both a and b
* deleted.csv - Contains the rows which are only in a

## License
[MIT](LICENSE)