# goshell

A few personal tools that are too small to deserve a dedicated repo, but still quite useful to me, and hopefully you, stranger.

If you'll decide to install some - make sure to have go and git installed.


## aqua

```
Usage of aqua:
  -a string
    	file a
  -b string
    	file b
```

Install:

```
go install github.com/MahouShoujoMivutilde/shellbin/cmd/aqua@latest
```


## hum

```
hum

  Is a tool for humanizing various things (now supports time)

  echo thing | hum (time) [-t]

Usage of hum:
  -t string
    	time format, fill with Mon Jan 2 15:04:05 MST 2006, see https://golang.org/src/time/format.go (default "2006-01-02 15:04:05.999999999 -0700")
```

Install:

```
go install github.com/MahouShoujoMivutilde/shellbin/cmd/hum@latest
```


## istextmt

```
istextmt

  Checks if file is a text file, and
    if it is
      prints filepath
    else
      doesn't print anything

  Designed to be used as filter for fd,
  it is also much faster than
    file --mime-type -b file.txt + case text/*...
  ...shenanigans.

Usage of istextmt:

Examples:
  check file is a text file
    echo "file.txt" | istextmt
    if it is, it will print it's name, otherwise - skip

  find only only text files with fd
    fd -t f | istextmt
```

Install:

```
go install github.com/MahouShoujoMivutilde/shellbin/cmd/istextmt@latest
```

On 17k files

* `fd ... | istextmt` takes 200ms

* `fd ... | istext` takes 2s

* `fd -x istext {}` takes 14s (legacy)

## sortlf

```
sortlf <diretory>
	like `ls`, but with the sorting algo from `lf`
	respects `lf_sortby`, `lf_reverse` and `lf_hidden` env. variables
	but it doesn't know about your `setlocal`'s and filter

DEPRECATED
	use `lf -remote "query $id files"` on lf r35+
```

Install:

```
go install github.com/MahouShoujoMivutilde/shellbin/cmd/sortlf@latest
```


## urlesc

```
urlesc

  Is a tool for escaping file path to make it safe for urls

Usage of urlesc:
  -u	unescape uri instead

Examples:
  escape path
    echo 'some/path/внезапно!@@"/dir' | urlesc
```

Install:

```
go install github.com/MahouShoujoMivutilde/shellbin/cmd/urlesc@latest
```


## zlgo

```
zlgo

  Is a tool for generating zalgo text

Usage of zlgo:

Examples:
  zalgofy stdio
    echo 'some text....' | zlgo
```

Install:

```
go install github.com/MahouShoujoMivutilde/shellbin/cmd/zlgo@latest
```

