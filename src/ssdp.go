package src

import (
  "fmt"
  "log"
  "os"
  "os/signal"
  "time"

  "github.com/koron/go-ssdp"
)

// SSDP : SSPD / DLNA Server
func SSDP() {

  showInfo(fmt.Sprintf("SSDP / DLNA:%t", Settings.SSDP))

  if Settings.SSDP == false {
    return
  }

  time.Sleep(10 * time.Second)
  ad, err := ssdp.Advertise(
    "upnp:"+System.AppName,                   // send as "ST"
    System.DeviceID+"::upnp:"+System.AppName, // send as "USN"
    System.URLBase+"/device.xml",             // send as "LOCATION"
    System.AppName,                           // send as "SERVER"
    1800)                                     // send as "maxAge" in "CACHE-CONTROL"

  if err != nil {
    ShowError(err, 000)
  }

  // Debug SSDP
  if System.Flag.Debug == 3 {
    ssdp.Logger = log.New(os.Stderr, "[SSDP] ", log.LstdFlags)
  }

  var aliveTick <-chan time.Time
  var ai = 10

  if ai > 0 {
    aliveTick = time.Tick(time.Duration(ai) * time.Second)
  } else {
    aliveTick = make(chan time.Time)
  }

  quit := make(chan os.Signal, 1)
  signal.Notify(quit, os.Interrupt)

loop:

  for {

    select {

    case <-aliveTick:
      ad.Alive()
    case <-quit:
      os.Exit(0)
      break loop

    }

  }

  ad.Bye()
  ad.Close()
}
