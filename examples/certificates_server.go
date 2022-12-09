package main

import (
	"fmt"

	"github.com/HewlettPackard/oneview-golang/ov"
	"github.com/HewlettPackard/oneview-golang/utils"
)

func main() {
	config, config_err := ov.LoadConfigFile("config.json")
	if config_err != nil {
		fmt.Println(config_err)
	}
	var (
		ClientOV                *ov.OVClient
		server_certificate_ip                 = config.ServerCertificateIp
		server_certificate_name               = "new_test_certificate"
		new_cert_base64data     utils.Nstring = "---BEGIN CERTIFICATE----END CERTIFICATE------"
	)
	// Use configuratin file to set the ip and  credentails
	ovc := ClientOV.NewOVClient(
		config.OVCred.UserName,
		config.OVCred.Password,
		config.OVCred.Domain,
		config.OVCred.Endpoint,
		config.OVCred.SslVerify,
		config.OVCred.ApiVersion,
		config.OVCred.IfMatch)

	server_cert, err := ovc.GetServerCertificateByIp(server_certificate_ip)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(server_cert)
	}
	server_cert.CertificateDetails[0].AliasName = server_certificate_name
	fmt.Println(server_cert.CertificateDetails[0].AliasName)
	server_cert.Type = "" //Making Type field as empty as it is not required

	er := ovc.CreateServerCertificate(server_cert)
	if er != nil {
		fmt.Println("............... Adding Server Certificate Failed:", er)
	} else {
		fmt.Println(".... Adding Server Certificate Success")
	}
	fmt.Println("#................... Server Certificate by Name ...............#")
	server_certn, err := ovc.GetServerCertificateByName(server_certificate_name)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(server_certn)
	}

	certificateDetails := new([]ov.CertificateDetail)
	certificateDetail_new := ov.CertificateDetail{Type: "CertificateDetailV2", AliasName: server_certificate_name, Base64Data: new_cert_base64data}
	*certificateDetails = append(*certificateDetails, certificateDetail_new)
	server_certn.CertificateDetails = *certificateDetails
	err = ovc.UpdateServerCertificate(server_certn)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("#.................... Server Certificate after Updating ...........#")
		server_cert_after_update, err := ovc.GetServerCertificateByName(server_certificate_name)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("..............Server Certificate Successfully updated.........")
			fmt.Println(server_cert_after_update)
		}
	}

	err = ovc.DeleteServerCertificate(server_certificate_name)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("#...................... Deleted Server Certificate Successfully .....#")
	}

}
