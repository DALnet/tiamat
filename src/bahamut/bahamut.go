/****************************************************/
/* DALnet Tiamat - Copyright (C) 2017 Kobi Shmueli. */
/* See the LICENSE file for more information.       */
/****************************************************/

package bahamut

import (
  "fmt"
  "net"
  "strings"
  "time"
  conf "../conf"
  types "../types"
  mem "../mem"
)

const MAX_PARAMS = 10
var connected bool = false
var synced bool = false
var conn net.Conn
var myconfig conf.ConfigType
var myuplink *types.ServerType
var Connect_time time.Time

func ssend(buf string, args ...interface{}) {
  if args == nil {
    fmt.Printf("Sending: " + buf + "\r\n");
    fmt.Fprintf(conn, buf + "\r\n")    
  } else {
    fmt.Printf("Sending: " + buf + "\r\n", args...);
    fmt.Fprintf(conn, buf + "\r\n", args...)
  }
}

func s_raw(number int, client *types.ClientType, buf string, args ...interface{}) {
  var newstr string

  if client != nil {
    newstr = fmt.Sprintf(":%s %d %s :", myconfig.ServerName, number, client.Nick);
    if args == nil {
      ssend(newstr + buf)
    } else {
      ssend(newstr + buf, args)
    }
  }
}

func got_notice(server *types.ServerType, client *types.ClientType, parv[MAX_PARAMS+1] string, parc int) {
  fmt.Printf("-Debug- Got notice\n");
  if !connected {
    fmt.Printf("-Debug- Sending server intro...\n");
    ssend("PASS %s :TS", myconfig.UplinkPass)
    ssend("CAPAB SSJOIN NOQUIT BURST UNCONNECT NICKIP TSMODE")
    ssend("SERVER %s 1 :%s", myconfig.ServerName, myconfig.ServerDesc)
    connected = true
  }
}

func got_ping(server *types.ServerType, client *types.ClientType, parv[MAX_PARAMS+1] string, parc int) {
  if parv[2] != "" {
    ssend("PONG %s :%s", parv[1], parv[2]);
  } else {
    ssend("PONG %s :%s", myconfig.ServerName, parv[1]);
  }
}

func got_pong(server *types.ServerType, client *types.ClientType, parv[MAX_PARAMS+1] string, parc int) {
  if !synced {
    ssend("GNOTICE :%s has synched in %s.", myconfig.ServerName, time.Since(Connect_time).String())
    synced = true
  }
}

func got_error(server *types.ServerType, client *types.ClientType, parv[MAX_PARAMS+1] string, parc int) {
  fmt.Printf("-Debug- Got error from uplink: %s\n", parv[1])
  connected = false
  synced = false
  myuplink = nil
}

func got_server(server *types.ServerType, client *types.ClientType, parv[MAX_PARAMS+1] string, parc int) {
  new := mem.New_Server(parv[1], parv[3])

  if server != nil {
    new.Uplink = server.Name
  } else {
    myuplink = new
    if !synced {
      ssend("PING :%s", myconfig.ServerName)
      ssend(":%s BURST 0", myconfig.ServerName)
    }
  }
}

func got_info(server *types.ServerType, client *types.ClientType, parv[MAX_PARAMS+1] string, parc int) {
  s_raw(371, client, "Tiamat v" + myconfig.Version);
  s_raw(374, client, "End of /INFO list.");
}

/*
 * got_nick
 * parv[0] = sender prefix
 * parv[1] = nickname
 * parv[2] = hopcount when new user; TS when nick change
 * parv[3] = TS
 * ---- new user only below ----
 * parv[4] = umode
 * parv[5] = username
 * parv[6] = hostname
 * parv[7] = server
 * parv[8] = serviceid
 * parv[9] = IP
 * parv[10] = ircname
 * -- endif
 */
func got_nick(server *types.ServerType, client *types.ClientType, parv[MAX_PARAMS+1] string, parc int) {
  if client != nil && parc == 2 {
    /* Nick change */
    client.Nick = parv[1]

    return
  }

  if client == nil && parc == 10 {
    new := mem.New_Client(parv[1], parv[5], parv[6])
    new.Server = parv[7]

    return
  }
}

/* got_quit
 * parv[0] = nick
 * parv[1] = reason
 */
func got_quit(server *types.ServerType, client *types.ClientType, parv[MAX_PARAMS+1] string, parc int) {
  mem.Del_Client(parv[0])
}

/* got_kill
 * parv[0] = sender
 * parv[1] = target nick
 * parv[2] = sourceserver!sourcehost!sourcenick (reason)
 */
func got_kill(server *types.ServerType, client *types.ClientType, parv[MAX_PARAMS+1] string, parc int) {
  mem.Del_Client(parv[1])
}

type BahamutFunc func(server *types.ServerType, client *types.ClientType, parv[MAX_PARAMS+1] string, parc int)

func Parser(d net.Conn, msg string) {
  var params[MAX_PARAMS+1] string
  var cmd string
  var client *types.ClientType
  var server *types.ServerType

  if msg == "" {
    return
  }

  conn = d;

  funclist := map[string]BahamutFunc {
    "NOTICE": got_notice,
    "ERROR": got_error,
    "PING": got_ping,
    "PONG": got_pong,
    "SERVER": got_server,
    "INFO": got_info,
    "NICK": got_nick,
    "QUIT": got_quit,
    "KILL": got_kill,
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
    params[0] = params[0][1:]
    if strings.Index(params[0], ".") > 0 {
      server = mem.Find_Server(params[0])
    } else {
      client = mem.Find_Client(params[0])
    }
  } else {
    cmd = params[0]
    server = nil
    client = nil
  }

  for key, value := range funclist {
    if key == cmd {
      value(server, client, params, cur_param)
      return
    }
  }

  fmt.Printf("Bahamut Error: Unknown command: %s\n", cmd)
}

func Init(realconfig conf.ConfigType) {
  myconfig = realconfig
  myuplink = nil
}
