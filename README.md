# Hello

Hello is an example project that demonstrates the development flow and common setup for projects released by Gravitational Inc.

Functionally it provides CLI, service and a library with a wrapper to output `Hello, gravity!` when called.

Current version: `0.0.1-alpha.1`

Guides on writing software at Gravitational:

* [Developers guide](https://github.com/gravitational/wiki/blob/master/docs/dev/org/devprocess.md)

# On writing documentation:

Lets aim for consistency here and every Gravitational project must:

- Execute `--help` when launched without parameters.
- Have `README.md` in its root directory with the User Manual containing the following sections:
  - Intro and Purpose
  - Building or Installing
  - Usage
- Have `docs` directory with additional docs, referred from `README.md`, for example:
  - `design.md` - developer design
  - `contribute.md` - instructions for other contributors 
