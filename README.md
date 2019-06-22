# mgosd

Tiny concurrent and scheduled dump creator for MongoDB.

## Usage

There are two ways of using **mgosd:**

1. Passing configuration as CLI arguments:
```bash
$ ./mgosd -h
+----------------------------------------+
| mgosd (c) Ildar Karymov, 2019          |
| https://github.com/ildarkarymoff/mgosd |
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