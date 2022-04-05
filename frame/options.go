package frame

import (
        "os"
        "fmt"
        "github.com/hlhv/scribe"
        // TODO: write custom implementation
        "github.com/akamensky/argparse"
)

var options struct {
        logLevel  scribe.LogLevel
}

func parseArgs (conf *Conf) {
        parser := argparse.NewParser ("", conf.Description)
        logLevel := parser.Selector ("l", "log-level", []string {
                "debug",
                "normal",
                "error",
                "none",
        }, &argparse.Options {
                Required: false,
                Default:  "normal",
                Help:     "The amount of logs to produce. Debug prints " +
                          "everything, and none prints nothing",
        })

        err := parser.Parse(os.Args)
        if err != nil {
                fmt.Print(parser.Usage(err))
                os.Exit(1)
        }

        switch *logLevel {
                case "debug":  options.logLevel = scribe.LogLevelDebug;  break
                default:
                case "normal": options.logLevel = scribe.LogLevelNormal; break
                case "error":  options.logLevel = scribe.LogLevelError;  break
                case "none":   options.logLevel = scribe.LogLevelNone;   break
        }
}
