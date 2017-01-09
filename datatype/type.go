package datatype

type ContainerInfo struct {
	Id         string
	Ip         string
	App        string
	Component  string
	ImageName  string
	Repository string
	HostIp     string
}

type RC struct {
	ApiVersion string
	Kind       string
	Metadata   struct {
		Name   string
		Labels struct {
			App       string
			Componnet string
			System    string
			Ip        string
		}
	}
	Spec struct {
		Replicas string
		Selector struct {
			App       string
			Component string
			System    string
			Ip        string
		}
		Template struct{}
	}
}
