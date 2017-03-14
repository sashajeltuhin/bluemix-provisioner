package softlayer

import (
	"fmt"
	//	"io/ioutil"
	//"os"
	//"path/filepath"
	//"regexp"
	//"time"
)

type sshMachineProvisioner struct {
	sshKey string
}

func (p sshMachineProvisioner) SSHKey() string {
	return p.sshKey
}

type bmProvisioner struct {
	sshMachineProvisioner
	client *Client
}

func GetProvisioner(username string, apikey string) (*bmProvisioner, error) {
	c := Client{}
	_, err := c.getAPISession(username, apikey)
	if err != nil {
		return nil, fmt.Errorf("Unable to initialize Softlayer API session. %v", err)
	}
	p := bmProvisioner{client: &c}
	return &p, nil
}

func (p bmProvisioner) ListDCs() (map[string]string, error) {
	return p.client.GetDatacenters()
}

func (p bmProvisioner) ProvisionACP(opts ACPOpts) error {
	fmt.Println("Provisioning...")
	return nil
}
