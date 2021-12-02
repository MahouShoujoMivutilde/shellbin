# goshell

A few personal tools that are too small to deserve full blown personal repo, but still quite useful to me, and hopefully you, stranger.

If you'll decide to install some - make sure to have golang and git installed.


## sortlf

```
sortlf <diretory>
	like `ls`, but with the sorting algo from `lf`
	respects `lf_sortby` and `lf_reverse` env. variables
```

Install:

```
go get github.com/MahouShoujoMivutilde/shellbin/cmd/sortlf
```


## hum

```
hum

  Is a tool for humanizing various thigs (now supports time)

  echo thing | hum -args...

Usage of hum:
  -t string
    	time format, fill with Mon Jan 2 15:04:05 MST 2006, see https://golang.org/src/time/format.go (default "2006-01-02 15:04:05.999999999 -0700")
```

Install:

```
go get github.com/MahouShoujoMivutilde/shellbin/cmd/hum
```



## fitrectg

```
fitrectg

  Fits some rectangle (meant for images) into given rectangle while
  preserving aspect ratio. Outputs new width and height as WxH.

Usage of fitrectg:
  -fh float
    	rectangle height
  -fw float
    	rectangle width
  -h float
    	current image height
  -w float
    	current image width

Examples:
  calculate new dimensions of image with width = 3600 and height = 2404
  to fit into 456x490 rectangle
    fitrectg -w 3600 -h 2404 -fw 456 -fh 490
```

Install:

```
go get github.com/MahouShoujoMivutilde/shellbin/cmd/fitrectg
```


## istext

```
istext

  Checks if file is a text file, and
    if it is
      prints filepath and exits with 0,
    else
      exits with 1

  Designed to be used as filter for fd,
  it is also much faster than
    file --mime-type -b file.txt + case text/*...
  ...shenanigans.

Usage of istext:

Examples:
  check file is a text file
    istext file.txt && echo 'this is text file' || echo 'this is not text'

  find only only text files with fd
    fd -t f -x istext {}
```

Install:

```
go get github.com/MahouShoujoMivutilde/shellbin/cmd/istext
```


## lsidups

lsidups is a barebone tool for finding image duplicates (or just similar images) from your terminal. Prints to stdout list of images grouped by similarity to allow later processing with other tools, like sxiv.

[Moved here](https://github.com/MahouShoujoMivutilde/lsidups).


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
go get github.com/MahouShoujoMivutilde/shellbin/cmd/urlesc
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
go get github.com/MahouShoujoMivutilde/shellbin/cmd/zlgo
```

