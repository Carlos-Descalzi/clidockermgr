package docker

import (
	"container/list"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"

	util "github.com/clidockermgr/util"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ServiceListener interface {
	ImagesUpdated()
	ContainersUpdated()
}

type ContainerSummary struct {
	container types.Container
	stats     types.Stats
}

type StatData struct {
	id    string
	stats types.Stats
}

type ServiceHandler struct {
	client           *client.Client
	active           bool
	activeContainers bool
	images           []types.ImageSummary
	containers       []ContainerSummary
	listeners        *list.List
}

func ServiceHandlerNew(client *client.Client) *ServiceHandler {
	handler := ServiceHandler{client: client, active: true, listeners: list.New()}
	go handler.UpdateContainers()
	go handler.UpdateImages()
	return &handler
}

func (s *ServiceHandler) AddListener(listener interface{}) {
	s.listeners.PushBack(listener)
}

func (s *ServiceHandler) NotifyContainersUpdated() {
	for v := s.listeners.Front(); v != nil; v = v.Next() {
		v.Value.(ServiceListener).ContainersUpdated()
	}
}

func (s *ServiceHandler) NotifyImagesUpdated() {
	for v := s.listeners.Front(); v != nil; v = v.Next() {
		v.Value.(ServiceListener).ImagesUpdated()
	}
}

func (s *ServiceHandler) Containers() []ContainerSummary {
	return s.containers
}

func (s *ServiceHandler) Images() []types.ImageSummary {
	return s.images
}

func (s *ServiceHandler) UpdateContainers() {
	for s.active {
		log.Print("Getting containers")

		containers, err := s.client.ContainerList(context.Background(), types.ContainerListOptions{All: !s.activeContainers})

		if err != nil {
			log.Print(err)
		} else {
			s.DoUpdateContainers(containers)
		}
		s.NotifyContainersUpdated()
		time.Sleep(time.Second)
	}
}

func (s *ServiceHandler) DoUpdateContainers(containers []types.Container) {
	var summaries []ContainerSummary

	var wg sync.WaitGroup
	wg.Add(len(containers))

	for i := range containers {
		summaries = append(summaries, ContainerSummary{})
		summaries[i].container = containers[i]

		go func(i int) {
			defer wg.Done()
			stats, err := s.client.ContainerStats(context.Background(), containers[i].ID, false)
			if err != nil {
				summaries[i].stats = *util.ParseStatsBody(stats.Body)
			}
		}(i)

	}
	wg.Wait()
	s.containers = summaries
}

func (s *ServiceHandler) UpdateImages() {
	for s.active {
		log.Print("Getting images")

		images, err := s.client.ImageList(context.Background(), types.ImageListOptions{})

		if err != nil {
			log.Printf("Error getting images: %s", err)
		} else {
			s.images = images
		}
		s.NotifyImagesUpdated()
		time.Sleep(time.Second)
	}
}

func (s *ServiceHandler) RemoveImage(imageId string) {
	s.client.ImageRemove(context.Background(), imageId, types.ImageRemoveOptions{})
}

func (s *ServiceHandler) RemoveContainer(containerId string) {
	s.client.ContainerRemove(context.Background(), containerId, types.ContainerRemoveOptions{})
}

func (s *ServiceHandler) KillContainer(containerId string) {
	s.client.ContainerKill(context.Background(), containerId, "9")
}

func (s *ServiceHandler) Logs(containerId string) string {
	reader, err := s.client.ContainerLogs(context.Background(), containerId, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})

	if err == nil {
		result, err := ioutil.ReadAll(reader)

		if err == nil {
			return string(result)
		}
	}
	return ""
}

func (s *ServiceHandler) InspectImage(imageId string) string {
	inspect, _, err := s.client.ImageInspectWithRaw(context.Background(), imageId)
	if err == nil {
		result, err := json.MarshalIndent(inspect, "", "    ")

		if err == nil {
			return string(result)
		}
	}
	return ""
}

func (s *ServiceHandler) InspectContainer(containerId string) string {
	inspect, err := s.client.ContainerInspect(context.Background(), containerId)
	if err == nil {
		result, err := json.MarshalIndent(inspect, "", "    ")
		if err == nil {
			return string(result)
		}
	}
	return ""
}
