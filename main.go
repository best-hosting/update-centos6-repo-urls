package main

import (
    "fmt"
    "os"
    "strings"
    "path/filepath"

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
        oldURL: "http://mirror.centos.org/centos/",
        newURL: "http://vault.centos.org/",
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

func main() {
    files, err := os.ReadDir("yum-repos")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    for _, d := range files {
        f := filepath.Join(iniDir, d.Name())
        if filepath.Ext(f) != ".repo" {
            fmt.Printf("Skipping not a repo file: %v\n", f)
            continue
        }
        fmt.Printf("Reading ini file %v\n", f)
        cfg, err := ini.Load(f)
        if err != nil {
            fmt.Printf("Fail to read file: %v", err)
            os.Exit(1)
        }

        b := replaceUrl(cfg)
        if b {
            cfg.SaveTo(filepath.Join(iniDir, d.Name() + "-local"))
        }
    }
}

