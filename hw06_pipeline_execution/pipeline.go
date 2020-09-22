package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func checkDone(in In, done In) (out Out) {
	pipe := make(Bi)
	go func() {
		defer close(pipe)
		for input := range in {
			select {
			case <-done:
				return
			case pipe <- input:
			}
		}
	}()
	out = pipe
	return
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	inputChan := in

	for _, stage := range stages {
		inputChan = checkDone(inputChan, done)
		inputChan = stage(inputChan)
	}

	return inputChan
}
