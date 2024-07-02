## rog-codegen

`rog-codegen` turns [OpenAPI 3.0 specs](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md) into Ruby on Rails code, cutting down on boilerplate so you can focus on your business logic and adding real value to your organization.

## Installation

`rog-codegen` requires [Go](https://go.dev/dl/) >= `1.20`.

You can install `rog-codegen` as a binary:

```bash
$ go install github.com/trungvitlonx/rog-codegen@latest
```

## Usage

`rog-codegen` is largely configured using a YAML configuration file, to simplify the number of flags that users need to remember.

This will create a YAML configuration file (`.rog.yaml`) and a sample of OpenAPI 3.0 specification file (`openapi.yaml`):

```bash
$ rog-codegen init
```

Then, to generate the code:

```bash
$ rog-codegen generate
```
By default, you will find the generated code at `./gen`.
