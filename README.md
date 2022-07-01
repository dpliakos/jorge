# Jorge - Secret manager

A tool that helps to manage multiple configuration file versions during your development process

## Problem

I usually am in a situation where I have to deal with multiple versions of configuration files with secrets  when I develop or debug an application. To solve this problem I use comments when the configuration file format supports it or I hold multiple copies in another place.
This tool automates the second approach and stores multiple versions of the configuration files and allows you to quickly access them


## Installation

### Build from source 

- `go get .`
- `make build`



## Usage

### Create a project

`jorge init`

Fill the path your project configuration file

### See the available environments

`jorge ls`

### Create and use an environment

`jorge use -n test01

### Change environment

`jorge use default`

### Commit your changes to the current env

`jorge commit`

### Restore the changes from a stored version to the current env

`jorge restore`
