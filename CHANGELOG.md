Build 27.1
==
(2020.07.03)
- CSS fixes [issue #16](https://github.com/Evedel/bow/issues/16)

Build 26.1
==
(2018.02.06)
- Added alpine based image, as requested in [pr #1](https://github.com/Evedel/bow/pull/12), changed labelling system.

Build 26
==
(2017.11.19)
- Fix: free format for service name in bearer token [issue #10](https://github.com/Evedel/bow/issues/10)

Build 25
==
(2017.09.01)
- Add: namespace level [issue #9](https://github.com/Evedel/bow/issues/9)
- Add: filters on graph page [issue #9](https://github.com/Evedel/bow/issues/9)
- Fix: list not only [first 100](https://github.com/docker/distribution/blob/b6e0cfbdaa1ddc3a17c95142c7bf6e42c5567370/registry/handlers/catalog.go#L16) but all of the images [issue #9](https://github.com/Evedel/bow/issues/9)
- Add: timewatch for handlers and daemons
- Add: script to fill test repo
- Fix: loops brake parentgraph when there were copies in repo

Build 24
==
(2017.02.18)
- Fix: Wrong permissions when creating /var/lib/bow/ as non-root user [issue #8](https://github.com/Evedel/bow/issues/8)
- Rename smurf names in checker, underscore to camel in db/schema/recursive

Build 23
==
(2017.02.18)
- Fix: Can't authenticate when using a password with percent symbol [issue #7](https://github.com/Evedel/bow/issues/7)
- Fix: Actualised Tests

Build 22
==
(2017.02.04)
- Fix: panic in parents search for committed images [issue #5](https://github.com/Evedel/bow/issues/5)
- Autoscroll to chosen element

Build 21
==
(2016.12.23)
- Added ugly button to force bow to update all info
- DB refactoring finished
  - db.tag : deleted, created alias, added precreation of non-existed buckets for new tag
  - db.catalog : the same as for db.tag
- QURL refactoring finished
  - Fixed isues where bow was not able to make bearer auth reqest for HEAD/DELET/headers requests
  - Nice & compact method capable to serv all reqest
- Fix: error when registry storage was changed externally
- Fix: random order of repos/names/tags
- Fix: error for non-existing catalogs just after creation
- Fix: db.upgrade for the most old versions

Build 20
==
(2016.12.04)
- Added possibility to login through bearer token
- Added explicit secure/insecure checkbox in repo config
- DB schema redesign
  - add '_info' bucket to store confings, move configs into it and create upgrade way
  - from 'repo.go' to general 'db' middleware
  - add db versions, check it in each restart, do upgrade if needed
- Splitted 'db' => 'put', 'get', 'delete', 'schema', 'upgrade' and 'alias'
- Splitted main.go => package 'handler' + package 'main'
- Added full test env (auth server, basic auth registry and token based auth) + gitignore
- Fixed error in blob size conversion

Build 19
==
(2016.11.01)  
- API compatibility covered by tests  
- Fixed time conversion bug for Manifest Daemon
- Upgrade moved from db to utils
- Size conversion moved from db to utils
- MakeQueryToRepo now operate with http status codes
