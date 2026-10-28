# Recomendações ao UERANSIM

**Criado em:** 28 de outubro de 2026 · **MVP:** MVP-13 (ponto L13, pedido do responsável) · **Base:** UERANSIM *commit* `85a0fbfd` (v3.2.6-66, 17/04/2024), o dos *patches* do projeto

Registro do que o projeto achou no UERANSIM contra as normas, para a lista pública de recomendações (`public/docs/ueransim-recommendations.md`). Cada item diz o que o UERANSIM faz, onde a norma diz o que fazer, como fazer e o que o projeto faz.

**Conferência na versão mais recente (28/10/2026).** Cada item foi conferido na v3.3.0 (07/06/2026) e no `master` (`4f45e18f`, 14/10/2026), lidos no GitHub (`aligungr/UERANSIM`). Os quatro continuam nas duas. Os *patches* do projeto aplicam no `85a0fbfd`; a v3.3.0 mudou `src/ue/nas/keys.cpp` e `src/ue/nas/mm/auth.cpp`, e portá-los fica para depois.

| Nº | Assunto | Norma | No projeto |
|---|---|---|---|
| UE-1 | AKMA no ME: a K_AKMA do UERANSIM está errada e o resto falta | TS 33.535 §4.4.3, §6.1, §6.2.1, Anexo A | *patches* `0001` a `0005` |
| UE-2 | Routing Indicator `"0000"` quando não configurado | TS 23.003 §2.2B; TS 24.501 §9.11.3.4 | configurações com o RID explícito |
| UE-3 | `Security Mode Complete` com o último `Registration Request` fora de um registro | TS 24.501 §4.4.6, §5.4.2.3 | o AMF ignora o conteúdo (MVP-12) |
| UE-4 | Sem segurança AS | TS 33.501 §5.2.2, §5.2.3, §6.9.4.4 | limitação registrada |

---

## UE-1 — AKMA no ME

- **O que o UERANSIM faz.** `[FREE5GC]` Desde abril de 2025 (*commits* `89f02d7c` e `27b496ce`, *"Fixes #766"*), o `src/ue/nas/keys.cpp` tem `CalculateAkmaKey`, com `KAKMA_derivation_function_code = 0x6A` e um só parâmetro, o SUPI. `src/ue/nas/mm/auth.cpp` a chama só no 5G AKA, não no EAP-AKA'. Não há A-TID, A-KID, K_AF nem interface para as aplicações. O `0x6A` é o FC da K_AUSF (TS 33.501 A.2).
- **Onde a norma diz.** `[3GPP]` TS 33.535 V19.0.0:
  - A.2: *"FC = 0x80; P0 = "AKMA"; L0 = length of "AKMA"; (i.e. 0x00 0x04); P1 = SUPI; L1 = length of SUPI"*, com a K_AUSF como chave e o SUPI na forma de TS 33.501 A.7.0;
  - A.3: o A-TID, com FC = 0x81 e P0 = `"A-TID"`; A.4: a K_AF, com FC = 0x82 e P0 = AF_ID; §6.1: o A-KID, em formato NAI, com o RID e o A-TID no *username* e o identificador da rede doméstica no *realm*;
  - §6.2.1 passo 1: *"The UE shall generate the AKMA Anchor Key (KAKMA) and the A-KID from the KAUSF before initiating communication with an AKMA Application Function"*;
  - §4.4.3: *"Applications on the UE shall not be able to get access to KAKMA"*; *"An application on the UE shall only get the KAF keys related to specific AF Identifiers (AF_IDs) that the application is authorized to get"*.
- **Como deve ser.** K_AKMA com FC = 0x80, P0 = `"AKMA"` e P1 = SUPI, nos dois métodos de autenticação; A-TID e A-KID; K_AF por AF_ID; entrega da K_AF só às aplicações autorizadas para aquele AF_ID.
- **No projeto.** *Patches* `af/ue/ueransim/patches/0001` a `0005` (GPL-3.0), para o `85a0fbfd`: K_AKMA, A-TID e A-KID quando o *Security Mode Command* põe o contexto novo em uso, nos dois métodos; K_AF sob demanda; serviço de chaves por aplicação (`ue-ueransim.md` §3); `nr-akma-selftest` contra os vetores do projeto. Verificados ponta a ponta com o free5GC (`docs/af-side/e2e.md`).

## UE-2 — Routing Indicator quando não configurado

- **O que o UERANSIM faz.** `[FREE5GC]` Sem `routingIndicator` na configuração do UE, `src/ue/nas/mm/identity.cpp` (linha 105 da v3.3.0 e do `master`) usa `"0000"`, quatro dígitos.
- **Onde a norma diz.** `[3GPP]`
  - TS 23.003 V20.1.0 §2.2B (o mesmo texto na V19.7.0): *"Each decimal digit present in the Routing Indicator shall be regarded as meaningful (e.g. value "012" is not the same as value "12"). If no Routing Indicator is configured on the USIM or the ME, this data field shall be set to the value 0 (i.e. only consist of one decimal digit of "0")."*
  - TS 24.501 V20.1.0 §9.11.3.4: *"If no Routing Indicator is configured in the USIM or the ME, the UE shall code bits 1 to 4 of octet 8 of the Routing Indicator as "0000" and the remaining digits as "1111"."*
- **Como deve ser.** O valor `"0"`, um dígito, codificado como `0`, `F`, `F`, `F`.
- **Impacto no AKMA.** O RID entra no A-KID (`rid0000…` contra `rid0…`) e escolhe o AAnF (TS 33.535 §6.7). Um UE sem RID configurado e uma rede que siga a norma chegam a A-KIDs diferentes.
- **No projeto.** Sem mudança: as configurações do UE têm `routingIndicator: '0000'` explícito.

## UE-3 — `Security Mode Complete` com o último `Registration Request`

- **O que o UERANSIM faz.** `[FREE5GC]` `src/ue/nas/mm/security.cpp` (linha 264 da v3.3.0 e do `master`): quando o último registro foi mandado sem contexto de segurança NAS, a marca `m_lastRegWithoutNsc` fica ligada depois do registro, e todo `Security Mode Complete` leva aquele `Registration Request` no `NAS message container`. Isso inclui o SMC de uma autenticação pedida pela rede doméstica, muito depois do registro (MVP-12).
- **Onde a norma diz.** `[3GPP]` TS 24.501 V20.1.0 §5.4.2.3: *"If, prior to receiving the SECURITY MODE COMMAND message, the UE without a valid 5G NAS security context had sent a REGISTRATION REQUEST message the UE shall include the entire REGISTRATION REQUEST message in the NAS message container IE of the SECURITY MODE COMPLETE message as described in clause 4.4.6."* §4.4.6 a) trata da proteção da mensagem NAS inicial: o SMC que segue aquele registro.
- **Como deve ser.** `[INTERPRETAÇÃO]` O conteúdo vai só no SMC do registro em que o `Registration Request` foi mandado. Desligar `m_lastRegWithoutNsc` quando o registro termina. A frase de §5.4.2.3, ao pé da letra, não tem esse limite; é o defeito D-Z do registro de defeitos, a sugerir à 3GPP.
- **No projeto.** O AMF ignora o conteúdo durante a autenticação pedida pela rede (`free5gc/patches/amf/0001`, `docs/free5gc-side/amf.md`). O UE não muda.

## UE-4 — Sem segurança AS

- **O que o UERANSIM faz.** `[FREE5GC]` Nem o UE nem o gNB têm o procedimento de modo de segurança do RRC nem a K_gNB (`src/ue/rrc`, `src/gnb/rrc`). O RRC e o plano de usuário entre o UE e o gNB vão sem cifração e sem proteção de integridade.
- **Onde a norma diz.** `[3GPP]` TS 33.501 V20.3.0:
  - §5.2.2: *"The UE shall support ciphering of user data between the UE and the gNB."* e *"The UE shall support ciphering of RRC and NAS-signalling."*;
  - §5.2.3: *"The UE shall support integrity protection and replay protection of user data between the UE and the gNB."* e o mesmo para o RRC e o NAS;
  - §6.9.4.4: a troca das chaves AS depois de uma autenticação.
- **Como deve ser.** O SMC do RRC e a segurança do PDCP, com as chaves derivadas da K_gNB (TS 33.501 Anexo A.9).
- **Impacto no AKMA.** Depois de uma nova autenticação, a troca das chaves AS não pode ser testada (MVP-12, K7).
- **No projeto.** Limitação registrada, sem mudança.
