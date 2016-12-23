package db

import(
  // "say"
  "strconv"
)

func GetRepos() (repos map[string]string){
  pairs := GetAllPairsFromBucket([]string{})
  delete(pairs, "_info");
  return pairs
}

func GetRepoPretty(repo string) (pretty map[string]string){
  return GetAllPairsFromBucket([]string{repo, "_info"})
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
  PutBucketToBucket([]string{name, "catalog"})
  PutBucketToBucket([]string{name, "_namesgraph"})
  PutBucketToBucket([]string{name, "_names"})
}

func DeleteRepo(repo string){
  DeleteBucket([]string{repo})
}

func PutTagDigest(er, en, et, shortsizedt, didgesr string) {
  uploadsnt := GetValueFromBucket([]string{ er, "catalog", en, et, "_uploads"}, shortsizedt)
  if uploadsnt == "" {
    PutSimplePairToBucket([]string{ er, "catalog", en, et, "_uploads"}, shortsizedt, "1")
  } else {
    val, _ := strconv.Atoi(string(uploadsnt))
    val++
    PutSimplePairToBucket([]string{ er, "catalog", en, et, "_uploads"}, shortsizedt, strconv.Itoa(val))
  }
  uploadsnn := GetValueFromBucket([]string{ er, "catalog", en,"_uploads"}, shortsizedt)
  if uploadsnn == "" {
    PutSimplePairToBucket([]string{ er, "catalog", en, "_uploads"}, shortsizedt, "1")
  } else {
    val, _ := strconv.Atoi(string(uploadsnn))
    val++
    PutSimplePairToBucket([]string{ er, "catalog", en, "_uploads"}, shortsizedt, strconv.Itoa(val))
  }
  PutSimplePairToBucket([]string{ er, "catalog", en, et}, "digest", didgesr)
}

func GetTags( er, en string) (tags []string){
  tagsdb := GetAllPairsFromBucket([]string{ er, "catalog", en})
  delete(tagsdb, "_valid")
  delete(tagsdb, "_uploads")
  for et, _ := range tagsdb {
    if valid := GetValueFromBucket([]string{ er, "catalog", en, et}, "_valid"); valid == "1" {
      tags = append(tags, et)
    }
  }
  return
}

func AddTags( er, en string, tags []string){
  tagsdb := GetAllPairsFromBucket([]string{ er, "catalog", en})
  delete(tagsdb, "_valid")
  delete(tagsdb, "_uploads")
  for etdb, _ := range tagsdb {
    PutSimplePairToBucket([]string{ er, "catalog", en, etdb}, "_valid", "0")
  }
  for _, etrp := range tags {
    // Just put new pairs, 'couse it will check and create bucket of tag in case it isn't exist
    PutSimplePairToBucket([]string{ er, "catalog", en, etrp}, "_valid", "1")
    PutBucketToBucket([]string{ er, "catalog", en, etrp, "_uploads"})
    PutBucketToBucket([]string{ er, "catalog", en, etrp, "history"})
    PutBucketToBucket([]string{ er, "catalog", en, etrp, "_totalsizehuman"})
    PutBucketToBucket([]string{ er, "catalog", en, etrp, "_totalsizebytes"})
    PutBucketToBucket([]string{ er, "catalog", en, etrp, "_parent"})
    PutBucketToBucket([]string{ er, "_names", en + ":" + etrp})
  }
}

func GetCatalog(er string) (catalog []string){
  catalogdb := GetAllPairsFromBucket([]string{ er, "catalog"})
  for en, _ := range catalogdb {
    if valid := GetValueFromBucket([]string{ er, "catalog", en}, "_valid"); valid == "1" {
      catalog = append(catalog, en)
    }
  }
  return
}

func AddCatalog(er string, catalog []string) {
  catalogdb := GetAllPairsFromBucket([]string{ er, "catalog"})
  for endb, _ := range catalogdb {
    PutSimplePairToBucket([]string{ er, "catalog", endb}, "_valid", "0")
  }
  for _, enrp := range catalog {
    PutSimplePairToBucket([]string{ er, "catalog", enrp}, "_valid", "1")
    PutBucketToBucket([]string{ er, "catalog", enrp, "_uploads"})
  }
}
