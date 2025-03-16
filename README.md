
# CLI Web Scraper

This project was made on purpose of learning [Go](https://go.dev/). It's probably really unoptimized, any criticizm is welcome.


## Installation & usage

  - Windows
```bash
  go build
  go install
  cli_web_scraper [url]
```

  - Linux
```bash
  go build
  ./cli_web_scraper [url]
```

## Packages used

 -  [errors - implements functions to manipulate errors](https://pkg.go.dev/errors)
 -	[fmt - implements formatted I/O with functions analogous to C's printf and scanf](https://pkg.go.dev/fmt)
 -	[os - provides a platform-independent interface to operating system functionality](https://pkg.go.dev/os)
 -	[io - provides basic interfaces to I/O primitives](https://pkg.go.dev/io)
 -  [log - implements a simple logging package](https://pkg.go.dev/log)
 -  [slices - defines various functions useful with slices of any type](https://pkg.go.dev/slices)
 -  [strings - implements simple functions to manipulate UTF-8 encoded strings](https://pkg.go.dev/strings)
 -  [time - provides functionality for measuring and displaying time](https://pkg.go.dev/time)
 -  [net - provides a portable interface for network I/O](https://pkg.go.dev/net)
 -  [net/http - provides HTTP client and server implementations](https://pkg.go.dev/net/http)
 -  [golang.org/x/net/html - implements an HTML5-compliant tokenizer and parser](https://pkg.go.dev/golang.org/x/net/html)
 -  [golang.org/x/net/html/atom - provides integer codes (also known as atoms) for a fixed set of frequently occurring HTML strings](https://pkg.go.dev/golang.org/x/net/html/atom)
 -	[github.com/spf13/cobra - framework for CLI apps](https://cobra.dev/)
