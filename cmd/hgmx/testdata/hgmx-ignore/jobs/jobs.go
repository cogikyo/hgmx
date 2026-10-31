package jobs

import "os"

func Stop() {
	//hgmx:ignore exit/main-only the supervisor restarts the worker
	os.Exit(1)
}

func Halt() {
	os.Exit(2) //hgmx:ignore exit/main-only
}

func Wait() {
	//hgmx:ignore exit/main-only nothing exits here
}

//hgmx:ignore exit/anywhere typo in the rule name
func Run() {}
