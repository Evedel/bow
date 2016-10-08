Bow
==
## As simple as possible frontend for your private docker registry
Pictures
==
![](develop/conf.png)  

![](develop/info.png)

![](develop/history.png)

![](develop/parents.png)

Features
==  
- v2 registry support only
- scary templates and interface in whole
- internal db (BoltBD) gives it ability to store info, and as result it responses much more faster then after direct api call
- app can pars, store and show info from registry such as:
 - image layers info:
   - name / tag
   - image size and pushes number
   - upload and push dates
 - image creating commands history
- it is possible to set multiple repositories and watch all registries in one place
- show statistics pretty, draw curves of uploads number and image sizes for tag with respects to dates
- find parent of image, in case, parent in the same repo (it is clickable!)
- __(new)__ show tree of parents for image/ build dependency tree for whole repo  
- __(the newest)__ now it supports insecure regestries
- __(killerfeature)__ enabled image deletion (registry --version >= 2.4.0)

Image deletion
==
To enable image deletion you need to:  
1. Run your registry with flag `-e REGISTRY_STORAGE_DELETE_ENABLED=true`  
Example:  
```
docker run -d -p 5000:5000 --restart=always --name registry \
  -v ./auth:/auth \
  -e "REGISTRY_AUTH=htpasswd" \
  -e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
  -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
  -e REGISTRY_STORAGE_DELETE_ENABLED=true \
  registry:2
```
2. Set up cron to run garbage collection  
Example:  
`10 * * * * docker exec registry bin/registry garbage-collect /etc/docker/registry/config.yml`

3. Also be aware, that there is a known [issue](https://github.com/docker/distribution/issues/1939) in docker registry 2.6 in earlier. It means, that if you delete an image from a repository, you will not able to push __the exactly same__ image in that repository. To fix it, you will need each time to perform rebuilding of image with `--no-cache` mode, or restarting of registry `docker restart registry`


See more:  
https://github.com/docker/docker-registry/issues/988#issuecomment-224280919  
https://docs.docker.com/registry/configuration/#delete  
https://docs.docker.com/registry/garbage-collection/#/how-garbage-collection-works

Prospects
==
I can say that this app almost fit my needs, so in all likelyhood, soon, I will not improve it hardly, but this is the list of ideas just for case:  
- ~~delete tags/images from repo by click~~
- "update repos" button (not wait for sleep time)
- repository info, size, number of pushes and so on
- make improvements on interface and visual side
- dynamically upload nice images from __icons8.com__ API

How to start use Bow
==
```
docker run -d \
   --name=Bow \
   -e BS_LOG_SILENT=yes \
   -v ~/db/bow:/var/lib/bow \
   -p 5001:19808 \
   evedel/bow
```
How to start contribute to Bow
==
If you have interest, you can easily start with
```
git clone https://github.com/evedel/bow.git
cd bow
docker-compose -f develop/devlinux.yml up -d
docker exec -it develop_golang_1 go get
docker exec -it develop_golang_1 go run main.go
```
Code and packages
==
This app is written on golang with use of standard packages and:  
https://github.com/boltdb/bolt -- BoltDB  
https://github.com/fatih/color -- to make cli shiny  
https://github.com/wader/disable_sendfile_vbox_linux -- to develop on docker-machine  
http://www.chartjs.org/ -- to draw best graphs ever  
https://github.com/fperucic/treant-js -- to draw parents graphs  
