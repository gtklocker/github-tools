ghnuke
======

This tool unsubscribes you from all repos belonging to a specific organization.
Usually when joining a new organization you are automatically subscribed to
updates from all its repo. Github does not provide an easy way to unsubscribe
from all the repos and worst of all the subscriptions for the organization's
public repos persist even if you leave the organization.

[Create a personal access token on Github](https://github.com/settings/tokens)

Set environment variable GHACCESSTOKEN to your personal access token and export it.

```
go run ./cmd/ghnuke -org yourorggoeshere -unwatch # add -unstar to also remove stars
```

ghpriv
======

Go through your least favorite public repos and allows you to:
- switch to private or
- delete them or
- keep them

```
go run ./cmd/ghpriv
```

ghunfollow
==========

Unfollows everyone.

```
go run ./cmd/ghunfollow
```
