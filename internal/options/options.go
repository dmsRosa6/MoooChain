package options

import (
	"log"
	"os"
	"strconv"
)

type Options struct {
	DebugChain bool
	CleanupChain bool
	log *log.Logger
}

func InitOptions(log *log.Logger) *Options{
	debugChain := false
	cleanup := false

	val := os.Getenv("DEBUG_CHAIN")
	
	if val != "" {
		convertedVal, err := strconv.ParseBool(val)

		if err != nil {
			log.Printf("Invalid DEBUG_CAHIN value %q, defaulting to FALSE", val)
		} else {
			debugChain = convertedVal
		}

	}

	val = os.Getenv("CLEANUP_DB")
	if val != "" {
		convertedVal, err := strconv.ParseBool(val)

		if err != nil {
			log.Printf("Invalid CLEANUP_DB value %q, defaulting to FALSE", val)
		} else {
			cleanup = convertedVal
		}

	}
	return &Options{
		DebugChain: debugChain,
		CleanupChain: cleanup,
		log: log,
	}
}

func (o Options) Print(){
	log.Printf("\n*****Options*****\n\nDebug Chain: %t\n\n*****************\n\n", o.DebugChain)
}