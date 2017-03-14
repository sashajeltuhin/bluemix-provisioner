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

type ACPOpts struct {
	APIUser string
	APIPass string
	DC      string
}

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acp",
		Short: "Provision ACP on Softlayer.",
		Long:  `Provision ACP on Softlayer.`,
	}

	cmd.AddCommand(DOCreateCmd())

	return cmd
}

func DOCreateCmd() *cobra.Command {
	opts := ACPOpts{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates infrastructure for a new ACP instance",
		Long:  `Creates infrastructure for a new ACP instance.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return makeInfra(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.APIUser, "user", "i", "", "API user name")
	cmd.Flags().StringVarP(&opts.APIPass, "apikey", "", "", "API key")
	cmd.Flags().StringVarP(&opts.DC, "dc", "", "", "Datacenter identifier")

	return cmd
}

func makeInfra(opts ACPOpts) error {
	reader := bufio.NewReader(os.Stdin)
	if opts.APIUser == "" {
		fmt.Print("Enter Your Softlayer User ID: \n")
		url, _ := reader.ReadString('\n')
		opts.APIUser = strings.Trim(url, "\n")
		opts.APIUser = strings.Replace(opts.APIUser, "\r", "", -1) //for Windows
	}
	if opts.APIPass == "" {
		fmt.Print("Enter Your Softlayer API key: \n")
		url, _ := reader.ReadString('\n')
		opts.APIPass = strings.Trim(url, "\n")
		opts.APIPass = strings.Replace(opts.APIPass, "\r", "", -1) //for Windows
	}

	fmt.Print("Provisioning\n")
	provisioner, errProv := GetProvisioner(opts.APIUser, opts.APIPass)
	if errProv != nil {
		return errProv
	}
	if opts.DC == "" {
		dcs, errDC := provisioner.ListDCs()
		if errDC != nil {
			return fmt.Errorf("Cannot load data centers. %v", errDC)
		}
		fmt.Print("Select Datacenter: \n")
		opts.DC = askForInput(dcs, reader)
	}
	err := provisioner.ProvisionACP(opts)

	if err != nil {
		return err
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
