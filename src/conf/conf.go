/****************************************************/
/* DALnet Tiamat - Copyright (C) 2017 Kobi Shmueli. */
/* See the LICENSE file for more information.       */
/****************************************************/

package conf

import (
  "fmt"
  "strconv"
  "github.com/vaughan0/go-ini"
)

type ConfigType struct {
  ServerName string
  ServerDesc string
  HttpPortInt int
  HttpPort string
  HttpsPort string
  UplinkName string
  UplinkHost string
  UplinkPort string
  UplinkPass string
  CookieDomain string
  Version string
}

var Config ConfigType

func Init(my_version string) {
  Config.ServerName = ""
  Config.ServerDesc = ""
  Config.HttpPortInt = 8080
  Config.HttpPort = "8080"
  Config.HttpsPort = "8443"
  Config.UplinkName = ""
  Config.UplinkHost = ""
  Config.UplinkPort = ""
  Config.UplinkPass = ""
  Config.CookieDomain = ""
  Config.Version = my_version
}

func Load() (res bool) {
  fmt.Print("Loading config... ")
  file, err := ini.LoadFile("tiamat.conf")
  if err != nil {
    fmt.Println("Error: " + err.Error())
    res = false
    return
  }
  Config.ServerName, _ = file.Get("server", "name")
  Config.ServerDesc, _ = file.Get("server", "description")
  temptxt, _ := file.Get("web", "port")
  Config.HttpPortInt, _ = strconv.Atoi(temptxt)
  Config.HttpPort, _ = file.Get("web", "port")
  Config.HttpsPort, _ = file.Get("web", "sslport")
  Config.CookieDomain, _ = file.Get("web", "cookiedomain")
  Config.UplinkName, _ = file.Get("uplink", "name")
  Config.UplinkHost, _ = file.Get("uplink", "host")
  Config.UplinkPort, _ = file.Get("uplink", "port")
  Config.UplinkPass, _ = file.Get("uplink", "pass")

  fmt.Println("Done.")
  res = true
  return
}

func Get() (res ConfigType) {
  res = Config
  return
}
