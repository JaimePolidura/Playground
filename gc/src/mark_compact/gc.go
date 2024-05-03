package markcompact

func StartGC() {
	stopThreads()
	mark()
	compact()
	startThreads()
}

func compact() {

}

func stopThreads() {
	//TDOO
}

func startThreads() {
	//TODO
}
