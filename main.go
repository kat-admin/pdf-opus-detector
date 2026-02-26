package main

// #################################################################################################
// ##### Imports ###################################################################################
// #################################################################################################

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

// #################################################################################################
// ##### Variables #################################################################################
// #################################################################################################

const (
	APPNAME = "pdf-opus-detector"
	APPVERS = "26.02.26"
)

type config struct {
	debug_mode bool
	ready_path string
	opus_path  string
	paid_path  string
}

var (
	CONFIG config
)

var Rst = "\033[0m"
var Hlc = "\033[32m"
var Red = "\033[31m"
var Grn = "\033[32m"
var Yel = "\033[33m"
var Blu = "\033[34m"
var Pur = "\033[35m"
var Cya = "\033[36m"
var Gra = "\033[37m"
var Whi = "\033[97m"

// #################################################################################################
// ##### Functions #################################################################################
// #################################################################################################

// #################################################################################################
// ##### Main Prog #################################################################################
// #################################################################################################

func main() {
	fmt.Println()
	fmt.Println(APPNAME + " - Version " + APPVERS + " - (c) Service 2 Solution GmbH")
	fmt.Println()

	inidata, err := ini.Load("pdf-opus-detector.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	section := inidata.Section("config")

	if section.Key("debug").String() == "true" {
		CONFIG.debug_mode = true
	} else {
		CONFIG.debug_mode = false
	}
	CONFIG.ready_path = section.Key("ready_path").String()
	CONFIG.opus_path = section.Key("opus_path").String()
	CONFIG.paid_path = section.Key("paid_path").String()

	fmt.Println("Fertig")
}
