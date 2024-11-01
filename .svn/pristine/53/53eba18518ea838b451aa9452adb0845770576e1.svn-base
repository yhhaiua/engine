package gxml

import (
	"testing"
)

type TestUser struct {
	Id          int    `xml:"id"`
	MapType     int    `xml:"mapType"`
	TypeName    string `xml:"typeName"`
	Name        string `xml:"name"`
	MapId       string `xml:"mapId" `
	PortMonster *PortM `xml:"portMonster" func:"UpMonster"`
}

type PortM struct {
	PortMonster string
}

func (t *TestUser) AfterLoad() {
	logger.Debugf("读取配置:%d", t.Id)
}

func (t *TestUser) UpMonster(str string) {
	t.PortMonster = &PortM{}
	t.PortMonster.PortMonster = str
}
func TestInitialize(t *testing.T) {

	path := "C:\\trunk\\data\\portGroup.xml"

	m := make(map[int]*TestUser)
	Initialize(path, &TestUser{}, m)
	r := len(m)
	logger.Infof("%d", r)
	//time.Sleep(time.Minute)
}
