package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	docker "github.com/clidockermgr/docker"
	"github.com/clidockermgr/input"
	ui "github.com/clidockermgr/ui"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/eiannone/keyboard"
)

func ShowTextPopup(app *ui.Application, title string, text string) {

	maxWidth, maxHeight := ui.ScreenSize()

	popupWidth := uint8(float32(maxWidth) * 0.75)
	popupHeight := uint8(float32(maxHeight) * 0.8)

	textView := ui.TextViewNew(text)
	container := ui.TitledContainerNew(title, textView, true)
	container.SetRect(ui.RectNew((maxWidth-popupWidth)/2, (maxHeight-popupHeight)/2, popupWidth, popupHeight))

	app.ShowPopup(container)
}

func ShowContainerInspect(app *ui.Application, inspect types.ContainerJSON) {

	result, err := json.MarshalIndent(inspect, "", "    ")
	if err == nil {
		var strResult = string(result)
		ShowTextPopup(app, "Container Details", strResult)
	}
}

func ShowImageInspect(app *ui.Application, inspect types.ImageInspect) {
	result, err := json.MarshalIndent(inspect, "", "    ")

	if err == nil {
		var strResult = string(result)
		ShowTextPopup(app, "Image Details", strResult)
	}
}

func SetupLog() {
	var logfile, err = os.OpenFile("dockermgr.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}
	log.SetOutput(logfile)

}

func BuildContainersView(app *ui.Application, client *client.Client, width uint8, height uint8) {
	var containerList = ui.ListNew()
	containerList.SetModel(docker.ContainerListModelNew(client))

	containerList.AddKeyHandler(input.KeyInputKey(keyboard.KeyCtrlV), func(input.KeyInput) {
		var item = containerList.SelectedItem().Value().(*types.Container)
		inspect, err := client.ContainerInspect(context.Background(), item.ID)
		if err == nil {
			ShowContainerInspect(app, inspect)
		}
	})
	containerList.AddKeyHandler(input.KeyInputKey(keyboard.KeyCtrlK), func(input.KeyInput) {
		var item = containerList.SelectedItem().Value().(*types.Container)
		client.ContainerKill(context.Background(), item.ID, "9")
	})
	containerList.AddKeyHandler(input.KeyInputKey(keyboard.KeyDelete), func(input.KeyInput) {
		var item = containerList.SelectedItem().Value().(*types.Container)
		client.ContainerDiff(context.Background(), item.ID)
	})

	var titledContainer1 = ui.TitledContainerNew("Containers", containerList, false)
	titledContainer1.SetRect(ui.RectNew(1, 1, width, height))
	app.Add(titledContainer1)

}

func BuildImagesView(app *ui.Application, client *client.Client, width uint8, height uint8) {
	var imageList = ui.ListNew()
	imageList.SetModel(docker.ImagesListModelNew(client))

	imageList.AddKeyHandler(input.KeyInputKey(keyboard.KeyCtrlV), func(input.KeyInput) {
		var item = imageList.SelectedItem().Value().(*types.ImageSummary)
		inspect, _, err := client.ImageInspectWithRaw(context.Background(), item.ID)
		if err == nil {
			ShowImageInspect(app, inspect)
		}
	})

	imageList.AddKeyHandler(input.KeyInputKey(keyboard.KeyDelete), func(input.KeyInput) {
		var item = imageList.SelectedItem().Value().(*types.ImageSummary)
		client.ImageRemove(context.Background(), item.ID, types.ImageRemoveOptions{})
	})

	var titledContainer2 = ui.TitledContainerNew("Images", imageList, false)
	titledContainer2.SetRect(ui.RectNew(1, height+1, width, height-1))
	app.Add(titledContainer2)
}

func main() {

	SetupLog()

	client, err1 := client.NewClientWithOpts(client.FromEnv)

	if err1 != nil {
		panic(err1)
	}
	maxWidth, maxHeight := ui.ScreenSize()
	areaHeight := maxHeight / 2

	var app = ui.ApplicationNew()

	BuildContainersView(app, client, maxWidth, areaHeight)
	BuildImagesView(app, client, maxWidth, areaHeight)

	app.Loop()

}
