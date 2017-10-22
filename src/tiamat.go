/****************************************************/
/* DALnet Tiamat - Copyright (C) 2017 Kobi Shmueli. */
/* See the LICENSE file for more information.       */
/****************************************************/

package main

import (
  "fmt"
  "net/http"
  "html/template"
  "io/ioutil"
  "crypto/rand"
  "encoding/base64"
  "net"
  "log"
  "time"
  conf "./conf"
  lang "./lang"
  ajax "./ajax"
  types "./types"
  mem "./mem"
  bahamut "./bahamut"
  funcs "./funcs"
)

const (
  MY_VERSION = "1.0.0a"
)

type FastCGIServer struct{}

type Templates struct {
    tmpl map[string]*template.Template
}

var (
    err error
    templates *Templates
    myconfig conf.ConfigType
)

func NewTemplatesFromGlob(glob string) (t *Templates, err error) {
    var (
        root *template.Template
    )
    if root, err = template.New("root").ParseGlob(glob); err != nil {
        return
    }
    t = &Templates{
        tmpl: map[string]*template.Template{
            "root": root,
        },
    }
    return
}

func (t *Templates) Split(from string, to string) (err error) {
    var (
        tmpl *template.Template
    )
    if tmpl, err = t.tmpl[from].Clone(); err != nil {
        return
    }
    t.tmpl[to] = tmpl
    return
}

func (t *Templates) New(to string) (err error) {
    var (
        tmpl *template.Template
    )
    tmpl = template.New("root")
    t.tmpl[to] = tmpl
    return
}

func (t *Templates) Get(key string) (tmpl *template.Template) {
    return t.tmpl[key]
}

func getunixtime() (res int32) {
  res = int32(time.Now().Unix())
  return
}

// GenerateRandomBytes returns securely generated random bytes. 
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
// Taken from https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/
func GenerateRandomBytes(n int) ([]byte, error) {
  b := make([]byte, n)
  _, err := rand.Read(b)
  // Note that err == nil only if we read len(b) bytes.
  if err != nil {
    return nil, err
  }

  return b, nil
}

// case the caller should not continue.
// Taken from https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/
func GenerateRandomString(s int) (string, error) {
  b, err := GenerateRandomBytes(s)
  return base64.URLEncoding.EncodeToString(b), err
}

func check_auth(w http.ResponseWriter, req *http.Request, p *types.Page, userinfo *types.UserInfo) {
  userinfo.User = "Guest"
  userinfo.Account = "none"
  userinfo.Groups = "all-guests"

  cookie_user, _ := req.Cookie("tiamat_user")
  cookie_key, _ := req.Cookie("tiamat_key")

  req.ParseForm()
  /* If we got a POST request, let's check it first... */
  if req.FormValue("ircnick") != "" {
    if(mem.Find_Client(req.FormValue("ircnick")) != nil) {
      p.ErrorMsg = "Nickname already in use";
      return
    }
    if !funcs.IsValidNick(req.FormValue("ircnick")) {
      p.ErrorMsg = "Invalid nickname provided";
      return
    }
    userinfo.User = req.FormValue("ircnick")
    n := mem.New_Client(userinfo.User, "", req.RemoteAddr)
    skey, _ := GenerateRandomString(10)
    n.CookieKey = skey
    var cookie_user http.Cookie
    var cookie_key http.Cookie
    cookie_user.Name = "tiamat_user"
    cookie_user.Value = userinfo.User
    cookie_user.MaxAge = 3600*24*14
    cookie_key.Name = "tiamat_key"
    cookie_key.Value = skey
    cookie_key.MaxAge = 3600*24*14
    if myconfig.CookieDomain != "" {
      cookie_user.Domain = myconfig.CookieDomain
      cookie_key.Domain = myconfig.CookieDomain
    }
    http.SetCookie(w, &cookie_user)
    http.SetCookie(w, &cookie_key)
    return
  }

  /* Check if we have a valid cookie... */
  if cookie_user != nil && mem.Find_Client(cookie_user.Value) != nil {
    n := mem.Find_Client(cookie_user.Value)
    if (n.CookieKey == cookie_key.Value) {
      userinfo.User = n.Nick
    }
  }

  return
}

func make_data_about(w http.ResponseWriter, req *http.Request, p *types.Page, userinfo types.UserInfo) {
  p.Data["MY_VERSION"] = MY_VERSION
}

func handler(w http.ResponseWriter, req *http.Request) {
  var p types.Page
  var userinfo types.UserInfo

  title := req.URL.Path[1:]
  fmt.Printf("Debugging: url=%s details: %+v\n", req.URL.Path, req)
  if title=="" {
    title = "index"
  }
  if title[len(title)-4:len(title)] == ".php" {
    title = title[0:len(title)-4]
  }
  p.Title = title
  p.ErrorMsg = ""

  template := templates.Get(title)
  if template == nil && title[0:4] != "ajax" {
    http.NotFound(w, req)
    return
  }

  check_auth(w, req, &p, &userinfo)

  p.FullName = userinfo.User
  p.User = userinfo.User
  p.Role = "User"
  p.LangList = lang.GetLangList()
  p.Lang = lang.GetAll()
  p.Data = make(map [string]string)

  if userinfo.User == "Guest" {
    template = templates.Get("login")
    if template == nil {
      http.NotFound(w, req)
      return
    }
  } else if title == "about" {
    make_data_about(w, req, &p, userinfo)
  } else if len(title)>4 && title[0:4] == "ajax" {
    ajax.Ajax_handler(w, req, &p, userinfo, myconfig)
    return
  }

  if err = template.Execute(w, p); err != nil {
    fmt.Fprintf(w, "Couldn't compile %s template! Error: %s\n", title, err.Error())
  }
}

func http_listen() {
  http.HandleFunc("/", handler)
  http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static-files/css"))))
  http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("static-files/img"))))
  http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static-files/js"))))
  http.Handle("/font-awesome/", http.StripPrefix("/font-awesome/", http.FileServer(http.Dir("static-files/font-awesome"))))

  fmt.Print("Listening to port " + myconfig.HttpPort + "... ");
  err := http.ListenAndServe(":" + myconfig.HttpPort, nil)
  if err != nil {
    log.Fatal("Error: ", err)
  }
}

func main() {
  fmt.Println("Running Tiamat v" + MY_VERSION)
  conf.Init(MY_VERSION)
  conf.Load()
  myconfig = conf.Get()
  mem.Init()
  bahamut.Init(myconfig)
  lang.Load()
  fmt.Print("Loading templates... ");
  if templates, err = NewTemplatesFromGlob("templates/base.html"); err != nil {
    fmt.Println("Couldn't load templates! Error: " + err.Error())
  }
  files, _ := ioutil.ReadDir("templates")
  for _, f := range files {
    if f.IsDir() {
      continue
    }
    fileext := f.Name()[len(f.Name())-4:len(f.Name())]
    fileonly := f.Name()[0:len(f.Name())-5]
    if fileext != "html" {
      continue
    }

    if f.Name() == "base.html" {
      // Base is special is already loaded...
      continue
    }

    fmt.Print(f.Name() + ": ")

    if f.Name() == "login.html" {
      // Login is special too...
      if err = templates.New(fileonly); err != nil {
        fmt.Println("Couldn't create template!")
      }
    } else if err = templates.Split("root", fileonly); err != nil {
      fmt.Println("Couldn't split template!")
    }
    if _, err = templates.Get(fileonly).ParseFiles("templates/" + f.Name()); err != nil {
      fmt.Println("Couldn't get template! Error: " + err.Error())
    }
    if err == nil {
      fmt.Print("OK. ");
    }
  }
  fmt.Println("Done.");

  if myconfig.UplinkName == "" {
    log.Fatal("Error: No uplink server is set!")
  }
  fmt.Print("Connecting to uplink " + myconfig.UplinkName + "... ");
  conn, err := net.Dial("tcp", myconfig.UplinkHost + ":" + myconfig.UplinkPort)
  if err != nil {
    log.Fatal("Error: ", err)
  }
  bahamut.Connect_time = time.Now()
  fmt.Println("Done.");
  go http_listen()
  line := ""
  for {
    //fmt.Println("-Debug- Looping...\n");
    buff := make([]byte, 1024)
    n, _ := conn.Read(buff)
    //fmt.Printf("-Debug- Receive: %s", buff[:n])
    for i := 0; i < n; i++ {
      if buff[i] == '\n' {
        //fmt.Printf("-Debug- Line: %s", line)
        bahamut.Parser(conn, line)
        line = ""
      }
      line = line + string(buff[i]);
    }
    if line == "" {
      time.Sleep(time.Millisecond)
    }
  }
  fmt.Println("Byebye!")
}
