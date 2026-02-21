package di_test

import (
	"fmt"
	"strings"
	"testing"

	di "git.apihub24.de/admin/generic-di"
	"github.com/google/uuid"
)

func init() {
	di.Injectable(newTextService)
	di.Injectable(newMessageService)
	di.Injectable(newConfiguration)
	di.Injectable(newGreetingService)
	di.Injectable(newBasicOverridableService)
}

type (
	configuration  struct{}
	messageService struct {
		texts *textService
	}

	textService struct {
		config *configuration
		id     string
	}

	greetingService interface {
		Greeting() string
		TakeID() string
	}

	overridableService interface {
		GetInstanceID() string
		GetValue() string
	}

	basicOverridableService struct {
		id string
	}

	basicOverridableServiceMock struct {
		id string
	}
)

func newBasicOverridableService() overridableService {
	return &basicOverridableService{
		id: uuid.NewString(),
	}
}

func newBasicOverridableServiceMock() overridableService {
	return &basicOverridableServiceMock{
		id: uuid.NewString(),
	}
}

func newConfiguration() *configuration {
	return &configuration{}
}

func newTextService() *textService {
	return &textService{
		config: di.Inject[*configuration](),
		id:     uuid.NewString(),
	}
}

func newGreetingService() greetingService {
	return newTextService()
}

func newMessageService() *messageService {
	return &messageService{
		texts: di.Inject[*textService](),
	}
}

func (ctx *basicOverridableService) GetInstanceID() string {
	return ctx.id
}

func (ctx *basicOverridableService) GetValue() string {
	return "i am original"
}

func (ctx *basicOverridableServiceMock) GetInstanceID() string {
	return ctx.id
}

func (ctx *basicOverridableServiceMock) GetValue() string {
	return "i am mock"
}

func (ctx *configuration) GetUserName() string {
	return "Markus"
}

func (ctx *textService) Greeting() string {
	return fmt.Sprintf("Hello %s", ctx.config.GetUserName())
}

func (ctx *textService) GetID() string {
	return ctx.id
}

func (ctx *textService) TakeID() string {
	return ctx.id
}

func (ctx *messageService) GetTextServiceID() string {
	return ctx.texts.GetID()
}

func TestInject(t *testing.T) {
	// testMutex.Lock()
	// defer testMutex.Unlock()

	msg := newMessageService()
	println(msg.texts.Greeting())
	if msg.texts.Greeting() != "Hello Markus" {
		t.Errorf("expect greeting Hello Markus")
	}
}

func TestInject_Duplicate(t *testing.T) {
	// testMutex.Lock()
	// defer testMutex.Unlock()

	msg1 := newMessageService()
	msg2 := newMessageService()
	println(msg1.texts.Greeting())
	println(msg2.texts.Greeting())
	if msg1.texts.Greeting() != "Hello Markus" {
		t.Errorf("expect greeting Hello Markus")
	}
	if msg2.texts.Greeting() != "Hello Markus" {
		t.Errorf("expect greeting Hello Markus")
	}
	if msg1.GetTextServiceID() != msg2.GetTextServiceID() {
		t.Errorf("expect same instance of textService")
	}
}

func TestInject_Parallel(t *testing.T) {
	for i := 0; i < 20; i++ {
		go func() {
			println(di.Inject[*textService]().GetID())
		}()
	}
}

func TestInject_MultipleInstances(t *testing.T) {
	textServiceA := di.Inject[*textService]("a")
	textServiceB := di.Inject[*textService]("b")
	if textServiceA.GetID() == textServiceB.GetID() {
		t.Errorf("expect a seperate instance textServiceA and textServiceB but there was identical")
	}
}

func TestTryInterface(t *testing.T) {
	greeter := di.Inject[greetingService]()
	if greeter.Greeting() != "Hello Markus" {
		t.Errorf("expect greeting Hello Markus")
	}
}

func TestTryInterface_MultipleInstances(t *testing.T) {
	greeterA := di.Inject[greetingService]("a")
	greeterB := di.Inject[greetingService]("b")
	if greeterA.Greeting() != "Hello Markus" {
		t.Errorf("expect greeting Hello Markus")
	}
	if greeterB.Greeting() != "Hello Markus" {
		t.Errorf("expect greeting Hello Markus")
	}
	if greeterA.TakeID() == greeterB.TakeID() {
		t.Errorf("expect greetingA and greeterB are the same Instance")
	}
}

func TestOverwriteInjectable(t *testing.T) {
	basic := di.Inject[overridableService]()
	basicID := basic.GetInstanceID()
	if basic.GetValue() != "i am original" {
		t.Errorf("wrong service instance get")
	}
	di.Replace(newBasicOverridableServiceMock)
	basic = di.Inject[overridableService]()
	if basic.GetInstanceID() == basicID {
		t.Errorf("basic and newOne are the same instance")
	}
	if basic.GetValue() != "i am mock" {
		t.Errorf("service not overwritten")
	}
}

func TestOverwriteInjectableInstance(t *testing.T) {
	basic := di.Inject[overridableService]()
	basicID := basic.GetInstanceID()
	di.ReplaceInstance(newBasicOverridableServiceMock())
	basic = di.Inject[overridableService]()
	if basic.GetInstanceID() == basicID {
		t.Errorf("basic and newOne are the same instance")
	}
	if basic.GetValue() != "i am mock" {
		t.Errorf("service not overwritten")
	}
}

func TestDestroyMatching(t *testing.T) {
	_ = di.Inject[greetingService]("abc")
	_ = di.Inject[greetingService]("def")
	_ = di.Inject[greetingService]("abc_def")
	di.DestroyAllMatching(func(key string) bool { return strings.HasPrefix(key, "abc") })
}

func TestDestroy(t *testing.T) {
	_ = di.Inject[textService]("a")
	di.Destroy[textService]("a")
}
