package main

import (
        "fmt"
	"log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"reflect"
        "github.com/aristanetworks/goeapi"
	"time"
)
type Device struct {
	Host string `yaml:"host"`
	Transport string `yaml:"transport"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	EnablePassword string `yaml:"enable_password"`
	Port int `yaml:"port"`
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
func monitor(node *goeapi.Node, cmds map[string]int) {
	fmt.Printf("Monitoring node %s\n",node)
	var last_ran map[string]time.Time
	last_ran = make(map[string]time.Time)
	for k,_ := range cmds {
		last_ran[k] = time.Now()
	}
	for {
		for k,_ := range last_ran {
			if time.Now().After(last_ran[k].Add(time.Second * time.Duration(cmds[k]))) {
				command := []string{k}
				_, err := node.Enable(command)
				if err != nil {
					log.Fatal(err)
				}
				last_ran[k] = time.Now()
			}
		}
		time.Sleep(1 * time.Second)
	}
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
	cfg := reflect.ValueOf(config)
	devices := cfg.FieldByName("Devices").Interface().([]Device)
	var commands []CommandSet
	var nodes []*goeapi.Node
	for _,v := range devices {
		dev := reflect.ValueOf(v)
		host := dev.FieldByName("Host").String()
		transport := dev.FieldByName("Transport").String()
		port := int(dev.FieldByName("Port").Int())
		username := dev.FieldByName("Username").String()
		password := dev.FieldByName("Password").String()
		enable_password := dev.FieldByName("EnablePassword").String()
		node, err := goeapi.Connect(transport,host,username,password,port)
		node.EnableAuthentication(enable_password)
		if err != nil {
			log.Fatal(err)
		}
		nodes = append(nodes,node)
	}
	commands = cfg.FieldByName("Commands").Interface().([]CommandSet)
	var cmds map[string]int
	cmds = make(map[string]int)
	for _,v := range commands {
		cmd := reflect.ValueOf(v)
		cmds[cmd.FieldByName("Command").String()] = int(cmd.FieldByName("Interval").Int())
	}
	for _,v := range nodes {
		go monitor(v,cmds)
	}
	select{}
}
