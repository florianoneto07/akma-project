# Plano de desenvolvimento por MVPs

**Estado:** **Aprovado** em 16 de outubro de 2026, com as decisões registradas na seção 1. Em execução: MVP-0 implementado (16/10/2026); MVP-1 implementado e verificado na VM do free5GC (17/10/2026); MVP-1b (correção do NRF do free5GC) implementado e verificado (17/10/2026); MVP-2 (UDM e AUSF) implementado e verificado (17/10/2026); MVP-3a (mensagens AKMA na tela do free5GC e do UERANSIM, com a derivação no ME antecipada do MVP-4) implementado (18/10/2026); MVP-3 (AF e Ua\* simplificado) implementado e verificado (18/10/2026); MVP-4 (K_AF no ME, cliente Ua\* no UE pelo plano de usuário) implementado e verificado (18/10/2026); MVP-5 (roteiro ponta a ponta e testes negativos) implementado e executado, com 63 verificações e nenhuma falha (18/10/2026), e reproduzido numa instalação limpa do free5GC v4.2.3 seguindo o tutorial (19/10/2026); MVP-6 (AF externo via NEF, acesso anônimo, remoção por OAM e UERANSIM em outra máquina) implementado e verificado, com 83 verificações e nenhuma falha (19/10/2026). MVP-7 (repositório público) concluído: publicado no GitHub e instalado numa VM limpa só pelo guia, com 83 verificações e nenhuma falha (20/10/2026). MVP-8 (Ua\* da norma e autorização por aplicação no UE): organização e especificação das fases (a) e (b) aprovadas; fase (a), o serviço de chaves do ME, implementada e verificada, com 88 verificações e nenhuma falha (20/10/2026); fase (b), os perfis TLS (K.2, PSK TLS 1.3 e PSK TLS 1.2), implementada e verificada, com 115 verificações e nenhuma falha (21/10/2026); fase (c), os perfis DTLS 1.3 e DTLS 1.2 com CoAP, implementada e verificada, com 126 verificações e nenhuma falha (22/10/2026); fase (d), o perfil OSCORE, implementada e verificada, com 141 verificações e nenhuma falha (22/10/2026); o CAPIF passou a MVP-9. MVP-12 (renovação da K_AKMA pelo AAnF) com a especificação aprovada (`docs/free5gc-side/renovacao-kakma.md`, pontos K1 a K12), implementado e verificado, com 161 verificações e nenhuma falha (28/10/2026). MVP-13 (V2 do repositório público) aprovado em 28/10/2026, com os pontos L1 a L13; em execução.
**Data:** 15 de outubro de 2026; revisão 3 em 16 de outubro de 2026; aprovação em 16 de outubro de 2026
**Base:** documentos `00` a `04`, Adendo E, ADR-0001-R2, ADR-0002, ADR-0003, `docs/baseline/ambiente-verificado.md`, `docs/baseline/conferencia-normativa.md`, `docs/literatura/analise-dos-artigos.md`, `docs/conformidade/desvios.md`
**Revisão 2 (15/10/2026):** conferido contra as normas do acervo (hoje em `specs/`).
**Revisão 3 (16/10/2026):** conferido contra TS 33.535, TS 29.535 (com o YAML), TS 23.502 V19.9.0 e TR 121 917; incorporada a análise dos artigos. Fica fora do acervo apenas TS 29.505 (D8).

Convenção: as marcações `[3GPP]`, `[FREE5GC]`, `[MVP]`, `[INTERPRETAÇÃO]`, `[LACUNA]` e `[DIVERGÊNCIA]` têm o significado do README. Cada MVP começa pela conferência, no texto da norma, das cláusulas que usa.

---

## 1. Decisões

Estado em 16/10/2026. "Aprovada" = recomendação aprovada com o plano; "decidida" = escolha explícita do responsável pelo projeto.

| Nº | Decisão | Resultado | Estado | Afeta |
|---|---|---|---|---|
| D1 | Onde ficam normas e artigos | Normas em `specs/` (`.docx` e PDF) e `specs/openapi/` (YAML e ABNF); artigos em `papers/`; `.gitignore` exclui `/specs/` e `/papers/`; `reference/` fica só com o oráculo da KDF. Os artigos são cópias licenciadas de editoras e não podem ir para um repositório | **executada** em 16/10/2026 (71 arquivos movidos, SHA-256 idênticos antes e depois) | primeiro commit |
| D2 | Raiz e unicidade do repositório | Um único repositório em `~/akma/akma-free5gc/`, com documentação e código. O README, que hoje diz que o código vive em outro repositório, será atualizado | aprovada | MVP-0 |
| D3 | Ferramentas na VM do AF | git (pelo usuário, via `apt`); Go 1.26.2 em diretório do usuário, sem `sudo`, mesma versão da VM do free5GC | aprovada | MVP-0 |
| D4 | Versão-alvo do free5GC | **v4.2.3** (instalada e testada) para MVP-0 a MVP-5; *rebase* para a versão corrente no início do trabalho *upstream* | aprovada | MVP-1 |
| D5 | Localização do lado UE | UERANSIM permanece na VM do free5GC. O ME (*patch* no `nr-ue`) deriva as chaves; o cliente Ua\* roda no *namespace* do UE e alcança o AF pelo plano de usuário (`uesimtun0` → UPF → N6 → `192.168.16.129`). O código-fonte fica no repositório da VM do AF e é documentado em `docs/af-side/` | aprovada | MVP-4 |
| D6 | Onde ficam os modelos do Naanf_AKMA | *Patch* no `free5gc/openapi` v1.2.4 com `AkmaKeyInfo`, `CtxRemove`, `AkmaAfKeyRequest`, `AkmaAfKeyData`, `ServiceDisableNotif` e um cliente `aanf/AKMA`, consumido por `replace` no AUSF, no AAnF e no AF | aprovada | MVP-1 |
| D7 | Onde fica o código comum (KDF, SUPI, A-KID, AF_ID, codecs) | Módulo Go `shared/akma` no repositório, consumido por `replace` pelo AUSF, pelo AAnF e pelo cliente Ua\*; migra para `free5gc/util` no *upstream* | aprovada | MVP-0 |
| D8 | TS 29.505 no baseline | Acrescentar. `[FREE5GC]` O `AuthenticationSubscription` do UDR já tem `akmaAllowed` e `routingId`. `[3GPP]` TS 33.535 §3.1 define *"AKMA subscription data"* como o dado que indica se o assinante pode usar AKMA, e §4.2.5 o põe no UDM; a semântica Stage 3 do campo está em TS 29.505 | aprovada; **norma obtida em 17/10/2026**: TS 29.505 V19.6.0 e `TS29505_Subscription_Data.yaml` em `specs/`. Tab. 5.4.2.2-1: `akmaAllowed` (O) indica se o assinante pode usar AKMA; `routingId` (C) *"shall be present if the akmaAllowed is set to "true""* e *"shall include the Routing Indicator contained in the SUCI"*; no YAML, cadeia `^[0-9]{1,4}$` | MVP-2 |
| D9 | Momento do registro no AAnF e efeito de falha | `[3GPP]` TS 33.535 §6.1 passo 3: K_AKMA e A-TID são gerados *"after the primary authentication procedure is successfully completed"*; a norma não condiciona a autenticação primária ao registro. `[INTERPRETAÇÃO]` O AUSF responde ao AMF normalmente e faz `AnchorKey_Register` em seguida, sem bloquear; falha no AAnF é registrada em log e não altera o resultado da autenticação | **decidida** em 16/10/2026 | MVP-2 |
| D12 | Versão de TS 23.502 | O baseline permanece **V19.8.0**; a cópia do acervo é a V19.9.0 e o projeto segue com ela para consulta. Registrado no Adendo E e em `especificacoes.md` | **decidida** em 16/10/2026 | registro do baseline |
| D13 | Limitações do NRF do free5GC v1.4.5 encontradas no MVP-1: descarta `aanfInfoList` (N1), recusa descoberta de `AANF` (N2), emite token com `aud` vazio (N3) — `docs/free5gc-side/registro-nrf.md` §5 | Corrigir o NRF com três *patches*, **agora**, numa unidade própria (MVP-1b) antes do MVP-2. Primeira decisão (17/10/2026): "mais adiante"; revista no mesmo dia, depois da conferência com TS 29.510, por recomendação do projeto aceita pelo responsável. Com o NRF corrigido, o AUSF do MVP-2 descobre o AAnF pelo NRF, com configuração local como alternativa (TS 33.535 §6.7), e o DV-8 é encerrado | **decidida** e **executada** em 17/10/2026 (MVP-1b) | MVP-1b, MVP-2 |
| D14 | Escopo do MVP-3a, pedido pelo responsável em 18/10/2026: as mensagens AKMA principais aparecem na tela do free5GC e do UERANSIM, no formato e nos níveis de log de cada um | `[FREE5GC]` O UERANSIM não tem nada de AKMA, e não há o que registrar sem derivação. A derivação de K_AKMA, A-TID e A-KID no ME, com a opção `akma` na configuração do UE, sai do MVP-4 e entra no MVP-3a; K_AF, o comando do `nr-cli`, o autoteste com os vetores e o cliente Ua\* ficam no MVP-4 | **decidida** em 18/10/2026 | MVP-3a, MVP-4 |
| D15 | Como o AAnF aparece na tela do free5GC | Opção `-aanf` no `run.sh`, no molde do `-bsf`, e `aanf` na lista do `force_kill.sh`, por *patch* (`free5gc/patches/run/`). O `run-aanf.sh` continua para uso avulso | **decidida** em 18/10/2026 | MVP-3a |
| D16 | Nível dos identificadores no log | `[MVP]` SUPI, A-KID e AF_ID em `info`, na categoria `AKMA` de UDM, AUSF e AAnF; chaves nunca, em nível algum. O free5GC já registra o SUPI em `info`, e o A-KID vai em claro no Ua\*. Até o MVP-2, o AUSF e o AAnF só mostravam os identificadores em `debug` | **decidida** em 18/10/2026 | MVP-3a |

Decisão que dependia de norma até então fora do acervo:

| Nº | Questão | Resultado | Estado | MVP |
|---|---|---|---|---|
| D10 | De onde vem o RID quando o UDM recebe SUPI, e não SUCI (por exemplo, registro por GUTI com reautenticação). `[3GPP]` TS 29.503 Tab. 6.3.6.2.3-1: `routingId` contém o RID *"contained in the SUCI"*; TS 33.535 §6.1: *"The AUSF shall use the RID received from the UDM"* — o RID que o AMF pode mandar ao AUSF (TS 29.509) não serve | `[3GPP]` Fonte: `routingId` do `AuthenticationSubscription` no UDR, que TS 29.505 Tab. 5.4.2.2-1 torna obrigatório com `akmaAllowed: true` e define como o RID do SUCI do assinante; vale com SUCI ou SUPI. `[MVP]` Monitoramento pedido pelo responsável: quando o UDM receber SUCI, comparar o RID do SUCI com o `routingId` do UDR e registrar em log qualquer divergência; durante a programação e os testes, procurar fonte melhor | **decidida** em 17/10/2026 | MVP-2 |

D11 (Digest do Ua\*) foi resolvida pela leitura de TS 24.109, TS 33.222 e TS 33.535 Anexo B; as escolhas entre textos divergentes estão na seção 6 (DV-5 e DV-6).

---

## 2. Crítica da sequência sugerida

A sequência MVP-0 → MVP-5 está essencialmente correta. Os ajustes:

1. **MVP-1, "três erros".** `ROAMING_AKMA_SERVICE_DENIED` só se aplica com a *feature* `RoamingRestriction`, fora do MVP. O MVP-1 implementa `K_AKMA_NOT_PRESENT` e `AKMA_CONTEXT_NOT_FOUND`, e testa que o terceiro **nunca** é emitido sem a *feature* negociada.
2. **Os modelos do Naanf_AKMA não existem no `free5gc/openapi`** (ambiente-verificado §3). Sem eles não há produtor (AAnF) nem consumidor (AUSF, AF). Entram no MVP-1, não no MVP-2.
3. **O `openapi` não precisa de `AkmaInd`, `RoutingId`, `AANF` nem `naanf-akma`** — já estão lá. O MVP-2 fica restrito à lógica do UDM e do AUSF; a linha "Adicionar `AkmaInd` e `RoutingId` ao modelo" do baseline §6.2 cai.
4. **O MVP-2 precisa de um UE para ser verificado, e o MVP-4 ainda não existe.** O banco de testes Go do free5GC (`test/ranUe.go`) já deriva K_AUSF em Go, no 5G-AKA e no EAP-AKA'. Com um *patch* pequeno nesse banco, o MVP-2 se verifica sozinho, nos dois métodos de autenticação primária, comparando a K_AF devolvida pelo AAnF com a derivada independentemente a partir da K_AUSF do UE de teste. O UERANSIM, por padrão, só exercita 5G-AKA.
5. **O MVP-3 não depende do MVP-2.** O AF pode ser testado contra um contexto semeado no AAnF por uma ferramenta que faz o papel do AUSF. MVP-2 e MVP-3 podem andar em qualquer ordem; se o MVP-2 parar em D8, D9 ou D10, o MVP-3 segue.
6. **A K_AF do lado UE é derivada pelo ME, não pela aplicação.** `[3GPP]` TS 33.535 §6.1 e §6.2.1 atribuem ao ME a derivação de K_AKMA, A-KID e K_AF. O MVP-4 coloca as três no `nr-ue` e entrega à aplicação (cliente Ua\*) apenas o A-KID e a K_AF do AF pedido. K_AKMA não sai do ME.
7. **Três implementações independentes da derivação.** Python (oráculo do projeto), Go (MVP-0) e C++ (MVP-4, no UERANSIM, sobre a KDF que o UERANSIM já tem). A comparação contra `ueauth.GetKDFValue` do free5GC acrescenta uma quarta implementação da KDF genérica, escrita por terceiros. Isso reduz o risco de o oráculo e o Go concordarem num erro comum de leitura.
8. **O "primeiro MVP funcional" é o MVP-5.** É nele que a pergunta-guia é respondida com evidência (seção 9). Os anteriores são incrementos verificáveis.
9. **O que a literatura muda no plano** (`docs/literatura/analise-dos-artigos.md`): nenhuma modificação de protocolo entra, mas duas verificações de implementação passam a critério de aceitação — o AAnF associa cada AF_ID autorizado à identidade do AF que o pede (contra o ataque de falsificação de AF de Akman et al., 2023), e o Ua\* confirma a posse de K_AF nos dois sentidos (sem isso a autenticação UE↔AF falha nas análises formais de Yang et al., 2022). Nenhum trabalho do acervo implementa AKMA sobre o free5GC; o mais próximo é a implementação sobre o Open5GCore (Boureanu et al., ASIA CCS '25), de código fechado.
10. **O conteúdo do MVP é do Release 17.** `[3GPP]` TR 121 917 situa AKMA e AKMA_TLS como itens de trabalho do Release 17, e TS 33.535 V19.0.0 é a 18.8.0 republicada. Os acréscimos do Release 18 (controle de *roaming*, notificação de desabilitação, DTLS, OSCORE, `AuthTrigger`) ficam no MVP-6+. A assimetria com o free5GC, que declara Release 17, não atinge o fluxo do MVP — argumento útil para a proposta *upstream*.

Sequência ajustada e dependências:

```
MVP-0  biblioteca comum: KDF, SUPI, A-KID, AF_ID, codecs
  │
  ├── MVP-1  AAnF autônomo + modelos Naanf_AKMA no openapi
  │     ├── MVP-2  UDM + AUSF → AnchorKey_Register     (verificado pelo banco de testes Go)
  │     │     └── MVP-3a  mensagens AKMA na tela + derivação no ME (UERANSIM)   (D14)
  │     └── MVP-3  AF + Ua* simplificado              (verificado com contexto semeado)
  │
  └── MVP-4  lado UE: ME no UERANSIM + cliente Ua*     (requer D5 e MVP-3)
                │
MVP-2 + MVP-3 + MVP-4 ──► MVP-5  ponta a ponta + testes negativos  = primeiro MVP funcional
```

Correspondência com as fases do README: F2 = MVP-0, F3 = MVP-1, F4 = MVP-2, F5 = MVP-3, F6 = MVP-4, F7 = MVP-5, F8+ = MVP-6+.

---

## 3. Estrutura do repositório proposta

Raiz conforme D2. Código e documentação separados pelos mesmos lados.

```
akma-free5gc/
├── specs/  specs/openapi/  papers/    normas e artigos, fora do controle de versão (D1, executada)
├── reference/kdf/                      oráculo Python (existente) + gerador de vetores
├── shared/akma/                        módulo Go comum (D7): kdf, supi, akid, afid, codec
├── free5gc/                            LADO free5GC
│   ├── aanf/                           NF nova, estrutura de NF do free5GC
│   ├── patches/openapi/                modelos e cliente do Naanf_AKMA (D6)
│   ├── patches/udm/
│   ├── patches/ausf/
│   ├── patches/test/                   banco de testes Go do free5GC
│   ├── config/aanfcfg.yaml
│   └── scripts/                        aplicar patches, compilar, provisionar assinante
├── af/                                 LADO AF / Ua*
│   ├── af/                             AF de teste
│   ├── ue/ueransim/patches/            ME: K_AKMA, A-TID, A-KID, K_AF
│   ├── ue/uaclient/                    cliente Ua*
│   └── scripts/                        ponta a ponta, testes negativos
└── docs/
    ├── 00–04, baseline/, decisions/, defects/, traceability/   (existentes)
    ├── plano/                          este documento
    ├── conformidade/desvios.md         criado no primeiro desvio aprovado
    ├── shared/                         KDF, identificadores, codificação, vetores
    ├── free5gc-side/                   AAnF, openapi, UDM, AUSF, assinatura, verificação
    ├── af-side/                        AF, Ua*, UE, cliente, scripts
    └── tutorial/                       incremental, atualizado ao fim de cada MVP
```

O código do lado free5GC é editado e versionado aqui e implantado na VM do free5GC por `rsync` e `git apply`. Toda alteração na VM do free5GC é anunciada antes; as destrutivas são confirmadas.

---

## 4. MVPs

### MVP-0 — Biblioteca AKMA comum e validação cruzada

**Estado (16/10/2026):** implementado. Critérios 2 a 6 atendidos; critério 1 atendido com `go test ./...` — a variante `-race` aguarda um compilador C na VM do AF (o detector de corridas exige cgo). **17/10/2026:** com o gcc instalado, `go test -race ./...` passa em `shared/akma`; critério 1 atendido por inteiro. Arquivos conforme a tabela abaixo, com duas diferenças de organização: as derivações ficaram em `shared/akma/derive/` (não em `kdf/akma.go`), para que `derive.KAf` receba um `afid.AfID` sem ciclo de importação; e a comparação com `ueauth` ficou num módulo próprio, `shared/crosscheck/`, para que a biblioteca não dependa de `free5gc/util`.

**Objetivo.** Implementar em Go as três derivações, a conversão de SUPI, a construção e a análise do A-KID e do AF_ID, e as codificações por camada, concordando octeto a octeto com o oráculo Python e com a KDF do free5GC.

**Pré-requisitos.** D1, D2, D3, D7. Estrutura `docs/` por lado criada. Oráculo Python estendido para emitir vetores em JSON.

**Cláusulas.** TS 33.220 B.1, B.2.0, B.2.1.2, B.2.1.3, B.2.2, H.1–H.3 · TS 33.535 A.1–A.4 (A.4 não diz como codificar o FQDN: a leitura por analogia com o `NAF_Id` continua `[INTERPRETAÇÃO]`), §6.1 (A-KID em NAI; RID vindo do UDM), D.3.5 (no CBOR, A-TID como `bstr`) · TS 33.501 A.7.0 · TS 23.003 §2.2 (MNC de 2 ou 3 dígitos conforme o MCC), §2.2B item 3 (RID), §28.2 · TS 29.522 Tab. 5.14.5.4.2-1 (`AKId`, `AfId`) · TS 29.571 (`Supi`) · TS 24.109 K.2 c)2) (codec base64 da K_AF) · RFC 4648, RFC 7542 §2.2.

**Pontos de leitura que o MVP-0 resolve** (seção 7): P9 — `AfId` é definido como `"<FQDN>.<Ua* security protocol id>"` com o identificador descrito como *"a string of 5 octet"*; a forma de dez caracteres hexadecimais vem do exemplo (`0100BC0001`, em maiúsculas). P10 — resolvido: TS 33.220 B.2.0 restringe P0 a *"static ASCII-encoded string"* apenas *"in this specification"*, e TS 33.535 A.4 não adota a restrição; o P0 de K_AF é o AF_ID completo.

**Arquivos.**

| Lado | Caminho | Conteúdo |
|---|---|---|
| shared | `shared/akma/go.mod` | módulo, Go 1.26.2; dependência `golang.org/x/text` (NFKC) |
| shared | `shared/akma/kdf/kdf.go`, `kdf/akma.go` | KDF genérica; `KAkma`, `ATid`, `KAf` |
| shared | `shared/akma/supi/supi.go` | `imsi-<dígitos>` (SBI) → dígitos (P0/P1); recusa NAI/GCI/GLI no MVP |
| shared | `shared/akma/akid/{rid,realm,akid}.go` | RID como cadeia; domínio `5gc.mnc<MNC>.mcc<MCC>.3gppnetwork.org`; montagem e análise do A-KID |
| shared | `shared/akma/afid/{uastar,afid}.go` | identificador Ua\* de 5 octetos; forma de transporte ↔ P0 |
| shared | `shared/akma/codec/key.go` | hex ancorado de 64 caracteres; base64 |
| shared | `shared/akma/testdata/vectors.json` | vetores **gerados pelo projeto**, rotulados |
| shared | `reference/kdf/akma_kdf.py`, `reference/kdf/gen_vectors.py` | oráculo e gerador |
| docs | `docs/shared/{kdf,identificadores,codificacao,vetores}.md` | |
| free5GC | — | nada é tocado |
| AF | — | nada é tocado |

**Estruturas de dados.**

```go
type Key [32]byte                       // K_AKMA, K_AF, A-TID: sempre 32 octetos internamente
type UaStarID [5]byte
type AfID struct { FQDN string; Proto UaStarID }
func (AfID) Transport() string           // "<FQDN>.<10 hex>"
func (AfID) KDFParam() []byte            // NFKC(FQDN) ‖ 5 octetos
func ParseAfID(s string) (AfID, error)
type HomeNetwork struct { MCC, MNC string } // MNC com 2 ou 3 dígitos, nunca inferido do IMSI
type AKID struct { RID string; ATID Key; Realm string }
type Options struct { ATIDEncoding, ATIDLabel string; AcceptLabels []string }  // ADR-0002 §4
func (AKID) Format(Options) string
func ParseAKID(s string, o Options) (AKID, error)
```

**Como testar.** `go test -race ./...` na VM do AF. O teste lê `vectors.json` e compara a cadeia `S` e a saída de cada derivação. Um teste à parte compara a KDF genérica com `ueauth.GetKDFValue` (`free5gc/util` v1.3.2) para os FC `0x80`, `0x81` e `0x82`.

**Critério de aceitação.**
1. `go test -race ./...` sem falhas.
2. Todos os vetores do oráculo reproduzidos octeto a octeto (cadeia `S` e saída), incluindo os três do baseline §2.5.
3. KDF genérica idêntica a `ueauth.GetKDFValue` em todos os vetores.
4. RID `"0"`, `"012"` e `"0000"` sobrevivem a montagem → análise sem alteração; `"012"` ≠ `"12"`.
5. Testes negativos: SUPI com prefixo `imsi-` recusado na camada da KDF; forma de transporte do AF_ID recusada como P0; rótulo `akid` aceito na leitura e nunca emitido; `kAkma` com 63 ou 65 caracteres, ou não hexadecimal, recusado; FQDN com caracteres de compatibilidade Unicode normalizado por NFKC.
6. Todo vetor no repositório rotulado "gerado pelo projeto; não é vetor 3GPP".

**Riscos.** Oráculo e Go escritos a partir das mesmas leituras — mitigado por `ueauth` agora e pelo C++ do UERANSIM no MVP-4. Versão de `golang.org/x/text` precisa ser compatível com o grafo do free5GC. Comprimento do MNC (2 ou 3 dígitos) não pode ser inferido do IMSI.

**O tutorial ganha.** `docs/tutorial/README.md` com pré-requisitos (free5GC v4.2.3, Go 1.26.2) e a seção "Conferir as derivações": como rodar o oráculo e os testes, e como ler um vetor.

---

### MVP-1 — AAnF autônomo na VM do free5GC

**Estado (17/10/2026):** implementado e verificado na VM do free5GC e a partir da VM do AF (`docs/free5gc-side/aanf.md` §9). `go test -race ./...` passa no módulo `free5gc/aanf`. Critérios:

| Nº | Resultado |
|---|---|
| 1 | Parcial no MVP-1 por limitação do NRF (D13, N1): o NRF descartava `aanfInfoList`, que o AAnF envia com `"0000"` preservado (`TestRegisterNFInstanceProfile`), e a descoberta de `AANF` respondia 400 (N2). **Atendido depois do MVP-1b**: perfil guardado com `aanfInfoList` e descoberta por `routing-indicator` funcionando |
| 2 | Atendido; `gpsi` → 400 `OPTIONAL_IE_INCORRECT` (`[INTERPRETAÇÃO]`, `aanf.md` §3.1) |
| 3 | Atendido (`TestReRegistrationReplacesPair`, `TestReRegistrationOverAPI`) |
| 4 | Atendido; K_AF igual ao vetor `testbed-ue` nos testes e na VM |
| 5 | Atendido nos testes e na VM, inclusive o AF_ID autorizado pedido de outra máquina |
| 6 | Atendido |
| 7 | Atendido: 8 códigos testados; 429, 502 e 503 justificados (`aanf.md` §6) |
| 8 | Atendido (conferido em todo teste de API) |
| 9 | Atendido: pedidos de `192.168.16.129` chegam ao AAnF |
| 10 | Atendido (`TestNoKeyInLogs` em nível `trace`; zero ocorrências nos logs da VM) |

Diferenças em relação à tabela de arquivos: `internal/context/{context,akma_context}.go` como previsto, mais `internal/util/problem.go` (`ProblemDetails`); `free5gc/scripts/prepare-openapi.sh` prepara a cópia do `openapi` com os *patches*; a ferramenta de teste ficou em `free5gc/aanf/tools/seed/` com os subcomandos `register` (papel do AUSF), `retrieve` (papel do AF), `remove` e `discover`. Decisão nova: DV-8 (`aud` vazio tolerado), encerrada no MVP-1b.

**Objetivo.** Uma NF AAnF, com a estrutura das NFs do free5GC, que se registra no NRF como `AANF`/`naanf-akma`, expõe as três operações do Naanf_AKMA com contexto em memória e responde aos erros de TS 29.535 e TS 29.500 — testada sem AUSF.

**Pré-requisitos.** MVP-0; D4, D6; decisões DV-1 e DV-7 (seção 6).

**Cláusulas.** TS 33.535 §3.1 (contexto AKMA), §4.1, §4.2.2 (AF autenticado e autorizado), §4.4.0, §5.2, §6.1 passos 4–5 e NOTA 1a, §6.2.1 passos 2, 7 e 8, §6.6, §6.7, §7.1 · TS 29.535 §§4.2.2.2–4.2.2.4, §4.2.1.3.2, §5.1.2.2.2 (`application/json`, `application/problem+json`), 5.1.3, 5.1.4, 5.1.6, §5.1.7.1 (erros de TS 29.500 Tab. 5.2.7.2-1; TS 29.501 §4.8), 5.1.7.3, 5.1.8, 5.1.9, Anexo A (YAML v1.2.0) · TS 29.522 §5.14.5.3–5.14.5.4, §5.14.7.3, Anexo A.12 · TS 29.510 (`NFProfile`, `NFService`, Tab. 6.1.6.2.73-1 `AanfInfo`: sem RIDs, o AAnF atende qualquer RID; chave do `aanfInfoList` com até 32 caracteres) · TS 29.500 §5.2.7.1 (Tab. 5.2.7.1-1), §5.2.7.2 (método não suportado → 501), §5.2.3.2.1 (cabeçalhos obrigatórios) · TS 29.501 §4.4.2, §6.2 (corpo JSON até 16 milhões de octetos) · TS 33.501 §13.1.0 (suporte a TLS), §13.4.1.0 (OAuth2 obrigatório de suportar; escopo em nível de serviço, escopos adicionais opcionais).

**Arquivos.**

| Lado | Caminho | Conteúdo |
|---|---|---|
| free5GC | `free5gc/patches/openapi/0001-naanf-akma-models.patch` | modelos do Naanf_AKMA sobre `openapi` v1.2.4 |
| free5GC | `free5gc/patches/openapi/0002-naanf-akma-client.patch` | cliente `aanf/AKMA` (usado pelo AUSF e pelo AF) |
| free5GC | `free5gc/aanf/cmd/main.go` | |
| free5GC | `free5gc/aanf/pkg/factory/{config,factory}.go` | leitura de `aanfcfg.yaml` |
| free5GC | `free5gc/aanf/pkg/service/init.go` | ciclo de vida |
| free5GC | `free5gc/aanf/internal/context/{context,akma_context}.go` | contexto da NF; armazenamento AKMA |
| free5GC | `free5gc/aanf/internal/sbi/{server,api_akma}.go` | rotas, validação, ProblemDetails, OAuth2 |
| free5GC | `free5gc/aanf/internal/sbi/processor/{register_anchorkey,retrieve_applicationkey,remove_context}.go` | |
| free5GC | `free5gc/aanf/internal/sbi/consumer/nrf_service.go` | registro, *heartbeat*, cancelamento |
| free5GC | `free5gc/aanf/tools/seed/` | ferramenta que faz o papel do AUSF, para testes |
| free5GC | `free5gc/config/aanfcfg.yaml` | |
| free5GC | `free5gc/scripts/{deploy,build-aanf,run-aanf}.sh` | |
| docs | `docs/free5gc-side/{aanf,openapi,registro-nrf}.md`, `docs/conformidade/desvios.md` | |
| VM free5GC | `~/free5gc/` | **nenhum arquivo do free5GC é alterado**; o binário vai para `~/free5gc/bin/aanf` e a configuração para `~/free5gc/config/aanfcfg.yaml` (arquivos novos) |

**Estruturas de dados.**

```go
type UeAkmaContext struct {
    Supi   string               // forma SBI, "imsi-..."; chave primária
    AKID   string               // A-KID corrente
    KAkma  akma.Key
    Gpsi   string               // vazio no MVP
    AfKeys map[string]*AfKeyRecord   // por AF_ID de transporte
}
type AfKeyRecord struct {
    AfID      afid.AfID
    KAf       akma.Key
    Expiry    time.Time          // instante absoluto; serializado em UTC
    AKIDSent  string             // §6.8.1 NOTA 1
    NotifURI  string             // RoamingRestriction; não usado no MVP
    NotifExp  time.Time
}
type Store struct {             // protegido por mutex
    bySupi map[string]*UeAkmaContext
    byAKID map[string]string     // A-KID → SUPI
}
```

Novo par (A-KID, K_AKMA) para um SUPI apaga o par anterior e o índice do A-KID antigo; os `AfKeyRecord` permanecem até `Expiry` (§6.1 NOTA 1a, §5.2).

Configuração (`aanfcfg.yaml`), além dos campos comuns às NFs do free5GC: `routingIndicators` (`AanfInfo`, cadeias), `akma.kafLifetime` (política do operador, §6.2.1 passo 7), `akma.authorizedAfs` (política local, §6.2.1 passo 2: cada AF_ID autorizado associado à identidade do AF que pode pedi-lo — no MVP, o endereço de origem; com TLS, a identidade do certificado), `akma.codec` (ADR-0002 §4), `akma.roamingStatusSubscription: false` (questão 1).

**Como testar.** Unitários do `Store` (substituição, expiração, índices). Testes de API com servidor h2c em memória para cada operação e cada código. Integração na VM do free5GC: NRF + AAnF, consulta do perfil, descoberta `target-nf-type=AANF&requester-nf-type=AUSF`, token OAuth2 de escopo `naanf-akma`. Da VM do AF: requisição à porta do AAnF.

**Critério de aceitação.**
1. Perfil no NRF com `nfType: AANF`, serviço `naanf-akma` versão `v1`, `ipEndPoints` em `192.168.16.130`, `aanfInfoList` com os RIDs configurados como cadeias (`"0000"` preservado).
2. `register-anchorkey` aceita `AkmaKeyInfo` válido e responde 200 com `AkmaKeyInfo` no corpo (TS 29.535 §4.2.2.2.2); recusa com 400 `supi` e `gpsi` juntos, nenhum dos dois, `gpsi` sem a *feature*, `kAkma` fora do padrão.
3. Novo registro para o mesmo SUPI torna o A-KID antigo desconhecido e mantém as K_AF já emitidas até a expiração.
4. `retrieve-applicationkey` devolve K_AF idêntica ao vetor do oráculo para (K_AKMA, AF_ID), `expiry` em RFC 3339 UTC no futuro, `supi` presente, `gpsi` ausente.
5. A-KID sem K_AKMA → 403 com `cause` `K_AKMA_NOT_PRESENT` (TS 29.535 §4.2.2.3.2). AF_ID fora da política local → 403 com `ProblemDetails` sem causa específica (P3). AF_ID autorizado pedido por AF a que ele não pertence → 403 (ataque de falsificação de AF; ver `docs/literatura/`). 204 conforme DV-7.
6. `remove-context` responde 204 e torna o contexto inexistente; SUPI desconhecido → 404 `AKMA_CONTEXT_NOT_FOUND`.
7. Matriz dos 11 códigos obrigatórios para POST (TS 29.500 Tab. 5.2.7.1-1): cada código ou tem teste, ou tem justificativa escrita de por que só é alcançável por infraestrutura ausente no MVP (por exemplo, 502 via SCP).
8. `ROAMING_AKMA_SERVICE_DENIED` nunca emitido.
9. Requisição feita a partir de `192.168.16.129` chega ao AAnF (ADR-0001-R2 §5.2).
10. Nenhuma chave aparece em log de nível `info` ou inferior em severidade.

**Riscos.** O *middleware* de OAuth2 das NFs do free5GC é por serviço; `retrieve-applicationkey` vindo do AF não terá token (o AF não alcança o NRF) e precisa ser autorizado por política local, operação por operação. O NRF do free5GC emite tokens com escopo do nome do serviço. `[3GPP]` TS 33.501 §13.4.1.0 define esse nível de serviço como a extensão básica do token e torna os escopos adicionais dependentes de configuração do produtor; os escopos finos de TS 29.535 (`naanf-akma:anchorkey` etc.) ficam opcionais no MVP. O NRF ignora `routing-indicator` para `AANF` (correto com uma instância). Colisão de porta em `192.168.16.130`.

**O tutorial ganha.** "Instalar o AAnF no free5GC": aplicar o *patch* do `openapi`, compilar, configurar, iniciar, verificar o registro no NRF.

---

### MVP-1b — Correção do NRF do free5GC

**Estado (17/10/2026):** implementado e verificado na VM do free5GC (`docs/free5gc-side/nrf.md`). Decisão D13.

**Objetivo.** Fazer o NRF do free5GC v4.2.3 (`nrf` v1.4.5) cumprir TS 29.510 nos três pontos encontrados no MVP-1, com um *patch* por ponto, cada um com teste e utilizável isoladamente *upstream*.

**Cláusulas.** TS 29.510 Tab. 6.1.6.2.2-1 (`aanfInfoList` no `NFProfile`), §5.2.2.2.2 (o NRF guarda o perfil, inclusive atributos que não conhece), Tab. 6.1.6.2.73-1 (`AanfInfo`: sem `routingIndicators`, qualquer RID), Tab. 6.2.3.2.3.1-1 (`routing-indicator` para `AUSF`, `AANF` e `UDM`), Tab. 6.3.5.2.4-1 (`aud` obrigatório: ID da instância ou tipo de NF do produtor) · TS 33.501 §13.4.1.1.2 · TS 33.535 §6.7.

**Arquivos.**

| Lado | Caminho | Conteúdo |
|---|---|---|
| free5GC | `free5gc/patches/nrf/0001-store-aanf-info-list.patch` | N1: `nnrfNFManagementOption` guarda `aanfInfoList` |
| free5GC | `free5gc/patches/nrf/0002-discover-aanf-by-routing-indicator.patch` | N2: `AANF` aceito na descoberta; filtro de `routing-indicator` sobre `aanfInfoList` |
| free5GC | `free5gc/patches/nrf/0003-access-token-audience.patch` | N3: `aud` = `targetNfInstanceId` ou `targetNfType` |
| free5GC | `free5gc/scripts/patch-nf.sh` | aplica ou reverte *patches* numa NF e recompila (reutilizado no MVP-2) |
| free5GC | `free5gc/aanf/internal/context/context.go` | AAnF passa a exigir `aud` (DV-8 encerrado) |
| VM free5GC | `~/free5gc/NFs/nrf` (três arquivos alterados, três de teste novos), `~/free5gc/bin/nrf` (original guardado em `bin/nrf.orig`) | |

**Critério de aceitação.**
1. Os testes do NRF passam com os *patches*, inclusive os três novos. **Atendido**, na VM do AF e na VM do free5GC.
2. O perfil do AAnF guardado no NRF contém `aanfInfoList` com `"0000"`. **Atendido**; com isso, o critério 1 do MVP-1 passa a atendido.
3. A descoberta `target-nf-type=AANF&routing-indicator=0000` devolve o AAnF; `0001` e `000` não. **Atendido.** A semântica do filtro foi conferida no MongoDB 5.0 com quatro perfis sintéticos: sem lista, com RID, com entrada sem RIDs e com duas entradas.
4. O token emitido para `AANF` traz `aud`, e o AAnF com verificação estrita aceita `register-anchorkey`. **Atendido.**
5. Regressão: as NFs do `run.sh` que não exigem `sudo` (`amf smf udr pcf udm nssf ausf chf nef`) se registram no NRF corrigido sem erro. **Atendido.** Fica para o responsável a regressão completa com UPF, gNB e UE, que pede `sudo`.

**Fora do escopo.** O caminho `complexQuery` da descoberta continua sem tratar `routing-indicator` para `AANF`. O NEF do free5GC pede o token de `nnrf-nfm` com `targetNfType=NEF` e não consegue se descadastrar; a falha é anterior ao projeto (log de 15/10/2026) e foi apenas registrada.

---

### MVP-2 — UDM e AUSF: indicação AKMA e AnchorKey_Register

**Estado (17/10/2026):** implementado e verificado na VM do free5GC (`docs/free5gc-side/{udm,ausf,assinatura-akma,verificacao}.md`).

| Critério | Resultado |
|---|---|
| 1 | Atendido no banco de testes: `TestAkma5GAKA` (RID `"0000"`) e `TestAkmaEAPAKAPrime` (RID `"012"`) |
| 2 | Atendido com UERANSIM: `rid0000.atidfTiRYw4XQKRmE8qEUX9JDLIbElNilKMynonfqfxgnjI=@5gc.mnc093.mcc208.3gppnetwork.org` |
| 3 | Atendido: `TestAkmaWithoutSubscription`; e `TestGenerateAuthDataWithoutAkma` no UDM |
| 4 | Atendido com UERANSIM: o reinício do UE trouxe A-KID novo, e o AAnF registrou `replaced previous: true` |
| 5 | Atendido: os quatro testes passam pelo `test.sh` |
| 6 | Atendido com UERANSIM e o AAnF parado: registro aceito e falha só no log do AUSF |
| 7 | Atendido: nenhuma cadeia de 64 hexadecimais nos logs |

Diferenças em relação à tabela de arquivos e decisões tomadas na implementação:

- **UDM.** A lógica ficou num arquivo novo, `internal/sbi/processor/akma.go`. A extração do RID do SUCI virou `suci.RoutingIndicator`, usada só na comparação de D10: o RID enviado é sempre o do UDR. Com `akmaAllowed` sem `routingId` válido, o UDM não envia a indicação e registra erro (TS 29.505 exige o `routingId`).
- **AUSF.** O domínio de origem vem do SUCI. Com SUPI, vem da PLMN única de `plmnSupportList` que prefixa o IMSI (`[INTERPRETAÇÃO]`, TS 23.003 §2.2). O contexto guarda `SupiOrSuci` em vez de `HomeNetwork`.
- **Configuração.** O `ausfcfg.yaml` é alterado por *patch* (`free5gc/patches/config/0001-ausf-akma.patch`).
- **Banco de testes.** O *patch* também muda `nfs_setup.go` (AAnF no processo, `akma.enable` no AUSF), `ngsetup_test.go` (sem SMF quando só os testes `TestAkma*` rodam) e `go.mod`. Os testes AKMA rodam sem `sudo`; os do critério 5 continuam pelo `test.sh`.
- **`openapi` do banco de testes.** O banco fixa `openapi` v1.2.5 (pré-lançamento). `prepare-openapi.sh` ganhou o argumento de versão e gera `free5gc/build/openapi-test`.
- **`patch-nf.sh`.** Passa a cobrir os alvos `test` e `config`, aplicados a partir da raiz do repositório do free5GC. Arquivos não rastreados não bloqueiam a aplicação.
- **Local do repositório.** Os `replace` do AUSF e do banco de testes apontam para `~/akma-free5gc`, ao lado de `~/free5gc`.

**Objetivo.** Depois de uma autenticação primária bem-sucedida (5G-AKA e EAP-AKA') de um assinante com assinatura AKMA, o UDM indica `akmaInd` e `routingId`, e o AUSF deriva K_AKMA, A-TID e A-KID e os registra no AAnF descoberto via NRF.

**Pré-requisitos.** MVP-1 e MVP-1b; D8 (TS 29.505, obtida); D9 e D10 decididas com o texto da norma (decididas em 16 e 17/10/2026); aviso antes de aplicar *patches* nas árvores de `NFs/ausf` e `NFs/udm` da VM do free5GC (com `free5gc/scripts/patch-nf.sh`).

**Cláusulas.** TS 33.535 §3.1 e §4.2.5 (assinatura AKMA), §6.1 passos 1–5, NOTAS 1 e 1a, *"The AUSF shall use the RID received from the UDM"*, §6.7 e NOTA 1 · TS 33.501 §6.1.3.1 passo 2 e §6.1.3.2.0 (o UDM *shall* incluir a indicação AKMA e o RID na resposta de `Nudm_UEAuthentication_Get`), §6.1.1.1 (K_AUSF pode ser guardada no AUSF quando AKMA é suportado), §6.2.2.1 (K_AUSF guardada ou substituída quando RES\* = XRES\*, passo 11 do 5G-AKA, ou quando o AUSF envia EAP-Success, passo 10 do EAP-AKA' — os dois pontos de acionamento do AKMA no AUSF), §14.2.2, A.7.0 · TS 29.503 Tab. 6.3.6.2.3-1, §5.4.2.2.2 · TS 29.509 Tab. 6.1.6.2.2-1 (o `routingIndicator` vindo do AMF não serve ao A-KID) · TS 29.505 (`akmaAllowed`) · TS 29.510 (descoberta) · TS 29.535 §4.2.2.2.

**Arquivos.**

| Lado | Caminho no free5GC (VM) | Alteração | *Patch* no repositório |
|---|---|---|---|
| free5GC | `NFs/udm/pkg/suci/suci.go` | exportar a extração do RID do SUCI | `free5gc/patches/udm/` |
| free5GC | `NFs/udm/internal/sbi/processor/generate_auth_data.go` | se `authSubs.AkmaAllowed`, preencher `AkmaInd` e `RoutingId` (cadeia) | idem |
| free5GC | `NFs/ausf/internal/context/context.go` | `AusfUeContext`: `AkmaInd`, `RoutingId`, PLMN doméstica do SUCI | `free5gc/patches/ausf/` |
| free5GC | `NFs/ausf/internal/sbi/processor/ue_authentication.go` | guardar a indicação em `UeAuthPostRequestProcedure`; acionar o AKMA no sucesso de `Auth5gAkaComfirmRequestProcedure` e de `EapAuthComfirmRequestProcedure` | idem |
| free5GC | `NFs/ausf/internal/sbi/processor/akma.go` (novo) | derivação e registro | idem |
| free5GC | `NFs/ausf/internal/sbi/consumer/aanf_service.go` (novo) | descoberta do AAnF com `routing-indicator`, token, `register-anchorkey`; URI local como alternativa (§6.7) | idem |
| free5GC | `NFs/ausf/pkg/factory/config.go`, `config/ausfcfg.yaml` | seção `akma` | idem |
| free5GC | `NFs/ausf/go.mod` | `replace` do `openapi` corrigido e de `shared/akma` | idem |
| free5GC | `test/ranUe.go`, `test/akma_test.go` (novo) | guardar K_AUSF do UE de teste; testes AKMA em 5G-AKA e EAP-AKA' | `free5gc/patches/test/` |
| free5GC | MongoDB, `authenticationSubscription` do assinante | `akmaAllowed: true` | `free5gc/scripts/provision-akma-subscriber.sh` |
| docs | — | — | `docs/free5gc-side/{udm,ausf,assinatura-akma,verificacao}.md` |
| AF | — | nada | — |

**Estruturas de dados.** `AusfUeContext` acrescido de `AkmaInd bool`, `RoutingId string`, `HomeNetwork akid.HomeNetwork`. Configuração `akma: { enable: bool, aanfUri: string }` no AUSF, com `aanfUri` vazio significando descoberta via NRF.

**Como testar.** Unitários do UDM (RID `"0"`, `"012"`, `"0000"` extraídos do SUCI sem perda). Unitário do AUSF com UDM e AAnF simulados: o `AkmaKeyInfo` enviado tem `supi` com prefixo (forma SBI), `aKId` e `kAkma` iguais aos derivados pelo oráculo a partir da mesma K_AUSF. Integração pelo banco de testes Go do free5GC: o UE de teste se registra; o teste deriva A-KID e K_AF a partir da própria K_AUSF e pede a K_AF ao AAnF. Integração com UERANSIM (5G-AKA).

**Critério de aceitação.**
1. No banco de testes, para 5G-AKA e para EAP-AKA', a K_AF devolvida pelo AAnF é igual à K_AF derivada pelo UE de teste a partir da sua K_AUSF.
2. Com UERANSIM (`imsi-208930000000001`, RID `"0000"`), o AAnF guarda o A-KID `rid0000.atid<44 caracteres base64>@5gc.mnc093.mcc208.3gppnetwork.org`.
3. Assinante sem `akmaAllowed`: `akmaInd` ausente e nenhuma chamada ao AAnF.
4. Reautenticação substitui o contexto no AAnF.
5. `TestRegistration`, `TestGUTIRegistration`, `TestReSynchronization` e `TestEAPAKAPrimeAuthentication` continuam passando.
6. Comportamento com AAnF indisponível conforme D9.
7. Nenhuma chave em log.

**Riscos.** RID quando não há SUCI (D10). Comprimento do MNC do domínio doméstico: tomado do SUCI, nunca do IMSI. A *webconsole* pode sobrescrever `akmaAllowed` ao editar o assinante — provisionar por `mongosh` e documentar. A política de tokens do NRF pode exigir ajuste para AUSF → AANF. O banco de testes do free5GC exige `sudo` (usa *namespaces*). Os *patches* deixam os submódulos `NFs/ausf` e `NFs/udm` com alterações locais — o usuário decide como versioná-las lá.

**O tutorial ganha.** "Habilitar AKMA no UDM e no AUSF" e "Provisionar um assinante AKMA", com a verificação pelo banco de testes.

---

### MVP-3a — Mensagens AKMA na tela do free5GC e do UERANSIM

**Estado (18/10/2026):** implementado; verificação na VM do free5GC em `docs/free5gc-side/verificacao.md`, seção 4.

**Objetivo.** Pedido do responsável (D14): acompanhar o AKMA pelas telas do `run.sh` e do UE, sem abrir arquivos de log, com as mensagens principais no formato e nos níveis de cada ferramenta.

**Pré-requisitos.** MVP-2; aviso antes de aplicar *patches* no UERANSIM, no `run.sh` e no `force_kill.sh` da VM do free5GC.

**Cláusulas.** Nenhuma cláusula 3GPP trata de log. A derivação antecipada no ME segue TS 33.535 §6.1 (*"before initiating communication with an AKMA Application Function"*), A.2, A.3 · TS 33.501 A.7.0 · TS 23.003 §2.2B, §28.2 · TS 29.522 Tab. 5.14.5.4.2-1.

**Arquivos.**

| Lado | Caminho | Conteúdo |
|---|---|---|
| free5GC | `free5gc/patches/udm/0001-*.patch` (regenerado) | categoria `AKMA`; indicação enviada em `info` |
| free5GC | `free5gc/patches/ausf/0001-*.patch` (regenerado) | indicação recebida, derivação com A-KID, AAnF selecionado, registro, AKMA desligado — em `info` |
| free5GC | `free5gc/aanf/` | identificadores em `info` (D16); linha de partida com RIDs, tempo de vida da K_AF e AF_ID autorizados |
| free5GC | `free5gc/patches/run/0001-aanf-option.patch` (novo) | `run.sh -aanf` e `aanf` no `force_kill.sh` (D15) |
| free5GC | `free5gc/scripts/{patch-nf,run-aanf}.sh` | alvo `run`; checagem só nos arquivos tocados; `GIN_MODE=release` |
| AF/UE | `af/ue/ueransim/patches/0001-akma-key-derivation.patch` | K_AKMA, A-TID e A-KID no ME; `[nas] [info]` com o A-KID |
| AF/UE | `af/ue/ueransim/patches/0002-free5gc-ue-akma.patch` | `akma: true` no `free5gc-ue.yaml` |
| AF/UE | `af/scripts/patch-ueransim.sh` | aplica, reverte e recompila só o `nr-ue` |
| docs | `docs/free5gc-side/logs-akma.md`, `docs/af-side/ue-ueransim.md` | catálogo, formato, verificação |

**Critério de aceitação.**
1. Com `./run.sh -aanf`, o AAnF sobe com as outras NFs, escreve no mesmo terminal e no mesmo `free5gc.log`, e para com elas no Ctrl-C.
2. No registro de um UE com assinatura AKMA, a tela do `run.sh` mostra, em `info` e na categoria `AKMA`: a indicação enviada pelo UDM; no AUSF, a indicação recebida, a derivação, o AAnF selecionado e o registro; no AAnF, o contexto registrado.
3. A tela do UE mostra `AKMA key material derived, A-KID[…]`, e esse A-KID é idêntico, caractere a caractere, ao do AUSF e ao do AAnF.
4. Sem `akma` na configuração do UE, o `nr-ue` se comporta como o original: nenhuma linha AKMA.
5. Nenhuma chave nas telas nem nos logs.
6. Os testes do UDM, do AUSF, do AAnF e os `TestAkma*` do banco de testes continuam passando.

**Riscos.** Duas instâncias do AAnF (pelo `run.sh -aanf` e pelo `run-aanf.sh`) disputam a mesma porta: documentado. A derivação no UE sem a conferência da K_AKMA, que só vem no MVP-4 (`ue-ueransim.md` §8).

**O tutorial ganha.** "Ver o AKMA na tela": `run.sh -aanf`, o UE com `akma: true` e o que procurar em cada terminal.

---

### MVP-3 — AF mínimo e Ua\* simplificado

**Estado (18/10/2026):** implementado e verificado com o AAnF da VM do free5GC. A especificação `docs/af-side/ua-star.md` foi aprovada pelo responsável antes do código; documentação em `docs/af-side/{af,teste-af}.md`.

| Critério | Resultado |
|---|---|
| 1 | Atendido: o cliente, com a K_AF do oráculo para o contexto semeado (`testbed-ue`), autentica no AF, que usou a K_AF recebida do AAnF |
| 2 | Atendido: `TestApplicationKey`; o AAnF registrou `AF_ID af1.akma.example.FF00000001` e o A-KID sem alteração |
| 3 | Atendido: K_AF certa aceita em três requisições; K_AF errada recebe `401` |
| 4 | Atendido: `rspauth` conferido pelo cliente; AF sem a K_AF recusado (`TestAFWithoutKaf`) |
| 5 | Atendido nos testes automáticos, com relógio controlado, nas políticas `refetch` e `reject` |
| 6 | Atendido: pedido sem `Authorization` recebe o desafio |
| 7 | Atendido: nenhuma chave nos logs do AF e do AAnF |

Diferenças em relação à tabela de arquivos:

- **Um módulo Go em `af/`**, e não um por programa. O Digest ficou em `af/uastar/`, compartilhado pelo AF (`af/af/`) e pelo cliente (`af/ue/uaclient/`), e não em `af/af/internal/uastar/`.
- **Cache da K_AF.** Fica em `af/af/internal/server/keys.go`, junto com a política de expiração (`kafExpiry`), e não num pacote `keycache`.
- **Configuração.** O `afcfg.yaml` fica em `af/config/`.
- **Vetores do Digest.** `af/uastar/testdata/gen_digest_vectors.py` é uma implementação Python independente; o `reference/` continua só com o oráculo da KDF (D1).
- **Cliente.** No MVP-3, recebe o A-KID e a K_AF de um vetor de `shared/akma` ou da linha de comando. O AAnF foi semeado com o `seed` do MVP-1.

**Objetivo.** Um AF na VM do AF que serve um Ua\* simplificado (formato do Digest de TS 24.109 K.2, identificador `FF 00 00 00 01`), obtém K_AF do AAnF por `AkmaAfKeyRequest`, a guarda até `expiry` e autentica um cliente que prova possuir K_AF.

**Pré-requisitos.** MVP-1 (AAnF alcançável da VM do AF); decisão DV-5 (seção 6); especificação do Ua\* simplificado escrita em `docs/af-side/ua-star.md` e aprovada antes do código; Go na VM do AF.

**Cláusulas.** TS 33.535 §4.4.1 (Ua\* carrega o A-KID e trata a expiração de K_AF), §6.2.1 passos 1, 2, 8 e 9, §6.4.2, §6.4.3, §6.5 (iniciação do AKMA pelo AF), A.4, B.1.2.2 (Stage 2 do perfil Digest) · TS 33.220 §4.5.2 (pré-requisitos de consistência do FQDN na derivação), H.2 NOTA 2, H.3 · TS 24.109 §5.2.1 (*product token*), §5.2.2.1 (Digest conforme RFC 7616; `username`; senha em base64; `realm` em duas partes), §5.2.6 (`qop=auth-int` opcional), Anexo K.2 · TS 33.222 §5.3.0 passo 2 (cabeçalho `Host` com o nome do AF) e NOTA 1b (MD5 desaconselhado) · TS 29.522 §5.14.5.3 · TS 29.535 §4.2.2.3.

**O que a leitura fixou.** `[3GPP]` O Ua\* simplificado é, na forma, o Digest de TS 24.109 §5.2 com os parâmetros AKMA de K.2 e de TS 33.535 B.1.2.2 — `User-Agent` com `3gpp-akma`, `username` = A-KID, senha = base64(K_AF), `realm` com o FQDN do AF como segunda parte e a primeira parte conforme DV-6 — sem o túnel TLS que K.2 a) exige. `[INTERPRETAÇÃO]` Por faltar o TLS, não é o perfil K.2; por isso mantém o identificador privado `FF 00 00 00 01` (ADR-0003). O identificador registrado `01 00 00 00 02` (TS 33.220 H.3, "HTTP digest authentication according to TS 24.109") existe, mas está definido para o GBA; usá-lo para chaves AKMA seria afirmar um perfil que TS 24.109 não especifica para o AKMA. O ADR-0003 precisa registrar essa alternativa e o motivo de não adotá-la (I15).

**Arquivos.**

| Lado | Caminho | Conteúdo |
|---|---|---|
| AF | `af/af/cmd/af/main.go` | |
| AF | `af/af/internal/config/` | FQDN, identificador Ua\*, URI do AAnF (configuração local, §6.7), endereço de escuta |
| AF | `af/af/internal/naanf/client.go` | h2c; usa o cliente do D6 |
| AF | `af/af/internal/keycache/cache.go` | K_AF por A-KID, até `expiry` |
| AF | `af/af/internal/uastar/{digest,server}.go` | desafio, verificação, `rspauth` |
| AF | `af/ue/uaclient/` | cliente Ua\*; no MVP-3 recebe A-KID e K_AF de vetores |
| free5GC | `free5gc/aanf/tools/seed/` (do MVP-1) | semeia o contexto no AAnF, na VM do free5GC |
| docs | `docs/af-side/{af,ua-star,teste-af}.md` | |

**Estruturas de dados.** `KafEntry{AKID string; KAf akma.Key; Expiry time.Time; Supi string}`; estado de desafio `{nonce, opaque, criado}`.

**Como testar.** Unitários do Digest contra os exemplos de RFC 7616 (algoritmo conforme DV-5) e com senha `base64(K_AF)`. Integração com o AAnF semeado: cliente com a K_AF correta autentica; com K_AF errada, não; A-KID desconhecido; AF fora da lista local; K_AF expirada leva a nova requisição ao AAnF, que devolve a mesma K_AF enquanto K_AKMA não mudar (§5.2, §6.4.3).

**Critério de aceitação.**
1. A K_AF obtida pelo AF é igual à do oráculo para o contexto semeado.
2. O `afId` enviado ao AAnF está na forma de transporte e o `aKId` é exatamente o recebido do cliente.
3. Autenticação Ua\* bem-sucedida com a K_AF correta e recusada com qualquer outra.
4. Confirmação de chave nos dois sentidos: o AF prova ao cliente que tem a K_AF (`rspauth` no `Authentication-Info`), e o cliente recusa resposta sem prova válida.
5. O AF nunca usa K_AF após `expiry` (teste com tempo de vida curto); com K_AF expirada, recusa ou pede nova ao AAnF, conforme TS 33.535 §6.4.2.
6. Pedido sem parâmetros AKMA recebe a mensagem de iniciação do AKMA (o desafio `WWW-Authenticate`), conforme TS 33.535 §6.5.
7. Nenhuma chave em log.

**Riscos.** Sem TLS, o A-KID trafega em claro e o AF não é autenticado perante o UE — limitação do perfil privado, declarada em `ua-star.md`. O AF precisa usar, na derivação, o mesmo FQDN que o UE usou (cabeçalho `Host` e segunda parte do `realm`); um AF alcançável por mais de um nome produz K_AF divergente (TS 33.220 §4.5.2). O Go não tem servidor Digest na biblioteca padrão; a implementação de RFC 7616 é do projeto e precisa de vetores da própria RFC.

**O tutorial ganha.** "Escrever um AF que usa AKMA": configurar o AF, autorizá-lo no AAnF, obter K_AF.

---

### MVP-4 — Lado UE: derivação no ME e envio do A-KID

**Estado (18/10/2026):** implementado e verificado na VM do free5GC, com o UE no próprio *namespace* e o cliente Ua\* indo ao AF pelo plano de usuário. Documentação em `docs/af-side/{ue-ueransim,ua-client,teste-ue}.md` e `docs/free5gc-side/plano-de-usuario.md`.

| Critério | Resultado |
|---|---|
| 1 | Atendido: o A-KID do UE é igual ao do AUSF, ao do AAnF e ao do AF, em duas autenticações |
| 2 | Atendido: três sessões Ua\* com a K_AF do ME, a partir de 10.60.0.1, com `rspauth` conferido |
| 3 | Atendido: `nr-akma-selftest` com 5 vetores conferidos e 3 de NFKC recusados na K_AF (DV-4), com K_AKMA e A-KID conferidos |
| 4 | Atendido: depois de `deregister normal` e novo registro, A-KID novo aceito; o antigo recebe `403 K_AKMA_NOT_PRESENT` do AAnF |
| 5 | Atendido: nenhum comando do `nr-cli` mostra a K_AKMA |

Diferenças em relação à tabela de arquivos e decisões tomadas na implementação:

- **Antecipação (D14).** K_AKMA, A-TID, A-KID e a opção `akma` entraram no MVP-3a (*patches* 0001 e 0002). O MVP-4 acrescentou o 0003 (K_AF, comando `akma-kaf` do `nr-cli`, `nr-akma-selftest`) e o 0004 (`gnbSearchList` no `veth`).
- **Um comando só.** O `akma-kaf <af-fqdn> [<id Ua*>]` devolve A-KID, AF_ID e K_AF.
- **Autoteste.** É um executável separado (`nr-akma-selftest`), que lê o `vectors.json` com o *parser* YAML do UERANSIM.
- **Cliente.** Ele mesmo chama o `nr-cli` (`-nr-cli`, `-ue`), e a K_AF nunca passa por linha de comando. `-bind` põe o endereço da sessão como origem.
- **Script.** `af/scripts/ue-netns.sh` (`up`, `ue`, `cli`, `ua`), no lugar de `ue-ua.sh`, também cria o *namespace* do UE.
- **Plano de usuário.** As falhas desde o MVP-1 tinham duas causas:
  - o UE rodava no *namespace* do host, onde o endereço da sessão é local e o kernel descarta os pacotes que saem da `upfgtp`;
  - faltava o NAT do N6.

  Correções: UE no *namespace* `ueransim-ue`, com `veth` e `linkIp: 10.100.0.1` no gNB (mudança manual), e `free5gc/scripts/n6-nat.sh`. Nada disso é persistente.
- **`patch-ueransim.sh`.** Confere a série numa cópia temporária, porque o 0003 depende do 0001, e compila um alvo por vez com o `cmake` 3.31 que configurou a árvore.

**Objetivo.** O `nr-ue` deriva K_AKMA, A-TID e A-KID após a autenticação primária e K_AF sob demanda; o cliente Ua\* usa o A-KID e a K_AF para autenticar-se no AF, e as chaves do UE e da rede coincidem.

**Antecipado no MVP-3a (D14).** A derivação de K_AKMA, A-TID e A-KID, a opção `akma` e o log já estão em `af/ue/ueransim/patches/0001-akma-key-derivation.patch` e `0002-free5gc-ue-akma.patch`. O MVP-4 acrescenta um `0003` com a K_AF, o comando do `nr-cli` e o autoteste com os vetores. O critério 1 (A-KID do UE igual ao do AAnF) já é verificável pela tela.

**Pré-requisitos.** MVP-0, MVP-3; D5; compilação do UERANSIM no local escolhido.

**Cláusulas.** TS 33.535 §6.1 (derivação pelo ME antes de comunicar com um AF), §6.2.1 passo 1, A.2–A.4 · TS 33.501 A.7.0 · TS 23.003 §2.2B, §28.2 · TS 24.109 K.2 (*product token* `3gpp-akma` no `User-Agent`).

**Arquivos (com D5 = A).**

| Lado | Caminho | Conteúdo |
|---|---|---|
| AF/UE | `af/ue/ueransim/patches/0001-akma-me.patch` sobre `src/ue/nas/keys.{hpp,cpp}`, `src/ue/types.hpp`, `src/ue/nas/mm/auth.cpp` ou `security.cpp`, `src/lib/app/cli_cmd.{hpp,cpp}`, `src/ue/app/cmd_handler.cpp`, leitura de configuração | derivações no ME; comando de `nr-cli` que devolve o A-KID e a K_AF para um AF_ID; opção `akma` na configuração do UE |
| AF/UE | `af/ue/uaclient/` | obtém (A-KID, K_AF) do ME pelo `nr-cli` e executa o Ua\* |
| AF/UE | `af/scripts/ue-ua.sh` | executa o cliente no *namespace* do UE |
| docs | `docs/af-side/{ue-ueransim,ua-client}.md` | |

**Estruturas de dados.** `struct AkmaContext { OctetString kAkma; OctetString aTid; std::string aKid; bool valid; }` no contexto de segurança NAS do UE.

**Como testar.** As derivações em C++ reproduzem os vetores do MVP-0 (comando de autoteste). Integração: o UE se registra; o A-KID devolvido pelo `nr-cli` é igual ao guardado no AAnF; o cliente Ua\* autentica no AF.

**Critério de aceitação.**
1. A-KID do UE igual, caractere a caractere, ao do AAnF.
2. Autenticação Ua\* bem-sucedida com a K_AF derivada pelo UE.
3. As derivações em C++ reproduzem todos os vetores do MVP-0.
4. Após novo registro com nova autenticação, o A-KID do UE muda e o antigo deixa de ser aceito pelo AAnF.
5. K_AKMA não é exposta pelo `nr-cli`.

**Riscos.** O comando do `nr-cli` é interface de teste, não interface 3GPP: entrega K_AF a qualquer processo local que o chame, enquanto TS 33.535 §4.4.3 exige que cada aplicação só obtenha K_AF dos AF_ID a que está autorizada (como garantir isso está fora do escopo 3GPP). Limitação a declarar no tutorial. NFKC indisponível no C++ do UERANSIM: FQDN restrito a ASCII no lado UE. Instante da derivação (após o *Security Mode Command* no 5G-AKA). *Toolchain* C++ no Ubuntu 18.04 (já resolvida na instalação atual).

**O tutorial ganha.** "Preparar o UE (UERANSIM) para AKMA" e "Executar o cliente Ua\*".

---

### MVP-5 — Ponta a ponta e testes negativos (primeiro MVP funcional)

**Estado (19/10/2026):** implementado e executado no laboratório (`docs/af-side/e2e.md`): `af/scripts/e2e.sh` terminou com código 0, 63 verificações e nenhuma falha, em 18/10/2026 na VM do free5GC e em 19/10/2026 numa instalação limpa (critério 4).

| Critério | Resultado |
|---|---|
| 1 | Atendido. O relatório mostra quatro coisas: o registro do UE; o contexto no AAnF com o A-KID igual ao do UE; a K_AF entregue ao AF com `expiry`; a autenticação Ua\* |
| 2 | Atendido: todos os testes negativos da lista com o resultado esperado (`e2e.md` §4) |
| 3 | Atendido: matriz revisão 4, com as etapas de TS 33.535 §6.1 e §6.2.1 ligadas a verificações do roteiro. A "tabela de etapas mínimas" do material de partida não está no repositório (I17, I23) e foi reconstruída a partir do texto da norma |
| 4 | Atendido em 19/10/2026, com as correções abaixo. VM nova 192.168.16.131 com o free5GC v4.2.3 e o UERANSIM instalados pelo guia oficial, sem *patches*. O tutorial foi seguido do início, partes 1 a 8, com um clone limpo do repositório na máquina do AF. `e2e.sh woca@192.168.16.131`: 63 verificações, nenhuma falha (execução `20261019-160723`; `e2e.md` §5.1, `verificacao.md` §5.1) |

Diferenças em relação à tabela de arquivos:

- **Execução remota.** A parte da máquina do free5GC ficou em `af/scripts/e2e-remote.sh`. Ela roda numa única sessão `ssh -t`, com uma senha de `sudo`, e sincroniza com a máquina do AF por arquivos-marcadores.
- **Negativos da API.** Os testes da API do AAnF estão em `af/scripts/negative/aanf_api.py`, em HTTP/1.1 cru. Os do `register-anchorkey` usam o `seed` com as opções novas `--gpsi` e `--k-akma`.
- **Opções de teste.** O AF ganhou `clockOffset`, para o desvio deliberado de relógio. O cliente Ua\* ganhou `-kaf-file`, para apresentar a K_AF anterior depois de uma nova autenticação.
- **Sem `free5gc/scripts/` novo.** O roteiro usa os *scripts* existentes.
- **Endereços (critério 4).** O `e2e.sh` toma o endereço do AAnF do `aanfUri` de `af/config/afcfg.yaml`. No *preflight*, confere que ele é o `sbi.bindingIPv4` do `aanfcfg.yaml` instalado. Os AF extras usam o `listen` do mesmo arquivo. O `e2e-remote.sh` lê o AAnF do `aanfcfg.yaml` da máquina do free5GC. As chamadas ao `seed` passam `--aanf`. Antes, o endereço da VM 192.168.16.130 estava fixo nos dois.

Divergências que a reprodução em instalação limpa achou no tutorial (critério 4), todas corrigidas:

| Seção | O que estava escrito | O que acontecia | Correção |
|---|---|---|---|
| 3.3, 4.1 | repositório clonado na máquina do free5GC só na 4.1, sem endereço | o `patch-nf.sh` da 3.3 não existe ainda na máquina do free5GC | nova seção 1.1: clone em `~/akma-free5gc` com o endereço do GitHub, SSH por chave, gNB numa máquina só, assinante OPc, Ubuntu 18.04 |
| 3.2, 3.4 | "copie e ajuste" o `aanfcfg.yaml` | o `deploy.sh` instala o do repositório, com 192.168.16.130, e o AAnF não sobe em outra máquina | dizer qual arquivo ajustar: o do repositório antes do *deploy* ou a cópia instalada |
| 3.5, 3.6, 6.3 | `seed discover/register/remove` sem caminho e sem `--aanf` | `command not found`; depois, `connection refused` em 192.168.16.130 | `~/akma-tools/seed` e `free5gc/build/bin/seed`, com `--aanf` |
| 6.3 | K_AF errada com `-kaf`; AF_ID fora da política sem comando | com `-vector`, o `uaclient` não aceita `-kaf` (`-akid is required`) | os comandos completos; segundo AF com outro `fqdn` para o 503 |
| 6.3, 7.3 | semear o `testbed-ue` e, depois, o Ua\* do UE | o `testbed-ue` tem o SUPI do UE; o contexto do UE é substituído e o AF responde 403 `AKMA_KEY_NOT_AVAILABLE` | aviso: nova autenticação do UE (7.4) antes da 7.3 |
| 8.2 | `af/scripts/e2e.sh` sem argumento | roda contra 192.168.16.130; com o argumento, o AAnF continuava fixo em 192.168.16.130 | `e2e.sh usuario@máquina` e os endereços tirados das configurações (item acima) |

**Objetivo.** Executar todo o fluxo obrigatório da tabela de etapas mínimas de forma roteirizada e reprodutível, com o plano de testes negativos.

**Pré-requisitos.** MVP-2, MVP-3, MVP-4.

**Cláusulas.** Todas as da tabela de etapas mínimas · TS 29.500 §5.2.7.1 · TS 29.535 §5.1.7.3 · questão 1 tratada como etapa configurável, desabilitada, registrada como desvio.

**Arquivos.** `af/scripts/e2e.sh` (orquestra as duas VMs por SSH, coleta evidências), `af/scripts/negative/`, `free5gc/scripts/` (apoio), `docs/af-side/e2e.md`, `docs/free5gc-side/verificacao.md`, `docs/conformidade/desvios.md` (consolidado), matriz de rastreabilidade revisão 4, `docs/tutorial/` completo.

**Testes negativos.** A-KID desconhecido; rótulo `akid` aceito; RID `"012"` × `"12"`; AF não autorizado; `afId` malformado (9 dígitos hexadecimais, sem ponto); `AkmaKeyInfo` com `supi` e `gpsi`; `kAkma` fora do padrão; `Content-Type` incorreto (415); corpo acima do limite (413); demais códigos da matriz do MVP-1; desvio deliberado de relógio no AF (aceita K_AF expirada ou recusa válida); reautenticação: A-KID antigo recusado, K_AF antiga ainda válida no AF até `expiry`; assinante sem AKMA; AAnF indisponível durante a autenticação primária (conforme D9); `remove-context` seguido de `retrieve`.

**Critério de aceitação.**
1. `e2e.sh` termina com código 0 e gera relatório com: registro do UE, contexto AKMA criado no AAnF com A-KID igual ao do UE, K_AF entregue ao AF com `expiry`, autenticação Ua\* bem-sucedida.
2. Todos os testes negativos com o resultado esperado.
3. Cada linha da tabela de etapas mínimas ligada a uma evidência (teste ou trecho de log) na matriz revisão 4.
4. O tutorial, seguido do início sobre o free5GC v4.2.3 instalado e sem os *patches*, reproduz o resultado.

**Riscos.** Instabilidade do ambiente (`gtp5g`, *namespace*, ordem de inicialização). Reprodução em instalação limpa exige *snapshot* da VM.

**O tutorial ganha.** A versão completa do MVP: do free5GC recém-instalado até a aplicação autenticada com K_AF.

---

### MVP-6 — AF externo via NEF, acesso anônimo e remoção por OAM

**Estado (19/10/2026):** **aprovado** pelo responsável, com P11 (b) e DV-9 (b). Na aprovação, o responsável acrescentou o item 7: UERANSIM em outra máquina. Especificação `docs/free5gc-side/nef.md` aprovada em 19/10/2026, com todos os pontos da seção 9 decididos; o *token* do AF externo vem do NRF no MVP-6 e do CAPIF no MVP-9 (MVP-8 até a renumeração de 20/10/2026). Também decidido: o AAnF passa para a porta 8010, fora das portas do free5GC.

**Estado (19/10/2026, fim do dia):** implementado e verificado no laboratório. `e2e.sh`: 83 verificações e nenhuma falha (execução `20261019-213720`; `docs/af-side/e2e.md` §5.2).

| Critério | Resultado |
|---|---|
| 1 | Atendido: a mesma K_AF via NEF e direto do AAnF [E02]; resposta do NEF sem `supi` e sem `gpsi` [E01] |
| 2 | Atendido: Ua\* do UE com o AF externo, pelo plano de usuário [R2n.ua-nef]; também com o UE em outra máquina (`docs/af-side/teste-nef.md`) |
| 3 | Atendido: sem *token* e com *token* de outro serviço, 401 [E04, E05]; AF_ID fora da política do NEF, 403 [E06]; fora da política do AAnF, 403 [E07]; A-KID desconhecido, `403 K_AKMA_NOT_PRESENT` [E08]; 400 e 415 [E09 a E11] |
| 4 | Atendido: acesso anônimo direto [N01c, R2n.ua-anon] e via NEF [E03, R2n.ua-nef-anon], sem `supi` nem `gpsi` |
| 5 | Atendido: código 0 com os grupos novos e as 63 verificações do MVP-5; NEF sem SUPI no log [R2n.nef-no-supi]; nenhuma chave nos logs [R7, L1] |
| 6 | Atendido em 19/10/2026: gNB e UE na VM 192.168.16.131, free5GC na 192.168.16.130 e AF na 192.168.16.129; mesmo A-KID no UE, no AUSF e no AAnF; Ua\* com o AF interno, o externo e o externo anônimo; nenhuma chave nos logs das três máquinas (`teste-nef.md`) |

Diferenças em relação ao que foi proposto:

- **gNB na máquina do free5GC.** Com o N2 e o N3 em 192.168.16.130, o gNB local precisa de `ngapIp` 192.168.16.130 e de `gtpIp` 10.100.0.1. Um endereço de *loopback* não alcança 192.168.16.130. O endereço da máquina conflita com a porta 2152 do UPF.
- **N2 e N3 sem *patch*.** A mudança é a configuração que o guia oficial manda fazer quando o UERANSIM fica em outra máquina, e o tutorial a descreve (parte 9.5).
- **Entrada de teste no NEF.** `af3.akma.example.FF00000001` está no `nefcfg.yaml` e não no `aanfcfg.yaml`, para o negativo E07.
- **Log do AF externo.** Diz `obtained from the NEF`.

**Por que este escopo.**

- `[3GPP]` Pelo histórico de TS 33.535 (Anexo E), três procedimentos que o MVP-5 não cobre já estavam no Release 17:
  - a remoção do contexto (§6.6, CR 0072, V17.2.0);
  - o acesso anônimo (§6.2.2, esclarecido pelo CR 0126, V17.6.0);
  - o pedido de K_AF via NEF (§6.3, esclarecido pelo CR 0147, V17.8.0).

  A seção 2, item 10, só deixou para o MVP-6+ os acréscimos do Release 18.
- `[INTERPRETAÇÃO]` Dos três, só o §6.3 depende de uma NF do free5GC que o projeto ainda não alterou: o NEF. Com ele, a contribuição *upstream* cobre o AKMA do Release 17 em todas as NFs do lado da rede (NRF, UDM, AUSF, AAnF e NEF). Os perfis TLS do Ua\* (AKMA_TLS, também Release 17 pelo TR 21.917 §16.3) ficam do lado do AF e do UE.
- `[FREE5GC]` O laboratório já tem o que o teste pede:
  - o NEF do free5GC v4.2.3 está no `run.sh` e protege cada grupo de APIs *northbound* com um *middleware* OAuth2 por nome de serviço;
  - `3gpp-akma` já é um `ServiceName` no `openapi` v1.2.4 e em TS 29.510;
  - a máquina do AF faz o papel do AF fora da rede do operador.

**O que já existe e não entra como trabalho.**

| Item da lista da seção 5 | Onde já está |
|---|---|
| `anonInd` no AAnF (§6.2.2) | MVP-1: resposta sem `supi` quando `anonInd` é `true`; teste em `free5gc/aanf/internal/sbi/api_akma_test.go` |
| operação `remove-context` (§6.6) | MVP-1, no AAnF; o `seed remove` faz o pedido |
| filtro `routing-indicator` para `AANF` no NRF | MVP-1b, `free5gc/patches/nrf/0002-discover-aanf-by-routing-indicator.patch` |

**Escopo.**

1. **NEF** (TS 33.535 §6.3; TS 29.522 §4.4.23.2 e §5.14). *Patch* no NEF do free5GC com a operação `POST {apiRoot}/3gpp-akma/v1/retrieve`:
   - o NEF autoriza o AF (passo 2) com o seu *middleware* OAuth2 para `3gpp-akma` e com uma política local de AF_ID;
   - o NEF seleciona o AAnF pelo NRF, com o RID do A-KID (§6.7);
   - o NEF manda `Naanf_AKMA_ApplicationKey_Get` ao AAnF (passo 3) e repassa o `anonInd`;
   - a resposta ao AF tem `kaf` e `expiry` e **nunca** o SUPI (passo 5);
   - `403 K_AKMA_NOT_PRESENT` do AAnF é repassado ao AF (TS 29.522 §4.4.23.2).
2. **AAnF**: autorizar os pedidos que chegam pelo NEF (decisão P11, abaixo).
3. **Acesso anônimo de ponta a ponta** (§6.2.2):
   - opção `anonymous` no AF, que manda `anonInd: true`, direto ao AAnF ou via NEF;
   - o AF trata o A-KID como identificador temporário do usuário.
4. **AF fora da rede do operador**:
   - `nefUri` em vez de `aanfUri` no `afcfg.yaml`;
   - *token* OAuth2 para `3gpp-akma`.
   - `[3GPP]` §6.3 NOTA: sem CAPIF, o AF é configurado localmente com o ponto de terminação da API.
5. **Remoção por OAM** (§6.6), sem código novo.
   - `[3GPP]` O consumidor que a norma nomeia é o OAM: TS 33.535 tab. 7.1.1-1; TS 29.535 §4.2.1.3.2 e tab. 4.2.2.1-1, corrigida pelo CR 0057 (V19.4.0, *"Incorrect consumer of ContextRemove service operation"*).
   - O tutorial documenta o `seed remove` nesse papel, e o item "acionamento pelo núcleo" sai da lista.
6. **Testes e documentação**: `e2e.sh` com grupos novos para o caminho via NEF e para o acesso anônimo; tutorial, parte 9.
7. **UERANSIM em outra máquina** (pedido do responsável, 19/10/2026): free5GC na VM 192.168.16.130; AF na VM 192.168.16.129; gNB, UE (com o ME) e cliente Ua\* na VM 192.168.16.131. O UE autentica no AF pelo Ua\*, e o Ua\* passa pelo plano de usuário: do UE ao gNB na .131, pelo N3 até o UPF na .130 e pelo N6 até o AF na .129.
   - `[FREE5GC]` Hoje o N2 do AMF (`ngapIpList`) e o N3 do UPF (`gtpu`) estão em endereços de *loopback* (127.0.0.18 e 127.0.0.8), que um gNB em outra máquina não alcança. Os dois passam para 192.168.16.130, no AMF, no UPF e no SMF. O gNB local da .130 continua funcionando com o mesmo endereço.
   - `[FREE5GC]` O guia oficial do free5GC instala o UERANSIM numa máquina separada. Com este item, o tutorial cobre as duas topologias.

**Fora do escopo: o GPSI ao AF externo.**

- `[3GPP]` §6.3 passo 5: por política local, o NEF traduz o SUPI em GPSI pelo `Nudm_SubscriberDataManagement` e o inclui *"optionally"*.
- `[FREE5GC]` O `id-translation-result` do UDM v4.2.3 só traduz GPSI → SUPI (`GetIdTranslationResultProcedure`).

O NEF do MVP-6 não devolve GPSI, o que a norma permite. A tradução SUPI → GPSI exige um *patch* no UDM e fica para depois (seção 5).

**Decisões.** Tomadas pelo responsável em 19/10/2026: P11 (b) e DV-9 (b), registrado em `desvios.md`. O detalhe fica em `docs/free5gc-side/nef.md`, aprovado antes da implementação, como a especificação do Ua\* no MVP-3.

| Nº | Questão | Opções | Recomendação |
|---|---|---|---|
| P11 | Como o AAnF autoriza o NEF. `[3GPP]` §6.3 passo 2: quem autoriza o AF é o NEF. `[MVP]` Hoje a política local do AAnF (§6.2.1 passo 2) liga cada AF_ID a endereços de origem, e o pedido via NEF sai do endereço do NEF | (a) o AAnF aceita do NEF qualquer AF_ID e confere só o *token* (`nfType` NEF), porque o NEF já autorizou o AF; (b) além do *token*, uma lista de AF_ID externos que o AAnF aceita via NEF | (b): a política do AAnF continua sendo a última barreira |
| DV-9 | Segurança entre o NEF e o AF externo. `[3GPP]` TS 33.501 §12.2 e §12.3: autenticação mútua por certificados de cliente e servidor e TLS *"shall"*; §12.4: autorização OAuth *"shall"*. `[FREE5GC]` O NEF tem OAuth2 com *tokens* do NRF e TLS só de servidor (`sbi.scheme: https`), sem certificado de cliente | (a) TLS com certificado de cliente no NEF, só para `3gpp-akma`; (b) `http`, como o resto do laboratório, registrado como DV-9 ao lado do DV-1 e encerrado junto com ele | (b), se (a) exigir mudar o servidor SBI comum do free5GC; senão (a). A conferir no código do NEF e decidir no `nef.md` |

**Pré-requisitos.** MVP-5.

**Cláusulas.** TS 33.535 §6.2.2, §6.3, §6.6, §6.7, §7.1.3, §7.1.5, §7.3.2 · TS 29.522 §4.4.23.2, §5.14 · TS 29.535 §4.2.2.3.2 (`anonInd`) · TS 33.501 §12.

**Arquivos.**

- `free5gc/patches/nef/` (API `3gpp-akma` e consumidor do AAnF) e o alvo `nef` em `patch-nf.sh`;
- `free5gc/aanf/` (P11) e `free5gc/config/aanfcfg.yaml` (porta 8010);
- `af/af/` (`anonymous`, `nefUri`, *token*);
- `af/scripts/e2e.sh` e `e2e-remote.sh`;
- configuração do AMF, do SMF e do UPF para o N2 e o N3 fora do *loopback* (item 7), por *patch* em `free5gc/patches/config/`;
- `docs/free5gc-side/nef.md`, `docs/free5gc-side/correcoes-free5gc.md`, `docs/af-side/af.md`;
- documentação da topologia com o UERANSIM em outra máquina, nos dois lados;
- matriz de rastreabilidade revisão 5; `desvios.md` (DV-9);
- tutorial, parte 9.

**Critérios de aceitação.**

1. O AF externo obtém pelo NEF, para o mesmo A-KID e AF_ID, a mesma K_AF que o AAnF entrega a um AF interno. A resposta do NEF não tem `supi` (§6.3 passo 5).
2. O UE autentica no AF externo pelo Ua\*, pelo plano de usuário.
3. Negativos no NEF:
   - AF sem *token*, ou com *token* de outro serviço: recusado pelo *middleware*;
   - AF_ID fora da política do NEF: recusado, com o código de TS 29.122 fixado no `nef.md`;
   - A-KID desconhecido: `403 K_AKMA_NOT_PRESENT` repassado;
   - pedido do NEF fora da política do AAnF: 403.
4. Acesso anônimo, direto e via NEF: `anonInd: true`; resposta sem `supi` e sem `gpsi` (TS 29.522 tab. 5.14.5.3.3-1, NOTA 3); Ua\* bem-sucedido.
5. `e2e.sh` termina com código 0, com os grupos novos e as 63 verificações do MVP-5. Nenhuma chave aparece no log do NEF.
6. Com o gNB e o UE na VM 192.168.16.131, o free5GC na 192.168.16.130 e o AF na 192.168.16.129:
   - o UE registra;
   - o A-KID é igual no UE, no AUSF e no AAnF;
   - o UE autentica no AF pelo Ua\*, direto e via NEF;
   - nenhuma chave aparece nos logs das três máquinas.

**Riscos.**

- O NEF é a NF menos exercitada do `run.sh` do free5GC.
- A mudança do N2 e do N3 (item 7) altera a configuração que o MVP-5 usa. O roteiro do MVP-5 tem de continuar passando com o gNB local.
- `[MVP]` O AF externo obtém o *token* no NRF do operador, como o `seed` já faz. Fora do laboratório, quem emite o *token* para um AF externo é o CAPIF ou o próprio NEF (TS 33.501 §12.4, §12.5).

**O tutorial ganha.** Parte 9: um AF fora da rede do operador, pelo NEF, e o acesso anônimo. As partes 4, 7 e 8 passam a cobrir o UERANSIM em outra máquina, como no guia oficial do free5GC.

---

### MVP-7 — Repositório público

**Estado (20/10/2026):** **aprovado** pelo responsável, com as decisões abaixo. Começa pela limpeza do código.

**Estado (20/10/2026, tarde):** itens 1 a 4 do escopo feitos; critérios 1 a 5 atendidos numa exportação de teste (`export-public.sh --dirty`); falta o critério 6.

**Estado (20/10/2026, noite):** **concluído**. O repositório público foi exportado do commit do projeto, publicado no GitHub (`1527b46`) e instalado numa VM limpa seguindo só o README e o `docs/install.md`. O guia teve uma correção: o `aanfUri` do `ausfcfg.yaml` entrou na tabela de endereços. Ela vai ao repositório público na próxima reexportação.

| Critério | Resultado |
|---|---|
| 1 | Atendido: nenhuma referência ao projeto fora de `docs/`; testes Go passam com `-race`; *patches* aplicam nas árvores originais |
| 2 | Atendido: `e2e.sh` e `deploy.sh` exigem a máquina (argumento ou `FREE5GC_HOST`) |
| 3 | Atendido: 122 arquivos; `LICENSE` Apache-2.0 e `NOTICE` na raiz; `COPYING` GPL-3.0 em `af/ue/ueransim/patches/` |
| 4 | Atendido: na pasta exportada, `build-aanf.sh` e `build-af.sh` compilam; `vet` e `go test -race` passam nos quatro módulos; os geradores reproduzem os vetores |
| 5 | Atendido: autoria, referências ao projeto e dados pessoais sem ocorrência; as cadeias de 64 hexadecimais são todas de vetores de teste |
| 6 | Atendido em 20/10/2026 na VM 192.168.16.132 (*snapshot* limpo), com o clone do GitHub nas duas máquinas: o UE registra; o A-KID é igual no UE, no AUSF e no AAnF; Ua\* direto e via NEF; `e2e.sh` com 83 verificações e nenhuma falha (execução `20261020-173639`; `docs/af-side/teste-instalacao-github.md`) |

Os documentos públicos ficam em `public/`: `README.md`, `docs/install.md` (o tutorial reescrito, com a parte 11, como desfazer), `docs/free5gc-fixes.md` (F1 a F15), `docs/aanf.md`, `docs/nef-akma-api.md`, `docs/ua-star.md` e `docs/testing.md`. O `export-public.sh` recusa exportar com mudanças não commitadas, sem `--dirty`, e confere no fim que a pasta não cita documentos que ficam aqui.

**Objetivo.** Uma pasta separada, com o próprio `git`, que o responsável publica no GitHub e anuncia no fórum do free5GC, para a comunidade começar a testar. Segue o molde do free5GC: o guia de instalação e o repositório com o que baixar. Requisitos do pedido (19/10/2026):

- só o tutorial e o necessário para compilar e aplicar o que o projeto desenvolveu: código, *patches* e *scripts*;
- nada do projeto: sem plano, MVP, ADR, desvios, literatura, matriz nem marcações `[3GPP]`;
- texto escrito como um desenvolvedor escreve, direto;
- todos os pontos do free5GC a corrigir para o AKMA (F1 a F15 de `docs/free5gc-side/correcoes-free5gc.md`). Cada correção diz onde a norma manda fazer e como fazer segundo a norma;
- autoria só do responsável, em todo arquivo, mensagem de *commit* e metadado.

**Decisões.** Tomadas pelo responsável em 20/10/2026.

| Questão | Decisão |
|---|---|
| Idioma dos documentos públicos | Inglês, o idioma do fórum e da documentação do free5GC. Os documentos do projeto continuam em português |
| Repositório | `github.com/florianoneto07/akma-free5gc`, que já existe e está vazio. Os caminhos de módulo Go já usam esse nome. Este repositório do projeto continua só local |
| Licença | Apache-2.0 (a do free5GC) para o repositório; GPL-3.0 para os *patches* do UERANSIM, que derivam de código GPL-3.0 |
| AF e cliente Ua\* de teste | Vão junto: sem eles não há como testar o AKMA de ponta a ponta |

**Escopo.**

1. **Limpeza do código, neste repositório.** Os comentários do código, dos *scripts*, das configurações e dos *patches* deixam de citar:
   - documentos do projeto;
   - números de MVP, DV, ADR e decisões (D9 a D16, P11, P12);
   - o laboratório.

   Passam a citar só a cláusula 3GPP, ou um texto neutro. O código continua tendo uma fonte só: este repositório. Os identificadores não mudam, exceto os que levavam o nome de etapas ou documentos do projeto: `afid.UaStarMVP` passa a `afid.UaStarDigest`, o vetor `baseline-2.5` a `worked-example` e o teste `TestSBaseline` a `TestSWorkedExample`.
2. ***Scripts* sem dado pessoal.** O `e2e.sh` e o `deploy.sh` deixam de ter `woca@192.168.16.130` como padrão. A máquina vem por argumento ou por variável de ambiente.
3. **Documentos públicos em inglês**, escritos neste repositório em `public/`, cujo conteúdo vai para a raiz do repositório público:
   - `README.md`: o que é, o estado, a arquitetura, o início rápido, as limitações e a licença;
   - guia de instalação: o tutorial reescrito, nas duas topologias (UERANSIM na máquina do free5GC e em outra máquina);
   - correções do free5GC, F1 a F15. Cada uma traz o que o free5GC faz, a cláusula da norma e como fazer segundo ela, inclusive o 401 sem ProblemDetails (F12);
   - AAnF, API AKMA do NEF, Ua\* e testes;
   - limitações, ditas sem rodeio:
     - sem TLS nas interfaces SBI e entre o NEF e o AF;
     - Ua\* simplificado, com identificador da faixa privada;
     - o NRF emite o *token* do AF externo;
     - `nr-cli` como interface de teste no UE;
     - *patches* para o free5GC v4.2.3 e o UERANSIM v3.2.6.
4. **Exportação.** `public/export-public.sh <pasta>` monta a pasta pública a partir deste repositório.
   - **Vai:** `shared/` (com os vetores de teste e a validação cruzada), `free5gc/aanf`, `free5gc/patches`, `free5gc/scripts`, `free5gc/config`, `af/af`, `af/uastar`, `af/ue`, `af/config`, `af/scripts`, `reference/kdf` (o oráculo da KDF), os documentos de `public/` e as licenças.
   - **Não vai:** `docs/`, `papers/`, `specs/`, o `README.md` e o `CHANGELOG.md` do projeto, `free5gc/build`, `af/build`, a história do git, as *tags* e o próprio *script*.
   - Na pasta exportada eu só rodo `git init`; o *commit* e o *push* são do responsável.

**Fora do escopo.**

- O *pull request* ao free5GC (seção 5, *Upstream*).
- O Ua\* da norma (MVP-8), o CAPIF (MVP-9) e o TLS nas SBI.

**Pré-requisitos.** MVP-6.

**Arquivos.**

- comentários em `shared/`, `free5gc/`, `af/` e `reference/`;
- `free5gc/patches/nef/0001-akma-api.patch`, regenerado se os comentários do NEF mudarem, e os demais *patches* com comentário a limpar;
- `af/scripts/e2e.sh` e `free5gc/scripts/deploy.sh`;
- `public/` (documentos e *script* de exportação).

**Critérios de aceitação.**

1. Nenhum arquivo fora de `docs/` cita documento do projeto, MVP, DV, ADR, decisão numerada ou o laboratório. Os testes Go passam em todos os módulos. Os *patches* aplicam no free5GC v4.2.3 e no UERANSIM v3.2.6.
2. Nenhum *script* tem usuário ou máquina pessoal como padrão.
3. A pasta exportada tem só o que o escopo lista, com a `LICENSE` Apache-2.0 na raiz e a GPL-3.0 junto dos *patches* do UERANSIM.
4. Na pasta exportada, sem nada deste repositório, tudo compila e os testes Go passam.
5. Varreduras na pasta exportada:
   - autoria: o responsável é o único autor citado;
   - projeto: nenhuma menção a plano, MVP, DV, ADR, matriz ou às marcações `[3GPP]`;
   - dados pessoais: nenhum caminho `/home/<usuário>`, nenhum usuário e nenhum e-mail além dos do autor;
   - chaves: nenhuma cadeia de 64 hexadecimais fora dos vetores de teste.
6. Numa VM limpa, seguindo só os documentos públicos:
   - o UE registra;
   - o A-KID é igual no UE, no AUSF e no AAnF;
   - o UE autentica no AF pelo Ua\*, direto e via NEF;
   - o `e2e.sh` termina com código 0.

**Riscos.**

- A VM limpa do critério 6 exige reverter o *snapshot* da .131 ou criar uma VM nova. Decidir com o responsável antes desse critério.
- Mudar um comentário dentro do *patch* do NEF obriga a regenerá-lo e a conferir que ele aplica. O binário na VM do free5GC não muda: comentários não alteram o código.
- Os *replace* do `go.mod` do NEF apontam para `../../../akma-free5gc/`. A pasta clonada precisa ter esse nome, que é o do repositório público, e o guia de instalação diz isso.

**O tutorial ganha.** Nada: o tutorial do projeto continua como está. O guia de instalação público é a sua reescrita em inglês.

---

### MVP-8 — Ua\* da norma e autorização por aplicação no UE

**Estado (20/10/2026):** organização **aprovada** pelo responsável. O escopo detalhado de cada fase sai de uma especificação aprovada antes do código, como `ua-star.md` no MVP-3 e `nef.md` no MVP-6.

**Estado (20/10/2026, noite):** especificação das fases (a) e (b) aprovada (`docs/af-side/ua-star-norma.md`, pontos U1 a U12). **Fase (a) implementada e verificada.** `e2e.sh` com 88 verificações e nenhuma falha (execução `20261020-215257`; `docs/af-side/e2e.md` §5.3). O *patch* é o `af/ue/ueransim/patches/0005-akma-key-service.patch`.

| Caso da especificação (§2.6) | Resultado |
|---|---|
| aplicação e AF_ID da política | atendido: K_AF entregue ao `uaclient` [R2k.given]; Ua\* passa [R2.ua] |
| AF_ID fora da lista da aplicação | atendido: `ERR NOT_AUTHORIZED` e `warn` no log do UE [R2k.af-id] |
| outro usuário | atendido: `nobody` recusado [R2k.uid] |
| outro executável do mesmo usuário | atendido à mão em 21/10/2026: `python3` recusado com `executable` na política (`e2e.md` §5.3) |
| AF_ID em duas aplicações | atendido: o `nr-ue` não sobe [R2k.duplicate] |
| `akma-kaf` sem `akmaCliKaf` | atendido: comando recusado [R2k.nr-cli] |
| nenhuma K_AKMA nas respostas; nenhuma K_AF em log | atendido: as respostas só têm A-KID, AF_ID e K_AF; nenhuma cadeia de 64 hexadecimais nos logs [R7] |

Diferenças em relação à especificação, registradas nela:

- o *socket* padrão leva o SUPI, `/run/ueransim/akma-<SUPI>.sock`, porque o `nr-ue -n` roda vários UEs num processo;
- entra a causa `UNAVAILABLE`, para quando o ME não consegue pausar as tarefas a tempo.

**Estado (21/10/2026):** **fase (b) implementada e verificada.** `e2e.sh` com 115 verificações e nenhuma falha (execução `20261021-164401`; `docs/af-side/e2e.md` §5.4). A implementação e as diferenças estão em `ua-star-norma.md` §3.8; o ADR-0003 foi revisto (U11).

| Caso da especificação (§3.7) | Resultado |
|---|---|
| K.2 em TLS 1.3 e 1.2, com o AF_ID do *cipher suite* | atendido [R2t.k2-tls13, R2t.k2-tls12, T01, T02] |
| K.2 com certificado ou `realm` de outro FQDN | atendido: certificado [T09]; `realm` no teste de unidade `TestK2RealmOfAnotherAF` |
| PSK TLS 1.3, `psk_dhe_ke`, AF_ID `…0100000007` | atendido [R2t.psk13, T03] |
| PSK TLS 1.2, suíte `D0 01`, *hint* `3GPP-AKMA`, AF_ID `…010001D001` | atendido [R2t.psk12, T04] |
| retomada de sessão | atendido nos três perfis [R2t.k2-resumed, R2t.psk12-resumed, T03, T04] |
| `psk_ke`, K_AF errada, A-KID desconhecido, AF_ID fora da política | recusados [R2t.psk-ke, T06, T07, T08, T10] |
| conferência independente com o OpenSSL | atendido [T05] |
| K_AF vencida, `refetch` e `reject` | atendido [T11, T12, T13] |
| nenhuma chave em log | atendido [L1, R7] |

Diferenças em relação à especificação, registradas nela (§3.8): wolfSSL 5.9.4 em vez de 5.8.4; sem a suíte `D0 02`, que o wolfSSL não tem; sem OCSP *stapling*; o UE dos perfis PSK só aceita PSK (U6); o UE pede a K_AF ao ME só depois do desafio. A `D0 02` e o OCSP *stapling* são recomendados (*should*) por TS 33.210 §6.2; as duas ausências foram registradas como DV-10, por decisão do responsável em 21/10/2026.

**Estado (21/10/2026, fase (c)):** especificação do DTLS **aprovada** (`ua-star-norma.md` §4), com os pontos U13 a U21 como propostos. Os principais:
- U13: o DTLS 1.2, que TS 33.210 exige de quem tem DTLS e TS 33.535 Anexo C não descreve;
- U14: CoAP dentro do DTLS;
- U15: a suíte obrigatória do CoAP, sem ECDHE, contra TS 33.210;
- U21: o OpenSSL 4.1 como conferência independente do DTLS 1.3.

**Estado (22/10/2026):** **fase (c) implementada e verificada.** `e2e.sh` com 126 verificações e nenhuma falha (execução `20261022-124514`; `docs/af-side/e2e.md` §5.5): o R2d e o grupo D, com o OpenSSL 4.1.0 *beta* 2 no papel do UE no D03. A implementação, as diferenças e dois achados estão em `ua-star-norma.md` §4.7: o wolfSSL perguntava a chave antes do *cookie* no DTLS 1.3, e o `ClientHello` do OpenSSL 4.1 passava do que o wolfSSL lê por padrão. Os casos e o resultado de cada um estão em §4.8. A suíte obrigatória do CoAP ficou de fora, como DV-11.

**Estado (22/10/2026, fase (d)):** especificação do OSCORE **aprovada** (`ua-star-norma.md` §7), com os pontos U22 a U35 como propostos. Os principais:
- U22: o identificador `01 00 00 00 08` do OSCORE do GBA, AF_ID `<FQDN>.0100000008`;
- U25 e U26: a forma do HPLMN-ID e os parâmetros do HKDF, que o Anexo D não fixa;
- U32: a renovação pelo Apêndice B.2 da RFC 8613, iniciada pelo UE;
- U33: CBOR, OSCORE e AES-CCM em Go próprio, sem o wolfSSL;
- U34: a libcoap como conferência independente;
- U35: cinco defeitos do Anexo D no registro de defeitos (D-J a D-N).

**Estado (22/10/2026):** **fase (d) implementada e verificada.** `e2e.sh` com 141 verificações e nenhuma falha (execução `20261022-201824`; `docs/af-side/e2e.md` §5.6): o R2o e o grupo O, com a libcoap 4.3.5b no papel do UE no O04. A implementação e as diferenças estão em `ua-star-norma.md` §7.11, e os casos com o resultado de cada um em §7.12. Na mesma data, o registro de defeitos foi conferido de novo contra as normas: os itens D-A a D-I são todos das normas, e o D-B passou a MVP-12.

**Pedido** (responsável, 20/10/2026):

- o Ua\* completo, não o simplificado, já com TLS, DTLS e OSCORE conforme a norma;
- entender e resolver a limitação do `nr-cli`, que entrega a K_AF de qualquer AF a qualquer processo da máquina.

**Decisões** (responsável, 20/10/2026):

| Questão | Decisão |
|---|---|
| Organização | Um MVP em quatro fases, cada uma com critérios e *commits* próprios: (a) serviço de chaves do ME com autorização por aplicação; (b) TLS; (c) DTLS 1.3; (d) OSCORE. O CAPIF passa a MVP-9 |
| `akma-kaf` do `nr-cli` | Fica atrás de uma opção de teste, desligada por padrão. O caminho normal passa a ser o serviço de chaves |
| Ua\* simplificado (`FF00000001`) | Fica como perfil de teste, desligado por padrão no AF |
| Normas | Baixadas para `specs/` em 20/10/2026: TS 33.310 V19.5.0 e TS 33.210 V19.3.0; RFC 4279, 5489, 5869, 6066, 6347, 7252, 8442, 8446, 8613, 8949, 9052, 9053, 9147, 9258 e 9325 em `specs/rfc/` |

**Situação normativa.**

- `[3GPP]` TS 33.535 §4.4.3, requisitos do UE:
  - as aplicações não acessam a K_AKMA (já atendido);
  - *"An application on the UE shall only get the KAF keys related to specific AF Identifiers (AF_IDs) that the application is authorized to get"*;
  - *"An application on the UE shall not be able to get access to the KAF keys that belong to other applications"*;
  - NOTA: *"How these requirements are satisfied is out of scope of 3GPP."* `[LACUNA]` A forma fica para o projeto.
- `[3GPP]` Perfis do Ua\* (TS 33.535 Anexos B, C e D; TS 33.222 §5.3 e §5.4; TS 24.109 Anexo K; TS 33.220 Anexos H, O e P):

  | Perfil | Identificador Ua\* (TS 33.220 H.3) |
  |---|---|
  | Digest dentro de TLS, AF autenticado por certificado (B.1.2; K.2) | `01 00 01 yy zz`, com o *cipher suite* negociado |
  | PSK TLS 1.2 (B.1.3.2.1; K.3.2) | `01 00 01 yy zz` |
  | PSK TLS 1.3 (B.1.3.2.2; K.3.3) | `01 00 00 00 07` |
  | DTLS 1.3, opcional (Anexo C) | `01 00 00 00 09` |
  | OSCORE, opcional (Anexo D) | `01 00 00 00 08` |

- `[3GPP]` TS 33.210 §6.2, o perfil TLS para o qual TS 33.310 Anexo E remete:
  - TLS 1.3: `psk_ke` *"shall not be supported"*, então o PSK vai com ECDHE (`psk_dhe_ke`);
  - TLS 1.2 com PSK: RFC 5489, com `TLS_ECDHE_PSK_WITH_AES_128_GCM_SHA256` (RFC 8442) obrigatório; o identificador fica `01 00 01 D0 01`;
  - com DTLS, o DTLS 1.2 *"shall"* e o DTLS 1.3 *"should"* ser suportados.
- `[FREE5GC]` O `crypto/tls` do Go não tem PSK nem DTLS. A pilha das fases (b) e (c) é escolhida na especificação.

**Pontos a decidir na especificação.**

1. `[LACUNA]` Como derivar da K_AF o *premaster secret* (TLS 1.2) e a PSK externa (TLS 1.3). TS 33.222 §5.4, TS 24.109 K.3 e TS 33.535 B.1.3 dizem só "derivar da K_AF"; TS 33.210 não diz.
2. `[LACUNA]` O separador da identidade PSK no TLS 1.2: TS 24.109 K.3.2 e) não o dá. O GBA usa `;` (TS 24.109 §5.3.3.1.2).
3. `[LACUNA]` DTLS 1.2: TS 33.210 §6.2.1 o exige quando há DTLS, e TS 33.535 C.1.1 diz que a identidade PSK de B.1, que inclui o TLS 1.2, vale para o DTLS; mas C.2 só descreve o procedimento do DTLS 1.3.
4. `[LACUNA]` Os parâmetros do HKDF em `OMS = HKDF(KAF, "AKMA-OSCORE")` (TS 33.535 D.3.3): *salt*, *info* e tamanho (RFC 5869).
5. `[LACUNA]` A forma do HPLMN-ID no A-KID em CBOR do OSCORE (D.3.5).
6. `[DIVERGÊNCIA]` TS 33.535 D.3.2 remete a "TS 33.220 clause Y.2.3" e "Y.3.3"; o conteúdo está no Anexo P.
7. `[DIVERGÊNCIA]` D.3.2 põe um *"CoAP Security protocol identifier"* no corpo do pedido; o CBOR de D.3.5 não o tem.
8. No K.2 continuam o DV-5 e o DV-6. O identificador Ua\* depende do *cipher suite*, então a política do AAnF precisa de um AF_ID por *cipher suite* aceito.
9. A identidade da aplicação no ME e o formato da política (`[INTERPRETAÇÃO]`, ponto da NOTA de §4.4.3).
10. A pilha de PSK TLS, DTLS e OSCORE. Verificado em 20/10/2026:
    - o OpenSSL 3.0.13 da VM do AF não tem `ECDHE-PSK-AES128-GCM-SHA256`, a suíte obrigatória de TS 33.210 §6.2.3 (só as variantes CBC e ChaCha20);
    - o OpenSSL só traz DTLS 1.3 a partir da 4.1, anunciada pelo projeto para outubro de 2026;
    - o wolfSSL tem TLS 1.3, DTLS 1.3 e PSK, sob GPLv3 ou licença comercial; a suíte `D0 01` aparece nos testes do próprio wolfSSL, mas ainda não foi conferida compilando;
    - não há biblioteca de OSCORE em Go conhecida. O RFC 8613 traz vetores de teste no Apêndice C.

**Fases.**

- **(a) Serviço de chaves do ME.** O `nr-ue` atende por um *socket* Unix. A aplicação é identificada pelo sistema operacional (`SO_PEERCRED` e o executável), e a política na configuração do UE dá os AF_IDs de cada aplicação. Um AF_ID não pode pertencer a duas aplicações. O `uaclient` passa a usar o serviço.
- **(b) TLS.** O K.2, com certificado do AF; o PSK TLS 1.3 e o PSK TLS 1.2.
- **(c) DTLS 1.3**, mais o DTLS 1.2, conforme a decisão do ponto 3.
- **(d) OSCORE.** CoAP, OSCORE (RFC 8613) e CBOR, com o recurso `/akma`.

Em todas as fases: o AF escolhe os perfis pela configuração, o AAnF autoriza os AF_IDs de cada perfil e o `e2e.sh` ganha os grupos novos.

**Pré-requisitos.** MVP-6. O critério 6 do MVP-7, a instalação limpa pelo GitHub, corre em paralelo noutra sessão e não bloqueia.

### MVP-12 — Renovação da K_AKMA pelo AAnF

**Estado (28/10/2026):** especificação **aprovada** pelo responsável, com os pontos K1 a K12 como propostos (`docs/free5gc-side/renovacao-kakma.md`).

**Estado (28/10/2026, tarde):** **implementado e verificado.** `e2e.sh` com 161 verificações e nenhuma falha (execução `20261028-162316`; `docs/af-side/e2e.md` §5.7), com as 20 do passo novo R6r. Duas decisões do responsável na implementação: a notificação ao AMF vai sem *token* de acesso (TS 29.500 §6.7.3; F20), e as recusas do AMF seguem a tabela 6.2.5.6-3 de TS 29.503 (`409` para registro e *handover* em curso). A implementação, as diferenças e os casos com o resultado de cada um estão em `renovacao-kakma.md` §10 e §11. Novos: DV-12 em `desvios.md`; D-O e D-P no registro de defeitos; F16 a F20 em `correcoes-free5gc.md`; `docs/free5gc-side/amf.md`; parte 10 do tutorial.

**Objetivo.** Quando a K_AKMA de um UE passa da idade da política local, o AAnF pede ao UDM uma nova autenticação primária (`Nudm_UECM_AuthTrigger`). O UDM notifica o AMF, e o AMF autentica o UE de novo (TS 33.535 §6.4.4; TS 33.501 §6.1.5). Com a nova autenticação, o AUSF registra no AAnF a nova K_AKMA, e o UE passa a ter um A-KID novo.

**Escopo.** Seções 3 e 4 da especificação:

- `openapi`, UDM e AMF do free5GC: F16 a F19 de `docs/free5gc-side/correcoes-free5gc.md`;
- AAnF: `kakmaRefresh` e o consumidor do UDM; `seed trigger-auth`;
- um passo novo do `e2e.sh`;
- desvio novo DV-12 (POST no `trigger-auth`) e defeito novo D-O.

**Critérios.** Os casos da seção 6 da especificação. Com `kakmaRefresh` desligado, o padrão, o `e2e.sh` atual passa igual.

**Pré-requisitos.** MVP-8.

### MVP-13 — V2 do repositório público

**Estado (28/10/2026):** **aprovada** pelo responsável, com os pontos L1 a L12 como propostos e o ponto L13, acrescentado por ele: uma terceira lista de recomendações, ao UERANSIM. As três listas se baseiam nas normas. Ao fim, o responsável recebe os comandos de *commit*, *tag* e *push* do repositório público. Em execução.

**Estado (28/10/2026, noite):** documentos prontos neste repositório; faltam a exportação, a verificação em VM limpa (L11) e a publicação.

- **L9, normas.** Baixadas as versões mais recentes: TS 29.503, TS 29.522 e TS 29.535 V20.1.0, com os YAML da Release 20; TS 33.501 V20.3.0, TS 23.502 V20.3.0, TS 24.501 V20.1.0, TS 23.003 V20.1.0, TS 33.210 V20.0.0 e TS 29.500 V20.0.0. TS 33.535, TS 24.109 e TS 33.220 não têm versão nova. Todos os itens continuam; D-H e D-I estão corrigidos em parte (registro de defeitos, conferência de 28/10/2026).
- **L7, sugestões.** D-Q a D-Z no registro de defeitos. O D-Z, de TS 24.501 §5.4.2.3, saiu da análise do UERANSIM (L13).
- **L6, free5GC.** O free5GC v4.3.0 saiu em 30/09/2026; F1 a F20 continuam nele (`correcoes-free5gc.md`). Os *patches* continuam para a v4.2.3.
- **L13, UERANSIM.** A v3.3.0 (07/06/2026) e o `master` (14/10/2026) têm os quatro itens de `docs/af-side/recomendacoes-ueransim.md`. Desde abril de 2025 o UERANSIM deriva uma K_AKMA, com o FC da K_AUSF (`0x6A`) e sem o P0 `"AKMA"`, o que não confere com TS 33.535 A.2.
- **Públicos.** `public/docs/3gpp-recommendations.md` e `public/docs/ueransim-recommendations.md` (novos), `public/CHANGELOG.md` (novo), `free5gc-fixes.md` (conferência na v4.3.0 e tabela por repositório), `README.md` (V2, MVP-12, recomendações) e `install.md` (parte 12, atualização a partir da V1).

**Objetivo.** Publicar a V2 em `github.com/florianoneto07/akma-free5gc`: o que o projeto fez depois da V1 (MVP-7, 20/10/2026) e as duas listas de recomendações que o responsável vai enviar ao free5GC e à 3GPP ao fim da implementação. Pedido do responsável em 28/10/2026. **Nome:** MVP-13 neste repositório; no repositório público, só V2, sem a numeração dos MVP (decisão de 28/10/2026). O `export-public.sh` recusa a palavra "MVP" na pasta exportada.

**O que mudou desde a V1.** Exportação de teste do *commit* `742ce41` (28/10/2026), comparada com a árvore pública da V1 (`e44188a`, a mesma do GitHub):

- 182 arquivos contra 122; 41 alterados; nenhum removido;
- código novo:
  - do MVP-8: `af/akmatls`, `af/wolftls` (GPL-3.0, com `COPYING`), `af/coap`, `af/oscore`, `af/cbor` e `af/akmaoscore`; os perfis do AF e do `uaclient`; o *patch* `0005` do UERANSIM (serviço de chaves do ME); `build-wolfssl.sh`, `make-test-ca.sh` e `oscore-check.py`;
  - do MVP-12: `free5gc/patches/amf`, `udm/0002` e `openapi/0003`; no AAnF, `kakmaRefresh` e o consumidor do UDM; `seed trigger-auth`;
- documentos públicos já atualizados em `public/` ao longo dos MVP-8 e 12: `README.md` (perfis Ua\*), `NOTICE` (wolfSSL sob GPL-3.0), `docs/ua-profiles.md` (novo), `install.md` (parte 11, renovação da K_AKMA), `free5gc-fixes.md` (F16 a F20), `aanf.md`, `testing.md` (161 verificações) e `ua-star.md`;
- a varredura da exportação de teste não achou referência ao projeto (MVP, DV, ADR, pontos K e U, documentos do projeto, laboratório). Só aparecem P0 e P1, nomes dos parâmetros da KDF das normas.

Falta, para a V2: o MVP-12 no `README.md`, um registro de mudanças público, a atualização a partir da V1 no guia, a lista do free5GC organizada para o envio e a lista da 3GPP, que ainda não existe em inglês.

**Pontos para aprovação.**

| Nº | Ponto | Proposta | Base e alternativas |
|---|---|---|---|
| L1 | Nome e versão | "V2" nos textos públicos. O responsável cria a *tag* `v2.0.0` no *commit* da V2 e, retroativamente, `v1.0.0` em `e44188a` (a V1 com a correção do guia). Um *release* no GitHub com o texto do registro de mudanças (L4) | nome: decisão de 28/10/2026; numeração: *Semantic Versioning*. Alternativa: só `v2`, sem *tag* para a V1 |
| L2 | Código | tudo o que o `export-public.sh` já leva, do *commit* final do MVP-13, sem tirar nada da V1. O Ua\* simplificado (`digest-http`) continua como perfil de teste, desligado por padrão | escopo do MVP-7, item 4 |
| L3 | `README.md` | acrescentar o MVP-12: a renovação da K_AKMA na tabela do que funciona, o *patch* do AMF na lista e no início rápido, as 161 verificações. Duas seções novas: "What's new in V2" e "Recommendations", com as duas listas (L6 a L8). As versões-alvo (free5GC v4.2.3, UERANSIM v3.2.6) e a data da última verificação | — |
| L4 | Registro de mudanças público | `CHANGELOG.md` novo, em inglês, com a V1 (20/10/2026) e a V2. Serve de texto para o *release* do GitHub | — |
| L5 | Atualização a partir da V1 | seção nova "Upgrading from V1" no `install.md`, com os comandos: desfazer cada conjunto de *patches* com o clone da V1 (NRF, UDM, AUSF, NEF, configuração, `run.sh`, banco de testes, UERANSIM), atualizar o clone, preparar de novo o `openapi`, aplicar os da V2 (o AMF é novo), os AF_IDs novos nas configurações do AAnF, do NEF e do UE, o wolfSSL, o AAnF e o `seed` | hoje o guia só diz *"revert with the old version first"* |
| L6 | Recomendações ao free5GC | `docs/free5gc-fixes.md` já tem F1 a F20. Acrescentar uma tabela por repositório do free5GC (`nrf`, `udm`, `ausf`, `amf`, `nef`, `openapi`, `webconsole`), com os itens e o *patch* de cada um, para um *issue* ou *pull request* por repositório. Antes de publicar, conferir cada item na versão mais recente do free5GC: "ainda presente em …" ou "corrigido em …" | `correcoes-free5gc.md`; pedido de 28/10/2026. Alternativa: a lista como está |
| L7 | Recomendações à 3GPP: conteúdo | documento novo `docs/3gpp-recommendations.md`, com: (a) os 16 defeitos D-A a D-P do registro; (b) como sugestões, as lacunas que o projeto teve de decidir: a codificação textual do A-TID (questão C, ADR-0002); a derivação da PSK do TLS 1.3 e do *premaster* do TLS 1.2 a partir da K_AF (U1); o separador da identidade PSK no TLS 1.2 (U2); o DTLS 1.2 (U13); o protocolo dentro do DTLS (U14); os erros do OSCORE (U30); a suíte obrigatória do CoAP com PSK contra TS 33.210 (DV-11); se o UDM responde ao `trigger-auth` antes ou depois do AMF (K4); o servidor de autorização do NEF sem CAPIF (TS 29.522 §6; F13). As sugestões entram antes no registro de defeitos, de D-Q em diante, com os mesmos campos. As divergências de enumeração de TS 23.502 vão no fim, como observações editoriais | pedido de 28/10/2026 ("erros ou sugestões"). Alternativa: só os defeitos D-A a D-P |
| L8 | Recomendações à 3GPP: forma | itens agrupados pelo grupo responsável (SA3: TS 33.535 e 33.220; CT1: TS 24.109; CT3: TS 29.522 e 29.535; CT4: TS 29.503; SA2: TS 23.502). Cada um com os campos de um *Change Request*: especificação e versão, cláusulas afetadas, texto atual (citado), motivo, mudança proposta, consequência se não mudar, o que este repositório faz, e a prioridade (correção ou editorial). O documento não diz como submeter; isso fica com o responsável | os campos do formulário de CR da 3GPP |
| L9 | Versões das normas | antes de publicar, conferir cada item na versão publicada mais recente de cada especificação: a última da Release 19 e, se existir, a da Release 20, baixadas em `specs/`. O documento público e o registro dizem em que versão cada item foi confirmado. Um item já corrigido sai do documento público e fica no registro, com a nota | o registro usa as versões do acervo: de V19.0.0 (TS 33.535, TS 24.109) a V19.7.0 (TS 29.503, TS 29.522) |
| L10 | Publicação | exportar do *commit* final para uma pasta temporária e sincronizar em `~/akma/akma-free5gc-public`, mantendo o `.git` (`rsync --delete`). Eu mostro `git status` e `git diff --stat` dali; o responsável faz um *commit* para a V2, as *tags* e o *push* | MVP-7: não faço *commit* no repositório público |
| L11 | Verificação em VM limpa | antes do *push*: um `git bundle` do repositório público instalado no *snapshot* limpo da 192.168.16.132, e o AF na máquina do AF com um clone do mesmo *bundle*, seguindo só os documentos públicos. O `e2e.sh` termina com 161 verificações e nenhuma falha (com `OPENSSL_DTLS13` e `LIBCOAP_CLIENT`). O responsável reverte o *snapshot* | critério 6 do MVP-7. Alternativa: depois do *push*, pelo clone do GitHub, como na V1 |
| L12 | Atualização V1 → V2 | (a) descrita (L5) e conferida contra os *patches* e configurações das duas versões, sem instalar a V1 antes. (b) Verificada numa VM com a V1 instalada, se houver um *snapshot* assim | (a) recomendada: instalar a V1 numa VM limpa antes da V2 dobra o tempo da verificação |
| L13 | Recomendações ao UERANSIM | documento novo `docs/ueransim-recommendations.md`, em inglês, com o que o projeto achou no UERANSIM contra as normas: o AKMA no ME, que os *patches* do repositório implementam (TS 33.535 §4.4.3, §6.1, §6.2, Anexo A), e os comportamentos do UE fora da norma vistos nos testes, cada um com a cláusula, o que o UERANSIM faz, como deve ser e o que o repositório faz. Cada item conferido na versão mais recente do UERANSIM. Registro no projeto: `docs/af-side/recomendacoes-ueransim.md` | acrescentado pelo responsável em 28/10/2026 |

**Fora do escopo.**

- O *pull request* ao free5GC e o *rebase* para uma versão mais nova (seção 5, *Upstream*).
- O CAPIF (MVP-9), o GPSI (MVP-10) e o TLS nas SBI (MVP-11).
- O envio das recomendações, que fica com o responsável.

**Pré-requisitos.** MVP-12.

**Arquivos.**

- `public/README.md`, `public/CHANGELOG.md` (novo), `public/docs/install.md`, `public/docs/free5gc-fixes.md`, `public/docs/3gpp-recommendations.md` (novo) e `public/docs/ueransim-recommendations.md` (novo);
- `docs/af-side/recomendacoes-ueransim.md` (novo);
- `docs/defects/registro-de-defeitos-3gpp.md` (D-Q em diante e as versões conferidas) e a cópia em `~/akma`;
- `docs/free5gc-side/correcoes-free5gc.md` (o repositório de cada item e a versão do free5GC conferida);
- este plano;
- `specs/` (versões novas das normas, fora do controle de versão).

**Critérios de aceitação.**

1. A exportação sai de um estado com *commit*, sem `--dirty`. Na pasta exportada, sem nada deste repositório: todos os módulos compilam; `go vet` e `go test -race` passam, com e sem a *tag* `wolfssl`; os *patches* aplicam no free5GC v4.2.3, nas duas versões do `openapi`, e no UERANSIM v3.2.6.
2. Varreduras na pasta exportada:
   - autoria: o responsável é o único autor citado;
   - projeto: nenhuma menção a MVP, DV, ADR, pontos K e U, documentos do projeto ou ao laboratório;
   - dados pessoais: nenhum caminho `/home/<usuário>`, nenhum usuário e nenhum e-mail além dos do autor;
   - chaves: nenhuma cadeia de 64 hexadecimais fora dos vetores de teste.
3. Os documentos públicos concordam entre si e com o código: o `README.md`, o `CHANGELOG.md`, o guia (com a parte 11 e a atualização a partir da V1), as três listas (todos os itens com a versão conferida) e o `testing.md` (161 verificações). Os *links* internos resolvem.
4. Na VM limpa (L11): o UE registra; o A-KID é igual no UE, no AUSF e no AAnF; o Ua\* passa por todos os perfis; a K_AKMA é renovada (R6r); o `e2e.sh` termina com 161 verificações e nenhuma falha.
5. O registro de defeitos e `correcoes-free5gc.md` dizem a versão conferida de cada item (L6, L9) e trazem as sugestões novas (L7).

**Riscos.**

- A instalação limpa leva horas: o guia compila o UERANSIM e o wolfSSL. O *snapshot* da .132 tem de ser revertido pelo responsável antes.
- O 3gpp.org pede um `User-Agent` de navegador para os *downloads*. As versões novas podem ter mudado a numeração das cláusulas; o documento público cita as da versão conferida.
- Uma versão nova do free5GC pode ter corrigido itens ou mudado o código que os *patches* tocam. A V2 continua no free5GC v4.2.3; a lista só registra o que mudou.
- O `build-wolfssl.sh` baixa o wolfSSL do site do fabricante, com o SHA-256 fixado. Se o arquivo sair do ar, a instalação limpa para nesse passo.

**O tutorial ganha.** Nada: o guia público é a reescrita do tutorial, em inglês.

---

## 5. Depois do MVP-6

**MVP-7 — Repositório público.** Aprovado em 20/10/2026; passou para a seção 4.

**MVP-8 — Ua\* da norma e autorização por aplicação no UE.** Aprovado em 20/10/2026; está na seção 4.

**MVP-9 — CAPIF no free5GC** (pedido do responsável em 19/10/2026, na confirmação do P12 do MVP-6; MVP-8 até a renumeração de 20/10/2026). O AF externo passa a obter o *token* do NEF pelo caminho da norma, e não do NRF.

- `[3GPP]` TS 29.522 §7.2 e TS 33.122 §6.5.2.3: com CAPIF, o CAPIF Core Function (CCF) é o servidor de autorização. O AF, como *API invoker*, pede o *token* em `{apiRoot}/capif-security/v1/securities/{securityId}/token` (`Obtain_Authorization`, TS 29.222 §5.6.2.3). O NEF, como AEF, confere no *token* o próprio identificador, o nome da API e o escopo.
- `[FREE5GC]` O free5GC não tem CCF (F13 de `docs/free5gc-side/correcoes-free5gc.md`).
- Escopo mínimo a detalhar na proposta do MVP-9:
  - um CCF com a entrada do *API invoker* (`CAPIF_API_Invoker_Management_API`), a descoberta (`CAPIF_Discover_Service_API`), a segurança (`CAPIF_Security_API`) e a publicação (`CAPIF_Publish_Service_API`);
  - o NEF como APF, publicando `3gpp-akma`, e como AEF, conferindo os *tokens* do CCF;
  - o AF como *API invoker*.
- Normas: TS 23.222 V19.10.0, TS 29.222 V19.8.0 (com os YAML em `specs/openapi/`) e TS 33.122 V19.4.0, em `specs/` desde 19/10/2026.
- Fecha a parte do DV-9 sobre a origem do *token*. O TLS entre o NEF e o AF continua com o DV-9.

**MVP-12 — Renovação da K_AKMA pelo AAnF** (decisão do responsável de 22/10/2026: MVP próprio, depois da fase (d) do MVP-8). Aprovado em 28/10/2026; está na seção 4.

**MVP-13 — V2 do repositório público** (pedido do responsável em 28/10/2026: logo depois do MVP-12, antes dos MVP-9 a 11). Uma nova exportação do repositório público, com o que o projeto fez depois do MVP-7 e as recomendações que o responsável vai enviar ao free5GC e à 3GPP ao fim da implementação. Proposto em 28/10/2026; a proposta, com os pontos L1 a L12, está na seção 4.

- **Código e guia.** Os perfis Ua\* do MVP-8 (TLS, DTLS e OSCORE, com o serviço de chaves do ME) e a renovação da K_AKMA do MVP-12, com os *patches* novos e o roteiro de testes.
- **Recomendações ao free5GC.** `docs/free5gc-fixes.md`, em inglês, com tudo de `docs/free5gc-side/correcoes-free5gc.md`: cada mudança com a cláusula da norma, a forma que ela pede e o *patch* do projeto, ou o motivo de estar aberta.
- **Recomendações à 3GPP.** Um documento público novo, em inglês, com os defeitos de `docs/defects/registro-de-defeitos-3gpp.md` e as divergências de enumeração de TS 23.502. Cada um com a especificação e a versão, a cláusula, o texto atual, o problema e a mudança proposta, que é o que um *Change Request* precisa.
- O envio ao free5GC e à 3GPP fica com o responsável. O escopo detalhado, a versão (*tag*) e os critérios, como a instalação limpa só pelo guia do MVP-7, saem da proposta do MVP-13.
- **Nome** (decisão do responsável de 28/10/2026): MVP-13 neste repositório; no repositório público, só V2, sem a numeração dos MVP.

**Demais itens.** A lista original desta seção, revista em 19/10/2026. A ordem é uma sugestão: cada MVP é proposto e aprovado à parte.

| Item | Situação | Proposta |
|---|---|---|
| Acesso anônimo (§6.2.2, `anonInd`) | AAnF pronto no MVP-1 | MVP-6, no AF e no NEF |
| `remove-context` (§6.6) | operação pronta no MVP-1; a norma nomeia o OAM como consumidor, não uma NF do núcleo | MVP-6, só documentação |
| Caminho via NEF (§6.3, `3gpp-akma`) | — | MVP-6 |
| Filtro `routing-indicator` para `AANF` no NRF | feito no MVP-1b | sai da lista |
| Ua\* conforme K.2 sobre TLS; TLS 1.3 PSK conforme K.3.3 | Release 17 (AKMA_TLS) | MVP-8 (seção 4), com o DTLS e o OSCORE; exige rever o ADR-0003 (identificadores de H.3 ao lado da faixa privada) |
| GPSI via `id-translation-result`; `AKMA_GPSI_Support` | o UDM v4.2.3 só traduz GPSI → SUPI | MVP-10: *patch* no UDM; GPSI ao AF externo pelo NEF e, com a *feature*, ao AF interno |
| TLS nas interfaces SBI (§4.4.0) | DV-1 e DV-9 | MVP-11, junto com a autenticação do AF |
| Pontos abertos do free5GC: escopos por tipo de consumidor (F4), `exp` e `aud` no `VerifyOAuth` (F5), *webconsole* (F9), 401 do NEF (F12), descadastro do NEF (F14), `complexQuery` (F15) | `docs/free5gc-side/correcoes-free5gc.md` | F4 e F5 junto com o TLS nas SBI; os demais com o *upstream* |
| NFKC no UE (DV-4); SUPI dos tipos NAI, GLI e GCI (DV-3) | desvios registrados | sob demanda |
| *Authentication Proxy* (TS 33.535 §4.7) | opcional na norma | sob demanda, depois do MVP-8 |
| *Upstream* | — | *rebase* para a versão corrente (D4); modelos do Naanf_AKMA gerados pela mesma ferramenta OpenAPI do free5GC; `shared/akma` migrado para `free5gc/util`. O repositório público (MVP-7) é o primeiro passo, e a V2 (MVP-13) traz a lista de recomendações; o *pull request* pode vir antes dos MVP seguintes: depende de quando o projeto quer abrir a conversa com o free5GC |
| *Roaming* (`ROAMING_STATUS`, `RoamingRestriction`, `Notify` com `3gpp-Sbi-Callback`); `ServiceDisableNotification` (§6.8) | Release 18. O *roaming* exige uma VPLMN, que o laboratório não tem. O `AuthTrigger` (§6.4.4) passou a MVP-12 | depois do *upstream*, sob demanda |
| Várias instâncias de AAnF; persistência do contexto; desempenho (ADR-0001-R2 §5.4); verificação formal | — | sob demanda |

---

## 6. Desvios e escolhas entre textos normativos

Registro oficial, com motivo e impacto: `docs/conformidade/desvios.md`. Cada um é implementado no MVP indicado e, então, marcado como implementado naquele registro.

| Nº | O que a norma exige | O que o projeto faz | Estado (16/10/2026) |
|---|---|---|---|
| DV-1 | TS 33.535 §4.4.0: interfaces AAnF↔AUSF, AAnF↔AF/NEF e AAnF↔UDM *"shall be confidentiality, integrity and replay protected"*; §4.2.2 e TS 29.535 §4.2.1.3.2: o AF *"shall be authenticated and authorized"* antes de receber K_AF. TS 33.501 §13.1.0: TLS dentro da PLMN *"unless network security is provided by other means"*; NOTA 2: com interfaces confiáveis, a proteção criptográfica é decisão do operador | `http` (h2c) em AUSF↔AAnF (*loopback*) e em AF↔AAnF (rede virtual entre VMs), declarando ambas como interfaces confiáveis do laboratório. Consequência: o AF não é autenticado pelo transporte, e um AF pode apresentar o AF_ID de outro (ataque de falsificação de AF, Akman et al. 2023). Mitigação no MVP: a política local associa cada AF_ID ao endereço de origem do AF; risco residual declarado. TLS no MVP-6+ | **implementado** (17/10/2026) · MVP-1 |
| DV-2 | TS 33.535 §6.2.1 passo 5 (*"shall"*) — questão 1. TS 29.535 condiciona a reação ao *roaming* à *feature* opcional `RoamingRestriction` (Tab. 5.1.8-1, 5.1.7.3-1) | subscrição de `ROAMING_STATUS` implementada como etapa configurável, desabilitada (orientação do material de partida; registrado em `desvios.md`) | **implementado** (17/10/2026), parcial: o campo só aceita `false` · MVP-1 |
| DV-3 | TS 33.501 A.7.0 define P0 também para SUPI do tipo NAI, GCI e GLI | apenas SUPI do tipo IMSI | aprovado com o plano · MVP-0 |
| DV-4 | TS 33.220 B.2.1.2: NFKC nos parâmetros de cadeia | sem NFKC no lado UE (C++): FQDN só ASCII | aprovado com o plano · MVP-4 |
| DV-5 | `[DIVERGÊNCIA]` TS 24.109 K.2 c) manda usar HTTP Digest *"(see RFC 3310)"* — o Digest AKA, em que a senha é derivada do RES. O mesmo K.2 b) remete à §5.2, que usa RFC 7616 (§5.2.2.1); TS 33.535 B.1.2.2 segue TS 33.222 §5.3.0, que também usa RFC 7616 e desaconselha MD5 (NOTA 1b). Com senha = base64(K_AF), só RFC 7616 é coerente | Digest conforme RFC 7616, `algorithm=SHA-256`, `qop=auth` (`auth-int` opcional, §5.2.6), `rspauth` obrigatório na resposta. Registrar a citação de RFC 3310 como defeito candidato (I18) | aprovado com o plano · MVP-3 |
| DV-6 | `[DIVERGÊNCIA]` Divergência E. TS 33.535 B.1.2.2 (Stage 2): o AF inclui `3GPP-bootstrapping-akma` no `WWW-Authenticate`, e o `realm` contém o FQDN do AF (passos 3 e 4); TS 24.109 K.2 b)2) repete a constante; TS 24.109 §5.2.2.1 dá a forma `<constante>@<FQDN>`. Só TS 24.109 K.2 c)3) usa `3gpp-akma@<FQDN>` | Segue TS 33.535: `realm = "3GPP-bootstrapping-akma@<FQDN do AF>"` — constante e FQDN de TS 33.535, separador `@` do Stage 3; primeira parte configurável. Diverge só de K.2 c)3). Registro de defeitos (D-C) e Adendo D a atualizar | **decidido** · MVP-3 |
| DV-7 | `[DIVERGÊNCIA]` TS 29.535 §4.2.2.3.2: 204 *"if the requested AKMA Application Key information for the UE does not exist"*, mas sem K_AKMA a resposta é 403 `K_AKMA_NOT_PRESENT`, e com K_AKMA a K_AF é sempre derivável — não há condição distinta para o 204 | não emitir 204; 403 `K_AKMA_NOT_PRESENT` quando não há K_AKMA. Registrar como defeito candidato (CT3) | **implementado** (17/10/2026) · MVP-1 |
| DV-8 | TS 33.501 §13.4.1.1.2 passo 2: o produtor verifica `exp` e `aud` do token. `[FREE5GC]` `oauth.VerifyOAuth` não verifica nenhum dos dois; o NRF v1.4.5 emite `aud` vazio | o AAnF verificava `exp` e tolerava `aud` vazio; com o NRF corrigido (MVP-1b), exige `aud` igual ao seu ID ou a `AANF` | **encerrado** em 17/10/2026 (MVP-1b) · MVP-1 |
| DV-9 | TS 33.501 §12.2 e §12.3: entre o NEF e o AF externo, TLS e autenticação mútua por certificados *"shall"*; §12.4: autorização OAuth *"shall"* | `http`, *token* OAuth2 do NRF para `3gpp-akma` e política local AF_ID ↔ endereço de origem no NEF, no padrão das APIs do NEF do free5GC. Sem TLS nem certificado de cliente | **decidido** (19/10/2026) · MVP-6 |

Não são desvios, e assim serão documentados: identificador Ua\* da faixa privada (TS 33.220 H.2 NOTA 2, ADR-0003); seleção do AAnF por configuração local (TS 33.535 §6.7); autorização do AF por política local em vez de token OAuth2 (TS 33.535 §6.2.1 passo 2; TS 29.535 §5.1.9: OAuth2 *"may"* ser usado); contexto em memória; GPSI não obtido (condicional à política local, §6.2.1 passo 3); A-KID transmitido a cada AF guardado mesmo sem subscrição de *roaming* (§6.2.1 passo 8 só o exige com a subscrição).

---

## 7. Pontos a conferir na norma antes de implementar

| Nº | Ponto | Estado | MVP |
|---|---|---|---|
| P1 | Resposta de `retrieve-applicationkey` para A-KID sem K_AKMA | **Resolvido.** 403 `K_AKMA_NOT_PRESENT` (TS 29.535 §4.2.2.3.2). O 204 da mesma cláusula não tem condição alcançável: DV-7 | — |
| P2 | Resposta de `register-anchorkey` | **Resolvido.** 200 com `AkmaKeyInfo` no corpo (TS 29.535 §4.2.2.2.2; YAML) | — |
| P3 | Causa e código para AF não autorizado | **Resolvido.** Nenhuma causa definida em TS 29.535 nem em TS 29.500 Tab. 5.2.7.2-1 para o caso; 403 com `ProblemDetails` opcional (Tab. 5.1.4.3.2-2). `[INTERPRETAÇÃO]` 403 com `ProblemDetails` sem `cause` | — |
| P4 | Escopos OAuth2 por operação | **Resolvido.** TS 33.501 §13.4.1.0; TS 29.535 §5.1.9 (OAuth2 *"may"*, por configuração local) | — |
| P5 | Tempo de vida de K_AF | **Resolvido.** TS 33.535 §4.4.0: tempo de vida máximo pela política de autenticação do operador; §5.2: K_AF pode continuar após nova K_AKMA até o fim da sessão ou da validade. Configurável (`akma.kafLifetime`) | — |
| P6 | Interfaces confiáveis | **Resolvido.** TS 33.501 §13.1.0 NOTA 2; DV-1 | — |
| P7 | D10 | depende de TS 29.505 (D8) | MVP-2 |
| P8 | Afirmações atribuídas a TS 33.535 e TS 29.535 | **Conferidas** em 16/10/2026 (`conferencia-normativa.md`) | — |
| P9 | Forma textual do identificador Ua\* dentro do `AfId` | **Resolvido** em 16/10/2026: dez caracteres hexadecimais, **maiúsculas na escrita** (como o exemplo de TS 29.522), as duas caixas na leitura | MVP-0 |
| P10 | Restrição de TS 33.220 B.2.0 (P0 *"static ASCII-encoded string"*) | **Resolvido.** TS 33.535 A.4 não a adota | — |

---

## 8. Inconsistências encontradas na documentação existente

| Nº | Onde | O que | Proposta |
|---|---|---|---|
| I1 | caminhos | O material de partida cita `~/akma/docs/`; os documentos estão em `~/akma/akma-free5gc/docs/`. Há cópias idênticas soltas em `~/akma/` (README, ADR-0001, matriz, registro de defeitos) e um `.tar.gz` idêntico | resolver com D2 |
| I2 | `specs/`, `papers/` | normas e artigos estavam em `reference/`, fora do `.gitignore` | **resolvido** (D1); falta só TS 29.505 (D8) |
| I3 | letras dos defeitos | As questões usam A–E (A = `atid`, B = codificação de chaves, C = codificação do A-TID, D = `trigger-auth`, E = `realm`); o registro de defeitos usa D-A a D-E com outra atribuição (D-B = `trigger-auth`, D-C = `realm`, D-D = `pattern`, D-E = separador `";"`). "D" designa coisas diferentes nos dois esquemas | tabela de correspondência no registro de defeitos |
| I4 | `00-baseline-normativo.md` §2.5 | a cadeia `S` do A-TID tem um espaço espúrio (`412d54494 4`); o oráculo imprime `412d544944` | corrigir |
| I5 | `00-baseline-normativo.md` §6.2 | prevê acrescentar `AkmaInd` e `RoutingId` ao `openapi`; já existem | anotar como superado |
| I6 | README | diz que o código vive em outro repositório | resolver com D2 |
| I7 | `scripts/bootstrap-repo.sh` | faz `git commit` com mensagem em português e cria repositório **público** no GitHub; conflita com a regra de commits do projeto | não executar; revisar ou remover |
| I8 | ADR-0001 | cita fases F11, F12 e F15, inexistentes na tabela de fases do README | correspondência F → MVP da seção 2 |
| I9 | Adendo B §B.1.1 | usa a marcação `[HIPÓTESE]`, fora do conjunto de categorias | já resolvido pela verificação; anotar |
| I10 | `reference/kdf/akma_kdf.py` | `kdf()` anotada como `-> bytes`, devolve tupla | corrigir no MVP-0 |

Encontradas na conferência contra as normas (detalhe e citações em `docs/baseline/conferencia-normativa.md`):

| Nº | Onde | O que | Proposta |
|---|---|---|---|
| I11 | Adendo D §D.1, ADR-0002 §3.1, registro de defeitos D-D | O precedente citado para `kAkma`/`kaf` é só TS 29.509 (padrão não ancorado). `TS29503_Nudm_UEAU.yaml` define `Kausf` como `^[A-Fa-f0-9]{64}$` — ancorado, na interface UDM→AUSF, e é a chave da qual K_AKMA deriva | acrescentar como precedente principal; a decisão (hexadecimal) não muda, fica mais forte |
| I12 | `00` §5.1; material de partida §4.9 | Atribui a TS 33.501 §6.2.2 a guarda de K_AUSF "quando AKMA é suportado". A menção a AKMA está em §6.1.1.1; §6.2.2.1 define quando guardar/substituir K_AUSF | corrigir as citações |
| I13 | Adendo A §A.5; material de partida §4.9 | Cita "TS 29.503 §5.2.3.5.2" para `immediateFlag`; a cláusula é §5.5.2.2.2, e o texto diz que o UDM ***may*** incluir o estado corrente *"if available"*. "Torna o passo 6 implementável de forma síncrona" é forte demais: não há garantia | corrigir; pesa na questão 1 (fora do MVP) |
| I14 | Adendo B §B.3 | Lista 4 cabeçalhos obrigatórios; a Tab. 5.2.3.2.1-1 de TS 29.500 V19.7.0 tem 21 entradas (a maioria para comunicação indireta) | corrigir a lista |
| I15 | `00` §2.4; ADR-0003 | A tabela de H.3 omite `(0x01,0x00,0x00,0x00,0x02)` — HTTP Digest conforme TS 24.109 — e outros. O ADR-0003 não registra por que essa alternativa não foi adotada | completar a tabela; acrescentar a alternativa ao ADR-0003 |
| I16 | `00` §2.4 item 2, matriz | Afirma como `[3GPP]` que o identificador Ua\* no `AfId` é representado "por dez dígitos hexadecimais"; o texto normativo diz *"a string of 5 octet"* e a forma hexadecimal vem do exemplo | reclassificar como `[INTERPRETAÇÃO]` (P9) |
| I17 | material de partida §4.6 | Registra "204 se não existir" e 403 `K_AKMA_NOT_PRESENT` para `ApplicationKey_Get` sem notar que o 204 não tem condição alcançável (o mesmo ocorre em TS 29.535 e na API do NEF em TS 29.522) | DV-7 |
| I18 | Adendo D §D.4.1 | Registra K.2 com Digest "RFC 3310" sem notar o conflito com TS 24.109 §5.2.2.1 e TS 33.222 §5.3.0 (RFC 7616) | novo defeito candidato (CT1), ligado a DV-5 |
| I19 | TS 29.522 (editorial) | Tab. 5.14.5.1-1 remete `ServiceDisableNotif` a "5.15.5.3.4" (é 5.14.5.3.4); "table 514.5.2-1" | anotar no registro de defeitos, prioridade baixa |

Encontradas com TS 33.535, TS 29.535, TS 23.502 e os artigos (16/10/2026):

| Nº | Onde | O que | Proposta |
|---|---|---|---|
| I20 | registro de defeitos D-C; Adendo D §D.4.1 | A decisão sobre o `realm` (seguir K.2 c)3)) não considerou TS 33.535 B.1.2.2, que manda incluir `3GPP-bootstrapping-akma` no `WWW-Authenticate` | decisão revista (DV-6); atualizar D-C e o Adendo D |
| I21 | `00` §6.3; material de partida §4.8 | Apresentam a guarda do A-KID transmitido a cada AF como parte incondicional do contexto; TS 33.535 §6.2.1 passo 8 a liga à subscrição de *roaming* | anotar; manter a guarda no MVP |
| I22 | README; `00` §2.5; `reference/kdf/README.md` | "AKMA+" cita o trabalho de Yang et al. (Zenodo 10.5281/zenodo.14649847); existe outro "AKMA+", de Boureanu et al. (ASIA CCS '25), com proposta diferente | identificar autores ao citar |
| I23 | material de partida §6.4 | "15 artigos": são 15 arquivos e 14 trabalhos (dois arquivos idênticos) | anotar |
| I24 | TS 33.535 §4.3.0 NOTA 1, §7.4 × TS 23.502 V19.9.0 §5.2.3.5.2 | A TS 33.535 aponta `Nudm_EventExposure_Subscribe` como serviço do UDM para o AKMA; a TS 23.502 não lista o AAnF como consumidor | acrescentar ao registro das divergências de enumeração |
| I25 | TS 29.535 (editorial) | `CtxRemove` sem `required` no YAML; `info.title` `3gpp-akma`; "retrieval notification" nas respostas 307/308 de `register-anchorkey` | registro de defeitos, prioridade baixa |
| I26 | `especificacoes.md` | TS 23.502 fixada em V19.8.0; o acervo tem V19.9.0 | **resolvido** (D12): baseline mantido, nota e Adendo E |

Confirmado sem alteração (detalhe em `docs/baseline/conferencia-normativa.md` §2): todo o fluxo de TS 33.535 §6.1 e §6.2.1, os Anexos A e B, e o Stage 3 de TS 29.535 com seu YAML, tal como o material de partida os descreve, salvo os itens acima; KDF e alocação de FC (TS 33.220 B.2); faixa privada `0xFF` (H.2 NOTA 2); B-TID e TMPI em base64 (TS 33.220 §4.5.2, B.4); NOTA 0 de TS 33.222 §5.4.0.2; separador `";"` (TS 24.109 K.3.3); senha base64(K_AF) (K.2 c)2)); divergências A e D; divergência E confirmada, com a decisão reaberta (DV-6); RID e domínio doméstico (TS 23.003); ausência de AKMA em TS 23.003 e TS 33.222; `akmaInd`/`routingId` e obrigatoriedade em TS 33.501; tipos de TS 29.522, TS 29.571 e TS 29.122; `AANF`, `naanf-akma`, `AanfInfo` e `routing-indicator` em TS 29.510; códigos obrigatórios de POST (TS 29.500); custom operations (TS 29.501); `id-translation-result`, `RoamingStatusReport` e `trigger-auth` em TS 29.503; números de CR citados.

**Situação em 16/10/2026 — correções aplicadas** (errata marcada no texto de cada documento; resumo no `CHANGELOG.md`, versão 2.1):

| Itens | Onde | Estado |
|---|---|---|
| I1, I6 | README | repositório único em `~/akma/akma-free5gc/` (D2). As cópias soltas em `~/akma/` (README, ADR-0001, matriz, registro de defeitos) e o `.tar.gz` ficam fora do repositório; removê-las fica a critério do responsável |
| I2, I26 | — | resolvidos por D1 e D12 |
| I3, I18–I20, I24, I25 | registro de defeitos | tabela de correspondência das letras; D-C revisto; D-D complementado; D-F a D-I acrescentados; inconsistência TS 33.535 × TS 23.502 registrada |
| I4, I5, I12, I15, I16, I21, I22 | `00-baseline-normativo.md` | corrigidos |
| I8 | ADR-0001 | estado "revisto pelo ADR-0001-R2" e correspondência das fases F → MVP |
| I9, I14 | Adendo B | corrigidos |
| I10 | `reference/kdf/akma_kdf.py` | anotação de tipo corrigida e comentários em inglês; saída idêntica à anterior |
| I11 | Adendo D, ADR-0002, registro D-D | precedente `Kausf` acrescentado |
| I13 | Adendo A | corrigido |
| I15 | ADR-0003 | §3.4 acrescentada |
| I16, I21 | matriz de rastreabilidade (revisão 3.1) | corrigidos, junto com as consequências do ADR-0001-R2 e as decisões DV-1 a DV-7 e D9 |
| I7 | `scripts/bootstrap-repo.sh` | removido |
| I17, I23 | material de partida | o texto de partida não está no repositório; as correções ficam registradas em `conferencia-normativa.md` e neste plano |

ADR-0002 e ADR-0003 passaram ao estado "Aceito" com a aprovação deste plano.

---

## 9. A pergunta-guia

> É possível executar o fluxo básico de AKMA em um free5GC modificado, derivar corretamente o material de chave conforme o 3GPP, e permitir que uma aplicação obtenha sua chave de aplicação através do AAnF?

| Parte da pergunta | Evidência | MVP |
|---|---|---|
| fluxo básico num free5GC modificado | autenticação primária real (UERANSIM e banco de testes Go) seguida de `AnchorKey_Register` — **respondido** em 18/10/2026 ([R2] a [R6] de `docs/af-side/e2e.md`) | MVP-2, MVP-5 |
| derivação correta conforme o 3GPP | concordância octeto a octeto entre Python, Go, C++ e a KDF do free5GC; A-KID igual no UE e no AAnF — **respondido** (vetores do MVP-0, [S1], [R2.akid-aanf]) | MVP-0, MVP-4 |
| aplicação obtém K_AF pelo AAnF | AF recebe K_AF com `expiry` e autentica o UE pelo Ua\* — **respondido** ([R2.ua], [N01b]) | MVP-3, MVP-5 |

Nenhuma etapa obrigatória da tabela de etapas mínimas é removida. A única etapa de estado aberto (`ROAMING_STATUS`) segue a orientação já dada e entra como desvio registrado (DV-2).
