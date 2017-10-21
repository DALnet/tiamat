/****************************************************/
/* DALnet Tiamat - Copyright (C) 2017 Kobi Shmueli. */
/* See the LICENSE file for more information.       */
/****************************************************/

package mem

import (
  "strings"
  types "../types"
)

var nicks map[string]interface{}

func FindNick(nick string) (found bool) {
  found = false

  for key, _ := range nicks {
    if(strings.ToLower(key) == strings.ToLower(nick)) {
      found = true
      return
    }
  }

  return
}

func GetNick(nick string)(*types.NickType) {
  var res types.NickType

  for key, value := range nicks {
    if(strings.ToLower(key) == strings.ToLower(nick)) {
      return value.(*types.NickType)
    }
  }

  return &res
}

func GetAllNicks() (res map[string]string) {
  res = make(map [string]string)
  for key, value := range nicks {
    res[key] = value.(string)
  }

  return
}

func AddNick(nick string, user string, host string) (*types.NickType) {
  n := new(types.NickType)
  n.Nick = nick
  n.User = user
  n.Host = host
  nicks[n.Nick] = n

  return n
}

func Init() () {
  nicks = make(map [string]interface{})

  return
}
