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

func (p bmProvisioner) ListQuotes(id int) (map[string]string, error) {
	return p.client.GetQuotes(id)
}

func (p bmProvisioner) ListPackages() (map[string]string, error) {
	return p.client.GetPackages()
}

func (p bmProvisioner) ListVMs(opts QueryOpts) (map[string]string, error) {
	return p.client.GetVMs(opts)
}

func (p bmProvisioner) DeleteVM(opts QueryOpts) error {
	return p.client.DeleteVM(opts)
}

func (p bmProvisioner) ListPackage(packageID int) (map[string]string, error) {
	return p.client.GetPackage(packageID)
}

func (p bmProvisioner) ListPackageTypes(id int) (map[string]string, error) {
	return p.client.GetPackageTypes(id)
}

func (p bmProvisioner) ProvisionACP(opts ACPOpts) error {
	return p.client.BuildOrder(opts)
}

func (p bmProvisioner) CreateHost(opts ACPOpts) (int, error) {
	fmt.Println("Provisioning Host...")
	return p.client.BuildVM(opts)
}
