# Holochain

[![Code Status](https://img.shields.io/badge/Code-Pre--Alpha-orange.svg)](https://github.com/metacurrency/holochain/milestones?direction=asc&sort=completeness&state=all)
[![Travis](https://img.shields.io/travis/metacurrency/holochain/master.svg)](https://travis-ci.org/metacurrency/holochain/branches)
[![Go Report Card](https://goreportcard.com/badge/github.com/metacurrency/holochain)](https://goreportcard.com/report/github.com/metacurrency/holochain)
[![Gitter](https://badges.gitter.im/metacurrency/holochain.svg)](https://gitter.im/metacurrency/holochain?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=body_badge)
[![In Progress](https://img.shields.io/waffle/label/metacurrency/holochain/in%20progress.svg)](http://waffle.io/metacurrency/holochain)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![Twitter Follow](https://img.shields.io/twitter/follow/holochain.svg?style=social&label=Follow)](https://twitter.com/holochain)

**Holographic storage for distributed applications.** A holochain is a monotonic distributed hash table (DHT) where every node enforces validation rules on data before publishing that data against the signed chains where the data originated.

In other words, a holochain functions very much **like a blockchain without bottlenecks** when it comes to enforcing validation rules, but is designed to  be fully distributed with each node only needing to hold a small portion of the data instead of everything needing a full copy of a global ledger. This makes it feasible to run blockchain-like applications on devices as lightweight as mobile phones.

**[Code Status:](https://github.com/metacurrency/holochain/milestones?direction=asc&sort=completeness&state=all)** Active development for **proof-of-concept stage**. Pre-alpha. Not for production use. We still expect to destructively restructure data chains at this time. These instructions are really for developers who want to build distributed apps on holochains, not so much for end users who should probably use a nice packaged installation.
<br/>

| Holochain Links: | [FAQ](https://github.com/metacurrency/holochain/wiki/FAQ) | [Developer Wiki](https://github.com/metacurrency/holochain/wiki) | [White Paper](http://ceptr.org/projects/holochain) | [GoDocs](https://godoc.org/github.com/metacurrency/holochain) |
|---|---|---|---|---|

**Table of Contents**

<!-- auto-generated by the markdown-toc addon for atom -->
<!-- TOC depthFrom:2 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Installation](#installation)
	- [Unix](#unix)
	- [Windows](#windows)
- [Docker Installation](#docker-installation)
- [Usage](#usage)
	- [Setting up a Holochain](#setting-up-a-holochain)
	- [Cloning Application DNA](#cloning-application-dna)
	- [Testing your Application](#testing-your-application)
	- [Generate New Chain From DNA](#generate-new-chain-from-dna)
	- [Accessing the Web UI](#accessing-the-web-ui)
	- [File Locations](#file-locations)
	- [Logging](#logging)
- [Architecture Overview and Documentation](#architecture-overview-and-documentation)
- [Development](#development)
	- [Contribute](#contribute)
	- [Dependencies](#dependencies)
	- [Tests](#tests)
- [Social](#social)
- [License](#license)
- [Acknowledgements](#acknowledgements)

<!-- /TOC -->

## Installation
Eiher use docker installation, or once you have a working environment (see OS specific instructions) you can install the Holochain command line interface with:
```bash
$ go get -d github.com/metacurrency/holochain
$ cd $GOPATH/src/github.com/metacurrency/holochain
$ make
```

Now that you've installed `hc` you'll need to do the first time configuration by running (just substitute your own email address.):
```bash
$ hc init 'name@example.org'
```
As a general user, you should only need to do this once, but as a developer, you will need to do this if you remove your `.holochain` directory during testing and such.

### Unix
(Unix includes macOS and Linux.)
In order to use Holochain, you'll need to have a working environment set up for [Go](http://golang.org) version 1.7 or later. See the [installation instructions for Go](http://golang.org/doc/install.html).

Most importantly you'll need to: (Almost all installation problems that have been reported stem from skipping one of these steps.)
1. Export the `$GOPATH` variable in your shell profile.
2. Add `$GOPATH/bin` to your `$PATH` in your shell profile.

For example, add the following to the end of your shell profile (usually `~/.bashrc` or `~/.bash_profile`):

    export GOPATH=`$HOME/go`
    export PATH=$GOPATH/bin:$PATH

### Windows
First you'll need to install some necessary programs if you don't already have them.
* [Install Go](https://golang.org/dl/) 1.7 or later.
* [Install Windows git](https://git-scm.com/downloads). Be sure to select the appropriate options so that git is accessible from the Windows command line.
* [Install GnuWin32 make](http://gnuwin32.sourceforge.net/packages/make.htm#download).

Next, in your Control Panel go to *System>Advanced system settings>Environment Variables...* and under *System Variables* do the following:
1. Add a new entry with the name `GOPATH` and the value `%USERPROFILE%\go` (Or your Go workspace folder).
2. Double-click Path, and in the window that pops up add the following entries:
    - `%GOPATH%\bin`
    - `C:\Go\bin` (Or wherever you installed Go to+`\bin`).
    - `C:\Program Files (x86)\GnuWin32\bin` (Or wherever you installed GnuWin32 make to+`\bin`).

## Docker Installation

Docker installation of the holochain core is suitable for holochain users, holochain developers and core developers

* Install the latest version of Docker for your OS. 
  * See [Docker Getting Started](https://docs.docker.com/engine/getstarted/step_one/)

* get our holochain repository from github:
  ```bash
  $ git clone https://github.com/metacurrency/holochain.git holochain
  $ cd holochain
  ```

* build the holochain core with all dependencies
  ```bash
  $ Scripts/docker.build
  ```
* to run holochain in your new environment, suitable to continue this walkthrough [usage](#usage)
  * to exit the holochain environment, press `Ctrl-D` or type `exit`

  ```bash
  $ Scripts/docker.run
  ```

#### docker for core development
* fork our github repository [https://github.com/metacurrency/holochain](https://github.com/metacurrency/holochain)
* use Scripts/docker.build to compile and test your latest source
* docker builds will overwrite the local docker metacurrency/holochain image

#### docker for holochain development
* checkout the holochain skeleton github repository [https://github.com/metacurrency/holoSkel](https://github.com/metacurrency/holoSkel)
* build and test scripts will use your local metacurrency/holochain image

## Usage
Once you've gotten everything working as described [above](#installation) you may want to use the `hc` command.

Since Holochain is intended to be used underneath distributed applications, end users won't have to do much through the command line.

For the most up-to-date information on how to use `hc`, run `hc help` or for sub-commands run `hc <COMMAND> help`.

For more detailed information, see [the wiki page](https://github.com/metacurrency/holochain/wiki/hc-Command)

### Setting up a Holochain
You've installed and built the distributed data integrity engine, but you haven't set up an application running on it yet. The basic flow involved in getting a chain running looks like this:

1. `hc clone`
2. `hc test`
3. `hc gen chain`
4. `hc web`

Instructions for each of these steps are below...

### Cloning Application DNA
You can load a pre-existing Holochain application DNA by running the following command (replacing SOURCE_PATH with a path to an application's DNA and CHAIN_NAME with the name you'd like it to be stored as).
```bash
$ hc clone <SOURCE_PATH> <CHAIN_NAME>
```
For example: `hc clone ./examples/sample sample`

You can source from files anywhere; such as a git repo you've cloned, a live chain you're already running in your `.holochain` directory, or one of the examples included in this repository.

Before you launch your chain, this is the chance for you to customize the application settings like the NAME, and the UUID

### Testing your Application
We have designed Holochain around test-driven development, so the DNA should contain tests to confirm that the rest of the DNA is functional.

You can run a chain's tests with:
```bash
$ hc test <CHAIN_NAME>
```
If you're a developer, you should be running this command as you make changes to your DNA files to leverage test-driven development. If the tests fail, then you know your application DNA is broken and you shouldn't think that your chain is going to work. And obviously, please do not send out applications that don't pass their own tests.

### Generate New Chain From DNA
After you have cloned a chain, you need to generate the genesis entries which start your new chain in order to use it.
```bash
$ hc gen chain <CHAIN_NAME>
```
The first entry is the DNA which is the hash of all the application code which confirms every person's chain starts with the the same code/DNA. The second block registers your keys so you have an address, identity, and signing keys for communicating on the chain.

### Accessing the Web UI
Holochains serve their UI via local web sockets. This let's interface developers have a lot of freedom to build HTML/JavaScript files and drop them in that chain's UI directory. You launch the web UI with:
```bash
$ hc web <CHAIN_NAME> [PORT]
```
In a web browser you can go to `localhost:<port>` (defaults to `3141`) to access UI files and send and receive JSON with exposed application functions.

### File Locations
By default `hc` stores all holochain data and configuration files to the `~/.holochain` directory.  You can override this with the `-path` flag or by setting the `HOLOPATH` environment variable, e.g.:
```bash
$ hc -path ~/mychains init '<my@other.identity>'
$ HOLOPATH=~/mychains hc
```
You can use the form: `hc -path=/your/path/here` but you must use the absolute path, as shell substitutions will not happen.

### Logging
The `-debug` flag will turn on a number of different kinds of debugging. You can also control exactly which of these logging types you wish to see in the chain's config.json file.  You can also set the DEBUG environment variable to 0 or 1 to temporarily override your settings to turn everything on or off.

## Architecture Overview and Documentation
Most architecture information is in the [Holochain Wiki](https://github.com/metacurrency/holochain/wiki/Architecture).

You can also look through auto-generated [reference API on GoDocs](https://godoc.org/github.com/metacurrency/holochain)

## Development
[![In Progress](https://img.shields.io/waffle/label/metacurrency/holochain/in%20progress.svg)](http://waffle.io/metacurrency/holochain)

We accept Pull Requests and welcome your participation.

Some helpful links:
* Come [chat with us on gitter](https://gitter.im/metacurrency/holochain)
* View our [Kanban on Waffle](https://waffle.io/metacurrency/holochain).
* View our  [Milestone](https://github.com/metacurrency/holochain/milestones?direction=asc&sort=due_date&state=all) progress.

If you'd like to get involved you can:
* Contact us on [Gitter](https://gitter.im/metacurrency/holochain) to set up a **pair coding session** with one of our developers to learn the lay of the land.
* **join our dev documentation calls** twice weekly on Tuesdays and Fridays.

Throughput graph:

[![Throughput Graph](http://graphs.waffle.io/metacurrency/holochain/throughput.svg)](https://waffle.io/metacurrency/holochain/metrics)

### Contribute
If you're going to contribute to our project we expect you to adhere to the following guidelines:

<!-- * Protocols for Inclusion. -->
We are committed to foster a vibrant thriving community, including growing a culture that breaks cycles of marginalization and dominance behavior. In support of this, some open source communities adopt [Codes of Conduct](http://contributor-covenant.org/version/1/3/0/).  We are still working on our social protocols, and empower each team to describe its own *Protocols for Inclusion*.  Until our teams have published their guidelines, please use the link above as a general guideline.

We use **test driven development**. When you add a new function or feature, be sure to add the tests that make sure it works.

All Go code should be formatted with [gofmt](https://blog.golang.org/go-fmt-your-code).
To make this easier consider using a [git-hook](https://gist.github.com/timotree3/d69b0fb90c8affbd705765abeabc489d#file-pre-commit) or configuring your editor with one of these: ([Emacs][], [vim][], [Sublime][], [Eclipse][])
[Emacs]: https://github.com/dominikh/go-mode.el
[vim]: https://github.com/fatih/vim-go
[Sublime]: https://github.com/DisposaBoy/GoSublime
[Eclipse]: https://github.com/GoClipse/goclipse

For Atom, you could try this [package](https://atom.io/packages/save-commands) but it requires some configuration.

### Dependencies
This project depends on various parts of [libp2p](https://github.com/libp2p/go-libp2p), which uses the [gx](https://github.com/whyrusleeping/gx) package manager. All of which will be automatically installed by make by following the [setup instructions](#installation) above.

### Tests
To compile and run all the tests:
```bash
$ cd $GOPATH/github.com/metacurrency/holochain
$ make test
```
If you want to use `go test` instead of `make test`, you'll need to do a couple extra things because of this project's dependency on `gx`:
* Before running `go test` you need to run `make work` to configure the imports properly.
* If you do this, before commiting you must also run `make pub` to revert the changes it makes.

## License
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)

Copyright (C) 2017, The MetaCurrency Project (Eric Harris-Braun, Arthur Brock, et. al.)

This program is free software: you can redistribute it and/or modify it under the terms of the license provided in the LICENSE file (GPLv3).  This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

**Note:** We are considering other 'looser' licensing options (like MIT license) but at this stage are using GPL while we're getting the matter sorted out.

## Acknowledgements
* **MetaCurrency & Ceptr**: Holochains are a sub-project of [Ceptr](http://ceptr.org) which is a semantic, distributed computing platform under development by the [MetaCurrency Project](http://metacurrency.org).
&nbsp;
* **Ian Grigg**: Some of our initial plans for this architecture were inspired in 2006 by [his paper about Triple Entry Accounting](http://iang.org/papers/triple_entry.html) and his work on [Ricardian Contracts](http://iang.org/papers/ricardian_contract.html).
&nbsp;
* **Juan Benet**: For all his work on IPFS and being a generally cool guy. Various functions like multihash, multiaddress, and such come from IPFS as well as the libP2P library which helped get peered node communications up and running.
&nbsp;
* **Crypto Pioneers** And of course the people who paved the road before us by writing good crypto libraries and *preaching the blockchain gospel*. Back in 2008, nobody understood what we were talking about when we started sharing our designs. The main reason people want it now, is because blockchains have opened their eyes to the power of decentralized architectures.
