# Identicon Generator

## Introduction

This is a small example of a simple identicon generator written in golang, the blog post can be found here:
[https://www.bartfokker.nl/posts/identicon/](https://www.bartfokker.nl/posts/identicon/)

## go get to get the package

```sh
go get github.com/lercher/identicon
```

## This Fork

This fork introduces more control over the png image generated,
adds more contrast to the chosen color and publishes its list of set pixel coordinates.

### Breaking Changes

The `Name` field is removed, b/c it is only converted to a string from the generator byte slice,
which needs not be a string up-front.

The `WriteImage` method was removed and replaced by `WritePNGImage` with more parameters.

### Added API

An identicon now publishes a slice of (X,Y) coordinates of set pixels as `Pixels`

## Usage

```go
import "github.com/lercher/identicon"

i := identicon.Generate([]byte("Simpson"))
// 50px per identicon pixel, i.e. 250x250 image written
// using dark colors in the byte range 0-127:
_ = i.WritePNGImage(w, 50, identicon.LightBackground(true))
log.Println(i.Pixels)
```
