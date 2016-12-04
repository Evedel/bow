package db

// import("say")

func GetRepos() (repos map[string]string){
  pairs := GetSimplePairsFromBucket([]string{})
  delete(pairs, "_info");
  return pairs
}

func GetRepoPretty(repo string) (pretty map[string]string){
  return GetSimplePairsFromBucket([]string{repo, "_info"})
}

func CreateRepo(params map[string][]string) {
  name := params["name"][0]
  PutSimplePairToBucket([]string{name, "_info"}, "host", params["host"][0])
  PutSimplePairToBucket([]string{name, "_info"}, "pass", params["pass"][0])
  PutSimplePairToBucket([]string{name, "_info"}, "user", params["user"][0])
  PutSimplePairToBucket([]string{name, "_info"}, "scheme", params["scheme"][0])
  PutSimplePairToBucket([]string{name, "_info"}, "name", name)

  if _, ok := params["secure"]; ok {
    PutSimplePairToBucket([]string{name, "_info"}, "secure", "true")
  } else {
    PutSimplePairToBucket([]string{name, "_info"}, "secure", "false")
  }
}

func DeleteRepo(repo string){
  DeleteBucket([]string{repo})
}
