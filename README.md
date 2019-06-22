# mgosd

Tiny concurrent and scheduled dump creator for MongoDB.

Current stable version: 1.1.0.

Licensed under GNU GPLv3.

## Table of Contents
1. [Installation](#installation)
2. [Usage](#usage)

<a name="installation" />

## Installation

If you're on Linux, just clone this repo and then you can find Linux executable in `bin/` directory.

If you're on Mac or Windows:
1. Install Go compiler (1.12.6+)
2. Install `mgo` dependency by running `go get sgithub.com/globalsign/mgo`
3. Compile **mgosd** by running `go build main.go`
4. Done!

<a name="usage" />

## Usage

There are two ways of using **mgosd:**

1. Configuration via CLI arguments:
```bash
$ ./mgosd -h
+----------------------------------------+
| mgosd (c) Ildar Karymov, 2019          |
| https://github.com/ildarkarymoff/mgosd |
| License: GNU GPLv3                     |
| Version: 1.1.0                         |
+----------------------------------------+
Usage of ./mgosd:
  -db string
    	Database name
  -host string
    	Database server address (default "127.0.0.1")
  -i string
    	Interval of dumping (default "12h")
  -login string
    	Database username (default "<empty>")
  -o string
    	Path to output directory (default $HOME)
  -password string
    	Database user password (default "<empty>")
  -port int
    	Database server port (default 27017)
```
2. Using JSON configuration:
```json
{
  "collections": [
    "users",
    "comments",
    "posts"
  ],
  "interval": "5s",
  "output": "dumps_test",
  "db": {
    "database": "mydb",
    "host": "localhost",
    "port": 27017
  }
}
```
```bash
$ ./mgosd config.json
```