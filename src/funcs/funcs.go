/****************************************************/
/* DALnet Tiamat - Copyright (C) 2017 Kobi Shmueli. */
/* See the LICENSE file for more information.       */
/****************************************************/

package funcs

import (
  "time"
  types "../types"
)

func getunixtime() (res int32) {
  res = int32(time.Now().Unix())
  return
}

func IsValidNick(nick string)(bool) {
  nicklen := len(nick)
  if nicklen < 1 || nicklen > types.NICKLEN { return false }
  if nick[0] == '-' { return false }
  if nick[0] >= '0' && nick[0] <= '9' { return false }

  for i := 0; i < nicklen; i++ {
    if nick[i] >= 'A' && nick[i] < '~' { continue }
    if nick[i] >= '0' && nick[i] <= '9' { continue }
    if nick[i] == '-' { continue }
    return false
  }

  return true
}
