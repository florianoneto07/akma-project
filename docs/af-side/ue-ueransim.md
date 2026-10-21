# UE (UERANSIM): K_AKMA, A-KID e K_AF no ME

**MVP:** MVP-3a (K_AKMA, A-KID e log) e MVP-4 (K_AF, `nr-cli`, autoteste, *namespace* do UE) · **Data:** 18 de outubro de 2026 · **Verificado** na VM do free5GC em 18/10/2026 (`teste-ue.md`)

## 1. O que o UE precisa fazer

`[3GPP]` TS 33.535 §6.1: *"The UE shall generate the AKMA Anchor Key (KAKMA) and the A-KID from the KAUSF before initiating communication with an AKMA Application Function."* O §6.2.1, passo 1, repete a exigência, acrescenta que o UE envia o A-KID ao AF e diz que *"The UE may derive KAF before sending the message or afterwards."* §4.4.3: *"Applications on the UE shall not be able to get access to KAKMA."*

| Valor | Derivação | Referência |
|---|---|---|
| K_AKMA | KDF(K_AUSF), FC `0x80`, P0 `"AKMA"`, P1 SUPI | TS 33.535 A.2 |
| A-TID | KDF(K_AUSF), FC `0x81`, P0 `"A-TID"`, P1 SUPI | TS 33.535 A.3 |
| A-KID | `rid<RID>.atid<A-TID>@5gc.mnc<MNC>.mcc<MCC>.3gppnetwork.org` | TS 33.535 §6.1; TS 29.522 Tab. 5.14.5.4.2-1; ADR-0002 |
| K_AF | KDF(K_AKMA), FC `0x82`, P0 AF_ID = FQDN do AF ‖ identificador Ua\* (5 octetos) | TS 33.535 A.4; TS 33.220 Anexo H |

- **SUPI:** o IMSI, só os dígitos (TS 33.501 A.7.0).
- **RID:** o do SUCI do UE (TS 23.003 §2.2B).
- **Domínio:** o da rede de origem, com o MNC em três dígitos (TS 23.003 §28.2).
- **A-TID:** em base64 com preenchimento, 44 caracteres (RFC 4648 §4), conforme o ADR-0002.
- **FQDN na K_AF:** em UTF-8 com NFKC (TS 33.220 B.2.1.2). `[MVP]` O UE não aplica NFKC e recusa FQDN não ASCII (DV-4).

`[FREE5GC]` O UERANSIM (a versão instalada na VM do free5GC, *commit* `85a0fbfd`) calcula a K_AUSF nos dois métodos de autenticação primária, mas não tem nada de AKMA.

## 2. Quando o UE deriva

`[INTERPRETAÇÃO]` O UE deriva K_AKMA e A-KID quando passa a usar o contexto de segurança NAS criado pela autenticação primária bem-sucedida. É o momento em que o `nr-ue` aceita o *Security Mode Command* que ativa esse contexto:

- **5G AKA:** o UE não recebe confirmação explícita de sucesso. O *Security Mode Command* que ativa o contexto novo é a primeira prova de que a rede aceitou o RES\*.
- **EAP-AKA':** o EAP-Success chega no mesmo *Security Mode Command*.

Isso atende ao "antes de comunicar com um AF" do §6.1: nenhuma aplicação fala com um AF antes de o UE se registrar. Cada autenticação primária nova gera K_AKMA e A-KID novos, como acontece no AUSF.

O material fica no contexto de segurança NAS (`UeKeys`) e é descartado com ele. A K_AF é derivada sob demanda, a cada pedido de uma aplicação, a partir da K_AKMA corrente.

## 3. Entrega à aplicação: serviço de chaves do ME

**Desde o MVP-8, fase (a)** (`ua-star-norma.md` §2; *patch* `0005`). `[3GPP]` TS 33.535 §4.4.3: a aplicação só obtém a K_AF dos AF_IDs a que tem direito e não acessa a K_AF de outra aplicação. A forma de garantir isso fica fora do escopo 3GPP.

**Como funciona.**

- **Canal.** O `nr-ue` atende num *socket* Unix, `/run/ueransim/akma-<SUPI>.sock`.
- **Quem pede.** O kernel diz ao ME o usuário e o processo de quem conecta (`SO_PEERCRED`). O ME confere o *uid* e, se a política trouxer, o executável do processo.
- **O que pode receber.** A política na configuração do UE dá a lista de AF_IDs de cada aplicação. Fora dela, o ME responde `ERR NOT_AUTHORIZED` e registra um `warn`.
- **Protocolo.** Uma linha por pedido:

```
KAF af1.akma.example FF00000001
OK a-kid=rid0000.atid…=@5gc.mnc093.mcc208.3gppnetwork.org af-id=af1.akma.example.FF00000001 k-af=<64 hexadecimais>
```

- A K_AKMA nunca sai do ME: nenhuma resposta a leva.
- **Erros:** `AKMA_DISABLED`, `NO_KEY_MATERIAL` (antes de uma autenticação primária), `BAD_REQUEST` (FQDN que não é nome DNS ASCII, ou identificador sem dez hexadecimais), `NOT_AUTHORIZED` e `UNAVAILABLE`.
- **Cliente.** O cliente Ua\* (`ua-client.md`) usa o serviço com `-me-socket`, que o `ue-netns.sh ua` passa.

**O `akma-kaf` do `nr-cli`.** Até o MVP-7, era o caminho da aplicação. Agora só funciona com `akmaCliKaf: true` (teste), e o `nr-ue` avisa na partida que a K_AF fica exposta a qualquer processo local. Sem a opção, o comando responde *"akma-kaf is disabled"*.

```bash
nr-cli imsi-208930000000001 -e "akma-kaf af1.akma.example"
```

## 4. Configuração

`[MVP]` Opção `akma` na configuração do UE, prevista no plano. O padrão é `false`, que preserva o comportamento original do UERANSIM. O *patch* `0002` liga a opção no `free5gc-ue.yaml`:

```yaml
# AKMA (3GPP TS 33.535): derive K_AKMA and the A-KID after each successful primary authentication
akma: true
```

Os demais dados vêm do que o UE já tem:

- `supi`: só IMSI, como na rede (DV-3); com outro tipo de SUPI, o UE registra um aviso e não deriva;
- `routingIndicator`: o mesmo que vai no SUCI. `[FREE5GC]` Sem o campo, o UERANSIM põe `"0000"` no SUCI, e o A-KID usa o mesmo valor;
- `mcc` e `mnc`: o comprimento do MNC vem da configuração, e não do IMSI.

**Política do serviço de chaves (MVP-8).** O *patch* `0005` acrescenta ao `free5gc-ue.yaml`:

```yaml
akmaKeyService:
  applications:
    - name: uaclient
      uid: 1000            # o usuário que roda a aplicação (id -u)
      afIds:
        - af1.akma.example.FF00000001
akmaCliKaf: false
```

- `executable`, opcional: o caminho absoluto do executável da aplicação.
- `socket`, opcional: outro caminho para o *socket*.
- Um AF_ID em duas aplicações faz o `nr-ue` recusar a configuração (TS 33.535 §4.4.3, item 3).

## 5. O UE no próprio *namespace* de rede

`[FREE5GC]` Com o UERANSIM e o free5GC na mesma máquina e o UE no *namespace* do host, o endereço da sessão PDU (10.60.0.1) é local ao host. O UPF entrega os pacotes do UE pela `upfgtp` com esse endereço de origem, e o kernel os descarta. Era o que impedia o plano de usuário desde o MVP-1: o UE se registrava e abria a sessão, mas não alcançava nada pelo `uesimtun0`. O *script* de instalação da VM (`~/free5gc-setup/05-ue-netns.sh`) já previa o *namespace*, mas nunca tinha rodado.

`[MVP]` Decisão D5: o UE roda no *namespace* `ueransim-ue`, ligado ao host por um par `veth`. O *script* `af/scripts/ue-netns.sh` cuida disso:

| Comando | O que faz |
|---|---|
| `sudo ue-netns.sh up` | cria o *namespace* e o `veth`: 10.100.0.1 no host, 10.100.0.2 no *namespace*, com rota padrão pelo host |
| `sudo ue-netns.sh ue` | roda o `nr-ue` no *namespace* |
| `sudo ue-netns.sh cli <argumentos>` | roda o `nr-cli` no *namespace*, com o seu usuário. O `nr-cli` fala com o `nr-ue` pela *loopback*, que é a do *namespace* |
| `sudo ue-netns.sh ua [opções]` | roda o cliente Ua\* no *namespace*, com o seu usuário (`ua-client.md`) |
| `sudo ue-netns.sh status`, `down` | mostra ou remove o *namespace* |

- **gNB:** o rádio simulado passa pelo `veth`. O gNB escuta em `linkIp: 10.100.0.1`; essa mudança no `free5gc-gnb.yaml` é manual, porque o arquivo já tem a alteração do endereço do AMF feita na instalação.
- **UE:** procura o gNB em `gnbSearchList: 10.100.0.1`, pelo *patch* `0004`.
- **Plano de usuário:** só o rádio passa pelo `veth`. O tráfego do UE sai pelo `uesimtun0` e pelo plano de usuário: o `nr-ue` instala a regra `from 10.60.0.x lookup rt_uesimtun0`. Para isso, o cliente usa o endereço da sessão como origem.
- **N6:** o encaminhamento e o NAT estão em `docs/free5gc-side/plano-de-usuario.md`.

O *namespace* e o `veth` não sobrevivem a uma reinicialização da VM: rode `up` de novo.

### 5.1 gNB e UE em outra máquina (MVP-6)

O *namespace* funciona igual numa máquina só do UERANSIM. No laboratório essa máquina é a 192.168.16.131, com o free5GC na 192.168.16.130. No `free5gc-gnb.yaml` dessa máquina:

- `ngapIp` e `gtpIp` passam a ser o endereço dela;
- o AMF passa a ser `192.168.16.130`;
- o `linkIp` continua `10.100.0.1`.

Ela não precisa de NAT: o N6 continua na máquina do free5GC.

Na máquina do free5GC, o gNB local, se houver, precisa de `ngapIp` igual ao endereço da máquina e de `gtpIp` 10.100.0.1:

- um endereço de *loopback* não alcança o AMF nem o UPF em 192.168.16.130;
- o próprio 192.168.16.130 conflita com a porta 2152 do UPF.

Procedimento: tutorial, parte 9.5. Resultado: `teste-nef.md`.

## 6. Autoteste com os vetores do projeto

`[MVP]` O executável `nr-akma-selftest` confere as derivações do ME (K_AKMA, A-KID, K_AF) contra os vetores do MVP-0 (`shared/akma/testdata/vectors.json`), lidos com o *parser* YAML do UERANSIM:

```bash
~/UERANSIM/build/nr-akma-selftest ~/akma-free5gc/shared/akma/testdata/vectors.json
```

Os três vetores de NFKC têm FQDN não ASCII. Neles, o autoteste confere K_AKMA e A-KID e espera a recusa da K_AF (DV-4). Com isso, a derivação tem quatro implementações que concordam:

- o oráculo Python;
- a biblioteca Go (`shared/akma`);
- a KDF do free5GC (`ueauth`);
- a KDF do UERANSIM, em C++.

## 7. Log

O UERANSIM imprime cada linha como `[data hora] [componente] [nível] mensagem`, com nível a partir de `debug`:

| Componente | Nível | Mensagem | Quando |
|---|---|---|---|
| `nas` | `info` | `AKMA key material derived, A-KID[<A-KID>]` | contexto novo ativado e `akma: true` |
| `nas` | `warn` | `AKMA key material not derived, the SUPI is not an IMSI` | `akma: true` com SUPI de outro tipo |
| `app` | `info` | `AKMA K_AF derived for AF_ID[<AF_ID>], A-KID[<A-KID>]` | pedido `akma-kaf` atendido |

O formato `texto[valor]` é o mesmo das mensagens do UERANSIM, como `Selected integrity[2] ciphering[0]`. K_AUSF, K_AKMA e K_AF nunca aparecem no log. O A-KID impresso é idêntico ao que o AUSF, o AAnF e o AF registram.

## 8. Arquivos tocados

Árvore: `~/UERANSIM`, na VM do free5GC (D5). Binários: `build/nr-ue`, com o original em `build/nr-ue.orig`, e `build/nr-akma-selftest`. *Patches*: `af/ue/ueransim/patches/`.

| Arquivo | Mudança | *Patch* |
|---|---|---|
| `src/ue/types.hpp` | `UeConfig::akma`; `UeKeys::kAkma` e `UeKeys::aKid`, copiados em `deepCopy` | `0001-akma-key-derivation.patch` |
| `src/ue.cpp` | lê `akma`; copia a opção para cada UE (`-n`) | idem |
| `src/ue/nas/keys.{hpp,cpp}` | `CalculateKAkma`, `CalculateATid`, `ConstructAKid`, `DeriveAkmaKeys`; codificação base64 | idem |
| `src/ue/nas/mm/mm.hpp`, `security.cpp` | `NasMm::deriveAkmaKeys`, chamada quando o contexto novo vira o corrente | idem |
| `config/free5gc-ue.yaml` | `akma: true` | `0002-free5gc-ue-akma.patch` |
| `src/ue/nas/keys.{hpp,cpp}` | `IsValidAfFqdn`, `CalculateKAf` | `0003-akma-kaf-nr-cli-selftest.patch` |
| `src/lib/app/cli_cmd.{hpp,cpp}` | comando `akma-kaf <af-fqdn> [<id Ua*>]` | idem |
| `src/ue/app/cmd_handler.cpp` | atende `akma-kaf` com a K_AKMA do contexto corrente | idem |
| `src/akma_selftest.cpp` (novo), `CMakeLists.txt` | `nr-akma-selftest` | idem |
| `config/free5gc-ue.yaml` | `gnbSearchList: 10.100.0.1` | `0004-free5gc-ue-gnb-on-veth.patch` |

A KDF é a do próprio UERANSIM (`crypto::CalculateKdfKey`), a mesma que ele usa para K_SEAF e K_AMF. O UERANSIM não tem base64, então o *patch* traz um codificador pequeno.

## 9. Como reproduzir numa instalação limpa

Na máquina do UERANSIM, com o repositório em `~/akma-free5gc`:

```bash
~/akma-free5gc/af/scripts/patch-ueransim.sh apply
```

O *script* faz quatro coisas:

- confere a série inteira numa cópia temporária dos arquivos, porque cada *patch* depende dos anteriores, e a aplica em ordem;
- guarda `build/nr-ue.orig`;
- recompila o `nr-ue` e o `nr-akma-selftest` em `cmake-build-release`, com o `cmake` que configurou a árvore. O UERANSIM exige o 3.17 ou mais novo, e na VM ele fica em `~/.local/cmake-3.31`;
- copia os binários para `build/`.

Alterações locais em arquivos que os *patches* não tocam, como o `free5gc-gnb.yaml`, não bloqueiam a aplicação. Para desfazer, use `revert`. Para trocar por uma série mais nova, reverta com a série antiga antes de atualizar o repositório.

Depois, a mudança manual do gNB e o *namespace* (seção 5):

```bash
sed -i 's/^linkIp: 127.0.0.1 /linkIp: 10.100.0.1 /' ~/UERANSIM/config/free5gc-gnb.yaml
```

```bash
sudo ~/akma-free5gc/af/scripts/ue-netns.sh up
```

## 10. Como verificar

O autoteste da seção 6, e a sequência de `teste-ue.md`:

- o UE se registra no *namespace*;
- o A-KID dele é igual ao do AAnF;
- o cliente Ua\* autentica no AF com a K_AF do ME, pelo plano de usuário;
- uma nova autenticação troca o A-KID.
