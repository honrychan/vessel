package middleware

import (
	//"github.com/ngaut/log"
	//	"github.com/containerops/vessel/module"

	"gopkg.in/macaron.v1"
)

func InitLog(runmode, path string) {
	//	log := module.GetLog()

	//	log.SetLevelByString("info")

	if runmode == "dev" {
		//		log.SetLevelByString("debug")
	}

	//	log.SetOutputByName(path)
}

func logger(runmode string) macaron.Handler {
	//	log := module.GetLog()
	return func(ctx *macaron.Context) {
		if runmode == "dev" {
			//			log.Debug("")
			//			log.Debug("----------------------------------------------------------------------------------")
		}

		//		log.Infof("[%s] [%s]", ctx.Req.Method, ctx.Req.RequestURI)
		//		log.Infof("[Header] %v", ctx.Req.Header)
	}
}
