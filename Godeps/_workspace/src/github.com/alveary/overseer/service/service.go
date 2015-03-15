package service

// Service ...
type Service struct {
	Name  string `json:"name"`
	Root  string `json:"root"`
	Alive string `json:"alive"`
	Fails int    `json:"fails"`
}

func (service *Service) AddFailure() {
	service.Fails = service.Fails + 1
}
