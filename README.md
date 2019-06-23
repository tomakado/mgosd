## Table of Contents
1. [About](#mgosd)
1. [Features](#features)
1. [Installation](#installation)
1. [Usage](#usage)
1. [Roadmap](#roadmap)

<a name="mgosd" />

# mgosd

Tiny concurrent and scheduled dump creator for MongoDB.

Why not using something like `cron`? So, cron is awesome Unix tool, but it's available on Unix systems only. **mgosd** allows you to do not waste your time on dealing with scheduled task systems in different environments.  

Current stable version: [v1.1.0](https://github.com/ildarkarymoff/mgosd/releases/tag/v1.1.0)

Licensed under GNU GPLv3.


<a name="features" />

## Features

* _Concurrent:_ mgosd process every collection in separate goroutine 
* _Scheduled:_ you can set different intervals of dumping  
* _Two ways of configuration:_ CLI arguments and JSON config file
* _Portable:_ you don't need to install MongoDB toolkit to use it
* _Cross-platform:_ unite configuration for all systems


<a name="installation" />

## Installation

If you're on Linux you can find binary executable [on Releases page](https://github.com/ildarkarymoff/mgosd/releases/). 

If you're on Mac or Windows:
1. Install Go compiler (1.12.6+)
2. Install `mgo` dependency by running `go get github.com/globalsign/mgo`
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

<a name="roadmap" />

## Roadmap

* Switch to official MongoDB driver for Go
* Ability to configure the schedule in absolute way (implicit time of day)
* Binary executables for Mac and Windows