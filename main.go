package main

// #################################################################################################
// ##### Imports ###################################################################################
// #################################################################################################

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

func readCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Öffnen der Datei %s: %w", filePath, err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = ';'          // Setzt das Trennzeichen auf Semikolon
	reader.FieldsPerRecord = -1 // Erlaubt variable Anzahl von Feldern pro Record

	// Erste Zeile (Header) überspringen
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Lesen der Kopfzeile von %s: %w", filePath, err)
	}

	data, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Fehler beim Lesen der CSV-Daten aus %s: %w", filePath, err)
	}

	return data, nil
}

func searchOpusFiles(opus_path string) ([]string, error) {
	var csvFiles []string

	fmt.Printf("Suche nach OPUS-Listen in %s\n", opus_path)

	err := filepath.WalkDir(opus_path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ".csv" {
			csvFiles = append(csvFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Fehler beim Durchsuchen des Verzeichnisses %s: %w", opus_path, err)
	}

	return csvFiles, nil
}

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

	csvFiles, err := searchOpusFiles(CONFIG.opus_path)
	if err != nil {
		log.Fatalf("Fehler beim Suchen der OPUS-Dateien: %v", err)
	}

	if len(csvFiles) > 0 {
		fmt.Println("\nVerarbeite OPUS-CSV-Dateien:")
		for _, file := range csvFiles {
			fmt.Printf("Lese Datei: %s\n", file)
			csvData, err := readCsvFile(file)
			if err != nil {
				log.Printf("Fehler beim Lesen der Datei %s: %v", file, err)
				continue
			}

			fmt.Printf("Gelesene Daten aus %s (ohne Kopfzeile):\n", file)
			for i, record := range csvData {
				fmt.Printf("  Record %d: %v   == %s / %s / %s\n", i+1, record[3], record[6], record[7], record[8])

				// Konto/Kundennummer   = record[1]
				// Rechnungsnummer      = record[3]
				// Soll                 = record[6]
				// Haben                = record[7]
				// Saldo                = record[8]
			}
		}
	}

	fmt.Println("Fertig")
}
