# PHARGO

PHAR-files reader and parser written in golang

## Info

Parser supports several signature algorithms:
* MD5
* SHA1
* SHA256
* SHA512
* ~~PGP~~
* ~~OPENSSL~~

Can read manifest version, alias and metadata. For every file inside PHAR-archive can read it contents, 
name, timestamp and metadata. Checks file CRC and signature of entire archive.

## Installation

1. Download and install:

```sh
$ go get -u github.com/deadkrolik/phargo
```

2. Import and use it:

```go
package main

import (
    "log"
    "time"

    "github.com/deadkrolik/phargo"
)

func main() {
    r := phargo.NewReader()
    
    //some limitations
    if false {
        r.SetOptions(phargo.Options{
            //metadata of every file and entire archive can be more than that number
            MaxMetaDataLength: 10000,
            //max length of first block of archive when looking for "HALT_COMPILER" string 
            MaxManifestLength: 1048576 * 100,
            //max length of name of the file
            MaxFileNameLength: 1000,
            //max length of archive alias in manifest
            MaxAliasLength: 1000,
        })
    }
    
    file, err := r.Parse("file.phar")
    if err != nil {
        log.Println(err)
        return
    }
    
    log.Println("Manifest version: ", file.Version)
    log.Println("File alias: ", file.Alias)
    log.Println("Serialized metadata: ", string(file.Metadata))
    
    for _, f := range file.Files {
        log.Println("File name: ", f.Name)
        log.Println("File metadata: ", string(f.Metadata))
        log.Println("File data len: ", len(f.Data))
        log.Println("File date: ", time.Unix(f.Timestamp, 0).String())
        log.Println("---")
    }
}
```

## Running the tests

Just run the command:

```sh
$ go test
```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
