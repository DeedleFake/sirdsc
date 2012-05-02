sirdsc
======

sirdsc is a simple that converts height-maps into single image random dot stereograms ([SIRDS][sirds]).

Prerequisites
-------------

 * [Go Programming Language][golang]

Installation
------------

Run the following command:

    go get github.com/DeedleFake/sirdsc

This will install the command 'sirdsc' into your GOPATH. For more information, run:

    go help gopath

Usage
-----

> sirdsc [options] &lt;src&gt; &lt;dest&gt;

Where &lt;src&gt; is an existing height-map file in a supported format, and &lt;dest&gt; is the file to which to write the generated file.

### Options ###

sirdsc accepts the following options:

<dl>
    <dt>-partsize=&lt;int&gt; (Default: 100)</dt>
    <dd>The size of the individual parts of the generated SIRDS. The generated image will be this many pixels wider than the height-map. If -partsize is set to 0 and used in conjuction with -pat it will be automatically detected from the width of the specified custom pattern.</dd>

    <dt>-depth=&lt;int&gt; (Default: 40)</dt>
    <dd>The maximum depth. A solid white pixel in the height-map results in this depth.</dd>

    <dt>-flat (Default: false)</dt>
    <dd>If specified, treat any non-black pixels as having the maximum depth.</dd>

    <dt>-pat=&lt;string&gt; (Default: "")</dt>
    <dd>If not equal to "", use the file specified as the pattern instead of generating a random one. To have -partsize automatically detected based on the width of the specified file, use -partsize=0.</dd>
</dl>

The following options only apply if the destination file is a JPEG:

<dl>
    <dt>-jpeg:quality=&lt;int&gt; (Default: 95)</dt>
    <dd>The quality of the output JPEG file as a percentage.</dd>
</dl>

Authors
-------

 * [DeedleFake](/DeedleFake)

[sirds]: http://www.wikipedia.com/wiki/SIRDS
[golang]: http://www.golang.org

<!--
    vim:ts=4 sw=4 et
-->
