//
// Test the go-chef/chef chef server api /organizations endpoints against a live server
//
package chefapi_client

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"

	chef "github.com/go-chef/chef"
)

// ChefInfo contains the values specified to use for a chefapi client connection
type chefInfo struct {
	user string
	keyfile string
	chefurl string
	certfile string
}

var flags chefInfo

// client for the chef server api
func Client() *chef.Client {
	// Pass in the chef-server api credentials.
	flagInit()

	// Create a client for access
	return buildClient(flags.user, flags.keyfile, flags.chefurl)
}

// client for the chef server api add organization
func OrgClient(organization string) *chef.Client {
	// Pass in the chef-server api credentials.
	flagInit()
	u, _ := url.Parse(flags.chefurl)
	u.Path = path.Join(u.Path, "organizations", organization)

	// Create a client for access
	fmt.Printf("ORGCLIENT PATH %+v\n", u.String()) // TODO: DEBUG
	return buildClient(flags.user, flags.keyfile, u.String() + "/")
}

// buildClient creates a connection to a chef server using the chef api.
// goiardi uses port 4545 by default, chef-zero uses 8889, chef-server uses 443
func buildClient(user string, keyfile string, baseurl string) *chef.Client {
	// TODO: debug msg
	fmt.Printf("KEYFILE path %+v\n",keyfile)
	key := clientKey(keyfile)
	client, err := chef.NewClient(&chef.Config{
		Name:    user,
		Key:     string(key),
		BaseURL: baseurl,
		SkipSSL: false,
		RootCAs: chefCerts(),
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue setting up client:", err)
		os.Exit(1)
	}
	return client
}

// clientKey reads the pem file containing the credentials needed to use the chef client.
func clientKey(filepath string) string {
	fmt.Printf("KEY path %+v\n",filepath)
	key, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read key.pem: %+v, %+v", filepath, err)
		os.Exit(1)
	}
	return string(key)
}

// chefCerts creats a cert pool for the self signed certs
// reference https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/
func chefCerts() *x509.CertPool {
	var certfile = flags.certfile
	certPool, _ := x509.SystemCertPool()
	if certPool == nil {
		certPool = x509.NewCertPool()
	}
	// Read in the cert file
	if certfile == "" {
		return certPool
	}
	certs, err := ioutil.ReadFile(certfile)
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", certfile, err)
	}
	// Append our cert to the system pool
	if ok := certPool.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}
	return certPool
}

func flagInit() {
	flags.user = os.Getenv("CHEFAPICHEFUSER")
        flags.keyfile = os.Getenv("CHEFAPIKEYFILE")
        flags.chefurl = os.Getenv("CHEFAPICHRURL")
        flags.certfile = os.Getenv("CHEFAPICERTFILE")
	fmt.Printf("flags %+v\n",flags)
	if flags.user == "" || flags.keyfile == "" || flags.chefurl == "" {
		log.Println("Env variables CHEFAPICHEFUSER, CHEFAPIKEYFILE, CHEFAPICHRURL must be set")
	}
	return
}
