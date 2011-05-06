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

> make PREFIX=&lt;your desired prefix&gt; install

Be sure to replace &lt;your desired prefix&gt; with whatever your desired prefix is.

Usage
-----

> sirdsc [options] &lt;src&gt; &lt;dest&gt;

Where &lt;src&gt; is an existing height-map file in a supported format, and &lt;dest&gt; is the file to which to write the generated file.

### Options ###

sirdsc accepts the following options:

 * -partsize=&lt;int&gt; (Default: 100)
   > The size of the individual parts of the generated SIRDS. The generated image will be this many pixels wider than the height-map.
 * -depth=&lt;int&gt; (Default: 10)
   > The maximum depth of the image pixels.

The following options only apply if the destination file is a JPEG:

 * -jpeg:quality (Default: 95)
   > The quality of the output JPEG file as a percentage.

Authors
-------

 * [DeedleFake](/DeedleFake)

[sirds]: http://www.wikipedia.com/wiki/SIRDS
[golang]: http://www.golang.org
