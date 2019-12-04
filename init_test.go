package pgo2

import (
    "github.com/pinguo/pgo2/iface"
    "github.com/pinguo/pgo2/logs"
)

type mockTarget struct {
    logs.Target
}

func (m *mockTarget) Process(item *logs.LogItem) {
    // fmt.Println("Process")
}

// Flush flush log to stdout
func (m *mockTarget) Flush(final bool) {
    //fmt.Println("Flush")
}

type mockPlugin struct {
}

func (t *mockPlugin) HandleRequest(ctx iface.IContext) {

}

type mockController struct {
    Controller
}

func (m *mockController) ActionIndex() {

}

func (m *mockController) ActionInfo() {

}
