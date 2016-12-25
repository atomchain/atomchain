package params

var SeedString = "atomchain"

var genesisHash = ""

var MainnetBootnodes [1]string = [1]string{
	"http://44.44.47.33:5000/atomchain",
}

var TestnetBootnodes [1]string = [1]string{
	"http://44.44.47.33:5000/atomchain",
}

//ServerAddr start Local listen port
var ServerAddr = ":1666"
