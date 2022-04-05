package frame

import (
        "github.com/hlhv/scribe"
        "github.com/hlhv/cellkit/client"
)

type Conf struct {
        Description string
        Run func(*client.Leash)
}

func Be (conf *Conf) {
        parseArgs(conf)
        leash := client.NewLeash()
        go conf.Run(leash)
        loop()
}

func loop () {
        for {
                scribe.ListenOnce()
        }
}
