package dockerlauncher

type Cmd struct {
	Command string `json:"command"`
}

type Container struct {
	ContainerID     string `json:"containerID"`
	ImageName       string `json:"imageName"`
	ContainerStatus string `json:"containerStatus"`
}

type GetContainersInfoReturn struct {
	Command    string      `json:"command"`
	Containers []Container `json:"containers"`
}

type UpdateImage struct {
	DeviceState struct {
		updateState string `json:"updateState"`
	} `json:"deviceState"`
}

type ErrorReturn struct {
	Message string `json:"message"`
}
