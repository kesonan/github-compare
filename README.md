![banner](./resource/banner.png)

# github-compare

A Command-line tool to statistics the GitHub repositories

## Install

```bash
$ go install github.com/anqiansong/github-compare@latest
```

## Example

### TableView

```bash
$ github-compare spf13/cobra urfave/cli junegunn/fzf antonmedv/fx
```
![preview](./resource/compare-preview.png)

### JSONView

```bash
$ github-compare spf13/cobra --json                                           
[
  {
    "age": "3187 days",
    "avgReleasePeriod": "199 days",
    "contributorCount": "246",
    "forkCount": "2327(0/d)",
    "fullName": "spf13/cobra",
    "homepage": "https://cobra.dev",
    "issue": "107/892",
    "language": "Go",
    "lastPushedAt": "8 hour(s) ago",
    "latestReleaseAt": "2 month(s) ago",
    "lastUpdatedAt": "1 hour(s) ago",
    "latestDayStarCount": "14 ⇊",
    "latestMonthStarCount": "477",
    "latestWeekStarCount": "110 ⇈",
    "license": "Apache License 2.0",
    "pull": "55/808",
    "releaseCount": "16",
    "starCount": "26774(8/d)",
    "watcherCount": "350"
  }
]
```

### YAMLView

```bash
$ github-compare spf13/cobra --yaml                                           
- age: 3187 days
  avgreleaseperiod: 199 days
  contributorcount: "246"
  forkcount: 2327(0/d)
  fullname: spf13/cobra
  homepage: https://cobra.dev
  issue: 107/892
  language: Go
  lastpushedat: 8 hour(s) ago
  latestreleaseat: 2 month(s) ago
  lastupdatedat: 1 hour(s) ago
  latestdaystarcount: 14 ⇊
  latestmonthstarcount: "477"
  latestweekstarcount: 110 ⇈
  license: Apache License 2.0
  pull: 55/808
  releasecount: "16"
  starcount: 26774(8/d)
  watchercount: "350"
```

### Export as a csv file

```bash
$ github-compare spf13/cobra urfave/cli junegunn/fzf antonmedv/fx -f data.csv
```
![csv](./resource/compare-csv.png)


## Usage

### Preparation

1. [Creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
2. Set access token token
   - Copied the access token and export to environment
   - Or you can set the access token by passing `github-compare` argument 
     `--token`

### Execute

```bash
# set access token through the flag --token
# github-compare zeromicro/go-zero --token ${accessToken}
# or export access token to environment (recommended)
$ export GITHUB_ACCESS_TOKEN=${GITHUB_ACCESS_TOKEN}
$ github-compare zeromicro/go-zero
```

### Commands

```bash
$ github-compare -h                                                    
A cli tool to compare two github repositories

Usage:
  github-compare [flags]

Flags:
  -f, --file string    output to a specified file
  -h, --help           help for github-compare
      --json           print with json style
      --table          print with table style(default) (default true)
  -t, --token string   github access token
      --yaml           print with yaml style
```

## Note

1. A GitHub personal access token is required.
2. `github-compare` accepts 1 to 4 repositories data queries.
3. If you prefer to export the access token to environment, you must use 
   environment key `GITHUB_ACCESS_TOKEN`

## Last

If this repository can help you, give a star please! 

Thanks all!

## License

[MIT License](License)