# Tutorial — AKMA sobre o free5GC

Este tutorial descreve o que fazer **depois de instalar o free5GC** pelo guia oficial para usar AKMA (3GPP TS 33.535). Ele cresce a cada MVP; esta versão corresponde ao MVP-3.

| Parte | Conteúdo | Disponível |
|---|---|---|
| 1 | Pré-requisitos | sim |
| 2 | Conferir as derivações de chave | sim |
| 3 | Instalar o AAnF no free5GC | sim |
| 4 | Habilitar AKMA no UDM e no AUSF; provisionar um assinante | sim |
| 5 | Ver o AKMA na tela: free5GC e UERANSIM | sim |
| 6 | Escrever um AF que usa AKMA | sim |
| 7 | Preparar o UE (UERANSIM) para o Ua\* | MVP-4 |
| 8 | Executar o fluxo de ponta a ponta | MVP-5 |

---

## 1. Pré-requisitos

| Item | Versão usada no projeto |
|---|---|
| free5GC | v4.2.3, instalado pelo guia oficial, com `TestRegistration` passando |
| Go | 1.26.2 — a mesma do free5GC |
| Python | 3.10 ou mais recente (para o oráculo de referência) |
| UERANSIM | v3.2.6 (a partir da parte 4; com *patch* a partir da parte 5) |

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

### 3.3 Corrigir o NRF

O NRF do free5GC v4.2.3 tem três falhas em relação a TS 29.510 que afetam o AAnF (`docs/free5gc-side/nrf.md`):

- descarta o `aanfInfoList` do perfil;
- recusa a descoberta do tipo `AANF`;
- emite tokens sem `aud`.

O AAnF exige `aud` no token. **Sem esta correção, o registro de chaves responde 401.** Na máquina do free5GC, com o free5GC parado:

```bash
free5gc/scripts/patch-nf.sh nrf apply
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

Com o free5GC em execução (pelo menos o NRF), inicie o AAnF na máquina do free5GC:

```bash
free5gc/scripts/run-aanf.sh start
```

Ou, a partir da outra máquina:

```bash
ssh usuario@maquina-do-free5gc bash -s -- start < free5gc/scripts/run-aanf.sh
```

O log fica em `~/free5gc/log/aanf/aanf.log`. As linhas esperadas são `OAuth2 setting receive from NRF: true` e `Start SBI server (listen on …:8000)`. `run-aanf.sh stop` encerra o AAnF, que antes se descadastra do NRF.

### 3.5 Verificar o registro e a descoberta no NRF

```bash
mongosh --quiet free5gc --eval 'printjson(db.NfProfile.find({nfType: "AANF"}, {_id: 0}).toArray())'
```

O perfil deve ter `nfType: 'AANF'`, o serviço `naanf-akma` (versão `v1`) no endereço configurado e o `aanfInfoList` com os RIDs. Para conferir a descoberta, rode na máquina do free5GC:

```bash
seed discover --nrf http://127.0.0.10:8000 --routing-indicator 0000
```

O AAnF deve aparecer na lista. Com um RID que ele não atende, a resposta é `null`. É assim que o AUSF vai encontrá-lo na parte 4 (TS 33.535 §6.7).

### 3.6 Testar as três operações

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
| `401` no `register` | sem `--nrf`; o NRF sem a correção da seção 3.3 (token sem `aud`); ou o `nrfCertPem` do AAnF não é o certificado do NRF em uso |
| `403` sem causa no `retrieve` | o AF_ID não está em `authorizedAfs`, ou o pedido saiu de um endereço que não está na lista daquele AF_ID |
| `403 (K_AKMA_NOT_PRESENT)` | o A-KID não foi registrado, foi substituído por um registro mais novo do mesmo SUPI ou o AAnF reiniciou (o contexto fica só em memória) |

Por fim, para remover o contexto: `seed remove --nrf http://127.0.0.10:8000 --supi imsi-208930000000001`.

---

## 4. Habilitar AKMA no UDM e no AUSF; provisionar um assinante

Depois de uma autenticação primária bem-sucedida de um assinante com assinatura AKMA:

- o UDM diz ao AUSF que o UE usa AKMA e qual é o seu Routing Indicator;
- o AUSF deriva K_AKMA e o A-KID da K_AUSF e os entrega ao AAnF da parte 3 (TS 33.535 §6.1).

Esta parte altera o UDM e o AUSF do free5GC com *patches* do repositório (`docs/free5gc-side/udm.md`, `ausf.md`).

### 4.1 Onde fica o repositório

O AUSF corrigido depende de código deste repositório: da biblioteca `shared/akma` e da cópia do `openapi` da parte 3. Os `go.mod` corrigidos apontam para `~/akma-free5gc` por caminho relativo. **Clone o repositório nesse local, ao lado de `~/free5gc`**, na máquina do free5GC, e prepare as duas cópias do `openapi`: a das NFs e a do banco de testes do free5GC.

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

**Pelo UERANSIM.** Inicie o free5GC (`./run.sh`), o AAnF (`run-aanf.sh start`), o gNB e o UE. No log do free5GC deve aparecer `AKMA key material of imsi-208930000000001 registered at the AAnF`. O log do AAnF (`~/free5gc/log/aanf/aanf.log`) mostra o A-KID guardado, da forma `rid0000.atid<44 caracteres>@5gc.mnc093.mcc208.3gppnetwork.org`. A parte 5 põe tudo isso numa tela só.

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

Na máquina do UERANSIM, aplique os *patches* do UE. O *script* guarda `build/nr-ue.orig` e recompila só o `nr-ue`; na primeira vez, leva alguns minutos.

```bash
~/akma-free5gc/af/scripts/patch-ueransim.sh apply
```

O *patch* liga a opção `akma: true` no `config/free5gc-ue.yaml`. Se você usa outro arquivo de configuração, acrescente essa linha a ele.

### 5.3 O que procurar

Inicie o gNB e o UE como na parte 4. Na tela do UE, logo depois do `Security Mode Command received`:

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
- `listen`: o endereço em que o AF escuta;
- `aanfUri`: o endereço do AAnF;
- `kafExpiry`: o que fazer quando a K_AF vence, `refetch` ou `reject`.

```bash
af/build/bin/af -c af/config/afcfg.yaml
```

Na tela, a partida mostra o AF_ID e o `realm`:

```
[INFO][AF][Main] AF_ID af1.akma.example.FF00000001; realm "3GPP-bootstrapping-akma@af1.akma.example"; ...
```

### 6.3 Testar sem UE

O UE do UERANSIM só entrega a K_AF a partir da parte 7. Até lá, faça o papel do UE com um vetor do projeto, em dois passos:

1. semeie no AAnF o contexto do vetor `testbed-ue`, como na seção 3.6 (`seed register`);
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
| K_AF errada (`-kaf` com outro valor) | `401` |
| A-KID que o AAnF não conhece (`-vector rid-012`) | `403 AKMA_KEY_NOT_AVAILABLE` |
| outro nome no `Host` (`-fqdn af2.akma.example`) | `421 MISDIRECTED_REQUEST` |
| AF_ID fora da política do AAnF | `503 AANF_UNAVAILABLE` |

Os resultados do laboratório estão em `docs/af-side/teste-af.md`.

### 6.4 Limitações

Sem TLS, o A-KID e o conteúdo trafegam em claro, e o AF só é autenticado pelo `rspauth` (`docs/af-side/ua-star.md` §10). É um protocolo de laboratório; o perfil normativo, com TLS, fica para depois do MVP.

