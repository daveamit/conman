# conman
CONfiguration MANager backed by ETCD

# Preface
Configurations are a integral part of any software, things get more complicated when the software we are talking about has a `microservice` based architecture.

With microservices, each service has to have its own configuration, it has to be isolated and secure in a sense that other services running on the same platform `must` not be able to `read` it. Also, the `service` / `software` must only be able to `read` it and not `modify` and so on.

Configuration `must` not be treaded as another piece of data because it is not. Also, `configuration` are `owned` by different set of user than the `data`, for example, suppose `identity` services needs a configuration called `password-hash-key`. The service will use this to hash passwords and hence is quite sensitive information and has to be `owned` by the a different `identity` source. It would not be wise to store it in same database or pass it as cli arguments (or environment variables), also a service in charge of sending email should not be able to `read` such information.

# Concept
> (tool) -> (authentication/authorization) -> (storage) -> (service)

That's it, in a nutshell. It translates to

> (etcdctl) -> (etcd) -> (service)

Don't get me wrong, `etcdctl` is awesome tool, but it lacks certain like help me validate my configuration schema, it is not that driven from file, it does not provide an intuitive way to map users to configuration from a single call.

Hence,
> (conman) -> (etcd) -> (service)

# What's in the box
* `REST`ful server (Feature set exposed as `REST`ful api)
* The `cli` (fully featured, runs on `linux`, `osx` and `windows`)
* driver (specially designed etcd wrapper with intuitive api)
  * `golang`
  * (will support more languages if need be)


# Feature set
## Authentication / Authorization
* Initialize and reset the `etcd`
* Login / Logout
## Manage settings
* Define configuration setting schema
* Set / Delete / Enable / Disable configuration settings
* Exporting / Importing configuration bundle
## Manage users
* Add / Delete / Enable / Disable users
* Permission descriptor for users