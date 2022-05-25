# github-compare

A Command-line tool to statistics the GitHub repositories

## Install

```bash
$ go install github.com/anqiansong/github-compare
```

## Example
```bash
$ github-compare zeromicro/go-zero go-kratos/kratos micro/micro go-kit/kit
```
![preview](./resource/compare-preview.png)

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