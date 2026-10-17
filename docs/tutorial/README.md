# Tutorial — AKMA sobre o free5GC

Este tutorial descreve o que fazer **depois de instalar o free5GC** pelo guia oficial para usar AKMA (3GPP TS 33.535). Ele cresce a cada MVP; esta versão corresponde ao MVP-1.

| Parte | Conteúdo | Disponível |
|---|---|---|
| 1 | Pré-requisitos | sim |
| 2 | Conferir as derivações de chave | sim |
| 3 | Instalar o AAnF no free5GC | sim |
| 4 | Habilitar AKMA no UDM e no AUSF; provisionar um assinante | MVP-2 |
| 5 | Escrever um AF que usa AKMA | MVP-3 |
| 6 | Preparar o UE (UERANSIM) | MVP-4 |
| 7 | Executar o fluxo de ponta a ponta | MVP-5 |

---

## 1. Pré-requisitos

| Item | Versão usada no projeto |
|---|---|
| free5GC | v4.2.3, instalado pelo guia oficial, com `TestRegistration` passando |
| Go | 1.26.2 — a mesma do free5GC |
| Python | 3.10 ou mais recente (para o oráculo de referência) |
| UERANSIM | v3.2.6 (a partir da parte 6) |

Topologia de referência (ADR-0001-R2): o free5GC e o AAnF numa máquina; o AF e este repositório em outra. Uma única máquina também serve para as partes 1 e 2.

Se o Go não estiver instalado, instale-o num diretório de usuário, sem privilégios:

```bash
mkdir -p ~/sdk/go1.26.2 && curl -fsSL https://go.dev/dl/go1.26.2.linux-amd64.tar.gz | tar -xz -C ~/sdk/go1.26.2 --strip-components=1
```

```bash
export PATH=$HOME/sdk/go1.26.2/bin:$PATH
```

Confira o SHA-256 do arquivo com o publicado em https://go.dev/dl/ antes de usar. O projeto usou `990e6b4bbba816dc3ee129eaeaf4b42f17c2800b88a2166c265ac1a200262282` para `go1.26.2.linux-amd64.tar.gz`.

---

## 2. Conferir as derivações de chave

O AKMA deriva três valores da chave K_AUSF que a autenticação primária do 5G deixa no AUSF e no UE (TS 33.535 Anexo A):

- **K_AKMA**, a chave âncora, guardada no AAnF;
- **A-TID**, parte do identificador A-KID que o UE apresenta à aplicação;
- **K_AF**, a chave que a aplicação (AF) recebe do AAnF.

Todas usam a KDF de TS 33.220 Anexo B.2 — a mesma que o free5GC já usa para K_SEAF. Detalhes em `docs/shared/kdf.md`.

### 2.1 Rodar o oráculo de referência

```bash
python3 reference/kdf/akma_kdf.py
```

A saída mostra a cadeia `S` e o resultado de cada derivação para o exemplo do baseline. O valor de K_AKMA deve ser `d6aa856e127fb8de3eb1a1d0cebf52a7df6b69b4a6c68c21c18c1c27d466dbba`.

### 2.2 Testar a biblioteca contra os vetores

```bash
cd shared/akma && go test ./...
```

```bash
cd shared/crosscheck && go test ./...
```

O primeiro compara a biblioteca Go com os vetores do oráculo; o segundo, com a KDF do próprio free5GC (`free5gc/util`). Os vetores são do projeto, não do 3GPP — ver `docs/shared/vetores.md`.

### 2.3 Dois erros que não dão mensagem de erro

1. **O SUPI vai à KDF sem o prefixo `imsi-`.** O free5GC trafega `imsi-208930000000001`; a KDF usa `208930000000001` (TS 33.501 A.7.0). Com o prefixo, todas as chaves saem erradas.
2. **O AF_ID da KDF não é o `afId` da API.** A API transporta `af1.akma.example.FF00000001`; a KDF usa o FQDN seguido de cinco octetos brutos `FF 00 00 00 01` (TS 33.535 A.4).

Os dois casos estão em `shared/akma/testdata/vectors.json`, seção `counterexamples`, com o resultado errado que produzem.

### 2.4 Um erro que só aparece para alguns assinantes

O Routing Indicator é uma cadeia de 1 a 4 dígitos em que cada dígito conta: `"0000"`, `"012"` e `"12"` são valores diferentes (TS 23.003 §2.2B). Guardá-lo como número apaga os zeros à esquerda e quebra o A-KID só dos assinantes cujo RID começa com zero. O UE do testbed usa `"0000"`.

---

## 3. Instalar o AAnF no free5GC

O AAnF (*AKMA Anchor Function*) é a NF que guarda a chave âncora K_AKMA de cada assinante e entrega a K_AF a cada aplicação autorizada (TS 33.535 §4.2.1). O free5GC não tem essa NF. Esta parte acrescenta uma, com a mesma estrutura das outras NFs, **sem alterar nenhum arquivo do free5GC**. Ela ganha um binário em `bin/` e uma configuração em `config/`. Os detalhes de cada decisão estão em `docs/free5gc-side/aanf.md`.

Nesta parte o AAnF ainda não recebe chaves do AUSF (isso é a parte 4). Ele é testado com a ferramenta `seed`, que faz o papel do AUSF e do AF.

### 3.1 Preparar o `openapi` e compilar

O `free5gc/openapi` v1.2.4 não tem os tipos de dados do serviço Naanf_AKMA. Dois *patches* do repositório os acrescentam numa cópia local do módulo, sem tocar no cache do Go (`docs/free5gc-side/openapi.md`). O *script* abaixo prepara essa cópia e compila o AAnF e o `seed` como binários estáticos em `free5gc/build/bin/`:

```bash
free5gc/scripts/build-aanf.sh
```

Para conferir antes de instalar:

```bash
cd free5gc/aanf && go test ./...
```

### 3.2 Configurar

Copie `free5gc/config/aanfcfg.yaml` e ajuste:

- `sbi.registerIPv4` e `sbi.bindingIPv4`: um endereço da máquina do free5GC que a máquina do AF alcance. Não use um *loopback* `127.0.0.x`, a menos que o AF rode na mesma máquina. Porta 8000, como as outras NFs.
- `nrfUri` e `nrfCertPem`: os mesmos das outras NFs (`http://127.0.0.10:8000`, `cert/nrf.pem`).
- `plmnSupportList`: a PLMN da sua rede.
- `routingIndicators`: os RIDs dos assinantes, **entre aspas** (`"0000"`). Com a lista vazia, o AAnF atende qualquer RID.
- `akma.authorizedAfs`: cada AF que pode pedir K_AF, com o AF_ID (`<FQDN>.<identificador Ua* em 10 hexadecimais>`) e o endereço de origem do AF.
- `nfInstanceId`: mantenha um UUID fixo. Com OAuth2 ligado no NRF, cada ID novo faz o NRF criar um certificado `cert/aanf_<ID>.pem`.

### 3.3 Instalar e iniciar

Copie `free5gc/build/bin/aanf` para `~/free5gc/bin/aanf` e a configuração para `~/free5gc/config/aanfcfg.yaml`. Se o repositório está em outra máquina, o *script* do projeto faz isso por SSH e não sobrescreve uma configuração já instalada:

```bash
free5gc/scripts/deploy.sh usuario@maquina-do-free5gc
```

Com o free5GC em execução (pelo menos o NRF), inicie o AAnF na máquina do free5GC:

```bash
free5gc/scripts/run-aanf.sh start
```

Ou, a partir da outra máquina:

```bash
ssh usuario@maquina-do-free5gc bash -s -- start < free5gc/scripts/run-aanf.sh
```

O log fica em `~/free5gc/log/aanf/aanf.log`. As linhas esperadas são `OAuth2 setting receive from NRF: true` e `Start SBI server (listen on …:8000)`. `run-aanf.sh stop` encerra o AAnF, que antes se descadastra do NRF.

### 3.4 Verificar o registro no NRF

```bash
mongosh --quiet free5gc --eval 'printjson(db.NfProfile.find({nfType: "AANF"}, {_id: 0}).toArray())'
```

O perfil deve ter `nfType: 'AANF'` e o serviço `naanf-akma`, versão `v1`, no endereço configurado. Duas limitações do NRF do free5GC v4.2.3 aparecem aqui (`docs/free5gc-side/registro-nrf.md` §5):

- o perfil guardado não tem `aanfInfoList`;
- a descoberta `target-nf-type=AANF` é recusada.

Por isso, nas próximas partes, o AUSF e o AF recebem o endereço do AAnF por configuração, o que TS 33.535 §6.7 permite.

### 3.5 Testar as três operações

Os valores abaixo são do vetor `testbed-ue` (`shared/akma/testdata/vectors.json`). O `seed` precisa de um *token* OAuth2 para registrar, e o NRF só é alcançável na máquina do free5GC. Para obtê-lo, o `seed` se registra no NRF como um AF temporário e se descadastra ao sair.

Na máquina do free5GC, no papel do AUSF:

```bash
seed register --nrf http://127.0.0.10:8000 --supi imsi-208930000000001 --k-ausf 141ef548f86e3a2a1ad949331caa9e6bb12401afa321bde4f9f23b84870dd654 --rid 0000 --mcc 208 --mnc 93
```

A saída deve trazer o A-KID `rid0000.atidiz1FN8DkMGyxwtfTVYsKUsWnzOJdBOLlQEiQlK2gHNA=@5gc.mnc093.mcc208.3gppnetwork.org`.

Na máquina do AF (o endereço dela precisa estar em `authorizedAfs`), no papel do AF:

```bash
seed retrieve --aanf http://192.168.16.130:8000 --af-id af1.akma.example.FF00000001 --akid 'rid0000.atidiz1FN8DkMGyxwtfTVYsKUsWnzOJdBOLlQEiQlK2gHNA=@5gc.mnc093.mcc208.3gppnetwork.org'
```

A K_AF deve ser `00c16845c2e4e1256755b84bb06d94e2e659d8ba53c59e874fc149bffd5ea04b`, a mesma que o oráculo da parte 2 calcula. O `expiry` deve estar uma hora à frente, em UTC.

Respostas que indicam configuração errada:

| Resposta | Causa provável |
|---|---|
| `401` no `register` | sem `--nrf`, ou o `nrfCertPem` do AAnF não é o certificado do NRF em uso |
| `403` sem causa no `retrieve` | o AF_ID não está em `authorizedAfs`, ou o pedido saiu de um endereço que não está na lista daquele AF_ID |
| `403 (K_AKMA_NOT_PRESENT)` | o A-KID não foi registrado, foi substituído por um registro mais novo do mesmo SUPI ou o AAnF reiniciou (o contexto fica só em memória) |

Por fim, para remover o contexto: `seed remove --nrf http://127.0.0.10:8000 --supi imsi-208930000000001`.
