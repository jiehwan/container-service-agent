package csaapi

type ContainerInfo struct {
	ContainerID     string `json:"container_id"`
	ContainerStatus string `json:"container_status"`
}
type ContainerLists struct {
	Cmd            string          `json:"cmd"`
	ContainerCount int             `json:"container_count"`
	Container      []ContainerInfo `json:"container"`
}

/*
// go-to-json output is follows.. but there is problem during init.
type ContainerLists struct {
	Cmd string `json:"cmd"`
	ContainerCount int `json:"container_count"`
	Container []struct {
		ContainerID string `json:"container_id"`
		ContainerStatus string `json:"container_status"`
	} `json:"container"`
}
*/
type ContainerService interface {
	GetContainersInfo() (ContainerLists, error)
}

/*
	STUB
*/
func GetContainersInfo() (ContainerLists, error) {

	send := ContainerLists{
		Cmd:            "getContainerLists",
		ContainerCount: 2,
		Container: []ContainerInfo{
			{
				ContainerID:     "api-1111",
				ContainerStatus: "running",
			},
			{
				ContainerID:     "api-2222",
				ContainerStatus: "exited",
			},
		},
	}

	return send, nil
}
