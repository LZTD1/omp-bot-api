package tests

//go:generate mockgen -destination=./mocks/repo_mock.go -package=mocks ozon-omp-api/internal/repo EventRepo
//go:generate mockgen -destination=./mocks/sender_mock.go -package=mocks ozon-omp-api/internal/sender EventSender
//go:generate mockgen -destination=./mocks/worker_pool.go -package=mocks ozon-omp-api/internal/repo WorkerPool

func Dummy() {
}
