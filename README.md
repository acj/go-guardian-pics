Summary
=======
A parser, written in the Go language, for the 24 Hours in Pictures RSS feed published by the Guardian. It also serves as a basic example for how to use Go's `xml` and `http` packages.

Usage
=====
To use the default feed URL embedded in the code:
./fetch_gallery

To parse a feed on disk:
./fetch_gallery file.xml
