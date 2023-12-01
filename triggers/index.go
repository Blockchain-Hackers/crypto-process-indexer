package triggers

// "fmt"
// "github.com/blockchain-hackers/indexer"

type Trigger interface {
	run()
}

var triggers = []Trigger{
	&EthSepoliaIndexer{},
}

// run triggers
func Run() {
	for _, trigger := range triggers {
		trigger.run()
	}
}
