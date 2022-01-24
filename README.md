# clidockermgr
A CLI docker client made just as excuse for learning Go

![screenshot](screenshot.png)

Key bindings:

- ESC: Closes the active popup, or exits the application
- TAB: Cycles focus across views
- h: Show help

Containers:

- v: View container details
- s: Opens a shell in an active container
- b: Opens a BASH shell if the command exists in the container.
- l: View container logs
- k: Kill a container
- delete: Deletes a container

Images:

- v: View image details
- delete: Deletes an image
- s: Runs a shell session with the selected image.
- b: Opens a BASH shell if the command exists with the selected image.
