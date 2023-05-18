package frouter

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type TestData struct {
	Name string
}

type TestController struct {
	testRouter *Router[TestData]
}

func NewTestController() *TestController {
	controller := &TestController{testRouter: New[TestData]()}
	controller.testRouter.RegisterSync(Init, func(instance TestData) error {
		fmt.Println("[init]name:" + instance.Name)
		return nil
	}, nil)
	controller.testRouter.RegisterAsync(Init, func(instance TestData) error {
		fmt.Println("[init]name:" + instance.Name)
		return nil
	}, nil)
	controller.testRouter.RegisterPreSync(func(instance TestData) error {
		fmt.Println("pre sync")
		return nil
	})
	controller.testRouter.RegisterPreAsync(func(instance TestData) error {
		fmt.Println("pre async")
		return nil
	})
	return controller
}

func TestNewController(t *testing.T) {
	controller := New[TestData]()
	controller.RegisterSync(Init, func(instance TestData) error {
		fmt.Println("init sync")
		return nil
	}, nil)
}

func TestNewControllerRollbackSync(t *testing.T) {
	controller := New[TestData]()
	controller.RegisterSync(Add, func(instance TestData) error {
		fmt.Println("Add sync 1")
		return nil
	}, func(instance TestData) {
		fmt.Println("Add sync 1 rollback")
	})
	controller.RegisterSync(Add, func(instance TestData) error {
		fmt.Println("Add sync 2")
		return nil
	}, func(instance TestData) {
		fmt.Println("Add sync 2 rollback")
	})
	controller.RegisterSync(Add, func(instance TestData) error {
		fmt.Println("Add sync 3")
		return errors.New("a sync 3 failed")
	}, func(instance TestData) {
		fmt.Println("Add sync 3 rollback")
	})
	_ = controller.Sync(Add, TestData{Name: "lai fu"})
	time.Sleep(3 * time.Second)
}

func TestNewControllerRollbackAsync(t *testing.T) {
	controller := New[TestData]()
	controller.RegisterAsync(Update, func(instance TestData) error {
		fmt.Println("Update sync 1")
		return nil
	}, func(instance TestData) {
		fmt.Println("Update sync 1 rollback")
	})
	controller.RegisterAsync(Update, func(instance TestData) error {
		fmt.Println("Update sync 2")
		//return nil
		return errors.New("Update sync 2 failed")
	}, func(instance TestData) {
		fmt.Println("Update sync 2 rollback")
	})
	controller.RegisterAsync(Update, func(instance TestData) error {
		fmt.Println("Update sync 3")
		//return errors.New("b sync 3 failed")
		return nil
	}, func(instance TestData) {
		fmt.Println("Update sync 3 rollback")
	})
	controller.Async(Update, TestData{Name: "lai fu"})
	time.Sleep(3 * time.Second)
}

func TestController_Sync(t *testing.T) {
	controller := NewTestController()
	fmt.Println("start")
	_ = controller.testRouter.Sync(Init, TestData{Name: "lai fu"})
}

func TestController_Async(t *testing.T) {
	controller := NewTestController()
	fmt.Println("start")
	controller.testRouter.Async(Init, TestData{Name: "lai fu"})
	controller.testRouter.Stop()
	controller.testRouter.Async(Init, TestData{Name: "lai fu"})
	controller.testRouter.Start()
	controller.testRouter.Async(Init, TestData{Name: "lai fu"})
	time.Sleep(2 * time.Second)
}
