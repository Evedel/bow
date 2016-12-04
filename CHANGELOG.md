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
