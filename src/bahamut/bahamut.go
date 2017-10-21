/****************************************************/
/* DALnet Tiamat - Copyright (C) 2017 Kobi Shmueli. */
/* See the LICENSE file for more information.       */
/****************************************************/

package bahamut

import (
  "fmt"
  "net"
  conf "../conf"
)

const MAX_PARAMS = 10
var connected bool = false
var conn net.Conn
var myconfig conf.ConfigType

func ssend(buf string, args ...interface{}) {
  fmt.Printf("Sending: " + buf + "\r\n", args...);
  fmt.Fprintf(conn, buf + "\r\n", args...)
}

func got_notice(params[MAX_PARAMS+1] string) {
  fmt.Printf("-Debug- Got notice\n");
  if !connected {
    fmt.Printf("-Debug- Sending server intro...\n");
    ssend("PASS %s :TS", myconfig.UplinkPass)
    ssend("CAPAB SSJOIN NOQUIT BURST UNCONNECT NICKIP TSMODE")
    ssend("SERVER %s 1 :%s", myconfig.ServerName, myconfig.ServerDesc)
    connected = true
  }
}

func got_ping(params[MAX_PARAMS+1] string) {
  if params[2] != "" {
    ssend("PONG %s :%s", params[1], params[2]);
  } else {
    ssend("PONG %s :%s", myconfig.ServerName, params[1]);
  }
}

func got_error(params[MAX_PARAMS+1] string) {
  fmt.Printf("-Debug- Got error from uplink: %s\n", params[1])
  connected = false
}

type BahamutFunc func(params[MAX_PARAMS+1] string)

func Parser(d net.Conn, msg string) {
  var params[MAX_PARAMS+1] string
  var cmd string

  if msg == "" {
    return
  }

  conn = d;

  funclist := map[string]BahamutFunc {
    "NOTICE": got_notice,
    "ERROR": got_error,
    "PING": got_ping,
  }

  cur_param := 0
  last := false
  for i := 0; i < len(msg); i++ {
    if msg[i] == '\n' || msg[i] == '\r' {
      continue
    }
    if msg[i] == ' ' && !last && cur_param<MAX_PARAMS {
      cur_param++
      continue
    }
    if cur_param!=0 && params[cur_param] == "" && msg[i] == ':' {
      last = true
      continue
    }
    params[cur_param] = params[cur_param] + string(msg[i])
  }

  for i := cur_param+1; i <= MAX_PARAMS; i++ {
    params[i] = ""
  }

//  for i := 0; i <= MAX_PARAMS; i++ {
//    fmt.Printf("-Debug- param[%d] = " + params[i]+ "\n", i)
//  }

  if(params[0][0] == ':') {
    cmd = params[1]
  } else {
    cmd = params[0]
  }

  for key, value := range funclist {
    if key == cmd {
      value(params)
      return
    }
  }

  fmt.Printf("Bahamut Error: Unknown command: %s\n", cmd)
}

func Init(realconfig conf.ConfigType) {
  myconfig = realconfig
}
