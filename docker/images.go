package docker

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/clidockermgr/ui"
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

	var id = i.image.ID[7:19]

	var repo = ""
	var tag = ""
	if len(i.image.RepoTags) > 0 {
		repo = i.image.RepoTags[len(i.image.RepoTags)-1]
		if strings.Contains(repo, ":") {
			var i = strings.Index(repo, ":")
			tag = repo[i+1:]
			repo = repo[0:i]
		}
	}

	var durationHs = uint64(time.Since(time.Unix(i.image.Created, 0)).Hours())

	var durationStr = ""
	if durationHs > 24 {
		durationStr += fmt.Sprintf("%d days, ", durationHs/24)
		durationHs %= 24
	}
	durationStr += fmt.Sprintf("%d hs", durationHs)

	return fmt.Sprintf("%s %-60s %-20s %s", id, repo, tag, durationStr)
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
