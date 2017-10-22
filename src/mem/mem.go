/****************************************************/
/* DALnet Tiamat - Copyright (C) 2017 Kobi Shmueli. */
/* See the LICENSE file for more information.       */
/****************************************************/

package mem

import (
  "strings"
  types "../types"
)

var clients map[string]interface{}
var servers map[string]interface{}

func Find_Client(nick string)(*types.ClientType) {
  for key, value := range clients {
    if(strings.ToLower(key) == strings.ToLower(nick)) {
      return value.(*types.ClientType)
    }
  }

  return nil
}

func New_Client(nick string, user string, host string) (*types.ClientType) {
  n := new(types.ClientType)
  n.Nick = nick
  n.User = user
  n.Host = host
  n.Server = ""
  clients[n.Nick] = n

  return n
}

func New_Server(name string, desc string) (*types.ServerType) {
  n := new(types.ServerType)
  n.Name = name
  n.Desc = desc
  servers[n.Name] = n

  return n
}

func Find_Server(server string)(*types.ServerType) {
  for key, value := range servers {
    if(strings.ToLower(key) == strings.ToLower(server)) {
      return value.(*types.ServerType)
    }
  }

  return nil
}

func Init() () {
  clients = make(map [string]interface{})
  servers = make(map [string]interface{})

  return
}
