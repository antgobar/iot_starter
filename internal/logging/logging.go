package logging

import "log"

func SetUp() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
