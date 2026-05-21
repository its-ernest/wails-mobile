package helloworld

type HelloService struct{}

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (s *HelloService) SayHello(name string) (map[string]string, error) {
	if name == "" {
		name = "World"
	}

	return map[string]string{
		"message": "Hello, " + name + "!",
	}, nil
}
