sirdsc
======

sirdsc is a simple that converts height-maps into single image random dot stereograms ([SIRDS][sirds]).

Prerequisites
-------------

 * [Go Programming Language][golang]

Installation
------------

Run the following commands:

> make
> make install

The default prefix is '/usr'. If you want to install sirdsc with a different prefix, simply run the following:

> make PREFIX="your desired prefix" install

Usage
-----

> sirdsc [options] <src> <dest>

Where <src> is an existing height-map file in a supported format, and <dest> is the file to which to write the generated file.

### Options ###

sirdsc accepts the following options:

 * -partsize=<int> (Default: 100)
   > The size of the individual parts of the generated SIRDS. The generated image will be this many pixels wider than the height-map.
 * -depth=<int> (Default: 10)
   > The maximum depth of the image pixels.

The following options only apply if the destination file is a JPEG:

 * -jpeg:quality (Default: 95)
   > The quality of the output JPEG file as a percentage.

Authors
-------

 * DeedleFake

[sirds]: http://www.wikipedia.com/wiki/SIRDS
[golang]: http://www.golang.org
