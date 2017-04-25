package softlayer

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"strings"

	"github.com/sashajeltuhin/bluemix-provisioner/provision/utils"
	"github.com/spf13/cobra"
)

type CredOpts struct {
	APIUser string
	APIPass string
}

type ACPOpts struct {
	CredOpts
	DC         string
	Domain     string
	HourlyBill bool
	NumLM      int
	NUMWin     int
	QuoteOnly  bool
	VerifyOnly bool
}

type ACPGroup struct {
	Name   string
	VMname string
	Count  int
	Prices []int
}

type QueryOpts struct {
	CredOpts
	ID     int
	Domain string
}

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acp",
		Short: "Provision ACP on Softlayer.",
		Long:  `Provision ACP on Softlayer.`,
	}

	cmd.AddCommand(ACPCreateCmd())
	cmd.AddCommand(DCCreateCmd())
	cmd.AddCommand(ShowQuotesCmd())
	cmd.AddCommand(ShowPackagesCmd())
	cmd.AddCommand(ShowVMsCmd())
	cmd.AddCommand(DeleteVMsCmd())

	return cmd
}

func ACPCreateCmd() *cobra.Command {
	opts := ACPOpts{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates infrastructure for a new ACP instance",
		Long:  `Creates infrastructure for a new ACP instance.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return makeInfra(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.APIUser, "user", "u", "", "API user name")
	cmd.Flags().StringVarP(&opts.APIPass, "apikey", "p", "", "API key")
	cmd.Flags().StringVarP(&opts.DC, "dc", "", "", "Datacenter identifier")
	cmd.Flags().StringVarP(&opts.Domain, "domain", "", "", "Domain for the ACP instance")
	cmd.Flags().BoolVarP(&opts.HourlyBill, "hourly", "", true, "Hourly billing. If set to false, monthly billing will apply to the order")
	cmd.Flags().BoolVarP(&opts.QuoteOnly, "quote", "", false, "Produce a quote only")
	cmd.Flags().BoolVarP(&opts.VerifyOnly, "verify", "", false, "Check that the order is put together correclty. Do not place yet")
	cmd.Flags().IntVarP(&opts.NumLM, "count-lm", "", 2, "Number of load managers. Default 2")
	cmd.Flags().IntVarP(&opts.NUMWin, "count-win", "", 4, "Number of web/app servers. Default 4")

	return cmd
}

func DCCreateCmd() *cobra.Command {
	opts := ACPOpts{}
	cmd := &cobra.Command{
		Use:   "createDC",
		Short: "Creates a single domain controller",
		Long:  `Creates a single domain controller for a new ACP instance.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return createDC(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.APIUser, "user", "u", "", "API user name")
	cmd.Flags().StringVarP(&opts.APIPass, "apikey", "", "", "API key")
	cmd.Flags().StringVarP(&opts.DC, "dc", "", "", "Datacenter identifier")
	cmd.Flags().StringVarP(&opts.Domain, "domain", "", "", "ACP Windows Domain name")

	return cmd
}

func ShowQuotesCmd() *cobra.Command {
	opts := QueryOpts{}
	cmd := &cobra.Command{
		Use:   "quotes",
		Short: "Displays a list of saved quotes",
		Long:  `Displays a list of saved quotes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return showQuotes(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.APIUser, "user", "u", "", "API user name")
	cmd.Flags().StringVarP(&opts.APIPass, "apikey", "", "", "API key")
	cmd.Flags().IntVarP(&opts.ID, "id", "", 0, "ID of the item to fetch")
	return cmd
}

func ShowVMsCmd() *cobra.Command {
	opts := QueryOpts{}
	cmd := &cobra.Command{
		Use:   "vms",
		Short: "Displays a list of active VMs",
		Long:  `Displays a list of active VMs.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return showVMs(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.APIUser, "user", "u", "", "API user name")
	cmd.Flags().StringVarP(&opts.APIPass, "apikey", "", "", "API key")
	cmd.Flags().IntVarP(&opts.ID, "id", "", 0, "ID of the item to fetch")
	cmd.Flags().StringVarP(&opts.Domain, "domain", "", "", "ACP Windows Domain name")
	return cmd
}

func DeleteVMsCmd() *cobra.Command {
	opts := QueryOpts{}
	cmd := &cobra.Command{
		Use:   "deletevm",
		Short: "Deletes VM with specified ID",
		Long:  `Deletes VM with specified ID.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return DeleteVM(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.APIUser, "user", "u", "", "API user name")
	cmd.Flags().StringVarP(&opts.APIPass, "apikey", "", "", "API key")
	cmd.Flags().IntVarP(&opts.ID, "id", "", 0, "ID of the item to delete")
	return cmd
}

func ShowPackagesCmd() *cobra.Command {
	opts := QueryOpts{}
	cmd := &cobra.Command{
		Use:   "packages",
		Short: "Displays a list of product packages",
		Long:  `Displays a list of product packages.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return showPackages(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.APIUser, "user", "u", "", "API user name")
	cmd.Flags().StringVarP(&opts.APIPass, "apikey", "", "", "API key")
	cmd.Flags().IntVarP(&opts.ID, "id", "", 0, "ID of the item to fetch")
	return cmd
}

func getCreds(opts *CredOpts) (*bufio.Reader, *bmProvisioner, error) {
	reader := bufio.NewReader(os.Stdin)
	opts.APIUser = os.Getenv("IBM_API_USER")
	if opts.APIUser == "" {
		fmt.Print("Enter Your Softlayer User ID: \n")
		url, _ := reader.ReadString('\n')
		opts.APIUser = strings.Trim(url, "\n")
		opts.APIUser = strings.Replace(opts.APIUser, "\r", "", -1) //for Windows
	}
	opts.APIPass = os.Getenv("IBM_API_KEY")
	if opts.APIPass == "" {
		fmt.Print("Enter Your Softlayer API key: \n")
		url, _ := reader.ReadString('\n')
		opts.APIPass = strings.Trim(url, "\n")
		opts.APIPass = strings.Replace(opts.APIPass, "\r", "", -1) //for Windows
	}
	provisioner, errProv := GetProvisioner(opts.APIUser, opts.APIPass)
	if errProv != nil {
		return nil, nil, errProv
	}
	return reader, provisioner, nil
}

func askForDC(opts *ACPOpts, reader *bufio.Reader, provisioner *bmProvisioner) error {
	if opts.DC == "" {
		dcs, errDC := provisioner.ListDCs()
		if errDC != nil {
			return fmt.Errorf("Cannot load data centers. %v", errDC)
		}
		fmt.Print("Select Datacenter: \n")
		opts.DC = askForInput(dcs, reader)
	}
	return nil
}

func askDomain(opts *ACPOpts, reader *bufio.Reader) {
	if opts.Domain == "" {
		fmt.Print("Enter Domain name for VMs of ACP instance: \n")
		url, _ := reader.ReadString('\n')
		opts.Domain = strings.Trim(url, "\n")
		opts.Domain = strings.Replace(opts.Domain, "\r", "", -1) //for Windows
	}
}

func showQuotes(opts QueryOpts) error {
	_, provisioner, _ := getCreds(&opts.CredOpts)
	quotes, errP := provisioner.ListQuotes(opts.ID)
	if errP != nil {
		return fmt.Errorf("Cannot load quotes. %v", errP)
	}
	fmt.Println("Number of quotes", len(quotes))
	return nil
}

func showVMs(opts QueryOpts) error {
	_, provisioner, _ := getCreds(&opts.CredOpts)
	vms, errP := provisioner.ListVMs(opts)
	if errP != nil {
		return fmt.Errorf("Cannot load VMs. %v", errP)
	}
	fmt.Println("Number of VMs", len(vms))
	return nil
}

func DeleteVM(opts QueryOpts) error {
	_, provisioner, _ := getCreds(&opts.CredOpts)
	errP := provisioner.DeleteVM(opts)
	if errP != nil {
		return fmt.Errorf("Cannot delete VMs. %v", errP)
	}

	return nil
}

func showPackages(opts QueryOpts) error {
	_, provisioner, _ := getCreds(&opts.CredOpts)
	var errP error
	var list map[string]string
	if opts.ID > 0 {
		list, errP = provisioner.ListPackage(opts.ID)
	} else {
		list, errP = provisioner.ListPackages()
	}

	if errP != nil {
		return fmt.Errorf("Cannot load packages. %v", errP)
	}
	fmt.Println("Number in list", len(list))
	return nil
}

func createDC(opts ACPOpts) error {
	reader, provisioner, _ := getCreds(&opts.CredOpts)
	err := askForDC(&opts, reader, provisioner)
	if err != nil {
		return err
	}
	id, e := provisioner.CreateHost(opts)
	if e != nil {
		return e
	}
	fmt.Println("VM orchestration started. ID = ", id)
	return nil
}

func makeInfra(opts ACPOpts) error {
	reader, provisioner, _ := getCreds(&opts.CredOpts)
	err := askForDC(&opts, reader, provisioner)
	if err != nil {
		return err
	}
	askDomain(&opts, reader)
	errProv := provisioner.ProvisionACP(opts)

	if errProv != nil {
		return errProv
	}
	fmt.Println("Done")
	return nil
}

func askForInput(objList map[string]string, reader *bufio.Reader) string {
	arrPairs := utils.SortMapbyVal(objList)
	count := len(objList)
	var arr = make([]string, count)
	for i := 0; i < count; i++ {
		arr[i] = arrPairs[i].Key
		fmt.Printf("%d - %s\n", i+1, arrPairs[i].Value)
	}

	objI, _ := reader.ReadString('\n')
	objIndex := strings.Trim(string(objI), "\n")
	index, _ := strconv.Atoi(objIndex)
	if index < 1 || index > len(objList) {
		fmt.Print("Invalid selection. Try again")
		return askForInput(objList, reader)
	} else {
		objID := arr[index-1]
		fmt.Println("You picked ", objList[objID])
		objID = strings.Trim(objID, "\"")
		return objID
	}
}
