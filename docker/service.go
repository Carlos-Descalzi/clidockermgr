package docker

import (
	"container/list"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/clidockermgr/util"
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
	diskUsage int64
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
	diskUsage        map[string]int64
}

func ServiceHandlerNew(client *client.Client) *ServiceHandler {
	handler := ServiceHandler{
		client:    client,
		active:    true,
		listeners: list.New(),
		diskUsage: make(map[string]int64),
	}

	go handler.UpdateContainers()
	go handler.UpdateImages()
	go handler.GetSystemStats()
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
	for i := range s.containers {
		val, ok := s.diskUsage[s.containers[i].container.ID]
		if ok {
			s.containers[i].diskUsage = val
		}
	}
	return s.containers
}

func (s *ServiceHandler) Images() []types.ImageSummary {
	return s.images
}

func (s *ServiceHandler) UpdateContainers() {
	for s.active {
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

func (s *ServiceHandler) UpdateImages() {
	for s.active {
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

func (s *ServiceHandler) GetSystemStats() {
	for s.active {
		diskUsage, err := s.client.DiskUsage(context.Background())

		if err != nil {
			log.Printf("Error getting disk usage: %s", err)
		} else {
			for i := range diskUsage.Containers {
				s.diskUsage[diskUsage.Containers[i].ID] = diskUsage.Containers[i].SizeRw
			}
		}
		time.Sleep(time.Second)
	}
}

func (s *ServiceHandler) DoUpdateContainers(containers []types.Container) {

	var summaries []ContainerSummary
	var wg sync.WaitGroup

	wg.Add(len(containers))

	for i := range containers {
		summaries = append(summaries, ContainerSummary{container: containers[i]})

		go func(i int) {
			defer wg.Done()
			stats, err := s.client.ContainerStats(context.Background(), containers[i].ID, false)
			if err == nil {
				summaries[i].stats = *util.ParseStatsBody(stats.Body)
			} else {
				log.Printf("Error getting stats for container %s %s", containers[i].ID, err)
			}
		}(i)

	}
	wg.Wait()
	s.containers = summaries
}

func (s *ServiceHandler) RemoveImage(imageId string) {
	_, err := s.client.ImageRemove(context.Background(), imageId, types.ImageRemoveOptions{})

	if err != nil {
		log.Print("Error removing image", err)
	}
}

func (s *ServiceHandler) RemoveContainer(containerId string) {
	err := s.client.ContainerRemove(context.Background(), containerId, types.ContainerRemoveOptions{})

	if err != nil {
		log.Print("Error removing container", err)
	}
}

func (s *ServiceHandler) KillContainer(containerId string) {
	err := s.client.ContainerKill(context.Background(), containerId, "9")

	if err != nil {
		log.Print("Error killing container", err)
	}
}

func (s *ServiceHandler) Logs(containerId string) string {
	reader, err := s.client.ContainerLogs(context.Background(), containerId, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})

	if err == nil {
		result, err := ioutil.ReadAll(reader)

		if err == nil {
			return string(result)
		}
	} else {
		log.Print("Error getting container logs", err)
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
	} else {
		log.Print("Error inspecting image", err)
	}
	return ""
}

func (s *ServiceHandler) InspectContainerRaw(containerId string) types.ContainerJSON {
	inspect, err := s.client.ContainerInspect(context.Background(), containerId)

	if err != nil {
		log.Print("Error inspecting container", err)
	}
	return inspect
}

func (s *ServiceHandler) InspectContainer(containerId string) string {
	inspect, err := s.client.ContainerInspect(context.Background(), containerId)
	if err == nil {
		result, err := json.MarshalIndent(inspect, "", "    ")
		if err == nil {
			return string(result)
		}
	} else {
		log.Print("Error inspecting container", err)
	}
	return ""
}
