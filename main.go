package main

import (
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
	h: Shows this help
    Container view:
        v: Displays container information
        s: Opens a shell in a container
        b: Opens a bash shell in a container, if command is present
        l: Shows container log
        k: Kills a container
        delete: Deletes a container
    Images view:
        s: Creates a container and runs shell for a given image
        b: Creates a container and runs bash shell for a given image if command is present
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

func ShowContainerInspect(app *ui.Application, client *docker.ServiceHandler, containerId string) {
	strResult := client.InspectContainer(containerId)
	ShowTextPopup(app, "Container Details", strResult)
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

func ShowLogs(app *ui.Application, client *docker.ServiceHandler, containerId string) {

	logs := client.Logs(containerId)
	ShowTextPopup(app, "Logs", logs)
}

func SetupLog() {
	var logfile, err = os.OpenFile("dockermgr.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}
	log.SetOutput(logfile)

}

func BuildContainersView(app *ui.Application, client *docker.ServiceHandler, width uint8, height uint8) {
	var containerList = ui.ListNew()

	containerList.SetModel(docker.ContainerListModelNew(client))

	containerList.AddKeyHandler(input.KeyInputChar('v'), func(input.KeyInput) {
		var item = containerList.SelectedItem().Value().(*types.Container)
		ShowContainerInspect(app, client, item.ID)
	})
	containerList.AddKeyHandler(input.KeyInputChar('k'), func(input.KeyInput) {
		var item = containerList.SelectedItem().Value().(*types.Container)
		client.KillContainer(item.ID)
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
		client.RemoveContainer(item.ID)
		containerList.Update()
	})
	containerList.AddKeyHandler(input.KeyInputChar('a'), func(input.KeyInput) {
		containerList.Model.SetProperty(docker.ContainerListModelOnlyActive, nil)
	})
	containerList.AddKeyHandler(input.KeyInputChar('h'), func(input.KeyInput) {
		ShowHelp(app)
	})

	var titledContainer1 = ui.TitledContainerNew("Containers", containerList, false)
	titledContainer1.SetRect(ui.RectNew(1, 1, width, height-1))
	app.Add(titledContainer1)

}

func ShowImageInspect(app *ui.Application, client *docker.ServiceHandler, imageId string) {
	result := client.InspectImage(imageId)
	ShowTextPopup(app, "Image Details", result)
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

func BuildImagesView(app *ui.Application, client *docker.ServiceHandler, width uint8, height uint8) {
	var imageList = ui.ListNew()
	imageList.SetModel(docker.ImagesListModelNew(client))

	imageList.AddKeyHandler(input.KeyInputChar('v'), func(input.KeyInput) {
		var item = imageList.SelectedItem().Value().(*types.ImageSummary)
		ShowImageInspect(app, client, item.ID)
	})
	imageList.AddKeyHandler(input.KeyInputKey(keyboard.KeyDelete), func(input.KeyInput) {
		var item = imageList.SelectedItem().Value().(*types.ImageSummary)
		client.RemoveImage(item.ID)
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

	service := docker.ServiceHandlerNew(client)
	service.RemoveContainer("")

	if err1 != nil {
		panic(err1)
	}
	maxWidth, maxHeight := ui.ScreenSize()
	areaHeight := maxHeight / 2

	var app = ui.ApplicationNew()

	BuildContainersView(app, service, maxWidth, areaHeight)
	BuildImagesView(app, service, maxWidth, areaHeight)

	app.Loop()

}
