# Tutorial — AKMA sobre o free5GC

Este tutorial descreve o que fazer **depois de instalar o free5GC** pelo guia oficial para usar AKMA (3GPP TS 33.535). Ele cresce a cada MVP; esta versão corresponde ao MVP-6. O MVP-5 foi o primeiro MVP funcional. Em 19/10/2026 ele foi seguido do início numa VM com o free5GC v4.2.3 recém-instalado, e o roteiro da parte 8 terminou com 63 verificações e nenhuma falha (critério 4 do MVP-5).

| Parte | Conteúdo | Disponível |
|---|---|---|
| 1 | Pré-requisitos | sim |
| 2 | Conferir as derivações de chave | sim |
| 3 | Instalar o AAnF no free5GC | sim |
| 4 | Habilitar AKMA no UDM e no AUSF; provisionar um assinante | sim |
| 5 | Ver o AKMA na tela: free5GC e UERANSIM | sim |
| 6 | Escrever um AF que usa AKMA | sim |
| 7 | Preparar o UE (UERANSIM) para o Ua\* | sim |
| 8 | Executar o fluxo de ponta a ponta | sim |
| 9 | AF fora da rede do operador, acesso anônimo e UERANSIM em outra máquina | sim |

---

## 1. Pré-requisitos

| Item | Versão usada no projeto |
|---|---|
| free5GC | v4.2.3, instalado pelo guia oficial, com `TestRegistration` passando |
| MongoDB | o do guia, com `mongosh` (partes 3 e 4) |
| Go | 1.26.2 — a mesma do free5GC |
| Python | 3.10 ou mais recente, na máquina do AF (oráculo de referência e testes negativos) |
| UERANSIM | v3.2.6, commit `85a0fbf` do guia do free5GC, na máquina do free5GC (a partir da parte 4; com *patch* a partir da parte 5) |

Topologia de referência (ADR-0001-R2): o free5GC, o AAnF e o UERANSIM numa máquina; o AF e este repositório em outra. Uma única máquina também serve para as partes 1 e 2.

Os comandos usam os endereços do laboratório: **192.168.16.130** para a máquina do free5GC e **192.168.16.129** para a do AF. Com outros endereços, troque-os nos comandos e nestes pontos de configuração:

- no `aanfcfg.yaml` (seção 3.2);
- no `--aanf` do `seed` (seção 3.6);
- no `afcfg.yaml` (seção 6.2);
- no argumento do `e2e.sh` (seção 8.2).

Se o Go não estiver instalado, instale-o num diretório de usuário, sem privilégios:

```bash
mkdir -p ~/sdk/go1.26.2 && curl -fsSL https://go.dev/dl/go1.26.2.linux-amd64.tar.gz | tar -xz -C ~/sdk/go1.26.2 --strip-components=1
```

```bash
export PATH=$HOME/sdk/go1.26.2/bin:$PATH
```

Confira o SHA-256 do arquivo com o publicado em https://go.dev/dl/ antes de usar. O projeto usou `990e6b4bbba816dc3ee129eaeaf4b42f17c2800b88a2166c265ac1a200262282` para `go1.26.2.linux-amd64.tar.gz`.

### 1.1 Antes da parte 3

**O repositório nas duas máquinas.**

- Na máquina do AF, clone onde preferir. Os comandos sem `~/` partem da raiz dele.
- Na do free5GC, clone em `~/akma-free5gc`, ao lado de `~/free5gc`. Os *scripts* das partes 3 a 7 rodam dali, e os `go.mod` corrigidos da parte 4 apontam para esse local:

```bash
git clone https://github.com/florianoneto07/akma-free5gc.git ~/akma-free5gc
```

**SSH por chave** da máquina do AF para a do free5GC. O `deploy.sh` (seção 3.4) e o `e2e.sh` (parte 8) não pedem senha de SSH.

**O UERANSIM na máquina do free5GC.** O guia instala o UERANSIM em outra máquina e manda apontar o gNB para o AMF (§6 do guia do UERANSIM). Numa máquina só, o ajuste é um: `amfConfigs` do `config/free5gc-gnb.yaml` passa de `127.0.0.1` para `127.0.0.18`, o `ngapIpList` do AMF. O restante fica como veio.

**O assinante.** Use o assinante padrão da *webconsole* (`imsi-208930000000001`), com o tipo de código de operador **OPc**, que já vem marcado. O guia manda trocar para "OP", mas o `free5gc-ue.yaml` do UERANSIM usa `opType: 'OPC'`, e prevalece a nota do próprio guia: os parâmetros têm de bater com os do UE.

`[FREE5GC]` **Ubuntu 18.04.** O guia supõe Ubuntu 20.04 ou mais recente. No 18.04, com kernel 5.4, o projeto precisou de quatro adaptações:

- MongoDB 5.0, do repositório `bionic`, porque o 8.2 do guia não tem pacote para 18.04;
- Node.js 22 da compilação para glibc 2.17, de unofficial-builds.nodejs.org, porque a oficial exige glibc 2.28;
- a *webconsole* pela alternativa manual do guia (D.2), porque o `make webconsole` executa `sudo corepack enable`;
- gcc 9 e CMake 3.17 ou mais recente para o UERANSIM.

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

O `deploy.sh` (seção 3.4) instala o `free5gc/config/aanfcfg.yaml` do repositório quando a máquina do free5GC ainda não tem um, e depois não o sobrescreve. Ajuste esse arquivo antes do *deploy*, ou ajuste depois a cópia instalada, `~/free5gc/config/aanfcfg.yaml`, na máquina do free5GC. O arquivo do repositório tem os endereços do laboratório. Os campos:

- `sbi.registerIPv4` e `sbi.bindingIPv4`: um endereço da máquina do free5GC que a máquina do AF alcance. Não use um *loopback* `127.0.0.x`, a menos que o AF rode na mesma máquina. Porta 8010: a 8000 é a das NFs do free5GC, e o NEF passa a usá-la no endereço da máquina (parte 9).
- `nrfUri` e `nrfCertPem`: os mesmos das outras NFs (`http://127.0.0.10:8000`, `cert/nrf.pem`).
- `plmnSupportList`: a PLMN da sua rede.
- `routingIndicators`: os RIDs dos assinantes, **entre aspas** (`"0000"`). Com a lista vazia, o AAnF atende qualquer RID.
- `akma.authorizedAfs`: cada AF que pode pedir K_AF, com o AF_ID (`<FQDN>.<identificador Ua* em 10 hexadecimais>`) e o endereço de origem do AF.
- `nfInstanceId`: mantenha um UUID fixo. Com OAuth2 ligado no NRF, cada ID novo faz o NRF criar um certificado `cert/aanf_<ID>.pem`.

### 3.3 Corrigir o NRF

O NRF do free5GC v4.2.3 tem três falhas em relação a TS 29.510 que afetam o AAnF (`docs/free5gc-side/nrf.md`):

- descarta o `aanfInfoList` do perfil;
- recusa a descoberta do tipo `AANF`;
- emite tokens sem `aud`.

O AAnF exige `aud` no token. **Sem esta correção, o registro de chaves responde 401.** Na máquina do free5GC, com o free5GC parado e o repositório em `~/akma-free5gc` (seção 1.1):

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh nrf apply
```

O *script*:

1. aplica os três *patches* de `free5gc/patches/nrf/` em `~/free5gc/NFs/nrf`;
2. guarda o binário original em `bin/nrf.orig`;
3. recompila com `make nrf`.

`patch-nf.sh nrf revert` desfaz tudo. Depois de corrigir, confirme que o free5GC continua funcionando como antes: `./run.sh` e o registro do UE do teste de instalação.

### 3.4 Instalar e iniciar

Copie `free5gc/build/bin/aanf` para `~/free5gc/bin/aanf` e a configuração para `~/free5gc/config/aanfcfg.yaml`. Se o repositório está em outra máquina, o *script* do projeto faz isso por SSH e não sobrescreve uma configuração já instalada:

```bash
free5gc/scripts/deploy.sh usuario@maquina-do-free5gc
```

Se o *deploy* instalou a configuração do repositório sem o ajuste da seção 3.2, ajuste-a agora na máquina do free5GC. Com o free5GC em execução (pelo menos o NRF), inicie o AAnF na máquina do free5GC:

```bash
~/akma-free5gc/free5gc/scripts/run-aanf.sh start
```

Ou, a partir da outra máquina:

```bash
ssh usuario@maquina-do-free5gc bash -s -- start < free5gc/scripts/run-aanf.sh
```

O log fica em `~/free5gc/log/aanf/aanf.log`. As linhas esperadas são `OAuth2 setting receive from NRF: true` e `Start SBI server (listen on …:8010)`. `run-aanf.sh stop` encerra o AAnF, que antes se descadastra do NRF.

### 3.5 Verificar o registro e a descoberta no NRF

```bash
mongosh --quiet free5gc --eval 'printjson(db.NfProfile.find({nfType: "AANF"}, {_id: 0}).toArray())'
```

O perfil deve ter `nfType: 'AANF'`, o serviço `naanf-akma` (versão `v1`) no endereço configurado e o `aanfInfoList` com os RIDs. Para conferir a descoberta, rode na máquina do free5GC o `seed`, que o `deploy.sh` deixou em `~/akma-tools/`, fora do `PATH`:

```bash
~/akma-tools/seed discover --nrf http://127.0.0.10:8000 --routing-indicator 0000
```

O AAnF deve aparecer na lista. Com um RID que ele não atende, a resposta é `null`. É assim que o AUSF vai encontrá-lo na parte 4 (TS 33.535 §6.7).

### 3.6 Testar as três operações

Os valores abaixo são do vetor `testbed-ue` (`shared/akma/testdata/vectors.json`). O `seed` precisa de um *token* OAuth2 para registrar, e o NRF só é alcançável na máquina do free5GC. Para obtê-lo, o `seed` se registra no NRF como um AF temporário e se descadastra ao sair.

`[MVP]` O `--aanf` dá o endereço do AAnF, o `sbi.registerIPv4` da seção 3.2 com a porta 8010. Sem ele, o `seed` usa o do laboratório (`http://192.168.16.130:8010`).

Na máquina do free5GC, no papel do AUSF:

```bash
~/akma-tools/seed register --nrf http://127.0.0.10:8000 --aanf http://192.168.16.130:8010 --supi imsi-208930000000001 --k-ausf 141ef548f86e3a2a1ad949331caa9e6bb12401afa321bde4f9f23b84870dd654 --rid 0000 --mcc 208 --mnc 93
```

A saída deve trazer o A-KID `rid0000.atidiz1FN8DkMGyxwtfTVYsKUsWnzOJdBOLlQEiQlK2gHNA=@5gc.mnc093.mcc208.3gppnetwork.org`.

Na máquina do AF (o endereço dela precisa estar em `authorizedAfs`), no papel do AF, com o `seed` da seção 3.1:

```bash
free5gc/build/bin/seed retrieve --aanf http://192.168.16.130:8010 --af-id af1.akma.example.FF00000001 --akid 'rid0000.atidiz1FN8DkMGyxwtfTVYsKUsWnzOJdBOLlQEiQlK2gHNA=@5gc.mnc093.mcc208.3gppnetwork.org'
```

A K_AF deve ser `00c16845c2e4e1256755b84bb06d94e2e659d8ba53c59e874fc149bffd5ea04b`, a mesma que o oráculo da parte 2 calcula. O `expiry` deve estar uma hora à frente, em UTC.

Respostas que indicam configuração errada:

| Resposta | Causa provável |
|---|---|
| `401` no `register` | sem `--nrf`; o NRF sem a correção da seção 3.3 (token sem `aud`); ou o `nrfCertPem` do AAnF não é o certificado do NRF em uso |
| `403` sem causa no `retrieve` | o AF_ID não está em `authorizedAfs`, ou o pedido saiu de um endereço que não está na lista daquele AF_ID |
| `403 (K_AKMA_NOT_PRESENT)` | o A-KID não foi registrado, foi substituído por um registro mais novo do mesmo SUPI ou o AAnF reiniciou (o contexto fica só em memória) |
| `connection refused` | `--aanf` ausente ou com o endereço de outra máquina |

Por fim, para remover o contexto, na máquina do free5GC:

```bash
~/akma-tools/seed remove --nrf http://127.0.0.10:8000 --aanf http://192.168.16.130:8010 --supi imsi-208930000000001
```

---

## 4. Habilitar AKMA no UDM e no AUSF; provisionar um assinante

Depois de uma autenticação primária bem-sucedida de um assinante com assinatura AKMA:

- o UDM diz ao AUSF que o UE usa AKMA e qual é o seu Routing Indicator;
- o AUSF deriva K_AKMA e o A-KID da K_AUSF e os entrega ao AAnF da parte 3 (TS 33.535 §6.1).

Esta parte altera o UDM e o AUSF do free5GC com *patches* do repositório (`docs/free5gc-side/udm.md`, `ausf.md`).

### 4.1 Onde fica o repositório

O AUSF corrigido depende de código deste repositório: da biblioteca `shared/akma` e da cópia do `openapi` da parte 3. Os `go.mod` corrigidos apontam para `~/akma-free5gc` por caminho relativo. **O repositório tem de estar nesse local, ao lado de `~/free5gc`** (seção 1.1). Na máquina do free5GC, prepare as duas cópias do `openapi`: a das NFs e a do banco de testes do free5GC.

```bash
~/akma-free5gc/free5gc/scripts/prepare-openapi.sh
```

```bash
~/akma-free5gc/free5gc/scripts/prepare-openapi.sh ~/akma-free5gc/free5gc/build/openapi-test v1.2.5-0.20260527003827-02dc71b4d94f
```

### 4.2 Aplicar os *patches*

Com o free5GC parado, aplique os *patches* do UDM e do AUSF. O *script* guarda os binários originais em `bin/udm.orig` e `bin/ausf.orig` e recompila.

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh udm apply
```

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh ausf apply
```

Ligue o AKMA na configuração do AUSF. O *patch* acrescenta `akma.enable: true` ao `ausfcfg.yaml`:

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh config apply
```

Sem `akma.aanfUri`, o AUSF encontra o AAnF no NRF pelo Routing Indicator, o que exige o NRF corrigido da seção 3.3. Para usar um endereço fixo, descomente `aanfUri` no `ausfcfg.yaml`.

Para usar os testes AKMA do banco de testes (seção 4.4), aplique também:

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh test apply
```

Qualquer um desses passos se desfaz com `revert` no lugar de `apply`.

### 4.3 Provisionar um assinante

A assinatura AKMA são dois atributos do assinante no UDR (TS 29.505): `akmaAllowed` e `routingId`. O `routingId` é o Routing Indicator do SUCI do UE, como cadeia: `"0000"` no UERANSIM do testbed. A *webconsole* não tem esses campos e os apaga quando o assinante é editado. Provisione pelo *script* (`docs/free5gc-side/assinatura-akma.md`):

```bash
~/akma-free5gc/free5gc/scripts/provision-akma-subscriber.sh imsi-208930000000001 0000
```

### 4.4 Verificar

**Pelo banco de testes do free5GC**, com o free5GC parado e sem `sudo`:

```bash
cd ~/free5gc/test && go test -v -count=1 -run TestAkma .
```

Três testes devem passar:

- `TestAkma5GAKA` e `TestAkmaEAPAKAPrime`: o UE de teste deriva a K_AF da própria K_AUSF, e ela é igual à que o AAnF entrega;
- `TestAkmaWithoutSubscription`: sem assinatura, nada é registrado.

**Pelo UERANSIM.** Inicie o free5GC (`./run.sh`), o AAnF (`~/akma-free5gc/free5gc/scripts/run-aanf.sh start`), o gNB e o UE. No log do free5GC deve aparecer `AKMA key material of imsi-208930000000001 registered at the AAnF`. O log do AAnF (`~/free5gc/log/aanf/aanf.log`) mostra o A-KID guardado, da forma `rid0000.atid<44 caracteres>@5gc.mnc093.mcc208.3gppnetwork.org`. A parte 5 põe tudo isso numa tela só.

Se o UE do UERANSIM ficar em `MM-DEREGISTERED/LIMITED-SERVICE` logo depois de iniciar, reinicie o UE. É uma corrida da seleção de célula do UERANSIM ("no SI"), sem relação com o AKMA.

Se o AAnF estiver parado, o UE se registra mesmo assim. O AUSF só registra a falha no log: o AKMA não interfere na autenticação primária.

---

## 5. Ver o AKMA na tela: free5GC e UERANSIM

Até aqui, o AKMA só aparecia no log de arquivo. Esta parte faz duas coisas:

- os eventos principais aparecem no terminal do `run.sh`, no formato das outras NFs, com o AAnF no mesmo terminal;
- o UE do UERANSIM passa a derivar K_AKMA e o A-KID e mostra o A-KID na tela.

Assim dá para acompanhar o mesmo UE nas duas telas. O catálogo das mensagens está em `docs/free5gc-side/logs-akma.md`, e o lado do UE em `docs/af-side/ue-ueransim.md`.

### 5.1 free5GC

Os *patches* do UDM e do AUSF da parte 4 já trazem as mensagens. Falta a opção `-aanf` do `run.sh`, com o free5GC parado:

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh run apply
```

Se o AAnF estiver rodando pelo `run-aanf.sh`, pare-o: os dois usariam a mesma porta.

```bash
~/akma-free5gc/free5gc/scripts/run-aanf.sh stop
```

Daqui em diante, inicie o free5GC assim:

```bash
cd ~/free5gc && ./run.sh -aanf
```

O AAnF aparece na mesma tela, com uma linha de partida que mostra o que ele atende:

```
[INFO][AANF][AKMA] Routing Indicators [0000]; K_AF lifetime 1h0m0s; AF_IDs authorized by local policy: [af1.akma.example.FF00000001]
```

### 5.2 UERANSIM

Na máquina do UERANSIM, aplique os *patches* do UE. O *script* guarda `build/nr-ue.orig` e recompila o `nr-ue` e o autoteste `nr-akma-selftest`; na primeira vez, leva alguns minutos.

```bash
~/akma-free5gc/af/scripts/patch-ueransim.sh apply
```

Os *patches* mudam duas coisas no `config/free5gc-ue.yaml`:

- ligam a opção `akma: true`;
- fazem o UE procurar o gNB em 10.100.0.1.

Se você usa outro arquivo de configuração, faça o mesmo nele.

O endereço 10.100.0.1 é o lado do host de um par `veth` que liga o host ao *namespace* de rede do UE. Com o UERANSIM e o free5GC na mesma máquina, o UE precisa desse *namespace* para que a sessão PDU funcione (parte 7; `docs/af-side/ue-ueransim.md` §5). Faça o gNB escutar nesse endereço e crie o *namespace*, uma vez por inicialização da máquina:

```bash
sed -i 's/^linkIp: 127.0.0.1 /linkIp: 10.100.0.1 /' ~/UERANSIM/config/free5gc-gnb.yaml
```

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh up
```

Daqui em diante, o UE roda no *namespace*:

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ue
```

### 5.3 O que procurar

Inicie o gNB como na parte 4 e o UE com `ue-netns.sh ue`. Na tela do UE, logo depois do `Security Mode Command received`:

```
[nas] [info] AKMA key material derived, A-KID[rid0000.atid<44 caracteres>@5gc.mnc093.mcc208.3gppnetwork.org]
```

Na tela do `run.sh`, as linhas com `[AKMA]`, nesta ordem:

1. `[UDM][AKMA] AKMA indication sent to the AUSF for imsi-…, Routing Indicator 0000`;
2. `[AUSF][AKMA] AKMA indication received from the UDM …`;
3. `[AUSF][AKMA] K_AKMA and A-KID derived for imsi-…: A-KID …`;
4. `[AUSF][AKMA] AAnF for Routing Indicator 0000: http://… (discovered at the NRF)`;
5. `[AANF][AKMA] AKMA context of imsi-… registered, A-KID …`;
6. `[AUSF][AKMA] AKMA key material of imsi-… registered at the AAnF`.

O A-KID tem de ser o mesmo, caractere a caractere, na tela do UE, na linha 3 e na linha 5. A cada nova autenticação, o UE e a rede trocam de A-KID juntos, e o AAnF acrescenta `(previous context replaced)`.

Quando um AF pedir a K_AF (parte 6), o AAnF mostra `K_AF of imsi-… provided to AF_ID … for A-KID …, expires …`. Um AF fora da política aparece em `WARN`.

Nenhuma chave aparece, em nenhuma das telas. No arquivo de log, as mesmas mensagens estão no formato `chave=valor`; para filtrá-las:

```bash
grep 'CAT="AKMA"' ~/free5gc/log/<data>/free5gc.log
```

Para voltar ao comportamento anterior: `patch-nf.sh run revert` e `patch-ueransim.sh revert`.

---

## 6. Escrever um AF que usa AKMA

Um AF que usa AKMA faz três coisas (TS 33.535 §6.2.1):

- recebe o A-KID do UE;
- pede ao AAnF a K_AF daquele A-KID para o próprio AF_ID;
- usa a K_AF para autenticar o UE.

O AF de teste do repositório faz isso com um Ua\* simplificado: o HTTP Digest de TS 24.109 K.2, sem TLS, com o identificador Ua\* `FF00000001` da faixa privada. O protocolo está em `docs/af-side/ua-star.md` e o AF em `docs/af-side/af.md`.

### 6.1 Autorizar o AF no AAnF

O AAnF só entrega K_AF a um AF_ID autorizado, a partir do endereço do AF. No `aanfcfg.yaml` da parte 3 já está:

```yaml
configuration:
  akma:
    authorizedAfs:
      - afId: af1.akma.example.FF00000001
        sourceAddresses:
          - 192.168.16.129
```

O AF_ID é o FQDN do AF seguido do identificador Ua\* (TS 33.535 A.4). Se mudar o FQDN ou a máquina do AF, mude aqui também.

### 6.2 Compilar e configurar o AF

Na máquina do AF (no laboratório, a VM do AF), com Go 1.26.2:

```bash
af/scripts/build-af.sh
```

A configuração fica em `af/config/afcfg.yaml`. Os campos que importam:

- `fqdn`: o nome pelo qual o UE chama o AF. Ele entra na K_AF, e o AF só aceita esse nome no cabeçalho `Host`;
- `listen`: o endereço em que o AF escuta. Os pedidos ao AAnF não saem dele: saem do endereço que o sistema escolhe pela rota até o AAnF, e é esse que tem de estar em `sourceAddresses` (seção 6.1). Numa máquina com uma interface só, como no laboratório, os dois coincidem;
- `aanfUri`: o endereço do AAnF, o `sbi.registerIPv4` da seção 3.2 com a porta 8010. No repositório está o do laboratório;
- `kafExpiry`: o que fazer quando a K_AF vence, `refetch` ou `reject`.

```bash
af/build/bin/af -c af/config/afcfg.yaml
```

Na tela, a partida mostra o AF_ID e o `realm`:

```
[INFO][AF][Main] AF_ID af1.akma.example.FF00000001; realm "3GPP-bootstrapping-akma@af1.akma.example"; ...
```

### 6.3 Testar sem UE

O UE do UERANSIM só entrega a K_AF na parte 7. Antes disso, faça o papel do UE com um vetor do projeto, em dois passos:

1. semeie no AAnF o contexto do vetor `testbed-ue`, como na seção 3.6 (`~/akma-tools/seed register`, com o `--aanf`);
2. rode o cliente Ua\* com o mesmo vetor:

   ```bash
   af/build/bin/uaclient -vector testbed-ue -n 3
   ```

O cliente pede a sessão, recebe o desafio, responde com o A-KID e a K_AF e confere o `rspauth` do AF:

```
session established (nc 00000001): AF authenticated by rspauth; K_AF valid until ...
```

Na tela do AF aparecem `K_AF for A-KID … obtained from the AAnF` e `UE with A-KID … authenticated`. Na do AAnF, `K_AF of imsi-… provided to AF_ID af1.akma.example.FF00000001`.

O que testar a seguir, e o que esperar:

| Caso | Resposta do AF |
|---|---|
| K_AF errada: `-akid` com o A-KID do vetor e `-kaf` com outro valor (com `-vector`, o `-kaf` não vale) | `401` |
| A-KID que o AAnF não conhece (`-vector rid-012`) | `403 AKMA_KEY_NOT_AVAILABLE` |
| outro nome no `Host` (`-vector testbed-ue -fqdn af2.akma.example`) | `421 MISDIRECTED_REQUEST` |
| AF_ID fora da política do AAnF: um segundo AF, com outra configuração, `fqdn: af2.akma.example` e outra porta em `listen`; o cliente com `-vector testbed-ue -fqdn af2.akma.example -connect <listen desse AF>` | `503 AANF_UNAVAILABLE` |

Os resultados do laboratório estão em `docs/af-side/teste-af.md`.

**Antes da parte 7.** O vetor `testbed-ue` tem o SUPI do UE do UERANSIM (`imsi-208930000000001`). Se o UE já se registrou, o `seed register` deste teste substituiu o contexto dele no AAnF. Faça então uma nova autenticação do UE (seção 7.4) antes da seção 7.3. Sem ela, o AF responde `403 AKMA_KEY_NOT_AVAILABLE` ao UE, porque o AAnF não tem mais o A-KID que o UE apresenta.

### 6.4 Limitações

Sem TLS, o A-KID e o conteúdo trafegam em claro, e o AF só é autenticado pelo `rspauth` (`docs/af-side/ua-star.md` §10). É um protocolo de laboratório; o perfil normativo, com TLS, fica para depois do MVP.

---

## 7. Preparar o UE (UERANSIM) para o Ua\*

Na parte 5, o UE passou a derivar K_AKMA e o A-KID. Agora ele deriva também a K_AF de um AF e a entrega a uma aplicação, que vai ao AF pela sessão PDU (TS 33.535 §6.2.1 passo 1). O lado UE está em `docs/af-side/ue-ueransim.md` e o cliente em `docs/af-side/ua-client.md`.

### 7.1 Abrir o plano de usuário

Com o UE no *namespace* da parte 5, falta o NAT do N6 na máquina do free5GC. Ele faz a resposta do AF, ou da Internet, voltar ao UE (`docs/free5gc-side/plano-de-usuario.md`). Não é persistente: repita depois de reiniciar a máquina.

```bash
sudo ~/akma-free5gc/free5gc/scripts/n6-nat.sh apply
```

Com o free5GC, o gNB e o UE no ar, o UE deve alcançar o AF e a Internet pelo `uesimtun0`:

```bash
sudo ip netns exec ueransim-ue ping -I uesimtun0 -c 3 192.168.16.129
```

### 7.2 Conferir as derivações do UE

O autoteste compara as derivações do UE com os vetores do projeto, os mesmos da parte 2:

```bash
~/UERANSIM/build/nr-akma-selftest ~/akma-free5gc/shared/akma/testdata/vectors.json
```

O resultado esperado é `5 passed, 3 refused (FQDN not ASCII), 0 failed`. Os três recusados têm FQDN não ASCII, que o UE não deriva, porque não aplica NFKC.

### 7.3 Autenticar no AF a partir do UE

Na máquina do free5GC, compile o cliente Ua\* com `~/akma-free5gc/af/scripts/build-af.sh`, ou copie o `af/build/bin/uaclient` da máquina do AF para o mesmo lugar. O AAnF precisa ter o contexto da última autenticação do UE: se a seção 6.3 semeou o vetor `testbed-ue` depois que o UE se registrou, faça antes a nova autenticação da seção 7.4. Com o AF da parte 6 no ar:

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -n 3
```

O cliente roda no *namespace* do UE, com o seu usuário:

- pede ao UE, pelo `nr-cli`, o A-KID e a K_AF de `af1.akma.example`;
- vai ao AF a partir do endereço da sessão PDU;
- confere o `rspauth`.

```
ME: K_AF derived for AF_ID af1.akma.example.FF00000001
AF af1.akma.example at 192.168.16.129:8080 from 10.60.0.1, A-KID rid0000.atid…
session established (nc 00000001): AF authenticated by rspauth; K_AF valid until …
```

O mesmo A-KID aparece nas telas do UE, do free5GC (AUSF e AAnF) e do AF. Na tela do UE aparece ainda `[app] [info] AKMA K_AF derived for AF_ID[af1.akma.example.FF00000001], …`.

### 7.4 Nova autenticação

Force uma nova autenticação primária:

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh cli imsi-208930000000001 -e "deregister normal"
```

O UE volta a se registrar sozinho, com autenticação nova. O que muda:

- o UE e a rede derivam um A-KID novo;
- o AAnF registra `(previous context replaced)`;
- o cliente Ua\* passa a usar o A-KID novo;
- um pedido de K_AF com o A-KID antigo recebe `403 K_AKMA_NOT_PRESENT` do AAnF.

A sequência completa e os resultados do laboratório estão em `docs/af-side/teste-ue.md`.

### 7.5 Limitações

O `nr-cli` é uma interface de teste: entrega a K_AF de qualquer AF a qualquer processo local da máquina. A K_AKMA nunca sai do UE (TS 33.535 §4.4.3).

---

## 8. Executar o fluxo de ponta a ponta

Com as partes 3 a 7 feitas, um único roteiro executa o AKMA inteiro e os testes negativos, nas duas máquinas, e escreve um relatório (`docs/af-side/e2e.md`).

### 8.1 Preparar

Na máquina do AF:

```bash
af/scripts/build-af.sh
```

```bash
free5gc/scripts/build-aanf.sh
```

A máquina do AF precisa de acesso SSH por chave à do free5GC (seção 1.1) e do `af/config/afcfg.yaml` da seção 6.2. Na do free5GC, o roteiro espera o que as partes anteriores deixaram:

- o repositório em `~/akma-free5gc`;
- os *patches* aplicados, o AAnF implantado e o `seed` em `~/akma-tools/`;
- o UERANSIM com os *patches* e o `linkIp: 10.100.0.1`;
- o `uaclient` em `~/akma-free5gc/af/build/bin/`.

O free5GC deve estar parado: o roteiro o inicia. O AF da parte 6 também deve estar parado, porque o roteiro inicia o seu na mesma porta.

### 8.2 Rodar

Na máquina do AF, com o usuário e o endereço da máquina do free5GC:

```bash
af/scripts/e2e.sh woca@192.168.16.130
```

Sem o argumento, o roteiro usa a máquina do free5GC do laboratório, `woca@192.168.16.130`. O endereço do AAnF vem do `aanfUri` de `af/config/afcfg.yaml`. Antes de começar, o roteiro confere que é o `sbi.bindingIPv4` do `aanfcfg.yaml` instalado na máquina do free5GC.

O roteiro pede a senha de `sudo` da máquina do free5GC uma vez e leva cerca de dois minutos. Ele:

- aplica o NAT do N6 e cria o *namespace* do UE;
- sobe o free5GC com o AAnF, o gNB e o UE;
- autentica no AF a partir do UE;
- força novas autenticações primárias: com o AAnF parado, sem assinatura AKMA e com a assinatura restaurada;
- roda os testes negativos;
- encerra tudo.

Cada verificação aparece como `PASS` ou `FAIL`, e o fim mostra o total:

```
report: …/af/build/e2e/<data-hora>/report.md
63 passed, 0 failed
```

### 8.3 Ler o relatório

O `report.md` traz o A-KID de cada autenticação, a tabela de verificações e as linhas AKMA do free5GC, do UE e do AF. Ali dá para seguir a pergunta-guia do projeto de ponta a ponta:

1. **Autenticação primária e registro.** O UE se autentica; o UDM indica AKMA; o AUSF deriva o A-KID e o registra no AAnF.
2. **A-KID igual.** O UE derivou o mesmo A-KID.
3. **Ua\*.** O AF pede a K_AF ao AAnF e autentica o UE com ela, e o UE confere o `rspauth` do AF.

Os testes negativos e o que cada um prova estão em `docs/af-side/e2e.md` §4. A rastreabilidade de cada etapa de TS 33.535 §6.1 e §6.2.1 até a verificação correspondente está na matriz, revisão 4.


---

## 9. AF fora da rede do operador, acesso anônimo e UERANSIM em outra máquina

Até aqui o AF ficou dentro da rede do operador e pediu a K_AF direto ao AAnF. Esta parte mostra três coisas:

- um AF de fora da rede pedindo a K_AF ao NEF (TS 33.535 §6.3);
- o acesso anônimo, em que o AF recebe a K_AF sem o SUPI (§6.2.2);
- o UERANSIM numa máquina separada, como no guia oficial do free5GC.

A especificação do NEF está em `docs/free5gc-side/nef.md`, e o desvio de segurança entre o NEF e o AF, em `docs/conformidade/desvios.md` (DV-9).

### 9.1 Preparar o NEF

O NEF do free5GC não tem a API de AKMA. Com o free5GC parado, na máquina do free5GC:

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh nef apply
```

O *patch* acrescenta ao NEF a operação `POST /3gpp-akma/v1/retrieve` (TS 29.522 §5.14) e recompila o NEF; o binário original fica em `bin/nef.orig`.

A configuração do NEF muda pelo *patch* de configuração da parte 4 (`free5gc/patches/config/0002-nef-akma.patch`). Se você aplicou a configuração antes desta versão do tutorial, desfaça-a com a cópia antiga dos *patches* (`patch-nf.sh config revert`), atualize o repositório e aplique de novo. No `~/free5gc/config/nefcfg.yaml`, confira:

- `sbi.registerIPv4` e `sbi.bindingIPv4`: o endereço da máquina do free5GC, que o AF externo alcança. A porta continua a 8000;
- `3gpp-akma` em `serviceList`;
- `akma.authorizedAfs`: cada AF_ID externo, com o endereço de origem do AF (§6.3 passo 2).

O AAnF precisa aceitar os pedidos que chegam pelo NEF. No `~/free5gc/config/aanfcfg.yaml`, a lista `akma.authorizedNefs` dá o endereço do NEF e os AF_ID externos aceitos dele (`docs/free5gc-side/aanf.md` §2). O AAnF escuta na porta 8010.

Suba o free5GC com `./run.sh -aanf`. Na tela do NEF deve aparecer `Start SBI server (listen on 192.168.16.130:8000)`; na do AAnF, `AF_IDs accepted from a NEF: [af1.akma.example.FF00000001]`.

### 9.2 Rodar o AF fora da rede do operador

O AF externo se registra no NRF para obter o *token* da API do NEF. No laboratório o NRF escuta num endereço de *loopback* da máquina do free5GC, e o AF o alcança por um túnel SSH (`nef.md` §2.1). Na máquina do AF, deixe o túnel aberto num terminal:

```bash
ssh -N -L 127.0.0.10:8000:127.0.0.10:8000 woca@192.168.16.130
```

A configuração do AF externo está em `af/config/afcfg-external.yaml`. Ela tem `nefUri` no lugar de `aanfUri`, além de `nrfUri` e da porta 8083. Em outro terminal:

```bash
af/build/bin/af -c af/config/afcfg-external.yaml
```

Na tela do AF:

```
[INFO][AF][Main] AF outside the operator network: K_AF through the NEF at http://192.168.16.130:8000; registered at the NRF http://127.0.0.10:8000 as AF … (OAuth2 true)
```

Com o UE no ar (parte 7), autentique o UE no AF externo:

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -connect 192.168.16.129:8083 -n 3
```

A saída do cliente é a mesma da seção 7.3. Na tela do free5GC aparecem as linhas do NEF e do AAnF:

```
[INFO][NEF][AKMA] K_AF request from AF_ID af1.akma.example.FF00000001 for A-KID rid0000.atid…
[INFO][AANF][AKMA] K_AF of imsi-208930000000001 provided to AF_ID af1.akma.example.FF00000001 for A-KID …, through the NEF at 192.168.16.130
[INFO][NEF][AKMA] K_AF for A-KID rid0000.atid… relayed to AF_ID af1.akma.example.FF00000001, expires …
```

As linhas do NEF não mostram o SUPI. A resposta do NEF ao AF também não o traz (§6.3 passo 5).

Ao encerrar o AF com Ctrl-C, ele se descadastra do NRF.

### 9.3 Acesso anônimo

Com `anonymous: true` na configuração, o AF pede a K_AF com `anonInd` e recebe só a K_AF e a validade, sem o SUPI (§6.2.2). Vale para o AF interno e para o externo. O A-KID serve ao AF como identificador temporário do usuário.

Na tela do NEF, o pedido anônimo aparece como `K_AF request … , anonymous`. O AAnF continua sabendo de quem é o contexto: o anonimato é em relação ao AF.

### 9.4 Remover o contexto do AAnF

Quem pede a remoção do contexto, pela norma, é o OAM (TS 29.535 §4.2.1.3.2 e tab. 4.2.2.1-1). No laboratório, o `seed` faz esse papel, na máquina do free5GC:

```bash
~/akma-tools/seed remove --nrf http://127.0.0.10:8000 --aanf http://192.168.16.130:8010 --supi imsi-208930000000001
```

A resposta é `removed`. Depois disso, o AF recebe `403 AKMA_KEY_NOT_AVAILABLE` até uma nova autenticação primária do UE (seção 7.4).

### 9.5 UERANSIM em outra máquina

O guia oficial do free5GC instala o UERANSIM numa máquina separada. No laboratório: free5GC em 192.168.16.130, AF em 192.168.16.129 e UERANSIM em 192.168.16.131.

**Na máquina do free5GC**, o N2 e o N3 saem do *loopback* para o endereço da máquina, como o guia oficial manda:

| Arquivo | Campo | Valor |
|---|---|---|
| `config/amfcfg.yaml` | `ngapIpList` | `192.168.16.130` |
| `config/smfcfg.yaml` | `endpoints` da interface N3 do UPF | `192.168.16.130` |
| `config/upfcfg.yaml` | `gtpu.ifList[].addr` | `192.168.16.130` |

Se a máquina do free5GC também roda um gNB, ele precisa de dois ajustes no `free5gc-gnb.yaml`:

- `ngapIp`: o endereço da máquina. Um gNB em *loopback* não alcança o AMF em 192.168.16.130 (`Network is unreachable`);
- `gtpIp`: `10.100.0.1`, o endereço da veth do *namespace*. Não pode ser `127.0.0.1`, pelo mesmo motivo. Também não pode ser o endereço da máquina, porque a porta 2152 desse endereço é do UPF.

**Na máquina do UERANSIM**, siga as partes 5 e 7 com o repositório em `~/akma-free5gc`: *patches* do UERANSIM, *namespace* do UE e `uaclient`. No `config/free5gc-gnb.yaml`:

| Campo | Valor |
|---|---|
| `linkIp` | `10.100.0.1` |
| `ngapIp`, `gtpIp` | o endereço da máquina do UERANSIM, `192.168.16.131` |
| `amfConfigs[].address` | o endereço da máquina do free5GC, `192.168.16.130` |

O NAT do N6 continua na máquina do free5GC: a máquina do UERANSIM não precisa de NAT. Na máquina do UERANSIM, um terminal por comando:

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh up
```

```bash
~/UERANSIM/build/nr-gnb -c ~/UERANSIM/config/free5gc-gnb.yaml
```

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ue
```

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -n 3
```

O gNB mostra `SCTP connection established (192.168.16.130:38412)` e `NG Setup procedure is successful`. O UE mostra o A-KID, que é o mesmo nas linhas do AUSF e do AAnF na tela do free5GC. Para o AF externo, use `ua -connect 192.168.16.129:8083`.

Resultado do laboratório em 19/10/2026: com o UE na máquina 192.168.16.131, o Ua\* passou com o AF interno, com o AF externo via NEF e com o AF externo anônimo (`docs/af-side/teste-nef.md`).
