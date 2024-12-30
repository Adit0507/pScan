package cmd

import (
	"bytes"
	"fmt"
	"github.com/Adit0507/pScan.com/scan"
	"io"
	"os"
	"strings"
	"testing"
)

func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	tf, err := os.CreateTemp("", "pScan")
	if err != nil {
		t.Fatal(err)
	}
	tf.Close()

	if initList {
		hl := &scan.HostLists{}
		for _, h := range hosts {
			hl.Add(h)
		}
		if err := hl.Save(tf.Name()); err != nil {
			t.Fatal(err)
		}
	}

	return tf.Name(), func() {
		os.Remove(tf.Name())
	}
}

func TestHostActions(t *testing.T) {
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	testCases := []struct {
		name           string
		args           []string
		expectedOut    string
		initList       bool
		actionFunction func(io.Writer, string, []string) error
	}{
		{
			name:           "AddAction",
			args:           hosts,
			expectedOut:    "Added host: host1\n Added host: host2 \n Added host: host3\n ",
			initList:       false,
			actionFunction: addAction,
		},
		{
			name:           "ListAction",
			expectedOut:    "host1 \n host2 \n host3 \n",
			initList:       true,
			actionFunction: listAction,
		},
		{
			name:           "DeleteAction",
			args:           []string{"host1", "host2"},
			expectedOut:    "Deleted host: host1 \n Deleted host: host2 \n ",
			initList:       true,
			actionFunction: deleteAction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tf, cleanup := setup(t, hosts, tc.initList)
			defer cleanup()
			
			var out bytes.Buffer
			if err := tc.actionFunction(&out, tf, tc.args); err != nil {
				t.Fatalf("expected no error, got %q\n", err)
			}

			if out.String() != tc.expectedOut{
				t.Errorf("expected output %q, got %q\n", tc.expectedOut, out.String())
			}

		})
	}
}

func TestIntegration(t *testing.T){
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	tf, cleanup := setup(t, hosts, false)
	defer cleanup()

	delHost := "host2"
	hostsEnd := []string{
		"host1",
		"host3",
	}

	var out bytes.Buffer
	expectedOut := ""
	for _, v := range hosts {
		expectedOut += fmt.Sprintf("Added host: %s\n", v)
	}
	expectedOut += strings.Join(hosts, "\n")
	expectedOut += fmt.Sprintln()
	expectedOut += fmt.Sprintf("Deleted host: %s\n", delHost)
	expectedOut += strings.Join(hostsEnd, "\n")
	expectedOut += fmt.Sprintln()

	// add hosts to list
	if err := addAction(&out, tf, hosts); err != nil {
		t.Fatalf("Expected no eror, got %q\n", err)
	}
	// list hosts
	if err := listAction(&out, tf, hosts); err != nil {
		t.Fatalf("Expected no eror, got %q\n", err)
	}

	// delete host2
	if err := deleteAction(&out, tf, []string{delHost}); err != nil {
		t.Fatalf("Expected no eror, got %q\n", err)
	}

	// list hosts after delete
	if err := listAction(&out, tf, hosts); err != nil {
		t.Fatalf("Expected no eror, got %q\n", err)
	}

	if out.String() != expectedOut {
		t.Errorf("Expected output %q, got %q\n", expectedOut, out.String())
	}


}