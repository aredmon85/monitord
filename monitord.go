package main

import (
        "fmt"
	"log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"reflect"
	//"os"
        "github.com/aristanetworks/goeapi"
        //"github.com/aristanetworks/goeapi/module"
)
type Device struct {
	Host string `yaml:"host"`
	Transport string `yaml:"transport"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
type CommandSet struct {
	Command string `yaml:"command"`
	Interval int `yaml:"interval"`
}
type MonitordConfig struct {
	Devices []Device `yaml:"devices"`
	Commands []CommandSet `yaml:"commands"`
}
func CreateConnection(device Device) (*goeapi.Node, error) {
	v := reflect.ValueOf(device)
	node, err := goeapi.Connect(
		v.FieldByName("Transport").String(),
		v.FieldByName("Host").String(),
		v.FieldByName("Username").String(),
		v.FieldByName("Password").String(),
		int(v.FieldByName("Port").Int()))
	return node, err
}
func main() {
	fd, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var config MonitordConfig
	if err := yaml.Unmarshal(fd, &config); err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%+v\n",config)

	//var hosts []string
	cfg := reflect.ValueOf(config)
	devices := cfg.FieldByName("Devices").Interface().([]Device)
	for _,v := range devices {
		dev := reflect.ValueOf(v)
		host := dev.FieldByName("Host")
		transport := dev.FieldByName("Transport")
		fmt.Printf("Host: %s\n",host)
		fmt.Printf("Transport: %s\n",transport)
		//fmt.Printf("Key: %s Value: %s\n",k,v)
	}
	/*
	hosts = goeapi.Connections()
	for i := 0; i<len(hosts); i++ {
		fmt.Println(hosts[i])
	}
	node, err := goeapi.ConnectTo("R1")
        if err != nil {
                panic(err)
        }

	// get the running config and print it
	commands := []string{"show version"}
	conf, err := node.Enable(commands)
	if err != nil {
		panic(err)
	}
	for k, v := range conf[0] {
		fmt.Println("k:", k, "v:", v)
	}
	*/
}
