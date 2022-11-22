package main

import (
	"fmt"

	"github.com/HewlettPackard/oneview-golang/ov"
	"github.com/HewlettPackard/oneview-golang/utils"
)

func main() {
	var (
		ClientOV *ov.OVClient
	)
	// Use configuratin file to set the ip and  credentails
	config, config_err := ov.LoadConfigFile("config.json")
	if config_err != nil {
		fmt.Println(config_err)
	}
	ovc := ClientOV.NewOVClient(
		config.OVCred.UserName,
		config.OVCred.Password,
		config.OVCred.Domain,
		config.OVCred.Endpoint,
		config.OVCred.SslVerify,
		config.OVCred.ApiVersion,
		config.OVCred.IfMatch)

	var (
		hypervisor_manager_ip           = config.HypervisorManagerConfig.IpAddress
		hypervisor_manager_display_name = "HM2"
		username                        = config.HypervisorManagerConfig.Username
		password                        = config.HypervisorManagerConfig.Password
	)
	scp, _ := ovc.GetScopeByName("ScopeTest")
	initialScopeUris := &[]utils.Nstring{(scp.URI)}
	//Adding Hypervisor Manager Server Certificate to Oneview for Secure conection
	server_cert, err := ovc.GetServerCertificateByIp(hypervisor_manager_ip)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Fetched Hypervisor Manager Server Certificate.", server_cert)
	}

	server_cert.CertificateDetails[0].AliasName = "Hypervisor Manager Server Certificate"
	server_cert.Type = ""
	er := ovc.CreateServerCertificate(server_cert)
	if er != nil {
		fmt.Println("............... Adding Server Certificate Failed: ", er)
	} else {
		fmt.Println("Imported Hypervisor Manager Server Certificate to Oneview for secure connection successfully.")
	}

	hypervisorManager := ov.HypervisorManager{
		DisplayName:      "HM1",
		Name:             hypervisor_manager_ip,
		Username:         username,
		Password:         password,
		Port:             443,
		InitialScopeUris: *initialScopeUris,
		Type:             "HypervisorManagerV2"}

	err = ovc.CreateHypervisorManager(hypervisorManager)
	if err != nil {
		fmt.Println("............... Create Hypervisor Manager Failed:", err)
	} else {
		fmt.Println(".... Create Hypervisor Manager Success")
	}

	fmt.Println("#................... Hypervisor Manager by Name ...............#")
	hypervisor_mgr, err := ovc.GetHypervisorManagerByName(hypervisor_manager_ip)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(hypervisor_mgr)
	}

	sort := "name:desc"
	hypervisor_mgr_list, err := ovc.GetHypervisorManagers("", "", "", sort)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("# ................... Hypervisor Managers List .................#")
		for i := 0; i < len(hypervisor_mgr_list.Members); i++ {
			fmt.Println(hypervisor_mgr_list.Members[i].Name)
		}
	}

	hypervisor_mgr.DisplayName = hypervisor_manager_display_name
	force := ""
	if config.OVCred.ApiVersion > 2400 {
		force = "true"
		err = ovc.UpdateHypervisorManager(hypervisor_mgr, force)
		fmt.Println("works")
	} else {
		err = ovc.UpdateHypervisorManager(hypervisor_mgr, force)
	}

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("#.................... Hypervisor Manager after Updating ...........#")
		hypervisor_mgr_after_update, err := ovc.GetHypervisorManagers("", "", "", sort)
		if err != nil {
			fmt.Println(err)
		} else {
			for i := 0; i < len(hypervisor_mgr_after_update.Members); i++ {
				fmt.Println(hypervisor_mgr_after_update.Members[i].Name)
			}
		}
	}
	err = ovc.DeleteHypervisorManager(hypervisor_manager_ip)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("#...................... Deleted Hypervisor Manager Successfully .....#")
	}
	//Create Hm for Automation
	err = ovc.CreateHypervisorManager(hypervisorManager)
	if err != nil {
		fmt.Println("............... Create Hypervisor Manager Failed:", err)
	} else {
		fmt.Println(".... Create Hypervisor Manager Success")
	}

}
