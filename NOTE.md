__BS_LOG_SILENT__ -- level of output logging  
nothing - ful logging, default  
`yes`   - only error and warn mesages  
`super` - without any output  

__BS_DB_PATH__ -- path to store your db  
  `/var/lib/bow`  

__BS_DB_NAME__ -- name of db and db file  
  `asapdrf.db`

__BS_SERVE_ADD__ -- port to serve  
  `:19808` - default port for app  

__BS_CHECKER_TIMEOUT__ -- timeout for daemons tictac in seconds  
  `300` -- default time to renew data

__BS_TIME_WATCH__ -- timewatch for daemons and handlers  
  `yes` - print time of execution in log (level2)  


BoltDB levels
```
{DB} --{ _info } => [ version:V ]
    \  
(b){repositories}   ,{ _info } => [ host:bow.example.com pass:test scheme:https user:test secure:false]
      \            /            
(br){reponame[N]}-*-{ _names }-*-{ imagename }---[ datetime => last_name ]
        \          \                                                   :
 (brc){catalog}     '--{ _namesgraph }                                 :
          \                                                            :
 (brcn){imagename[N]}--*--[ _valid => 0 || 1 ]                         :
            \           \-{ _namepair }-[namespace : name]             :
             \           \-{ _uploads }---[ date => count ]            :
              \                                                        :
      (brcnt){tags}--*                                                 :
                      \                                                :
 [ _valid => 0 || 1 ]--*--{ _uploads }---[ date => count ]             :
                        \                                              :
                         \--[ digest => header:digest ]................:
                          \
                           \--{ history }--[ datetime => (command + blob:{sha256, size} ]
                            *--{ _totalsizehuman }--[ datetime => size ]
                            |--{ _totalsizebytes }--[ datetime => size ]
                            '--[ _parent => name:tag ]
```
