package holochain

import (
	"bytes"
	"fmt"
	ic "github.com/libp2p/go-libp2p-crypto"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	d := SetupTestDir()
	defer CleanupTestDir(d)

	Convey("we can detect an uninitialized directory", t, func() {
		So(IsInitialized(filepath.Join(d, DefaultDirectoryName)), ShouldBeFalse)
	})

	agent := "Fred Flintstone <fred@flintstone.com>"

	s, err := Init(filepath.Join(d, DefaultDirectoryName), AgentIdentity(agent))

	Convey("when initializing service in a directory", t, func() {
		So(err, ShouldBeNil)

		Convey("it should return a service with default values", func() {
			So(s.DefaultAgent.Identity(), ShouldEqual, AgentIdentity(agent))
			So(fmt.Sprintf("%v", s.Settings), ShouldEqual, "{true true bootstrap.holochain.net:10000 false}")
		})

		p := filepath.Join(d, DefaultDirectoryName)
		Convey("it should create agent files", func() {
			a, err := LoadAgent(p)
			So(err, ShouldBeNil)
			So(a.Identity(), ShouldEqual, AgentIdentity(agent))
		})

		Convey("we can detect that it was initialized", func() {
			So(IsInitialized(filepath.Join(d, DefaultDirectoryName)), ShouldBeTrue)
		})

		Convey("it should create an agent file", func() {
			a, err := ReadFile(p, AgentFileName)
			So(err, ShouldBeNil)
			So(string(a), ShouldEqual, agent)
		})
	})
}

func TestLoadService(t *testing.T) {
	d, service := setupTestService()
	root := service.Path
	defer CleanupTestDir(d)
	Convey("loading service from disk should set up the struct", t, func() {
		s, err := LoadService(root)
		So(err, ShouldBeNil)
		So(s.Path, ShouldEqual, root)
		So(s.Settings.DefaultPeerModeDHTNode, ShouldEqual, true)
		So(s.Settings.DefaultPeerModeAuthor, ShouldEqual, true)
		So(s.DefaultAgent.Identity(), ShouldEqual, AgentIdentity("Herbert <h@bert.com>"))
	})

}

func TestValidateServiceConfig(t *testing.T) {
	svc := ServiceConfig{}

	Convey("it should fail without one peer mode set to true", t, func() {
		err := svc.Validate()
		So(err.Error(), ShouldEqual, SysFileName+": At least one peer mode must be set to true.")
	})

	svc.DefaultPeerModeAuthor = true

	Convey("it should validate", t, func() {
		err := svc.Validate()
		So(err, ShouldBeNil)
	})

}

func TestConfiguredChains(t *testing.T) {
	d, s, h := SetupTestChain("test")
	defer CleanupTestDir(d)

	Convey("Configured chains should return a hash of all the chains in the Service", t, func() {
		chains, err := s.ConfiguredChains()
		So(err, ShouldBeNil)
		So(chains["test"].nucleus.dna.UUID, ShouldEqual, h.nucleus.dna.UUID)
	})
}

func TestServiceGenChain(t *testing.T) {
	d, s, h := SetupTestChain("test")
	defer CleanupTestDir(d)

	Convey("it should return a list of the chains", t, func() {
		list := s.ListChains()
		So(list, ShouldEqual, "installed holochains:\n    test <not-started>\n")
	})
	Convey("it should start a chain and return a holochain object", t, func() {
		h2, err := s.GenChain("test")
		So(err, ShouldBeNil)
		So(h2.nucleus.dna.UUID, ShouldEqual, h.nucleus.dna.UUID)
		list := s.ListChains()
		So(list, ShouldEqual, fmt.Sprintf("installed holochains:\n    test %v\n", h2.dnaHash))
	})
}

func TestCloneNew(t *testing.T) {
	d, s, h0 := SetupTestChain("test")
	defer CleanupTestDir(d)

	name := "test2"
	root := filepath.Join(s.Path, name)

	orig := filepath.Join(s.Path, "test")

	agent, err := LoadAgent(s.Path)
	if err != nil {
		panic(err)
	}

	Convey("it should clone a chain by copying and creating an new UUID", t, func() {
		err = s.Clone(orig, root, agent, CloneWithNewUUID, InitializeDB)
		So(err, ShouldBeNil)

		So(DirExists(root, ChainDataDir), ShouldBeTrue)
		So(FileExists(root, ChainDataDir, StoreFileName), ShouldBeTrue)

		h, err := s.Load(name) // reload to confirm that it got saved correctly
		So(err, ShouldBeNil)

		So(h.Name(), ShouldEqual, "test2")
		So(h.nucleus.dna.UUID, ShouldNotEqual, h0.nucleus.dna.UUID)

		agent, err := LoadAgent(s.Path)
		So(err, ShouldBeNil)
		So(h.agent.Identity(), ShouldEqual, agent.Identity())
		So(ic.KeyEqual(h.agent.PrivKey(), agent.PrivKey()), ShouldBeTrue)
		So(ic.KeyEqual(h.agent.PubKey(), agent.PubKey()), ShouldBeTrue)

		So(compareFile(filepath.Join(orig, "dna", "zySampleZome"), filepath.Join(h.DNAPath(), "zySampleZome"), "zySampleZome.zy"), ShouldBeTrue)

		So(h.rootPath, ShouldEqual, root)
		So(h.UIPath(), ShouldEqual, root+"/ui")
		So(h.DNAPath(), ShouldEqual, root+"/dna")
		So(h.DBPath(), ShouldEqual, root+"/db")

		So(compareFile(filepath.Join(orig, "ui"), h.UIPath(), "index.html"), ShouldBeTrue)
		So(compareFile(filepath.Join(orig, "dna", "zySampleZome"), filepath.Join(h.DNAPath(), "zySampleZome"), "profile.json"), ShouldBeTrue)
		So(compareFile(filepath.Join(orig, "dna"), h.DNAPath(), "properties_schema.json"), ShouldBeTrue)
		So(compareFile(orig, h.rootPath, ConfigFileName+".toml"), ShouldBeTrue)

		So(compareFile(filepath.Join(orig, ChainTestDir), filepath.Join(h.rootPath, ChainTestDir), "testSet1.json"), ShouldBeTrue)

		So(h.nucleus.dna.Progenitor.Identity, ShouldEqual, "Herbert <h@bert.com>")
		pk, _ := agent.PubKey().Bytes()
		So(string(h.nucleus.dna.Progenitor.PubKey), ShouldEqual, string(pk))
	})
}

func TestCloneJoin(t *testing.T) {
	d, s, h0 := SetupTestChain("test")
	defer CleanupTestDir(d)

	name := "test2"
	root := filepath.Join(s.Path, name)

	orig := filepath.Join(s.Path, "test")

	agent, err := LoadAgent(s.Path)
	if err != nil {
		panic(err)
	}

	Convey("it should clone a chain by copying and without creating a new UUID", t, func() {
		err = s.Clone(orig, root, agent, CloneWithSameUUID, InitializeDB)
		So(err, ShouldBeNil)

		So(DirExists(root, ChainDataDir), ShouldBeTrue)
		So(FileExists(root, ChainDataDir, StoreFileName), ShouldBeTrue)

		h, err := s.Load(name) // reload to confirm that it got saved correctly
		So(err, ShouldBeNil)

		So(h.Name(), ShouldEqual, "test")
		So(h.nucleus.dna.UUID, ShouldEqual, h0.nucleus.dna.UUID)
		agent, err := LoadAgent(s.Path)
		So(err, ShouldBeNil)
		So(h.agent.Identity(), ShouldEqual, agent.Identity())
		So(ic.KeyEqual(h.agent.PrivKey(), agent.PrivKey()), ShouldBeTrue)

		So(ic.KeyEqual(h.agent.PubKey(), agent.PubKey()), ShouldBeTrue)
		src, _ := ReadFile(orig, "dna", "zySampleZome.zy")
		dst, _ := ReadFile(root, "zySampleZome.zy")
		So(string(src), ShouldEqual, string(dst))
		So(FileExists(h.UIPath(), "index.html"), ShouldBeTrue)
		So(FileExists(h.DNAPath(), "zySampleZome", "profile.json"), ShouldBeTrue)
		So(FileExists(h.DNAPath(), "properties_schema.json"), ShouldBeTrue)
		So(FileExists(h.rootPath, ConfigFileName+".toml"), ShouldBeTrue)

		So(h.nucleus.dna.Progenitor.Identity, ShouldEqual, "Progenitor Agent <progenitore@example.com>")
		pk := []byte{8, 1, 18, 32, 193, 43, 31, 148, 23, 249, 163, 154, 128, 25, 237, 167, 253, 63, 214, 220, 206, 131, 217, 74, 168, 30, 215, 237, 231, 160, 69, 89, 48, 17, 104, 210}
		So(string(h.nucleus.dna.Progenitor.PubKey), ShouldEqual, string(pk))

	})
}

func TestCloneNoDB(t *testing.T) {
	d, s, _ := SetupTestChain("test")
	defer CleanupTestDir(d)

	name := "test2"
	root := filepath.Join(s.Path, name)

	orig := filepath.Join(s.Path, "test")

	agent, err := LoadAgent(s.Path)
	if err != nil {
		panic(err)
	}

	Convey("it should create a chain from the examples directory", t, func() {
		err = s.Clone(orig, root, agent, CloneWithNewUUID, SkipInitializeDB)
		So(err, ShouldBeNil)

		So(DirExists(root, ChainDataDir), ShouldBeFalse)
		So(FileExists(root, ChainDNADir, "zySampleZome", "profile.json"), ShouldBeTrue)
	})
}

func TestGenDev(t *testing.T) {
	d, s := setupTestService()
	defer CleanupTestDir(d)
	name := "test"
	root := filepath.Join(s.Path, name)

	Convey("we detected unconfigured holochains", t, func() {
		f, err := s.IsConfigured(name)
		So(f, ShouldEqual, "")
		So(err.Error(), ShouldEqual, "No DNA file in "+filepath.Join(root, ChainDNADir)+"/")
		_, err = s.load("test", "json")
		So(err.Error(), ShouldEqual, "open "+filepath.Join(root, ChainDNADir, DNAFileName+".json")+": no such file or directory")

	})

	Convey("when generating a dev holochain", t, func() {
		h, err := s.GenDev(root, "json", InitializeDB)
		So(err, ShouldBeNil)

		f, err := s.IsConfigured(name)
		So(err, ShouldBeNil)
		So(f, ShouldEqual, "json")

		h, err = s.Load(name)
		So(err, ShouldBeNil)

		lh, err := s.load(name, "json")
		So(err, ShouldBeNil)
		So(lh.nodeID, ShouldEqual, h.nodeID)
		So(lh.nodeIDStr, ShouldEqual, h.nodeIDStr)

		So(lh.Config.Port, ShouldEqual, DefaultPort)
		So(h.Config.PeerModeDHTNode, ShouldEqual, s.Settings.DefaultPeerModeDHTNode)
		So(h.Config.PeerModeAuthor, ShouldEqual, s.Settings.DefaultPeerModeAuthor)
		So(h.Config.BootstrapServer, ShouldEqual, s.Settings.DefaultBootstrapServer)
		So(h.Config.EnableMDNS, ShouldEqual, s.Settings.DefaultEnableMDNS)

		So(DirExists(root), ShouldBeTrue)
		So(DirExists(h.DNAPath()), ShouldBeTrue)
		So(DirExists(h.TestPath()), ShouldBeTrue)
		So(DirExists(h.UIPath()), ShouldBeTrue)
		So(FileExists(h.TestPath(), "sampleScenario", "listener.json"), ShouldBeTrue)
		So(FileExists(h.DNAPath(), "zySampleZome", "profile.json"), ShouldBeTrue)
		So(FileExists(h.UIPath(), "index.html"), ShouldBeTrue)
		So(FileExists(h.UIPath(), "hc.js"), ShouldBeTrue)
		So(FileExists(h.rootPath, ConfigFileName+".json"), ShouldBeTrue)

		Convey("we should not be able re generate it", func() {
			_, err = s.GenDev(root, "json", SkipInitializeDB)
			So(err.Error(), ShouldEqual, "holochain: "+root+" already exists")
		})
	})
}

func TestSaveFromScaffold(t *testing.T) {
	d, s := setupTestService()
	defer CleanupTestDir(d)
	name := "test"
	root := filepath.Join(s.Path, name)

	Convey("it should write out a scaffold file to a directory tree with JSON encoding", t, func() {
		scaffoldReader := bytes.NewBuffer([]byte(BasicTemplateScaffold))

		scaffold, err := s.SaveFromScaffold(scaffoldReader, root, "appName", "json", false)
		So(err, ShouldBeNil)
		So(scaffold, ShouldNotBeNil)
		So(scaffold.ScaffoldVersion, ShouldEqual, ScaffoldVersion)
		So(scaffold.DNA.Name, ShouldEqual, "appName")
		So(DirExists(root), ShouldBeTrue)
		So(DirExists(root, ChainDNADir), ShouldBeTrue)
		So(DirExists(root, ChainUIDir), ShouldBeTrue)
		So(DirExists(root, ChainTestDir), ShouldBeTrue)
		So(DirExists(root, ChainTestDir, scaffold.Scenarios[0].Name), ShouldBeTrue)
		So(FileExists(root, ChainTestDir, scaffold.Scenarios[0].Name, scaffold.Scenarios[0].Roles[0].Name+".json"), ShouldBeTrue)
		So(FileExists(root, ChainTestDir, scaffold.Scenarios[0].Name, scaffold.Scenarios[0].Roles[1].Name+".json"), ShouldBeTrue)
		So(FileExists(root, ChainTestDir, scaffold.Scenarios[0].Name, "_config.json"), ShouldBeTrue)

		So(DirExists(root, ChainDNADir, "sampleZome"), ShouldBeTrue)
		So(FileExists(root, ChainDNADir, "sampleZome", "sampleEntry.json"), ShouldBeTrue)
		So(FileExists(root, ChainDNADir, "sampleZome", "sampleZome.js"), ShouldBeTrue)
		So(FileExists(root, ChainDNADir, DNAFileName+".json"), ShouldBeTrue)
		So(FileExists(root, ChainDNADir, "properties_schema.json"), ShouldBeTrue)
		So(FileExists(root, ChainTestDir, "sample.json"), ShouldBeTrue)
		So(FileExists(root, ChainUIDir, "index.html"), ShouldBeTrue)
		So(FileExists(root, ChainUIDir, "hc.js"), ShouldBeTrue)
	})

	Convey("it should write out a scaffold file to a directory tree with toml encoding", t, func() {
		scaffoldReader := bytes.NewBuffer([]byte(BasicTemplateScaffold))

		root2 := filepath.Join(s.Path, name+"2")

		scaffold, err := s.SaveFromScaffold(scaffoldReader, root2, "appName", "toml", false)
		So(err, ShouldBeNil)
		So(scaffold, ShouldNotBeNil)
		So(scaffold.ScaffoldVersion, ShouldEqual, ScaffoldVersion)
		So(DirExists(root2), ShouldBeTrue)
		So(FileExists(root2, ChainDNADir, DNAFileName+".toml"), ShouldBeTrue)
		// the reset of the files are still saved as json...
	})

	Convey("it should write out a scaffold file to a directory tree with binary UI files", t, func() {
		scaffoldReader := bytes.NewBuffer([]byte(TestingAppScaffold()))

		_, err := s.SaveFromScaffold(scaffoldReader, root+"3", "appName2", "json", false)
		root3 := filepath.Join(s.Path, name+"3")

		So(err, ShouldBeNil)
		So(DirExists(root3, ChainUIDir), ShouldBeTrue)
		So(FileExists(root3, ChainUIDir, "index.html"), ShouldBeTrue)
		So(FileExists(root3, ChainUIDir, "logo.png"), ShouldBeTrue)
	})

}

func TestMakeConfig(t *testing.T) {
	d, s := setupTestService()
	defer CleanupTestDir(d)
	h := &Holochain{encodingFormat: "json", rootPath: d}
	Convey("make config should produce default values", t, func() {
		err := makeConfig(h, s)
		So(err, ShouldBeNil)
		So(h.Config.Port, ShouldEqual, DefaultPort)
		So(h.Config.EnableMDNS, ShouldBeFalse)
		So(h.Config.BootstrapServer, ShouldNotEqual, "")
		So(h.Config.Loggers.App.Format, ShouldEqual, "%{color:cyan}%{message}")

	})

	Convey("make config should produce default config from OS env overridden values", t, func() {
		os.Setenv("HOLOCHAINCONFIG_PORT", "12345")
		os.Setenv("HOLOCHAINCONFIG_ENABLEMDNS", "true")
		os.Setenv("HCLOG_PREFIX", "prefix:%{color:cyan}")
		os.Setenv("HOLOCHAINCONFIG_BOOTSTRAP", "_")
		err := makeConfig(h, s)
		So(err, ShouldBeNil)

		So(h.Config.Port, ShouldEqual, 12345)
		So(h.Config.EnableMDNS, ShouldBeTrue)
		So(h.Config.Loggers.App.Format, ShouldEqual, "%{color:cyan}%{message}")
		So(h.Config.Loggers.App.Prefix, ShouldEqual, "prefix:")
		So(h.Config.Loggers.App.PrefixColor, ShouldEqual, h.Config.Loggers.App.GetColor("cyan"))
		So(h.Config.BootstrapServer, ShouldEqual, "")
	})
}

func TestMakeScaffold(t *testing.T) {
	d, s := setupTestService()
	defer CleanupTestDir(d)
	name := "test"
	root := filepath.Join(s.Path, name)
	h, err := s.GenDev(root, "json", InitializeDB)
	if err != nil {
		panic(err)
	}
	Convey("make scaffold should produce a scaffold file for holochain", t, func() {
		scaffoldBlob, err := s.MakeScaffold(h)
		So(err, ShouldBeNil)
		scaffoldReader := bytes.NewBuffer(scaffoldBlob)
		if err != nil {
			panic(err)
		}
		root = filepath.Join(s.Path, "appFromScaffold")
		scaffold, err := s.SaveFromScaffold(scaffoldReader, root, "appFromScaffold", "json", false)
		So(err, ShouldBeNil)
		So(scaffold, ShouldNotBeNil)
		So(scaffold.ScaffoldVersion, ShouldEqual, ScaffoldVersion)
		So(scaffold.DNA.Name, ShouldEqual, "appFromScaffold")

		So(DirExists(root), ShouldBeTrue)
		So(DirExists(root, ChainDNADir), ShouldBeTrue)
		So(DirExists(root, ChainUIDir), ShouldBeTrue)
		So(DirExists(root, ChainTestDir), ShouldBeTrue)
		So(DirExists(root, ChainTestDir, scaffold.Scenarios[0].Name), ShouldBeTrue)
		So(FileExists(root, ChainTestDir, scaffold.Scenarios[0].Name, scaffold.Scenarios[0].Roles[0].Name+".json"), ShouldBeTrue)
		So(FileExists(root, ChainTestDir, scaffold.Scenarios[0].Name, scaffold.Scenarios[0].Roles[1].Name+".json"), ShouldBeTrue)
		So(FileExists(root, ChainTestDir, scaffold.Scenarios[0].Name, TestConfigFileName), ShouldBeTrue)

		So(DirExists(root, ChainDNADir, "jsSampleZome"), ShouldBeTrue)
		So(FileExists(root, ChainDNADir, "jsSampleZome", "profile.json"), ShouldBeTrue)
		So(FileExists(root, ChainDNADir, "jsSampleZome", "jsSampleZome.js"), ShouldBeTrue)
		So(FileExists(root, ChainDNADir, DNAFileName+".json"), ShouldBeTrue)
		So(FileExists(root, ChainDNADir, "properties_schema.json"), ShouldBeTrue)
		So(FileExists(root, ChainTestDir, "testSet1.json"), ShouldBeTrue)
		So(FileExists(root, ChainTestDir, "testSet1.json"), ShouldBeTrue)
		So(FileExists(root, ChainUIDir, "index.html"), ShouldBeTrue)
		So(FileExists(root, ChainUIDir, "hc.js"), ShouldBeTrue)
		So(FileExists(root, ChainUIDir, "logo.png"), ShouldBeTrue)
	})
}

func TestLoadTestFiles(t *testing.T) {
	d, _, h := SetupTestChain("test")
	defer CleanupTestDir(d)

	Convey("it should fail if there's no test data", t, func() {
		tests, err := LoadTestFiles(d)
		So(tests, ShouldBeNil)
		So(err.Error(), ShouldEqual, "no test files found in: "+d)
	})

	Convey("it should load test files", t, func() {
		path := h.TestPath()
		tests, err := LoadTestFiles(path)
		So(err, ShouldBeNil)
		So(len(tests), ShouldEqual, 2)
	})
}

func TestGetScenarioData(t *testing.T) {
	d, _, h := SetupTestChain("test")
	defer CleanupTestDir(d)
	Convey("it should return list of scenarios", t, func() {
		scenarios, err := GetTestScenarios(h)
		So(err, ShouldBeNil)
		_, ok := scenarios["sampleScenario"]
		So(ok, ShouldBeTrue)
		_, ok = scenarios["foo"]
		So(ok, ShouldBeFalse)
	})
	Convey("it should return list of scenarios in a role", t, func() {
		scenarios, err := GetTestScenarioRoles(h, "sampleScenario")
		So(err, ShouldBeNil)
		So(fmt.Sprintf("%v", scenarios), ShouldEqual, `[listener speaker]`)
	})
}
