# Revisão das normas — segunda passada do MVP-14

**MVP:** MVP-14 · **Data:** 8 de novembro de 2026 · **Estado:** segunda passada feita, por leitura das normas e do código; nada foi executado no laboratório. Pedida pelo responsável no ponto Q15, com o Q1 aprovado (plano, MVP-14). Os achados 11 a 15 da seção 6 foram decididos na mesma data (pontos Q16 a Q20 do plano): os achados 11 a 14 entram no MVP-15, com o F31, e o 15 fica com o Q10.

Este documento é a tabela de conformidade do projeto: cada parágrafo com *shall* ou *should* das cláusulas que tocam o projeto, com o estado e a evidência. A primeira passada, com a revisão dos desvios, está no plano (MVP-14). As opções de texto para os defeitos da norma estão em `docs/defects/propostas-para-a-norma.md` (Q4).

## 1. Método

- **Escopo.** TS 33.535 inteira e, nas normas que ela chama, as cláusulas que o projeto usa (Q15): TS 33.501, TS 29.535, TS 29.522, TS 29.503, TS 24.109, TS 33.220, TS 33.222, TS 33.122, TS 23.222, TS 29.222, TS 33.210 e TS 33.310. Pelo Q15, entram também as normas que o projeto implementa e que a lista não trazia: TS 29.510, TS 29.500, TS 29.501, TS 29.122, TS 23.003, TS 23.502, TS 29.505 e TS 29.571. Cada seção diz as cláusulas lidas.
- **Versões.** A mais recente de cada norma, lida no arquivo da 3GPP em 06/11/2026 (plano, MVP-14, tabela de versões). Os textos da Release 20 estão em `specs/` e os YAML em `specs/openapi-r20/`, fora do controle de versão (Q1). TS 23.502 fica na V19.9.0, porque o arquivo da V20.3.0 passa do limite de 10 MB da ferramenta de download.
- **Extração.** Os parágrafos com *shall* ou *should* saíram do texto de cada norma, com a cláusula a que pertencem, pelo título mais próximo acima. Uma linha de tabela da norma conta como um parágrafo. Na coluna "Tipo", S é *shall* (ou *shall not*) e R é *should*; um parágrafo com os dois conta como S.
- **Leitura.** Cada parágrafo foi lido e comparado com o código do repositório, com as cópias do free5GC v4.2.3 e do UERANSIM em `free5gc/build` e `af/build`, e com as especificações dos MVPs. As evidências citam as verificações do `e2e.sh` (`docs/af-side/e2e.md`) e os testes Go.
- **Tabelas.** As quatro normas do AKMA (TS 33.535, TS 29.535, TS 29.522 e TS 24.109) têm uma linha por parágrafo. As outras agrupam, na mesma linha, os parágrafos seguidos de uma cláusula com o mesmo estado e a mesma evidência; toda linha **em parte**, **não feito** ou **desvio** traz o começo do primeiro parágrafo.
- **Estados.**
  - **feito**: o projeto faz o que o parágrafo pede, ou o comportamento pedido vem do free5GC sem mudança e o laboratório o exercita;
  - **em parte**: falta uma parte, dita na evidência;
  - **não feito**: o projeto não faz;
  - **desvio**: registrado em `desvios.md`;
  - **free5GC**: requisito de uma NF do free5GC que o projeto não altera nem revisa, como a autenticação primária;
  - **não se aplica**: GBA, SCP, SEPP, outra PLMN, SNPN, recursos opcionais que o laboratório não tem, ou frases para quem escreve a norma.

## 2. Resumo

| Norma | Parágrafos | Feito | Em parte | Não feito | Desvio | free5GC | Não se aplica |
|---|---|---|---|---|---|---|---|
| TS 33.535 | 97 | 73 | 11 | 4 | 7 | 0 | 2 |
| TS 29.535 | 57 | 38 | 2 | 8 | 3 | 0 | 6 |
| TS 29.522 | 40 | 20 | 4 | 11 | 0 | 0 | 5 |
| TS 24.109 | 35 | 23 | 0 | 0 | 3 | 0 | 9 |
| TS 33.501 | 126 | 49 | 3 | 0 | 1 | 63 | 10 |
| TS 29.503 | 81 | 11 | 0 | 34 | 0 | 9 | 27 |
| TS 33.220 | 38 | 25 | 0 | 0 | 1 | 0 | 12 |
| TS 33.222 | 49 | 33 | 1 | 0 | 0 | 0 | 15 |
| TS 33.122 | 110 | 72 | 4 | 8 | 0 | 0 | 26 |
| TS 23.222 | 9 | 6 | 0 | 3 | 0 | 0 | 0 |
| TS 29.222 | 117 | 51 | 0 | 49 | 1 | 0 | 16 |
| TS 33.210 | 52 | 38 | 2 | 2 | 0 | 0 | 10 |
| TS 33.310 | 121 | 42 | 3 | 66 | 1 | 0 | 9 |
| TS 29.510 | 190 | 37 | 3 | 14 | 0 | 13 | 123 |
| TS 29.500 | 165 | 60 | 9 | 2 | 0 | 0 | 94 |
| TS 29.501 | 30 | 30 | 0 | 0 | 0 | 0 | 0 |
| TS 29.122 | 40 | 13 | 1 | 4 | 0 | 0 | 22 |
| TS 23.003 | 19 | 3 | 1 | 0 | 3 | 2 | 10 |
| TS 23.502 | 15 | 0 | 0 | 0 | 0 | 3 | 12 |
| **Total** | **1391** | **624** | **44** | **205** | **20** | **90** | **408** |

Onde está o que não é **feito**:

- **TS 33.535 e TS 29.535.** O *roaming* e o aviso de serviço desligado (DV-2, mantido até haver uma VPLMN pelo Q5), o GPSI (MVP-10), o 204 (DV-7) e quatro pontos novos: o Routing Indicator padrão (achado 1), o `anonInd` no NEF (achado 2), a escolha do AAnF pelo AF de dentro (achado 11) e os cabeçalhos de TS 29.500 (achado 14).
- **TS 29.503.** Os 34 parágrafos **não feitos** são o `id-translation-result` de SUPI para GPSI (F10, MVP-10) e o `ROAMING_STATUS` (F30, DV-2).
- **TS 29.222, TS 23.222 e TS 33.122.** Os eventos, a revogação, os *logs* e a auditoria (MVP-17, Q7) e, no DV-14, a atualização da publicação e da entrada, a entrada com aprovação e CAPIF-5.
- **TS 33.310.** Dos 66, 61 são a inscrição automática de certificados (§10, CMPv2), fora do escopo desde o MVP-11; os outros são o OCSP (MVP-18, Q9).
- **TS 29.510.** O *heart-beat* (F29) e a operação de chave dos *tokens* (F28), da primeira passada, e o reuso do resultado da descoberta (achado 13).
- **TS 29.500.** Os cabeçalhos `User-Agent`, `Server` e `3gpp-Sbi-*` (achado 14).
- **TS 33.222 e TS 33.310 §6.1.1.** O certificado do AF do K.2 (achado 12).

## 3. As normas do AKMA

### 3.1 TS 33.535 V19.0.0 (todas as cláusulas)

Um parágrafo por linha. As cláusulas sem *shall* nem *should* (os serviços de §7, as listas de §4.5) estão na tabela da primeira passada, no plano.

| Ref. | Tipo | Requisito | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §4.2.2 ¶1 | S | *"AF shall be authenticated and authorized by the operator network before providing the KAF to the AF."* | feito | MVP-11: o AAnF e o NEF reconhecem o AF pelo certificado ([S15], [CP04], [S16.token]) |
| §4.3.0 ¶1 | S | *"The AAnF interacts with the AUSF and the AF using Service-based Interfaces. When the AF is located in the operator's network, the AAnF shall use …"* | feito | AF de dentro direto no AAnF ([R2.ua]); AF de fora pelo NEF ([R2c.ua-capif]) |
| §4.4.0 ¶1 | S | *"AKMA shall reuse the same UE subscription and the same credentials used for 5G access."* | feito | a assinatura e a credencial do registro do UE ([R2]) |
| §4.4.0 ¶2 | S | *"AKMA shall reuse the 5G primary authentication procedure and methods specified in TS 33.501 [2] for the sake of implicit authentication for AKMA …"* | feito | 5G-AKA ([R2]) e EAP-AKA' (`TestAkmaEAPAKAPrime`), os métodos que o free5GC tem |
| §4.4.0 ¶3 | S | *"The SBA interface between the AAnF and the AUSF shall be confidentiality, integrity and replay protected."* | feito | TLS mútuo na SBI desde o MVP-11 ([S01] a [S07]) |
| §4.4.0 ¶4 | S | *"The SBA interface between AAnF and AF/NEF shall be confidentiality, integrity and replay protected."* | feito | TLS mútuo na SBI desde o MVP-11 ([S01] a [S07]) |
| §4.4.0 ¶5 | S | *"The SBA interface between AAnF and UDM shall be confidentiality, integrity and replay protected."* | feito | TLS mútuo na SBI desde o MVP-11 ([S01] a [S07]) |
| §4.4.0 ¶6 | S | *"The AKMA Application Key (KAF) shall be provided with a maximum lifetime based on the operator’s local authentication policy."* | feito | `akma.kafLifetime` |
| §4.4.1 ¶1 | S | *"Ua\* protocol shall be able to carry AKMA Key Identifier (A-KID) ."* | feito | os perfis K.2, PSK TLS, DTLS e OSCORE levam o A-KID e tratam a expiração ([R2t], [R2d], [R2o], [T11] a [T13]) |
| §4.4.1 ¶2 | S | *"The UE and the AKMA AF shall be able to secure the reference point Ua\* using the AKMA Application Key derived from the AKMA Anchor Key."* | feito | os perfis K.2, PSK TLS, DTLS e OSCORE levam o A-KID e tratam a expiração ([R2t], [R2d], [R2o], [T11] a [T13]) |
| §4.4.1 ¶3 | S | *"The Ua\* protocol shall be able to handle the expiration of KAF."* | feito | os perfis K.2, PSK TLS, DTLS e OSCORE levam o A-KID e tratam a expiração ([R2t], [R2d], [R2o], [T11] a [T13]) |
| §4.4.2 ¶1 | S | *"A-KID shall be globally unique."* | feito | RID, A-TID e domínio de origem; a colisão do A-TID a norma não trata (§6.1 NOTA 2) |
| §4.4.2 ¶2 | S | *"A-KID shall be usable as a key identifier in protocols used in the reference point Ua\*."* | feito | *username* do Digest e identidade PSK ([R2t]) |
| §4.4.2 ¶3 | S | *"AKMA AF shall be able to identify the AAnF serving the UE from the A-KID."* | **em parte** | achado 11: o AF de dentro usa o AAnF da configuração, sem olhar o RID do A-KID |
| §4.4.3 ¶1 | S | *"Applications on the UE shall not be able to get access to KAKMA."* | feito | serviço de chaves do ME ([R2k.given], [R2k.af-id], [R2k.uid], [R2k.duplicate]) |
| §4.4.3 ¶2 | S | *"An application on the UE shall only get the KAF keys related to specific AF Identifiers (AF\_IDs) that the application is authorized to get."* | feito | serviço de chaves do ME ([R2k.given], [R2k.af-id], [R2k.uid], [R2k.duplicate]) |
| §4.4.3 ¶3 | S | *"An application on the UE shall not be able to get access to the KAF keys that belong to other applications."* | feito | serviço de chaves do ME ([R2k.given], [R2k.af-id], [R2k.uid], [R2k.duplicate]) |
| §4.6.1 ¶1 | S | *"The roaming subscriber shall be able to utilize the AKMA feature provided by the home network."* | **desvio** | DV-2, mantido até haver uma VPLMN (Q5, opção c) |
| §4.6.1 ¶2 | S | *"The home network shall be able to control whether its subscriber is authorized to use the service in the visited network."* | **desvio** | DV-2, mantido até haver uma VPLMN (Q5, opção c) |
| §4.7.2 ¶1 | S | *"Confidentiality and integrity protection can be provided for the reference point between the AP and the AS using NDS/IP mechanisms as specified in TS …"* | não se aplica | sem *Authentication Proxy*, opcional (§4.7) |
| §5.2 ¶1 | S | *"AKMA Application Keys KAF shall use explicit lifetimes based on the operator's policy. The lifetime of KAF shall be sent by the AAnF as described in …"* | feito | `expiry` na resposta ([N01b]); a K_AF anterior vale até expirar ([R3.ua-old]) |
| §6.1 ¶1 | S | *"In the response, the UDM may also indicate to the AUSF whether the AKMA Anchor key needs to be generated for the UE. If the AKMA indication is …"* | feito | `udm/0001` ([R2]; `TestSetAkmaIndication`) |
| §6.1 ¶2 | S | *"If the AUSF receives the AKMA indication from the UDM, the AUSF shall store the KAUSF and generate the AKMA Anchor Key (KAKMA) and the A-KID from …"* | feito | `ausf/0001` ([R2.akid-ausf]); a K_AUSF fica no contexto do AUSF |
| §6.1 ¶3 | S | *"The UE shall generate the AKMA Anchor Key (KAKMA) and the A-KID from the KAUSF before initiating communication with an AKMA Application Function."* | feito | UERANSIM `0001` ([R2.akid-aanf]) |
| §6.1 ¶4 | S | *"After AKMA key material is generated, the AUSF selects the AAnF as defined in clause 6.7, and shall send the generated A-KID and KAKMA to the AAnF …"* | feito | `ausf/0001` e AAnF ([R2.akid-aanf], [R3.replaced]) |
| §6.1 ¶5 | S | *"A-KID shall be in NAI format as specified in clause 2.2 of IETF RFC 7542 [6], i.e. username@realm. The username part shall include the RID and the …"* | feito | `shared/akma/akid` (D-A; ADR-0002) |
| §6.1 ¶6 | S | *"The A-TID shall be derived from KAUSF as specified in Annex A.3."* | feito | vetores do MVP-0, em quatro implementações |
| §6.1 ¶7 | S | *"The AUSF shall use the RID received from the UDM as described in step 2 to derive A-KID."* | feito | RID do UDM (D10) |
| §6.1 ¶8 | S | *"KAKMA shall be derived from KAUSF as specified in Annex A.2. Since KAKMA and A-TID in A-KID are both derived from KAUSF based on primary …"* | feito | vetores do MVP-0, em quatro implementações |
| §6.2.1 ¶1 | S | *"The UE shall generate the AKMA Anchor Key (KAKMA) and the A-KID from the KAUSF before initiating communication with an AKMA Application Function. …"* | feito | [R2.ua] |
| §6.2.1 ¶2 | S | *"If the AF does not have an active context associated with the A-KID, then the AF selects the AAnF as defined in clause 6.7, and sends a …"* | **em parte** | o pedido e a seleção, sim; a *Service Disable URI* e a subscrição implícita, não (DV-2) |
| §6.2.1 ¶3 | S | *"The AAnF shall check whether the AAnF can provide the service to the AF based on the configured local policy or based on the authorization …"* | feito | política pelo certificado e pelos escopos ([S14], [S15]) |
| §6.2.1 ¶4 | S | *"The AAnF shall verify whether the subscriber is authorized to use AKMA based on the presence of the UE specific KAKMA key identified by the A-KID."* | feito | 403 `K_AKMA_NOT_PRESENT` ([N06], [B1]) |
| §6.2.1 ¶5 | S | *"If KAKMA is present in AAnF, the AAnF shall continue with step 3."* | feito | 403 `K_AKMA_NOT_PRESENT` ([N06], [B1]) |
| §6.2.1 ¶6 | S | *"If KAKMA is not present in the AAnF, the AAnF shall continue with step 8 with an error response."* | feito | 403 `K_AKMA_NOT_PRESENT` ([N06], [B1]) |
| §6.2.1 ¶7 | S | *"Once receiving the request from the AF, if the AAnF determines this specific AF needs GPSI, according to its local policy, the AAnF sends a …"* | **em parte** | nenhum AF pede GPSI na política local; a busca do GPSI não existe (MVP-10) |
| §6.2.1 ¶8 | S | *"The UDM responds with the GPSI of the UE. The AAnF shall store the received GPSI as part of UE’s AKMA context."* | **não feito** | MVP-10 |
| §6.2.1 ¶9 | S | *"Once receiving the request from the AF, the AAnF shall send a Nudm\_EventExposure\_Subscribe request to UDM with SUPI/GPSI to request the …"* | **desvio** | DV-2 |
| §6.2.1 ¶10 | S | *"The UDM shall send the Nudm\_EventExposure\_Subscribe response to the AAnF with the information of roaming status."* | **desvio** | DV-2; o UDM do free5GC não tem o evento (F30) |
| §6.2.1 ¶11 | S | *"Once the AAnF receives the roaming status from the UDM, it checks the local policy and determines whether to provide service to the UE. If yes, the …"* | **em parte** | a derivação e a validade, sim ([N01b]); a política de *roaming*, não (DV-2) |
| §6.2.1 ¶12 | S | *"The key derivation of KAF shall be performed as specified in Annex A.4."* | feito | Anexo A.4 |
| §6.2.1 ¶13 | S | *"If the AAnF determines to provide AKMA service to the UE, the AAnF sends Naanf\_AKMA\_ApplicationKey\_Get response to the AF with SUPI/GPSI, KAF and the …"* | **em parte** | SUPI, K_AF e validade, sim; a falha por *roaming*, não (DV-2); os A-KID enviados ficam guardados (`AKIDSent`) |
| §6.2.1 ¶14 | S | *"The AF sends the Application Session Establishment Response to the UE. If the information in step 8 indicates failure of AKMA key request, the AF …"* | feito | `403 AKMA_KEY_NOT_AVAILABLE` ([R5.ua]) e `503 AANF_UNAVAILABLE` ([R4.ua]) |
| §6.2.2 ¶1 | S | *"in step 8, the AAnF sends Naanf\_AKMA\_ApplicationKey\_AnonUser\_Get response to the AF with KAF and the KAF expiration time. The AAnF shall store the …"* | feito | [N01c], [E03] |
| §6.3 ¶1 | S | *"When the AF is about to request AKMA Application Key for the UE from the AAnF, e.g. when UE initiates application session establishment request as in …"* | **em parte** | A-KID, AF_ID e `anonInd`, sim; a HPLMN pelo A-KID, não: um operador configurado (achado 11) |
| §6.3 ¶2 | S | *"The AAnF shall process the request in the same way as specified in clause 6.2.1 with following changes:"* | **em parte** | como em §6.2.1, salvo o DV-2 |
| §6.3 ¶3 | S | *"If KAKMA is present in AAnF, the AAnF shall continue with step 4 in this clause."* | feito | [E08] |
| §6.3 ¶4 | S | *"If KAKMA is not present in the AAnF, the AAnF shall continue with step 5 in this clause with an error response."* | feito | [E08] |
| §6.3 ¶5 | S | *"Once receiving the request from the AF, AAnF shall request the UE roaming status report from UDM as specified in clause 6.2.1, step 5-6. If the AAnF …"* | **em parte** | K_AF e validade, sim; o estado de *roaming*, não (DV-2); com `anonInd`, o NEF não recebe o SUPI (achado 2) |
| §6.3 ¶6 | S | *"The NEF forwards the response to the AF, the response contains the KAF, the KAF expiration time (KAF exptime) and optionally GPSI (external ID) or …"* | **em parte** | K_AF e validade, sem o SUPI ([R2n.nef-no-supi]); o GPSI, não (MVP-10) |
| §6.4.1 ¶1 | S | *"KAKMA shall be re-keyed by running a successful primary authentication as described in clause 6.1."* | feito | MVP-12 ([R6r]) |
| §6.6.1 ¶1 | S | *"AAnF shall delete AKMA Context (e.g. SUPI, A-KID, KAKMA, GPSI and KAF expiration time) from its local database identified by SUPI."* | feito | [R6.removed] |
| §6.7 ¶1 | S | *"Internal AFs and the NEF performs AAnF instance selection that handles the AKMA request. The AF/NEF shall utilize the NRF to discover the AAnF …"* | feito | o NEF pelo NRF; o AF de dentro pela configuração, que a frase permite |
| §6.7 ¶2 | S | *"The AUSF performs AAnF selection to allocate an AAnF Instance to send the AKMA key material related to the UE. The AUSF shall utilize the NRF to …"* | feito | pelo NRF, com `akma.aanfUri` como alternativa ([R2]) |
| §6.7 ¶3 | S | *"The NF specified in clause 6.6 performs AAnF instance selection that handles the AKMA request. The NF shall utilize the NRF to discover the AAnF …"* | feito | o OAM é o `seed remove`, com o endereço configurado |
| §6.7 ¶4 | R | *"The AAnF selection functionality in NF consumer or in SCP should consider the following factor:"* | **em parte** | o AUSF e o NEF consideram o RID; o AF de dentro, não (achado 11) |
| §6.7 ¶5 | S | *"Internal AFs, the NEF and the AUSF shall select the same AAnF set based on the UE’s Routing Indicator."* | **em parte** | um AAnF só; com o RID padrão, o AUSF e o NEF não acham AAnF (achado 1) |
| §6.7 ¶6 | S | *"In the case of delegated discovery and selection in SCP, the AAnF NF consumer shall send all available factors to the SCP."* | não se aplica | sem SCP |
| §6.8.1 ¶1 | S | *"If AF(s) are determined at step 4, the AAnF shall send notifications to the subscribed AF(s) about AKMA roaming via …"* | **não feito** | §6.8 depende do DV-2 (Q5) |
| §6.8.1 ¶2 | S | *"The AF shall send the response and based on the notification and internal policy, the AF may stop the UE AKMA service, or may stop the encryption."* | **não feito** | §6.8 depende do DV-2 (Q5) |
| §6.8.2 ¶1 | S | *"Once receiving the notification from AAnF, NEF shall send Nnef\_AKMA\_ServiceDisableNotification to AF."* | **não feito** | §6.8 depende do DV-2 (Q5) |
| §A.1.1 ¶1 | S | *"All key derivations for AKMA shall be performed using the key derivation function (KDF) specified in Annex B.2.2 of TS 33.220 [4]."* | feito | vetores do MVP-0 (`shared/akma`, oráculo Python, UERANSIM); no ME, sem NFKC (DV-4) |
| §A.2 ¶1 | S | *"When deriving a KAKMA from KAUSF, the following parameters shall be used to form the input S to the KDF:"* | feito | vetores do MVP-0 (`shared/akma`, oráculo Python, UERANSIM); no ME, sem NFKC (DV-4) |
| §A.2 ¶2 | S | *"The input key KEY shall be the KAUSF."* | feito | vetores do MVP-0 (`shared/akma`, oráculo Python, UERANSIM); no ME, sem NFKC (DV-4) |
| §A.2 ¶3 | S | *"SUPI shall be the same value as parameter P0 in Annex A.7.0 of TS 33.501 [2]."* | **desvio** | só SUPI IMSI (DV-3) |
| §A.3 ¶1 | S | *"When deriving the A-TID from KAUSF, the following parameters shall be used to form the input S to the KDF:"* | feito | vetores do MVP-0 (`shared/akma`, oráculo Python, UERANSIM); no ME, sem NFKC (DV-4) |
| §A.3 ¶2 | S | *"The input key KEY shall be KAUSF."* | feito | vetores do MVP-0 (`shared/akma`, oráculo Python, UERANSIM); no ME, sem NFKC (DV-4) |
| §A.3 ¶3 | S | *"SUPI shall be the same value as parameter P0 in Annex A.7.0 of TS 33.501 [2]."* | **desvio** | só SUPI IMSI (DV-3) |
| §A.4 ¶1 | S | *"When deriving a KAF from KAKMA, the following parameters shall be used to form the input S to the KDF:"* | feito | vetores do MVP-0 (`shared/akma`, oráculo Python, UERANSIM); no ME, sem NFKC (DV-4) |
| §A.4 ¶2 | S | *"The input key KEY shall be KAKMA."* | feito | vetores do MVP-0 (`shared/akma`, oráculo Python, UERANSIM); no ME, sem NFKC (DV-4) |
| §B.1.2.2 ¶1 | S | *"At step 2, if the clients supports AKMA with this protocol then the client shall add the constant string "3gpp-akma" to the "User-Agent" HTTP header …"* | feito | `TestHostAndProductToken` |
| §B.1.2.2 ¶2 | S | *"At step 3, if the AF selects AKMA for deriving the key, then the AF shall include the "3GPP-bootstrapping-akma" within the WWW-Authenticate header …"* | **desvio** | o `realm` leva a constante, na forma do DV-6; sem GBA, a escolha não se aplica |
| §B.1.2.2 ¶3 | S | *"At step 4, on receiving the response from the AF, the client shall verify that the FQDN in the realm attribute corresponds to the FQDN of the AF it …"* | feito | `TestChallengeChecks`, `TestK2RealmOfAnotherAF` |
| §B.1.2.2 ¶4 | S | *"At step 5 given AKMA has been selected for keying, the client shall send a response with an Authorization header field where Digest is inserted using …"* | feito | Digest com o A-KID e a K_AF (DV-5) |
| §B.1.2.2 ¶5 | S | *"At step 6 given AKMA has been selected for keying, the AF shall verify the value of the password attribute using KAF retrieved from AAnF using the …"* | feito | `TestKeyNotAvailable` |
| §B.1.3.2.1 ¶1 | S | *"At step 2, the AF shall include a constant string "3GPP-AKMA" is used as PSK-identity hint to indicate that AKMA based keying is supported."* | feito | [R2t.psk12], [T04] |
| §B.1.3.2.1 ¶2 | S | *"At step 3, the UE may use an AKMA generated key if support was indicated by the AF (even if GBA-based keys were also indicated as supported by the …"* | feito | [R2t.psk12], [T04] |
| §B.1.3.2.1 ¶3 | S | *"At step 4, if the AF receives the "3GPP-AKMA" prefix and the A-KID in the ClientKeyExchange messages it fetches the AF specific shared secret (KAF) …"* | feito | [R2t.psk12], [T04] |
| §B.1.3.2.2 ¶1 | S | *"In step 1, the PSK identities in the ClientHello shall include a prefix indicating the PSK-identity name space (i.e. "3GPP-AKMA") and the A-KID to …"* | feito | [R2t.psk13], [T03]; o separador `;` (D-E) |
| §B.1.3.2.2 ¶2 | S | *"In step 2 if the AF is willing to establish a TLS tunnel using PSK authentication with AKMA keys, then the AF shall indicate the index of the AKMA …"* | feito | [R2t.psk13], [T03]; o separador `;` (D-E) |
| §B.1.3.2.2 ¶3 | S | *"The UE and NAF shall derive the TLS external PSK from KAF."* | feito | [R2t.psk13], [T03]; o separador `;` (D-E) |
| §C.1.1 ¶1 | R | *"UE hosts the DTLS client. The UE should be able to send the AKMA PSK identity to the AF to indicate which key (KAF) the UE intends to use to secure …"* | feito | [R2d.dtls13], [R2d.dtls12], [D01] a [D03] |
| §C.1.2 ¶1 | R | *"DTLS should be supported by the AF for the UE-AF reference point (Ua\*)."* | feito | [R2d.dtls13], [R2d.dtls12], [D01] a [D03] |
| §C.1.2 ¶2 | R | *"The AF should be able to require that a certain key (i.e., KAF) used to secure the Ua reference point based on DTLS."* | feito | [R2d.dtls13], [R2d.dtls12], [D01] a [D03] |
| §C.2.2 ¶1 | R | *"AKMA PSK identity should be delivered via DTLS message."* | feito | [R2d.dtls13], [R2d.dtls12], [D01] a [D03] |
| §D.2.2 ¶1 | S | *"To utilise AKMA as described in this document the UE shall be equipped with an CoAP capable client implementing the particular features of AKMA as …"* | feito | [R2o.oscore], [O01] a [O10] |
| §D.2.3 ¶1 | S | *"To utilise AKMA as described in this document the AF shall support the features of AKMA as specified in this document."* | feito | [R2o.oscore], [O01] a [O10] |
| §D.3.2 ¶1 | S | *"In Step 1, the CoAP Client (UE) shall send a CoAP request to the AF. This is the Application Session Establishment Request in Step 1 in clause 6.2. …"* | feito | [R2o.oscore], [O01] a [O10] |
| §D.3.2 ¶2 | S | *"URI of the AKMA resource on the AF. The URI shall have the format of <AF\_IP\_or\_FQDN>/akma, where AF\_IP\_or\_FQDN indicates the IP address or the FQDN …"* | feito | [R2o.oscore], [O01] a [O10] |
| §D.3.2 ¶3 | S | *"The CoAP Server (AF) shall respond to the CoAP Client (UE) with a CoAP response. This is the Application Session Establishment Response in Step 5 in …"* | feito | [R2o.oscore], [O01] a [O10] |
| §D.3.3 ¶1 | S | *"The OSCORE security context used in AKMA OSCORE Ua\* protocol is similar to the GBA OSCORE security context specified in TS 33.220 [4], clause Y.3.3 …"* | feito | [R2o.oscore], [O01] a [O10] |
| §D.3.3 ¶2 | S | *"where HKDF shall be the HMAC-based Key Derivation Function specified in IETF RFC 5869 [16]."* | feito | [R2o.oscore], [O01] a [O10] |
| §D.3.4 ¶1 | S | *"Moreover even if KAF remains constant upon a new application session establishment (Step 1 in clause D.3.2) or a renegotiation of the OSCORE key …"* | feito | [R2o.oscore], [O01] a [O10] |
| §D.3.5 ¶1 | S | *"IETF CoAP and OSCORE shall use the IETF Concise Binary Object Representation (CBOR) specified in the IETF RFC 8949 [15] for payload encoding for …"* | feito | [R2o.oscore], [O01] a [O10] |
| §D.3.5 ¶2 | S | *"The CoAP media type for CBOR encoding shall be:"* | feito | [R2o.oscore], [O01] a [O10] |
| §D.3.5 ¶3 | S | *"The Request Payload in the Application Session Request shall be formatted as a CBOR Array as follows:"* | feito | [R2o.oscore], [O01] a [O10] |
| §D.3.5 ¶4 | S | *"The Response Payload in the Application Session Response shall be formatted as a CBOR Array as follows:"* | feito | [R2o.oscore], [O01] a [O10] |

### 3.2 TS 29.535 V20.1.0 (§4 e §5; o YAML do Anexo A)

O texto da Release 20 é o mesmo da baseline V19.5.0 nos requisitos; só a redação mudou (plano, MVP-14, tabela de versões).

| Ref. | Tipo | Requisito | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §4.2.1.3.2 ¶1 | S | *"shall be authenticated and authorized by the operator network before receiving the KAF from the AAnF;"* | feito | MVP-11 ([S15], [CP04]) |
| §4.2.2.2.2 ¶1 | S | *"In order to store the AKMA related key material, the NF service consumer shall send an HTTP POST request message to the resource URI …"* | feito | `ausf/0001`; `TestRegisterAnchorKey` |
| §4.2.2.2.2 ¶2 | S | *"The AkmaKeyInfo data structure shall include:"* | feito | `ausf/0001`; `TestRegisterAnchorKey` |
| §4.2.2.2.2 ¶3 | S | *"If the AAnF cannot successfully fulfil the received HTTP POST request due to an internal error or an error in the HTTP POST request, the AAnF shall …"* | feito | `TestRegisterBadRequests` |
| §4.2.2.2.2 ¶4 | S | *"If the AAnF determines the received HTTP POST request needs to be redirected, the AAnF shall send an HTTP redirect response as specified in clause …"* | não se aplica | sem redirecionamento: um AAnF, sem SCP |
| §4.2.2.2.2 ¶5 | S | *"Upon successful reception of an HTTP POST request, the AAnF shall store the key material information and respond to the NF service consumer with a …"* | feito | 200 com `AkmaKeyInfo` |
| §4.2.2.3.2 ¶1 | S | *"The NF service consumer shall invoke the Naanf\_AKMA\_ApplicationKey\_Get service operation to retrieve the AKMA Application Key information. The NF …"* | feito | [N01b] |
| §4.2.2.3.2 ¶2 | S | *"If the request corresponds to a Naanf\_AKMA\_ApplicationKey\_AnonUser\_Get request, then the AkmaAfKeyRequest shall contain the "anonInd" attribute set …"* | **em parte** | o AF, sim ([N01c]); o NEF repassa o `anonInd` (achado 2) |
| §4.2.2.3.2 ¶3 | S | *"If the AAnF determines the received HTTP POST request needs to be redirected, the AAnF shall send an HTTP redirect response as specified in clause …"* | não se aplica | sem redirecionamento: um AAnF, sem SCP |
| §4.2.2.3.2 ¶4 | S | *"If the AAnF cannot successfully fulfil the received HTTP POST request due to an internal error or an error in the HTTP POST request, the AAnF shall …"* | feito | §5.1.7 |
| §4.2.2.3.2 ¶5 | S | *"The AAnF shall also verify the presence of the UE specific KAKMA key identified by the A-KID."* | feito | [N06] |
| §4.2.2.3.2 ¶6 | S | *"If KAKMA is not present in the AAnF, the AAnF shall reply with an HTTP "403 Forbidden" status code and the response message body including a …"* | feito | [N06] |
| §4.2.2.3.2 ¶7 | S | *"If KAKMA is present in the AAnF, the AAnF shall continue and process the request as specified below."* | feito | [N06] |
| §4.2.2.3.2 ¶8 | S | *"Upon the reception of the HTTP POST request, the AAnF shall respond with an HTTP "200 OK" status code and the response message body containing the …"* | feito | `kaf`, `expiry` e `supi`; `[INTERPRETAÇÃO]` sem o escopo `supi-access`, sem `supi` (tab. 5.1.9-1) |
| §4.2.2.3.2 ¶9 | S | *"If the requested AKMA Application Key information for the UE does not exist, the AAnF shall respond with "204 No Content"."* | **desvio** | DV-7 |
| §4.2.2.3.2 ¶10 | S | *"If the NF service consumer is an NEF, and if UE identifier is required to relay to the AF based on local policy, the NEF invokes the …"* | **não feito** | MVP-10 (F10) |
| §4.2.2.3.2 ¶11 | S | *"If the AAnF identifies that the request is targeting a UE that is currently roaming (e.g. the AAnF have received a roaming status information …"* | **desvio** | DV-2 |
| §4.2.2.4.2 ¶1 | S | *"The NF service consumer shall invoke the Naanf\_AKMA\_ContextRemove service operation to request the AAnF to remove the AKMA related key material. The …"* | feito | `TestRemoveContext`; [R6.removed], [R6.not-found] |
| §4.2.2.4.2 ¶2 | S | *"If errors occur when processing the HTTP POST request, the AAnF shall send an HTTP error response as specified in clause 5.1.7."* | feito | `TestRemoveContext`; [R6.removed], [R6.not-found] |
| §4.2.2.4.2 ¶3 | S | *"If the AAnF determines the received HTTP POST request needs to be redirected, the AAnF shall send an HTTP redirect response as specified in clause …"* | não se aplica | sem redirecionamento: um AAnF, sem SCP |
| §4.2.2.4.2 ¶4 | S | *"Upon the reception of the HTTP POST request, if the AKMA context (e.g. A-KID, KAKMA) for the UE has been removed successfully, the AAnF shall send an …"* | feito | `TestRemoveContext`; [R6.removed], [R6.not-found] |
| §4.2.2.4.2 ¶5 | S | *"If the AKMA Context resource does not exist, the AAnF shall respond with "404 Not Found" and the "cause" attribute of the "ProblemDetails" data …"* | feito | `TestRemoveContext`; [R6.removed], [R6.not-found] |
| §4.2.2.5.2 ¶1 | S | *"In order to notify a previously subscribed service consumer on AKMA service disablement, the AAnF shall invoke the Naanf\_AKMA\_Notify service …"* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §4.2.2.5.2 ¶2 | S | *"Upon reception of the notification, the NF service consumer shall acknowledge the successful reception of the notification with a "204 No Content" …"* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §4.2.2.5.2 ¶3 | S | *"If errors occur when processing the HTTP POST request, the NF service consumer shall send an appropriate HTTP error response as specified in clause …"* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §5.1.1 ¶1 | S | *"The Naanf\_AKMA service shall use the Naanf\_AKMA API."* | feito | `{apiRoot}/naanf-akma/v1` |
| §5.1.1 ¶2 | S | *"The API URI of the Naanf\_AKMA API shall be:"* | feito | `{apiRoot}/naanf-akma/v1` |
| §5.1.1 ¶3 | S | *"The request URIs used in HTTP requests from the NF service consumer towards the NF service producer shall have the Resource URI structure defined in …"* | feito | `{apiRoot}/naanf-akma/v1` |
| §5.1.1 ¶4 | S | *"The {apiRoot} shall be set as described in 3GPP TS 29.501 [5]."* | feito | `{apiRoot}/naanf-akma/v1` |
| §5.1.1 ¶5 | S | *"The <apiName> shall be "naanf-akma"."* | feito | `{apiRoot}/naanf-akma/v1` |
| §5.1.1 ¶6 | S | *"The <apiVersion> shall be "v1"."* | feito | `{apiRoot}/naanf-akma/v1` |
| §5.1.1 ¶7 | S | *"The <apiSpecificResourceUriPart> shall be set as described in clause 5.1.3 and 5.1.4."* | feito | `{apiRoot}/naanf-akma/v1` |
| §5.1.2.1 ¶1 | S | *"HTTP/2, IETF RFC 9113 [11], shall be used as specified in clause 5 of 3GPP TS 29.500 [4]."* | feito | HTTP/2 sobre TLS; `[INTERPRETAÇÃO]` o servidor também aceita HTTP/1.1, que os testes negativos usam |
| §5.1.2.1 ¶2 | S | *"HTTP/2 shall be transported as specified in clause 5.3 of 3GPP TS 29.500 [4]."* | feito | HTTP/2 sobre TLS; `[INTERPRETAÇÃO]` o servidor também aceita HTTP/1.1, que os testes negativos usam |
| §5.1.2.2.2 ¶1 | S | *"JSON, IETF RFC 8259 [12], shall be used as content type of the HTTP bodies specified in the present specification as specified in clause 5.4 of 3GPP …"* | feito | `application/json` e `application/problem+json` |
| §5.1.2.2.2 ¶2 | S | *""Problem Details" JSON object shall be used to indicate additional details of the error in a HTTP response body and shall be signalled by the content …"* | feito | `application/json` e `application/problem+json` |
| §5.1.2.3 ¶1 | S | *"The mandatory HTTP custom header fields specified in clause 5.2.3.2 of 3GPP TS 29.500 [4] shall be supported, and the optional HTTP custom header …"* | **em parte** | achado 14: o AAnF não trata os cabeçalhos de TS 29.500 tab. 5.2.3.2.1-1 |
| §5.1.4.2.2 ¶1 | S | *"This operation shall support the response data structures and response codes specified in tables 5.1.4.2.2-1 and 5.1.4.2.2-2."* | feito | `api_akma.go` |
| §5.1.4.3.2 ¶1 | S | *"This operation shall support the response data structures and response codes specified in tables 5.1.4.3.2-1 and 5.1.4.3.2-2."* | feito | `api_akma.go` |
| §5.1.4.3.2 ¶2 | S | *"n/a / / / 204 No Content / If the requested data does not exist, the AAnF shall respond with "204 No Content"."* | **desvio** | DV-7 |
| §5.1.4.3.2 ¶3 | S | *"NOTE 1: The mandatory HTTP error status codes for the POST method listed in Table 5.2.7.1-1 of 3GPP TS 29.500 [4] also apply. / NOTE 2: For the …"* | feito | os códigos de TS 29.500 tab. 5.2.7.1-1 (`aanf.md` §6); `supi`, salvo no acesso anônimo |
| §5.1.4.3.2 ¶4 | R | *"Location / string / M / 1 / Contains an alternative URI representing the endpoint of an alternative AAnF (service) instance towards which the request …"* | não se aplica | redirecionamento |
| §5.1.4.4.2 ¶1 | S | *"This operation shall support the response data structures and response codes specified in tables 5.1.4.4.2-1 and 5.1.4.4.2-2."* | feito | `api_akma.go` |
| §5.1.4.4.2 ¶2 | S | *"n/a / / / 204 No Content / Successful case: The AKMA context matching the "CtxRemove" in the request body was deleted, the AAnF shall respond with …"* | feito | `api_akma.go` |
| §5.1.4.4.2 ¶3 | S | *"RedirectResponse / O / 0..1 / 307 Temporary Redirect / Temporary redirection, during remove procedure. The response shall include a Location header …"* | não se aplica | redirecionamento |
| §5.1.4.4.2 ¶4 | S | *"RedirectResponse / O / 0..1 / 308 Permanent Redirect / Permanent redirection, during remove procedure. The response shall include a Location header …"* | não se aplica | redirecionamento |
| §5.1.5.1 ¶1 | S | *"Notifications shall comply to clause 6.2 of 3GPP TS 29.500 [4] and clause 4.6.2.3 of 3GPP TS 29.501 [5]."* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §5.1.5.2.2 ¶1 | S | *"The Callback URI "{notifUri}" shall be used with the callback URI variables defined in table 5.1.5.2.2-1."* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §5.1.5.2.3.1 ¶1 | S | *"This method shall support the request data structures specified in table 5.1.5.2.3.1-1 and the response data structures and response codes specified …"* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §5.1.5.2.3.1 ¶2 | S | *"NOTE 1: The mandatory HTTP error status codes for the HTTP POST method listed in table 5.2.7.1-1 of 3GPP TS 29.500 [4] shall also apply. / NOTE 2: …"* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §5.1.6.2.2 ¶1 | S | *"suppFeat / SupportedFeatures / C / 0..1 / Contains the list of supported feature(s) among the ones defined in clause 5.1.8. / This attribute shall be …"* | feito | `negotiatedFeatures` |
| §5.1.6.2.2 ¶2 | S | *"NOTE: When the "AKMA\_GPSI\_Support" feature is supported, the presence of "supi" attribute and the "gpsi" attribute are mutually exclusive. When the …"* | feito | só `supi`; `gpsi` recebe 400 (`aanf.md` §3.1) |
| §5.1.6.2.3 ¶1 | S | *"NOTE: In current release of specification, only "supi" can be used to indicate the AKMA context to be removed. The "supi" attribute shall be included."* | feito | sem `supi`, 400 (D-I) |
| §5.1.7.1 ¶1 | S | *"For the Naanf\_AKMA API, HTTP error responses shall be supported as specified in clause 4.8 of 3GPP TS 29.501 [5]. Protocol errors and application …"* | feito | `aanf.md` §6; `TestRoutingErrors` |
| §5.1.8 ¶1 | S | *"The optional features in table 5.1.8-1 are defined for the Naanf\_AKMA API. They shall be negotiated using the extensibility mechanism defined in …"* | feito | nenhuma *feature* suportada (`TestRegisterFeatureNegotiation`) |
| §5.1.9 ¶1 | S | *"If OAuth2 is used, an NF Service Consumer, prior to consuming services offered by the Naanf\_AKMA API, shall obtain a "token" from the authorization …"* | feito | AUSF, NEF e AF pedem o *token* ao NRF ([S14]) |
| §A.2 ¶1 | S | *"Parameters to request to delete the AKMA context for the UE, the "supi" attribute shall be"* | feito | `TestRemoveContext` |

### 3.3 TS 29.522 V20.1.0 (§4.4.23, §5.2, §5.14, §6 e §7)

| Ref. | Tipo | Requisito | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §4.4.23.2 ¶1 | S | *"In order to retrieve the AKMA application key, the AF shall send an HTTP POST request message to the resource URI "{apiRoot}/3gpp-akma/v1/retrieve". …"* | feito | `af/af/internal/nef` ([E01]); `notifUri` só com `RoamingRestriction` |
| §4.4.23.2 ¶2 | S | *"Upon receipt of the corresponding HTTP POST message from the AF, if the AF's request is authorized by the NEF, then the NEF shall interact with the …"* | feito | `nef/0001`; sem SUPI ([R2n.nef-no-supi]); o GPSI, opcional, não vai |
| §4.4.23.2 ¶3 | S | *"If the NEF receives a response from the AAnF with an HTTP "403 Forbidden" status code and the response message body including a ProblemDetails data …"* | **em parte** | repassa `K_AKMA_NOT_PRESENT` ([E08]); `ROAMING_AKMA_SERVICE_DENIED` nunca chega (DV-2) |
| §4.4.23.3 ¶1 | S | *"In order to notify an AF of AKMA service disablement, the NEF shall send an HTTP POST request message to the AF with the request body including the …"* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §4.4.23.3 ¶2 | S | *"Upon reception of this notification request, the AF shall acknowledge its successful reception by sending a response message with an HTTP "204 No …"* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §4.4.23.3 ¶3 | S | *"On failure, the AF shall take proper error handling actions, as specified in clause 5.14.7, and respond to the NEF with an appropriate error status …"* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §5.2 ¶1 | S | *"The usage of HTTP, content type and URI structure definition, as specified in clauses 5.2.2, 5.2.3 and 5.2.4 of 3GPP TS 29.122 [4] respectively, …"* | **em parte** | HTTP/1.1 e HTTP/2 sobre TLS; o HTTP/2 por ALPN, sem o *upgrade* (achado 4, DV-15) |
| §5.2 ¶2 | S | *"The notification, error handling, feature negotiation, HTTP custom headers as specified in clauses 5.2.5, 5.2.6, 5.2.7, 5.2.8 of 3GPP TS 29.122 [4] …"* | **em parte** | erros, *features* e cabeçalhos, sim; as notificações, não (DV-2) |
| §5.2 ¶3 | S | *"The conventions for Open API specification files as specified in clause 5.2.9 of 3GPP TS 29.122 [4] shall be applicable for NEF Northbound APIs."* | não se aplica | convenção para quem escreve o OpenAPI |
| §5.14.1 ¶1 | S | *"The Nnef\_AKMA service shall use the AKMA API."* | feito | `nef/0001` |
| §5.14.1 ¶2 | S | *"The API URI of AKMA API shall be:"* | feito | `nef/0001` |
| §5.14.1 ¶3 | S | *""apiName" shall be set to "3gpp-akma"."* | feito | `nef/0001` |
| §5.14.1 ¶4 | S | *""apiVersion" shall be set to "v1" for the current version defined in the present document."* | feito | `nef/0001` |
| §5.14.3.2.2 ¶1 | S | *"This operation shall support the request and response data structures and response codes specified in tables 5.14.3.2.2-1 and 5.14.3.2.2-2."* | feito | `nef/0001` |
| §5.14.3.2.2 ¶2 | S | *"n/a / / / 204 No Content / If the requested data does not exist, the NEF shall respond with "204 No Content"."* | feito | `[INTERPRETAÇÃO]` o NEF repassa o 204 do AAnF, que nunca o manda (DV-7) |
| §5.14.3.2.2 ¶3 | S | *"n/a / / / 307 Temporary Redirect / Temporary redirection. The response shall include a Location header field containing an alternative URI of the …"* | não se aplica | redirecionamento |
| §5.14.3.2.2 ¶4 | S | *"n/a / / / 308 Permanent Redirect / Permanent redirection. The response shall include a Location header field containing an alternative URI of the …"* | não se aplica | redirecionamento |
| §5.14.4.1 ¶1 | S | *"Notifications shall comply to clause 5.2.5 of 3GPP TS 29.122 [4]."* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §5.14.4.2.2 ¶1 | S | *"The Callback URI "{notifUri}" shall be used with the callback URI variables defined in table 5.14.4.2.2-1."* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §5.14.4.2.3 ¶1 | S | *"This method shall support the request data structures specified in table 5.14.4.2.3-1 and the response data structures and response codes specified …"* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §5.14.4.2.3 ¶2 | S, R | *"n/a / / / 307 Temporary Redirect / Temporary redirection. / The response shall include a Location header field containing an alternative URI …"* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §5.14.4.2.3 ¶3 | S, R | *"n/a / / / 308 Permanent Redirect / Permanent redirection. / The response shall include a Location header field containing an alternative URI …"* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §5.14.4.2.3 ¶4 | S | *"NOTE: The mandatory HTTP error status codes for the HTTP POST method listed in table 5.2.6-1 of 3GPP TS 29.122 [4] shall also apply."* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §5.14.4.2.3 ¶5 | R | *"Location / string / M / 1 / Contains an alternative URI representing the end point of an alternative AF towards which the notification should be …"* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §5.14.4.2.3 ¶6 | R | *"Location / string / M / 1 / Contains an alternative URI representing the end point of an alternative AF towards which the notification should be …"* | **não feito** | §6.8 de TS 33.535 (DV-2, Q5) |
| §5.14.5.3.2 ¶1 | S | *"suppFeat / SupportedFeatures / C / 0..1 / Contains the list of supported feature(s) among the ones defined in clause 5.14.6. / This attribute shall …"* | feito | `nef/0001` |
| §5.14.5.3.3 ¶1 | S | *"suppFeat / SupportedFeatures / O / 0..1 / Contains the list of supported feature(s) among the ones defined in clause 5.14.6. / This attribute shall …"* | feito | `nef/0001` |
| §5.14.5.3.3 ¶2 | S | *"NOTE 1: Properties marked with a feature as defined in clause 5.14.6 are applicable as described in clause 5.2.7 of 3GPP TS 29.122 [4]. If no feature …"* | feito | `nef/0001` |
| §5.14.5.4.2 ¶1 | S | *"The simple data types defined in table 5.14.5.4.2-1 shall be supported."* | feito | `shared/akma/afid` e `shared/akma/akid` (D-A) |
| §5.14.5.4.2 ¶2 | S | *"AKId / string / AKMA Key Identifier shall be in NAI format as specified in clause 2.2 of IETF RFC 7542 [40], which is formatted as the following …"* | feito | `shared/akma/afid` e `shared/akma/akid` (D-A) |
| §5.14.7.1 ¶1 | S | *"HTTP error handling shall be supported as specified in clause 5.2.6 of 3GPP TS 29.122 [4]."* | feito | TS 29.122 §5.2.6 (F12, MVP-11) |
| §5.14.7.1 ¶2 | S | *"In addition, the requirements in the following clauses shall apply."* | feito | TS 29.122 §5.2.6 (F12, MVP-11) |
| §6 ¶1 | S | *"The security aspects of NEF shall be supported as defined in clause 12 of 3GPP TS 33.501 [6]."* | feito | TS 33.501 §12 (MVP-11) |
| §6 ¶2 | S | *"When CAPIF is not supported and/or not used, the access to the NEF Northbound APIs shall be authorized by means of OAuth2 protocol (see IETF RFC 6749 …"* | não se aplica | o laboratório usa o CAPIF; o modo sem CAPIF fica desligado, sem servidor de autorização (Q3, opção a) |
| §6 ¶3 | S | *"When the NEF supports CAPIF for external exposure and CAPIF is used, the provisions of clause 7.2 shall apply for authorizing the access to the NEF …"* | feito | §7.2 |
| §7.1 ¶1 | S | *"When CAPIF is used with an NEF that is used for external exposure, the NEF shall support the following as defined in 3GPP TS 29.222 [12]:"* | **em parte** | AEF, APF e o método 3, sim; CAPIF-5, não (DV-14; `[INTERPRETAÇÃO]` implantação centralizada) |
| §7.2 ¶1 | S | *"When CAPIF is used for managing the exposure of the NEF northbound APIs, before invoking an API exposed by the NEF, the AF, acting as an API Invoker, …"* | feito | [CP03], [CP04] |
| §7.2 ¶2 | S | *"If PKI or TLS-PSK is used as the selected security method between the AF and the NEF, upon API invocation, the NEF shall retrieve the authorization …"* | não se aplica | só OAuth2 |
| §7.2 ¶3 | S | *"If OAuth2 is used as the selected security method between the AF and the NEF, the AF, prior to consuming services offered by the NEF northbound APIs, …"* | feito | [CP05], [CP06], [CP14], [CP15] |
| §7.2 ¶4 | S | *"It is the NEF responsibility to check whether the AF is authorized to use an API based on the "token". Once the NEF verifies the "token", it shall …"* | feito | [CP05], [CP06], [CP14], [CP15] |

### 3.4 TS 24.109 V19.0.0 (Anexo K e as partes de §5.2 e §5.3 que ele usa)

| Ref. | Tipo | Requisito | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §5.2.1 ¶1 | S | *"The UE shall indicate to an application server (i.e. a NAF) that it supports 3GPP-bootstrapping based HTTP Digest authentication by including a …"* | não se aplica | GBA; o K.2 b) 1) troca pelo *product token* `3gpp-akma` |
| §5.2.1 ¶2 | S | *"The UE may indicate multiple GBA modes by inserting multiple "product" tokens in the User-Agent header field. The User-Agent header field with GBA …"* | não se aplica | GBA; o K.2 b) 1) troca pelo *product token* `3gpp-akma` |
| §5.2.1 ¶3 | S | *"Upon receiving GBA related "product" tokens, the application server if it supports NAF functionality may decide to authenticate the UE using …"* | não se aplica | GBA; o K.2 b) 1) troca pelo *product token* `3gpp-akma` |
| §5.2.2.1 ¶1 | S | *"HTTP Digest authentication as specified in RFC 7616 [36] shall be used with previously bootstrapped security association as follows:"* | feito | RFC 7616 (DV-5) |
| §5.2.2.1 ¶2 | S | *"the "username" parameter shall be the bootstrapping transaction identifier (B-TID);"* | não se aplica | GBA: B-TID, GBA_U, USS |
| §5.2.2.1 ¶3 | S | *"the password used in the digest calculations shall be"* | não se aplica | GBA: B-TID, GBA_U, USS |
| §5.2.2.1 ¶4 | S | *"3 the "realm" parameter shall contain two parts delimited by "@" sign. The first part is the constant string "3GPP-bootstrapping" (in case of a …"* | não se aplica | GBA: B-TID, GBA_U, USS |
| §5.2.2.1 ¶5 | S | *"In the case of GBA\_U, the NAF shall indicate to the UE which NAF specific key was selected to be used by setting the first part of the realm to …"* | não se aplica | GBA: B-TID, GBA_U, USS |
| §5.2.2.1 ¶6 | S | *"Both the UE and the NAF shall verify upon receiving each of the HTTP responses and HTTP requests that the second part of the realm attribute is equal …"* | feito | como o K.2 d) |
| §5.2.2.1 ¶7 | S | *"In the case of GBA\_U, if the HTTPS client application resides in the ME, then the application shall use only the ME-based key i.e. Ks\_ext\_NAF (the …"* | não se aplica | GBA: B-TID, GBA_U, USS |
| §5.2.2.1 ¶8 | S | *"In the case of GBA\_U, the operator may indicate the type of the key to be used in the Ua reference point in the NAF specific USS as specified in 3GPP …"* | não se aplica | GBA: B-TID, GBA_U, USS |
| §5.2.4 ¶1 | S | *"NAF shall indicate to the UE that bootstrapped security association is required by sending an HTTP response with code 401 "Unauthorized" and include …"* | feito | 401 com `WWW-Authenticate` (`TestInitiation`) |
| §5.2.5 ¶1 | S | *"The NAF shall indicate to the UE that the existing bootstrapped security association used in the last HTTP request sent by the UE has expired and …"* | feito | K_AF vencida: `refetch` ou `reject` ([N22], [N23]) |
| §K.2 ¶1 | S | *"the UE and the AF shall support the TLS version as specified in annex E of 3GPP TS 33.310 [25]. See clause 5.3.1 in 3GPP TS 33.222 [5] for the …"* | feito | [R2t.k2-tls13], [R2t.k2-tls12], [T01], [T02], [T09] |
| §K.2 ¶2 | S | *"when the UE starts communication via Ua\* reference point with the AF, it shall establish a TLS tunnel with the AF. The AF is authenticated to the UE …"* | feito | [R2t.k2-tls13], [R2t.k2-tls12], [T01], [T02], [T09] |
| §K.2 ¶3 | S | *"the UE shall indicate to an AF that it supports AKMA based HTTP Digest authentication by including a "product" token, that is a constant string …"* | feito | [R2t.k2-tls13], [R2t.k2-tls12], [T01], [T02], [T09] |
| §K.2 ¶4 | S | *"the AF may decide to authenticate the UE using AKMA-based shared secret by executing the authentication procedure. This is indicated in the "realm" …"* | **desvio** | contém a constante; a forma do `realm` é o DV-6 |
| §K.2 ¶5 | S | *"the UE shall generate the HTTP request and the AF shall authenticate the HTTP request using HTTP Digest. HTTP Digest authentication (see RFC 3310 …"* | **desvio** | DV-5 |
| §K.2 ¶6 | S | *"the "username" parameter shall be the A-KID;"* | feito | [R2t.k2-tls13], [R2t.k2-tls12], [T01], [T02], [T09] |
| §K.2 ¶7 | S | *"the password used in the digest calculations shall be KAF (AKMA Application Key) with the KAF Base64 encoded as specified in RFC 4648 [37]; and"* | feito | [R2t.k2-tls13], [R2t.k2-tls12], [T01], [T02], [T09] |
| §K.2 ¶8 | S | *"the "realm" parameter shall contain two parts delimited by "@" sign where the first part is the constant string "3gpp-akma" and the latter part shall …"* | **desvio** | DV-6 |
| §K.2 ¶9 | S | *"both the UE and the AF shall verify upon receiving each of the HTTP responses and HTTP requests that the second part of the realm attribute is equal …"* | feito | [R2t.k2-tls13], [R2t.k2-tls12], [T01], [T02], [T09] |
| §K.3.2 ¶1 | S | *"The PSK TLS handshake shall be used with bootstrapped security association as follows:"* | feito | [R2t.psk12], [R2t.psk13], [T03] a [T06] |
| §K.3.2 ¶2 | S | *"the ClientHello message shall contain one or more PSK-based ciphersuites;"* | feito | [R2t.psk12], [R2t.psk13], [T03] a [T06] |
| §K.3.2 ¶3 | S | *"the ClientHello message shall contain the server\_name TLS extension and it shall contain the hostname of the AF;"* | feito | [R2t.psk12], [R2t.psk13], [T03] a [T06] |
| §K.3.2 ¶4 | S | *"the ServerHello message shall contain a PSK-based ciphersuite selected by the AF;"* | feito | [R2t.psk12], [R2t.psk13], [T03] a [T06] |
| §K.3.2 ¶5 | S | *"the ServerKeyExchange shall be sent by the server and it shall contain the psk\_identity\_hint field and it shall contain the static string "3GPP-AKMA";"* | feito | [R2t.psk12], [R2t.psk13], [T03] a [T06] |
| §K.3.2 ¶6 | S | *"the ClientKeyExchange shall contain the psk\_identity field and it shall contain a prefix "3GPP-AKMA" and the A-KID. If the UE has a choice between …"* | feito | prefixo e A-KID, com `;` (D-S); sem GBA, a escolha não se aplica |
| §K.3.2 ¶7 | S | *"the UE and AF shall derive the TLS premaster secret from KAF (AKMA Application Key)."* | feito | [R2t.psk12], [R2t.psk13], [T03] a [T06] |
| §K.3.3 ¶1 | S | *"The PSK TLS handshake shall be used with bootstrapped security association as follows:"* | feito | [R2t.psk12], [R2t.psk13], [T03] a [T06] |
| §K.3.3 ¶2 | S | *"The UE shall include in the ClientHello message:"* | feito | [R2t.psk12], [R2t.psk13], [T03] a [T06] |
| §K.3.3 ¶3 | S | *"PSK identities within the psk\_identities field. The psk\_identity parameters within the psk\_identities field shall contain a prefix indicating the PSK …"* | feito | [R2t.psk12], [R2t.psk13], [T03] a [T06] |
| §K.3.3 ¶4 | S | *"The UE shall derive the TLS external PSK from the AF specific key KAF."* | feito | [R2t.psk12], [R2t.psk13], [T03] a [T06] |
| §K.3.3 ¶5 | S | *"If the AF is willing to establish a TLS tunnel using PSK authentication with AKMA keys, the AF shall reply with the ServerHello message and indicate …"* | feito | [R2t.psk12], [R2t.psk13], [T03] a [T06] |
| §K.3.3 ¶6 | S | *"The AF shall derive the TLS external PSK from the AF specific key KAF."* | feito | [R2t.psk12], [R2t.psk13], [T03] a [T06] |

## 4. As normas que TS 33.535 chama

### 4.1 TS 33.501 V20.3.0 (§6.1, §12, §13.1.0, §13.3, §13.4.1.0 a §13.4.1.1, §14.1.2 a §14.3, A.7.0)

A Release 20 é igual à baseline V19.7.0 nessas cláusulas, salvo a NOTA 1 de §13.1.0 e uma frase do CCA.

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §6.1.1.1 | General | ¶1 (1 S) | feito | o free5GC tem os dois métodos; o UERANSIM só o 5G-AKA, e o banco de testes, o EAP-AKA' |
| §6.1.1.1 | General | ¶2 a ¶3 (2 S) | não se aplica | USIM e UICC simulados pelo UERANSIM |
| §6.1.1.1 | General | ¶4 (1 S) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.1.3 | Granularity of anchor key binding to serving net | ¶1 a ¶4 (4 S) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.1.4.2 | Construction of the serving network name by the  | ¶1 a ¶3 (3 S) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.1.4.3 | Construction of the serving network name by the  | ¶1 a ¶4 (4 S) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.2 | Initiation of authentication and selection of au | ¶1 a ¶7 (7 S) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.2 | Initiation of authentication and selection of au | ¶8 (1 S) | não se aplica | *Disaster Roaming* |
| §6.1.2 | Initiation of authentication and selection of au | ¶9 a ¶10 (2 S) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.2 | Initiation of authentication and selection of au | ¶11 (1 S) | não se aplica | *Disaster Roaming* |
| §6.1.3.1 | Authentication procedure for EAP-AKA' | ¶1 a ¶2 (2 S) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.3.1 | Authentication procedure for EAP-AKA' | ¶3 (1 S) | feito | `udm/0001` ([R2]) |
| §6.1.3.1 | Authentication procedure for EAP-AKA' | ¶4 a ¶14 (11 S) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.3.2.0 | 5G AKA | ¶1 a ¶2 (2 S) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.3.2.0 | 5G AKA | ¶3 (1 S) | feito | `udm/0001` ([R2]) |
| §6.1.3.2.0 | 5G AKA | ¶4 a ¶15 (12 S) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.3.2.2 | RES\* verification failure in SEAF or AUSF or bot | ¶1 a ¶5 (5 S) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.3.3.1 | Synchronization failure or MAC failure in USIM | ¶1 a ¶2 (2 S) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.4.1 | Introduction | ¶1 (1 S) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.4.1a | Linking authentication confirmation to Nudm\_UECM | ¶1 a ¶4 (4 S) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.4.2 | Guidance on linking authentication confirmation  | ¶1 (0 S, 1 R) | free5GC | autenticação primária do free5GC e do UERANSIM, que o projeto não altera |
| §6.1.5.2 | Security mechanisms | ¶1 a ¶3 (3 S) | feito | MVP-12: `amf/0001`, `udm/0002` ([R6r]); só o acesso 3GPP |
| §12.1 | General | ¶1 (0 S, 1 R) | não se aplica | `Nnef_UEId` |
| §12.2 | Mutual authentication | ¶1 a ¶2 (2 S) | feito | MVP-11 ([S16.token], [S16.check], [CP04]) |
| §12.3 | Protection of the NEF – AF interface | ¶1 (1 S) | feito | MVP-11 ([S16.token], [S16.check], [CP04]) |
| §12.3 | Protection of the NEF – AF interface | ¶2 (1 S) | **em parte** | o perfil de TS 33.210 §6.2, salvo o OCSP *stapling* (achado 3). *"Security profiles for TLS implementation and usage shall follow the provisions given in clause 6.2 of TS …"* |
| §12.4 | Authorization of Application Function’s requests | ¶1 (1 S) | feito | MVP-11 ([S16.token], [S16.check], [CP04]) |
| §12.5 | Support for CAPIF | ¶1 (1 S) | feito | MVP-11 ([S16.token], [S16.check], [CP04]) |
| §13.1.0 | General | ¶1 (1 S) | feito | MVP-11 ([S01] a [S09]) |
| §13.1.0 | General | ¶2 (1 S) | **em parte** | o perfil de TS 33.210 §6.2 e o SNI, salvo o OCSP *stapling* (achado 3). *"The TLS profile shall follow the profile given in clause 6.2 of TS 33.210 [3] with the restriction that it …"* |
| §13.1.0 | General | ¶3 (1 S) | feito | MVP-11 ([S01] a [S09]) |
| §13.3.0 | Static authorization | ¶1 (1 S) | não se aplica | o laboratório usa *tokens* |
| §13.3.1.1 | Direct communication | ¶1 a ¶3 (3 S) | feito | MVP-11 (`nrf/0005`; [S10]) |
| §13.3.1.3 | Authorization of discovery request and error han | ¶1 a ¶3 (3 S) | feito | MVP-11 (`nrf/0005`; [S10]) |
| §13.3.1.3 | Authorization of discovery request and error han | ¶4 (1 S) | não se aplica | *roaming* |
| §13.3.1.3 | Authorization of discovery request and error han | ¶5 (1 S) | feito | MVP-11 (`nrf/0005`; [S10]) |
| §13.3.2.1 | Direct communication | ¶1 a ¶3 (3 S) | feito | MVP-11 (`nrf/0005`; [S10]) |
| §13.3.2.4 | Error handling | ¶1 (1 S) | feito | MVP-11 (`nrf/0005`; [S10]) |
| §13.4.1.0 | General | ¶1 (1 S) | feito | MVP-11 (`nrf/0005`, `openapi/0006`; [S10] a [S14]) |
| §13.4.1.1.1 | OAuth 2.0 roles | ¶1 a ¶6 (6 S) | feito | MVP-11 (`nrf/0005`, `openapi/0006`; [S10] a [S14]) |
| §13.4.1.1.2 | Service Request Process | ¶1 a ¶3 (3 S) | feito | MVP-11 (`nrf/0005`, `openapi/0006`; [S10] a [S14]) |
| §13.4.1.1.2 | Service Request Process | ¶4 (1 S) | **em parte** | sem fatias no AKMA; o escopo é conferido num só perfil do tipo (F27). *"The NRF shall additionally verify the S-NSSAIs of the NF Service Consumer and check whether there are …"* |
| §13.4.1.1.2 | Service Request Process | ¶5 a ¶17 (13 S) | feito | MVP-11 (`nrf/0005`, `openapi/0006`; [S10] a [S14]) |
| §13.4.1.1.2 | Service Request Process | ¶18 a ¶19 (2 S) | não se aplica | sem NF Set |
| §13.4.1.1.2 | Service Request Process | ¶20 (1 S) | feito | MVP-11 (`nrf/0005`, `openapi/0006`; [S10] a [S14]) |
| §13.4.1.1.3 | Access token requests in deployments with severa | ¶1 (1 S) | não se aplica | um NRF |
| §14.3.2 | Nnrf\_AccessToken\_Get Service Operation | ¶1 (1 S) | feito | `nrf/0005` |
| §A.7.0 | Parameters for the input S to the KDF | ¶1 a ¶2 (2 S) | free5GC | K_AMF do free5GC |
| §A.7.0 | Parameters for the input S to the KDF | ¶3 (1 S) | **desvio** | o P0 do SUPI, que o AKMA usa: só IMSI (DV-3). *"For P0, when the SUPI type is IMSI, P0 shall be set to IMSI as defined in clause 2.2 of TS 23.003 [19]. For …"* |

### 4.2 TS 29.503 V20.1.0 (as operações que o AKMA usa)

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §5.2.2.2.10 | Identifier Translation | ¶1 a ¶7 (7 S) | **não feito** | F10; MVP-10. *"Figure 5.2.2.2.10-1 shows a scenario where the NF service consumer (e.g. NEF, AAnF) sends a request to the …"* |
| §5.3.2.15.2 | Reauthentication Notify | ¶1 (1 S) | feito | MVP-12 (`udm/0002`, `amf/0001`; [R6r]); POST (DV-12) |
| §5.3.2.16.2 | Auth Trigger | ¶1 a ¶2 (2 S) | feito | MVP-12 (`udm/0002`, `amf/0001`; [R6r]); POST (DV-12) |
| §5.4.2.2.2 | Authentication Information Retrieval | ¶1 (1 S) | free5GC | AUSF do free5GC |
| §5.4.2.2.2 | Authentication Information Retrieval | ¶2 a ¶3 (2 S) | free5GC | UDM do free5GC |
| §5.4.2.2.3 | FN-RG Authentication | ¶1 a ¶2 (1 S, 1 R) | free5GC | UDM do free5GC |
| §5.5.2.2.2 | Subscription to Notification of event occurrence | ¶1 a ¶2 (2 S) | **não feito** | o UDM não atende `ROAMING_STATUS` por SUPI (F30); o AAnF não subscreve (DV-2). *"The NF service consumer sends a POST request to the parent resource (collection of subscriptions) …"* |
| §5.5.2.2.2 | Subscription to Notification of event occurrence | ¶3 a ¶6 (4 S) | não se aplica | grupos, EPC e outros eventos |
| §5.5.2.2.2 | Subscription to Notification of event occurrence | ¶7 (1 S) | **não feito** | o UDM não atende `ROAMING_STATUS` por SUPI (F30); o AAnF não subscreve (DV-2). *"The response, based on operator policy, may contain the expiry time, as determined by the UDM, after which …"* |
| §5.5.2.2.2 | Subscription to Notification of event occurrence | ¶8 a ¶20 (13 S) | não se aplica | grupos, EPC e outros eventos |
| §5.5.2.2.2 | Subscription to Notification of event occurrence | ¶21 a ¶22 (2 S) | **não feito** | o UDM não atende `ROAMING_STATUS` por SUPI (F30); o AAnF não subscreve (DV-2). *"If the user does not exist, HTTP status code "404 Not Found" shall be returned including additional error …"* |
| §5.5.2.2.2 | Subscription to Notification of event occurrence | ¶23 (1 S) | não se aplica | grupos, EPC e outros eventos |
| §5.5.2.2.2 | Subscription to Notification of event occurrence | ¶24 a ¶25 (2 S) | **não feito** | o UDM não atende `ROAMING_STATUS` por SUPI (F30); o AAnF não subscreve (DV-2). *"On failure, the appropriate HTTP status code indicating the error shall be returned and appropriate …"* |
| §5.5.2.3.2 | Unsubscribe to notifications of event occurrence | ¶1 a ¶2 (2 S) | **não feito** | o UDM não atende `ROAMING_STATUS` por SUPI (F30); o AAnF não subscreve (DV-2). *"If there is no valid subscription available (e.g. due to an unknown SubscriptionId value), HTTP status code …"* |
| §5.5.2.4.2 | Event Occurrence Notification | ¶1 a ¶2 (2 S) | **não feito** | o UDM não atende `ROAMING_STATUS` por SUPI (F30); o AAnF não subscreve (DV-2). *"The UDM sends a POST request to the callbackReference URI as provided by the NF service consumer during the …"* |
| §5.5.2.4.3 | Monitoring Revocation Notification | ¶1 a ¶5 (5 S) | não se aplica | revogação, restauração e atualização da subscrição |
| §5.5.2.4.4 | UDR-initiated Data Restoration | ¶1 a ¶2 (1 S, 1 R) | não se aplica | revogação, restauração e atualização da subscrição |
| §5.5.2.4.5 | Subscription Update Notification | ¶1 a ¶2 (2 S) | não se aplica | revogação, restauração e atualização da subscrição |
| §6.1.3.12.2 | Resource Definition | ¶1 (1 S) | **não feito** | F10; MVP-10. *"This resource shall support the resource URI variables defined in table 6.1.3.12.2-1."* |
| §6.1.3.12.3.1 | GET | ¶1 a ¶6 (5 S, 1 R) | **não feito** | F10; MVP-10. *"This method shall support the URI query parameters specified in table 6.1.3.12.3.1-1."* |
| §6.2.3.13.2 | Resource Definition | ¶1 (1 S) | feito | MVP-12 (`udm/0002`, `amf/0001`; [R6r]); POST (DV-12) |
| §6.2.3.13.4.2.2 | Operation Definition | ¶1 a ¶2 (2 S) | feito | MVP-12 (`udm/0002`, `amf/0001`; [R6r]); POST (DV-12) |
| §6.2.5.6 | Re-AuthenticationNotification Notification | ¶1 a ¶3 (3 S) | feito | MVP-12 (`udm/0002`, `amf/0001`; [R6r]); POST (DV-12) |
| §6.3.6.2.3 | Type: AuthenticationInfoResult | ¶1 (1 S) | free5GC | UDM do free5GC |
| §6.3.6.2.3 | Type: AuthenticationInfoResult | ¶2 (1 S) | feito | `udm/0001`; o `routingId` vem do UDR e é comparado com o do SUCI (D10) |
| §6.3.6.2.3 | Type: AuthenticationInfoResult | ¶3 (1 S) | free5GC | UDM do free5GC |
| §6.3.6.2.3 | Type: AuthenticationInfoResult | ¶4 (1 S) | feito | `udm/0001`; o `routingId` vem do UDR e é comparado com o do SUCI (D10) |
| §6.3.6.2.3 | Type: AuthenticationInfoResult | ¶5 a ¶6 (2 S) | free5GC | UDM do free5GC |
| §6.4.3.2.2 | Resource Definition | ¶1 a ¶2 (2 S) | **não feito** | F30; DV-2. *"This resource shall support the resource URI variables defined in table 6.4.3.2.2-1."* |
| §6.4.3.2.3.1 | POST | ¶1 a ¶4 (4 S) | **não feito** | F30; DV-2. *"This method shall support the URI query parameters specified in table 6.4.3.2.3.1-1."* |
| §6.4.6.2.8 | Type: RoamingStatusReport | ¶1 a ¶3 (3 S) | **não feito** | F30; DV-2. *"newServingPlmn / PlmnId / M / 1 / The new Serving PLMN. / If the "purged" attribute is present with value …"* |

### 4.3 TS 33.220 V19.1.0 (Anexos B, H e P)

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §B.2.0 | General | ¶1 a ¶3 (3 S) | feito | `shared/akma/kdf` e o oráculo Python |
| §B.2.1.2 | Character string encoding | ¶1 (1 S) | **desvio** | a biblioteca Go aplica NFKC; o ME recusa FQDN não ASCII (DV-4). *"A character string shall be encoded to an octet string according to UTF-8 encoding rules as specified in IETF …"* |
| §B.2.1.3 | Non-negative integer encoding | ¶1 a ¶2 (2 S) | feito | `shared/akma/kdf` e o oráculo Python |
| §B.2.2 | FC value allocations | ¶1 (1 S) | feito | `shared/akma/kdf` e o oráculo Python |
| §B.3 | NAF specific key derivation in GBA and GBA\_U | ¶1 a ¶2 (2 S) | não se aplica | GBA |
| §B.4 | Derivation of TMPI | ¶1 a ¶2 (2 S) | não se aplica | GBA |
| §B.5 | Derivation of passwd and Ks | ¶1 a ¶3 (3 S) | não se aplica | GBA |
| §B.6 | NAF specific key derivation in GBA\_Digest | ¶1 a ¶2 (2 S) | não se aplica | GBA |
| §H.1 | Definition | ¶1 (1 S) | feito | ADR-0003; `shared/akma/afid` |
| §H.2 | Organization Octet | ¶1 a ¶2 (2 S) | feito | ADR-0003; `shared/akma/afid` |
| §H.3 | Ua security protocol identifiers for 3GPP specif | ¶1 (1 S) | feito | ADR-0003; `shared/akma/afid` |
| §P.2.2 | Requirements on the UE | ¶1 a ¶2 (2 S) | feito | pelo Anexo D de TS 33.535 (MVP-8, fase d) |
| §P.2.3 | Requirements on the NAF | ¶1 (1 S) | feito | pelo Anexo D de TS 33.535 (MVP-8, fase d) |
| §P.2.3 | Requirements on the NAF | ¶2 (1 S) | não se aplica | chaves do GBA_U |
| §P.3.2 | Procedures | ¶1 a ¶3 (3 S) | feito | pelo Anexo D de TS 33.535 (MVP-8, fase d) |
| §P.3.2 | Procedures | ¶4 (1 S) | não se aplica | lista KISS, do GBA |
| §P.3.2 | Procedures | ¶5 a ¶6 (2 S) | feito | pelo Anexo D de TS 33.535 (MVP-8, fase d) |
| §P.3.3 | OSCORE Security context | ¶1 a ¶3 (3 S) | feito | pelo Anexo D de TS 33.535 (MVP-8, fase d) |
| §P.3.4 | Refresh of OSCORE key material | ¶1 (1 S) | não se aplica | renovação pelo Ub, do GBA |
| §P.3.5 | OSCORE Ua protocol payload encoding | ¶1 a ¶4 (4 S) | feito | pelo Anexo D de TS 33.535 (MVP-8, fase d) |

### 4.4 TS 33.222 V19.0.0 (§5.3 e §5.4)

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §5.3.0 | Procedures | ¶1 (0 S, 1 R) | não se aplica | notas sobre os modos do GBA |
| §5.3.0 | Procedures | ¶2 a ¶4 (3 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.3.0 | Procedures | ¶5 a ¶6 (2 S) | não se aplica | modos e USS do GBA |
| §5.3.0 | Procedures | ¶7 a ¶8 (2 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.3.0 | Procedures | ¶9 (1 S) | não se aplica | modos e USS do GBA |
| §5.3.0 | Procedures | ¶10 (1 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.3.0 | Procedures | ¶11 (1 S) | não se aplica | modos e USS do GBA |
| §5.3.1.0 | General | ¶1 (1 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.3.1.0 | General | ¶2 (0 S, 1 R) | feito | a revogação é opcional no Ua\*; o OCSP do AF fica para o MVP-18 (DV-10) |
| §5.3.1.0 | General | ¶3 (1 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.3.1.3 | Authentication of the AP/AS | ¶1 (1 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.3.1.3 | Authentication of the AP/AS | ¶2 (1 S) | **em parte** | achado 12: o certificado do AF do K.2 fica fora de TS 33.310 §6.1.1. *"The AP/AS certificate profile shall comply with the requirements for TLS certificates in clause 6.1 of TS …"* |
| §5.3.1.4 | Authentication Failures | ¶1 a ¶2 (2 S) | feito | sem certificado de cliente; sem reautenticação pela rede |
| §5.3.1.5 | Set-up of Security parameters | ¶1 (1 S) | feito | *tickets* presos à expiração da K_AF (U8) |
| §5.3.1.6 | Error cases | ¶1 (1 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.4.0.0 | General | ¶1 a ¶2 (2 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.4.0.0 | General | ¶3 (0 S, 1 R) | não se aplica | notas sobre os modos do GBA |
| §5.4.0.1 | TLS 1.2 | ¶1 (1 S) | feito | `[INTERPRETAÇÃO]` outras suítes no `ClientHello` (`ua-star-norma.md` §3.8, item 5) |
| §5.4.0.1 | TLS 1.2 | ¶2 (0 S, 1 R) | não se aplica | notas sobre os modos do GBA |
| §5.4.0.1 | TLS 1.2 | ¶3 a ¶6 (4 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.4.0.1 | TLS 1.2 | ¶7 (1 S) | não se aplica | GBA: Ks_NAF, USS, nova inscrição |
| §5.4.0.1 | TLS 1.2 | ¶8 a ¶9 (2 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.4.0.1 | TLS 1.2 | ¶10 (1 S) | não se aplica | GBA: Ks_NAF, USS, nova inscrição |
| §5.4.0.2 | TLS 1.3 | ¶1 (1 S) | feito | `[INTERPRETAÇÃO]` outras suítes no `ClientHello` (`ua-star-norma.md` §3.8, item 5) |
| §5.4.0.2 | TLS 1.3 | ¶2 (1 S) | não se aplica | GBA: Ks_NAF, USS, nova inscrição |
| §5.4.0.2 | TLS 1.3 | ¶3 (1 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.4.0.2 | TLS 1.3 | ¶4 (1 S) | não se aplica | GBA: Ks_NAF, USS, nova inscrição |
| §5.4.0.2 | TLS 1.3 | ¶5 (0 S, 1 R) | não se aplica | notas sobre os modos do GBA |
| §5.4.0.2 | TLS 1.3 | ¶6 (1 S) | não se aplica | GBA: Ks_NAF, USS, nova inscrição |
| §5.4.0.2 | TLS 1.3 | ¶7 a ¶8 (2 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.4.0.2 | TLS 1.3 | ¶9 (1 S) | não se aplica | GBA: Ks_NAF, USS, nova inscrição |
| §5.4.0.2 | TLS 1.3 | ¶10 (1 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.4.1.0 | General | ¶1 a ¶2 (2 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.4.1.1 | Protection mechanisms | ¶1 (1 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.4.1.3 | Authentication Failures | ¶1 (1 S) | feito | pelo Anexo B de TS 33.535 (MVP-8, fase b; `ua-star-norma.md` §3) |
| §5.4.1.3 | Authentication Failures | ¶2 (1 S) | não se aplica | GBA: Ks_NAF, USS, nova inscrição |
| §5.4.1.4 | Set-up of Security parameters | ¶1 (1 S) | feito | *tickets* presos à expiração da K_AF (U8) |

### 4.5 TS 33.122 V19.4.0 (§4, §6.1, §6.3, §6.5, §6.6 a §6.8, Anexo C)

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §4.2 | Common security requirements | ¶1 (1 S) | **em parte** | o NEF é o único ponto de entrada; o mecanismo de ocultação (TS 23.222 §8.13) não existe (DV-14). *"[CAPIF-SEC-4.2-a] The CAPIF shall provide mechanisms to hide the topology of the PLMN trust domain from the …"* |
| §4.2 | Common security requirements | ¶2 a ¶3 (2 S) | não se aplica | provedor de terceiros |
| §4.2 | Common security requirements | ¶4 (1 S) | feito | TLS em todas as interfaces do CAPIF |
| §4.2 | Common security requirements | ¶5 (1 S) | **em parte** | uma implantação, o CCF ao lado do NEF (C1); as outras, não (DV-14). *"[CAPIF-SEC-4.2-e] API invoker authentication and authorization shall support all deployment models listed in …"* |
| §4.2 | Common security requirements | ¶6 (0 S, 1 R) | feito | TLS em todas as interfaces do CAPIF |
| §4.3 | Security requirements on the CAPIF-1/1e referenc | ¶1 a ¶10 (10 S) | feito | CAPIF-1e com TLS mútuo e a credencial de entrada ([CP00] a [CP03]) |
| §4.4 | Security requirements on the CAPIF-2/2e referenc | ¶1 a ¶7 (7 S) | feito | CAPIF-2e com TLS mútuo, o *token* e o `check-authentication` ([CP04], [S16.check]) |
| §4.5 | Security requirements on the CAPIF-3/4/5 referen | ¶1 a ¶5 (5 S) | feito | CAPIF-3/4 com TLS mútuo; o provedor pelo certificado ([S17.nef]) |
| §4.5 | Security requirements on the CAPIF-3/4/5 referen | ¶6 (1 S) | **não feito** | sem a entrada com aprovação (DV-14). *"[CAPIF-SEC-4.5-f] The CAPIF core function shall be able to request explicit grant of new API invoker’s …"* |
| §4.5 | Security requirements on the CAPIF-3/4/5 referen | ¶7 a ¶8 (2 S) | **não feito** | sem CAPIF-5 (DV-14). *"[CAPIF-SEC-4.5-g] The CAPIF core function shall be able to authenticate the API Management function’s …"* |
| §4.6 | Security requirements on the CAPIF-3e/4e/5e refe | ¶1 a ¶8 (8 S) | não se aplica | outro domínio, interconexão, ocultação e RNAA |
| §4.6A | Security requirements on the CAPIF-6/6e referenc | ¶1 a ¶4 (4 S) | não se aplica | outro domínio, interconexão, ocultação e RNAA |
| §4.7 | Security requirements on the CAPIF-7/7e referenc | ¶1 a ¶5 (5 S) | não se aplica | outro domínio, interconexão, ocultação e RNAA |
| §4.8 | Security requirements on the CAPIF-8 reference p | ¶1 a ¶4 (4 S) | não se aplica | outro domínio, interconexão, ocultação e RNAA |
| §6.1 | Security procedures for API invoker onboarding | ¶1 a ¶6 (6 S) | feito | [CP00], [CP01], [CP08] |
| §6.3.1.1 | General | ¶1 a ¶2 (2 S) | feito | certificados de §6.1.3a e §6.1.4a, com CRL (MVP-11) |
| §6.3.1.1 | General | ¶3 (1 S) | **em parte** | TLS no perfil, salvo o OCSP *stapling* (DV-16). *"TLS [9] shall be used to provide integrity protection, replay protection and confidentiality protection for …"* |
| §6.3.1.2 | Security method negotiation | ¶1 a ¶4 (4 S) | feito | [CP03]; perfis de AF no CCF (C20) |
| §6.3.1.3 | API discovery | ¶1 (1 S) | feito | [CP03]; perfis de AF no CCF (C20) |
| §6.3.1.4 | Topology hiding | ¶1 (1 S) | não se aplica | sem ocultação |
| §6.5.1 | General | ¶1 (1 S) | feito | método 3 ([CP04] a [CP07]); TLS mútuo no passo 5 (D-AJ) |
| §6.5.2.3 | Method 3 – TLS with OAuth token | ¶1 a ¶11 (11 S) | feito | método 3 ([CP04] a [CP07]); TLS mútuo no passo 5 (D-AJ) |
| §6.6 | Security procedures for CAPIF-3/4/5/6 reference  | ¶1 (1 S) | **em parte** | TLS no perfil, salvo o OCSP *stapling* (achado 3). *"TLS shall be used to provide integrity protection, replay protection and confidentiality protection. The …"* |
| §6.6 | Security procedures for CAPIF-3/4/5/6 reference  | ¶2 (1 S) | feito | MVP-11, fase (c) ([S17.nef]) |
| §6.7 | Security procedures for updating security method | ¶1 (1 S) | **não feito** | sem troca do método de segurança (DV-14). *"As specified in TS 23.222 [3], the CAPIF core function shall receive updates to AEF authentication and …"* |
| §6.8 | Security procedure for API invoker offboarding | ¶1 a ¶4 (4 S) | feito | [CP18], [CP19] |
| §6.8 | Security procedure for API invoker offboarding | ¶5 a ¶8 (4 S) | **não feito** | MVP-17 (Q7). *"The CAPIF core function shall send Event notification message to the API exposing function to indicate that …"* |
| §C.1 | General | ¶1 a ¶5 (5 S) | feito | MVP-9 (`capif.md` §5.7 e §5.8); DV-13 |
| §C.2.2 | Token claims | ¶1 (1 S) | feito | MVP-9 (`capif.md` §5.7 e §5.8); DV-13 |
| §C.2.2 | Token claims | ¶2 (1 S) | não se aplica | RNAA e *redirect_uri* |
| §C.2.2 | Token claims | ¶3 a ¶4 (2 S) | feito | MVP-9 (`capif.md` §5.7 e §5.8); DV-13 |
| §C.3.2 | Access token request | ¶1 a ¶2 (2 S) | feito | MVP-9 (`capif.md` §5.7 e §5.8); DV-13 |
| §C.3.2 | Access token request | ¶3 (1 S) | não se aplica | RNAA e *redirect_uri* |
| §C.3.3 | Access token response | ¶1 (1 S) | feito | `Bearer`; `[INTERPRETAÇÃO]` o valor não distingue maiúsculas (RFC 6749 §5.1) |
| §C.4.1 | Client Credentials Grant | ¶1 (1 S) | feito | MVP-9 (`capif.md` §5.7 e §5.8); DV-13 |
| §C.5 | Using the token to access API exposing functions | ¶1 a ¶2 (2 S) | feito | MVP-9 (`capif.md` §5.7 e §5.8); DV-13 |
| §C.6 | Token revocation | ¶1 a ¶3 (3 S) | feito | MVP-9 (`capif.md` §5.7 e §5.8); DV-13 |
| §C.7.1 | Access token validation | ¶1 (1 S) | feito | MVP-9 (`capif.md` §5.7 e §5.8); DV-13 |

### 4.6 TS 23.222 V20.1.0 (§6.3 e as cláusulas de §8 que o projeto usa)

O Stage 2 descreve os procedimentos quase sem *shall*: as capacidades de §6.3.3 e §6.3.4 e os procedimentos de §8 estão na tabela da primeira passada e no DV-14.

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §8.1.2.2 | Onboard API invoker response | ¶1 (1 S) | feito | MVP-9 (`capif.md`) |
| §8.5.2.2 | Service API get response | ¶1 (1 S) | feito | MVP-9 (`capif.md`) |
| §8.7.2.1 | Service API discover request | ¶1 (1 S) | feito | MVP-9 (`capif.md`) |
| §8.7.2.2 | Service API discover response | ¶1 (1 S) | feito | MVP-9 (`capif.md`) |
| §8.8.2.2 | Event subscription response | ¶1 (1 S) | **não feito** | eventos: MVP-17 (Q7). *"NOTE: Shall be present if the Result information element indicates that the event subscription operation is …"* |
| §8.8.2.7 | Event subscription update request | ¶1 (1 S) | **não feito** | eventos: MVP-17 (Q7). *"NOTE: At least one of these information elements shall be present."* |
| §8.16.2.3 | Obtain API invoker information request | ¶1 (1 S) | feito | MVP-9 (`capif.md`) |
| §8.16.2.4 | Obtain API invoker information response | ¶1 (1 S) | feito | MVP-9 (`capif.md`) |
| §8.22.2.2 | Query service API log response | ¶1 (1 S) | **não feito** | auditoria: MVP-17 (Q7). *"NOTE: Information element shall be present when result indicates success."* |

### 4.7 TS 29.222 V20.1.0 (§5.2 a §5.9, §6.2 e §7)

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §5.2.2.2.2 | Consumer discovering service API using Discover\_ | ¶1 a ¶4 (4 S) | feito | `Discover_Service_API` ([CP02]) |
| §5.3.2.2.2 | API publishing function publishing service APIs  | ¶1 a ¶3 (3 S) | feito | `Publish`, `Unpublish` e `Get_Service_API` ([R1.nef-published]) |
| §5.3.2.2.2 | API publishing function publishing service APIs  | ¶4 (1 S) | não se aplica | ocultação de topologia |
| §5.3.2.2.2 | API publishing function publishing service APIs  | ¶5 (1 S) | feito | `Publish`, `Unpublish` e `Get_Service_API` ([R1.nef-published]) |
| §5.3.2.2.3 | CAPIF core function publishing service APIs on o | ¶1 a ¶6 (6 S) | não se aplica | interconexão de CCF |
| §5.3.2.3.2 | Consumer un-publishing service APIs from CAPIF c | ¶1 a ¶3 (3 S) | feito | `Publish`, `Unpublish` e `Get_Service_API` ([R1.nef-published]) |
| §5.3.2.3.2 | Consumer un-publishing service APIs from CAPIF c | ¶4 (1 S) | não se aplica | ocultação de topologia |
| §5.3.2.3.2 | Consumer un-publishing service APIs from CAPIF c | ¶5 (1 S) | feito | `Publish`, `Unpublish` e `Get_Service_API` ([R1.nef-published]) |
| §5.3.2.4.2 | Consumer retrieving service APIs from CAPIF core | ¶1 a ¶4 (4 S) | feito | `Publish`, `Unpublish` e `Get_Service_API` ([R1.nef-published]) |
| §5.3.2.5.2 | Consumer updating published service APIs on CAPI | ¶1 a ¶3 (3 S) | **não feito** | sem `Update_Service_API` (DV-14). *"To update information of published service APIs, the consumer (e.g. API publishing function) shall send an …"* |
| §5.3.2.5.2 | Consumer updating published service APIs on CAPI | ¶4 (1 S) | não se aplica | ocultação de topologia |
| §5.3.2.5.2 | Consumer updating published service APIs on CAPI | ¶5 (1 S) | **não feito** | sem `Update_Service_API` (DV-14). *"if errors occur when processing the request, the CAPIF core function shall respond to the consumer with an …"* |
| §5.4.2.2.2 | Subscribing to CAPIF events using Subscribe\_Even | ¶1 a ¶5 (5 S) | **não feito** | eventos: MVP-17 (Q7). *"To subscribe to CAPIF events reporting, the Subscriber shall send an HTTP POST request message to the CCF …"* |
| §5.4.2.3.2 | Unsubscribing from CAPIF events using Unsubscrib | ¶1 a ¶4 (4 S) | **não feito** | eventos: MVP-17 (Q7). *"To unsubscribe from CAPIF events, the Subscriber shall send an HTTP DELETE request message targeting the …"* |
| §5.4.2.4.2 | Notifying CAPIF events using Notify\_Event servic | ¶1 a ¶4 (4 S) | **não feito** | eventos: MVP-17 (Q7). *"To notify on CAPIF events, the CCF shall send an HTTP POST request message using the Notification Destination …"* |
| §5.4.2.5.2 | Update Subscription to CAPIF events using Update | ¶1 a ¶4 (4 S) | **não feito** | eventos: MVP-17 (Q7). *"To update an existing subscription to CAPIF events, the Subscriber shall send an HTTP PUT/PATCH request …"* |
| §5.5.2.2.2 | API Invoker on-boarding itself as a recognized u | ¶1 a ¶4 (4 S) | feito | entrada automática ([CP01]) |
| §5.5.2.2.2 | API Invoker on-boarding itself as a recognized u | ¶5 a ¶7 (3 S) | **não feito** | sem a entrada com aprovação, 202 (DV-14). *"the CCF shall acknowledge the reception of the on-boarding request to the API Invoker by returning an HTTP …"* |
| §5.5.2.2.2 | API Invoker on-boarding itself as a recognized u | ¶8 (1 S) | feito | entrada automática ([CP01]) |
| §5.5.2.3.2 | API Invoker off-boarding itself from being a rec | ¶1 a ¶4 (4 S) | feito | saída ([CP18]) |
| §5.5.2.4.2 | Notifying API Invoker's onboarding creation/upda | ¶1 a ¶3 (3 S) | não se aplica | só com a entrada com aprovação |
| §5.5.2.5.2 | API Invoker updating its details on CAPIF using  | ¶1 a ¶7 (7 S) | **não feito** | sem a atualização da entrada (DV-14). *"To update the API Invoker's onboarding details at the CCF, the service consumer (e.g., API Invoker) shall …"* |
| §5.6.2.2.1 | General | ¶1 (1 S) | feito | [CP03] a [CP05], [CP11] a [CP13] |
| §5.6.2.2.2 | Request service API security method from CAPIF u | ¶1 a ¶3 (3 S) | feito | [CP03] a [CP05], [CP11] a [CP13] |
| §5.6.2.3.1 | General | ¶1 (1 S) | feito | [CP03] a [CP05], [CP11] a [CP13] |
| §5.6.2.3.2 | Obtain authorization using Obtain\_Authorization  | ¶1 a ¶3 (3 S) | feito | [CP03] a [CP05], [CP11] a [CP13] |
| §5.6.2.3.2 | Obtain authorization using Obtain\_Authorization  | ¶4 a ¶5 (2 S) | não se aplica | RNAA |
| §5.6.2.3.2 | Obtain authorization using Obtain\_Authorization  | ¶6 a ¶8 (3 S) | feito | [CP03] a [CP05], [CP11] a [CP13] |
| §5.6.2.4.2 | Obtain API invoker's security information using  | ¶1 a ¶3 (3 S) | feito | [CP03] a [CP05], [CP11] a [CP13] |
| §5.6.2.5.2 | Invalidate authorization using Revoke\_Authorizat | ¶1 a ¶6 (6 S) | **não feito** | `Revoke_Authorization`: MVP-17 (Q7). *"To invalidate authorization of an API invoker for all service APIs, the service consumer shall send an HTTP …"* |
| §5.8.2.2.2 | Logging service API invocations using Log\_API\_In | ¶1 a ¶4 (4 S) | **não feito** | `Log_API_Invocation`: MVP-17 (Q7). *"To log service API invocations at the CAPIF core function, the API exposing function shall send an HTTP POST …"* |
| §5.9.2.2.2 | Query API invocation information logs using Quer | ¶1 a ¶4 (4 S) | **não feito** | `Query_API_Invocation_Log`: MVP-17 (Q7). *"To query service API invocation logs at the CAPIF core function, the API management function shall send an …"* |
| §6.2.2.2.2 | API invoker initiating authentication using Init | ¶1 a ¶3 (3 S) | feito | `check-authentication` ([CP04]) |
| §6.2.2.3.2 | CAPIF core function initiating revocation using  | ¶1 a ¶3 (3 S) | **não feito** | revogação no AEF: MVP-17 (Q7). *"To revoke authorization, the CAPIF core function shall send an HTTP POST message to the API exposing function …"* |
| §7.1 | General | ¶1 a ¶2 (2 S) | feito | TS 29.122 §5.2; HTTP/1.1 e HTTP/2 sobre TLS |
| §7.3 | Usage of HTTP | ¶1 (1 S) | feito | TS 29.122 §5.2; HTTP/1.1 e HTTP/2 sobre TLS |
| §7.3 | Usage of HTTP | ¶2 (1 S) | **desvio** | DV-15. *"A functional entity desiring to use HTTP/2 shall use the HTTP upgrade mechanism to negotiate applicable HTTP …"* |
| §7.4 | Content type | ¶1 (1 S) | feito | TS 29.122 §5.2; HTTP/1.1 e HTTP/2 sobre TLS |
| §7.5.1 | Resource URI structure | ¶1 (1 S) | feito | TS 29.122 §5.2; HTTP/1.1 e HTTP/2 sobre TLS |
| §7.5.2 | Custom operations URI structure | ¶1 (1 S) | feito | TS 29.122 §5.2; HTTP/1.1 e HTTP/2 sobre TLS |
| §7.6 | Notifications | ¶1 (1 S) | **não feito** | sem notificações do CAPIF: MVP-17 (Q7). *"shall support the delivery of notifications using a separate HTTP connection towards an address;"* |
| §7.8 | Feature negotiation | ¶1 (1 S) | feito | TS 29.122 §5.2; HTTP/1.1 e HTTP/2 sobre TLS |
| §7.9 | HTTP custom headers | ¶1 (1 S) | não se aplica | os cabeçalhos de TS 29.122 §5.2.8 são os de carga e sobrecarga, opcionais |
| §7.10 | Conventions for Open API specification files | ¶1 (1 S) | não se aplica | convenção para quem escreve o OpenAPI |
| §7.11 | CAPIF vendor-specific extensions | ¶1 a ¶2 (2 S) | feito | `[INTERPRETAÇÃO]` o CCF ignora atributos e parâmetros que não conhece |

### 4.8 TS 33.210 V20.0.0 (§6.2 e §6.3)

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §6.2.1 | General | ¶1 (0 S, 1 R) | não se aplica | introdução |
| §6.2.1 | General | ¶2 a ¶6 (5 S) | feito | `crypto/tls` do Go (`sbitls`, CCF, AF) e wolfSSL (Ua\*): `tls-sbi.md` §6, `ua-star-norma.md` §3.8 |
| §6.2.2 | Profiling for TLS 1.3 | ¶1 a ¶9 (7 S, 2 R) | feito | `crypto/tls` do Go (`sbitls`, CCF, AF) e wolfSSL (Ua\*): `tls-sbi.md` §6, `ua-star-norma.md` §3.8 |
| §6.2.2 | Profiling for TLS 1.3 | ¶10 (0 S, 1 R) | **não feito** | OCSP *stapling*: achado 3, DV-10, DV-16; MVP-18. *"The TLS Certificate Status Request extension (i.e., "OCSP stapling") as defined in RFC 6066 [57] and RFC 8446 …"* |
| §6.2.2 | Profiling for TLS 1.3 | ¶11 (1 S) | feito | `crypto/tls` do Go (`sbitls`, CCF, AF) e wolfSSL (Ua\*): `tls-sbi.md` §6, `ua-star-norma.md` §3.8 |
| §6.2.3 | Profiling for TLS 1.2 | ¶1 a ¶16 (14 S, 2 R) | feito | `crypto/tls` do Go (`sbitls`, CCF, AF) e wolfSSL (Ua\*): `tls-sbi.md` §6, `ua-star-norma.md` §3.8 |
| §6.2.3 | Profiling for TLS 1.2 | ¶17 (0 S, 1 R) | **não feito** | OCSP *stapling*: achado 3, DV-10, DV-16; MVP-18. *"The OCSP Status (a.k.a. certificate status request) extension, defined in RFC 6066 [57] should be supported."* |
| §6.2.3 | Profiling for TLS 1.2 | ¶18 (1 S) | **em parte** | `D0 01`, sim; `D0 02`, não (DV-10). *"If pre-shared key (psk) cipher suites are implemented in TLS, then RFC 5489 [64] shall apply and the …"* |
| §6.3.2.1 | General | ¶1 a ¶5 (5 S) | não se aplica | sem JWE |
| §6.3.2.2 | N32 interface | ¶1 (1 S) | não se aplica | sem JWE |
| §6.3.3.1 | General | ¶1 a ¶6 (6 S) | feito | `ES256` com `kid`; `none` e `jwk` recusados (`openapi/0006`, `free5gc/ccf/internal/jws`) |
| §6.3.3.1 | General | ¶7 (0 S, 1 R) | feito | o `kid` só é comparado com a impressão da chave; `x5u` e `jku` não são seguidos (RFC 8725 §3.10) |
| §6.3.3.2 | N32 interface | ¶1 (1 S) | não se aplica | N32 e CCA |
| §6.3.3.3 | Client credentials assertion | ¶1 a ¶2 (2 S) | não se aplica | N32 e CCA |
| §6.3.3.4 | Access token | ¶1 (0 S, 1 R) | **em parte** | o NRF manda `x5t#S256`; o CCF, não (DV-16; MVP-15). *"The "x5t#S256" header parameter should be supported."* |

### 4.9 TS 33.310 V19.5.0 (§6.1, §6.1a, §6.1b, §6.2a, §6.3, §7.1 e §10)

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §6.1.1 | Common rules to all certificates | ¶1 a ¶7 (7 S) | feito | `make-sba-certs.sh` e `make-capif-certs.sh`: ECDSA P-256 e P-384, SHA-256 e SHA-384 |
| §6.1.1 | Common rules to all certificates | ¶8 (1 S) | **em parte** | SBA e CCF, sim; o certificado do AF do K.2 não tem O (achado 12). *"(C=<country>), O=<Organization Name>, CN=<Some distinguishing name>. Organization and CN shall be in UTF8 …"* |
| §6.1.1 | Common rules to all certificates | ¶9 (1 S) | feito | `make-sba-certs.sh` e `make-capif-certs.sh`: ECDSA P-256 e P-384, SHA-256 e SHA-384 |
| §6.1.1 | Common rules to all certificates | ¶10 (0 S, 1 R) | **não feito** | OCSP: MVP-18. *"OCSP as specified in subclause 6.1b should be supported for certificate revocation verification."* |
| §6.1.1 | Common rules to all certificates | ¶11 a ¶12 (2 S) | **em parte** | SBA e CCF, sim; o do AF do K.2 marca `basicConstraints` como crítica (achado 12). *"Certificate extensions which are not mandated by this specification but which are mentioned within RFC5280 …"* |
| §6.1.3a | TLS entity certificate profile | ¶1 (1 S) | feito | certificados do CCF e dos *invokers* (`make-capif-certs.sh`) |
| §6.1.3a | TLS entity certificate profile | ¶2 (0 S, 1 R) | não se aplica | SIP |
| §6.1.3a | TLS entity certificate profile | ¶3 a ¶4 (2 S) | feito | certificados do CCF e dos *invokers* (`make-capif-certs.sh`) |
| §6.1.3c.2 | General SBA Certificate profile | ¶1 (1 S) | feito | `make-sba-certs.sh`, que confere o perfil ([S02]) |
| §6.1.3c.3 | NF Certificate profile | ¶1 a ¶13 (11 S, 2 R) | feito | `make-sba-certs.sh`, que confere o perfil ([S02]) |
| §6.1.3c.4 | SCP certificate profile | ¶1 a ¶3 (1 S, 2 R) | não se aplica | SCP e SEPP |
| §6.1.3c.5.2 | SEPP intra-domain certificate profile | ¶1 (1 S) | não se aplica | SCP e SEPP |
| §6.1.3c.5.3.2 | SEPP inter-domain certificate profile for inter- | ¶1 (1 S) | não se aplica | SCP e SEPP |
| §6.1.4a | TLS client/server CA certificate profile | ¶1 (0 S, 1 R) | feito | CA do operador e CA do CCF |
| §6.1a | CRL profile | ¶1 a ¶6 (6 S) | feito | CRL v2 assinada com ECDSA P-384 (`free5gc/crl`, CCF) |
| §6.1a | CRL profile | ¶7 (1 S) | **desvio** | DV-18; §7.1 deixa as entidades TLS usarem LDAP ou HTTP (achado 15). *"CRL retrieval with LDAPv3�(RFC 4510 [67], RFC 4517 [68], RFC 4523 [69], RFC 4512 [70]) shall be supported as …"* |
| §6.1b | OCSP profile | ¶1 a ¶7 (7 S) | **não feito** | OCSP: MVP-18. *"Hash algorithm for use before signing OCSP response: SHA-256 and SHA-384 shall be supported, MD5 MD2, and …"* |
| §6.2a | TLS profiling | ¶1 (1 S) | não se aplica | TLS entre operadores |
| §6.2a.1 | TLS profile | ¶1 a ¶3 (3 S) | feito | comportamento do `crypto/tls` |
| §6.3.1 | Path validation profiling | ¶1 a ¶4 (4 S) | feito | `sbitls`, falha fechada (`TestRevocation`; [S06], [S07]) |
| §6.3.1 | Path validation profiling | ¶5 (1 S) | não se aplica | NE e SEG |
| §6.3.1 | Path validation profiling | ¶6 (1 S) | feito | `sbitls`, falha fechada (`TestRevocation`; [S06], [S07]) |
| §7.1 | Repositories | ¶1 (1 S) | feito | CRL por HTTP, fora do TLS; achado 15 |
| §7.1 | Repositories | ¶2 (0 S, 1 R) | não se aplica | repositórios de CA combinadas |
| §7.1 | Repositories | ¶3 (1 S) | feito | CRL por HTTP, fora do TLS; achado 15 |
| §10.2.2 | Architecture | ¶1 a ¶4 (4 S) | **não feito** | inscrição automática: os certificados saem dos *scripts*. *"The assumption is that the OAM system is trusted for the operator CA/RA, i.e., the trust between the OAM …"* |
| §10.2.3 | Procedure | ¶1 a ¶14 (14 S) | **não feito** | inscrição automática: os certificados saem dos *scripts*. *"If the initial trust has been established by initial digital certificate or OAM signature, the public root …"* |
| §10.3.1.1 | General Requirements | ¶1 a ¶12 (12 S) | **não feito** | inscrição automática: os certificados saem dos *scripts*. *"The following requirements shall apply to CMPv2 usage in Service Based Architecture:"* |
| §10.3.1.2 | Profile for PKIMessage | ¶1 a ¶2 (2 S) | **não feito** | inscrição automática: os certificados saem dos *scripts*. *"The support and usage of the optional protection field of type PKIProtection is required by this profile. The …"* |
| §10.3.1.3 | Profile for PKIHeader Field | ¶1 a ¶5 (5 S) | **não feito** | inscrição automática: os certificados saem dos *scripts*. *"The sender field shall contain the identity of the NF as the end entity. This identity shall be identical to …"* |
| §10.3.1.4.1 | General | ¶1 (1 S) | **não feito** | inscrição automática: os certificados saem dos *scripts*. *"The NF Instance certificate enrolment shall support the following CMPv2 PKI message bodies:"* |
| §10.3.1.4.2 | Initialization Request | ¶1 a ¶5 (5 S) | **não feito** | inscrição automática: os certificados saem dos *scripts*. *"The Initialization Request as specified in IETF RFC 4210 [4] shall contain exactly one CertReqMessages as …"* |
| §10.3.1.4.3 | Initialization Response | ¶1 a ¶4 (4 S) | **não feito** | inscrição automática: os certificados saem dos *scripts*. *"The Initialization Response as specified in RFC 4210 [4] shall contain exactly one generated NF certificate, …"* |
| §10.3.1.4.4 | Certification request and Certification Response | ¶1 a ¶3 (2 S, 1 R) | **não feito** | inscrição automática: os certificados saem dos *scripts*. *"The structure and content of these messages is identical to initialization requests and responses, thus the …"* |
| §10.3.1.4.5 | Key Update Request and Key Update Response | ¶1 a ¶2 (1 S, 1 R) | **não feito** | inscrição automática: os certificados saem dos *scripts*. *"The PKIMessage sent by the NF shall be signed with the private key which is related to the last received …"* |
| §10.3.1.4.6 | Certificate Confirm Request and Confirmation Res | ¶1 a ¶3 (3 S) | **não feito** | inscrição automática: os certificados saem dos *scripts*. *"Initialization responses and key update responses shall always be followed by a Certificate Confirm request …"* |
| §10.3.2 | CMPv2 transport | ¶1 a ¶2 (2 S) | **não feito** | inscrição automática: os certificados saem dos *scripts*. *"Transport of CMPv2 messages between end entities (network elements) and RA/CA shall be done using HTTP-based …"* |
| §10.3.3 | Trusted Network Function instances identifiers | ¶1 (0 S, 1 R) | **não feito** | inscrição automática: os certificados saem dos *scripts*. *"Operator RA/CA should be able to verify that the nfinstanceID in the certificate signing request (�ir� and …"* |

## 5. Normas que a lista não trazia

### 5.1 TS 29.510 V20.1.0 (registro, *heart-beat*, descoberta, *token* e chave)

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §5.2.2.2.2 | NF (other than NRF) registration to NRF | ¶1 (1 S) | feito | o AAnF registra-se com PUT ([R1]) |
| §5.2.2.2.2 | NF (other than NRF) registration to NRF | ¶2 (1 S) | feito | UUID em minúsculas |
| §5.2.2.2.2 | NF (other than NRF) registration to NRF | ¶3 a ¶4 (2 S) | feito | o AAnF registra-se com PUT ([R1]) |
| §5.2.2.2.2 | NF (other than NRF) registration to NRF | ¶5 a ¶11 (7 S) | free5GC | NRF do free5GC; `AANF` aceito desde o `nrf/0002` |
| §5.2.2.2.2 | NF (other than NRF) registration to NRF | ¶12 (0 S, 1 R) | **não feito** | o AAnF não manda OPTIONS antes do registro (*should*). *"Before an NF Instance registers its NF Profile in NRF, the NF Instance should check the capabilities of the …"* |
| §5.2.2.2.3 | NRF registration to another NRF | ¶1 a ¶3 (3 S) | não se aplica | NRF entre NRF e dados compartilhados |
| §5.2.2.2.4 | Shared Data registration to NRF | ¶1 a ¶8 (8 S) | não se aplica | NRF entre NRF e dados compartilhados |
| §5.2.2.3.1A | NF Profile Complete replacement | ¶1 a ¶6 (6 S) | não se aplica | o AAnF não atualiza o perfil |
| §5.2.2.3.1B | NF Profile Partial Update | ¶1 a ¶9 (8 S, 1 R) | não se aplica | o AAnF não atualiza o perfil |
| §5.2.2.3.1C | Shared Data Complete replacement | ¶1 a ¶8 (8 S) | não se aplica | NRF entre NRF e dados compartilhados |
| §5.2.2.3.1D | Shared data Partial Update | ¶1 a ¶10 (9 S, 1 R) | não se aplica | NRF entre NRF e dados compartilhados |
| §5.2.2.3.2 | NF Heart-Beat | ¶1 a ¶9 (8 S, 1 R) | **não feito** | sem *heart-beat* no AAnF e nas NFs do free5GC (achado 6, F29). *"Each NF that has previously registered in NRF shall contact the NRF periodically (heart-beat), by invoking …"* |
| §5.2.2.4.2 | NF Instance Deregistration | ¶1 a ¶4 (4 S) | feito | o AAnF descadastra-se ao parar; o NEF desde o `nef/0004` (F14) |
| §5.2.2.4.3 | Shared Data Deregistration | ¶1 a ¶5 (5 S) | não se aplica | NRF entre NRF e dados compartilhados |
| §5.3.2.2.1 | General | ¶1 a ¶2 (1 S, 1 R) | **em parte** | o AAnF guarda o UDM; o AUSF e o NEF descobrem o AAnF a cada pedido (achado 13). *"Before a service consumer invokes this service operation, it shall consider if it is possible to reuse the …"* |
| §5.3.2.2.1 | General | ¶3 (1 S) | não se aplica | SCP |
| §5.3.2.2.1 | General | ¶4 (0 S, 1 R) | **em parte** | o AAnF guarda o UDM; o AUSF e o NEF descobrem o AAnF a cada pedido (achado 13). *"The NF Service Consumer should avoid to send multiple NF discovery requests with the same query parameters to …"* |
| §5.3.2.2.2 | Service Discovery in the same PLMN | ¶1 a ¶6 (6 S) | free5GC | NRF do free5GC; o AAnF por RID (`nrf/0002`) |
| §5.3.2.2.2 | Service Discovery in the same PLMN | ¶7 a ¶14 (8 S) | não se aplica | outra PLMN, SCP, NF específicas |
| §5.3.2.2.3 | Service Discovery in a different PLMN | ¶1 a ¶8 (7 S, 1 R) | não se aplica | outra PLMN, SCP, NF específicas |
| §5.3.2.2.4 | Service Discovery with intermediate redirecting  | ¶1 a ¶6 (6 S) | não se aplica | outra PLMN, SCP, NF específicas |
| §5.3.2.2.5 | Service Discovery with intermediate forwarding N | ¶1 a ¶6 (6 S) | não se aplica | outra PLMN, SCP, NF específicas |
| §5.3.2.2.6 | Service Discovery with resolution of the target  | ¶1 a ¶2 (1 S, 1 R) | não se aplica | outra PLMN, SCP, NF específicas |
| §5.4.2.2.1 | General | ¶1 a ¶13 (12 S, 1 R) | feito | MVP-11 (`nrf/0005`; [S10], [S13]) |
| §5.4.2.2.2 | Access Token request with intermediate forwardin | ¶1 a ¶6 (6 S) | não se aplica | vários NRF |
| §5.4.2.2.3 | Access Token request with intermediate redirecti | ¶1 a ¶6 (6 S) | não se aplica | vários NRF |
| §5.4.2.3.1 | General | ¶1 a ¶4 (4 S) | **não feito** | sem `Nnrf_AccessToken_RetrieveKey` (achado 5, F28). *"The NF Service Consumer shall send a POST request to resource URI of the "Retrieve Key Request" custom …"* |
| §6.3.5.2.2 | Type: AccessTokenReq | ¶1 a ¶6 (6 S) | feito | `nrf/0005` |
| §6.3.5.2.2 | Type: AccessTokenReq | ¶7 a ¶25 (19 S) | não se aplica | PLMN, SNPN, fatias, NWDAF |
| §6.3.5.2.3 | Type: AccessTokenRsp | ¶1 a ¶4 (4 S) | feito | `nrf/0005` |
| §6.3.5.2.4 | Type: AccessTokenClaims | ¶1 a ¶5 (5 S) | feito | `nrf/0005` |
| §6.3.5.2.4 | Type: AccessTokenClaims | ¶6 a ¶15 (10 S) | não se aplica | PLMN, SNPN, fatias, NWDAF |
| §6.3.5.2.4 | Type: AccessTokenClaims | ¶16 (1 S) | feito | `nrf/0005` |
| §6.3.5.2.4 | Type: AccessTokenClaims | ¶17 a ¶18 (2 S) | não se aplica | PLMN, SNPN, fatias, NWDAF |

### 5.2 TS 29.500 V20.0.0 (§5, §6.6 e §6.7)

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §5.1 | Protocol Stack Overview | ¶1 (1 S) | feito | HTTP/2; `[INTERPRETAÇÃO]` o AAnF e o CCF também aceitam HTTP/1.1 |
| §5.2.1 | General | ¶1 (1 S) | feito | HTTP/2; `[INTERPRETAÇÃO]` o AAnF e o CCF também aceitam HTTP/1.1 |
| §5.2.2.1 | General | ¶1 (1 S) | feito | `Authorization`, `Content-Type`, `Allow`, `WWW-Authenticate` (F12) |
| §5.2.2.2 | Mandatory to support HTTP standard headers | ¶1 (1 S) | feito | `Authorization`, `Content-Type`, `Allow`, `WWW-Authenticate` (F12) |
| §5.2.2.2 | Mandatory to support HTTP standard headers | ¶2 (1 S) | **em parte** | o AAnF manda `User-Agent` com o tipo ao UDM; o AUSF, o NEF e o AF mandam ao AAnF o padrão do gerador OpenAPI (achado 14). *"User-Agent / IETF RFC 9110 [11] / This header shall be mainly used to identify the NF type of the HTTP/2 …"* |
| §5.2.2.2 | Mandatory to support HTTP standard headers | ¶3 (1 S) | não se aplica | `Via`: SCP e SEPP |
| §5.2.2.2 | Mandatory to support HTTP standard headers | ¶4 (1 S) | feito | `Authorization`, `Content-Type`, `Allow`, `WWW-Authenticate` (F12) |
| §5.2.2.2 | Mandatory to support HTTP standard headers | ¶5 (1 S) | **em parte** | o AAnF manda `User-Agent` com o tipo ao UDM; o AUSF, o NEF e o AF mandam ao AAnF o padrão do gerador OpenAPI (achado 14). *"NOTE 1: The value of NF type in the User-Agent header shall comply with the enumeration value (case …"* |
| §5.2.2.2 | Mandatory to support HTTP standard headers | ¶6 a ¶7 (2 S) | feito | `Authorization`, `Content-Type`, `Allow`, `WWW-Authenticate` (F12) |
| §5.2.2.2 | Mandatory to support HTTP standard headers | ¶8 (1 S) | não se aplica | `Via`: SCP e SEPP |
| §5.2.2.2 | Mandatory to support HTTP standard headers | ¶9 a ¶10 (2 S) | feito | `Authorization`, `Content-Type`, `Allow`, `WWW-Authenticate` (F12) |
| §5.2.2.2 | Mandatory to support HTTP standard headers | ¶11 a ¶12 (1 S, 1 R) | **não feito** | sem `Server` nas respostas de erro do AAnF e do CCF (*should*; achado 14). *"Server / IETF RFC 9110 [11] / This header should be inserted by the originator of an HTTP error response (see …"* |
| §5.2.3.1 | General | ¶1 a ¶9 (8 S, 1 R) | não se aplica | sintaxe dos cabeçalhos |
| §5.2.3.2.1 | General | ¶1 a ¶7 (5 S, 2 R) | **em parte** | achado 14: nenhum cabeçalho `3gpp-Sbi-*` tratado; a maioria só vale com SCP, SEPP ou *binding*. *"The 3GPP NF Services shall support the HTTP custom headers specified in Table 5.2.3.2.1-1 below. A …"* |
| §5.2.3.2.3 | 3gpp-Sbi-Callback | ¶1 (0 S, 1 R) | não se aplica | SCP, SEPP, *binding*, sobrecarga e CCA |
| §5.2.3.2.5 | 3gpp-Sbi-Routing-Binding | ¶1 a ¶6 (6 S) | não se aplica | SCP, SEPP, *binding*, sobrecarga e CCA |
| §5.2.3.2.6 | 3gpp-Sbi-Binding | ¶1 a ¶7 (7 S) | não se aplica | SCP, SEPP, *binding*, sobrecarga e CCA |
| §5.2.3.2.7 | 3gpp-Sbi-Discovery | ¶1 a ¶5 (5 S) | não se aplica | SCP, SEPP, *binding*, sobrecarga e CCA |
| §5.2.3.2.9 | 3gpp-Sbi-Oci | ¶1 a ¶6 (6 S) | não se aplica | SCP, SEPP, *binding*, sobrecarga e CCA |
| §5.2.3.2.10 | 3gpp-Sbi-Lci | ¶1 a ¶5 (5 S) | não se aplica | SCP, SEPP, *binding*, sobrecarga e CCA |
| §5.2.3.2.11 | 3gpp-Sbi-Client-Credentials | ¶1 a ¶7 (7 S) | não se aplica | SCP, SEPP, *binding*, sobrecarga e CCA |
| §5.2.3.2.12 | 3gpp-Sbi-Nrf-Uri | ¶1 a ¶3 (3 S) | não se aplica | SCP, SEPP, *binding*, sobrecarga e CCA |
| §5.2.3.2.15 | 3gpp-Sbi-Originating-Network-Id | ¶1 a ¶2 (2 S) | não se aplica | SCP, SEPP, *binding*, sobrecarga e CCA |
| §5.2.3.2.16 | 3gpp-Sbi-Access-Scope | ¶1 (1 S) | não se aplica | SCP, SEPP, *binding*, sobrecarga e CCA |
| §5.2.3.2.20 | 3gpp-Sbi-Nrf-Uri-Callback | ¶1 (1 S) | não se aplica | SCP, SEPP, *binding*, sobrecarga e CCA |
| §5.2.3.2.21 | 3gpp-Sbi-NF-Peer-Info | ¶1 a ¶9 (9 S) | não se aplica | SCP, SEPP, *binding*, sobrecarga e CCA |
| §5.2.3.2.22 | 3gpp-Sbi-Source-NF-Client-Credentials | ¶1 (1 S) | não se aplica | SCP, SEPP, *binding*, sobrecarga e CCA |
| §5.2.7.1 | General | ¶1 a ¶4 (4 S) | feito | `aanf.md` §6; `TestRoutingErrors` |
| §5.2.7.2 | NF as HTTP Server | ¶1 a ¶13 (11 S, 2 R) | feito | `aanf.md` §6; `TestRoutingErrors` |
| §5.2.7.3 | NF as HTTP Client | ¶1 a ¶5 (3 S, 2 R) | feito | `aanf.md` §6; `TestRoutingErrors` |
| §5.2.7.4 | SCP/SEPP | ¶1 a ¶5 (3 S, 2 R) | não se aplica | SCP e SEPP |
| §5.3 | Transport Protocol | ¶1 a ¶2 (2 S) | feito | TCP |
| §6.6.2 | Feature negotiation | ¶1 a ¶10 (9 S, 1 R) | feito | *features* (`TestRegisterFeatureNegotiation`) |
| §6.6.3 | Vendor-specific extensions | ¶1 a ¶2 (1 S, 1 R) | não se aplica | extensões de fornecedor |
| §6.6.4 | Extensibility for Query parameters | ¶1 a ¶4 (2 S, 2 R) | não se aplica | extensões de fornecedor |
| §6.7.1 | General | ¶1 (1 S) | não se aplica | SEPP e CCA |
| §6.7.2 | Transport layer security protection of messages | ¶1 (1 S) | feito | `openapi/0006` ([S11], [S12]) |
| §6.7.3 | Authorization of NF service access | ¶1 a ¶16 (16 S) | feito | `openapi/0006` ([S11], [S12]) |
| §6.7.4.1 | General | ¶1 (1 S) | não se aplica | SEPP e CCA |
| §6.7.4.2 | N32 Procedures | ¶1 a ¶2 (2 S) | não se aplica | SEPP e CCA |
| §6.7.5.1 | General | ¶1 a ¶3 (3 S) | não se aplica | SEPP e CCA |
| §6.7.5.2 | Authorization of NF Service Consumers for data a | ¶1 a ¶9 (9 S) | não se aplica | SEPP e CCA |
| §6.7.5.3 | Authorization of requesting ML models on behalf  | ¶1 a ¶2 (2 S) | não se aplica | SEPP e CCA |

### 5.3 TS 29.501 V19.5.0 (§4.4, §4.6.2.3, §4.8 e §4.10)

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §4.4.1 | Resource URI structure | ¶1 a ¶6 (4 S, 2 R) | feito | URIs com o FQDN (MVP-11); operações customizadas; `ProblemDetails` |
| §4.4.2 | Custom operations URI structure | ¶1 a ¶2 (2 S) | feito | URIs com o FQDN (MVP-11); operações customizadas; `ProblemDetails` |
| §4.4.3 | Callback URI structure | ¶1 (1 S) | feito | URIs com o FQDN (MVP-11); operações customizadas; `ProblemDetails` |
| §4.6.2.3 | Notifications | ¶1 a ¶6 (4 S, 2 R) | feito | a notificação do UDM ao AMF (MVP-12); o AAnF não notifica (DV-2) |
| §4.8.1 | Error Response Status Code | ¶1 (1 S) | feito | URIs com o FQDN (MVP-11); operações customizadas; `ProblemDetails` |
| §4.8.2 | Error Response Body | ¶1 a ¶6 (1 S, 5 R) | feito | URIs com o FQDN (MVP-11); operações customizadas; `ProblemDetails` |
| §4.8.3 | Extending ProblemDetails for API Backward Compat | ¶1 a ¶2 (1 S, 1 R) | feito | URIs com o FQDN (MVP-11); operações customizadas; `ProblemDetails` |
| §4.10 | Scopes definition for OAuth2.0 access token | ¶1 a ¶6 (2 S, 4 R) | feito | escopos de TS 29.535 tab. 5.1.9-1 (F4) |

### 5.4 TS 29.122 V19.7.0 (§5.2, que TS 29.522 §5.2 e TS 29.222 §7 aplicam)

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §5.2.2.1 | General | ¶1 (1 S) | **em parte** | HTTP/1.1 e HTTP/2 sobre TLS; o HTTP/2 por ALPN (achado 4, DV-15). *"For T8 APIs, support of HTTP/1.1 (IETF RFC 9112 [16], IETF RFC 9110 [17], IETF RFC 9111 [20]) over TLS is …"* |
| §5.2.2.2 | Usage of the HTTP PATCH method | ¶1 a ¶2 (2 S) | não se aplica | PATCH |
| §5.2.3 | Content type | ¶1 a ¶2 (2 S) | feito | JSON, URIs, erros e *features* do NEF |
| §5.2.3 | Content type | ¶3 a ¶5 (3 S) | não se aplica | PATCH |
| §5.2.3 | Content type | ¶6 (1 S) | feito | JSON, URIs, erros e *features* do NEF |
| §5.2.4.1 | Resource URI structure | ¶1 a ¶3 (3 S) | feito | JSON, URIs, erros e *features* do NEF |
| §5.2.4.2 | Custom operations URI structure | ¶1 a ¶2 (2 S) | feito | JSON, URIs, erros e *features* do NEF |
| §5.2.4.3 | Callback URI structure | ¶1 (1 S) | feito | JSON, URIs, erros e *features* do NEF |
| §5.2.5.1 | General | ¶1 (1 S) | **não feito** | sem notificações do NEF ao AF (DV-2). *"The SCEF and SCS/AS shall support the delivery of Notifications using a separate HTTP connection towards an …"* |
| §5.2.5.2 | Notification Delivery using a separate HTTP conn | ¶1 a ¶3 (2 S, 1 R) | **não feito** | sem notificações do NEF ao AF (DV-2). *"If a delivery of notifications is required for an API, the SCS/AS shall provide a notification destination …"* |
| §5.2.5.3 | Notification Test Event | ¶1 a ¶2 (2 S) | não se aplica | *features* opcionais das notificações |
| §5.2.5.4 | Notification Delivery using Websocket | ¶1 a ¶13 (11 S, 2 R) | não se aplica | *features* opcionais das notificações |
| §5.2.6 | Error handling | ¶1 a ¶3 (3 S) | feito | JSON, URIs, erros e *features* do NEF |
| §5.2.7 | Feature negotiation | ¶1 (1 S) | feito | JSON, URIs, erros e *features* do NEF |
| §5.2.8.3.4 | Nbi-Api-Lci | ¶1 (1 S) | não se aplica | carga e sobrecarga, opcionais |
| §5.2.8.3.5 | Nb-Api-Oci | ¶1 (1 S) | não se aplica | carga e sobrecarga, opcionais |

### 5.5 TS 23.003 V20.1.0 (§2.2B, §28.2, §28.3.2.3 e §28.7.2)

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §2.2B | Subscription Concealed Identifier (SUCI) | ¶1 (1 S) | **desvio** | SUPI de outros tipos (DV-3). *"When the SUPI type is a Network Specific Identifier (NSI), a GLI or a GCI, the Home Network Identifier …"* |
| §2.2B | Subscription Concealed Identifier (SUCI) | ¶2 (1 S) | **em parte** | cada dígito conta (`akid.RID`); o valor padrão `0` fica sem AAnF (achado 1). *"Each decimal digit present in the Routing Indicator shall be regarded as meaningful (e.g. value "012" is not …"* |
| §2.2B | Subscription Concealed Identifier (SUCI) | ¶3 a ¶4 (2 S) | free5GC | SUCI: free5GC e UERANSIM |
| §28.2 | Home Network Domain | ¶1 a ¶2 (2 S) | feito | `shared/akma/akid`: o MNC com três dígitos |
| §28.2 | Home Network Domain | ¶3 a ¶4 (2 S) | não se aplica | SNPN |
| §28.3.2.3.1 | General | ¶1 (1 S) | não se aplica | NRF de outra PLMN; SNPN |
| §28.3.2.3.2 | Format of NRF FQDN | ¶1 (1 S) | feito | `nrf.5gc.mnc093.mcc208.3gppnetwork.org` (MVP-11) |
| §28.3.2.3.2 | Format of NRF FQDN | ¶2 (1 S) | não se aplica | NRF de outra PLMN; SNPN |
| §28.3.2.3.3 | NRF URI | ¶1 a ¶6 (6 S) | não se aplica | NRF de outra PLMN; SNPN |
| §28.7.2 | NAI format for SUPI | ¶1 a ¶2 (2 S) | **desvio** | DV-3. *"The NAI for SUPI shall have the form username@realm as specified in clause 2.2 of IETF RFC 7542 [126]."* |

### 5.6 TS 23.502 V19.9.0 (§4.15.3, §4.17.4 e §5.2.3)

A V20.3.0 não foi comparada: o arquivo passa do limite de 10 MB da ferramenta de download. As enumerações do AKMA foram conferidas nela em 28/10/2026 (registro de defeitos).

| Cláusula | Assunto | Parágrafos | Estado | Evidência ou motivo |
|---|---|---|---|---|
| §4.15.3.1 | Monitoring Events | ¶1 a ¶2 (1 S, 1 R) | não se aplica | eventos e dados de outras NFs |
| §4.15.3.2.1 | AMF service operations information flow | ¶1 (1 S) | não se aplica | eventos e dados de outras NFs |
| §4.17.4 | NF/NF service discovery by NF service consumer i | ¶1 a ¶3 (3 S) | free5GC | NRF do free5GC |
| §5.2.3.3.1 | General | ¶1 a ¶9 (9 S) | não se aplica | eventos e dados de outras NFs |

### 5.7 Outras

| Norma e cláusula | Requisito | Estado | Evidência ou motivo |
|---|---|---|---|
| TS 29.505 V19.6.0 tab. 5.4.2.2-1 | `akmaAllowed`; `routingId` *"shall be present if the akmaAllowed is set to "true""* | feito | `provision-akma-subscriber.sh` (MVP-2); a *webconsole* apaga os campos (F9) |
| TS 29.571 V19.7.0 tab. 5.3.2-1, `Supi` | IMSI, NAI, GCI ou GLI | **desvio** | só IMSI (DV-3) |
| TS 29.571 V19.7.0 tab. 5.3.2-1, `NfInstanceId` | UUID versão 4, *"should be formatted as lower-case"* | feito | `sba-entities.conf` |
| TS 33.501 V20.3.0 §14.3.3 | `Nnrf_AccessToken_RetrieveKey` | **não feito** | achado 5 (F28) |
| TS 33.501 V20.3.0 §13.4.1.0 NOTA 1b | a chave do NRF no produtor, *"for example by OAM, by using the service operation specified in clause 14.3.3 or in any other implementation specific way"* | feito | `nrfCertPem` na configuração |

## 6. Achados da segunda passada

A segunda passada confirma os achados 1 a 10 da primeira e acrescenta os que seguem, decididos em 08/11/2026 como propostos (plano, MVP-14, Q16 a Q20). Os números continuam os da primeira passada.

11. `[3GPP]` TS 33.535 §4.4.2: *"AKMA AF shall be able to identify the AAnF serving the UE from the A-KID."* §6.7: a seleção *"should consider"* o Routing Indicator. §6.3 passo 1: *"the AF discovers the HPLMN of the UE based on the A-KID"*. O AF de dentro só usa o AAnF de `aanfUri` (`af/af/internal/config`), e o AF de fora manda todo A-KID ao CCF e ao NEF configurados, sem olhar o domínio. §6.7 permite a configuração local, mas o AF não consegue servir a mais de um AAnF nem a mais de uma HPLMN. Correção: sem `aanfUri`, o AF de dentro descobre o AAnF no NRF com o RID do A-KID, como o AUSF e o NEF; o AF de fora escolhe o CCF pelo domínio do A-KID, numa lista da configuração, e recusa o domínio que não conhece. Custo pequeno.
12. `[3GPP]` TS 33.222 §5.3.1.3: *"The AP/AS certificate profile shall comply with the requirements for TLS certificates in clause 6.1 of TS 33.310"*. O certificado que `af/scripts/make-test-ca.sh` faz para o AF do K.2 tem só o CN, sem o O que TS 33.310 §6.1.1 pede, e marca `basicConstraints` como crítica, que §6.1.1 manda deixar não crítica quando a norma não a exige. A CA de teste também não tem O. Correção: o perfil de `make-capif-certs.sh`, que segue §6.1.3a e §6.1.4a. Custo pequeno.
13. `[3GPP]` TS 29.510 §5.3.2.2.1: *"Before a service consumer invokes this service operation, it shall consider if it is possible to reuse the results from a previous searching"*; *"should reuse the previous result"*. O AUSF e o NEF descobrem o AAnF no NRF a cada pedido. Correção: guardar o resultado durante a `validityPeriod` da resposta. Custo pequeno.
14. `[3GPP]` TS 29.500 tab. 5.2.2.2-1: o `User-Agent` identifica o tipo de NF do cliente e *"should be included in every HTTP/2 request"*; o `Server` *"should be inserted by the originator of an HTTP error response"*. Tab. 5.2.3.2.1-1: *"The 3GPP NF Services shall support the HTTP custom headers specified in Table 5.2.3.2.1-1"*, o que TS 29.535 §5.1.2.3 repete para o AAnF. O AUSF, o NEF e o AF mandam ao AAnF o `User-Agent` padrão do gerador OpenAPI (`OpenAPI-Generator/1.0.0/go`), sem o tipo de NF; só o AAnF, no pedido ao UDM, manda o tipo (`AANF-<id>`). O AAnF e o CCF não mandam `Server`, e nenhum `3gpp-Sbi-*` é tratado. A maioria desses cabeçalhos só tem papel com SCP, SEPP ou *binding*; `3gpp-Sbi-Message-Priority` e `3gpp-Sbi-Callback` (nas notificações) têm papel sem eles. Correção: o `User-Agent` com o tipo de NF nos pedidos do AUSF, do NEF e do AF ao AAnF, e o `Server` nas respostas de erro do AAnF e do CCF (pequeno); os `3gpp-Sbi-*`, com o *upstream*, porque valem para todas as NFs do free5GC. Registro proposto para a segunda parte: F31.
15. `[3GPP]` TS 33.310 §7.1: *"SEGs shall use LDAP to access the CRL and cross-certificate repositories. TLS entities shall use LDAP or HTTP to access the CRL repositories."* A frase reforça a leitura (b) do T27 do MVP-11: para as entidades TLS, o HTTP basta. O *"shall be supported as the primary method"* de §6.1a continua, e o Q10, aprovado, põe o LDAPv3 no MVP-18; o DV-18 ganha a citação.

## 7. O que a segunda passada refez da primeira

- **TS 33.310 na tabela das normas do plano.** A primeira passada dava §6.1.1 como coberto. O certificado do AF do K.2 não segue o perfil (achado 12): a linha passa a **em parte**.
- **Motivo do DV-4.** A primeira passada dizia que o FQDN do AF vai no `server_name`, que só aceita ASCII. Vale para o K.2, o PSK TLS e o DTLS, que mandam o SNI. No OSCORE, o FQDN vai no URI do CoAP, e a RFC 3986 §3.2.2 admite nome com *percent-encoding* de UTF-8, ao lado dos A-labels. O motivo fica: na prática, o FQDN do AF é um nome DNS em ASCII; o NFKC só pesaria para um AF do perfil OSCORE com nome não ASCII.
- **DV-18.** Ganha a citação de §7.1 (achado 15), sem mudar a decisão do Q10.
- **Ordem dos MVPs.** Com o Q5 na opção (c), o *roaming* sai da fila: o MVP-16 fica reservado para quando houver uma VPLMN. Com o Q7, o MVP-17 traz também os *logs* e a auditoria.

## 8. Registros

Feitos nesta data, com o Q11 aprovado:

- `desvios.md`: a decisão de cada desvio na revisão do MVP-14; o DV-15 ampliado ao NEF (achado 4); o DV-16 ampliado à SBI, a CAPIF-2e e a CAPIF-3/4 (achado 3); o DV-18 com §7.1 (achado 15).
- `correcoes-free5gc.md`: F28 (`Nnrf_AccessToken_RetrieveKey`), F29 (*heart-beat*) e F30 (`ROAMING_STATUS` no UDM).
- Registro de defeitos, e a cópia em `~/akma`: D-AL a D-AO e o D-AE ampliado a TS 29.122 §5.2.2.
- Matriz de rastreabilidade, revisão 12 (Q14).

Decididos depois, na mesma data, com as recomendações aceitas: o Q3 (opção a, sem servidor de autorização para o modo sem CAPIF), o Q6 (opção b no DV-3, no DV-4 e no DV-17; opção a na `D0 02`), o Q13 (opção a) e os achados 11 a 15 (Q16 a Q20). Registros que vieram dessas decisões:

- `desvios.md`: a decisão do DV-3, do DV-4, do DV-9, do DV-10 e do DV-17;
- `correcoes-free5gc.md`: o F31 (os cabeçalhos de TS 29.500) e a conferência de F21 a F31 no free5GC v4.3.0, em que os onze continuam presentes.

O registro do DV-19 não é preciso: pelo Q8, o NEF deixa de repassar o `anonInd` no MVP-15.
