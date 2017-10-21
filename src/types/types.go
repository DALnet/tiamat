/****************************************************/
/* DALnet Tiamat - Copyright (C) 2017 Kobi Shmueli. */
/* See the LICENSE file for more information.       */
/****************************************************/

package types

import (
  "html/template"
)

type Page struct {
    Title string
    Body  []byte
    ErrorMsg string
    User string
    FullName string
    Role string
    LangList map[string]int
    Lang map[string]string
    Data map[string]string
    JSData template.JS
}

type UserInfo struct {
  Id string
  User string
  Account string
  Groups string
  Email string
}

type NickType struct {
  Nick string
  User string
  Host string
  CookieKey string
}
