# SZChat API (Não oficial)

[![Go Reference](https://pkg.go.dev/badge/github.com/raykavin/szchat-go.svg)](https://pkg.go.dev/github.com/raykavin/szchat-go)
[![Go Version](https://img.shields.io/badge/go-1.26+-blue)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Uma biblioteca cliente em Go para a **API do SZChat Chat Center**.

Cobre os módulos principais: Contatos, Grupos de Contatos, Canais, Agentes, Mensagens
entre Agentes, Equipes, Administradores, Mensagens, Configuração da Aplicação,
Tabulações e Pausas.


> **Aviso:** este é um projeto de código aberto independente e não possui
> qualquer vínculo, patrocínio ou endosso por parte da [Fortics](https://fortics.com.br/),
> fornecedora oficial do SZChat. "SZChat" é uma marca de sua respectiva
> proprietária, citada aqui apenas para fins de identificação da API
> consumida por esta biblioteca.

---

## Instalação

```bash
go get github.com/raykavin/szchat-go
```

Requer **Go 1.26+**. Única dependência externa: [github.com/raykavin/gobox](https://github.com/raykavin/gobox).

---

## Início Rápido

```go
import szchat "github.com/raykavin/szchat-go"

// baseURL é obrigatório e deve ser apenas a URL do seu tenant SZChat, sem o
// path de versão — o "/api/v4" é adicionado automaticamente pela biblioteca.
client, err := szchat.NewClient("https://seu-tenant.sz.chat", "agente@example.com", "sua-senha")
if err != nil {
    log.Fatal(err)
}

ctx := context.Background()

// Listar contatos (primeira página)
page, err := client.ContactAPI.List(ctx, szchat.ContactListFilter{})
if err != nil {
    log.Fatal(err)
}

for _, c := range page.Data {
    fmt.Println(c.ID, c.Name, c.Email)
}
```

---

## Configuração

A URL base **não possui valor padrão** e deve ser sempre informada explicitamente
como o primeiro argumento de `NewClient`. Ela é específica de cada tenant e **não**
deve incluir o path de versão da API: a biblioteca anexa `/api/v4` automaticamente
para manter a compatibilidade interna com a versão da API suportada (ex: passe
`https://seu-tenant.sz.chat`, e a biblioteca chamará `https://seu-tenant.sz.chat/api/v4/...`).

```go
// URL base obrigatória — apenas o host do tenant, sem "/api/v4"
client, err := szchat.NewClient("https://seu-tenant.sz.chat", "agente@example.com", "sua-senha")

// Device token enviado no login, usado pelo SZChat para push notifications
client, err = szchat.NewClient("https://seu-tenant.sz.chat", "agente@example.com", "sua-senha",
    szchat.WithDeviceToken("meu-device-token"),
)

// Cliente HTTP customizado (TLS, proxies, timeouts customizados, etc.)
httpClient := &http.Client{Timeout: 60 * time.Second}
client, err = szchat.NewClient("https://seu-tenant.sz.chat", "agente@example.com", "sua-senha",
    szchat.WithHTTPClient(httpClient),
)

// Política de retry customizada para falhas transitórias (429/502/503/504 e erros de rede)
client, err = szchat.NewClient("https://seu-tenant.sz.chat", "agente@example.com", "sua-senha",
    szchat.WithRetry(5, 200*time.Millisecond, 2*time.Second),
)
```

---

## Autenticação

O Client se autentica de forma **lazy**: a primeira requisição dispara um login
usando as credenciais passadas para `NewClient`, e o token bearer é renovado
automaticamente sempre que uma requisição retorna não autorizada — nenhuma
configuração adicional é necessária.

```go
client, err := szchat.NewClient("https://seu-tenant.sz.chat", "agente@example.com", "sua-senha")

// Chamar Login explicitamente normalmente não é necessário, mas está disponível:
loginResp, err := client.Login(ctx)

// Obter o perfil do usuário autenticado
me, err := client.Me(ctx)

// Encerrar a sessão atual ("api", "web" ou "all")
err = client.Logout(ctx, "api")
```

---

## Módulos

### Contatos

```go
// Listar contatos (paginado)
page, err := client.ContactAPI.List(ctx, szchat.ContactListFilter{
    ListOptions: szchat.ListOptions{Limit: 20},
})

// Buscar com filtros
page, err = client.ContactAPI.Search(ctx, szchat.ContactListFilter{
    Name:     "Maria",
    Platform: "Whatsapp",
})

// Obter um contato
contact, err := client.ContactAPI.Get(ctx, contactID)

// Criar um contato (Extra mescla campos de canal específicos do tenant no payload)
contact, err = client.ContactAPI.Create(ctx, szchat.ContactRequest{
    Name:     "Maria Silva",
    Whatsapp: "5511999999999",
    Extra:    map[string]any{"Generic_MyBot": "external-id-123"},
})

// Atualizar um contato
contact, err = client.ContactAPI.Update(ctx, contactID, szchat.ContactRequest{
    Name: "Maria S. Silva",
})

// Atualizar apenas campos customizados
fields, err := client.ContactAPI.UpdateFields(ctx, contactID, map[string]any{
    "cpf": "123.456.789-00",
})

// Excluir um contato
err = client.ContactAPI.Delete(ctx, contactID)

// Salvar/atualizar uma anotação
annotation, err := client.ContactAPI.SaveAnnotation(ctx, szchat.ContactAnnotationRequest{
    ContactID:   contactID,
    Observation: "Retornar ligação na segunda-feira",
    AgentID:     agentID,
})

// Estatísticas de atendimento em uma plataforma
stats, err := client.ContactAPI.AttendanceStats(ctx, contactID, "Whatsapp")
```

---

### Grupos de Contatos

```go
// Listar grupos (paginado)
page, err := client.ContactGroupAPI.List(ctx, szchat.ListOptions{Limit: 20})

// Criar, renomear, excluir
group, err := client.ContactGroupAPI.Create(ctx, "Clientes VIP")
group, err = client.ContactGroupAPI.Update(ctx, group.ID, "Clientes VIP 2026")
err = client.ContactGroupAPI.Delete(ctx, group.ID)

// Grupos aos quais um contato específico pertence
groups, err := client.ContactGroupAPI.ListByContact(ctx, contactID)
```

---

### Canais

```go
// Listar canais configurados
channels, err := client.ChannelAPI.List(ctx)
for _, ch := range channels {
    fmt.Println(ch.ID, ch.Platform, ch.Number)
}

// Todas as plataformas disponíveis para novos canais
platforms, err := client.ChannelAPI.ListPlatforms(ctx)

// Apenas plataformas com pelo menos um canal ativo
active, err := client.ChannelAPI.ListActivePlatforms(ctx)
```

---

### Agentes

```go
// Listar agentes (paginado)
page, err := client.AgentAPI.List(ctx, szchat.AgentListFilter{
    ListOptions: szchat.ListOptions{Limit: 20},
    Name:        "João",
})

// Criar, obter, atualizar, excluir (Delete aceita múltiplos ids para exclusão em lote)
agent, err := client.AgentAPI.Create(ctx, szchat.AgentRequest{
    Name:     "João Souza",
    Email:    "joao@example.com",
    Password: "s3cr3t",
})
agent, err = client.AgentAPI.Get(ctx, agent.ID)
agent, err = client.AgentAPI.GetByEmail(ctx, "joao@example.com")
agent, err = client.AgentAPI.Update(ctx, agent.ID, szchat.AgentRequest{Name: "João S."})
err = client.AgentAPI.Delete(ctx, agent.ID)

// Status online entre agentes/administradores
online, err := client.AgentAPI.ListOnlineStatus(ctx, "Whatsapp", "")
online2, err := client.AgentAPI.ListOnlineAgents(ctx, true, "1")

// Dados do próprio agente autenticado
teams, err := client.AgentAPI.MyTeams(ctx)
attendances, err := client.AgentAPI.MyAttendances(ctx)
attendancesPlus, err := client.AgentAPI.MyAttendancesPlus(ctx)
grades, err := client.AgentAPI.MyGrades(ctx)

// Alternar a equipe ativa e atualizar o timestamp de última interação
msg, err := client.AgentAPI.ToggleTeam(ctx, teamID)
err = client.AgentAPI.UpdateLastInteraction(ctx)
```

---

### Mensagens entre Agentes

```go
// Enviar uma mensagem interna para outro agente
resp, err := client.AgentTalkAPI.SendMessage(ctx, szchat.AgentTalkMessageRequest{
    AgentTo: peerAgentID,
    Type:    "text",
    Message: "Pode assumir esse contato?",
})

// Listar todas as conversas
conversations, err := client.AgentTalkAPI.ListConversations(ctx)

// Histórico paginado com um agente específico
history, err := client.AgentTalkAPI.SearchConversation(ctx, peerAgentID, szchat.ListOptions{
    Limit: 30,
})
for _, msg := range history.Data {
    fmt.Println(msg.AgentFrom, "->", msg.AgentTo, ":", msg.Message)
}
```

---

### Equipes

```go
// Listar equipes (paginado)
page, err := client.TeamAPI.List(ctx, szchat.TeamListFilter{Limit: 20})

// Resumo leve {_id, name} de todas as equipes
summaries, err := client.TeamAPI.Resume(ctx, true)

// Filtrar equipes por ids
page, err = client.TeamAPI.FilterByIDs(ctx, szchat.TeamFilterByIDsRequest{
    CampaignIDs: []string{teamID1, teamID2},
})

// Criar, atualizar, excluir
team, err := client.TeamAPI.Create(ctx, szchat.Team{
    Name:           "Suporte",
    RuleAttendance: "sequential",
})
updateResp, err := client.TeamAPI.Update(ctx, team.ID, team)
err = client.TeamAPI.Delete(ctx, team.ID)
```

---

### Administradores

```go
// Listar administradores e grupos de permissão (ambos paginados)
page, err := client.AdminAPI.List(ctx, szchat.ListOptions{Limit: 20})
groups, err := client.AdminAPI.ListGroups(ctx, szchat.ListOptions{})

// Criar, obter, atualizar (parcial), excluir
admin, err := client.AdminAPI.Create(ctx, szchat.AdminRequest{
    Name:    "Usuário Admin",
    Email:   "admin@example.com",
    GroupID: groupID,
})
admin, err = client.AdminAPI.Get(ctx, admin.ID)
updateResp, err := client.AdminAPI.Update(ctx, admin.ID, map[string]any{"name": "Admin U."})
err = client.AdminAPI.Delete(ctx, admin.ID)
```

---

### Mensagens

```go
// Enviar uma mensagem
resp, err := client.MessageAPI.Send(ctx, szchat.SendMessageRequest{
    PlatformID: contactPlatformID,
    ChannelID:  channelID,
    Type:       "text",
    Message:    "Olá!",
})

// Enviar e criar/atualizar o contato de destino na mesma chamada
resp, err = client.MessageAPI.SendPlus(ctx, szchat.SendMessagePlusRequest{
    SendMessageRequest: szchat.SendMessageRequest{
        PlatformID: "5511999999999",
        ChannelID:  channelID,
        Type:       "text",
        Message:    "Seja bem-vindo!",
    },
    ContactVariables: &szchat.SendMessageContactVariables{
        Name: "Maria Silva",
    },
})

// Anotar/remover uma nota em uma mensagem
note, err := client.MessageAPI.AddAnnotation(ctx, messageID, map[string]any{"note": "Escalado"})
err = client.MessageAPI.RemoveAnnotation(ctx, messageID)
```

`MessageAPI.Read`, `MessageAPI.Pending` e `MessageAPI.ReplaceVars` envolvem
endpoints cujo formato de payload não é documentado publicamente pelo SZChat,
por isso aceitam e retornam `map[string]any` sem tipagem adicional.

---

### Configuração da Aplicação

```go
// Configuração geral do tenant
config, err := client.ApplicationAPI.Get(ctx)
updated, err := client.ApplicationAPI.Update(ctx, szchat.ApplicationUpdateRequest{
    Timezone: "America/Sao_Paulo",
    Language: "pt-BR",
})

// Faixa de notas de atendimento
grades, err := client.ApplicationAPI.GetAttendanceGrades(ctx)
updated, err = client.ApplicationAPI.UpdateAttendanceGrades(ctx, szchat.ApplicationAttendanceGrades{
    InitialGrade: 1,
    FinalGrade:   5,
})

// Mensagens automáticas (opção inválida, atendimento iniciado/concluído)
messages, err := client.ApplicationAPI.GetAutoMessages(ctx)
updated, err = client.ApplicationAPI.UpdateAutoMessages(ctx, szchat.ApplicationMessages{
    StartedService: &szchat.ApplicationAutoMessage{Enable: true, Message: "Já estamos atendendo você!"},
})
```

---

### Tabulações

```go
// Listar tabulações (paginado)
page, err := client.TabulationAPI.List(ctx, szchat.ListOptions{Limit: 20})

// Criar, renomear, excluir
tab, err := client.TabulationAPI.Create(ctx, "Resolvido")
tab, err = client.TabulationAPI.Update(ctx, tab.ID, "Resolvido - Acompanhar")
err = client.TabulationAPI.Delete(ctx, tab.ID)
```

---

### Pausas

```go
// Listar motivos de pausa (paginado)
page, err := client.PauseAPI.List(ctx, szchat.ListOptions{Limit: 20})

// Criar, atualizar, excluir
pause, err := client.PauseAPI.Create(ctx, szchat.Pause{
    Name:    "Almoço",
    MaxTime: 3600,
    Active:  true,
})
pause, err = client.PauseAPI.Update(ctx, pause.ID, pause)
err = client.PauseAPI.Delete(ctx, pause.ID)
```

---

## Paginação

Endpoints paginados retornam `*PaginatedResponse[T]`, um tipo genérico que
envolve o formato de paginação estilo Laravel usado em toda a API do SZChat.

```go
type PaginatedResponse[T any] struct {
    CurrentPage  int
    Data         []T
    Total        int
    PerPage      int
    LastPage     int
    From         int
    To           int
    Path         string
    FirstPageURL string
    LastPageURL  string
    NextPageURL  *string
    PrevPageURL  *string
}
```

Os parâmetros de paginação seguem uma convenção de nomes consistente entre os
filtros de listagem, através do `ListOptions` embutido:

| Campo      | Parâmetro   | Descrição                          |
|------------|-------------|-------------------------------------|
| `Page`     | `page`      | Número da página                   |
| `Limit`    | `limit`     | Quantidade de itens por página      |
| `Paginate` | `paginate`  | Se deve paginar ("0" ou "1")        |

Padrão completo de iteração:

```go
opts := szchat.ListOptions{Limit: 50}
for {
    page, err := client.ContactAPI.List(ctx, szchat.ContactListFilter{ListOptions: opts})
    if err != nil {
        log.Fatal(err)
    }
    // processar page.Data ...
    if page.CurrentPage >= page.LastPage {
        break
    }
    opts.Page = page.CurrentPage + 1
}
```

---

## Tratamento de Erros

Todos os métodos retornam `*APIError` em falhas no nível HTTP.

```go
page, err := client.ContactAPI.List(ctx, szchat.ContactListFilter{})
if err != nil {
    switch {
    case szchat.IsUnauthorized(err):
        // HTTP 401 - credenciais inválidas ou sessão expirada
    case szchat.IsNotFound(err):
        // HTTP 404 - recurso não encontrado
    case szchat.IsValidationError(err):
        // HTTP 422 - falha de validação
    case szchat.IsConflict(err):
        // HTTP 409 - recurso em conflito
    case szchat.IsRateLimited(err):
        // HTTP 429 - muitas requisições, aguarde e tente novamente
    default:
        var apiErr *szchat.APIError
        if errors.As(err, &apiErr) {
            fmt.Printf("Status: %d\n", apiErr.StatusCode)
            fmt.Printf("Mensagem: %s\n", apiErr.Message)
            fmt.Printf("Erros de campo: %v\n", apiErr.FieldErrors)
        }
    }
}
```

O client também retenta automaticamente falhas transitórias (429/502/503/504 e
erros de rede) e renova o token bearer automaticamente em um único 401, então a
maioria dos chamadores só precisa tratar o erro final e definitivo.

---

## Executando os Testes

```bash
go test ./... -v
```

Todos os testes usam `net/http/httptest` — nenhum serviço externo ou variável de
ambiente é necessário.

---

## Cobertura de Endpoints

### Contatos
- [x] `GET /contacts` — Listar contatos (paginado, com filtros)
- [x] `GET /contacts/search` — Buscar contatos
- [x] `GET /contacts/{id}` — Obter um contato
- [x] `POST /contacts` — Criar um contato
- [x] `PUT /contacts/{id}` — Atualizar um contato
- [x] `PUT /contacts/update_fields/{id}` — Atualizar campos customizados de um contato
- [x] `DELETE /contacts/{id}` — Excluir um contato
- [x] `POST /contacts/annotation` — Salvar anotação de contato
- [x] `GET /contacts/{id}/attendances` — Obter estatísticas de atendimento do contato

### Grupos de Contatos
- [x] `GET /contacts/groups` — Listar grupos de contatos
- [x] `POST /contacts/groups` — Criar grupo de contatos
- [x] `PUT /contacts/groups/{id}` — Atualizar grupo de contatos
- [x] `DELETE /contacts/groups/{id}` — Excluir grupo de contatos
- [x] `GET /contacts/groups/contact/{contact_id}` — Listar grupos por contato

### Canais
- [x] `GET /channels` — Listar canais
- [x] `GET /channels/platforms` — Listar plataformas disponíveis
- [x] `GET /channels/platforms/active` — Listar plataformas ativas

### Agentes
- [x] `GET /agents` — Listar agentes (paginado, com filtros)
- [x] `POST /agents` — Criar um agente
- [x] `GET /agents/{id}` — Obter um agente
- [x] `GET /agents/email/{email}` — Obter um agente por email
- [x] `PUT /agents/{id}` — Atualizar um agente
- [x] `DELETE /agents/{id}` — Excluir um ou mais agentes
- [x] `GET /online-status` — Listar status online dos agentes
- [x] `GET /user/agents/online` — Listar agentes/administradores online
- [x] `GET /user/agents/campaigns` — Obter minhas equipes
- [x] `GET /user/agents/attendances` — Obter meus atendimentos
- [x] `GET /user/agents/attendances_plus` — Obter meus atendimentos (detalhado)
- [x] `GET /user/agents/grades` — Obter minhas notas
- [x] `POST /user/agents/toggle/campaign` — Alternar minha equipe ativa
- [x] `POST /user/agents/updateLastInteraction` — Atualizar meu timestamp de última interação

### Mensagens entre Agentes
- [x] `POST /user/agents/messages/send` — Enviar mensagem interna
- [x] `GET /user/agents/messages` — Listar minhas conversas
- [x] `GET /user/agents/messages/read/{agent_id}` — Obter histórico de conversa com um par

### Equipes
- [x] `GET /campaigns` — Listar equipes (paginado)
- [x] `GET /campaigns/resume/{paginate}` — Listar resumo de equipes
- [x] `POST /campaigns/filterByIds` — Filtrar equipes por ids
- [x] `POST /campaigns` — Criar uma equipe
- [x] `PUT /campaigns/{id}` — Atualizar uma equipe
- [x] `DELETE /campaigns/{id}` — Excluir uma equipe

### Administradores
- [x] `GET /admins` — Listar administradores (paginado)
- [x] `GET /admins/groups` — Listar grupos de permissão de administradores
- [x] `POST /admins` — Criar um administrador
- [x] `GET /admins/{id}` — Obter um administrador
- [x] `PUT /admins/{id}` — Atualizar um administrador
- [x] `DELETE /admins/{id}` — Excluir um administrador

### Mensagens
- [x] `POST /message/send` — Enviar uma mensagem
- [x] `POST /message/send_plus` — Enviar mensagem e criar/atualizar contato
- [x] `POST /message/read` — Ler mensagens de uma sessão
- [x] `POST /message/pending` — Obter contagem de mensagens pendentes
- [x] `POST /message/replace_vars` — Substituir variáveis de template
- [x] `POST /message/{id}/annotation` — Adicionar anotação a uma mensagem
- [x] `DELETE /message/{id}/annotation` — Remover anotação de uma mensagem

### Configuração da Aplicação
- [x] `GET /application` — Obter configuração do tenant
- [x] `PUT /application` — Atualizar configuração do tenant
- [x] `GET /application/attendance` — Obter faixa de notas de atendimento
- [x] `PUT /application/attendance` — Atualizar faixa de notas de atendimento
- [x] `GET /application/messages` — Obter configuração de mensagens automáticas
- [x] `PUT /application/messages` — Atualizar configuração de mensagens automáticas

### Tabulações
- [x] `GET /tabulations` — Listar tabulações (paginado)
- [x] `POST /tabulations` — Criar uma tabulação
- [x] `PUT /tabulations/{id}` — Atualizar uma tabulação
- [x] `DELETE /tabulations/{id}` — Excluir uma tabulação

### Pausas
- [x] `GET /pauses` — Listar pausas (paginado)
- [x] `POST /pauses` — Criar uma pausa
- [x] `PUT /pauses/{id}` — Atualizar uma pausa
- [x] `DELETE /pauses/{id}` — Excluir uma pausa

### Autenticação
- [x] `POST /auth/login` — Login
- [x] `GET /auth/me` — Obter perfil do usuário autenticado
- [x] `GET /auth/logout` — Logout
- [x] `GET /auth/refresh` — Renovar token

**Total: 58 endpoints cobertos**

---

## Contribuindo

Contribuições para o szchat-go são bem-vindas! Aqui estão algumas formas de ajudar:

- **Reportar bugs e sugerir funcionalidades** abrindo issues no GitHub
- **Enviar pull requests** com correções de bugs ou novas funcionalidades
- **Melhorar a documentação** para ajudar outros usuários e desenvolvedores

---

## Licença

O szchat-go é distribuído sob a **Licença MIT**.
Para os termos e condições completos da licença, veja o arquivo [LICENSE](LICENSE) no repositório.
