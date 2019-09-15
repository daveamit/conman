[![Build Status](https://cloud.drone.io/api/badges/daveamit/conman/status.svg)](https://cloud.drone.io/daveamit/conman)
[![Go Report Card](https://goreportcard.com/badge/github.com/daveamit/conman)](https://goreportcard.com/report/github.com/daveamit/conman)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/daveamit/conman/blob/master/LICENSE)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-89%25-brightgreen.svg?longCache=true&style=flat)</a>

# conman
CONfiguration MANager

# Preface
Configurations are a integral part of any software, things get more complicated when the software we are talking about has a `microservice` based architecture.

With microservices, each service has to have its own configuration, it has to be isolated and secure in a sense that other services running on the same platform `must` not be able to `read` it. Also, the `service` / `software` must only be able to `read` it and not `modify` and so on.

Configuration `must` not be treaded as another piece of data because it is not. Also, `configuration` are `owned` by different set of user than the `data`, for example, suppose `identity` services needs a configuration called `password-hash-key`. The service will use this to hash passwords and hence is quite sensitive information and has to be `owned` by the a different `identity` source. It would not be wise to store it in same database or pass it as cli arguments (or environment variables), also a service in charge of sending email should not be able to `read` such information.

# Concept
> (tool) -> (authentication/authorization) -> (storage) -> (push) -> (service)

That's it, in a nutshell. It translates to

Hence,
> (conman) -> (storage) -> (conman-driver) (service)

# What's in the box
* `REST`ful server (Feature set exposed as `REST`ful api)
* The `cli` (fully featured, runs on `linux`, `osx` and `windows`)
* driver (specially designed etcd wrapper with intuitive api)
  * `golang`
  * (will support more languages if need be)


# Feature set
## Authentication / Authorization
* Initialize and reset the `storage`
* Login / Logout
## Manage settings
* Define configuration setting schema
* Set / Delete / Enable / Disable configuration settings
* Exporting / Importing configuration bundle
## Manage users
* Add / Delete / Enable / Disable users
* Permission descriptor for users