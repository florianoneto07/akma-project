# Registro do AAnF no NRF do free5GC

**MVP:** MVP-1 · **Data:** 17 de outubro de 2026 · **Verificado na VM do free5GC** (192.168.16.130) em 17/10/2026

## 1. Nenhum arquivo do NRF é alterado

O AAnF usa o NRF do free5GC v4.2.3 (`nrf` v1.4.5) como ele é. A configuração do NRF (`~/free5gc/config/nrfcfg.yaml`) não muda. O NRF, porém, grava estado próprio quando o AAnF se registra; isso está listado na seção 4.

## 2. Perfil publicado

`[3GPP]` TS 29.510 `NFProfile` e `NFService`. Implementação em `free5gc/aanf/internal/sbi/consumer/nrf_service.go` (`BuildNfProfile`) e `internal/context/context.go` (`nfService`).

| Atributo | Valor | Origem |
|---|---|---|
| `nfInstanceId` | `configuration.nfInstanceId` (fixo no `aanfcfg.yaml` do projeto) | `[MVP]` fixo para que o NRF não acumule certificados, ver seção 4 |
| `nfType` | `AANF` | `[3GPP]` TS 29.510 Tab. 6.1.6.3.3-1 |
| `nfStatus` | `REGISTERED` | |
| `ipv4Addresses` | `sbi.registerIPv4` (`192.168.16.130`) | ADR-0001-R2 §5.2: endereço alcançável a partir da VM do AF |
| `plmnList` | `plmnSupportList` | |
| `nfServices[0].serviceName` | `naanf-akma` | `[3GPP]` TS 29.535 §5.1.1 |
| `versions` | `apiVersionInUri: v1`, `apiFullVersion: 1.2.0` | `[3GPP]` versão do `TS29535_Naanf_AKMA.yaml` implementado. `[FREE5GC]` As NFs do free5GC publicam a versão do arquivo de configuração; o AAnF publica a da API |
| `ipEndPoints` | `192.168.16.130`, TCP, porta 8000 | |
| `allowedNfTypes` | `AUSF`, `NEF`, `AF` | `[3GPP]` consumidores do Naanf_AKMA, TS 29.535 Tab. 4.2.2.1-1. `[FREE5GC]` O NRF só emite token para esses tipos |
| `aanfInfoList` | `{"1": {"routingIndicators": ["0000"]}}`, omitido sem RIDs configurados | `[3GPP]` TS 29.510 Tab. 6.1.6.2.73-1: RIDs como cadeias; sem `AanfInfo`, o AAnF atende qualquer RID. Chave do mapa com até 32 caracteres |

`[FREE5GC]` Sem *heartbeat*: o NRF v1.4.5 aceita `heartBeatTimer`, mas não expira perfis, e as NFs do free5GC (o AUSF, por exemplo) não enviam `NFUpdate` periódico. O AAnF segue o mesmo padrão. Ele se descadastra (`NFDeregister`) ao receber `SIGTERM` ou `SIGINT`.

## 3. OAuth2

`[FREE5GC]` O NRF devolve no perfil, em `customInfo.oauth2`, se o OAuth2 está ligado (`nrfcfg.yaml`: `oauth: true` na VM). O AAnF lê esse valor na resposta 201 (criação) e também na 200 (atualização). O AUSF v1.4.5 só o lê na 201: se reiniciar com o mesmo ID antes de o NRF limpar o perfil, deixa de exigir token. Para verificar os tokens, o AAnF usa o certificado do NRF configurado em `nrfCertPem` (`cert/nrf.pem`). A verificação completa está em `aanf.md` §5.

## 4. Estado que o NRF grava

| Onde | O quê | Quando some |
|---|---|---|
| MongoDB, base `free5gc`, coleção `NfProfile` | o perfil do AAnF | no `NFDeregister`; o NRF descarta a coleção inteira ao encerrar |
| `~/free5gc/cert/aanf.pem` | chave pública RSA gerada pelo NRF na primeira vez que vê o tipo `AANF` | permanece (arquivo novo no diretório do free5GC; `git status` o mostra como não rastreado) |
| `~/free5gc/cert/aanf_<nfInstanceId>.pem` | certificado assinado pelo NRF para a instância | apagado pelo NRF no `NFDeregister` |

Com `nfInstanceId` aleatório, cada partida do AAnF criaria um certificado novo. O ID fixo deixa um único arquivo. A ferramenta `seed`, quando registra um AF temporário para obter token, também usa ID fixo. O NRF cria e apaga `cert/af_<id>.pem`; o `cert/af.pem` já existia na instalação.

## 5. Limitações do NRF v1.4.5 encontradas

Nos três casos o AAnF segue a norma; quem não segue é o NRF. Em 17/10/2026, o responsável pelo projeto decidiu **corrigir o NRF com um *patch*, mais adiante**. Até lá:

| Nº | O que o NRF faz | Onde | Efeito |
|---|---|---|---|
| N1 | Descarta `aanfInfoList` ao gravar o perfil: só copia os `*Info` que conhece (`udrInfo`, `udmInfo`, `ausfInfo`…) | `internal/context/management_data.go`, `nnrfNFManagementOption` | o perfil guardado não tem os RIDs; o critério 1 do MVP-1 fica verificado no que o AAnF envia (teste `TestRegisterNFInstanceProfile`), não no que o NRF guarda |
| N2 | A descoberta recusa `target-nf-type=AANF` com 400 (`Loss mandatory parameter`): a lista de tipos aceitos é fixa e não tem `AANF` | `internal/sbi/processor/nf_discovery.go`, `validateQueryParameters` | AUSF e AF não descobrem o AAnF pelo NRF; usam o endereço configurado, permitido por TS 33.535 §6.7 (*"unless AAnF information is available by other means, e.g. locally configured"*) |
| N3 | Emite tokens com `aud` vazio: copia `targetNfInstanceId`, que os clientes do free5GC não enviam | `internal/sbi/processor/access_token.go` | o AAnF tolera `aud` vazio (desvio DV-8) |

## 6. Como reproduzir numa instalação limpa

Nada a fazer no NRF além de iniciá-lo: o AAnF se registra sozinho ao partir (`aanf.md` §8). O `nrfCertPem` do `aanfcfg.yaml` deve apontar para o mesmo `cert/nrf.pem` usado pelas outras NFs.

## 7. Como verificar

Com NRF e AAnF em execução, na máquina do free5GC:

```bash
mongosh --quiet free5gc --eval 'printjson(db.NfProfile.find({nfType: "AANF"}, {_id: 0}).toArray())'
```

O resultado deve conter `nfType: 'AANF'`, o serviço `naanf-akma` com `apiVersionInUri: 'v1'` e `ipEndPoints` em `192.168.16.130:8000`; `aanfInfoList` não aparece (N1). No log do AAnF (`~/free5gc/log/aanf/aanf.log`) deve aparecer `OAuth2 setting receive from NRF: true`. Ao parar o AAnF, aparecem `Deregister from NRF successfully` no log e um `NfProfile` sem o AAnF.

Resultado em 17/10/2026: perfil conforme a tabela da seção 2, sem `aanfInfoList`. `customInfo.oauth2: true`. `NFDeregister` aceito (204), e `cert/aanf_<id>.pem` apagado pelo NRF. A descoberta respondeu 400 (N2).
