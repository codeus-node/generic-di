# generic-di

Go Dependency Injection with Generics

## Example

### Struct

configuration.go

```go
package main

import "github.com/codeus-node/generic-di"

func init() {
	// register the Struct Constructor Function for DI
	di.Injectable(NewConfiguration)
}

type Configuration struct {
	UserName string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		UserName: "Markus",
	}
}
```

greeter.go

```go
package main

import (
	"fmt"
	"github.com/codeus-node/generic-di"
)

func init() {
	di.Injectable(NewGreeter)
}

type Greeter struct {
	config *Configuration
}

func NewGreeter() *Greeter {
	return &Greeter{
		// here was the Configuration from configuration.go injected
		config: di.Inject[*Configuration](),
	}
}

func (ctx *Greeter) Greet() string {
	return fmt.Sprintf("Hello, %s", ctx.config.UserName)
}
```

message_service.go

```go
package main

import "github.com/codeus-node/generic-di"

func init() {
	di.Injectable(NewMessageService)
}

type MessageService struct {
	greeter *Greeter
}

func NewMessageService() *MessageService {
	return &MessageService{
		// here was the Greeter from greeter.go injected
		greeter: di.Inject[*Greeter](),
	}
}

func (ctx *MessageService) Welcome() string {
	return ctx.greeter.Greet()
}
```

main.go

```go
package main

import di "github.com/codeus-node/generic-di"

func main() {
	msgService := di.Inject[*MessageService]()
	// prints the message "Hello, Markus"
	println(msgService.Welcome())
}
```

### Interface

message_service.go

```go
package main

import "github.com/codeus-node/generic-di"

func init() {
	di.Injectable(newMessageService)
}

type IMessageService interface {
	Welcome() string
}

type messageService struct {
	greeter *Greeter
}

func NewMessageService() IMessageService {
	return &messageService{
		// here was the Greeter from greeter.go injected
		greeter: di.Inject[*Greeter](),
	}
}

func (ctx *messageService) Welcome() string {
	return ctx.greeter.Greet()
}
```

main.go

```go
package main

import di "github.com/codeus-node/generic-di"

func main() {
	msgService := di.Inject[IMessageService]()
	// prints the message "Hello, Markus"
	println(msgService.Welcome())
}
```

## Replace Instance

services/message_service.go

```go
package services

import "github.com/codeus-node/generic-di"

func init() {
	di.Injectable(newMessageService)
}

type IMessageService interface {
	Welcome() string
}

type messageService struct {
	greeter *Greeter
}

func NewMessageService() IMessageService {
	return &messageService{
		// here was the Greeter from greeter.go injected
		greeter: di.Inject[*Greeter](),
	}
}

func (ctx *messageService) Welcome() string {
	return ctx.greeter.Greet()
}
```

main.go

```go
package main

import di "github.com/codeus-node/generic-di"

func main() {
	msgService := di.Inject[services.IMessageService]()
	// prints the message "Hello, Markus"
	println(msgService.Welcome())
}
```

message_service_test.go

```go
package services_test

func init() {
	// replace the instance of IMessageService with the Mock
	di.Replace(newMessageServiceMock)
}

type messageServiceMock struct {}

func newMessageServiceMock() services.IMessageService {
	return &messageServiceMock{}
}

func (svc *messageServiceMock) Welcome() string {
	return "Hello, Mock"
}

func TestMessageServiceMocking(t *testing.T) {
	service := di.Inject[services.IMessageService]()
	if service.Welcome() != "Hello, Mock" {
		t.Errorf("expect Hello, Mock but get %s", service.Welcome())
	}
}
```
