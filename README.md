# Jorge

A tool that helps to manage multiple configuration file versions of an application under development

![lightblue gopher dressed as medieval librarian](./docs/logo_375x368.png)

## Description

### Problem 

I usually am in a situation where I have to deal with multiple versions of configuration files with secrets (e.g. `.env` or `appsettings.development.json`) when I develop or debug an application. 
To solve this problem I usually use comments when the configuration file format supports it or I hold multiple copies in another place.

### Solution

This tool enables the developer to store multiple versions of a configuration file. It is designed to be used during development (and debuging).
The different versions of the configuration file are stored under a `.jorge` directory at the project root. Jorge replaces the configuration file that the project is using with one that is stored under `.jorge` when you want to use it.

### Typical use case

1. You already have a project with a configuration file. Let's say `.env`
2. Initialize a jorge project `jorge init --config .env` 
3. Now you want to test your project with different, but temporary, configuration (e.g. to use another backend server or email service)
4. You create a different jorge environment `jorge use -n staging-server` or `jorge use -n staging-services` etc
5. You do you and your testing was finished, so you want to return to the development settings
6. You change back to original environment `jorge use default` (note that the env `default` is created with the `jorge init` command)

### Why Jorge?

Because of [a badass jorge](https://www.litcharts.com/lit/the-name-of-the-rose/characters/jorge-of-burgos) who knew how to protect secrets

## Installation

### Build from source 

- `go get`
- `make build`
- `make install` - requires `sudo`
- Use it as a cli tool `jorge --version`

### Build using docker

- `docker build -t jorge .`
- Use it as container `docker run --rm -v "$PWD":/root/projectRoot jorge jorge --version`



## Usage

Create a project

`jorge init`

Fill the path your project configuration file

See the available environments

`jorge ls`

Create and use an environment

`jorge use -n test01`

Change environment

`jorge use default`

Commit your changes to the current env

`jorge commit`

Restore the changes from a stored version to the current env

`jorge restore`


## Reference

```
Manages different versions of a configuration file

Usage:
  jorge [command]

Available Commands:
  commit      Stores the current config file
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initializes a jorge environment
  ls          List the available environments
  restore     Restores the current configuration file with the copy that is saved in the .jorge dir
  use         Selects or creates an environment

Flags:
  -d, --debug     Prints debug messages
  -h, --help      help for jorge
  -v, --version   version for jorge
```
