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
	"time"

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
	debug_mode  bool
	ignore_days int
	ready_path  string
	opus_list   string
	paid_path   string
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

var csvData [][]string

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

func moveFile(src, dst string) error {
	// fmt.Printf("  -> Verschiebe Datei von %s nach %s\n", src, dst)

	// Sicherstellen, dass Zielordner existiert
	err := os.MkdirAll(filepath.Dir(dst), os.ModePerm)
	if err != nil {
		log.Fatalf("  Fehler beim Erstellen des Zielordners: %v", err)
	}

	// Datei verschieben
	err = os.Rename(src, dst)
	if err != nil {
		log.Fatalf("  Fehler beim Verschieben der Datei: %v", err)
	}

	return nil
}

func searchInvoiceInOpusList(invoiceNumber string) bool {
	for _, record := range csvData {
		// Konto/Kundennummer   = record[1]
		// Rechnungsnummer      = record[3]
		// Soll                 = record[6]
		// Haben                = record[7]
		// Saldo                = record[8]

		if record[3] == invoiceNumber {
			// fmt.Printf(" :: gefunden in OPUS-Liste: %v   == %s / %s / %s\n", record[3], record[6], record[7], record[8])
			return true
		}
	}

	return false
}

func searchInvoices() error {
	fmt.Printf("Lese Datei: %s\n", CONFIG.opus_list)
	csvData, _ = readCsvFile(CONFIG.opus_list)
	fmt.Printf("Gelesene Daten aus %s (ohne Kopfzeile): %d\n", CONFIG.opus_list, len(csvData))

	fmt.Printf("Suche nach Rechnungen in %s\n", CONFIG.ready_path)

	threshold := time.Now().AddDate(0, 0, (CONFIG.ignore_days * -1)) // CONFIG.ignore_days Tage in der Vergangenheit

	err := filepath.WalkDir(CONFIG.ready_path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ".pdf" {
			info, err := d.Info()
			if err != nil {
				log.Printf("Fehler beim Abrufen der Dateiinfo für %s: %v", path, err)
				return nil // Weiter mit der nächsten Datei
			}

			if info.ModTime().Before(threshold) {
				filename := d.Name()
				if len(filename) >= 8 { // Rechnungsnummer ist 8 Zeichen lang, + .pdf sind 4, gesamt 12.
					invoiceNumber := filename[len(filename)-12 : len(filename)-4]
					fmt.Printf(" - %-31s", d.Name())

					if !searchInvoiceInOpusList(invoiceNumber) {
						merr := moveFile(filepath.Join(CONFIG.ready_path, filename), filepath.Join(CONFIG.paid_path, filename))

						if merr == nil {
							fmt.Printf(" %snach %s verschoben%s\n", Grn, CONFIG.ready_path, Rst)
						} else {
							fmt.Printf(" %s(Fehler beim verschieben: %w)%s\n", Red, merr, Rst)
						}
					} else {
						fmt.Printf(" %s(ignoriert - in OPUS-Liste gefunden)%s\n", Gra, Rst)
					}
				} else {
					fmt.Printf(" - %-31s %s(ignoriert - Name zu kurz)%s\n", d.Name(), Gra, Rst)
				}
			} else {
				fmt.Printf(" - %-31s %s(ignoriert - jünger als %d Tage)%s\n", d.Name(), Gra, CONFIG.ignore_days, Rst)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf(" - Fehler beim Durchsuchen des Verzeichnisses %s: %w", CONFIG.ready_path, err)
	}

	return nil
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
	CONFIG.ignore_days = section.Key("ignore_days").MustInt()
	CONFIG.opus_list = section.Key("opus_list").String()
	CONFIG.ready_path = section.Key("ready_path").String()
	CONFIG.paid_path = section.Key("paid_path").String()

	err = searchInvoices()
	if err != nil {
		log.Fatalf("Fehler beim Suchen der Rechnungen: %v", err)
	}

	fmt.Println("Fertig")
}
