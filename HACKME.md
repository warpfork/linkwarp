HACKME
======

### Y u no use fs package?

Well, we did, in some places.  The walk functions, in particular, remain handy.

However, were we able to use the `fs` package to near-total exclusion of the `os` package?  No.

As much as I would love to, especially to be able to have abilities for test mocking... we just can't.
The `fs` package leaves a lot to be desired beyond the most basic operations.

In particular, ability to either describe nor read nor write symlinks is pretty much entirely absent...
and give that's the bulk of what we do in this program, uh, yeah, that's a bit of a limit, isn't it.
