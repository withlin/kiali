package config

import (
	"os"
	"testing"
)

func TestEnvVar(t *testing.T) {
	defer os.Setenv(EnvServerAddress, os.Getenv(EnvServerAddress))
	defer os.Setenv(EnvServerPort, os.Getenv(EnvServerPort))
	os.Setenv(EnvServerAddress, "test-address")
	os.Setenv(EnvServerPort, "12345")

	conf := NewConfig()

	if conf.Server.Address != "test-address" {
		t.Error("server address is wrong")
	}
	if conf.Server.Port != 12345 {
		t.Error("server port is wrong")
	}
}

func TestDefaults(t *testing.T) {
	conf := NewConfig()

	if conf.Server.Address != "" {
		t.Error("server address default is wrong")
	}

	if conf.Server.Port != 20000 {
		t.Error("server port default is wrong")
	}
}

func TestMarshalUnmarshalStaticContentRootDirectory(t *testing.T) {
	testConf := Config{
		Server: Server{
			StaticContentRootDirectory: "/tmp",
		},
	}

	yamlString, err := Marshal(&testConf)
	if err != nil {
		t.Errorf("Failed to marshal: %v", err)
	}
	if yamlString != "server:\n  static_content_root_directory: /tmp\n" {
		t.Errorf("Failed to marshal - StaticContentRootDirectory to static_content_root_directory: [%v]", yamlString)
	}
	conf, err := Unmarshal(yamlString)
	if err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}
	if conf.Server.StaticContentRootDirectory != "/tmp" {
		t.Errorf("Failed to unmarshal static content root directory:\n%v", conf)
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	testConf := Config{
		Server: Server{
			Address: "foo-test",
			Port:    321,
			StaticContentRootDirectory: "/tmp",
		},
	}

	yamlString, err := Marshal(&testConf)
	if err != nil {
		t.Errorf("Failed to marshal: %v", err)
	}
	if yamlString == "" {
		t.Errorf("Failed to marshal - empty yaml string")
	}

	conf, err := Unmarshal(yamlString)
	if err != nil {
		t.Errorf("Failed to unmarshal: %v", err)
	}

	if conf.Server.Address != "foo-test" {
		t.Errorf("Failed to unmarshal server address:\n%v", conf)
	}
	if conf.Server.Port != 321 {
		t.Errorf("Failed to unmarshal server port:\n%v", conf)
	}
	if conf.Server.StaticContentRootDirectory != "/tmp" {
		t.Errorf("Failed to unmarshal static content root directory:\n%v", conf)
	}
}

func TestLoadSave(t *testing.T) {
	testConf := Config{
		Server: Server{
			Address: "foo-test",
			Port:    321,
		},
	}

	filename := "/tmp/config_test.yaml"
	defer os.Remove(filename)

	err := SaveToFile(filename, &testConf)
	if err != nil {
		t.Errorf("Failed to save to file: %v", err)
	}

	conf, err := LoadFromFile(filename)
	if err != nil {
		t.Errorf("Failed to load from file: %v", err)
	}

	t.Logf("Config from file\n%v", conf)

	if conf.Server.Address != "foo-test" {
		t.Errorf("Failed to unmarshal server address:\n%v", conf)
	}
	if conf.Server.Port != 321 {
		t.Errorf("Failed to unmarshal server port:\n%v", conf)
	}
}

func TestError(t *testing.T) {
	_, err := Unmarshal("bogus-yaml")
	if err == nil {
		t.Errorf("Unmarshal should have failed")
	}

	_, err = LoadFromFile("bogus-file-name")
	if err == nil {
		t.Errorf("Load should have failed")
	}
}
