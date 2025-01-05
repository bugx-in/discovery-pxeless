package fpi

import (
	"errors"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Fact struct {
	FactName  string `json:"name"`
	FactValue string `json:"value"`
}

type DiscoveryImage struct {
	HostName  string `json:"host_name" binding:"required"`
	Ip        string `json:"ip" binding:"required"`
	Cidr      string `json:"cidr" binding:"required"`
	Gateway   string `json:"gateway" binding:"required"`
	Dns       string `json:"dns"`
	ProxyUrl  string `json:"proxy_url"`
	ProxyType string `json:"proxy_type"`
	CountDown string `json:"countdown"`
	Ssh       string `json:"ssh_enabled"`
	RootPw    string `json:"ssh_password"`
	PxAuto    string `json:"auto"`
	Facts     []Fact `json:"facts"`
}

// New discovery image struct with default values.
func NewDiscoveryImage() *DiscoveryImage {
	return &DiscoveryImage{
		CountDown: "5",
		PxAuto:    "1",
		Ssh:       "1",
		RootPw:    "discovery_admin",
	}
}

func (d *DiscoveryImage) ValidateImage() (string, error) {
	// Validate FDQN.
	if !strings.Contains(d.HostName, ".") {
		log.Error().Msgf("Hostname is not a FQDN: %s", d.HostName)
		return "", errors.New("hostname is not a FQDN")
	}
	return "", nil
}

// Execute a command.
func executeCommand(args ...string) (string, error) {
	// Prepare command.
	//lint:ignore SA1005
	cmd := exec.Command("bash -c", args...)

	// Do not wait.
	if err := cmd.Start(); err != nil {
		log.Error().Msgf("Error creating image: %s", err.Error())
		return "", err
	}

	return "Command is being executed.", nil
}

// Convert facts list to string.
func formatFacts(f []Fact) string {
	var facts string

	for index, fact := range f {
		// Add fact name and value
		facts = facts + " fdi.pxfactname" + strconv.Itoa(index+1) + "=" + fact.FactName + " "
		facts = facts + "fdi.pxfactvalue" + strconv.Itoa(index+1) + "=" + fact.FactValue
	}

	return facts
}

// Generate a discovery image.
func (d *DiscoveryImage) GenerateDiscoveryImage(imagesPath string) (string, error) {
	// Get the remaster script and the ISO base image.
	// The ISO base image will be remastered to add the host details.
	remasterBin := os.Getenv("DISCOVERY_REMASTER")
	discoveryBaseImage := os.Getenv("DISCOVERY_BASE_IMAGE")

	// Generate the command to execute to create the image.
	cmdOpts := []string{remasterBin,
		discoveryBaseImage,
		"\"" +
			"fdi.pxip=" + d.Ip + "/" + d.Cidr,
		"fdi.pxgw=" + d.Gateway,
		"fdi.pxdns=" + d.Dns,
		"proxy.url=" + d.ProxyUrl,
		"proxy.type=" + d.ProxyType,
		"fdi.countdown=" + d.CountDown,
		"fdi.pxauto=" + d.PxAuto,
		"fdi.ssh=" + d.Ssh,
		"fdi.rootpwd=" + d.RootPw,
		"facts=" + formatFacts(d.Facts) +
			"\"",
		imagesPath + "/" + d.HostName + ".iso"}

	cmd := strings.Join(cmdOpts, " ")

	// Generate the image.
	log.Debug().Msgf("Generatin ISO image with command: %s", cmd)

	return executeCommand(cmd)
}
