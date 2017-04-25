package softlayer

import (
	//	"bufio"
	"fmt"
	//	"os"
	//"context"
	//"io/ioutil"
	"strconv"
	s "strings"

	"github.com/Jeffail/gabs"

	//	"github.com/sashajeltuhin/bluemix-provisioner/provision/utils"
	bmtypes "github.com/softlayer/softlayer-go/datatypes"
	filter "github.com/softlayer/softlayer-go/filter"
	bmservices "github.com/softlayer/softlayer-go/services"
	bmsession "github.com/softlayer/softlayer-go/session"
)

type NodeConfig struct {
	Image    string
	Name     string
	Domain   string
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
		return nil, fmt.Errorf("Cannot load datacenters %v \n", err)
	}
	var objMap map[string]string = make(map[string]string)
	for i := 0; i < len(dcs); i++ {
		var dc bmtypes.Location = dcs[i]
		id := strconv.Itoa(*dc.Id)
		fmt.Println("DC ID:", id, *dc.LongName)
		objMap[id] = *dc.LongName
	}
	return objMap, nil
}

func (c *Client) GetPackageTypes(id int) (map[string]string, error) {
	var types []bmtypes.Product_Package_Type
	var err error
	serv := bmservices.GetProductPackageTypeService(c.apiSession)
	types, err = serv.GetAllObjects()
	if err != nil {
		return nil, fmt.Errorf("Cannot load types %v", err)
	}
	var objMap map[string]string = make(map[string]string)
	for i := 0; i < len(types); i++ {
		var p bmtypes.Product_Package_Type = types[i]
		fmt.Println("Package Type:", *p.Id, *p.KeyName, *p.Name)
		objMap[string(*p.Id)] = *p.KeyName
	}
	if id > 0 {
		packages, errP := serv.Id(id).GetPackages()
		if errP != nil {
			return nil, fmt.Errorf("Cannot load packages %v", errP)
		}
		for i := 0; i < len(packages); i++ {
			var p bmtypes.Product_Package = packages[i]
			fmt.Println("")
			fmt.Println("Package:", *p.Id, *p.Name)
			prods, errProd := c.GetProducts(p)
			if errProd != nil {
				fmt.Errorf("Cannot get products for package %s. %v \n", p.Name, errProd)
			} else {
				fmt.Println("Products:")
				for ii := 0; ii < len(prods); ii++ {
					prod := prods[ii]
					fmt.Println(*prod.Id, *prod.KeyName, *prod.Description)
					fmt.Println("Prices:")
					for y := 0; y < len(prod.Prices); y++ {
						price := prod.Prices[y]
						fmt.Println("ID ", *price.Id)
						if price.OrderOptions != nil {
							for o := 0; o < len(price.OrderOptions); o++ {
								fmt.Println("Option", *price.OrderOptions[o].Name, *price.OrderOptions[o].Description)
							}
						}
						if price.RecurringFee != nil {
							fmt.Println("Recurring fee: ", *price.RecurringFee)
						}
						if price.OneTimeFee != nil {
							fmt.Println("One-time fee: ", *price.OneTimeFee)
						}
						if price.HourlyRecurringFee != nil {
							fmt.Println("Hourly fee: ", *price.HourlyRecurringFee)
						}
						if price.UsageRate != nil {
							fmt.Println("Hourly fee: ", *price.UsageRate)
						}
					}
				}
			}
		}
	}
	return objMap, nil
}

func (c *Client) GetPackages() (map[string]string, error) {
	var packages []bmtypes.Product_Package
	var err error
	serv := bmservices.GetProductPackageService(c.apiSession)
	packages, err = serv.GetAllObjects()
	if err != nil {
		return nil, fmt.Errorf("Cannot load packages %v", err)
	}
	var objMap map[string]string = make(map[string]string)
	for i := 0; i < len(packages); i++ {
		var p bmtypes.Product_Package = packages[i]
		fmt.Println("Package:", *p.Id, *p.Name)
		objMap[string(*p.Id)] = *p.Name
		prods, errProd := c.GetProducts(p)
		if errProd != nil {
			fmt.Errorf("Cannot get products for package %s. %v \n", p.Name, errProd)
		} else {
			fmt.Println("Products:")
			for ii := 0; ii < len(prods); ii++ {
				prod := prods[ii]
				fmt.Println(*prod.Id, *prod.KeyName, *prod.Description)
			}
		}
	}
	return objMap, nil
}

func (c *Client) GetPackage(packageID int) (map[string]string, error) {

	var objMap map[string]string = make(map[string]string)
	serv := bmservices.GetProductPackageService(c.apiSession)
	packageObj, err := serv.Id(packageID).GetObject()
	if err != nil {
		return nil, fmt.Errorf("Cannot get package %d info %v", packageID, err)
	}
	products, errProd := serv.Id(packageID).GetItems()
	if errProd != nil {
		return nil, fmt.Errorf("Cannot get products for package %d info %v", packageID, errProd)
	}

	//	packPrices, errpc := serv.Id(packageID).Mask("id;item.description;categories.id").GetItemPrices()
	//	if errpc != nil {
	//		return nil, fmt.Errorf("Cannot get prices for package %d info %v", packageID, errpc)
	//	}

	config, errCat := serv.Id(packageID).Mask("isRequired;itemCategory").GetConfiguration()
	if errCat != nil {
		return nil, fmt.Errorf("Cannot get categories for package %d info %v", packageID, errCat)
	}

	fmt.Println("Package:", *packageObj.Id, *packageObj.Name)

	fmt.Println("Required categories:")
	var requiredCats []int
	for c := 0; c < len(config); c++ {
		conf := config[c]
		if conf.IsRequired != nil && *conf.IsRequired > 0 {
			requiredCats = append(requiredCats, *conf.ItemCategory.Id)
			fmt.Printf("Required Configs: %d - %s\n", *conf.ItemCategory.Id, *conf.ItemCategory.Name)
			//			fmt.Println("Required Category Items:")
			//			for pc := 0; pc < len(packPrices); pc++ {
			//				itemprice := packPrices[pc]
			//				objMap[strconv.Itoa(*itemprice.Id)] = *itemprice.Item.Description
			//				for ii := 0; ii < len(itemprice.Categories); ii++ {
			//					itempriceCat := itemprice.Categories[ii]
			//					if utils.ContainsInt(requiredCats, *itempriceCat.Id) {
			//						fmt.Printf("%d - %s (%d)\n", *itemprice.Id, *itemprice.Item.Description, *itempriceCat.Id)
			//					}
			//				}

			//			}
		}
	}

	fmt.Println("Products:")
	for i := 0; i < len(products); i++ {
		var p bmtypes.Product_Item = products[i]
		fmt.Println("Product:", *p.Id, *p.KeyName, *p.Description)
		objMap[strconv.Itoa(*p.Id)] = *p.KeyName
		fmt.Println("Prices:")
		for y := 0; y < len(p.Prices); y++ {
			price := p.Prices[y]
			fmt.Println("ID ", *price.Id)

			if price.Quantity != nil {
				fmt.Println("Quantity: ", *price.Quantity)
			}
			if price.RecurringFee != nil {
				fmt.Println("Recurring fee: ", *price.RecurringFee)
			}

			if price.OneTimeFee != nil {
				fmt.Println("One-time fee: ", *price.OneTimeFee)
			}

			if price.HourlyRecurringFee != nil {
				fmt.Println("Hourly fee: ", *price.HourlyRecurringFee)
			}

		}
	}
	return objMap, nil
}

func (c *Client) GetProducts(pack bmtypes.Product_Package) (products []bmtypes.Product_Item, err error) {
	serv := bmservices.GetProductPackageService(c.apiSession)
	return serv.Id(*pack.Id).GetItems()
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

func (c *Client) BuildOrder(opts ACPOpts) error {
	numCats := 1
	//category fields
	hourly := opts.HourlyBill
	num := 1
	vmname := "ACP-WIN-%d"
	items := []int{1641, //GUEST_CORES_2 2 x 2.0 GHz Cores
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
	packageID := 46

	serv := bmservices.GetProductOrderService(c.apiSession)
	var container bmtypes.Container_Product_Order
	container.OrderContainers = []bmtypes.Container_Product_Order{}
	//services
	for c := 0; c < numCats; c++ {
		var child bmtypes.Container_Product_Order
		child.Location = &opts.DC
		fmt.Println("Location:", *child.Location)
		child.PackageId = &packageID
		child.UseHourlyPricing = &hourly
		child.Quantity = &num
		//service prices
		child.Prices = []bmtypes.Product_Item_Price{}
		for p := 0; p < len(items); p++ {
			child.Prices = append(child.Prices, bmtypes.Product_Item_Price{Id: &items[p]})
		}
		//VMs for the service
		child.VirtualGuests = []bmtypes.Virtual_Guest{}
		for i := 0; i < num; i++ {
			var vm bmtypes.Virtual_Guest
			sname := fmt.Sprintf(vmname, i+1)
			vm.Hostname = &sname
			vm.Domain = &opts.Domain
			child.VirtualGuests = append(child.VirtualGuests, vm)
		}
		container.OrderContainers = append(container.OrderContainers, child)
	}

	var err error
	var receipt bmtypes.Container_Product_Order_Receipt
	var temp bmtypes.Container_Product_Order
	if opts.QuoteOnly {
		receipt, err = serv.PlaceQuote(&container)
		if receipt.OrderDetails != nil {
			temp = *receipt.OrderDetails
		}
	} else if opts.VerifyOnly {
		temp, err = serv.VerifyOrder(&container)
	} else {
		save := false
		receipt, err = serv.PlaceOrder(&container, &save)
		if receipt.OrderDetails != nil {
			temp = *receipt.OrderDetails
		}
	}
	if err != nil {
		return fmt.Errorf("Verify order failed: %v", err)
	}

	fmt.Println("Order Summary:")
	c.showOrderContainer(temp)

	account := bmservices.GetAccountService(c.apiSession)

	fmt.Println("Account VMs:")
	vms, errvm := account.Mask("id;hostname").GetVirtualGuests()
	if errvm != nil {
		return fmt.Errorf("Cannot retrieve VMs for the account: %v", errvm)
	}
	if len(vms) == 0 {
		fmt.Println("No VMs yet")
	}
	for i := 0; i < len(vms); i++ {
		fmt.Printf("%d - %s\n", *vms[i].Id, *vms[i].Hostname)
	}

	return nil

	//isPingable

	//get virtual guest status
	//getObject - network component - ipAddressBindings - Network_Component_IpAddress - Network_Subnet_IpAddress
	//Account::getAllBillingItems

	//Cancel: Billing Item cancel service

}

func (c *Client) BuildVM(opts ACPOpts) (int, error) {
	vmService := bmservices.GetVirtualGuestService(c.apiSession)
	var obj bmtypes.Virtual_Guest

	newVM, err := vmService.CreateObject(&obj)
	if err != nil {
		fmt.Errorf("Cannot create VM %v \n", err)
	}
	return *newVM.Id, nil
}

func (c *Client) GetQuotes(id int) (map[string]string, error) {
	var quotes []bmtypes.Billing_Order_Quote
	var err error
	accService := bmservices.GetAccountService(c.apiSession)
	quotes, err = accService.GetQuotes()
	if err != nil {
		return nil, fmt.Errorf("Cannot load quotes %v \n", err)
	}
	var objMap map[string]string = make(map[string]string)
	for i := 0; i < len(quotes); i++ {
		var q bmtypes.Billing_Order_Quote = quotes[i]
		fmt.Println(*q.Id, *q.Name, *q.CreateDate, *q.Status)
		objMap[strconv.Itoa(*q.Id)] = *q.Name
	}
	if id > 0 {
		billservice := bmservices.GetBillingOrderQuoteService(c.apiSession)
		var fl bool = false
		container, _ := billservice.Id(id).GetRecalculatedOrderContainer(nil, &fl)
		c.showOrderContainer(container)
	}

	return objMap, nil
}

func (c *Client) showOrderContainer(container bmtypes.Container_Product_Order) error {
	fmt.Println("Order Container:")

	if container.ProratedInitialCharge != nil {
		fmt.Println("Prorated Initial Charge:", *container.ProratedInitialCharge)
	}

	if container.ProratedOrderTotal != nil {
		fmt.Println("Prorated Order Total:", *container.ProratedOrderTotal)
	}

	if container.PostTaxRecurring != nil {
		fmt.Println("PostTaxRecurring Charge", *container.PostTaxRecurring)
	}

	if container.PostTaxRecurringHourly != nil {
		fmt.Println("PostTaxRecurringHourly Charge", *container.PostTaxRecurringHourly)
	}

	if container.PostTaxRecurringMonthly != nil {
		fmt.Println("PostTaxRecurringMonthly Charge", *container.PostTaxRecurringMonthly)
	}

	if container.OrderContainers != nil {
		fmt.Println("Child Container")
		for i := 0; i < len(container.OrderContainers); i++ {
			child := container.OrderContainers[i]
			c.showOrderContainer(child)
		}
	}

	if container.Quantity != nil {
		fmt.Println("Quantity:", *container.Quantity)
	}

	if container.Location != nil {
		dcService := bmservices.GetLocationDatacenterService(c.apiSession)
		id, _ := strconv.Atoi(*container.Location)
		dc, _ := dcService.Id(id).GetObject()
		fmt.Println("Location:", *dc.Id, *dc.Name)
	}

	if container.PackageId != nil {
		fmt.Println("PackageId:", *container.PackageId)
	}

	if container.ImageTemplateId != nil {
		fmt.Println("ImageTemplateId:", *container.ImageTemplateId)
	}

	if container.VirtualGuests != nil {
		fmt.Println("VirtualGuests:")
		for g := 0; g < len(container.VirtualGuests); g++ {
			guest := container.VirtualGuests[g]
			if guest.Id != nil {
				fmt.Println("VM ID:", *guest.Id)
			}
			if guest.Domain != nil {
				fmt.Println("VM domain:", *guest.Domain)
			}
			if guest.Hostname != nil {
				fmt.Println("VM host name:", *guest.Hostname)
			}
		}
	}

	if container.Prices != nil && len(container.Prices) > 0 {
		fmt.Println("Prices:")
		itemPriceSvc := bmservices.GetProductItemPriceService(c.apiSession)
		for ii := 0; ii < len(container.Prices); ii++ {
			pr := container.Prices[ii]
			if pr.Id != nil {
				fmt.Println("Price ID:", *pr.Id)
				prod, prodErr := itemPriceSvc.Id(*pr.Id).GetItem()
				if prodErr == nil {
					fmt.Println("Product Data:")
					fmt.Println(*prod.Id, *prod.KeyName, *prod.Description)
				}
			}

			if pr.ItemId != nil {
				fmt.Println("Item ID:", *pr.ItemId)
			}

			if pr.RecurringFee != nil {
				fmt.Println("RecurringFee:", *pr.RecurringFee)
			}

			if pr.Quantity != nil {
				fmt.Println("Quantity:", *pr.Quantity)
			}
			if pr.HourlyRecurringFee != nil {
				fmt.Println("HourlyRecurringFee:", *pr.HourlyRecurringFee)
			}
		}
	}
	return nil
}

func (c *Client) GetVMs(opts QueryOpts) (map[string]string, error) {
	account := bmservices.GetAccountService(c.apiSession)
	var vms []bmtypes.Virtual_Guest
	var errvm error
	fmt.Println("Account VMs:")
	mask := "id;hostname;primaryIpAddress;primaryBackendIpAddress"
	if opts.Domain != "" {
		vms, errvm = account.Mask(mask).Filter(filter.Path("virtualGuests.domain").Eq(opts.Domain).Build()).GetVirtualGuests()
	} else {
		vms, errvm = account.Mask(mask).GetVirtualGuests()
	}
	if errvm != nil {
		return nil, fmt.Errorf("Cannot retrieve VMs for the account: %v", errvm)
	}
	if len(vms) == 0 {
		fmt.Println("No VMs yet")
	}
	var objMap map[string]string = make(map[string]string)
	for i := 0; i < len(vms); i++ {
		fmt.Printf("%d - %s. Public IP: %s  Private IP:  %s\n", *vms[i].Id, *vms[i].Hostname, *vms[i].PrimaryIpAddress, *vms[i].PrimaryBackendIpAddress)
		objMap[strconv.Itoa(*vms[i].Id)] = *vms[i].Hostname
	}

	return objMap, nil
}

func (c *Client) DeleteVM(opts QueryOpts) error {
	vmserv := bmservices.GetVirtualGuestService(c.apiSession)

	ok, err := vmserv.Id(opts.ID).DeleteObject()

	if ok {
		return nil
	} else {
		return fmt.Errorf("Unable to delete VM %d. %v\n", opts.ID, err)
	}
}

//	fmt.Println("Account billing items:")
//	billed, errbill := account.GetAllBillingItems()
//	if errbill != nil {
//		return fmt.Errorf("Cannot retrieve billing items for the account: %v", errbill)
//	}
//	for b := 0; b < len(billed); b++ {
//		item := billed[b]
//		hostName := ""
//		orderItemId := 0
//		if item.HostName != nil {
//			hostName = *item.HostName
//		}
//		if item.OrderItemId != nil {
//			orderItemId = *item.OrderItemId
//		}
//		fmt.Printf("Billing ID: %d - %s  Hostname - %s, OrderItemID - %d\n", *item.Id, *item.Description, hostName, orderItemId)
//	}
