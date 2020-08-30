package usecase

import "peterparada.com/online-bookmarks/domain"

type pingUsecase struct {
}

func NewPingUsecase() domain.PingUsecase {
	return &pingUsecase{}
}

func (p *pingUsecase) Get() string {
	return "PONG"
}
