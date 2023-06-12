<h1 align="center">
  <br>
  <a href="https://leofvo.github.io/gosurp/"><img src="docs/public/logo.png" width="200px" alt="gosurp"></a>
</h1>
<h4 align="center">GoSurp is a cli tool to work on email and DNS usurpation. It is written in Go and is distributed as a single binary without any external dependencies.</h4>
<p align="center">
<img src="https://img.shields.io/github/go-mod/go-version/leofvo/gosurp?filename=go.mod">
<a href="https://github.com/leofvo/gosurp/releases"><img src="https://img.shields.io/github/downloads/leofvo/gosurp/total">
<a href="https://github.com/leofvo/gosurp/graphs/contributors"><img src="https://img.shields.io/github/contributors-anon/leofvo/gosurp">
<a href="https://github.com/leofvo/gosurp/releases/"><img src="https://img.shields.io/github/release/leofvo/gosurp">
<a href="https://github.com/leofvo/gosurp/issues"><img src="https://img.shields.io/github/issues-raw/leofvo/gosurp">
<a href="https://github.com/leofvo/gosurp/discussions"><img src="https://img.shields.io/github/discussions/leofvo/gosurp">
<a href="https://twitter.com/leofvo"><img src="https://img.shields.io/twitter/follow/leofvo.svg?logo=twitter"></a>

</p>
      
<p align="center">
  <a href="#how-to-use">How</a> •
  <a href="#install">Install</a> •
  <a href="https://leofvo.github.io/gosurp/">Documentation</a> •
</p>

## How to use

```bash
gosurp -h
```

This will display help for the tool. Here are all the switches it supports.

### Inspect DNS Record

First of all, Gosurp allow you to inpect DNS records.

```bash
gosurp inspect domain <yourdomain.com>
```

You can also look for specific DNS records using:

```bash
gosurp inspect (spf/dkim/dmarc) <yourdomain.com>
```

### Email usurpation

First of all, to allow you working on building the good mail without sending it, you can use the following command to start a local SMTP server that will print the mail content in the console like a debugger.

```bash
gosurp smtp listen -vvv
# or
gosurp smtp listen --hostname localhost --port 2525 -vvv
```

Then, you can use the following command to send a mail using the local SMTP server.

```bash
gosurp smtp send -vvv --from attacker@fsociety.local --to victim@localhost --port 2525
```

You can find some complete recipes in the [documentation](https://leofvo.github.io/gosurp/recipes).

## Install

If you have a [Go](https://golang.org/) environment ready to go (at least go 1.19), it's as easy as:

```bash
go install github.com/LeoFVO/gosurp@latest
```

PS: You need at least go 1.19 to compile gosurp.

<details>
  <summary>Docker</summary>
  ```bash
  docker pull ghcr.io/leofvo/gosurp:latest  
  docker run gosurp:latest
  ```
</details>

<details>
  <summary>Binary Releases</summary>
We are now shipping binaries for each of the releases so that you don't even have to build them yourself! How wonderful is that!

If you're stupid enough to trust binaries that I've put together, you can download them from the [releases](https://github.com/LeoFVO/gosurp/releases) page.

</details>
<details>
  <summary>Build from source</summary>

#### Prerequisites

Since this tool is written in [Go](https://golang.org/) you need to install the Go language/compiler/etc. Full details of installation and set up can be found [on the Go language website](https://golang.org/doc/install). Once installed you have two options. You need at least go 1.19 to compile gosurp.

#### Clone the repository

```bash
git clone git@github.com:LeoFVO/gosurp.git
```

#### Compiling

`gosurp` has external dependencies, and so they need to be pulled in first:

```bash
go get && go build
```

This will create a `gosurp` binary for you. If you want to install it in the `$GOPATH/bin` folder you can run:

```bash
go install
```

</details>

## Setup

### Documentation

The documentation is available at [https://leofvo.github.io/gosurp/](https://leofvo.github.io/gosurp/).

In order to deploy documentation for your project, you need to allow github actions to deploy github pages. To do so, go to your repository settings > Pages, and in the `Build and deployment` section, select `Github Actions` as the source.

# License

gosurp is distributed under [MIT License](https://github.com/leofvo/gosurp/blob/main/LICENSE)
