package main

import (
    "fmt"
    "os"
    "strings"
    "path/filepath"
    "time"
    "io"

    "gopkg.in/ini.v1"
)

type Name string
type Links struct { oldURL, newURL string }

var newBaseUrl = map[Name]Links {
    "epel": Links {
        oldURL: "https://download.fedoraproject.org/pub/epel/",
        newURL: "https://archives.fedoraproject.org/pub/archive/epel/",
    },
    "base": Links {
        oldURL: "http://mirror.centos.org/centos/$releasever",
        newURL: "http://vault.centos.org/6.10",
    },
}

func replaceUrl(cfg *ini.File) (changed bool) {
    for _, sec := range cfg.Sections() {
        if _, ok := newBaseUrl[Name(sec.Name())]; ok && sec.HasKey("baseurl") {
            l := newBaseUrl[Name(sec.Name())]
            x := sec.Key("baseurl").Value()
            y := strings.ReplaceAll(sec.Key("baseurl").Value(), l.oldURL, l.newURL)
            if x != y {
                changed = true
                fmt.Printf("Replacing in section '%v' baseurl '%v' with '%v'\n", sec.Name(), x, y)
                sec.Key("baseurl").SetValue(y)
            }
        }
    }
    return
}

var iniDir = "./yum-repos"
//var iniDir = "/etc/yum.repos.d"

func main() {
    files, err := os.ReadDir(iniDir)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    for _, d := range files {
        fp := filepath.Join(iniDir, d.Name())
        if filepath.Ext(fp) != ".repo" {
            fmt.Printf("Skipping not a repo file: %v\n", fp)
            continue
        }
        fmt.Printf("Reading ini file %v\n", fp)
        cfg, err := ini.Load(fp)
        if err != nil {
            fmt.Printf("Fail to read file: %v", err)
            os.Exit(1)
        }

        b := replaceUrl(cfg)
        if b {
            h0, err := os.Open(fp)
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
            defer h0.Close()

            bkp := fp + ".bkp_" + time.Now().Format("02-01-06_03:04:05")
            h1, err := os.Create(bkp)
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
            defer h1.Close()
            io.Copy(h1, h0)

            cfg.SaveTo(fp)
        }
    }
}

