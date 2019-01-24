package main

import "C"
import (
	"github.com/mrosset/util/file"
	"github.com/mrosset/util/json"
	"github.com/mrosset/via/pkg"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/gl"
	"log"
	"os"
	"path/filepath"
)

const (
	PACKAGE = "org.golang.todo.via"
	HOME    = "/data/data/org.golang.todo.via/files"
	SRCPATH = "/data/data/org.golang.todo.via/files/src/via"
	GITURL  = "https://github.com/mrosset/via"
)

var (
	configPath = filepath.Join(HOME, "config.json")
)

func init() {
	os.Setenv("HOME", HOME)
}

func getConfig() (*via.Config, error) {
	var (
		config = &via.Config{}
	)
	if !file.Exists(SRCPATH) {
		log.Printf("cloning %s -> %s", GITURL, SRCPATH)
		if err := via.Clone(SRCPATH, GITURL); err != nil {
			return nil, err
		}
	}
	log.Printf("reading %s", filepath.Join(SRCPATH, "plans/config.json"))
	if err := json.Read(config, filepath.Join(SRCPATH, "plans/config.json")); err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	app.Main(func(a app.App) {

		config, err := getConfig()
		if err != nil {
			log.Println(err)
		}

		log.Printf("Branch: %s", config.Branch)
		// if err := gurl.Download(HOME, "https://raw.githubusercontent.com/mrosset/plans/aarch64-via-linux-gnu-android/config.json"); err != nil {
		//	log(err.Error())
		// }
		if !file.Exists("/data/data/org.golang.todo.via/files") {
			log.Println("we don't have a data directory")
		}
		log.Println("started glxctx")
		var glctx gl.Context
		sz := size.Event{}
		for {
			select {
			case e := <-a.Events():
				switch e := a.Filter(e).(type) {
				case lifecycle.Event:
					glctx, _ = e.DrawContext.(gl.Context)
				case size.Event:
					sz = e
				case paint.Event:
					if glctx == nil {
						continue
					}
					onDraw(glctx, sz)
					a.Publish()
				}
			}
		}
	})
}

func onDraw(glctx gl.Context, sz size.Event) {
	glctx.Clear(gl.COLOR_BUFFER_BIT)
}
