package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/xuri/excelize/v2"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	PrincipioAtivo                  = "SUBSTÂNCIA"
	CNPJ                            = "CNPJ"
	Laboratorio                     = "LABORATÓRIO"
	CodigoGGREM                     = "CÓDIGO GGREM"
	Registro                        = "REGISTRO"
	EAN1                            = "EAN 1"
	EAN2                            = "EAN 2"
	EAN3                            = "EAN 3"
	Produto                         = "PRODUTO"
	Apresentacao                    = "APRESENTAÇÃO"
	ClasseTerapeutica               = "CLASSE TERAPÊUTICA"
	Tipo                            = "TIPO DE PRODUTO (STATUS DO PRODUTO)"
	RegimePreco                     = "REGIME DE PREÇO"
	PFSemImpostos                   = "PF Sem Impostos"
	PF0                             = "PF 0%"
	PF12                            = "PF 12%"
	PF12ALC                         = "PF 12% ALC"
	PF17                            = "PF 17%"
	PF17ALC                         = "PF 17% ALC"
	PF175                           = "PF 17,5%"
	PF175ALC                        = "PF 17,5% ALC"
	PF18                            = "PF 18%"
	PF18ALC                         = "PF 18% ALC"
	PF19                            = "PF 19%"
	PF19ALC                         = "PF 19% ALC"
	PF195                           = "PF 19,5%"
	PF195ALC                        = "PF 19,5% ALC"
	PF20                            = "PF 20%"
	PF20ALC                         = "PF 20% ALC"
	PF205                           = "PF 20,5%"
	PF205ALC                        = "PF 20,5% ALC"
	PF21                            = "PF 21%"
	PF21ALC                         = "PF 21% ALC"
	PF22                            = "PF 22%"
	PF22ALC                         = "PF 22% ALC"
	PF225                           = "PF 22,5%"
	PF225ALC                        = "PF 22,5% ALC"
	PF23                            = "PF 23%"
	PF23ALC                         = "PF 23% ALC"
	PMVGSemImpostos                 = "PMVG Sem Impostos"
	PMVG0                           = "PMVG 0%"
	PMVG12                          = "PMVG 12%"
	PMVG12ALC                       = "PMVG 12% ALC"
	PMVG17                          = "PMVG 17%"
	PMVG17ALC                       = "PMVG 17% ALC"
	PMVG175                         = "PMVG 17,5%"
	PMVG175ALC                      = "PMVG 17,5% ALC"
	PMVG18                          = "PMVG 18%"
	PMVG18ALC                       = "PMVG 18% ALC"
	PMVG19                          = "PMVG 19%"
	PMVG19ALC                       = "PMVG 19% ALC"
	PMVG195                         = "PMVG 19,5%"
	PMVG195ALC                      = "PMVG 19,5% ALC"
	PMVG20                          = "PMVG 20%"
	PMVG20ALC                       = "PMVG 20% ALC"
	PMVG205                         = "PMVG 20,5%"
	PMVG205ALC                      = "PMVG 20,5% ALC"
	PMVG21                          = "PMVG 21%"
	PMVG21ALC                       = "PMVG 21% ALC"
	PMVG22                          = "PMVG 22%"
	PMVG22ALC                       = "PMVG 22% ALC"
	PMVG225                         = "PMVG 22,5%"
	PMVG225ALC                      = "PMVG 22,5% ALC"
	PMVG23                          = "PMVG 23%"
	PMVG23ALC                       = "PMVG 23% ALC"
	RestricaoHospitalar             = "RESTRIÇÃO HOSPITALAR"
	CAP                             = "CAP"
	Confaz87                        = "CONFAZ 87"
	ICMS0                           = "ICMS 0%"
	AnaliseRecursal                 = "ANÁLISE RECURSAL"
	ListaConcessaoCreditoTributario = "LISTA DE CONCESSÃO DE CRÉDITO TRIBUTÁRIO (PIS/COFINS)"
	Comercializacao2024             = "COMERCIALIZAÇÃO 2024"
	Tarja                           = "TARJA"
)

var cabecalho = []string{
	PrincipioAtivo,
	CNPJ,
	Laboratorio,
	CodigoGGREM,
	Registro,
	EAN1,
	EAN2,
	EAN3,
	Produto,
	Apresentacao,
	ClasseTerapeutica,
	Tipo,
	RegimePreco,
	PFSemImpostos,
	PF0,
	PF12,
	PF12ALC,
	PF17,
	PF17ALC,
	PF175,
	PF175ALC,
	PF18,
	PF18ALC,
	PF19,
	PF19ALC,
	PF195,
	PF195ALC,
	PF20,
	PF20ALC,
	PF205,
	PF205ALC,
	PF21,
	PF21ALC,
	PF22,
	PF22ALC,
	PF225,
	PF225ALC,
	PF23,
	PF23ALC,
	PMVGSemImpostos,
	PMVG0,
	PMVG12,
	PMVG12ALC,
	PMVG17,
	PMVG17ALC,
	PMVG175,
	PMVG175ALC,
	PMVG18,
	PMVG18ALC,
	PMVG19,
	PMVG19ALC,
	PMVG195,
	PMVG195ALC,
	PMVG20,
	PMVG20ALC,
	PMVG205,
	PMVG205ALC,
	PMVG21,
	PMVG21ALC,
	PMVG22,
	PMVG22ALC,
	PMVG225,
	PMVG225ALC,
	PMVG23,
	PMVG23ALC,
	RestricaoHospitalar,
	CAP,
	Confaz87,
	ICMS0,
	AnaliseRecursal,
	ListaConcessaoCreditoTributario,
	Comercializacao2024,
	Tarja,
}

type Metadados struct {
	Data            string   `json:"data"`
	DataAtualizacao string   `json:"data-atualizacao,omitempty"`
	Observacoes     []string `json:"observacoes"`
}

type Medicamento map[string]interface{}

type Output struct {
	Metadados     Metadados         `json:"metadados"`
	Medicamentos  []Medicamento     `json:"medicamentos"`
	Laboratorios  map[string]string `json:"laboratorios"`
	Apresentacoes []string          `json:"apresentacoes"`
}

func writeJSONFile(output Output, infilePath string) error {
	outfilePath := strings.TrimSuffix(infilePath, filepath.Ext(infilePath)) + ".json"
	jsonFile, err := os.Create(outfilePath)
	if err != nil {
		return fmt.Errorf("failed to create json file: %w", err)
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}

	fmt.Printf("Arquivo %s criado!\n", outfilePath)
	return nil
}

func processExcelFile(infilePath, data, dataAtualizacao string) (Output, error) {
	f, err := excelize.OpenFile(infilePath)
	if err != nil {
		return Output{}, fmt.Errorf("failed to open excel file: %w", err)
	}

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return Output{}, fmt.Errorf("failed to get rows from sheet: %w", err)
	}

	var planilhaObservacoes []string
	laboratoriosList := make(map[string]string)
	var apresentacaoList []string
	var medicamentosList []Medicamento
	linhaCabecalho := -1

	for i, row := range rows {
		if linhaCabecalho == -1 {
			if len(row) > 0 && row[0] == PrincipioAtivo {
				linhaCabecalho = i
			} else if len(row) > 0 {
				planilhaObservacoes = append(planilhaObservacoes, row[0])
			}
			continue
		}

		medicamento := make(Medicamento)
		for j, header := range cabecalho {
			var value any
			if j < len(row) {
				value = strings.TrimSpace(row[j])
			} else {
				value = ""
			}

			medicamento[header] = processaValorCelula(value, header)
		}

		medicamentosList = append(medicamentosList, medicamento)

		if medicamento[CNPJ] != nil {
			cnpjLaboratorio := medicamento[CNPJ].(string)
			if _, ok := laboratoriosList[cnpjLaboratorio]; !ok {
				laboratoriosList[cnpjLaboratorio] = medicamento[Laboratorio].(string)
			}
		}

		if medicamento[Apresentacao] != nil {
			apresentacao := removeAccents(medicamento[Apresentacao].(string))
			found := slices.Contains(apresentacaoList, apresentacao)
			if !found {
				apresentacaoList = append(apresentacaoList, apresentacao)
			}
		}
	}

	output := Output{
		Metadados: Metadados{
			Data:            data,
			DataAtualizacao: dataAtualizacao,
			Observacoes:     planilhaObservacoes,
		},
		Medicamentos:  medicamentosList,
		Laboratorios:  laboratoriosList,
		Apresentacoes: apresentacaoList,
	}

	return output, nil
}

func main() {
	data := flag.String("data", time.Now().Format("2006-01-02"), "Data da planilha no formato AAAA-MM-DD")
	dataAtualizacao := flag.String("data-atualizacao", "", "Data de atualização da planilha no formato AAAA-MM-DD")
	flag.Parse()

	if *dataAtualizacao == "" {
		*dataAtualizacao = *data
	}

	if len(flag.Args()) != 1 {
		log.Fatal("Uso: go run main.go [flags] <arquivo.xlsx>")
	}

	infilePath := flag.Args()[0]
	ext := filepath.Ext(infilePath)
	if ext != ".xlsx" {
		log.Fatal("O arquivo de entrada deve ser .xlsx")
	}

	output, err := processExcelFile(infilePath, *data, *dataAtualizacao)
	if err != nil {
		log.Fatal(err)
	}

	if err := writeJSONFile(output, infilePath); err != nil {
		log.Fatal(err)
	}
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

func processaValorCelula(value any, header string) any {
	valoresRegex := regexp.MustCompile(`^(PF |PMVG )([0-2]|S)`)
	realRegex := regexp.MustCompile(`^[0-9]+([\.,])[0-9]+\*?$`)

	if header == PrincipioAtivo && value == nil {
		return ""
	}

	strValue, ok := value.(string)
	if !ok {
		return value
	}

	if strValue == "-" || strValue == "" {
		return nil
	} else if valoresRegex.MatchString(header) && realRegex.MatchString(strValue) {
		strValue = strings.Replace(strValue, ",", ".", 1)
		strValue = strings.Replace(strValue, "*", "", 1)
		floatValue, err := strconv.ParseFloat(strValue, 64)
		if err == nil {
			return floatValue
		}
	} else if header == CAP || header == Confaz87 || header == ICMS0 || header == RestricaoHospitalar || header == Comercializacao2024 {
		return strings.ToLower(strValue) == "sim"
	}

	return strValue
}
