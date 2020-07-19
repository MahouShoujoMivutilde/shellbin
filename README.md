# goshell

A few personal tools that are too small to deserve full blown personal repo, but still quite useful to me, and hopefully you, stranger.

If you'll decide to install some - make sure to have golang and git installed.



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

```
lsidups

  Is a tool for finding image dupicates (or just similar images).
  Outputs images grouped by similarity (one filepath per line) to stdio
  so you can process them as you please.

Usage of lsidups:
  -e value
    	image extensions (with dots) to look for (default .jpg,.jpeg,.png,.gif)
  -i string
    	directory to search (recursively) for duplicates, when set to - can take list of images
    	to compare from stdio (one filepath per line, like from find & fd...) (default ".")
  -v	show time it took to complete key parts of the search

Examples:
  find duplicates in ~/Pictures
    lsidups -i ~/Pictures > dups.txt

  or compare just selected images
    fd 'mashu' -e png --changed-within 2weeks ~/Pictures > yourlist.txt
    lsidups -i - < yourlist.txt > dups.txt

  then process them in any image viewer that can read stdio (sxiv, imv...)
    sxiv -io < dups.txt
```

Install:

```
go get github.com/MahouShoujoMivutilde/shellbin/cmd/lsidups
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

