package manager

import (
	"lyrid-sd/model"
	"sync"
)

type NodeManager struct {
	RouteMap map[string]model.Router
}

var instance *NodeManager
var once sync.Once

func GetInstance() *NodeManager {
	once.Do(func() {
		instance = &NodeManager{}
	})
	return instance
}

func (manager *NodeManager) Init() {
	manager.RouteMap = make(map[string]model.Router)
}

func (manager *NodeManager) Run() {

}

func (manager *NodeManager) Add(r model.Router) {
	manager.RouteMap[r.GetPort()] = r
}
