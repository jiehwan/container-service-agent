package csac

type ContainerInfo struct {
	ImageName       string `json:"ImageName"`
	ContainerName   string `json:"ContainerName"`
	ContainerStatus string `json:"ContainerStatus"`
}
type ContainerLists struct {
	Cmd            string          `json:"Cmd"`
	DeviceID       string          `json:"DeviceID"`
	ContainerCount int             `json:"ContainerCount"`
	Container      []ContainerInfo `json:"ContainerInfo"`
}
