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
