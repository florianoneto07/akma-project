# Ambiente verificado

**Data da verificação:** 15 de outubro de 2026
**Escopo:** fixação da base de código do free5GC e do ambiente das duas VMs, e resposta à questão em aberto nº 5 (presença de `AANF`, `naanf-akma`, `AkmaInd` e `RoutingId` no `free5gc/openapi`).
**Método:** inspeção somente leitura, por SSH, da VM do free5GC; nenhuma alteração foi feita nela.

Este documento preenche a seção "Base de código" de `especificacoes.md`, que estava pendente desde o ADR-0001 §7.

---

## 1. Máquinas

| Item | VM free5GC | VM do AF (repositório) |
|---|---|---|
| Endereço | `192.168.16.130` (`ens33`) | `192.168.16.129` (`ens33`) |
| Nome do host | `ubuntu` | `akma` |
| Sistema | Ubuntu 18.04.6 LTS | — |
| Kernel (`uname -r`) | `5.4.0-150-generic` | `7.0.0-38-generic` |
| Go | `go1.26.2 linux/amd64` | **ausente** |
| git | 2.17.1 | **ausente** |
| Python | — | 3.12.3 |
| pandoc / python-docx | — | **ausentes** |
| Relógio | sincronizado (`systemd-timesyncd`), fuso America/Los_Angeles | sincronizado (NTP ativo), fuso America/Bahia |

`[FREE5GC]` Desvio de relógio medido entre as VMs por SSH: cerca de +0,16 s (remota adiantada). A medida é grosseira — inclui a assimetria da conexão SSH —, mas fica abaixo do critério de 1 s do ADR-0001 §7. Os fusos horários diferem; isso não afeta `expiry` desde que o instante seja serializado conforme RFC 3339 com deslocamento explícito (recomenda-se UTC, sufixo `Z`).

Latência ICMP entre as VMs: cerca de 0,5 ms.

---

## 2. free5GC

| Item | Valor |
|---|---|
| Repositório | `https://github.com/free5gc/free5gc.git` |
| Versão | **v4.2.3** |
| *Commit* | `3b34a08e93a9b334f0f4005d3a3a9f79b66d59b9` (2026-06-30, "chore: update submodule") |
| Estado da árvore | um arquivo modificado localmente: `cert/nrf.pem` |
| Teste de instalação | `TestRegistration` passou (registro em `~/free5gc-setup/test-registration.log`) |

### 2.1 Submódulos

| Submódulo | Versão | *Commit* |
|---|---|---|
| `NFs/amf` | v1.4.5 | `2962b9fdbdff1fb04de906063438e7acce761ce4` |
| `NFs/ausf` | v1.4.5 | `b74ee7dafec5975e3440ae95125994565cfb0496` |
| `NFs/bsf` | v1.0.3 | `1754bde834ed8e8eed822a555b4f84ac53eebcff` |
| `NFs/chf` | v1.2.5 | `3564da89afc476f2414c66aa34d9df1574b95449` |
| `NFs/n3iwf` | v1.3.5 | `f149cdca714cbec8236586db636c581cf63ccc8c` |
| `NFs/nef` | v1.2.5 | `6e27bc78a4696835da6d99d065696690701ba995` |
| `NFs/nrf` | v1.4.5 | `8e567d4c161c95dfb5765b5399b6b2086c64985a` |
| `NFs/nssf` | v1.4.4 | `cb7305b08e19297447fb5765c9db6e046076b4c8` |
| `NFs/pcf` | v1.4.5 | `90d88f15fe18780bb6c4b2cf3bbd59d636dd130d` |
| `NFs/smf` | v1.4.5 | `796a2de8cfedc208e6b709da9e8fe5a93026434e` |
| `NFs/tngf` | v1.1.4 | `b267d8c391148e26eddddfd9e81d00440bc12a08` |
| `NFs/udm` | v1.4.5 | `bdcc8bc316bb4e2a525233289d30f9945c844ea5` |
| `NFs/udr` | v1.4.4 | `b649f90d39a99007dcbf081fcb3479477e28b6eb` |
| `NFs/upf` | v1.2.12 | `0012cebc4c3758882a89f6485d1d0d95ec3659c7` |
| `webconsole` | v1.4.5 | `ea04cadd91550fadd5c90939cd5fafe86c253b6c` |

### 2.2 Módulos Go compartilhados

| NF | `free5gc/openapi` | `free5gc/util` | diretiva `go` |
|---|---|---|---|
| ausf, udm, udr, nrf, nef | v1.2.4 | v1.3.2 | 1.26.2 |
| amf | v1.2.5-0.20260527003827-02dc71b4d94f | v1.3.2 | 1.26.2 |

Nenhum `go.mod` contém diretiva `replace`.

### 2.3 Outros componentes na VM do free5GC

| Componente | Versão | Observação |
|---|---|---|
| `gtp5g` | v0.9.14 | módulo carregado (`lsmod`) |
| UERANSIM | `v3.2.6-66-g85a0fbfd` (`85a0fbfdd166da84c39ed5b4a5435afd40725878`) | compilado em `~/UERANSIM/build`; **instalado na VM do free5GC**, com o UE em *network namespace* próprio (`ueransim-ue`, script `~/free5gc-setup/05-ue-netns.sh`) |

`[FREE5GC]` Configuração do UE em `~/UERANSIM/config/free5gc-ue.yaml`: `supi: imsi-208930000000001`, `mcc: 208`, `mnc: 93`, `routingIndicator: '0000'`, `protectionScheme: 0`. O RID configurado tem quatro dígitos zero — exatamente o caso de teste de zeros à esquerda do Adendo A §A.2.

### 2.4 Endereçamento SBI corrente

| NF | `registerIPv4` / `bindingIPv4` | Porta |
|---|---|---|
| NRF | `127.0.0.10` | 8000 |
| AUSF | `127.0.0.9` | 8000 |
| UDM | `127.0.0.3` | 8000 |

`[FREE5GC]` `nrfcfg.yaml` tem `oauth: true`: as NFs obtêm token OAuth2 do NRF e o validam nas requisições recebidas. O NRF escuta apenas em *loopback*; não é alcançável a partir da VM do AF.

`[FREE5GC]` `ausfcfg.yaml`: `plmnSupportList` com 208/93 e 123/45; `eapAkaSupiImsiPrefix: false` — o AUSF já remove o prefixo `imsi-` na identidade do EAP-AKA'.

---

## 3. Questão em aberto nº 5 — resultado

### 3.1 `free5gc/openapi` v1.2.4 (versão usada por AUSF, UDM, UDR e NRF)

| Elemento | Presente? | Onde |
|---|---|---|
| `NFType` `AANF` | **sim** | `models/model_nf_type.go` (`NFType_AANF`), `models/model_nrf_nf_management_nf_type.go` (`NrfNfManagementNfType_AANF`) |
| `ServiceName` `naanf-akma` | **sim** | `models/model_service_name.go` (`ServiceName_NAANF_AKMA`) |
| `AanfInfo.routingIndicators` | **sim**, `[]string` | `models/model_aanf_info.go` |
| `NFProfile.aanfInfoList` | **sim**, `map[string]AanfInfo` | perfis de NFManagement e de NFDiscovery |
| `AuthenticationInfoResult.akmaInd` | **sim**, `bool` | `models/model_udm_ueau_authentication_info_result.go` (`AkmaInd`) |
| `AuthenticationInfoResult.routingId` | **sim**, `string` | idem (`RoutingId`) |
| Parâmetro de descoberta `routing-indicator` | **sim** | `nrf/NFDiscovery/api_nf_instances_store.go` |
| `AuthenticationSubscription.akmaAllowed` (UDR) | **sim**, `bool` | `models/model_authentication_subscription.go` |
| `AuthenticationSubscription.routingId` (UDR) | **sim**, `string` | idem |
| Modelos do Naanf_AKMA — `AkmaKeyInfo`, `CtxRemove`, `AkmaAfKeyRequest`, `AkmaAfKeyData`, `ServiceDisableNotif` | **não** | — |
| Cliente/servidor do Naanf_AKMA | **não** | não há diretório `aanf/` |

`[FREE5GC]` Os modelos declaram como fonte TS 29.503 V17.10.0 (API 1.2.2) e TS 29.504 V17.12.0: os modelos correspondem ao Release 17. `RoutingId` e `AanfInfo.RoutingIndicators` são cadeias, o que preserva zeros à esquerda.

**Consequência para o escopo:** o *patch* no `free5gc/openapi` **não** precisa acrescentar `AANF`, `naanf-akma`, `AkmaInd`, `RoutingId` nem `AanfInfo`. Precisa apenas dos modelos (e, se adotado, do cliente) do Naanf_AKMA. A questão nº 5 fica encerrada.

### 3.2 Versões mais recentes publicadas

| Módulo | Versão mais recente | Data | AKMA |
|---|---|---|---|
| `free5gc/free5gc` | v4.3.0 | 2026-09-30 | — |
| `free5gc/openapi` | v1.3.0 | 2026-08-19 | mesmos elementos de 3.1, com modelos **renomeados** (ex.: `Nrf_NFMgmt_NFType_AANF`, `model_udr_dr_authentication_subscription.go`); acrescenta `ServiceName` `3gpp-akma`; continua sem modelos do Naanf_AKMA |
| `free5gc/ausf` | v1.4.6 | 2026-09-04 | nenhuma referência a AKMA |
| `free5gc/udm` | v1.4.6 | 2026-08-28 | nenhuma referência a AKMA |

`[FREE5GC]` A renomeação dos modelos em `openapi` v1.3.0 significa que *patches* escritos contra v1.2.4 não se aplicam diretamente sobre v4.3.0. Ver decisão pendente sobre a versão-alvo no plano de desenvolvimento.

---

## 4. Outras constatações sobre o código do free5GC

### 4.1 NRF

`[FREE5GC]` `nrf/internal/sbi/processor/nf_profile_validation.go` aceita `AANF` como tipo de NF registrável.

`[FREE5GC]` `nrf/internal/sbi/processor/nf_discovery.go` (linhas 1006–1041) aplica o filtro `routing-indicator` apenas quando `target-nf-type` é `AUSF` ou `UDM`. Para `AANF` o parâmetro é **ignorado**: a descoberta devolve todas as instâncias de AAnF. Com uma única instância de AAnF o resultado é correto; com várias instâncias segmentadas por RID, não. TS 29.510 prevê o parâmetro para `AANF`.

### 4.2 AUSF (`ausf` v1.4.5)

`[FREE5GC]` Pontos de integração confirmados em `internal/sbi/processor/ue_authentication.go`: `UeAuthPostRequestProcedure` (recebe `UdmUeauAuthenticationInfoResult`), `Auth5gAkaComfirmRequestProcedure` (sucesso do 5G-AKA), `EapAuthComfirmRequestProcedure` (sucesso do EAP-AKA'). `AusfUeContext` (`internal/context/context.go`) guarda `Kausf` como cadeia hexadecimal; no EAP-AKA', `Kausf` é extraída de `EMSK[0:32]`. O AUSF consome hoje apenas NRF e UDM (`internal/sbi/consumer/`).

### 4.3 UDM (`udm` v1.4.5)

`[FREE5GC]` `internal/sbi/processor/generate_auth_data.go` monta `UdmUeauAuthenticationInfoResult` sem `AkmaInd` nem `RoutingId`, e já lê do UDR o `AuthenticationSubscription` do assinante — o modelo que contém `akmaAllowed`.

`[FREE5GC]` `pkg/suci/suci.go` analisa o SUCI e captura o Routing Indicator como cadeia (`Suci.RoutingIndicator`, regex `\d{1,4}`), mas a função de análise não é exportada; `ToSupi` devolve apenas o SUPI.

### 4.4 `ueauth.GetKDFValue` (`free5gc/util` v1.3.2) confrontada com TS 33.220 B.2.0

| Ponto verificado | Resultado |
|---|---|
| Saída integral de 256 bits ou truncada | **Integral.** `hmac.New(sha256.New, key)` e `Sum(nil)`: 32 octetos |
| `KDFLen` em dois octetos *big-endian* | **Sim.** `binary.BigEndian.PutUint16(uint16(len(input)))`. Sem verificação de limite: um parâmetro acima de 65535 octetos teria o comprimento truncado silenciosamente (TS 33.220 B.1 proíbe esse caso; não afeta o AKMA) |
| Normalização NFKC | **Não aplicada.** A função concatena os octetos recebidos; a codificação dos parâmetros é responsabilidade do chamador |
| Forma de `FC` | cadeia hexadecimal (`"6C"`); não há constantes para `0x80`–`0x82` |

`[INTERPRETAÇÃO]` A função é utilizável pelo AUSF para K_AKMA e A-TID, cujos parâmetros (`"AKMA"`, `"A-TID"`, dígitos do IMSI) são ASCII e, portanto, invariantes sob NFKC. Para K_AF — cujo P0 contém o FQDN do AF, possivelmente internacionalizado — o chamador deve aplicar NFKC antes; isso afeta o AAnF e o UE, não o AUSF.

`[3GPP]` Conferido em TS 33.220 V19.1.0: B.2.0 (`S = FC ‖ P0 ‖ L0 ‖ … ‖ Pn ‖ Ln`, HMAC-SHA-256), B.2.1.2 (UTF-8 + NFKC), B.2.1.3 (comprimento em dois octetos), B.1 (parâmetro de no máximo 65535 octetos), B.2.2 (FC `0x80`–`0x82` alocados a TS 33.535).

### 4.5 Banco de testes do free5GC

`[FREE5GC]` `test/ranUe.go` deriva `Kausf` em Go tanto no 5G-AKA quanto no EAP-AKA' (linhas 255–342). Os testes `TestRegistration` e `TestEAPAKAPrimeAuthentication` exercitam os dois métodos de autenticação primária sem UERANSIM.

---

## 5. Acervo normativo disponível

Atualizado em 15 e 16/10/2026, após a cópia das especificações para o repositório (hoje em `specs/`, D1). As versões abaixo foram conferidas na capa de cada documento e no campo `info.version`/`externalDocs` de cada arquivo OpenAPI; todas coincidem com `especificacoes.md`, exceto TS 23.502.

| Especificação | Arquivo | Estado |
|---|---|---|
| TS 23.003 V19.7.0 | `specs/23003-j70.docx` | presente |
| TS 24.109 V19.0.0 | `specs/24109-j00.docx` | presente |
| TS 29.122 V19.7.0 | `specs/29122-j70.docx` + `specs/openapi/TS29122_*.yaml` | presente |
| TS 29.500 V19.7.0 | `specs/29500-j70.docx` + `specs/openapi/TS29500_CustomHeaders.abnf` | presente |
| TS 29.501 V19.5.0 | `specs/29501-j50.docx` | presente |
| TS 29.503 V19.7.0 | `specs/29503-j70.docx` + `specs/openapi/TS29503_*.yaml` | presente |
| TS 29.509 V19.6.0 | `specs/29509-j60.docx` + `specs/openapi/TS29509_*.yaml` | presente |
| TS 29.510 V19.7.0 | `specs/29510-j70.docx` + `specs/openapi/TS29510_*.yaml` | presente |
| TS 29.522 V19.7.0 | `specs/29522-j70_{0,1,2}_*.docx` (Anexo A.12 = `TS29522_AKMA.yaml`) | presente |
| TS 29.571 V19.7.0 | `specs/29571-j70.docx` + `specs/openapi/TS29571_CommonData.yaml` | presente |
| TS 33.220 V19.1.0 | `specs/33220-j10.docx` | presente |
| TS 33.222 V19.0.0 | `specs/33222-j00.docx` | presente |
| TS 33.501 V19.7.0 | `specs/33501-j70.docx` | presente |
| TS 33.535 V19.0.0 | `specs/33535-j00.docx` | presente (16/10/2026) |
| TS 29.535 V19.5.0 e `TS29535_Naanf_AKMA.yaml` v1.2.0 | `specs/29535-j50.docx`, `specs/openapi/TS29535_Naanf_AKMA.yaml` | presentes (16/10/2026); o YAML declara V19.5.0 |
| **TS 23.502** | `specs/23502-j90.docx` | presente, na **V19.9.0 (2026-09)**; o baseline permanece V19.8.0 (Adendo E) |
| ETSI TR 121 917 V17.0.1 (TR 21.917 Rel-17) | `specs/tr_121917v170001p.pdf` | presente |
| Artigos | `papers/` (15 PDF + `INDICE.txt`) | presentes; os SHA-256 conferem com o índice; dois arquivos são idênticos (14 artigos distintos) |
| TS 29.505 | — | ausente (decisão D8 do plano) |

`.gitignore` exclui `/specs/` e `/papers/` (D1, 16/10/2026). As normas não são redistribuídas pelo repositório, como declara o README, e vários artigos são cópias licenciadas de editoras (o rodapé de cada página registra *"Authorized licensed use limited to"* a instituição que os baixou).

## 6. Pendências deste registro

| Item | Estado |
|---|---|
| TS 29.505 | ausente — ver D8 do plano |
| Versão de TS 23.502 | baseline mantido em V19.8.0; cópia do acervo V19.9.0 — Adendo E |
| git e Go na VM do AF | instalados para o MVP-0 (Go 1.26.2 em `~/sdk/go1.26.2`); gcc instalado em 17/10/2026 (detector de corridas) |
| Alcançabilidade da porta do AAnF a partir da VM do AF | **verificada** em 17/10/2026: `192.168.16.130:8000` responde a pedidos de `192.168.16.129`; UFW desabilitado na VM do free5GC (`/etc/ufw/ufw.conf`: `ENABLED=no`) — `docs/free5gc-side/aanf.md` §9 |
| `curl` na VM do AF | ausente; os testes de protocolo do AAnF usam h2c nos testes Go |
| Alterações na VM do free5GC depois desta verificação | MVP-1: arquivos novos do AAnF (`docs/free5gc-side/aanf.md` §1). MVP-1b: *patches* aplicados em `~/free5gc/NFs/nrf` e `bin/nrf` recompilado, com o original em `bin/nrf.orig` (`docs/free5gc-side/nrf.md`). MVP-2: *patches* em `NFs/udm`, `NFs/ausf`, `test/` e `config/ausfcfg.yaml`, `bin/udm` e `bin/ausf` recompilados (originais em `*.orig`), assinante `imsi-208930000000001` com `akmaAllowed`/`routingId` e repositório copiado para `~/akma-free5gc` (`docs/free5gc-side/{udm,ausf,assinatura-akma}.md`) |
