package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

type NCSIReg struct {
	ActiveWebProbeHost      *string
	ActiveWebProbePath      *string
	ActiveWebProbeContent   *string
	ActiveWebProbeHostV6    *string
	ActiveWebProbePathV6    *string
	ActiveWebProbeContentV6 *string
	ActiveDnsProbeHost      *string
	ActiveDnsProbeContent   *string
	ActiveDnsProbeHostV6    *string
	ActiveDnsProbeContentV6 *string
	EnableActiveProbing     *uint64
}

func (n *NCSIReg) Print() {
	fmt.Println("-------------------------------------------------")
	fmt.Println("           System NCSI registry values")
	fmt.Println("-------------------------------------------------")
	fmt.Println(fmt.Sprintf("[ActiveWebProbeHost]: %#v", *n.ActiveWebProbeHost))
	fmt.Println(fmt.Sprintf("[ActiveWebProbePath]: %#v", *n.ActiveWebProbePath))
	fmt.Println(fmt.Sprintf("[ActiveWebProbeContent]: %#v", *n.ActiveWebProbeContent))
	fmt.Println(fmt.Sprintf("[ActiveWebProbeHostV6]: %#v", *n.ActiveWebProbeHostV6))
	fmt.Println(fmt.Sprintf("[ActiveWebProbePathV6]: %#v", *n.ActiveWebProbePathV6))
	fmt.Println(fmt.Sprintf("[ActiveWebProbeContentV6]: %#v", *n.ActiveWebProbeContentV6))
	fmt.Println(fmt.Sprintf("[ActiveDnsProbeHost]: %#v", *n.ActiveDnsProbeHost))
	fmt.Println(fmt.Sprintf("[ActiveDnsProbeContent]: %#v", *n.ActiveDnsProbeContent))
	fmt.Println(fmt.Sprintf("[ActiveDnsProbeHostV6]: %#v", *n.ActiveDnsProbeHostV6))
	fmt.Println(fmt.Sprintf("[ActiveDnsProbeContentV6]: %#v", *n.ActiveDnsProbeContentV6))
	fmt.Println(fmt.Sprintf("[EnableActiveProbing]: %d", *n.EnableActiveProbing))
	fmt.Println("-------------------------------------------------")
}

func (n *NCSIReg) ToServer() *Server {
	return &Server{
		ActiveWebProbeHost:      *n.ActiveWebProbeHost,
		ActiveWebProbePath:      *n.ActiveWebProbePath,
		ActiveWebProbeContent:   *n.ActiveWebProbeContent,
		ActiveWebProbeHostV6:    *n.ActiveWebProbeHostV6,
		ActiveWebProbePathV6:    *n.ActiveWebProbePathV6,
		ActiveWebProbeContentV6: *n.ActiveWebProbeContentV6,
		ActiveDnsProbeHost:      *n.ActiveDnsProbeHost,
		ActiveDnsProbeContent:   *n.ActiveDnsProbeContent,
		ActiveDnsProbeHostV6:    *n.ActiveDnsProbeHostV6,
		ActiveDnsProbeContentV6: *n.ActiveDnsProbeContentV6,
	}
}

func (n *NCSIReg) setSystemNCSIReg() error {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\NlaSvc\Parameters\Internet`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()
	if n.ActiveWebProbeHost != nil {
		err = key.SetStringValue("ActiveWebProbeHost", *n.ActiveWebProbeHost)
		if err != nil {
			return err
		}
	}
	if n.ActiveWebProbePath != nil {
		err = key.SetStringValue("ActiveWebProbePath", *n.ActiveWebProbePath)
		if err != nil {
			return err
		}
	}
	if n.ActiveWebProbeContent != nil {
		err = key.SetStringValue("ActiveWebProbeContent", *n.ActiveWebProbeContent)
		if err != nil {
			return err
		}
	}
	if n.ActiveWebProbeHostV6 != nil {
		err = key.SetStringValue("ActiveWebProbeHostV6", *n.ActiveWebProbeHostV6)
		if err != nil {
			return err
		}
	}
	if n.ActiveWebProbePathV6 != nil {
		err = key.SetStringValue("ActiveWebProbePathV6", *n.ActiveWebProbePathV6)
		if err != nil {
			return err
		}
	}
	if n.ActiveWebProbeContentV6 != nil {
		err = key.SetStringValue("ActiveWebProbeContentV6", *n.ActiveWebProbeContentV6)
		if err != nil {
			return err
		}
	}
	if n.ActiveDnsProbeHost != nil {
		err = key.SetStringValue("ActiveDnsProbeHost", *n.ActiveDnsProbeHost)
		if err != nil {
			return err
		}
	}
	if n.ActiveDnsProbeContent != nil {
		err = key.SetStringValue("ActiveDnsProbeContent", *n.ActiveDnsProbeContent)
		if err != nil {
			return err
		}
	}
	if n.ActiveDnsProbeHostV6 != nil {
		err = key.SetStringValue("ActiveDnsProbeHostV6", *n.ActiveDnsProbeHostV6)
		if err != nil {
			return err
		}
	}
	if n.ActiveDnsProbeContentV6 != nil {
		err = key.SetStringValue("ActiveDnsProbeContentV6", *n.ActiveDnsProbeContentV6)
		if err != nil {
			return err
		}
	}
	// EnableActiveProbing must be 1
	err = key.SetDWordValue("EnableActiveProbing", 1)
	if err != nil {
		return err
	}
	return nil
}

func GetSystemNCSIReg() (*NCSIReg, error) {
	var ncsiReg NCSIReg
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\NlaSvc\Parameters\Internet`, registry.QUERY_VALUE)
	if err != nil {
		return nil, err
	}
	defer key.Close()

	activeWebProbeHost, _, err := key.GetStringValue("ActiveWebProbeHost")
	if err == nil {
		ncsiReg.ActiveWebProbeHost = &activeWebProbeHost
	}
	activeWebProbePath, _, err := key.GetStringValue("ActiveWebProbePath")
	if err == nil {
		ncsiReg.ActiveWebProbePath = &activeWebProbePath
	}
	activeWebProbeContent, _, err := key.GetStringValue("ActiveWebProbeContent")
	if err == nil {
		ncsiReg.ActiveWebProbeContent = &activeWebProbeContent
	}
	activeWebProbeHostV6, _, err := key.GetStringValue("ActiveWebProbeHostV6")
	if err == nil {
		ncsiReg.ActiveWebProbeHostV6 = &activeWebProbeHostV6
	}
	activeWebProbePathV6, _, err := key.GetStringValue("ActiveWebProbePathV6")
	if err == nil {
		ncsiReg.ActiveWebProbePathV6 = &activeWebProbePathV6
	}
	activeWebProbeContentV6, _, err := key.GetStringValue("ActiveWebProbeContentV6")
	if err == nil {
		ncsiReg.ActiveWebProbeContentV6 = &activeWebProbeContentV6
	}
	activeDnsProbeHost, _, err := key.GetStringValue("ActiveDnsProbeHost")
	if err == nil {
		ncsiReg.ActiveDnsProbeHost = &activeDnsProbeHost
	}
	activeDnsProbeContent, _, err := key.GetStringValue("ActiveDnsProbeContent")
	if err == nil {
		ncsiReg.ActiveDnsProbeContent = &activeDnsProbeContent
	}
	activeDnsProbeHostV6, _, err := key.GetStringValue("ActiveDnsProbeHostV6")
	if err == nil {
		ncsiReg.ActiveDnsProbeHostV6 = &activeDnsProbeHostV6
	}
	activeDnsProbeContentV6, _, err := key.GetStringValue("ActiveDnsProbeContentV6")
	if err == nil {
		ncsiReg.ActiveDnsProbeContentV6 = &activeDnsProbeContentV6
	}
	enableActiveProbing, _, err := key.GetIntegerValue("EnableActiveProbing")
	if err == nil {
		ncsiReg.EnableActiveProbing = &enableActiveProbing
	}
	return &ncsiReg, nil
}
