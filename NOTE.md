BS_LOG_SILENT -- level of output logging  
nothing - ful logging, default  
`yes`   - only error mesages  
`super` - without any output  

BS_DB_PATH -- path to store your db  
  `/var/lib/bow`  

BS_DB_NAME -- name of db and db file  
  `asapdrf.db`

BS_SERVE_ADD -- port to serve  
  `:19808` - default port for app  

BoltDB levels
```
{DB}
    \  
(b){repositories}
      \
(br){reponame[N]}
        \   
 (brc){catalog}
          \---{ _names }---[ id => last_name ]
 (brcn){imagename[N]}--*--[ _valid => 0 || 1 ]
            \           \
             \           \- { _uploads }---[ date => count ]
              \            
      (brcnt){tags}--*
                      \
 [ _valid => 0 || 1 ]--*--{ _uploads }---[ date => count ]
                        \--[ digest => header:digest ]
                         \--{ history }--[ datetime => (command + blob:{sha256, size} ]
                          *--{ _totalsizehuman }--[ datetime => size ]
                          |--{ _totalsizebytes }--[ datetime => size ]
```

//IDEAS  
- [x] Delete b{repositories} and brc{catalog} (static db fields not changed and used)
Make db.Put(args []string, level int, key string, value string ){
  for i 1 to level-1{
    read bucket @ args[i]
  }
  create if not exist bucket @ level
  put key value
}
Same for get and delete
- [ ] refactor all db midleware to use common get put delete
