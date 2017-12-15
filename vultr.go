package vultr

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	vultr "github.com/JamesClonk/vultr/lib"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
)

type Driver struct {
	*drivers.BaseDriver
	APIKey            string
	APIEndpoint       string
	MachineID         string
	PrivateIP         string
	OSID              int
	RegionID          int
	PlanID            int
	SSHKeyID          string
	VultrPublicKey    string
	ROSVersion        string
	ReservedIP        string
	IPv6              bool
	Backups           bool
	PrivateNetworking bool
	PxeScriptID       int
	BootScriptID      int
	CustomPxeScript   bool
	UserDataFile      string
	SnapshotID        string
	VultrTag          string
	FirewallGroupID   string
	client            *vultr.Client
}

const (
	defaultOS          = 159
	defaultRegion      = 1
	defaultPlan        = 201
	defaultSSHuser     = "root"
	defaultROSVersion  = "v1.0.2"
	clientMaxRetries   = 4
	defaultAPIEndpoint = ""
)

// GetCreateFlags registers the flags this driver adds to
// "docker hosts create"
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "VULTR_API_KEY",
			Name:   "vultr-api-key",
			Usage:  "Vultr API key",
		},
		mcnflag.StringFlag{
			EnvVar: "VULTR_API_ENDPOINT",
			Name:   "vultr-api-endpoint",
			Usage:  "Vultr API endpoint",
			Value:  defaultAPIEndpoint,
		},
		mcnflag.StringFlag{
			EnvVar: "VULTR_SSH_USER",
			Name:   "vultr-ssh-user",
			Usage:  "Vultr SSH username",
			Value:  defaultSSHuser,
		},
		mcnflag.IntFlag{
			EnvVar: "VULTR_REGION",
			Name:   "vultr-region-id",
			Usage:  "Vultr region ID. Default: 1 (New Jersey)",
			Value:  defaultRegion,
		},
		mcnflag.IntFlag{
			EnvVar: "VULTR_PLAN",
			Name:   "vultr-plan-id",
			Usage:  "Vultr plan ID. Default: 201 (1024 MB RAM).",
			Value:  defaultPlan,
		},
		mcnflag.IntFlag{
			EnvVar: "VULTR_OS",
			Name:   "vultr-os-id",
			Usage:  "Vultr operating system ID. Default: RancherOS.",
			Value:  defaultOS,
		},
		mcnflag.StringFlag{
			EnvVar: "VULTR_ROS_VERSION",
			Name:   "vultr-ros-version",
			Usage:  "RancherOS version to use (eg. v0.6.0). Default: v1.0.2",
			Value:  defaultROSVersion,
		},
		mcnflag.IntFlag{
			EnvVar: "VULTR_PXE_SCRIPT",
			Name:   "vultr-pxe-script",
			Usage:  "ID of a PXE script in your Vultr account.",
		},
		mcnflag.IntFlag{
			EnvVar: "VULTR_BOOT_SCRIPT",
			Name:   "vultr-boot-script",
			Usage:  "ID of a boot script in your Vultr account. Mutually exclusive of --vultr-pxe-script.",
		},
		mcnflag.StringFlag{
			EnvVar: "VULTR_SSH_KEY",
			Name:   "vultr-ssh-key-id",
			Usage:  "ID of an existing SSH key in your Vultr account.",
		},
		mcnflag.StringFlag{
			EnvVar: "VULTR_RESERVED_IP",
			Name:   "vultr-reserved-ip",
			Usage:  "ID of a reserved IP in your Vultr account.",
		},
		mcnflag.BoolFlag{
			EnvVar: "VULTR_IPV6",
			Name:   "vultr-ipv6",
			Usage:  "Enable IPv6 for the VPS.",
		},
		mcnflag.BoolFlag{
			EnvVar: "VULTR_PRIVATE_NETWORKING",
			Name:   "vultr-private-networking",
			Usage:  "Enable private networking for the VPS.",
		},
		mcnflag.BoolFlag{
			EnvVar: "VULTR_BACKUPS",
			Name:   "vultr-backups",
			Usage:  "Enable automatic backups for the VPS.",
		},
		mcnflag.StringFlag{
			EnvVar: "VULTR_USERDATA",
			Name:   "vultr-userdata",
			Usage:  "Path to a file containing cloud-init user data.",
		},
		mcnflag.StringFlag{
			EnvVar: "VULTR_SNAPSHOT",
			Name:   "vultr-snapshot-id",
			Usage:  "ID of an existing Snapshot in your Vultr account.",
		},
		mcnflag.StringFlag{
			EnvVar: "VULTR_TAG",
			Name:   "vultr-tag",
			Usage:  "Tag to assign to the VPS.",
		},
		mcnflag.StringFlag{
			EnvVar: "VULTR_FIREWALL_GROUP",
			Name:   "vultr-firewall-group",
			Usage:  "ID of existing firewall group to assign.",
		},
	}
}

func NewDriver(hostName, storePath string) *Driver {
	d := &Driver{
		OSID:     defaultOS,
		PlanID:   defaultPlan,
		RegionID: defaultRegion,
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
	}
	return d
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return "vultr"
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.APIKey = flags.String("vultr-api-key")
	d.APIEndpoint = flags.String("vultr-api-endpoint")
	d.OSID = flags.Int("vultr-os-id")
	d.ROSVersion = flags.String("vultr-ros-version")
	d.RegionID = flags.Int("vultr-region-id")
	d.PlanID = flags.Int("vultr-plan-id")
	d.PxeScriptID = flags.Int("vultr-pxe-script")
	d.BootScriptID = flags.Int("vultr-boot-script")
	d.SSHKeyID = flags.String("vultr-ssh-key-id")
	d.ReservedIP = flags.String("vultr-reserved-ip")
	d.IPv6 = flags.Bool("vultr-ipv6")
	d.PrivateNetworking = flags.Bool("vultr-private-networking")
	d.Backups = flags.Bool("vultr-backups")
	d.UserDataFile = flags.String("vultr-userdata")
	d.SnapshotID = flags.String("vultr-snapshot-id")
	d.VultrTag = flags.String("vultr-tag")
	d.FirewallGroupID = flags.String("vultr-firewall-group")
	d.SwarmMaster = flags.Bool("swarm-master")
	d.SwarmHost = flags.String("swarm-host")
	d.SwarmDiscovery = flags.String("swarm-discovery")
	d.SSHUser = flags.String("vultr-ssh-user")
	d.SSHPort = 22

	if d.APIKey == "" {
		return fmt.Errorf("Vultr driver requires the --vultr-api-key option")
	}
	return nil
}

func (d *Driver) PreCreateCheck() error {
	if d.UserDataFile != "" {
		if d.OSID == 159 {
			return fmt.Errorf("--vultr-userdata does currently not support 'Custom OS' (OS ID 159)")
		}
		if _, err := os.Stat(d.UserDataFile); os.IsNotExist(err) {
			return fmt.Errorf("Unable to find user data file at %s", d.UserDataFile)
		}
	}

	log.Info("Validating Vultr VPS parameters...")

	if d.PxeScriptID != 0 && d.BootScriptID != 0 {
		return fmt.Errorf("--vultr-pxe-script and --vultr-boot-script are mutually exclusive")
	}

	if d.PxeScriptID != 0 && d.OSID != 159 {
		return fmt.Errorf("--vultr-pxe-script requires the 'Custom OS' (OS ID 159)")
	}

	if d.BootScriptID != 0 && d.OSID == 159 {
		return fmt.Errorf("--vultr-boot-script can't be used with the 'Custom OS' (OS ID 159)")
	}

	if d.SnapshotID != "" && d.OSID == defaultOS {
		//	reassign OSID to Snapshot OSID 164, if OSID is the defaultOS.
		//	And allow user to specify an OSID, in case there is an API update in the future.
		d.OSID = 164
	}

	if d.SSHKeyID != "" {
		key, err := d.getPublicKeyByID(d.SSHKeyID)
		if err != nil {
			return err
		}

		log.Infof("Using existing SSH public key: %s", key.Name)
		d.VultrPublicKey = key.Key
	}

	if err := d.validateRegion(); err != nil {
		return err
	}

	if err := d.validatePlan(); err != nil {
		return err
	}

	if err := d.validateApiCredentials(); err != nil {
		return err
	}

	return nil
}

func (d *Driver) Create() error {
	if d.SSHKeyID == "" {
		log.Debug("Generating SSH key...")
		key, err := d.createSSHKey()
		if err != nil {
			return err
		}
		d.SSHKeyID = key.ID
	}

	log.Info("Creating Vultr VPS")
	var userdata string
	var err error
	if d.OSID == 159 {
		log.Info("Using PXE boot")
		if d.PxeScriptID != 0 {
			d.CustomPxeScript = true
		} else {
			log.Infof("Provisioning RancherOS (%s). SSH user set to 'rancher'.", d.ROSVersion)
			d.SSHUser = "rancher"
			if err := d.createBootScript(); err != nil {
				return err
			}

			log.Debugf("Created RancherOS PXE script: ID %d", d.PxeScriptID)
		}

		userdata, err = d.getCloudConfig()
		if err != nil {
			return err
		}
	} else if d.UserDataFile != "" {
		buf, err := ioutil.ReadFile(d.UserDataFile)
		if err != nil {
			return err
		}
		userdata = string(buf)
	}

	if userdata != "" {
		log.Debugf("Using the following cloud-init user data:")
		log.Debugf("%s", userdata)
	}

	var scriptID int
	switch {
	case d.PxeScriptID != 0:
		scriptID = d.PxeScriptID
	case d.BootScriptID != 0:
		scriptID = d.BootScriptID
	}

	client := d.getClient()
	machine, err := client.CreateServer(
		d.MachineName,
		d.RegionID,
		d.PlanID,
		d.OSID,
		&vultr.ServerOptions{
			SSHKey:               d.SSHKeyID,
			IPV6:                 d.IPv6,
			PrivateNetworking:    d.PrivateNetworking,
			AutoBackups:          d.Backups,
			Script:               scriptID,
			UserData:             userdata,
			Snapshot:             d.SnapshotID,
			Hostname:             d.MachineName,
			DontNotifyOnActivate: true,
			Tag:                  d.VultrTag,
			FirewallGroupID:      d.FirewallGroupID,
			ReservedIP:           d.ReservedIP,
		})
	if err != nil {
		return err
	}

	d.MachineID = machine.ID
	log.Info("Waiting for IP address to become available...")
	for {
		machine, err = client.GetServer(d.MachineID)
		if err != nil {
			return err
		}
		d.IPAddress = machine.MainIP
		d.PrivateIP = machine.InternalIP

		if d.IPAddress != "" && d.IPAddress != "0" && d.IPAddress != "0.0.0.0" {
			break
		}
		log.Debug("IP address not yet available")
		time.Sleep(2 * time.Second)
	}

	if d.PrivateIP == "0" {
		d.PrivateIP = ""
	}

	log.Infof("Created Vultr VPS ID: %s, Public IP: %s, Private IP: %s",
		d.MachineID,
		d.IPAddress,
		d.PrivateIP,
	)

	return nil
}

func (d *Driver) getPublicKeyByID(id string) (*vultr.SSHKey, error) {
	keys, err := d.getClient().GetSSHKeys()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		if key.ID == id {
			return &key, nil
		}
	}

	return nil, fmt.Errorf("Vultr SSH key with ID %s doesn't exist", id)
}

func (d *Driver) createSSHKey() (*vultr.SSHKey, error) {
	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return nil, err
	}

	publicKey, err := ioutil.ReadFile(d.publicSSHKeyPath())
	if err != nil {
		return nil, err
	}

	key, err := d.getClient().CreateSSHKey(d.MachineName, string(publicKey))
	if err != nil {
		return &key, err
	}

	return &key, nil
}

func (d *Driver) GetURL() (string, error) {
	s, err := d.GetState()
	if err != nil {
		return "", err
	}

	if s != state.Running {
		return "", drivers.ErrHostIsNotRunning
	}

	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

func (d *Driver) GetIP() (string, error) {
	if d.IPAddress == "" || d.IPAddress == "0" || d.IPAddress == "0.0.0.0" {
		return "", fmt.Errorf("IP address is not set")
	}

	return d.IPAddress, nil
}

func (d *Driver) GetState() (state.State, error) {
	machine, err := d.getClient().GetServer(d.MachineID)
	if err != nil {
		return state.Error, err
	}

	switch machine.Status {
	case "pending":
		return state.Starting, nil
	case "active":
		switch machine.ServerState {
		case "ok":
			switch machine.PowerStatus {
			case "running":
				return state.Running, nil
			case "stopped":
				return state.Stopped, nil
			}
		default:
			return state.Starting, nil
		}
	}
	return state.None, nil
}

func (d *Driver) Start() error {
	if vmState, err := d.GetState(); err != nil {
		return err
	} else if vmState == state.Running || vmState == state.Starting {
		log.Infof("Host is already running or starting")
		return nil
	}

	log.Debugf("starting %s", d.MachineName)
	return d.getClient().StartServer(d.MachineID)
}

func (d *Driver) Stop() error {
	if vmState, err := d.GetState(); err != nil {
		return err
	} else if vmState == state.Stopped {
		log.Infof("Host is already stopped")
		return nil
	}

	log.Debugf("stopping %s", d.MachineName)
	return d.getClient().HaltServer(d.MachineID)
}

func (d *Driver) Remove() error {
	client := d.getClient()
	log.Debugf("removing %s", d.MachineName)
	if err := client.DeleteServer(d.MachineID); err != nil {
		if strings.Contains(err.Error(), "Invalid server") {
			log.Infof("VPS doesn't exist, assuming it is already deleted")
		} else {
			return err
		}
	}

	if d.PxeScriptID != 0 && !d.CustomPxeScript {
		if err := client.DeleteStartupScript(strconv.Itoa(d.PxeScriptID)); err != nil {
			if strings.Contains(err.Error(), "Check SCRIPTID") {
				log.Infof("PXE script doesn't exist, assuming it is already deleted")
			} else {
				return err
			}
		}
	}

	if d.VultrPublicKey == "" {
		if err := client.DeleteSSHKey(d.SSHKeyID); err != nil {
			if strings.Contains(err.Error(), "Invalid SSH Key") {
				log.Infof("SSH key doesn't exist, assuming it is already deleted")
			} else {
				return err
			}
		}
	}

	return nil
}

func (d *Driver) Restart() error {
	if vmState, err := d.GetState(); err != nil {
		return err
	} else if vmState == state.Stopped {
		log.Infof("Host is already stopped, use start command to run it")
		return nil
	}

	log.Debugf("restarting %s", d.MachineName)
	return d.getClient().RebootServer(d.MachineID)
}

func (d *Driver) Kill() error {
	if vmState, err := d.GetState(); err != nil {
		return err
	} else if vmState == state.Stopped {
		log.Infof("Host is already stopped")
		return nil
	}

	log.Debugf("killing %s", d.MachineName)
	return d.getClient().HaltServer(d.MachineID)
}

func (d *Driver) getClient() *vultr.Client {
	if d.client == nil {
		opts := &vultr.Options{
			MaxRetries: clientMaxRetries,
			Endpoint:   d.APIEndpoint,
		}
		d.client = vultr.NewClient(d.APIKey, opts)
	}

	return d.client
}

func (d *Driver) publicSSHKeyPath() string {
	return d.GetSSHKeyPath() + ".pub"
}

func (d *Driver) GetSSHKeyPath() string {
	// don't set SSHKeyPath when using an existing SSH key
	if d.SSHKeyPath == "" && d.VultrPublicKey == "" {
		d.SSHKeyPath = d.ResolveStorePath("id_rsa")
	}

	return d.SSHKeyPath
}

func (d *Driver) instanceIsRunning() bool {
	st, err := d.GetState()
	if err != nil {
		log.Debug(err)
	}

	if st == state.Running {
		return true
	}

	log.Debug("VPS not yet started")
	return false
}

func (d *Driver) validateApiCredentials() error {
	_, err := d.getClient().GetAccountInfo()
	if err != nil {
		return err
	}

	return nil
}

func (d *Driver) validateRegion() error {
	regions, err := d.getClient().GetRegions()
	if err != nil {
		return err
	}

	for _, region := range regions {
		if region.ID == d.RegionID {
			return nil
		}
	}

	return fmt.Errorf("Region ID %d is invalid", d.RegionID)
}

func (d *Driver) validatePlan() error {
	plans, err := d.getClient().GetAvailablePlansForRegion(d.RegionID)
	if err != nil {
		return err
	}

	for _, v := range plans {
		if v == d.PlanID {
			return nil
		}
	}

	return fmt.Errorf("PlanID %d not available in the chosen region. Available plans for RegionID %d: %v", d.PlanID, d.RegionID, plans)
}

// RancherOS - Create iPXE script
func (d *Driver) createBootScript() error {
	content := `#!ipxe
set base-url http://releases.rancher.com/os/%s
kernel ${base-url}/vmlinuz rancher.state.dev=LABEL=RANCHER_STATE rancher.state.autoformat=[/dev/vda] rancher.state.formatzero rancher.cloud_init.datasources=[ec2]
initrd ${base-url}/initrd
boot`

	content = fmt.Sprintf(content, d.ROSVersion)
	log.Debugf("Using the following PXE script:")
	log.Debugf("%s", content)
	script, err := d.getClient().CreateStartupScript(d.MachineName, content, "pxe")
	if err != nil {
		return err
	}
	d.PxeScriptID, err = strconv.Atoi(script.ID)
	if err != nil {
		return err
	}
	return nil
}

// RancherOS - Generate cloud-config userdata string that will
// provision the SSH Key to the VPS and configure private networking
func (d *Driver) getCloudConfig() (string, error) {
	type userData struct {
		HostName     string
		SSHkey       string
		PrivateNet   bool
		CustomScript bool
	}

	const tpl = `#cloud-config
hostname: {{.HostName}}
ssh_authorized_keys:
  - {{.SSHkey}}{{if not .CustomScript}}
write_files:
  - path: /opt/rancher/bin/start.sh
    permissions: "0755"
    owner: root
    content: |
      #!/bin/sh
      mount | grep /dev/vda >/dev/null
      RETVAL=$?
      if [ $RETVAL -eq 0 ]; then
        exit 0
      fi
      sudo dd if=/dev/zero of=/dev/vda bs=1M count=1
      logger -t start.sh "Prepared /dev/vda for use as Rancher state disk. Rebooting."
      sudo reboot
rancher:
  network:
    interfaces:
      eth0:
        dhcp: true{{if .PrivateNet}}
      eth1:
        address: $private_ipv4/16
        mtu: 1450{{end}}{{end}}
`
	var buffer bytes.Buffer
	var publicKey string

	if d.VultrPublicKey != "" {
		publicKey = d.VultrPublicKey
	} else {
		keyByte, err := ioutil.ReadFile(d.publicSSHKeyPath())
		if err != nil {
			return "", err
		}
		publicKey = string(keyByte)
	}

	config := userData{HostName: d.MachineName, SSHkey: publicKey, PrivateNet: d.PrivateNetworking, CustomScript: d.CustomPxeScript}
	tmpl, err := template.New("cloud-config").Parse(tpl)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&buffer, config)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
