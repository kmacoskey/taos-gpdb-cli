package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/kmacoskey/taos-gpdb-cli/request"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	// . "github.com/kmacoskey/taos-gpdb-cli"
)

var _ = Describe("Main", func() {

	var (
		pathToTaosCLI   string
		session         *gexec.Session
		configFile      *os.File
		configFilePath  string
		mockServer      *http.Server
		validPortConfig []byte
		validHostConfig []byte
	)

	BeforeSuite(func() {
		var err error

		// Mock taos server

		validServerPort := "8080"
		validServerHost := "localhost"

		router := mux.NewRouter()
		router.HandleFunc("/cluster/{id}", mockGetCluster).Methods("GET")
		mockServer = &http.Server{Addr: fmt.Sprintf(":%s", validServerPort), Handler: router}
		go func() {
			if err := mockServer.ListenAndServe(); err != nil {
				log.Printf("Httpserver: ListenAndServe() error: %s", err)
			}
		}()

		// CLI execution

		pathToTaosCLI, err = gexec.Build("github.com/kmacoskey/taos-gpdb-cli")
		Expect(err).NotTo(HaveOccurred())

		// CLI configuration

		configFile, err = TempFile(os.TempDir(), "taos-cli-test-config-file.*.yml")
		Expect(err).NotTo(HaveOccurred())
		configFilePath = configFile.Name()

		validPortConfig = []byte(fmt.Sprintf("port: %s", validServerPort))
		validHostConfig = []byte(fmt.Sprintf("host: %s", validServerHost))

	})

	AfterSuite(func() {
		var err error

		gexec.CleanupBuildArtifacts()
		err = mockServer.Shutdown(nil)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Getting a cluster", func() {
		Context("When everything goes ok", func() {
			BeforeEach(func() {
				var err error

				err = ioutil.WriteFile(configFilePath, validHostConfig, os.ModeAppend)
				Expect(err).NotTo(HaveOccurred())
				err = ioutil.WriteFile(configFilePath, validPortConfig, os.ModeAppend)
				Expect(err).NotTo(HaveOccurred())

				command := exec.Command(pathToTaosCLI, "--config", configFilePath, "get", "-i", "64a3a483-e273-48ba-a42a-75d93cd2d6a9")
				session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
			})
			AfterEach(func() {
				session.Kill()
				os.Remove(configFile.Name())
			})
			It("Should output the returned cluster", func() {
				Eventually(session.Out).Should(gbytes.Say(`{"request_id":"5676fe91-80c7-4bb7-a2ab-5c55c3f2b4e5","status":"200","data":{"type":"cluster","Attributes":{"id":"5676fe91-80c7-4bb7-a2ab-5c55c3f2b4e5","name":"Fake Cluster","status":"provision_success","message":"Apply complete! Resources: 0 added, 0 changed, 0 destroyed.","TerraformOutputs":{"foo":{"sensitive":false,"type":"string","value":"bar"}}}}}`))
			})
		})
	})

})

func mockGetCluster(w http.ResponseWriter, r *http.Request) {
	terraformOutputs := map[string]request.TerraformOutput{
		"foo": request.TerraformOutput{
			Sensitive: false,
			Type:      "string",
			Value:     "bar",
		},
	}

	clusterResponseAttributes := request.ClusterResponseAttributes{
		Id:               "5676fe91-80c7-4bb7-a2ab-5c55c3f2b4e5",
		Name:             "Fake Cluster",
		Status:           "provision_success",
		Message:          "Apply complete! Resources: 0 added, 0 changed, 0 destroyed.",
		TerraformOutputs: terraformOutputs,
	}

	clusterResponseData := request.ClusterResponseData{
		Type:       "cluster",
		Attributes: clusterResponseAttributes,
	}

	clusterResponse := request.ClusterResponse{
		RequestId: "5676fe91-80c7-4bb7-a2ab-5c55c3f2b4e5",
		Status:    "200",
		Data:      clusterResponseData,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(clusterResponse)
}

// ================================================
//  _                        __ _ _
// | |_ ___ _ __ ___  _ __  / _(_) | ___
// | __/ _ \ '_ ` _ \| '_ \| |_| | |/ _ \
// | ||  __/ | | | | | |_) |  _| | |  __/
//  \__\___|_| |_| |_| .__/|_| |_|_|\___|
//                   |_|
// ================================================

// This nice feature hasn't made it's way into golang yet:
// https://go-review.googlesource.com/c/go/+/105675

// Random number state.
// We generate random temporary file names so that there's a good
// chance the file doesn't exist yet - keeps the number of tries in
// TempFile to a minimum.
var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextRandom() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

// TempFile creates a new temporary file in the directory dir,
// opens the file for reading and writing, and returns the resulting *os.File.
// The filename is generated by taking pattern and adding a random
// string to the end. If pattern includes a "*", the random string
// replaces the last "*".
// If dir is the empty string, TempFile uses the default directory
// for temporary files (see os.TempDir).
// Multiple programs calling TempFile simultaneously
// will not choose the same file. The caller can use f.Name()
// to find the pathname of the file. It is the caller's responsibility
// to remove the file when no longer needed.
func TempFile(dir, pattern string) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	var prefix, suffix string
	if pos := strings.LastIndex(pattern, "*"); pos != -1 {
		prefix, suffix = pattern[:pos], pattern[pos+1:]
	} else {
		prefix = pattern
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, prefix+nextRandom()+suffix)
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	return
}
