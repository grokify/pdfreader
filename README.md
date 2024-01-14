# pdfreader

## Introduction

The pdfreader library for Go is a library to read contents of PDF files.

Basically it turns out that this will result in a PDF to SVG converter at the first stage.

## Details

PDF files are basically something that is usually at the end of some workflow and that is intended to conserve informations in a way that allows the informations to be accessed as they where intended (e.g. in terms of typographical layout). These informations need to be fetched. The project here tries to make this possible with a library for Go.

If you are not a Go programmer, just move away or play with an example application.

Currently everything is at it's premature state and there is no production-ready library to be expected. Well, the things work usually fine for many tasks.

If you are willing to make experiments, just checkout at http://code.google.com/p/pdfreader/source/checkout

## Basic design principles

* Using this library with a malformed PDF might crash the program. This is intentional.
* Keep things simple - no reason to produce billions of lines of code.
* Make the crash to be late. As late as possible. If there is something really wrong it will crash earlier or later. Why using a "safe" programming language if not using it and adding useless tests for validity of input?
* Avoid endless recursions. There are many places where this could occur in PDF-files. A fixing of issue 226 with golang would help, but the gurus of Google did decide to do different. So be prepared to have no real fun with the implementation language. See Philosophy-page.

## Example

This shows an SVG displayed in Inkscape that was converted from a PDF:

[`image-convert.png`](image-convert.png")

## Philosophy

It's as of today (end of 2009) sometimes a little bit strange to write larger programs with Go. The language and the libraries are still in development and they may change at every point of a day. The thing is more a moving target than something you can build cathedrals on. Does it sound bad enough for you? It might become even worse if you insist on the chance to get something changed in Go that fits your predeclared mind about what a programming language is or how it has to be. The gods of Go will follow their own ways.

For my mind they crippled recursion - especially recursion of local/unnamed functions. I guess they recognized the problem but they did not manage to jump to the attitude to add a new keyword that would be needed to support recursion for unnamed functions (How else to call a function that has no name?).

Other people do have other issues. So some people are not willing to type ";" (semicolon) at end of statements. Well, the Go-gods decided to remove the need for semicolons more or less. One could say this is minimalism, but even a lot earlier they decided to not have something like a = b < 0 ? "less zero" : "greater/equal zero" This is C-like and the same statement now needs several lines of code to be done. Others than me would cry about the infix- and postfix-operators that are now statements, but these are really not that useful in my mind.

So well, the things are interesting.

I do use * code generators written in Perl * an autoimporter for needed packages * a dependency-finder for the Makefiles to make work with the sources more suitable. I hope that at some day I do not need the code generators in some special cases (as I use they now, they will be useful forever ;) ) and that I do not need the autoimporter anymore. The autoimporter would be better located inside the compiler. How these two things are organized will change what is needed for the dependency-finder. Until a really good solution is out there, I'd say the things should stay the same as they are now (Dec. 22 2009 update).

Please comment if you have some ideas. I'll update this page as far as I have time.