package workers

// Worker is a type that will do an asynchronous job every X time.
type Worker interface {
	// Run starts the worker.
	Run()
}

// Workers is a collection of workers.
type Workers []Worker

// Run starts all workers in their own goroutines.
func (w Workers) Run() {
	// TODO: Implement graceful shutdown
	for _, worker := range w {
		go worker.Run()
	}
}
