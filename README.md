linkwarp
========

`linkwarp` is a symlink farmer to help you manage your `$PATH`.

In short: it will...

1. look for applications and their executables;
2. when it finds them, put a symlink to each executable in the output directory;
3. that's it!

The idea is when `linkwarp` is done, you just add that one output directory to your `$PATH` environment variable,
and now all your programs are ready to be launched.

`linkwarp` does a simple search pattern guided by zero-config heuristics.  It's fast and it does what you mean.

That's it!


Installation
------------

```
go get github.com/warptools/linkwarp/cmd/linkwarp@latest
```


Usage
-----

```
linkwarp
```

or, if you're fancy, and want to specify search path and output path:

```
linkwarp ./searchpath ./outputpath
```


Configuration
-------------

There isn't any, yet.  (PRs welcome.)

`linkwarp` aims to "DTRT" by default by following the simplest conventions that work.  Usually, this is enough.

The default heuristics are this:

- within the search path: look at any directory up to two dirs deep to consider where applications might be.
- if a directory contains a child dir called 'bin/', it's an application directory.
- find every executable file inside `bin/`: symlink those into the output dir.
- _in the case of conflict_: do the simplest thing that works: compare the application directory name; the one that looks "bigger" using a human-friendly sort wins.
  (Meaning: if the application directory names are "foo-v1" and "foo-v2", then "foo-v2" wins!  This usually does what you want, with no special intervention.)



Relationships
-------------

### Warpforge

Linkwarp was made under the same inspirations as power the [Warpforge](http://warpforge.io/) project,
and specifically, implements the conventions described in the [WarpSys Execution Path Management documentation](https://warpforge.notion.site/Execution-Path-Management-e9e1844bcfc44d528eed09107d2ebadc).
However, it's not directly entangled with those projects in any way -- you can use `linkwarp` anywhere.
