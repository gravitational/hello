# Hello

This is a tiny web app skeleton which shows the basic project structure and development workflow we use.
Features are few:

* CLI 
* web service 
* library shared by the two

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
