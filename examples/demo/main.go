package main

import (
	"flag"
	"github.com/garyburd/twister/web"
	"github.com/garyburd/twister/server"
	"template"
	"net"
	"log"
)

func homeHandler(req *web.Request) {
	homeTempl.Execute(
		req.Respond(web.StatusOK, web.HeaderContentType, "text/html"),
		req)
}

func main() {
	flag.Parse()
	h := web.SetErrorHandler(coreErrorHandler,
		web.ProcessForm(10000, true, web.NewHostRouter(nil).
			Register("www.example.com", web.NewRouter().
			Register("/", "GET", homeHandler).
			Register("/core/file", "GET", web.FileHandler("static/file.txt")).
			Register("/static/<path:.*>", "GET", web.DirectoryHandler("static/")).
			Register("/chat", "GET", chatFrameHandler).
			Register("/chat/ws", "GET", chatWsHandler).
			Register("/mp", "GET", mpGetHandler, "POST", mpPostHandler).
			Register("/core/", "GET", coreHandler).
			Register("/core/a/<a>/", "GET", coreHandler).
			Register("/core/b/<b>/c/<c>", "GET", coreHandler).
			Register("/core/c", "POST", coreHandler))))

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Listen", err)
		return
	}
	defer listener.Close()
	err = (&server.Server{Listener: listener, Handler: h, Logger: server.VerboseLogger}).Serve()
	if err != nil {
		log.Fatal("Server", err)
	}
}

var homeTempl = template.MustParse(homeStr, template.FormatterMap{"": template.HTMLFormatter})

const homeStr = `
<html>
<head>
</head>
<body>
<ul>
<li><a href="/core">Core functionality</a>
<li><a href="/chat">Chat using WebSockets</a>
<li><a href="/mp">Multipart Form</a>
</ul>
</body>
</html>`