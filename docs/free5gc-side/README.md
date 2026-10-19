# Lado free5GC

Tudo o que é feito **dentro** de uma instalação do free5GC para habilitar AKMA: o AAnF, os *patches* de AUSF, UDM e `free5gc/openapi`, configuração, registro no NRF e dados de assinatura.

Cada documento desta pasta registra:

1. qual arquivo do free5GC foi tocado (caminho completo na instalação);
2. o que mudou e por quê, com a cláusula 3GPP;
3. como reproduzir numa instalação limpa;
4. como verificar que funcionou.

Base de código: free5GC v4.2.3 (`docs/baseline/ambiente-verificado.md`).

| Documento | MVP | Estado |
|---|---|---|
| `aanf.md` — a NF AAnF: configuração, operações, erros, verificação | MVP-1 | escrito (17/10/2026) |
| `openapi.md` — *patches* com os modelos e o cliente do Naanf_AKMA | MVP-1 | escrito (17/10/2026) |
| `registro-nrf.md` — perfil no NRF, OAuth2, limitações do NRF v1.4.5 | MVP-1 | escrito (17/10/2026) |
| `nrf.md` — *patches* no NRF: `aanfInfoList`, descoberta de `AANF`, `aud` do token | MVP-1b | escrito (17/10/2026) |
| `udm.md` — indicação AKMA e Routing Indicator (D10) | MVP-2 | escrito (17/10/2026) |
| `ausf.md` — K_AKMA, A-KID, descoberta do AAnF e registro (D9) | MVP-2 | escrito (17/10/2026) |
| `assinatura-akma.md` — `akmaAllowed` e `routingId` no UDR | MVP-2 | escrito (17/10/2026) |
| `verificacao.md` — banco de testes, regressão, UERANSIM e instalação limpa | MVP-2, MVP-3a, MVP-5 | escrito (17/10/2026); seções 4 e 5 em 18/10/2026; seção 5.1 em 19/10/2026 |
| `logs-akma.md` — mensagens AKMA na tela do `run.sh`; opção `-aanf` | MVP-3a | escrito (18/10/2026) |
| `plano-de-usuario.md` — do UE ao AF pelo N6: diagnóstico e NAT | MVP-4 | escrito (18/10/2026) |

O que é comum aos dois lados — KDF, formatos de identificador, codificação, vetores — fica em `docs/shared/`.
