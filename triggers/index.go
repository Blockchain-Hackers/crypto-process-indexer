package triggers

// "fmt"
// "github.com/blockchain-hackers/indexer"

type Trigger interface {
	run()
	// they should process events if they find one
	processEvent(Event)
	// EventName() string
}

type Event struct {
	EventName string
	Data      map[string]interface{}
}

var triggers = []Trigger{
	&ChainlinkPriceFeed{},
	&EthSepoliaIndexer{},
}

// run triggers
func Run() {
	for _, trigger := range triggers {
		go trigger.run()
	}
}
