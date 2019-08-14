package routing

import "net/http"

type Route struct {
	Method      string
	Path        string
	Name        string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		Method:      "GET",
		Path:        "/start",
		Name:        "start",
		HandlerFunc: Start,
	},

	Route{
		Method:      "POST",
		Path:        "/peers",
		Name:        "peers",
		HandlerFunc: PeersAlive,
	},
	Route{
		Method:      "GET",
		Path:        "/peers",
		Name:        "peers",
		HandlerFunc: PeersAlive,
	},

	Route{
		Method:      "GET",
		Path:        "/chat",
		Name:        "chat",
		HandlerFunc: Chat,
	},
	Route{
		Method:      "POST",
		Path:        "/chat",
		Name:        "chat",
		HandlerFunc: Chat,
	},
	Route{
		Method:      "POST",
		Path:        "/chat/recv",
		Name:        "chatRecv",
		HandlerFunc: ChatBeatRecv,
	},

	Route{
		Method:      "GET",
		Path:        "/hello",
		Name:        "hello",
		HandlerFunc: Hello,
	},
	Route{
		Method:      "GET",
		Path:        "/wit",
		Name:        "wit/{textInput}",
		HandlerFunc: Wit,
	},
}
