const XLSX = require("xlsx")
const fs = require("fs")
const accents = require("remove-accents")
const prompt = require("prompt")
const tmp = require("tmp")
const zl = require("zip-lib")

const PRINCIPIO_ATIVO = "SUBSTÂNCIA"
const CNPJ = "CNPJ"
const LABORATORIO = "LABORATÓRIO"
const CODIGO_GGREM = "CÓDIGO GGREM"
const REGISTRO = "REGISTRO"
const EAN_1 = "EAN 1"
const EAN_2 = "EAN 2"
const EAN_3 = "EAN 3"
const PRODUTO = "PRODUTO"
const APRESENTACAO = "APRESENTAÇÃO"
const CLASSE_TERAPEUTICA = "CLASSE TERAPÊUTICA"
const TIPO = "TIPO DE PRODUTO (STATUS DO PRODUTO)"
const REGIME_DE_PRECO = "REGIME DE PREÇO"
const PF_SEM_IMPOSTOS = "PF Sem Impostos"
const PF_0 = "PF 0%"
const PF_12 = "PF 12%"
const PF_12_ALC = "PF 12% ALC"
const PF_17 = "PF 17%"
const PF_17_ALC = "PF 17% ALC"
const PF_175 = "PF 17,5%"
const PF_175_ALC = "PF 17,5% ALC"
const PF_18 = "PF 18%"
const PF_18_ALC = "PF 18% ALC"
const PF_19 = "PF 19%"
const PF_19_ALC = "PF 19% ALC"
const PF_195 = "PF 19,5%"
const PF_195_ALC = "PF 19,5% ALC"
const PF_20 = "PF 20%"
const PF_20_ALC = "PF 20% ALC"
const PF_205 = "PF 20,5%"
const PF_21 = "PF 21%"
const PF_21_ALC = "PF 21% ALC"
const PF_22 = "PF 22%"
const PF_22_ALC = "PF 22% ALC"
const PMVG_SEM_IMPOSTOS = "PMVG Sem Imposto"
const PMVG_0 = "PMVG 0%"
const PMVG_12 = "PMVG 12%"
const PMVG_12_ALC = "PMVG 12% ALC"
const PMVG_17 = "PMVG 17%"
const PMVG_17_ALC = "PMVG 17% ALC"
const PMVG_175 = "PMVG 17,5%"
const PMVG_175_ALC = "PMVG 17,5% ALC"
const PMVG_18 = "PMVG 18%"
const PMVG_18_ALC = "PMVG 18% ALC"
const PMVG_19 = "PMVG 19%"
const PMVG_19_ALC = "PMVG 19% ALC"
const PMVG_195 = "PMVG 19,5%"
const PMVG_195_ALC = "PMVG 19,5% ALC"
const PMVG_20 = "PMVG 20%"
const PMVG_20_ALC = "PMVG 20% ALC"
const PMVG_205 = "PMVG 20,5%"
const PMVG_21 = "PMVG 21%"
const PMVG_21_ALC = "PMVG 21% ALC"
const PMVG_22 = "PMVG 22%"
const PMVG_22_ALC = "PMVG 22% ALC"
const RESTRICAO_HOSPITALAR = "RESTRIÇÃO HOSPITALAR"
const CAP = "CAP"
const CONFAZ_87 = "CONFAZ 87"
const ICMS_0 = "ICMS 0%"
const ANALISE_RECURSAL = "ANÁLISE RECURSAL"
const LISTA_CONCESSAO_CREDITO_TRIBUTARIO = "LISTA DE CONCESSÃO DE CRÉDITO TRIBUTÁRIO (PIS/COFINS)"
const COMERCIALIZACAO_2022 = "COMERCIALIZAÇÃO 2022"
const TARJA = "TARJA"

const cabecalho = [
  PRINCIPIO_ATIVO,
  CNPJ,
  LABORATORIO,
  CODIGO_GGREM,
  REGISTRO,
  EAN_1,
  EAN_2,
  EAN_3,
  PRODUTO,
  APRESENTACAO,
  CLASSE_TERAPEUTICA,
  TIPO,
  REGIME_DE_PRECO,
  PF_SEM_IMPOSTOS,
  PF_0,
  PF_12,
  PF_12_ALC,
  PF_17,
  PF_17_ALC,
  PF_175,
  PF_175_ALC,
  PF_18,
  PF_18_ALC,
  PF_19,
  PF_19_ALC,
  PF_195,
  PF_195_ALC,
  PF_20,
  PF_20_ALC,
  PF_205,
  PF_21,
  PF_21_ALC,
  PF_22,
  PF_22_ALC,
  PMVG_SEM_IMPOSTOS,
  PMVG_0,
  PMVG_12,
  PMVG_12_ALC,
  PMVG_17,
  PMVG_17_ALC,
  PMVG_175,
  PMVG_175_ALC,
  PMVG_18,
  PMVG_18_ALC,
  PMVG_19,
  PMVG_19_ALC,
  PMVG_195,
  PMVG_195_ALC,
  PMVG_20,
  PMVG_20_ALC,
  PMVG_205,
  PMVG_21,
  PMVG_21_ALC,
  PMVG_22,
  PMVG_22_ALC,
  RESTRICAO_HOSPITALAR,
  CAP,
  CONFAZ_87,
  ICMS_0,
  ANALISE_RECURSAL,
  LISTA_CONCESSAO_CREDITO_TRIBUTARIO,
  COMERCIALIZACAO_2022,
  TARJA,
]

const today = new Date()

const tmpobj = tmp.fileSync()
const zip = new zl.Zip()

const properties = [
  {
    name: "data",
    validator: /^202[0-9]-(0?[1-9]|1[012])-(0?[1-9]|([1-2][0-9]|3[0-1]))$/,
    warning: "Uma data válida deve ser informada",
    description: "Informe a data da planilha",
    default:
      today.getFullYear() +
      "-" +
      (today.getMonth() + 1) +
      "-" +
      today.getDate(),
    required: true,
  },
  {
    name: "data-atualizacao",
    validator: /^202[0-9]-(0?[1-9]|1[012])-(0?[1-9]|([1-2][0-9]|3[0-1]))$/,
    warning: "Uma data válida deve ser informada",
    description: "Se for atualização, informe a data de atualização",
    default: "",
    required: false,
  },
]

prompt.start()

prompt.get(properties, function (err, metadados) {
  if (err) {
    return onErr(err)
  }

  var wb = XLSX.readFile(process.argv[2], { type: "array" })
  var tabelaParsed = XLSX.utils.sheet_to_json(wb.Sheets[wb.SheetNames[0]], {
    header: cabecalho,
    skipHeader: true,
  })

  let cnpjRegex = new RegExp("^[0-9]{2}(.[0-9]{3}){2}/[0-9]{4}-[0-9]{2}$")
  let valoresRegex = new RegExp("^(PF |PMVG )([0-2]|S)")
  let realRegex = new RegExp("^[0-9]+(.|,)[0-9]+")

  let planilhaObservacoes = []
  let laboratoriosList = {}
  let apresentacaoList = []
  let medicamentosList = []

  let linhaCabecalho = 0

  for (let i = 0; i < tabelaParsed.length; i++) {
    if (tabelaParsed[i][CNPJ] === undefined) {
      planilhaObservacoes.push(tabelaParsed[i][PRINCIPIO_ATIVO])
      continue
    } else if (!linhaCabecalho) {
      for (let [key, value] of Object.entries(tabelaParsed[i])) {
        let trimedValue = value.trim()
        if (key !== trimedValue) {
          console.error(
            'Cabeçalho inválido: "%s" encontrado, "%s" esperado',
            trimedValue,
            key
          )
          process.exit(1)
        } else {
          linhaCabecalho = i + 1
        }
      }
      continue
    }

    if (!linhaCabecalho) {
      console.error("Cabeçalho não encontrado")
      process.exit(1)
    }

    let cnpjLaboratorio = tabelaParsed[i][CNPJ].trim()
    if (!cnpjRegex.test(cnpjLaboratorio) && !linhaCabecalho) {
      continue
    }

    cabecalho.forEach((c) => {
      if (tabelaParsed[i][c] === undefined) {
        tabelaParsed[i][c] = ""
      }
    })

    if (Object.keys(tabelaParsed[i]).length !== cabecalho.length) {
      console.error(
        "Quantidade de colunas inválida: %d encontradas. %d esperadas.",
        Object.keys(tabelaParsed[i]).length,
        cabecalho.length
      )
      process.exit(1)
    }

    for (let [key, value] of Object.entries(tabelaParsed[i])) {
      let trimedValue = value.trim()
      trimedValue = trimedValue.replace(/\s+/i, " ")

      if (trimedValue === "-" || !trimedValue) {
        trimedValue = null
      }

      if (valoresRegex.test(key) && trimedValue) {
        if (!realRegex.test(trimedValue)) {
          console.error(
            "Valor invlálido na coluna: " + key + ". Esperado decimal."
          )
          process.exit(1)
        }
        trimedValue = parseFloat(trimedValue.replace(",", "."))
      } else if (
        key === CAP ||
        key === CONFAZ_87 ||
        key === ICMS_0 ||
        key === RESTRICAO_HOSPITALAR ||
        key === COMERCIALIZACAO_2022
      ) {
        trimedValue = trimedValue.toLowerCase() === "sim" ? true : false
      } else if (key === PRINCIPIO_ATIVO && !trimedValue) {
        trimedValue = ""
      }

      tabelaParsed[i][key] = trimedValue
    }

    let laboratorio = tabelaParsed[i][LABORATORIO].trim()
    let apresentacao = accents.remove(tabelaParsed[i][APRESENTACAO].trim())

    if (laboratoriosList[cnpjLaboratorio] == undefined) {
      laboratoriosList[cnpjLaboratorio] = laboratorio
    }

    if (apresentacaoList.indexOf(apresentacao) < 0) {
      apresentacaoList.push(apresentacao)
    }

    medicamentosList.push(tabelaParsed[i])
  }

  metadados.observacoes = planilhaObservacoes

  fs.writeFileSync(
    tmpobj.name,
    JSON.stringify({
      metadados: metadados,
      medicamentos: medicamentosList,
      laboratorios: laboratoriosList,
      apresentacoes: apresentacaoList,
    })
  )

  zip.addFile(tmpobj.name, "cmed.json")

  zip.archive(process.argv[3]).then(
    function () {
      tmpobj.removeCallback()
      console.log("Arquivo " + process.argv[3] + " criado!")
    },
    function (err) {
      tmpobj.removeCallback()
      console.log(err)
    }
  )
})

function onErr(err) {
  console.log(err)
  return 1
}
