package main

import (
	"archive/zip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"
)

func TestRemoveAccents(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"sem acentos", "palavra", "palavra"},
		{"com acentos", "palavrà", "palavra"},
		{"frase com acentos", "uma frase com acentuação", "uma frase com acentuacao"},
		{"string vazia", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := removeAccents(tc.input)
			if result != tc.expected {
				t.Errorf("esperado: %s, obtido: %s", tc.expected, result)
			}
		})
	}
}

func TestProcessaValorCelula(t *testing.T) {
	testCases := []struct {
		name     string
		value    interface{}
		header   string
		expected interface{}
	}{
		{"string normal", "valor", "QUALQUER", "valor"},
		{"string vazia", "", "QUALQUER", nil},
		{"hífen", "-", "QUALQUER", nil},
		{"booleano verdadeiro", "Sim", "RESTRIÇÃO HOSPITALAR", true},
		{"booleano falso", "Não", "RESTRIÇÃO HOSPITALAR", false},
		{"numérico com vírgula", "12,34", "PF 12%", 12.34},
		{"numérico com ponto", "56.78", "PMVG 17%", 56.78},
		{"numérico com asterisco", "90,12*", "PF 18% ALC", 90.12},
		{"princípio ativo nulo", nil, "SUBSTÂNCIA", ""},
		{"CAP verdadeiro", "Sim", "CAP", true},
		{"CAP falso", "Não", "CAP", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := processaValorCelula(tc.value, tc.header)
			if result != tc.expected {
				t.Errorf("esperado: %v, obtido: %v", tc.expected, result)
			}
		})
	}
}

func TestConvertIntToExcelColumn(t *testing.T) {
	testCases := []struct {
		name     string
		input    int
		expected string
	}{
		{"zero", 0, "A"},
		{"um", 1, "B"},
		{"vinte e seis", 25, "Z"},
		{"vinte e sete", 26, "AA"},
		{"cinquenta e um", 50, "AY"},
		{"setenta e cinco", 74, "BW"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convertIntToExcelColumn(tc.input)
			if result != tc.expected {
				t.Errorf("esperado: %s, obtido: %s", tc.expected, result)
			}
		})
	}
}

func TestWriteJSONFile(t *testing.T) {
	// Create a temporary directory for the test file
	tempDir := t.TempDir()
	infilePath := filepath.Join(tempDir, "test-file.xlsx")
	outputFilePath := filepath.Join(tempDir, "test-file.json")

	// Create a sample output
	output := Output{
		Metadados: Metadados{
			Data:            "2025-07-03",
			DataAtualizacao: "2025-07-04",
			Observacoes:     []string{"Observação 1", "Observação 2"},
		},
		Medicamentos: []Medicamento{
			{"SUBSTÂNCIA": "IBUPROFENO", "CNPJ": "12.345.678/0001-90"},
			{"SUBSTÂNCIA": "PARACETAMOL", "CNPJ": "98.765.432/0001-10"},
		},
		Laboratorios: map[string]string{
			"12.345.678/0001-90": "LAB A",
			"98.765.432/0001-10": "LAB B",
		},
		Apresentacoes: []string{"COM REV", "GOTAS"},
	}

	// Write the JSON file
	if err := writeJSONFile(output, infilePath); err != nil {
		t.Fatalf("writeJSONFile failed: %v", err)
	}

	// Check if the file exists
	if _, err := os.Stat(outputFilePath); os.IsNotExist(err) {
		t.Fatalf("expected JSON file %s to exist, but it does not", outputFilePath)
	}
}

func TestWriteZipFile(t *testing.T) {
	// Create a temporary directory for the test file
	tempDir := t.TempDir()
	infilePath := filepath.Join(tempDir, "test-file.xlsx")
	outputZipPath := filepath.Join(tempDir, "test-file.zip")
	jsonFileName := "test-file.json" // The name of the JSON file inside the zip

	// Create a sample output
	expectedOutput := Output{
		Metadados: Metadados{
			Data:            "2025-07-03",
			DataAtualizacao: "2025-07-04",
			Observacoes:     []string{"Observação 1", "Observação 2"},
		},
		Medicamentos: []Medicamento{
			{"SUBSTÂNCIA": "IBUPROFENO", "CNPJ": "12.345.678/0001-90"},
			{"SUBSTÂNCIA": "PARACETAMOL", "CNPJ": "98.765.432/0001-10"},
		},
		Laboratorios: map[string]string{
			"12.345.678/0001-90": "LAB A",
			"98.765.432/0001-10": "LAB B",
		},
		Apresentacoes: []string{"COM REV", "GOTAS"},
	}

	// Write the JSON file (which is now zipped)
	if err := writeZipFile(expectedOutput, infilePath); err != nil {
		t.Fatalf("writeJSONFile failed: %v", err)
	}

	// Check if the zip file exists
	if _, err := os.Stat(outputZipPath); os.IsNotExist(err) {
		t.Fatalf("expected zip file %s to exist, but it does not", outputZipPath)
	}

	// Open the zip file
	r, err := zip.OpenReader(outputZipPath)
	if err != nil {
		t.Fatalf("failed to open zip file: %v", err)
	}
	defer r.Close()

	// Find the JSON file inside the zip
	var jsonFile *zip.File
	for _, f := range r.File {
		if f.Name == jsonFileName {
			jsonFile = f
			break
		}
	}

	if jsonFile == nil {
		t.Fatalf("JSON file %s not found inside the zip archive", jsonFileName)
	}

	// Read the content of the JSON file
	rc, err := jsonFile.Open()
	if err != nil {
		t.Fatalf("failed to open JSON file in zip: %v", err)
	}
	defer rc.Close()

	jsonData, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("failed to read JSON data from zip: %v", err)
	}

	// Unmarshal the JSON data
	var actualOutput Output
	if err := json.Unmarshal(jsonData, &actualOutput); err != nil {
		t.Fatalf("failed to unmarshal JSON data: %v", err)
	}

	// Compare the unmarshaled data with the expected output
	if !reflect.DeepEqual(actualOutput, expectedOutput) {
		t.Errorf("unexpected output. got %+v, want %+v", actualOutput, expectedOutput)
	}
}

func TestProcessExcelFile(t *testing.T) {
	// Create a temporary directory for the test file
	tempDir := t.TempDir()
	infilePath := filepath.Join(tempDir, "test.xlsx")

	// Create a new Excel file
	f := excelize.NewFile()
	sheetName := "Sheet1"
	f.SetSheetName(f.GetSheetName(0), sheetName)

	// Write observations
	f.SetCellValue(sheetName, "A1", "Observação 1")
	f.SetCellValue(sheetName, "A2", "Observação 2")

	// Write cabecalho
	cabecalho := []string{
		"SUBSTÂNCIA", "CNPJ", "LABORATÓRIO", "CÓDIGO GGREM", "REGISTRO",
		"EAN 1", "EAN 2", "EAN 3", "PRODUTO", "APRESENTAÇÃO",
		"CLASSE TERAPÊUTICA", "TIPO DE PRODUTO (STATUS DO PRODUTO)", "REGIME DE PREÇO",
		"PF Sem Impostos", "PF 0%", "PF 12%", "PF 12% ALC", "PF 17%", "PF 17% ALC",
		"PF 17,5%", "PF 17,5% ALC", "PF 18%", "PF 18% ALC", "PF 19%", "PF 19% ALC",
		"PF 19,5%", "PF 19,5% ALC", "PF 20%", "PF 20% ALC", "PF 20,5%", "PF 20,5% ALC",
		"PF 21%", "PF 21% ALC", "PF 22%", "PF 22% ALC", "PF 22,5%", "PF 22,5% ALC",
		"PF 23%", "PF 23% ALC", "PMVG Sem Impostos", "PMVG 0%", "PMVG 12%", "PMVG 12% ALC",
		"PMVG 17%", "PMVG 17% ALC", "PMVG 17,5%", "PMVG 17,5% ALC", "PMVG 18%",
		"PMVG 18% ALC", "PMVG 19%", "PMVG 19% ALC", "PMVG 19,5%", "PMVG 19,5% ALC",
		"PMVG 20%", "PMVG 20% ALC", "PMVG 20,5%", "PMVG 20,5% ALC", "PMVG 21%",
		"PMVG 21% ALC", "PMVG 22%", "PMVG 22% ALC", "PMVG 22,5%", "PMVG 22,5% ALC",
		"PMVG 23%", "PMVG 23% ALC", "RESTRIÇÃO HOSPITALAR", "CAP", "CONFAZ 87",
		"ICMS 0%", "ANÁLISE RECURSAL", "LISTA DE CONCESSÃO DE CRÉDITO TRIBUTÁRIO (PIS/COFINS)",
		"COMERCIALIZAÇÃO 2024", "TARJA",
	}
	// Ensure the header has the correct number of columns
	if len(cabecalho) != 73 {
		t.Fatalf("expected 73 header columns, got %d", len(cabecalho))
	}
	// Set the header row
	f.SetSheetRow(sheetName, "A3", &cabecalho)

	// Write header
	for i, header := range cabecalho {
		cell, _ := excelize.CoordinatesToCellName(i+1, 3)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data rows
	f.SetCellValue(sheetName, "A4", "IBUPROFENO")         // SUBSTÂNCIA
	f.SetCellValue(sheetName, "B4", "12.345.678/0001-90") // CNPJ
	f.SetCellValue(sheetName, "C4", "LAB A")              // LABORATÓRIO
	f.SetCellValue(sheetName, "J4", "COM REV")            // APRESENTAÇÃO
	f.SetCellValue(sheetName, "N4", "10,50")              // PF Sem Impostos
	f.SetCellValue(sheetName, "BN4", "Sim")               // RESTRIÇÃO HOSPITALAR
	f.SetCellValue(sheetName, "BO4", "Não")               // CAP
	f.SetCellValue(sheetName, "BP4", "Sim")               // CONFAZ 87

	f.SetCellValue(sheetName, "A5", "PARACETAMOL")        // SUBSTÂNCIA
	f.SetCellValue(sheetName, "B5", "98.765.432/0001-10") // CNPJ
	f.SetCellValue(sheetName, "C5", "LAB B")              // LABORATÓRIO
	f.SetCellValue(sheetName, "J5", "GOTAS")              // APRESENTAÇÃO
	f.SetCellValue(sheetName, "N5", "25,00*")             // PF Sem Impostos
	f.SetCellValue(sheetName, "BN5", "Não")               // RESTRIÇÃO HOSPITALAR
	f.SetCellValue(sheetName, "BO5", "Sim")               // CAP

	// Save the temporary file
	if err := f.SaveAs(infilePath); err != nil {
		t.Fatalf("failed to save temporary excel file: %v", err)
	}

	// Expected output
	expectedOutput := Output{
		Metadados: Metadados{
			Data:            "2025-07-03",
			DataAtualizacao: "2025-07-04",
			Observacoes:     []string{"Observação 1", "Observação 2"},
		},
		Medicamentos: []Medicamento{
			{
				"SUBSTÂNCIA":                          "IBUPROFENO",
				"CNPJ":                                "12.345.678/0001-90",
				"LABORATÓRIO":                         "LAB A",
				"CÓDIGO GGREM":                        nil,
				"REGISTRO":                            nil,
				"EAN 1":                               nil,
				"EAN 2":                               nil,
				"EAN 3":                               nil,
				"PRODUTO":                             nil,
				"APRESENTAÇÃO":                        "COM REV",
				"CLASSE TERAPÊUTICA":                  nil,
				"TIPO DE PRODUTO (STATUS DO PRODUTO)": nil,
				"REGIME DE PREÇO":                     nil,
				"PF Sem Impostos":                     10.50,
				"PF 0%":                               nil,
				"PF 12%":                              nil,
				"PF 12% ALC":                          nil,
				"PF 17%":                              nil,
				"PF 17% ALC":                          nil,
				"PF 17,5%":                            nil,
				"PF 17,5% ALC":                        nil,
				"PF 18%":                              nil,
				"PF 18% ALC":                          nil,
				"PF 19%":                              nil,
				"PF 19% ALC":                          nil,
				"PF 19,5%":                            nil,
				"PF 19,5% ALC":                        nil,
				"PF 20%":                              nil,
				"PF 20% ALC":                          nil,
				"PF 20,5%":                            nil,
				"PF 20,5% ALC":                        nil,
				"PF 21%":                              nil,
				"PF 21% ALC":                          nil,
				"PF 22%":                              nil,
				"PF 22% ALC":                          nil,
				"PF 22,5%":                            nil,
				"PF 22,5% ALC":                        nil,
				"PF 23%":                              nil,
				"PF 23% ALC":                          nil,
				"PMVG Sem Impostos":                   nil,
				"PMVG 0%":                             nil,
				"PMVG 12%":                            nil,
				"PMVG 12% ALC":                        nil,
				"PMVG 17%":                            nil,
				"PMVG 17% ALC":                        nil,
				"PMVG 17,5%":                          nil,
				"PMVG 17,5% ALC":                      nil,
				"PMVG 18%":                            nil,
				"PMVG 18% ALC":                        nil,
				"PMVG 19%":                            nil,
				"PMVG 19% ALC":                        nil,
				"PMVG 19,5%":                          nil,
				"PMVG 19,5% ALC":                      nil,
				"PMVG 20%":                            nil,
				"PMVG 20% ALC":                        nil,
				"PMVG 20,5%":                          nil,
				"PMVG 20,5% ALC":                      nil,
				"PMVG 21%":                            nil,
				"PMVG 21% ALC":                        nil,
				"PMVG 22%":                            nil,
				"PMVG 22% ALC":                        nil,
				"PMVG 22,5%":                          nil,
				"PMVG 22,5% ALC":                      nil,
				"PMVG 23%":                            nil,
				"PMVG 23% ALC":                        nil,
				"RESTRIÇÃO HOSPITALAR":                true,
				"CAP":                                 false,
				"CONFAZ 87":                           true,
				"ICMS 0%":                             nil,
				"ANÁLISE RECURSAL":                    nil,
				"LISTA DE CONCESSÃO DE CRÉDITO TRIBUTÁRIO (PIS/COFINS)": nil,
				"COMERCIALIZAÇÃO 2024":                                  nil,
				"TARJA":                                                 nil,
			},
			{
				"SUBSTÂNCIA":                          "PARACETAMOL",
				"CNPJ":                                "98.765.432/0001-10",
				"LABORATÓRIO":                         "LAB B",
				"CÓDIGO GGREM":                        nil,
				"REGISTRO":                            nil,
				"EAN 1":                               nil,
				"EAN 2":                               nil,
				"EAN 3":                               nil,
				"PRODUTO":                             nil,
				"APRESENTAÇÃO":                        "GOTAS",
				"CLASSE TERAPÊUTICA":                  nil,
				"TIPO DE PRODUTO (STATUS DO PRODUTO)": nil,
				"REGIME DE PREÇO":                     nil,
				"PF Sem Impostos":                     25.00,
				"PF 0%":                               nil,
				"PF 12%":                              nil,
				"PF 12% ALC":                          nil,
				"PF 17%":                              nil,
				"PF 17% ALC":                          nil,
				"PF 17,5%":                            nil,
				"PF 17,5% ALC":                        nil,
				"PF 18%":                              nil,
				"PF 18% ALC":                          nil,
				"PF 19%":                              nil,
				"PF 19% ALC":                          nil,
				"PF 19,5%":                            nil,
				"PF 19,5% ALC":                        nil,
				"PF 20%":                              nil,
				"PF 20% ALC":                          nil,
				"PF 20,5%":                            nil,
				"PF 20,5% ALC":                        nil,
				"PF 21%":                              nil,
				"PF 21% ALC":                          nil,
				"PF 22%":                              nil,
				"PF 22% ALC":                          nil,
				"PF 22,5%":                            nil,
				"PF 22,5% ALC":                        nil,
				"PF 23%":                              nil,
				"PF 23% ALC":                          nil,
				"PMVG Sem Impostos":                   nil,
				"PMVG 0%":                             nil,
				"PMVG 12%":                            nil,
				"PMVG 12% ALC":                        nil,
				"PMVG 17%":                            nil,
				"PMVG 17% ALC":                        nil,
				"PMVG 17,5%":                          nil,
				"PMVG 17,5% ALC":                      nil,
				"PMVG 18%":                            nil,
				"PMVG 18% ALC":                        nil,
				"PMVG 19%":                            nil,
				"PMVG 19% ALC":                        nil,
				"PMVG 19,5%":                          nil,
				"PMVG 19,5% ALC":                      nil,
				"PMVG 20%":                            nil,
				"PMVG 20% ALC":                        nil,
				"PMVG 20,5%":                          nil,
				"PMVG 20,5% ALC":                      nil,
				"PMVG 21%":                            nil,
				"PMVG 21% ALC":                        nil,
				"PMVG 22%":                            nil,
				"PMVG 22% ALC":                        nil,
				"PMVG 22,5%":                          nil,
				"PMVG 22,5% ALC":                      nil,
				"PMVG 23%":                            nil,
				"PMVG 23% ALC":                        nil,
				"RESTRIÇÃO HOSPITALAR":                false,
				"CAP":                                 true,
				"CONFAZ 87":                           nil,
				"ICMS 0%":                             nil,
				"ANÁLISE RECURSAL":                    nil,
				"LISTA DE CONCESSÃO DE CRÉDITO TRIBUTÁRIO (PIS/COFINS)": nil,
				"COMERCIALIZAÇÃO 2024":                                  nil,
				"TARJA":                                                 nil,
			},
		},
		Laboratorios: map[string]string{
			"12.345.678/0001-90": "LAB A",
			"98.765.432/0001-10": "LAB B",
		},
		Apresentacoes: []string{"COM REV", "GOTAS"},
	}

	// Parse the date strings to time.Time
	data, err := time.Parse("2006-01-02", "2025-07-03")
	if err != nil {
		t.Fatalf("failed to parse data: %v", err)
	}
	dataAtualizacao, err := time.Parse("2006-01-02", "2025-07-04")
	if err != nil {
		t.Fatalf("failed to parse dataAtualizacao: %v", err)
	}
	// Process the file
	output, err := processExcelFile(infilePath, data, dataAtualizacao)
	if err != nil {
		t.Fatalf("processExcelFile failed: %v", err)
	}

	// Compare the output with the expected result
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("unexpected output. got %+v, want %+v", output, expectedOutput)
	}
}
