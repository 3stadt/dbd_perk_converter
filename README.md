# Project Title

## Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Usage](#usage)
- [Contributing](../CONTRIBUTING.md)

## About <a name = "about"></a>

This is a converter tool that copies perk icons from dead by daylight to a new folder while converting the name and file structure so it can be used with [3stadt/dbd_perk_slot_machine/](https://github.com/3stadt/dbd_perk_slot_machine/).

## Getting Started <a name = "getting_started"></a>

After cloning the repo open your terminal, naivage to the folder containing the file `main.go` and execute `go run main.go --help`.

Alternatively you can build and run the file as any other go file

### Prerequisites

[go 1.17](https://golang.org/dl/) or later.

## Usage <a name = "usage"></a>

The tool takes three input parameters:

- `-i` Specifies the input directory which holds the source perk images.
- `-o` Specifies the output directory where the renamed files should be copied to.
- `-d` Specify the config file holding the target perk names.
    - It is recommended to use the included `*.txt` files.