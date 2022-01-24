package docker

import (
	"fmt"
	"log"

	"github.com/clidockermgr/ui"
	"github.com/clidockermgr/util"

	"github.com/docker/docker/api/types"
)

const (
	ContainerListModelOnlyActive = 1
)

type ContainerListModelItem struct {
	container types.Container
	usedMem   int
	maxMem    int
}

func (c ContainerListModelItem) Value() interface{} {
	return &c.container
}

func (i ContainerListModelItem) String() string {

	var image = i.container.Image

	if image[0:6] == "sha256" {
		image = image[7:19]
	}

	if len(image) > 40 {
		image = "..." + image[len(image)-37:]
	}

	var command = i.container.Command

	if len(command) > 30 {
		command = "..." + command[len(command)-27:]
	}

	var status = i.container.Status

	var memString = fmt.Sprintf("%s / %s", util.FormatMemory(i.usedMem), util.FormatMemory(i.maxMem))

	return fmt.Sprintf("%s %-40s %-30s %-30s %s", i.container.ID[0:12], image, command, status, memString)
}

type ContainerListModel struct {
	ui.BaseListModel
	dockerClient *ServiceHandler
	items        []ContainerSummary
	onlyActive   bool
	active       bool
}

func ContainerListModelNew(client *ServiceHandler) *ContainerListModel {

	var model = ContainerListModel{dockerClient: client, active: true}
	model.Init()
	model.Update()
	client.AddListener(&model)
	return &model
}

func (m *ContainerListModel) SetProperty(property int, value interface{}) {
	switch property {
	case ContainerListModelOnlyActive:
		m.onlyActive = !m.onlyActive
		m.Update()
	}
}

func (m *ContainerListModel) Update() {
	m.items = m.dockerClient.Containers()
	m.NotifyChanged()
}

func (m ContainerListModel) ItemCount() int {
	return len(m.items)
}

func (m ContainerListModel) Item(index int) ui.ListItem {

	return &ContainerListModelItem{
		m.items[index].container,
		int(m.items[index].stats.MemoryStats.Usage),
		int(m.items[index].stats.MemoryStats.Limit)}
}
func (m *ContainerListModel) ImagesUpdated() {
}
func (m *ContainerListModel) ContainersUpdated() {
	m.Update()
	log.Printf("Model changed %d", len(m.items))
	m.NotifyChanged()
}
