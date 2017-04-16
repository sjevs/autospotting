package autospotting

import (
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

var logger, debug *log.Logger

// Run starts processing all AWS regions looking for AutoScaling groups
// enabled and taking action by replacing more pricy on-demand instances with
// compatible and cheaper spot instances.
func Run(cfg Config) {

	setupLogging(cfg)

	debug.Println(cfg)

	// use this only to list all the other regions
	ec2Conn := connectEC2(cfg.MainRegion)

	allRegions, err := getRegions(ec2Conn)

	if err != nil {
		logger.Println(err.Error())
		return
	}

	processRegions(allRegions, cfg)

}

func disableLogging() {
	setupLogging(Config{LogFile: ioutil.Discard})
}

func setupLogging(cfg Config) {
	logger = log.New(cfg.LogFile, "", cfg.LogFlag)

	if os.Getenv("AUTOSPOTTING_DEBUG") == "true" {
		debug = log.New(cfg.LogFile, "", cfg.LogFlag)
	} else {
		debug = log.New(ioutil.Discard, "", 0)
	}

}

// processAllRegions iterates all regions in parallel, and replaces instances
// for each of the ASGs tagged with 'spot-enabled=true'.
func processRegions(regions []string, cfg Config) {

	var wg sync.WaitGroup

	for _, r := range regions {

		wg.Add(1)
		r := region{name: r, conf: cfg}

		go func() {

			if r.enabled() {
				logger.Printf("Enabled to run in %s, processing region.\n", r.name)
				r.processRegion()
			} else {
				debug.Println("Not enabled to run in", r.name)
				debug.Println("List of enabled regions:", cfg.Regions)
			}

			wg.Done()
		}()
	}
	wg.Wait()
}

func connectEC2(region string) *ec2.EC2 {

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	return ec2.New(sess,
		aws.NewConfig().WithRegion(region))
}

// getRegions generates a list of AWS regions.
func getRegions(ec2conn ec2iface.EC2API) ([]string, error) {
	var output []string

	logger.Println("Scanning for available AWS regions")

	resp, err := ec2conn.DescribeRegions(&ec2.DescribeRegionsInput{})

	if err != nil {
		logger.Println(err.Error())
		return nil, err
	}

	debug.Println(resp)

	for _, r := range resp.Regions {

		if r != nil && r.RegionName != nil {
			debug.Println("Found region", *r.RegionName)
			output = append(output, *r.RegionName)
		}
	}
	return output, nil
}
