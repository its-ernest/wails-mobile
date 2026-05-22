package wailsmobile

type AppService struct{}

func NewService() *AppService {
	return &AppService{}
}

func (s *AppService) SayHello(name string) (map[string]string, error) {
	if name == "" {
		name = "World"
	}

	return map[string]string{
		"message": "Hello, " + name + "!",
	}, nil
}
