package holochain

import (
	"bytes"
	"fmt"
	. "github.com/HC-Interns/holochain-proto/hash"
	ic "github.com/libp2p/go-libp2p-crypto"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var Crash bool

func Panix(on string) {
	if Crash {
		panic(on)
	}
}

func MakeTestDirName() string {
	t := time.Now()
	d := "holochain_test" + strconv.FormatInt(t.Unix(), 10) + "." + strconv.Itoa(t.Nanosecond())
	return d
}

func MakeTestSeed(id string) io.Reader {
	return strings.NewReader(id + "1234567890123456789012345678901234567890")
}

func setupTestService() (d string, s *Service) {
	d = SetupTestDir()
	identity := "Herbert <h@bert.com>"
	s, err := Init(filepath.Join(d, DefaultDirectoryName), AgentIdentity(identity), MakeTestSeed(identity))
	if err != nil {
		panic(err)
	}
	s.Settings.DefaultBootstrapServer = "localhost:3142"
	return
}

func SetupTestService() (d string, s *Service) {
	return setupTestService()
}

// Ask the kernel for a free open port that is ready to use
func getFreePort() (port int, err error) {
	port = -1
	err = nil

	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	defer l.Close()
	port = l.Addr().(*net.TCPAddr).Port
	return
}

func setupTestChain(name string, count int, s *Service) (h *Holochain) {
	path := filepath.Join(s.Path, name)

	a := s.DefaultAgent
	if count > 0 {
		var err error
		identity := string(a.Identity()) + fmt.Sprintf("%d", count)
		a, err = NewAgent(LibP2P, AgentIdentity(identity), MakeTestSeed(identity))
		if err != nil {
			panic(err)
		}
	}

	h, err := s.MakeTestingApp(path, "toml", InitializeDB, CloneWithSameUUID, a)
	if err != nil {
		panic(err)
	}
	h.Config.DHTPort, err = getFreePort()
	if err != nil {
		panic(err)
	}
	return
}

func chainTestSetup() (hs HashSpec, key ic.PrivKey, now time.Time) {
	a, _ := NewAgent(LibP2P, "agent id", MakeTestSeed(""))
	key = a.PrivKey()
	hc := Holochain{agent: a}
	dna := DNA{DHTConfig: DHTConfig{HashType: "sha2-256"}}
	hc.nucleus = NewNucleus(&hc, &dna)
	hP := &hc
	hP.PrepareHashType()
	hs = hP.hashSpec
	return
}

func SetupTestChain(n string) (d string, s *Service, h *Holochain) {
	d, s = setupTestService()
	h = setupTestChain(n, 0, s)
	return
}

func prepareTestChain(h *Holochain) {
	_, err := h.GenChain()
	if err != nil {
		panic(err)
	}

	err = h.Activate()
	if err != nil {
		panic(err)
	}
}

func PrepareTestChain(n string) (d string, s *Service, h *Holochain) {
	d, s, h = SetupTestChain(n)
	prepareTestChain(h)
	return
}

func SetupTestDir() string {
	n := MakeTestDirName()
	d, err := ioutil.TempDir("", n)
	if err != nil {
		panic(err)
	}
	return d
}

func CleanupTestDir(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		panic(err)
	}
}

func CleanupTestChain(h *Holochain, d string) {
	h.Close()
	CleanupTestDir(d)
}

func ShouldLog(log *Logger, fn func(), messages ...string) {
	var buf bytes.Buffer
	w := log.w
	log.w = &buf
	e := log.Enabled
	log.Enabled = true
	fn()
	for _, message := range messages {
		matched := strings.Index(buf.String(), message) >= 0
		if matched {
			So(matched, ShouldBeTrue)
		} else {
			So(buf.String(), ShouldEqual, message)
		}
	}
	log.Enabled = e
	log.w = w
}

func compareFile(path1 string, path2 string, fileName string) bool {
	src, err := ReadFile(path1, fileName)
	if err != nil {
		panic(err)
	}
	dst, _ := ReadFile(path2, fileName)
	if err != nil {
		panic(err)
	}
	return (string(src) == string(dst)) && (string(src) != "")
}

func SetIdentity(h *Holochain, identity AgentIdentity) {
	agent := h.Agent()
	SetAgentIdentity(agent, identity)
	h.nodeID, h.nodeIDStr, _ = agent.NodeID()
}

func SetAgentIdentity(agent Agent, identity AgentIdentity) {
	agent.SetIdentity(identity)
	var seed io.Reader
	seed = MakeTestSeed(string(identity))
	agent.GenKeys(seed)
}

func NormaliseJSON(json string) string {
	json = strings.Replace(json, "\n", "", -1)
	json = strings.Replace(json, "    ", "", -1)
	return strings.Replace(json, ": ", ":", -1)
}
