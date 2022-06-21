package main

var guild string

func GuildLogin(name string) {
	_, rsp := Send(MethodGuild, map[string]any{"gld": name})
	if rsp.Error == nil {
		guild = name
		Write("gld", "Successfully connected to guild %s!", name)
	} else {
		Error("gld", "%s", *rsp.Error)
	}
}
