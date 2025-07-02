package main

import "testing"

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
