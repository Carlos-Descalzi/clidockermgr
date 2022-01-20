package docker

import (
	"context"
	"log"

	ui "github.com/clidockermgr/ui"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ImageItem struct {
	image types.ImageSummary
}

func (i ImageItem) Value() interface{} {
	return &i.image
}

func (i ImageItem) String() string {
	if len(i.image.RepoTags) > 0 {
		return i.image.RepoTags[len(i.image.RepoTags)-1]
	}
	return i.image.ID[0:12]
}

type ImagesListModel struct {
	ui.BaseListModel
	dockerClient *client.Client
	items        []types.ImageSummary
}

func ImagesListModelNew(dockerClient *client.Client) *ImagesListModel {
	var model = ImagesListModel{dockerClient: dockerClient}
	model.Init()
	model.Update()
	return &model
}

func (m *ImagesListModel) Update() {
	items, err := m.dockerClient.ImageList(context.Background(), types.ImageListOptions{})

	if err != nil {
		log.Printf("Error getting images: %s", err)
	}

	m.items = items
}

func (m ImagesListModel) ItemCount() int {
	return len(m.items)
}

func (m ImagesListModel) Item(index int) ui.ListItem {
	return &ImageItem{m.items[index]}
}
