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

// Wrap structure to adding default values.
func NewDiscoveryImage() (*DiscoveryImage) {
	return &DiscoveryImage{
		CountDown: "5",
		PxAuto: "1",
		Ssh: "1",
		RootPw: "discovery_admin",
	}
}

func DiscoveryImageValidate(d *DiscoveryImage) (string, error) {
	// Validate FDQN.
	if !strings.Contains(d.HostName, ".") {
		log.Error().Msgf("Hostname is not a FQDN: %s", d.HostName)
		return "", errors.New("hostname is not a FQDN")
	}

	return "", nil
}

// Execute a command.
func executeCommand(wait bool, command string, args ...string) string {
	// Prepare command.
	cmd := exec.Command(command, args...)

	// If wait wait until command returns the result.
	if wait {
		output, err := cmd.Output()
		if err != nil {
			log.Error().Msgf("Error creating image: %s", err.Error())
			return ""
		}

		return string(output)
	}

	// If not wait start the command and return.
	if err := cmd.Start(); err != nil {
		log.Error().Msgf("Error creating image: %s", err.Error())
		return ""
	}

	return "Command is being executed."
}

// Convert facts list to string
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
func GenerateDiscoveryImage(d *DiscoveryImage, imagesPath string) string {
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

	return executeCommand(false, "bash", "-c", cmd)
}
