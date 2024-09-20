## <a id="readme">README</a>

![Go language badge](http://shields.io/badge/language-Go-cyan.svg
"Go language badge")
![GPLv3+ license badge](http://shields.io/badge/license-GPLv3+-blue.svg
"GPLv3+ license badge")

<details>
<summary>Table of Contents</summary>
<p/>

**Table of Contents**
* <a href="#name" alt="name">Name</a>
* <a href="#version" alt="version">Version</a>
* <a href="#synopsis" alt="synopsis">Synopsis</a>
* <a href="#description" alt="description">Description</a>
* <a href="#installation" alt="installation">Installation</a>
* <a href="#funcs" alt="functions">Public Functions</a>
  * <a href="#capture" alt="capture output">CaptureOutput</a>
  * <a href="#filepath" alt="">FilePathInCwd</a>
  * <a href="#ignore" alt="ignore unused">IgnoreUnused</a>
  * <a href="#setlog"
       alt="set global zerolog to file">SetGlobalZerologToFile</a>
* <a href="#dependencies" alt="dependencies">Dependencies</a>
* <a href="#incompat" alt="incompatibilities">Incompatibilities</a>
* <a href="#bugs" alt="bugs and limitations">Bugs and Limitations</a>
* <a href="#thanks" alt="acknowledgements">Acknowledgements</a>
* <a href="#author" alt="author">Author</a>
* <a href="#copyright"
     alt="copyright and license">Copyright and License</a>
</details>

### <a name="name">Name</a>

veil - minor Go additions

The veil package is thus named because it is a thin
[veil][veilimg]
over existing Go code.

The veil package contains minor[^1] enhancements to the Go
[standard library][stdlib] and other wonderful open source packages.

### <a name="version">Version</a>

This documentation is for veil version **v1.0.0**

### <a name="synopsis">Synopsis</a>

```go
package main

import (
    "log"

    "github.com/kjmjonline/veil"
    "github.com/rs/zerolog"
)

func main() {
    // Let's use FilePathInCwd to get the fully qualified file name
    // of the file to write logs to in the current working directory
    var filepath string
    var err error
    if filepath, err = veil.FilePathInCwd("my-project.log"); err != nil {
        log.Fatal(err)
    }

    // Now we will initialize the zerlogo logging to this file
    if err = veil.SetGlobalZerologToFile(filepath, zerolog.DebugLevel); err != nil {
        log.Fatal(err)
    }

    log.Print(
        "At this point we can log using the global zerolog logger",
    )

    ...

    // We can also use veil.CaptureOutput to capture and test any output
}
```

### <a name="description">Description</a>

This package contains some trivial functions that are needed by kjmjonline.

It contains a zerolog helper function [SetGlobalZerologToFile][setlog]
to setup the zerolog logging facility.

Fully qualified file paths for files in the current working directory can
be determined with the [FilePathInCwd][filepath] function.

The [CaptureOutput][capture] function can be used to get any ouput from a
function. This can be used, for instance, to verify that the function is
working correctly.

The Go compiler does an excellent job of catching any unused symbols in
your code. You can, however, use the [IgnoreUnused][ignore] function to
silence these errors, for instance during development.

### <a name="installation">Installation</a>

```bash
go get -u -v github.com/kjmjonline/veil@latest
```

This will install veil and its dependencies to your `go/pkg`
directory.

### <a id="funcs">Public Functions</a>

#### <a id="capture">CaptureOutput</a>

Captures, and returns, the merged `stdout` and `stderr` output of a
function.

```go
package main

import (
    "fmt"

    "github.com/kjmjonline/veil"
)

func sayHello() {
    fmt.Print("Hello, stranger!")
}

func main() {
    greeting, err := veil.CaptureOutput(sayHello)
    // `greeting` will contain "Hello, stranger!" here

    if err == nil {
        fmt.Println(greeting)
    }
}
```

#### <a name="filepath">FilePathInCwd</a>

Returns the full path to the given _fileName_ in the current work directory
(i.e., equivalent to `cwd` in &#42;nix systems).

The current directory is the directory that this program was started from.
This may be different from the directory that the executable is in.

```go
package main

import (
    "fmt"

    "github.com/kjmjonline/veil"
)

func main() {
    fileName := "blort.txt"

    // i.e., the path plus the filename
    var filePath string

    var err error
    if filePath, err = veil.FilePathInCwd(fileName); err != nil {
        panic(err)
    }
    fmt.Print(filePath)

    // If the current directory (the directory this program was
    // started from) is "/tmp" then the`filePath` returned would be:
    //   "/tmp/blort.txt".
}
```

#### <a name="ignore">IgnoreUnused</a>

Silences Go errors caused when code contains any unused constants,
variables, and/or functions.

To silence these errors pass the name of each unused identifier to this
function.

```go
package main

import (
    // We can ignore unused package import errors by preceding the
    // unused package with un underscore character:
    _ "fmt"

    "github.com/kjmjonline/veil"
)

const unusedConst = 42

type unusedType struct{}

var trulyUnusedVar string

func someUnusedFunc() string {
    return "this function is not used"
}

func main() {
    # We can ignore unused types like this:
    var _ unusedType

    # And here we ignore any unused constants, variables, and functions
    veil.IgnoreUnused(
        unusedConst,
        trulyUnusedVar,
        someUnusedFunc,
        )
}
```

#### <a name="setlog">SetGlobalZerologToFile</a>

This function sets up the global zerolog logger.

The log file is created, or is appended to if it already exists.

Logging is set up to create log entries with the current time timestamp,
and file name and line numbers where the log entries were created. The
timestamps use the [RFC 3339 Nano][rfc3339] time format, which has
sub-second precision.

The log messages are colored. If you do not want colored logs you could
create a utility to remove the color escape character sequences.

```go
package main

import (
    sl "log"

    "github.com/kjmjonline/veil"
    "github.com/rs/zerolog/log"
    "github.com/rs/zerolog/pkgerrors"
)

func main() {
    err := veil.SetGlobalZerologToFile("mylog", zerolog.DebugLevel)
    if err != nil {
        sl.Fatal(err)
    }

    // log.Print() is similar to log.Debug().Msg()
    log.Print("global zerolog has been created")

    err := errors.New("Whoa!")
    wrapped := errors.WithStack(err)
    log.Error().Stack().Err(wrapped).Msg("An error with a stacktrace")
}
```

### <a name="dependencies">Dependencies</a>

veil uses some packages that are not part of the Go standard library.
These libraries are _automatically_ installed when veil is installed.

They are:
* github.com/rs/zerolog

What!? That's it!

### <a name="bugs">Bugs and Limitations</a>

This package has no known bugs or limitations.

Please report any incompatibilities to <justinhanekom7@outlook.com>.

### <a name="thanks">Acknowledgements</a>

The CaptureOutput, FilePathInCwd, and IgnoreUnused functions were
adopted from other sources that have been lost to the mists of time.

The functions have been changed slightly.

Please report any missing acknowledgements to <justinhanekom7@outlook.com>.

### <a name="author">Author</a>

Justin Hanekom

### <a name="copyright">Copyright and License</a>

Copyright (c) 2024 Justin Hanekom

This file is part of veil - minor enhancements to Go libraries.

veil is free software: you can redistribute it and/or modify it
under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

veil is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with veil. If not, see <https://www.gnu.org/licenses/>.

[veilimg]:  https://unsplash.com/photos/a-woman-wearing-a-veil-and-earrings-6UeL4IWhKfY "image of a woman wearing a veil"

[stdlib]:   http://pkg.go.dev/std/ "Go standard library"
[capture]:  #capture  "CaptureOutput function"
[filepath]: #filepath "FilePathInCwd function"
[ignore]:   #ignore   "IgnoreUnused function"
[setlog]:   #setlog   "SetGlobalZerologToFile function"
[rfc3339]:  https://www.ietf.org/archive/id/draft-ietf-sedate-datetime-extended-09.html "RFC 3339 timestamp format"

[^1]:       _miniscule_, really!
