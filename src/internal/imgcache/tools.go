package imgcache

import (
  "crypto/md5"
  "encoding/hex"
)

func strToMD5(str string) string {
  md5Hasher := md5.New()
  md5Hasher.Write([]byte(str))
  return hex.EncodeToString(md5Hasher.Sum(nil))
}

func indexOfString(str string, slice []string) int {

  for i, v := range slice {
    if str == v {
      return i
    }
  }

  return -1
}

func removeStringFromSlice(str string, slice []string) []string {

  var i = indexOfString(str, slice)

  if i != -1 {
    slice = append(slice[:i], slice[i+1:]...)
  }

  return slice
}
