package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// The `!java` command generates a random Java variable assignment
type JavaCommander struct{}

func (jc JavaCommander) Handle(n ApiNooter, msg Message) {
	n.NootMessage(genJava())
}

func genJava() string {
	inWords := []string{
		"Http", "Attribute", "Order", "Https", "Composite", "Invalid",
		"Supported", "Abstract", "Common", "Concrete", "Autowire", "Simple",
		"Aware", "Aspect", "Principal", "Driven", "Interruptible", "Batch",
		"Remote", "Stateless", "Session", "Based", "Meta", "Data", "Readable",
		"Reflective", "Xml", "Generic", "Interface", "Advisable", "Observable",
		"Identifiable", "Iterable", "Distributed", "Notification", "Failure",
		"Type", "Swift", "Rust", "Go", "Protocol", "Trait",
	}
	outWords := []string{
		"Handler", "Factory", "Wrapper", "Visitor", "Thread", "Pool",
		"Serializer", "Model", "Method", "Configuration", "Interceptor",
		"Exception", "Error", "Property", "Value", "Identifier",
		"Policy", "Container", "Info", "Descriptor", "Bean", "Singleton",
		"Parameter", "Adapter", "Bridge", "Decorator", "Facade", "Proxy",
		"Worker", "Interpreter", "Iterator", "Observer", "Request",
		"Transformer", "Interface", "State", "Template", "Task", "Resolver",
		"Definition", "Getter", "Setter", "Listener", "Processor", "Printer",
		"Prototype", "Protocol", "Trait", "Composer", "Event", "Helper", "Utils",
		"Exporter", "Importer", "Serializer", "Factory", "Callback", "Context",
		"Tests", "Annotation", "Service", "Factory", "Dispatcher", "Client",
		"Server", "Map", "List", "Collection", "Queue", "Manager", "Database",
		"Response", "Broadcaster", "Watcher", "Publisher", "Consumer", "Producer",
		"Factory",
	}

	rand.Seed(time.Now().UnixNano())

	// max halfopen interval
	maxIn := len(inWords)
	maxOut := len(outWords)

	in := inWords[rand.Intn(maxIn)]

	out := []string{}
	maxOutWords := 7
	for i := 0; i < maxOutWords; i++ {
		if rand.Intn(2) == 0 {
			out = append(out, outWords[rand.Intn(maxOut)])
		}
	}

	className := strings.Join(append([]string{in}, out...), "")
	varName := strings.Join(append([]string{strings.ToLower(in)}, out...), "")
	return fmt.Sprintf("%s %s = new %s();", className, varName, className)
}
