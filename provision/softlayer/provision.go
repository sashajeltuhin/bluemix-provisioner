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

func (p bmProvisioner) ListDCs(mode string) (map[string]string, error) {
	return p.client.GetDatacenters(mode)
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

func (p bmProvisioner) CreateHost(opts ACPOpts) ([]int, error) {
	fmt.Println("Provisioning Host...")
	return p.client.BuildVM(opts)
}

func getWinSpec(node *ACPNode, orderMode string) {
	if orderMode == "single" { //CreateObject method
		node.OS = "WIN_2012-STD_64"
		node.CPU = 2
		node.Mem = 4096
	} else {
		//order container
		node.Package = 46
		node.Prices = []int{1641, //GUEST_CORES_2 2 x 2.0 GHz Cores
			1645,   //RAM_2_GB 2 GB
			1639,   //GUEST_DISK_100_GB_SAN 100 GB (SAN)
			905,    //REBOOT_REMOTE_CONSOLE Reboot / Remote Console
			274,    // 1_GBPS_PUBLIC_PRIVATE_NETWORK_UPLINKS 1 Gbps Public & Private Network Uplinks
			1800,   //BANDWIDTH_0_GB_2 0 GB Bandwidth
			21,     //1_IP_ADDRESS 1 IP Address
			175777, //OS_WINDOWS_2012_R2_FULL_STD_64_BIT Windows Server 2012 R2 Standard Edition (64 bit)
			55,     //MONITORING_HOST_PING Host Ping
			57,     //NOTIFICATION_EMAIL_AND_TICKET Email and Ticket
			58,     //AUTOMATED_NOTIFICATION Automated Notification
			420,    //UNLIMITED_SSL_VPN_USERS_1_PPTP_VPN_USER_PER_ACCOUNT Unlimited SSL VPN Users & 1 PPTP VPN User per account
			418}    //NESSUS_VULNERABILITY_ASSESSMENT_REPORTING Nessus Vulnerability Assessment & Reporting
	}
}

func (p bmProvisioner) GetBootSpec(orderMode string) ACPNode {
	node := ACPNode{}
	node.Name = "boot"
	node.VMname = "ACP-BOOT"
	p.GetLinuxSpec(&node, orderMode)
	node.ScriptUri = "https://raw.githubusercontent.com/sashajeltuhin/bluemix-provisioner/master/provision/softlayer/scripts/bootstrap.ps1"
	return node
}

func (p bmProvisioner) GetDCSpec(orderMode string) ACPNode {
	node := ACPNode{}
	node.Name = "dc"
	node.VMname = "ACP-DC"
	getWinSpec(&node, orderMode)
	node.ScriptUri = "https://raw.githubusercontent.com/sashajeltuhin/bluemix-provisioner/master/provision/softlayer/scripts/bootDC.ps1"
	return node
}

func (p bmProvisioner) GetWebAppSpec(orderMode string) ACPNode {
	node := ACPNode{}
	node.Name = "webapp"
	getWinSpec(&node, orderMode)
	node.ScriptUri = "https://raw.githubusercontent.com/sashajeltuhin/bluemix-provisioner/master/provision/softlayer/scripts/webapp.ps1"
	return node
}

func (p bmProvisioner) GetLinuxSpec(node *ACPNode, orderMode string) {
	if orderMode == "single" { //CreateObject method
		node.OS = "CENTOS_7_64"
		node.CPU = 2
		node.Mem = 4096
	} else {
		//order container
		node.Package = 46
		node.Prices = []int{1641, //GUEST_CORES_2 2 x 2.0 GHz Cores
			1645,  //RAM_2_GB 2 GB
			2202,  //GUEST_DISK_25_GB_SAN 25 GB (SAN)
			905,   //REBOOT_REMOTE_CONSOLE Reboot / Remote Console
			274,   // 1_GBPS_PUBLIC_PRIVATE_NETWORK_UPLINKS 1 Gbps Public & Private Network Uplinks
			1800,  //BANDWIDTH_0_GB_2 0 GB Bandwidth
			21,    //1_IP_ADDRESS 1 IP Address
			45466, //OS_CENTOS_7_X_MINIMAL_64_BIT CentOS 7.x - Minimal Install (64 bit)
			55,    //MONITORING_HOST_PING Host Ping
			57,    //NOTIFICATION_EMAIL_AND_TICKET Email and Ticket
			58,    //AUTOMATED_NOTIFICATION Automated Notification
			420,   //UNLIMITED_SSL_VPN_USERS_1_PPTP_VPN_USER_PER_ACCOUNT Unlimited SSL VPN Users & 1 PPTP VPN User per account
			418}   //NESSUS_VULNERABILITY_ASSESSMENT_REPORTING Nessus Vulnerability Assessment & Reporting
	}
}
