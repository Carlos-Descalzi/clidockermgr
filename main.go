package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/clidockermgr/docker"
	"github.com/clidockermgr/input"
	"github.com/clidockermgr/ui"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/eiannone/keyboard"
)

const HelpText = `Keys:

    tab: Switch focus between UI elements
    ESC: Closes active popup, or exits the application
    Container view:
        v: Displays container information
        k: Kills a container
        delete: Deletes a container
        s: Opens a shell in a container
    Images view:
        v: Displays image information
        delete: Deletes an image
`

func ShowTextPopup(app *ui.Application, title string, text string) {

	maxWidth, maxHeight := ui.ScreenSize()

	popupWidth := uint8(float32(maxWidth) * 0.75)
	popupHeight := uint8(float32(maxHeight) * 0.8)

	textView := ui.TextViewNew(text)
	container := ui.TitledContainerNew(title, textView, true)
	container.SetRect(ui.RectNew((maxWidth-popupWidth)/2, (maxHeight-popupHeight)/2, popupWidth, popupHeight))

	app.ShowPopup(container)
}

func ShowContainerInspect(app *ui.Application, client *client.Client, containerId string) {
	inspect, err := client.ContainerInspect(context.Background(), containerId)
	if err == nil {
		result, err := json.MarshalIndent(inspect, "", "    ")
		if err == nil {
			var strResult = string(result)
			ShowTextPopup(app, "Container Details", strResult)
		}
	}
}

func ShowHelp(app *ui.Application) {
	ShowTextPopup(app, "Help", HelpText)
}

func RunCommand(command string, args ...string) {
	var cmd = exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ui.CursorOn()
	ui.ClearScreen()
	cmd.Run()
	ui.CursorOff()
}

func DoExecContainer(containerId string, command string) {
	RunCommand("docker", "exec", "-it", containerId, command)
}

func ExecShell(containerId string) {
	DoExecContainer(containerId, "sh")
}

func ExecBashShell(containerId string) {
	DoExecContainer(containerId, "bash")
}

func ShowLogs(app *ui.Application, client *client.Client, containerId string) {

	var reader, err = client.ContainerLogs(context.Background(), containerId, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})

	if err == nil {
		var result, err2 = ioutil.ReadAll(reader)

		if err2 == nil {
			ShowTextPopup(app, "Logs", string(result))
		} else {
			log.Print(err2)
		}
	} else {
		log.Print(err)
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

	containerList.AddKeyHandler(input.KeyInputChar('v'), func(input.KeyInput) {
		var item = containerList.SelectedItem().Value().(*types.Container)
		ShowContainerInspect(app, client, item.ID)
	})
	containerList.AddKeyHandler(input.KeyInputChar('k'), func(input.KeyInput) {
		var item = containerList.SelectedItem().Value().(*types.Container)
		client.ContainerKill(context.Background(), item.ID, "9")
		containerList.Update()
	})
	containerList.AddKeyHandler(input.KeyInputChar('d'), func(input.KeyInput) {
		var item = containerList.SelectedItem().Value().(*types.Container)
		client.ContainerRemove(context.Background(), item.ID, types.ContainerRemoveOptions{})
		containerList.Update()
	})
	containerList.AddKeyHandler(input.KeyInputChar('s'), func(input.KeyInput) {
		var item = containerList.SelectedItem().Value().(*types.Container)
		ExecShell(item.ID)
	})
	containerList.AddKeyHandler(input.KeyInputChar('b'), func(input.KeyInput) {
		var item = containerList.SelectedItem().Value().(*types.Container)
		ExecBashShell(item.ID)
	})
	containerList.AddKeyHandler(input.KeyInputChar('l'), func(input.KeyInput) {
		var item = containerList.SelectedItem().Value().(*types.Container)
		ShowLogs(app, client, item.ID)
	})
	containerList.AddKeyHandler(input.KeyInputKey(keyboard.KeyDelete), func(input.KeyInput) {
		var item = containerList.SelectedItem().Value().(*types.Container)
		client.ContainerRemove(context.Background(), item.ID, types.ContainerRemoveOptions{})
	})
	containerList.AddKeyHandler(input.KeyInputChar('a'), func(input.KeyInput) {
		containerList.Model.SetProperty(docker.ContainerListModelOnlyActive, nil)
	})
	containerList.AddKeyHandler(input.KeyInputChar('h'), func(input.KeyInput) {
		ShowHelp(app)
	})

	var titledContainer1 = ui.TitledContainerNew("Containers", containerList, false)
	titledContainer1.SetRect(ui.RectNew(1, 1, width, height))
	app.Add(titledContainer1)

}

func ShowImageInspect(app *ui.Application, client *client.Client, imageId string) {

	inspect, _, err := client.ImageInspectWithRaw(context.Background(), imageId)
	if err == nil {
		result, err := json.MarshalIndent(inspect, "", "    ")

		if err == nil {
			var strResult = string(result)
			ShowTextPopup(app, "Image Details", strResult)
		}
	}

}

func RemoveImage(app *ui.Application, client *client.Client, imageId string) {
	_, err := client.ImageRemove(context.Background(), imageId, types.ImageRemoveOptions{})

	if err != nil {
		log.Print(err)
	}
}

func DoRunImage(image types.ImageSummary, command string) {
	var name = ""

	if len(image.RepoTags) > 0 {
		name = image.RepoTags[len(image.RepoTags)-1]
	} else {
		name = image.ID
	}

	RunCommand("docker", "run", "-it", "--entrypoint", command, name)
}

func RunShell(image types.ImageSummary) {
	DoRunImage(image, "sh")
}

func RunBashShell(image types.ImageSummary) {
	DoRunImage(image, "bash")
}

func BuildImagesView(app *ui.Application, client *client.Client, width uint8, height uint8) {
	var imageList = ui.ListNew()
	imageList.SetModel(docker.ImagesListModelNew(client))

	imageList.AddKeyHandler(input.KeyInputChar('v'), func(input.KeyInput) {
		var item = imageList.SelectedItem().Value().(*types.ImageSummary)
		ShowImageInspect(app, client, item.ID)
	})
	imageList.AddKeyHandler(input.KeyInputKey(keyboard.KeyDelete), func(input.KeyInput) {
		var item = imageList.SelectedItem().Value().(*types.ImageSummary)
		RemoveImage(app, client, item.ID)
		imageList.Update()
	})
	imageList.AddKeyHandler(input.KeyInputChar('s'), func(input.KeyInput) {
		var item = imageList.SelectedItem().Value().(*types.ImageSummary)
		RunShell(*item)
	})
	imageList.AddKeyHandler(input.KeyInputChar('b'), func(input.KeyInput) {
		var item = imageList.SelectedItem().Value().(*types.ImageSummary)
		RunBashShell(*item)
	})
	imageList.AddKeyHandler(input.KeyInputChar('h'), func(input.KeyInput) {
		ShowHelp(app)
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
