package main

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestInMemory(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// Setup
	err := os.Setenv("DB", string(dbMemory))
	assert.NoError(t, err)

	// Start tests
	e2eTest(t, "localhost:9000")

	log.Info("Finished in memory test")
	// No cleanup needed
}

func TestDisk(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// Setup
	err := os.Setenv("DB", string(dbDisk))
	assert.NoError(t, err)
	os.TempDir()
	err = os.Setenv("DISKPATH", os.TempDir()+"/counter-test-disk-data")
	assert.NoError(t, err)

	// Start tests
	e2eTest(t, "localhost:9001")

	log.Info("Finished on disk test")
	err = os.RemoveAll(os.TempDir() + "/counter-test-disk-data")
	assert.NoError(t, err)
}

func TestBadger(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// Setup
	err := os.Setenv("DB", string(dbBadger))
	assert.NoError(t, err)
	os.TempDir()
	err = os.Setenv("DISKPATH", os.TempDir()+"/counter-test-disk-data")
	assert.NoError(t, err)

	// Start tests
	e2eTest(t, "localhost:9002")

	log.Info("Finished badger test")
	err = os.RemoveAll(os.TempDir() + "/counter-test-badger-data")
	assert.NoError(t, err)
}

func TestEtcd(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// You can run etcd as following:
	// docker run --network host  gcr.io/etcd-development/etcd
	// then set ETCDHOST=localhost:2379

	host := os.Getenv("ETCDHOST")
	if host == "" {
		t.Skip("Skipping etcd test as DBHOST is not set up")
	}

	// Setup
	err := os.Setenv("DB", string(dbEtcd3))
	assert.NoError(t, err)
	err = os.Setenv("DBHOST", host)
	assert.NoError(t, err)

	// Start tests
	e2eTest(t, "localhost:9003")

	log.Info("Finished etcd test")
	// No cleanup needed
}

func e2eTest(t *testing.T, addr string) {
	// ensure env is set before calling this func
	err := os.Setenv("ADDRESS", addr)
	assert.NoError(t, err)
	url := "http://" + addr

	gitHash = "E2E in Progress"

	// Startup the main program
	go main()

	time.Sleep(time.Second)

	// http client
	client := &http.Client{}

	// Test root for checking the server is working
	req, err := http.NewRequest(http.MethodGet, url+"/", nil)
	assert.NoError(t, err)
	resp, err := client.Do(req)
	assert.NoError(t, err)

	bytes, err := ioutil.ReadAll(resp.Body)
	text := string(bytes)
	assert.NoError(t, err)
	assert.Contains(t, text, gitHash)

	// Get NX Counter
	req, err = http.NewRequest(http.MethodGet, url+"/test/yeet", nil)
	assert.NoError(t, err)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Create Counter
	req, err = http.NewRequest(http.MethodPost, url+"/test/yeet", nil)
	assert.NoError(t, err)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	token := resp.Header.Get("Authorization")
	assert.True(t, strings.HasPrefix(token, "Bearer "))
	token = strings.TrimPrefix(token, "Bearer ")
	id, err := uuid.Parse(token)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	text = string(body)
	assert.Contains(t, text, "/test/yeet")
	assert.Contains(t, text, "0")
	assert.Contains(t, text, token)

	// Get 0 Counter
	resp, err = http.Get(url + "/test/yeet")
	assert.NoError(t, err)

	assert.Empty(t, resp.Header.Get("Authorization"))

	bytes, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	text = string(bytes)

	assert.Contains(t, text, "/test/yeet")
	assert.Contains(t, text, "0")
	assert.NotContains(t, text, token)

	// Increment counter without token
	req, err = http.NewRequest(http.MethodPut, url+"/test/yeet", nil)
	assert.NoError(t, err)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Increment counter without correctly
	req, err = http.NewRequest(http.MethodPut, url+"/test/yeet", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+id.String())
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bytes, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	text = string(bytes)
	assert.Contains(t, text, "1")
	assert.Contains(t, text, "/test/yeet")

	// And again
	req, err = http.NewRequest(http.MethodPut, url+"/test/yeet", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+id.String())
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bytes, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	text = string(bytes)
	assert.Contains(t, text, "2")
	assert.Contains(t, text, "/test/yeet")

	// And delete the counter without perm
	req, err = http.NewRequest(http.MethodPut, url+"/test/yeet", nil)
	assert.NoError(t, err)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// And delete the counter with perm
	req, err = http.NewRequest(http.MethodDelete, url+"/test/yeet", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+id.String())
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// should be gone now
	req, err = http.NewRequest(http.MethodDelete, url+"/test/yeet", nil)
	assert.NoError(t, err)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
