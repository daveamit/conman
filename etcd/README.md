# conman
CONfiguration MANager implementation backed by ETCD


# Concept
> (tool) -> (authentication/authorization) -> (storage) -> (push) -> (service)

That's it, in a nutshell. It translates to

> (etcdctl) -> (etcd) -> (service)

Don't get me wrong, `etcdctl` is an awesome tool, but it lacks certain features like help me validate my configuration schema, it is not that file driven, it does not provide an intuitive way to map users to configuration from a single (or rather simple) call.

Hence,
> (conman) -> (etcd) -> (conman-driver) (service)
