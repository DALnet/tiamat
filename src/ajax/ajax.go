/****************************************************/
/* DALnet Tiamat - Copyright (C) 2017 Kobi Shmueli. */
/* See the LICENSE file for more information.       */
/****************************************************/

package ajax

import (
  "fmt"
  "net/http"
  "encoding/json"
  types "../types"
  conf "../conf"
)

func ajax_test(w http.ResponseWriter, req *http.Request, p *types.Page, userinfo types.UserInfo, myconfig conf.ConfigType) {
  var myjson []interface{}

  myjson[0] = "test"

  /* Send the json to the client... */
  real_json, err := json.Marshal(myjson)
  if err != nil {
      p.ErrorMsg = "Error #32148: " + err.Error()
      return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(real_json)
}

type AjaxFunc func(http.ResponseWriter, *http.Request, *types.Page, types.UserInfo, conf.ConfigType)

func Ajax_handler(w http.ResponseWriter, req *http.Request, p *types.Page, userinfo types.UserInfo, myconfig conf.ConfigType) {
  funclist := map[string]AjaxFunc {
    "ajax_test": ajax_test,
  }

  for key, value := range funclist {
    if key == p.Title {
      value(w, req, p, userinfo, myconfig)
      if p.ErrorMsg != "" {
        fmt.Fprintf(w, "Ajax Error in %s: %s\n", p.Title, p.ErrorMsg)
      }
      return
    }
  }
  //http.NotFound(w, req)
  fmt.Fprintf(w, "Couldn't find ajax!")
}
