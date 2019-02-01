# Gogitignore

Simple gitignore downloader from https://github.com/github/redmeros

## Overview

What is very annoying for me is to creating gitignore files. When we do have very good templates in github repository.

This simple (and my first public package) is very simple. The only thing its doing is downloading and merging files from https://github.com/github/gitignore search is done through github public api, and is done only in main dir and in Global dir.

## Usage

```shell
gogitignore -c visualstudiocode -c go -S
```

This command will fetch
https://github.com/github/gitignore/blob/master/Global/VisualStudioCode.gitignore
and
https://github.com/github/gitignore/blob/master/Go.gitignore

and write it to a `.gitignore` file

## Installation

### Via go tool

Download package

```shell
go get github.com/redmeros/gogitignore
```

Install

```shell
go install github.com/redmeros/gogitignore
```

### Manual

download latest binary from
[releases page](https://github.com/redmeros/gogitignore/releases)

put this binary into your PATH

## Contribution

Feel free to contribute PR's are welcome