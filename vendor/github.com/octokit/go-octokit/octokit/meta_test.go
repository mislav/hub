package octokit

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMeta(t *testing.T) {
	setup()
	defer tearDown()

	stubGet(t, "/meta", "meta", nil)
	info, result := client.Meta(&MetaURL)

	assert.False(t, result.HasError())

	assert.True(t, info.VerifiablePasswordAuthentication)
	assert.Equal(t, info.GithubServicesSha, "2e886f407696261bd5adfc99b16d36d5e7b50241")

	_, ipNet, _ := net.ParseCIDR("192.30.252.0/22")
	nets := []*net.IPNet{ipNet}
	assert.Equal(t, info.Hooks, nets)

	// example Git in meta.json is the same as Hooks
	assert.Equal(t, info.Git, nets)

	nets = nets[:0]
	ss := []string{"192.30.252.153/32", "192.30.252.154/32"}
	for _, s := range ss {
		_, ipNet, _ := net.ParseCIDR(s)
		nets = append(nets, ipNet)
	}
	assert.Equal(t, info.Pages, nets)

	ss = []string{
		"54.80.154.161",
		"54.80.168.241",
		"54.196.81.106",
		"54.158.161.132",
		"54.226.70.38",
	}
	ipSlice := make([]net.IP, 0, len(ss))
	for _, s := range ss {
		ip := net.ParseIP(s)
		ipSlice = append(ipSlice, ip)
	}
	assert.Equal(t, info.Importer, ipSlice)

	//Error case
	var invalid = Hyperlink("{")
	metaErr, resultErr := client.Meta(&invalid)
	assert.True(t, resultErr.HasError())
	assert.Equal(t, metaErr, APIInfo{})
}
