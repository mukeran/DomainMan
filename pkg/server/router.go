package server

import "DomainMan/pkg/server/handlers"

func registerRouter() {
	handlers.AccessTokenHandler{}.Register(GIN.Group("access_token"))
	handlers.ConfigHandler{}.Register(GIN.Group("config"))
	handlers.DomainHandler{}.Register(GIN.Group("domain"))
	handlers.RegistrarHandler{}.Register(GIN.Group("registrar"))
	handlers.SuffixHandler{}.Register(GIN.Group("suffix"))
	handlers.WhoisHandler{}.Register(GIN.Group("whois"))
}
