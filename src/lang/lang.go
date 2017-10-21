/****************************************************/
/* DALnet Tiamat - Copyright (C) 2017 Kobi Shmueli. */
/* See the LICENSE file for more information.       */
/****************************************************/

package lang

import (
  "fmt"
  "io/ioutil"
  "encoding/json"
)

var langtexts map[string]interface{}

func Load() (res bool) {
  fmt.Print("Loading language files... ")
  data, err := ioutil.ReadFile("lang/english.json")
  if err != nil {
    fmt.Println("Error: " + err.Error())
    res = false
    return
  }
  err = json.Unmarshal(data, &langtexts)
  if err != nil {
    fmt.Println("Error: " + err.Error())
    res = false
    return
  }

  fmt.Println("Done.")
  res = true
  return
}

func Reply(txt string) (res string) {
  res = langtexts[txt].(string)
  return
}

func GetAll() (res map[string]string) {
  res = make(map [string]string)
  for key, value := range langtexts {
    res[key] = value.(string)
  }
  return
}

func GetLangList() (res map[string]int) {
  res = make(map [string]int)
  res["English"] = 0
  //TODO: For future use...
  //res["Hebrew"] = 1
  return
}
