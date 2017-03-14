package softlayer

import (
	//	"bufio"
	"fmt"
	//	"os"
	//"context"
	//"io/ioutil"
	s "strings"

	"github.com/Jeffail/gabs"

	bmtypes "github.com/softlayer/softlayer-go/datatypes"
	bmservices "github.com/softlayer/softlayer-go/services"
	bmsession "github.com/softlayer/softlayer-go/session"
	//bmsl "github.com/softlayer/softlayer-go/sl"
)

type NodeConfig struct {
	Image    string
	Name     string
	Region   string
	Size     string
	UserData string
}

// Client for provisioning machines on AWS
type Client struct {
	apiSession *bmsession.Session
}

func (c *Client) getAPISession(username string, apikey string) (*bmsession.Session, error) {
	if c.apiSession == nil {
		c.apiSession = bmsession.New(username, apikey)
	}
	return c.apiSession, nil
}

func (c *Client) GetDatacenters() (map[string]string, error) {
	var dcs []bmtypes.Location
	var err error
	loc := bmservices.GetLocationService(c.apiSession)
	dcs, err = loc.GetDatacenters()
	if err != nil {
		return nil, fmt.Errorf("Cannot load datacenters %v", err)
	}
	var objMap map[string]string = make(map[string]string)
	for i := 0; i < len(dcs); i++ {
		var dc bmtypes.Location = dcs[i]
		objMap[*dc.Name] = *dc.LongName
	}
	return objMap, nil
}

func (c *Client) parseObj(body []byte, objNode string, idfield string, namefield string) (map[string]string, error) {
	var objMap map[string]string = make(map[string]string)
	jsonParsed, err := gabs.ParseJSON(body)
	if err != nil {
		return nil, err
	}
	if jsonParsed.Exists(objNode) {
		ids := s.TrimLeft(jsonParsed.Path(fmt.Sprintf("%s.%s", objNode, idfield)).String(), "[")
		ids = s.TrimRight(ids, "]")
		idarray := s.Split(ids, ",")
		names := s.TrimLeft(jsonParsed.Path(fmt.Sprintf("%s.%s", objNode, namefield)).String(), "[")
		names = s.TrimRight(names, "]")
		namearray := s.Split(names, ",")
		count := len(idarray)
		for i := 0; i < count; i++ {
			objMap[idarray[i]] = namearray[i]
		}
	}
	return objMap, nil
}
