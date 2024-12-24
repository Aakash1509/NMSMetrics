package models

type Credential struct {
	ID        int64  `json:"profile_id"`
	Protocol  string `json:"profile_protocol"`
	Name      string `json:"user_name"`
	Password  string `json:"user_password"`
	Community string `json:"community"`
	Version   string `json:"version"`
}

type Input struct {
	EventType       string       `json:"event.type"`
	DeviceType      string       `json:"device.type"`
	IP              string       `json:"ip"`
	Port            int          `json:"port"`
	Credentials     []Credential `json:"discovery.credential.profiles"`
	MetricGroup     string       `json:"metric.group.name"`
	ProfileProtocol string       `json:"profile.protocol"`
	UserName        string       `json:"user.name"`
	UserPassword    string       `json:"user.password"`
	Community       string       `json:"community"`
	Version         string       `json:"version"`
}
