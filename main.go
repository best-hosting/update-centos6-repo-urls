package main

import (
    "fmt"
    "os"
    "strings"

    "gopkg.in/ini.v1"
)

type Name string
type Links struct { oldURL, newURL string }

var newBaseUrl = map[Name]Links {
    "epel": Links {
        oldURL: "https://download.fedoraproject.org/pub/epel/",
        newURL: "https://archives.fedoraproject.org/pub/archive/epel/",
    },
}

func main() {
    cfg, err := ini.Load("my.ini")
    if err != nil {
        fmt.Printf("Fail to read file: %v", err)
        os.Exit(1)
    }

    v := cfg.Sections()
    for i := 0; i < len(v); i++ {
        if v[i].HasKey("name") {
            if v[i].Key("enabled").MustInt(1) == 1 {
                fmt.Printf("enabled yum repo %v\n", v[i].Name())
                for n, l := range newBaseUrl {
                    if v[i].Name() == string(n) && v[i].HasKey("baseurl") {
                        fmt.Printf("Try to replace in %v\n", v[i].Key("baseurl").Value())
                        fmt.Printf("New baseurl: %v\n", strings.ReplaceAll(v[i].Key("baseurl").Value(), l.oldURL, l.newURL))
                    }
                }
            } else {
                fmt.Printf("disabled yum repo %v\n", v[i].Name())
            }
        } else {
            fmt.Printf("Not a yum repo section %v\n", v[i].Name())
        }
    }
    return

    // Classic read of values, default section can be represented as empty string
    fmt.Println("App Mode:", cfg.Section("").Key("app_mode").String())
    fmt.Println("Data Path:", cfg.Section("paths").Key("data").String())

    // Let's do some candidate value limitation
    fmt.Println("Server Protocol:",
        cfg.Section("server").Key("protocol").In("http", []string{"http", "https"}))
    // Value read that is not in candidates will be discarded and fall back to given default value
    fmt.Println("Email Protocol:",
        cfg.Section("server").Key("protocol").In("smtp", []string{"imap", "smtp"}))

    // Try out auto-type conversion
    fmt.Printf("Port Number: (%[1]T) %[1]d\n", cfg.Section("server").Key("http_port").MustInt(9999))
    fmt.Printf("Enforce Domain: (%[1]T) %[1]v\n", cfg.Section("server").Key("enforce_domain").MustBool(false))
    
    // Now, make some changes and save it
    cfg.Section("").Key("app_mode").SetValue("production")
    cfg.SaveTo("my.ini.local")
}

