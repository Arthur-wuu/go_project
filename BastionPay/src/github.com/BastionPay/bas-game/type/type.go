package _type

//
//type FirstRcv struct {
//	Type        string    `yaml:"type"`
//	UpNumber    string    `yaml:"upnumber"`
//}
//
//type SecondRcv struct {
//	Type        string    `yaml:"type"`
//	Count       string    `yaml:"count"`
//	Message     string    `yaml:"message"`
//	State       string    `yaml:"state"`
//}

type MsgRcv struct {
	Type     string `yaml:"type"`
	UpNumber string `yaml:"upnumber"`
	Message  string `yaml:"message"`
	State    string `yaml:"state"`
	Count    string `yaml:"count"`
}
