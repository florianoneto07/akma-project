# Tutorial — AKMA sobre o free5GC

Este tutorial descreve o que fazer **depois de instalar o free5GC** pelo guia oficial para usar AKMA (3GPP TS 33.535). Ele cresce a cada MVP; esta versão corresponde ao MVP-9 (CAPIF), que veio depois do MVP-12. O MVP-5 foi o primeiro MVP funcional. Em 19/10/2026 ele foi seguido do início numa VM com o free5GC v4.2.3 recém-instalado, e o roteiro da parte 8 terminou com 63 verificações e nenhuma falha (critério 4 do MVP-5).

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
| 10 | Renovar a K_AKMA: o AAnF pede uma nova autenticação | sim |
| 11 | CAPIF: o AF obtém o *token* do CCF | sim |

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
[INFO][AANF][AKMA] Routing Indicators [0000]; K_AF lifetime 1h0m0s; AF_IDs authorized by local policy: [af1.akma.example.0100000007 af1.akma.example.0100011301 …]
```

### 5.2 UERANSIM

Na máquina do UERANSIM, aplique os *patches* do UE. O *script* guarda `build/nr-ue.orig` e recompila o `nr-ue` e o autoteste `nr-akma-selftest`; na primeira vez, leva alguns minutos.

```bash
~/akma-free5gc/af/scripts/patch-ueransim.sh apply
```

Os *patches* mudam três coisas no `config/free5gc-ue.yaml`:

- ligam a opção `akma: true`;
- fazem o UE procurar o gNB em 10.100.0.1;
- desde o MVP-8, põem a política do serviço de chaves do ME: cada aplicação só recebe a K_AF dos seus AF_IDs (TS 33.535 §4.4.3; `docs/af-side/ue-ueransim.md` §3). O exemplo autoriza o `uaclient` do usuário de *uid* 1000 para os AF_IDs dos perfis do AF da parte 6: um por perfil e por *cipher suite*, e o de teste, `af1.akma.example.FF00000001`.

Se o usuário que roda o cliente Ua\* não tem *uid* 1000 (`id -u`), troque o `uid`. Se você usa outro arquivo de configuração, faça o mesmo nele.

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

O AF de teste do repositório serve os perfis de TS 24.109 Anexo K e de TS 33.535 Anexos C e D, cada um na sua porta (`docs/af-side/ua-star-norma.md`):

| Perfil | Porta no laboratório | Como o UE se autentica | AF_ID |
|---|---|---|---|
| `digest-tls` (K.2) | 8443 | HTTP Digest com a K_AF dentro de TLS; o AF se autentica pelo certificado | `<FQDN>.010001<cipher suite>` |
| `psk-tls13` (K.3.3) | 8444 | TLS 1.3 com a K_AF como chave pré-compartilhada | `<FQDN>.0100000007` |
| `psk-tls12` (K.3.2) | 8445 | TLS 1.2 `ECDHE-PSK` com a K_AF | `<FQDN>.010001D001` |
| `psk-dtls13` (TS 33.535 Anexo C) | UDP 5684 | DTLS 1.3 com a K_AF como chave pré-compartilhada; CoAP dentro | `<FQDN>.0100000009` |
| `psk-dtls12` | UDP 5685 | DTLS 1.2 `ECDHE-PSK` com a K_AF; CoAP dentro | `<FQDN>.010001D001` |
| `oscore` (TS 33.535 Anexo D) | UDP 5683 | CoAP protegido pelo OSCORE, com o contexto derivado da K_AF | `<FQDN>.0100000008` |
| `digest-http` (teste, desligado por padrão) | 8080 | o Ua\* simplificado: o Digest, sem TLS (`docs/af-side/ua-star.md`) | `<FQDN>.FF00000001` |

O AF está em `docs/af-side/af.md`.

### 6.1 Autorizar o AF no AAnF

O AAnF só entrega K_AF a um AF_ID autorizado, a partir do endereço do AF. Cada perfil, e no K.2 cada *cipher suite*, tem o seu AF_ID (TS 33.220 H.3). No `aanfcfg.yaml` da parte 3 já estão os dez AF_IDs de `af1.akma.example`:

```yaml
configuration:
  akma:
    authorizedAfs:
      - afId: af1.akma.example.0100011301   # K.2 em TLS 1.3, TLS_AES_128_GCM_SHA256
        sourceAddresses:
          - 192.168.16.129
      # ... 0100011302, 0100011303, 010001C02B, 010001C02C (K.2),
      #     0100000007 (PSK TLS 1.3), 010001D001 (PSK TLS e DTLS 1.2),
      #     0100000009 (PSK DTLS 1.3), 0100000008 (OSCORE), FF00000001 (teste)
```

O AF_ID é o FQDN do AF seguido do identificador Ua\* (TS 33.535 A.4). Se mudar o FQDN ou a máquina do AF, mude em todas as entradas.

### 6.2 Compilar e configurar o AF

Na máquina do AF (no laboratório, a VM do AF), com Go 1.26.2, `gcc`, `make` e `unzip`. Primeiro o wolfSSL, que faz os perfis PSK: o `crypto/tls` do Go não tem as suítes PSK de TS 33.210 nem DTLS. O *script* baixa o wolfSSL 5.9.4, confere o SHA-256 e o compila em `af/build/wolfssl`:

```bash
af/scripts/build-wolfssl.sh
```

Depois a CA de teste e o certificado do AF para o K.2, em `af/build/cert`:

```bash
af/scripts/make-test-ca.sh
```

E o AF e o cliente:

```bash
af/scripts/build-af.sh
```

Com o wolfSSL, os binários ficam sob GPLv3, como o próprio wolfSSL. Sem ele (`AF_WOLFSSL=no`), o AF não tem os perfis PSK e recusa partir com eles na configuração.

A configuração fica em `af/config/afcfg.yaml`. Os campos que importam:

- `fqdn`: o nome pelo qual o UE chama o AF. Ele entra na K_AF, e o AF só aceita esse nome no cabeçalho `Host` e no `server_name` do TLS;
- `profiles`: um servidor por perfil, com `type` e `listen`; o `digest-tls` tem ainda `certificate` e `key`, relativos ao arquivo. Os pedidos ao AAnF não saem dos `listen`: saem do endereço que o sistema escolhe pela rota até o AAnF, e é esse que tem de estar em `sourceAddresses` (seção 6.1). Numa máquina com uma interface só, como no laboratório, os dois coincidem;
- `aanfUri`: o endereço do AAnF, o `sbi.registerIPv4` da seção 3.2 com a porta 8010. No repositório está o do laboratório;
- `kafExpiry`: o que fazer quando a K_AF vence, `refetch` ou `reject`. Nos perfis PSK, o AF também encerra a sessão TLS quando a K_AF vence.

```bash
af/build/bin/af -c af/config/afcfg.yaml
```

Na tela, a partida mostra cada perfil:

```
[INFO][AF][Main] Built with wolfSSL: PSK profiles available
[INFO][AF][Main] Ua* server digest-tls listening on 192.168.16.129:8443 for https://af1.akma.example/ua/v1/session, AF_ID af1.akma.example.010001<cipher suite> (TS 24.109 K.2)
[INFO][AF][Main] Ua* server psk-tls13 listening on 192.168.16.129:8444 for https://af1.akma.example/ua/v1/session, AF_ID af1.akma.example.0100000007 (TS 24.109 K.3.3)
[INFO][AF][Main] Ua* server psk-tls12 listening on 192.168.16.129:8445 for https://af1.akma.example/ua/v1/session, AF_ID af1.akma.example.010001<cipher suite> (TS 24.109 K.3.2)
```

### 6.3 Testar sem UE

O UE do UERANSIM só entrega a K_AF na parte 7, e ele a deriva para qualquer perfil. Antes disso, faça o papel do UE com um vetor do projeto. O vetor traz a K_AF de um AF_ID só, `af1.akma.example.FF00000001`, a do perfil de teste. Ligue-o, tirando o comentário destas linhas de `af/config/afcfg.yaml`, e reinicie o AF:

```yaml
    - type: digest-http
      listen: 192.168.16.129:8080
```

Então:

1. semeie no AAnF o contexto do vetor `testbed-ue`, como na seção 3.6 (`~/akma-tools/seed register`, com o `--aanf`);
2. rode o cliente Ua\* com o mesmo vetor:

   ```bash
   af/build/bin/uaclient -vector testbed-ue -n 3
   ```

O cliente pede a sessão, recebe o desafio, responde com o A-KID e a K_AF e confere o `rspauth` do AF:

```
session established (nc 00000001): AF authenticated by rspauth; A-KID …; AF_ID af1.akma.example.FF00000001; K_AF valid until ...
```

Na tela do AF aparecem `K_AF for A-KID … and AF_ID af1.akma.example.FF00000001 obtained from the AAnF` e `UE with A-KID … authenticated`. Na do AAnF, `K_AF of imsi-… provided to AF_ID af1.akma.example.FF00000001`.

O que testar a seguir, e o que esperar:

| Caso | Resposta do AF |
|---|---|
| K_AF errada: `-akid` com o A-KID do vetor e `-kaf` com outro valor (com `-vector`, o `-kaf` não vale) | `401` |
| A-KID que o AAnF não conhece (`-vector rid-012`) | `403 AKMA_KEY_NOT_AVAILABLE` |
| outro nome no `Host` (`-vector testbed-ue -fqdn af2.akma.example`) | `421 MISDIRECTED_REQUEST` |
| AF_ID fora da política do AAnF: um segundo AF, com outra configuração, `fqdn: af2.akma.example` e outras portas; o cliente com `-vector testbed-ue -fqdn af2.akma.example -connect <listen desse AF>` | `503 AANF_UNAVAILABLE` |
| certificado de outro nome no K.2: `-profile digest-tls -ca af/build/cert/ca.pem -vector testbed-ue -fqdn af2.akma.example` | o cliente encerra o TLS antes de mandar o A-KID |
| PSK com a K_AF de outro AF_ID: `-profile psk-tls13 -vector testbed-ue` | o *handshake* falha |

Os resultados do laboratório estão em `docs/af-side/teste-af.md` e, para os perfis TLS, no grupo T de `docs/af-side/e2e.md`. Desligue de novo o `digest-http` depois do teste.

**Antes da parte 7.** O vetor `testbed-ue` tem o SUPI do UE do UERANSIM (`imsi-208930000000001`). Se o UE já se registrou, o `seed register` deste teste substituiu o contexto dele no AAnF. Faça então uma nova autenticação do UE (seção 7.4) antes da seção 7.3. Sem ela, o AF recusa o UE, porque o AAnF não tem mais o A-KID que o UE apresenta.

### 6.4 Limitações

- **`digest-http`.** Sem TLS, o A-KID e o conteúdo trafegam em claro, e o AF só é autenticado pelo `rspauth` (`docs/af-side/ua-star.md` §10). É só para teste.
- **Certificado.** O do K.2 vem de uma CA de teste, sem revogação nem OCSP.
- **SBI.** O pedido do AF ao AAnF continua sem TLS (DV-1).

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

Na máquina do free5GC, ponha o cliente Ua\* e a CA de teste do AF onde o *script* os procura: copie da máquina do AF o `af/build/bin/uaclient` para `~/akma-free5gc/af/build/bin/` e o `af/build/cert/ca.pem` para `~/akma-free5gc/af/build/cert/`. O cliente compilado na máquina do AF é estático e roda na máquina do free5GC. O AAnF precisa ter o contexto da última autenticação do UE: se a seção 6.3 semeou o vetor `testbed-ue` depois que o UE se registrou, faça antes a nova autenticação da seção 7.4. Com o AF da parte 6 no ar:

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -profile psk-tls13 -n 2
```

O cliente roda no *namespace* do UE, com o seu usuário:

- pede o A-KID e a K_AF ao serviço de chaves do ME, em `/run/ueransim/akma-<SUPI>.sock`, para o AF_ID do perfil;
- vai ao AF a partir do endereço da sessão PDU;
- no TLS 1.3 com PSK, o *handshake* autentica os dois lados com a K_AF.

```
AF af1.akma.example at 192.168.16.129:8444 from 10.60.0.1, profile psk-tls13
ME: K_AF derived for AF_ID af1.akma.example.0100000007
session established: AF authenticated by the TLS handshake with K_AF; A-KID rid0000.atid…; AF_ID af1.akma.example.0100000007; TLS 1.3, cipher suite 1301; K_AF valid until …
```

Os outros perfis:

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -profile psk-tls12 -reconnect -n 2
```

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -profile digest-tls -ca ~/akma-free5gc/af/build/cert/ca.pem -n 3
```

E os perfis DTLS, com CoAP sobre UDP pela sessão PDU:

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -profile psk-dtls13 -reconnect -n 2
```

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -profile psk-dtls12 -reconnect -n 2
```

E o perfil OSCORE: um `POST /akma` e dois pedidos protegidos no mesmo contexto; com `-rekey`, o contexto é renovado antes de cada pedido:

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -profile oscore -n 2
```

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -profile oscore -rekey -n 2
```

No TLS 1.2, o ME deriva a K_AF depois que o AF escolheu a suíte (`AF_ID af1.akma.example.010001D001`), e com `-reconnect` a segunda conexão retoma a sessão (`resumed true`). No DTLS 1.3, o AF_ID é `af1.akma.example.0100000009`; o AF responde ao primeiro `ClientHello` com um *cookie* e só pede a K_AF ao AAnF depois que o UE o devolve. No OSCORE, o AF_ID é `af1.akma.example.0100000008`, e a saída do `-rekey` mostra o *ID Context* novo de cada renovação. No K.2, o cliente confere o certificado do AF contra `af1.akma.example` e a CA, e faz o Digest dentro do TLS; o AF_ID leva a suíte da conexão, como `af1.akma.example.0100011301`. As opções estão em `docs/af-side/ua-client.md`.

O mesmo A-KID aparece nas telas do UE, do free5GC (AUSF e AAnF) e do AF. Na tela do UE aparece ainda `K_AF for AF_ID[af1.akma.example.0100000007] given to application [uaclient] (uid 1000), …`.

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

Desde o MVP-8, o ME separa as aplicações pelo usuário; o executável da política é conferência de melhor esforço (`docs/af-side/ua-star-norma.md` §2.2). O `nr-cli akma-kaf`, que entrega a K_AF de qualquer AF a qualquer processo local, só funciona com `akmaCliKaf: true`, para testes. A K_AKMA nunca sai do UE (TS 33.535 §4.4.3).

---

## 8. Executar o fluxo de ponta a ponta

Com as partes 3 a 7 feitas, um único roteiro executa o AKMA inteiro e os testes negativos, nas duas máquinas, e escreve um relatório (`docs/af-side/e2e.md`).

### 8.1 Preparar

Na máquina do AF:

```bash
af/scripts/build-wolfssl.sh
```

```bash
af/scripts/make-test-ca.sh
```

```bash
af/scripts/build-af.sh
```

```bash
free5gc/scripts/build-aanf.sh
```

A máquina do AF precisa de acesso SSH por chave à do free5GC (seção 1.1), do `openssl` e do `af/config/afcfg.yaml` da seção 6.2. Na do free5GC, o roteiro espera o que as partes anteriores deixaram:

- o repositório em `~/akma-free5gc`;
- os *patches* aplicados, o AAnF implantado e o `seed` em `~/akma-tools/`;
- o UERANSIM com os *patches* e o `linkIp: 10.100.0.1`;
- o `uaclient` compilado com o wolfSSL em `~/akma-free5gc/af/build/bin/` (seção 7.3).

O roteiro sobe o seu AF com os perfis de `afcfg.yaml` e mais o `digest-http` na porta 8080, e copia a CA de teste para a máquina do free5GC. O free5GC deve estar parado: o roteiro o inicia. O AF da parte 6 também deve estar parado, porque o do roteiro usa as mesmas portas.

### 8.2 Rodar

Na máquina do AF, com o usuário e o endereço da máquina do free5GC:

```bash
af/scripts/e2e.sh woca@192.168.16.130
```

Sem o argumento, o roteiro usa a máquina do free5GC do laboratório, `woca@192.168.16.130`. O endereço do AAnF vem do `aanfUri` de `af/config/afcfg.yaml`. Antes de começar, o roteiro confere que é o `sbi.bindingIPv4` do `aanfcfg.yaml` instalado na máquina do free5GC.

O roteiro pede a senha de `sudo` da máquina do free5GC uma vez e leva cerca de quatro minutos. Ele:

- aplica o NAT do N6 e cria o *namespace* do UE;
- sobe o free5GC com o AAnF, o gNB e o UE;
- autentica no AF a partir do UE, por todos os perfis: o K.2 em TLS 1.3 e 1.2, o PSK em TLS 1.3 e 1.2, o PSK em DTLS 1.3 e 1.2, com CoAP, e o OSCORE, com retomada de sessão e renovação do contexto;
- confere o PSK TLS 1.3 com o `openssl s_client` no papel do UE, e a expiração da K_AF numa sessão TLS; com um OpenSSL 4.1 em `OPENSSL_DTLS13`, confere também o PSK DTLS 1.3, e com a libcoap em `LIBCOAP_CLIENT`, o OSCORE (`docs/af-side/e2e.md` §3);
- força novas autenticações primárias: com o AAnF parado, sem assinatura AKMA e com a assinatura restaurada;
- renova a K_AKMA pelo AAnF, com `kakmaRefresh: 20s`: o UDM notifica o AMF, que autentica de novo o UE registrado (parte 10);
- roda os testes negativos;
- encerra tudo.

Cada verificação aparece como `PASS` ou `FAIL`, e o fim mostra o total:

```
report: …/af/build/e2e/<data-hora>/report.md
161 passed, 0 failed, 0 skipped
```

Sem o `OPENSSL_DTLS13`, a verificação D03 aparece como pulada; sem o `LIBCOAP_CLIENT`, a O04 e a O05. Sem os dois: `158 passed, 0 failed, 3 skipped`.

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

O segundo *patch* da série traz o NEF do CAPIF (parte 11). Com a seção `capif` do `nefcfg.yaml`, que o *patch* de configuração da parte 4 acrescenta, o NEF lê na partida a chave pública dos *tokens* do CCF. Faça já as chaves do CCF, com o endereço da máquina do free5GC:

```bash
~/akma-free5gc/free5gc/scripts/make-capif-certs.sh ~/free5gc/cert 192.168.16.130
```

Sem o CCF no ar, o NEF tenta publicar a API a cada 5 s, com um aviso na primeira falha, e continua aceitando os *tokens* do NRF (`allowNrfTokens: true`).

A configuração do NEF muda pelo *patch* de configuração da parte 4 (`free5gc/patches/config/0002-nef-akma.patch`). Se você aplicou a configuração antes desta versão do tutorial, desfaça-a com a cópia antiga dos *patches* (`patch-nf.sh config revert`), atualize o repositório e aplique de novo. No `~/free5gc/config/nefcfg.yaml`, confira:

- `sbi.registerIPv4` e `sbi.bindingIPv4`: o endereço da máquina do free5GC, que o AF externo alcança. A porta continua a 8000;
- `3gpp-akma` em `serviceList`;
- `akma.authorizedAfs`: cada AF_ID externo, com o endereço de origem do AF (§6.3 passo 2): os mesmos dez AF_IDs de `af1.akma.example` do AAnF.

O AAnF precisa aceitar os pedidos que chegam pelo NEF. No `~/free5gc/config/aanfcfg.yaml`, a lista `akma.authorizedNefs` dá o endereço do NEF e os AF_ID externos aceitos dele (`docs/free5gc-side/aanf.md` §2). O AAnF escuta na porta 8010.

Suba o free5GC com `./run.sh -aanf`. Na tela do NEF deve aparecer `Start SBI server (listen on 192.168.16.130:8000)`; na do AAnF, `AF_IDs accepted from a NEF: [af1.akma.example.0100000007 af1.akma.example.0100011301 …]`.

### 9.2 Rodar o AF fora da rede do operador

O AF externo se registra no NRF para obter o *token* da API do NEF. No laboratório o NRF escuta num endereço de *loopback* da máquina do free5GC, e o AF o alcança por um túnel SSH (`nef.md` §2.1). Na máquina do AF, deixe o túnel aberto num terminal:

```bash
ssh -N -L 127.0.0.10:8000:127.0.0.10:8000 woca@192.168.16.130
```

A configuração do AF externo está em `af/config/afcfg-external.yaml`. Ela tem `nefUri` no lugar de `aanfUri`, além de `nrfUri`, e os mesmos perfis do AF interno, nas portas TCP 8453 a 8455 e UDP 5693 a 5695. Em outro terminal:

```bash
af/build/bin/af -c af/config/afcfg-external.yaml
```

Na tela do AF:

```
[INFO][AF][Main] AF outside the operator network: K_AF through the NEF at http://192.168.16.130:8000; registered at the NRF http://127.0.0.10:8000 as AF … (OAuth2 true)
```

Com o UE no ar (parte 7), autentique o UE no AF externo:

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -profile psk-tls13 -connect 192.168.16.129:8454 -n 1
```

A saída do cliente é a mesma da seção 7.3. Na tela do free5GC aparecem as linhas do NEF e do AAnF:

```
[INFO][NEF][AKMA] K_AF request from AF_ID af1.akma.example.0100000007 for A-KID rid0000.atid…
[INFO][AANF][AKMA] K_AF of imsi-208930000000001 provided to AF_ID af1.akma.example.0100000007 for A-KID …, through the NEF at 192.168.16.130
[INFO][NEF][AKMA] K_AF for A-KID rid0000.atid… relayed to AF_ID af1.akma.example.0100000007, expires …
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
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -profile psk-tls13 -n 2
```

O gNB mostra `SCTP connection established (192.168.16.130:38412)` e `NG Setup procedure is successful`. O UE mostra o A-KID, que é o mesmo nas linhas do AUSF e do AAnF na tela do free5GC. Para o AF externo, use `ua -profile psk-tls13 -connect 192.168.16.129:8454`.

Resultado do laboratório em 19/10/2026: com o UE na máquina 192.168.16.131, o Ua\* passou com o AF interno, com o AF externo via NEF e com o AF externo anônimo (`docs/af-side/teste-nef.md`).

---

## 10. Renovar a K_AKMA: o AAnF pede uma nova autenticação

A K_AKMA só muda com uma nova autenticação primária (TS 33.535 §6.4.1). Pela política local do operador, o AAnF pode pedir ao UDM essa autenticação (§6.4.4): o UDM avisa o AMF, o AMF autentica de novo o UE registrado, e o AUSF registra no AAnF a K_AKMA e o A-KID novos. O UE deriva os seus no `Security Mode Command` que o AMF manda em seguida. Nada muda no UE do UERANSIM.

O free5GC não tem essa cadeia. Esta parte aplica os *patches* do UDM, do AMF e do `openapi` (`docs/free5gc-side/renovacao-kakma.md`, `udm.md` §6, `amf.md`) e liga a política do AAnF (`aanf.md` §11).

### 10.1 Atualizar o `openapi` e o UDM, corrigir o AMF

Na máquina do free5GC, com o free5GC parado, prepare de novo as duas cópias do `openapi` (seção 4.1). O *script* aplica os *patches* até o `0003`:

```bash
~/akma-free5gc/free5gc/scripts/prepare-openapi.sh
```

```bash
~/akma-free5gc/free5gc/scripts/prepare-openapi.sh ~/akma-free5gc/free5gc/build/openapi-test v1.2.5-0.20260527003827-02dc71b4d94f
```

Se você seguiu a parte 4 com esta versão do repositório, o UDM já tem os dois *patches* (`patch-nf.sh udm status` diz `patches applied`). Se a seguiu com uma versão anterior, o UDM só tem o `0001`, e o `apply` recusa as mudanças locais. Desfaça com uma cópia do `0001` e aplique de novo:

```bash
mkdir -p ~/udm-0001 && cp ~/akma-free5gc/free5gc/patches/udm/0001-*.patch ~/udm-0001/
```

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh udm revert ~/udm-0001
```

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh udm apply
```

Depois, o AMF. O *script* guarda `bin/amf.orig` e recompila:

```bash
~/akma-free5gc/free5gc/scripts/patch-nf.sh amf apply
```

### 10.2 Ligar a política do AAnF

Compile e instale de novo o AAnF e o `seed` (seções 3.1 e 3.4):

```bash
free5gc/scripts/build-aanf.sh
```

```bash
free5gc/scripts/deploy.sh usuario@maquina-do-free5gc
```

O `deploy.sh` não sobrescreve a configuração instalada. Acrescente ao bloco `akma` do `~/free5gc/config/aanfcfg.yaml` a idade a partir da qual o AAnF pede a renovação:

```yaml
  akma:
    kakmaRefresh: 1m
```

`0`, ou a linha ausente, desliga a renovação. Com um valor, quando um AF pede a K_AF de uma K_AKMA registrada há mais tempo que isso, o AAnF entrega a K_AF e pede a renovação ao UDM em segundo plano. Pede uma vez por K_AKMA, e de novo depois de 5 minutos se ela não mudar. Um minuto serve para ver a renovação; numa rede de verdade, o valor é a política de autenticação do operador.

### 10.3 Ver a renovação

Suba o free5GC com o AAnF (`./run.sh -aanf`, parte 5), o gNB e o UE no *namespace* (parte 7). Espere um minuto depois do registro do UE e faça um Ua\* com um AF_ID que o AF ainda não usou com esse A-KID. O AF guarda cada K_AF até ela vencer e, para um AF_ID já usado, não consultaria o AAnF:

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh ua -profile psk-tls13 -n 1
```

O AF recebe a K_AF, e a tela do free5GC mostra a cadeia, na ordem em que apareceu no laboratório em 28/10/2026 (140 ms, do pedido do AAnF ao contexto novo):

```
[INFO][AANF][AKMA] K_AKMA of imsi-208930000000001 registered 1m5s ago: refresh requested from the UDM (TS 33.535 cl. 6.4.4)
[INFO][AANF][AKMA] UDM for Nudm_UECM_AuthTrigger: http://127.0.0.3:8000 (discovered at the NRF)
[INFO][AMF][AKMA] Re-authentication of imsi-208930000000001 requested by the UDM: started (ngKSI 1)
[INFO][UDM][AKMA] AuthTrigger for imsi-208930000000001 from AANF at 127.0.0.1: re-authentication notified to the AMF …
[INFO][AANF][AKMA] Refresh of K_AKMA of imsi-208930000000001 accepted by the UDM: new primary authentication requested
[INFO][AUSF][AKMA] K_AKMA and A-KID derived for imsi-208930000000001: A-KID rid0000.atid…
[INFO][AMF][AKMA] Re-authentication of imsi-208930000000001 requested by the UDM: completed, new KAMF in use (ngKSI 1)
[INFO][AANF][AKMA] AKMA context of imsi-208930000000001 registered, A-KID rid0000.atid… (previous context replaced)
```

O UE mostra `AKMA key material derived, A-KID[…]` com o A-KID novo, igual ao do AUSF e ao do AAnF. A sessão PDU continua: o próximo Ua\* passa pelo mesmo `uesimtun0`, com o A-KID novo. O AF ainda aceita a K_AF anterior até ela vencer (TS 33.535 §6.4.2).

### 10.4 Pedir a renovação à mão

O `seed` manda o mesmo pedido que o AAnF, sem esperar a política. Na máquina do free5GC:

```bash
~/akma-tools/seed trigger-auth --nrf http://127.0.0.10:8000 --supi imsi-208930000000001
```

A saída é `accepted (204)`. Um segundo pedido para o mesmo SUPI em menos de 60 s recebe `UDM answered 403 (REAUTHENTICATION_NOT_ALLOWED; …)`: é a política do UDM (TS 33.501 §6.1.5.2, NOTA A).

### 10.5 Limitações

- Só o acesso 3GPP, e sem a troca das chaves AS, que é opcional (TS 33.501 §6.9.4.4): o UERANSIM não tem segurança AS.
- O UE em CM-IDLE (resposta `202` e autenticação pendente) só é verificado nos testes de unidade do AMF: o UERANSIM não fica em CM-IDLE.
- Depois de reiniciar o UDM, o pedido recebe `404 CONTEXT_NOT_FOUND` até o UE se registrar de novo: o UDR do free5GC não guarda o endereço de notificação do AMF.

Para desfazer: `kakmaRefresh: 0` desliga a renovação. `patch-nf.sh amf revert` tira o *patch* do AMF. Para deixar o UDM só com o `0001` da parte 4, `patch-nf.sh udm revert` seguido de `patch-nf.sh udm apply ~/udm-0001`.

## 11. CAPIF: o AF obtém o *token* do CCF

Na parte 9, o AF fora da rede do operador obtém o *token* da API `3gpp-akma` no NRF. TS 29.522 §6 admite isso por configuração local quando o CAPIF não é usado. Com o CAPIF, o *token* vem do CAPIF Core Function (CCF): o NEF publica a API no CCF e confere o *token*, e o AF é o *API invoker* (TS 29.522 §7.2; TS 33.122 §6.5.2.3).

O free5GC não tem CCF. Esta parte instala o CCF do projeto (`free5gc/ccf/`) e põe o AF no papel de *API invoker*. A especificação é `docs/free5gc-side/capif.md`.

### 11.1 Instalar o CCF

Na máquina do AF, compile o CCF. O binário é estático:

```bash
free5gc/scripts/build-ccf.sh
```

```bash
scp free5gc/build/bin/ccf usuario@maquina-do-free5gc:free5gc/bin/ccf
```

Na máquina do free5GC, copie a configuração:

```bash
cp ~/akma-free5gc/free5gc/config/ccfcfg.yaml ~/free5gc/config/
```

As chaves e os certificados vêm de `make-capif-certs.sh` (seção 9.1), com o endereço que o AF alcança. Se o endereço da máquina não for 192.168.16.130, troque-o também em `capif1e.bindingIPv4` do `ccfcfg.yaml`. O perfil `af1-external` lista as APIs e os AF_ID que o AF pode usar (`capif.md` §5.1).

### 11.2 Atualizar o `openapi`, a configuração e o `run.sh`

Com o free5GC parado, prepare de novo as duas cópias do `openapi` (seção 10.1). Agora o *script* aplica os *patches* até o `0004`, que traz os clientes do CAPIF.

Se você seguiu as partes 4 e 9 com esta versão do repositório, o NEF, a configuração e o `run.sh` já têm os *patches* do CAPIF: `patch-nf.sh nef status`, `config status` e `run status` dizem `patches applied`.

Se as seguiu com uma versão anterior, o `apply` recusa as mudanças locais. Nesse caso, antes de atualizar o repositório, guarde uma cópia dos *patches* aplicados:

```bash
mkdir -p ~/old-patches && cp -a ~/akma-free5gc/free5gc/patches/{nef,config,run} ~/old-patches/
```

Desfaça os três com essa cópia, atualize o repositório e aplique de novo:

```bash
for t in run config nef; do ~/akma-free5gc/free5gc/scripts/patch-nf.sh $t revert ~/old-patches/$t; done
```

```bash
for t in run config nef; do ~/akma-free5gc/free5gc/scripts/patch-nf.sh $t apply; done
```

O `run.sh` ganha a opção `-ccf`, que põe o CCF logo antes do NEF.

### 11.3 Subir o free5GC e emitir a credencial

```bash
./run.sh -aanf -ccf
```

Na tela devem aparecer:
- do CCF, `Start CAPIF-1e server (https, listen on 192.168.16.130:8020)`;
- do NEF, `3gpp-akma published at the CCF as <apiId>`.

O AF precisa da informação de entrada (TS 33.122 §6.1 passo 1): a CA do CCF e uma credencial de entrada. Na máquina do free5GC, emita a credencial do perfil `af1-external`:

```bash
cd ~/free5gc && ./bin/ccf credential -c config/ccfcfg.yaml --profile af1-external --validity 24h > ~/af1-external.jwt
```

Na máquina do AF, copie os dois arquivos para `af/build/capif/`, de onde `af/config/afcfg-capif.yaml` os lê:

```bash
mkdir -p af/build/capif && scp usuario@maquina-do-free5gc:free5gc/cert/ccf-ca.pem usuario@maquina-do-free5gc:af1-external.jwt af/build/capif/
```

### 11.4 Rodar o AF como *API invoker*

```bash
af/build/bin/af -c af/config/afcfg-capif.yaml
```

O AF entra no CCF, descobre o NEF, obtém o método `OAUTH` e pede o `check-authentication` ao NEF:

```
[INFO][AF][CAPIF] Onboarded at the CCF as API invoker 6cdc32e3-…
[INFO][AF][CAPIF] 3gpp-akma discovered at http://192.168.16.130:8000, AEF nef-aef
[INFO][AF][CAPIF] Security method OAUTH at AEF nef-aef
[INFO][AF][CAPIF] The NEF has the authentication data of API invoker 6cdc32e3-…
```

O AF guarda o estado em `af/build/capif/af1-external.json`, com permissão `0600`, e o usa na próxima partida. No primeiro pedido de K_AF de um UE (parte 7), o AF pede o *token* ao CCF, e o NEF o aceita:

```
[INFO][AF][CAPIF] Access token obtained, expires in 600 s
[INFO][NEF][CAPIF] Access token of API invoker 6cdc32e3-… accepted for AF_ID af1.akma.example.FF00000001
```

O `e2e.sh` (parte 8) faz tudo isso sozinho: sobe o CCF, emite a credencial, roda o AF do modo CAPIF na 8086, o segundo *invoker* `af/scripts/negative/capif_api.py` e o Ua\* do UE com esse AF (R2c; `docs/af-side/e2e.md` §4.2).

### 11.5 Limitações

- O CAPIF é parcial: só o caminho do AKMA (DV-14).
- CAPIF-3/4 e CAPIF-2e ficam sem TLS até o MVP-11 (DV-1 e DV-9).
- Os certificados não têm revogação, nem CRL nem OCSP, e o *token* não leva `x5t#S256` (DV-16).
- O CCF guarda tudo na memória. Depois de reiniciá-lo, o NEF publica de novo, e o AF entra de novo quando o CCF não o conhece mais.
- Um *token* já emitido vale até `exp` (10 min), mesmo depois da saída do *invoker*.

Para desfazer: `./run.sh -aanf`, sem `-ccf`, deixa o CCF parado, e o AF volta à parte 9 com `af/config/afcfg-external.yaml`. Com `allowNrfTokens: true`, o NEF aceita os dois modos.
