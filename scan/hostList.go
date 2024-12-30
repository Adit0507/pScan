package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

var (
	ErrExists    = errors.New("Host already in the list!")
	ErrNotExists = errors.New("Host is not in the list")
)

// represents a list of hosts to run port scan
type HostLists struct {
	Hosts []string
}

func (hl *HostLists) search(host string) (bool, int) {
	sort.Strings(hl.Hosts)

	i := sort.SearchStrings(hl.Hosts, host)
	if i < len(hl.Hosts) && hl.Hosts[i] == host {
		return true, i
	}

	return false, -1
}

func (hl *HostLists) Add(host string) error {
	// check if host already exists
	if found, _ := hl.search(host); found {
		return fmt.Errorf("%w: %s", ErrExists, host)
	}

	hl.Hosts = append(hl.Hosts, host)
	return nil
}

// delete given host from list
func (hl *HostLists) Remove(host string) error {
	if founnd, i := hl.search(host); founnd {
		hl.Hosts = append(hl.Hosts[:i], hl.Hosts[i+1:]...)
	}

	return fmt.Errorf("%w: %s", ErrNotExists, host)
}

// loads hosts from hostsFile
func (hl *HostLists) Load(hostsFile string) error {
	f, err := os.Open(hostsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts, scanner.Text())
	}

	return nil
}

func (hl *HostLists) Save(hostsFile string) error {
	output := ""

	for _, h := range hl.Hosts {
		output += fmt.Sprintln(h)
	}

	return os.WriteFile(hostsFile, []byte(output), 0644)
}
