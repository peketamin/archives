go-digital-signage
========

[![Build Status](https://travis-ci.org/peketamin/go-digital-signage.svg?branch=master)](https://travis-ci.org/peketamin/go-digital-signage)

Overview
--------
Simple web server for pages like digital-signage.
Display some clipped HTMLs and images randomly in a DB.

This sample product **made for myself to study golang** :flushed: .


Description
--------
Pages in sqlite3 database are displayed in cycles.
Next page is picked up from remaining pages those are not displayed yet in a round.


Demo
--------
![Anime Gif](https://github.com/peketamin/go-digital-signage/blob/5dd11f4e8f294879b4907d17abab9cdc1ba914f0/demo.gif)


Usage
--------
Open `http://localhost:8080` on your browser and access to,

- `/`: view page randomly that pickuped from never watched pages in a round.
- `/add`: add new data
- `/[0-9]+`: view a specified page
- `/edit/[0-9]+`: edit a specified page


### How to modify the duration of displaying pages
Add parameter when start the program.
This means rotation pages every 15 seconds.

```bash
$ ./go-digital-signage 15
```

(Default: 86400 sec.)


<!-- 
VS.
--------
-->


Requirements
--------
### [The Go Programming Language](http://golang.org/)
as known as *golang*.

### [SQLite3](http://www.sqlite.org/).

Download from [SQLite Download Page](http://www.sqlite.org/download.html) precompiled archive and install them.

User of Linux or Mac may also choose following ways,

- Linux: using package manager of your distribution.
- Mac: using [Homebrew — The missing package manager for OS X](http://brew.sh/) or [The MacPorts Project -- Home](http://www.macports.org/) 


Install
--------
Download zip or `git clone` this.

```bash
$ cd go-digital-signage
$ go build
$ ./go-digital-signage
```

<!-- 
Contribution
--------

1. Fork ([https://github.com/peketamin/go-digital-signage/fork](https://github.com/peketamin/go-digital-signage/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request
-->


TODO
--------
- Listing items page
- Disable flag into page data struct
- Authentification


Libraries
--------
- [github.com/jinzhu/gorm](https://github.com/jinzhu/gorm)
- [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
- [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
- [github.com/lib/pq](https://github.com/lib/pq)


Licence
--------
[MIT](https://github.com/tcnksm/tool/blob/master/LICENCE)


Author
--------
[peketamin](https://github.com/peketamin)


Thanks
----
- [tk0miya](https://twitter.com/tk0miya) ([blockdiag](http://blockdiag.com) developer)
- japanese: [わかりやすいREADME.mdを書く | SOTA](http://deeeet.com/writing/2014/07/31/readme/)
