# CMED Parser

Este projeto consiste em um parser da tabela de preços de medicamentos da CMED (Câmara de Regulação do Mercado de Medicamentos), que converte os dados de um arquivo `.xlsx` para um formato `.json` estruturado.

## Funcionalidades

- Lê os dados de medicamentos a partir de um arquivo `.xlsx` fornecido.
- Extrai metadados da planilha, como observações e datas.
- Analisa e processa cada linha da tabela de medicamentos.
- Realiza a limpeza e padronização dos dados, convertendo valores monetários para números e campos "Sim"/"Não" para booleanos.
- Agrega listas únicas de laboratórios (por CNPJ) e apresentações de medicamentos.
- Gera um arquivo `.json` único contendo todos os dados processados de forma organizada.

## Formato de Entrada

O programa aceita como entrada apenas arquivos no formato `.xlsx`. Caso o seu arquivo esteja em outro formato (como `.xls` ou `.ods`), é necessário convertê-lo para `.xlsx` antes de usar o parser. Você pode fazer essa conversão utilizando o Microsoft Excel, LibreOffice Calc ou outro programa compatível.

## Como Usar

Existem duas maneiras de executar o parser: utilizando o `go run` para compilar e executar o código-fonte diretamente, ou utilizando o binário pré-compilado específico para o seu sistema operacional.

### Usando `go run`

Para executar o parser diretamente do código-fonte, utilize o seguinte comando:

```bash
go run main.go [flags] <caminho/para/arquivo.xlsx>
```

### Usando o binário pré-compilado

Você pode encontrar os binários pré-compilados para Linux, macOS e Windows na seção de [Releases](https://github.com/elvisdiniz/cmed-parser/releases) do projeto. Após baixar e descompactar o arquivo correspondente ao seu sistema operacional, você pode executar o parser da seguinte forma:

#### Linux/macOS

```bash
./cmed-parser-[linux|macos]-amd64 [flags] <caminho/para/arquivo.xlsx>
```

#### Windows

```bash
.\cmed-parser-windows-amd64.exe [flags] <caminho/para/arquivo.xlsx>
```

### Flags

- `--data`: (Opcional) Especifica a data da planilha no formato `AAAA-MM-DD`. Se omitido, utiliza a data atual.
- `--data-atualizacao`: (Opcional) Especifica a data de atualização da planilha no formato `AAAA-MM-DD`. Se omitido, utiliza o mesmo valor da flag `--data`.

### Exemplo

#### `go run`

```bash
go run main.go --data 2024-07-25 ./lista-de-precos.xlsx
```

#### Binário (Linux/macOS)

```bash
./cmed-parser-linux-amd64 --data 2024-07-25 ./lista-de-precos.xlsx
```

Este comando irá processar o arquivo `lista-de-precos.xlsx` e gerar um novo arquivo chamado `lista-de-precos.json` no mesmo diretório.

## Estrutura do JSON de Saída

O arquivo de saída (`.json`) é estruturado da seguinte forma:

```json
{
  "metadados": {
    "data": "2024-07-25",
    "data-atualizacao": "2024-07-25",
    "observacoes": [
      "Observação 1 da planilha...",
      "Observação 2 da planilha..."
    ]
  },
  "medicamentos": [
    {
      "SUBSTÂNCIA": "NOME DA SUBSTÂNCIA",
      "CNPJ": "00.000.000/0000-00",
      "LABORATÓRIO": "NOME DO LABORATÓRIO",
      "CÓDIGO GGREM": "0000000000000",
      "REGISTRO": "0000000000000",
      "EAN 1": "0000000000000",
      // ... outros campos do medicamento
      "PF 18%": 123.45,
      "RESTRIÇÃO HOSPITALAR": true,
      // ...
    }
  ],
  "laboratorios": {
    "00.000.000/0000-00": "NOME DO LABORATÓRIO",
    "00.000.000/0000-01": "NOME DE OUTRO LABORATÓRIO"
  },
  "apresentacoes": [
    "APRESENTACAO SEM ACENTO 1",
    "APRESENTACAO SEM ACENTO 2"
  ]
}
```
