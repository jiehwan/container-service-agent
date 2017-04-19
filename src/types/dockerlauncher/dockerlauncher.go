package dockerlauncher

type DeviceState struct {
	UpdateState string `json:"updateState"`
}

type UpdateImage struct {
	State DeviceState `json:"deviceState"`
}

type Container struct {
	ContainerID     string `json:"containerID"`
	ImageName       string `json:"imageName"`
	ContainerStatus string `json:"containerStatus"`
}

type ErrorReturn struct {
	Message string `json:"message"`
}

type Cmd struct {
	Command string `json:"command"`
}

type UpdateImageParameters struct {
	Command     string      `json:"command"`
	UpdateParam UpdateImage `json:"updateParam'`
}

type GetContainersInfoReturn struct {
	Command    string      `json:"command"`
	Containers []Container `json:"containers"`
}

type GetUpdateImageReturn struct {
	Command string      `json:"command"`
	State   DeviceState `json:"deviceState"`
}
