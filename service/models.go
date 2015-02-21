package tiroler

type Novel struct {
	Ncode    string    `param:"ncode"`
	Tcode    string    `param:"tcode"`
	Title    string    `param:"title"`
}

// Store all our greets in a big list in memory, because, let's be honest, who's
// actually going to use a service that only allows you to post 140-character
// messages?
var Novel = []Novel{
}

