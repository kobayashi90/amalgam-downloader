# Amalgam-Fansub Detektiv Conan Downloader
This project is a small cli download tool to download the newer Detektiv Conan
episodes from https://amalgam-fansubs.moe/. 

# Installation
## Download releases
Download the binary for your system from the [release page](https://gitlab.com/mauamy/amalgamdetektivconandownloader/-/tags).

## Build it yourself
If you want to build it yourself, simply clone this repository and use the makefile
for building and installing it.
```bash
# local/test linux build
$ make build

# install it to your GOPATH (linux)
$ make install

# uninstall from your GOPATH (linux)
$ make uninstall

# build for windows
$ make windows

# build for mac
$ make mac
```

# Usage
#### Show Help
```bash
$ adcl -h
```
#### List Episodes
```bash
$ adcl list
$ adcl l
```
#### Download Episodes
```bash
$ adcl download <episode_numbers>
$ adcl d <episode_numbers>
```
*<episode_numbers>* is a separated list of episode numbers: **1 2 3 4**.
In addition you can provide ranges within this list: **1 2 3-8 10**. 
```bash
# example
$ adcl d 710 840-845 870
```
