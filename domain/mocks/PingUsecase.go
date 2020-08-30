package mocks

type PingUsecase struct {
}

func (p *PingUsecase) Get() string {
	return "PONG"
}
