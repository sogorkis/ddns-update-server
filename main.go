package main

import (
	"context"
	"flag"
	"fmt"
	dns "google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var help = flag.Bool("help", false, "Show help")
var addr = flag.String("addr", ":8080", "Bind address, e.g. :8080")
var credentialsFilePath = flag.String("credentials_file_path", "", "JSON service account credentials file path")
var gcpProject = flag.String("gcp_project", "", "GCP project name")
var managedZone = flag.String("managed_zone", "", "DNS managed zone")
var rrsetName = flag.String("rrset_name", "", "DNS rrset name")
var rrsetType = flag.String("rrset_type", "", "DNS rrset type")

func rootPathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "dns-update-server is running")
}

func updateDnsPathHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")

	if ip == "" {
		http.Error(w, "FAILURE, details: missing ip", http.StatusBadRequest)
	} else {
		err := updateDns(r.Context(), ip)

		if err != nil {
			http.Error(w, fmt.Sprintf("FAILURE, details: %s", err), http.StatusInternalServerError)
		} else {
			fmt.Fprintf(w, "SUCCESS")
		}
	}
}
func updateDns(ctx context.Context, ip string) error {
	log.Printf("Updating rrset '%s', type '%s' to ip: '%s'", *rrsetName, *rrsetType, ip)

	client, err := dns.NewService(ctx, option.WithCredentialsFile(*credentialsFilePath))
	if err != nil {
		return err
	}

	_, err = client.ResourceRecordSets.Patch(*gcpProject, *managedZone, *rrsetName, *rrsetType, &dns.ResourceRecordSet{
		Rrdatas: []string{ip},
	}).Do()
	if err != nil {
		return err
	}

	log.Print("Update of DNS was successful")

	return nil
}

func main() {
	flag.Parse()

	// print help if requested
	if *help || *credentialsFilePath == "" || *gcpProject == "" || *managedZone == "" || *rrsetName == "" || *rrsetType == "" {
		flag.Usage()
		os.Exit(0)
	}

	//log startup messages
	log.Println("Starting ddns-update-server")

	http.HandleFunc("/", rootPathHandler)
	http.HandleFunc("/update-dns", updateDnsPathHandler)
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
