package imgcache

import (
  "fmt"
  "io"
  "io/ioutil"
  "net/http"
  "net/url"
  "os"
  "path/filepath"
  "strings"
  "sync"
)

// Cache : Cache strcut
type Cache struct {
  path     string
  cacheURL string
  caching  bool
  images   map[string]string
  Queue    []string
  Cache    []string
  Image    imageFunc
  sync.RWMutex
}

type imageFunc struct {
  GetURL  func(string) string
  Caching func()
  Remove  func()
}

// New : New cahce
func New(path, chacheURL string, caching bool) (c *Cache, err error) {

  c = &Cache{}

  c.images = make(map[string]string)
  c.path = path
  c.cacheURL = chacheURL
  c.caching = caching
  c.Queue = []string{}
  c.Cache = []string{}

  var queue []string

  c.Image.GetURL = func(src string) (cacheURL string) {

    c.Lock()
    defer c.Unlock()

    src = strings.Trim(src, "\r\n")

    if c.caching == false {
      return src
    }

    u, err := url.Parse(src)
    if err != nil || len(filepath.Ext(u.Path)) == 0 {
      return src
    }

    var filename = fmt.Sprintf("%s%s", strToMD5(src), filepath.Ext(u.Path))
    if cacheURL, ok := c.images[fmt.Sprintf("%s%s", strToMD5(src), filepath.Ext(u.Path))]; ok {
      return cacheURL
    }

    if indexOfString(filename, c.Cache) == -1 {

      if indexOfString(src, c.Queue) == -1 {
        c.Queue = append(c.Queue, src)
      }

    } else {
      c.images[filename] = c.cacheURL + filename
      src = c.cacheURL + filename
    }

    /*
       if _, err := os.Stat(c.path + filename); err != nil {
         //c.images[filename] = c.cacheURL + filename
         if indexOfString(src, c.Queue) == -1 {
           c.Queue = append(c.Queue, src)
         }
       } else {
         c.images[filename] = c.cacheURL + filename
       }
    */

    return src
  }

  c.Image.Caching = func() {

    c.Lock()
    defer c.Unlock()

    var filename string

    for _, src := range c.Queue {

      resp, err := http.Get(src)
      if err != nil {
        continue
      }
      defer resp.Body.Close()

      if resp.StatusCode != http.StatusOK {
        continue
      }

      filename = fmt.Sprintf("%s%s%s%s", c.path, string(os.PathSeparator), strToMD5(src), filepath.Ext(src))

      file, err := os.Create(filename)
      if err != nil {
        continue
      }

      defer file.Close()

      _, err = io.Copy(file, resp.Body)
      if err != nil {
        continue
      }

      u, err := url.Parse(src)
      if err == nil {
        c.images[fmt.Sprintf("%s%s", strToMD5(src), filepath.Ext(u.Path))] = c.cacheURL + filename
      }

      queue = append(queue, src)

    }

    for _, q := range queue {
      c.Queue = removeStringFromSlice(q, c.Queue)
    }

  }

  c.Image.Remove = func() {

    c.Lock()
    defer c.Unlock()

    files, err := ioutil.ReadDir(c.path)
    if err != nil {
      return
    }

    for _, file := range files {

      switch c.caching {

      case true:
        if _, ok := c.images[file.Name()]; !ok {
          os.RemoveAll(c.path + file.Name())
        }

      case false:
        os.RemoveAll(c.path + file.Name())
      }

    }

  }

  files, err := ioutil.ReadDir(c.path)
  if err != nil {
    return
  }

  for _, file := range files {
    c.Cache = append(c.Cache, file.Name())
  }

  return
}
