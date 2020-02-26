# Identicon Generator

## Introduction

This is a small example of a simple identicon generator written in golang, the blog post can be found here: https://www.bartfokker.nl/posts/identicon/

## go get to get the package

```sh
go get github.com/lercher/identicon
```

## This Fork

This fork introduces more control over the generated png image
and adds more contrast to the color chosen.

### Breaking Changes

The `Name` field is removed, b/c it is only converted to a string from the generator byte slice,
which needs not be a string up-front.

The `WriteImage` method was removed and replaced by `WritePNGImage` with more parameters.

## Usage 

```go
import "github.com/lercher/identicon"

i := identicon.Generate([]byte("Simpson"))
_ = i.WritePNGImage(w, pw, identicon.LightBackground(true)) 
```
