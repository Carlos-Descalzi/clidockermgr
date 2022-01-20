package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	ui "github.com/clidockermgr/ui"
	util "github.com/clidockermgr/util"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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
	dockerClient *client.Client
	items        []types.Container
	stats        map[string]types.Stats
}

func ParseStats(statsBody []byte) *types.Stats {
	var stats types.Stats
	json.Unmarshal(statsBody, &stats)
	return &stats
}

func Update(model ContainerListModel) {
	for true {
		time.Sleep(100)
		for c := range model.items {
			stats, err := model.dockerClient.ContainerStats(context.Background(), model.items[c].ID, false)
			if err != nil {
				log.Printf("Error getting stats %s", err)
			} else {
				var buf = bufio.NewReader(stats.Body)
				var result, _ = buf.ReadBytes(byte(0))
				var stats = ParseStats(result)
				model.stats[model.items[c].ID] = *stats
			}
		}
		model.NotifyChanged()
	}
}

func ContainerListModelNew(client *client.Client) *ContainerListModel {

	var model = ContainerListModel{dockerClient: client, stats: make(map[string]types.Stats)}
	model.Init()
	model.Update()
	go Update(model)
	return &model
}

func (m *ContainerListModel) Update() {
	items, err := m.dockerClient.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}
	m.items = items
}

func (m ContainerListModel) ItemCount() int {
	return len(m.items)
}

func (m ContainerListModel) Item(index int) ui.ListItem {
	var container = m.items[index]
	var stats = m.stats[container.ID]
	return &ContainerListModelItem{m.items[index], int(stats.MemoryStats.Usage), int(stats.MemoryStats.Limit)}
}
