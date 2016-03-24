package workers

type Worker interface {
	Run()
}

type Workers []Worker

func (w Workers) Run() {
	// TODO: Implement graceful shutdown
	for _, worker := range w {
		go worker.Run()
	}
}
