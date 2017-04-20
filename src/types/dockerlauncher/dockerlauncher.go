package dockerlauncher

type DeviceState struct {
	CurrentState string `json:"currentState"`
}

type Container struct {
	ContainerId     string `json:"containerId"`
	ContainerName   string `json:"containerName"`
	ImageName       string `json:"imageName"`
	ContainerStatus string `json:"containerStatus"`
}

type ErrorReturn struct {
	Message string `json:"message"`
}

type Cmd struct {
	Command string `json:"command"`
}
type UpdateParam struct {
	ImageName     string `json:"imageName"`
	ContainerName string `json:"containerName"`
}

type UpdateImageParameters struct {
	Command string      `json:"command"`
	Param   UpdateParam `json:"updateParam'`
}

type GetContainersInfoReturn struct {
	Containers []Container `json:"containers"`
}

type UpdateImageReturn struct {
	State DeviceState `json:"deviceState"`
}
