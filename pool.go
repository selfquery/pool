package pool

import "sync"

func CreatePool(jobs []Details, size int) pool {
	return pool{
		Jobs: jobs,
		Size: size,
		Work: make(chan Details),
	}
}

type Details interface {
	Process()
	Run(wait *sync.WaitGroup)
}

type pool struct {
	Jobs []Details
	Size int
	Work chan Details

	Group sync.WaitGroup
}

func (P *pool) Worker() {
	for job := range P.Work {
		job.Run(&P.Group)
	}
}

func (P *pool) Run() {
	for i := 0; i < P.Size; i++ {
		go P.Worker()
	}

	P.Group.Add(len(P.Jobs))
	for _, job := range P.Jobs {
		P.Work <- job
	}

	close(P.Work)
	P.Group.Wait()
}
