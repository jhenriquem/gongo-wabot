# 🤖 Gongo WaBot 
gongo-wabot é um bot para WhatsApp desenvolvido para automatizar respostas baseadas em "stages" (etapas de conversa). Ele permite interações dinâmicas e organizadas, facilitando a comunicação automatizada.

Inicialmente, criei ele para um fim comercial, porém decidi publicá-lo no GitHub com o objetivo de ajudar aqueles que querem desenvolver um projeto semelhante e estão com dificuldades 🚀 . Ele foi construido para uma livraria, porém é facil modificalo para outro fins. 

 ### 📌 Para que serve?

O gongo-wabot foi criado para responder mensagens automaticamente no WhatsApp, seguindo um fluxo de conversa estruturado. Ele pode ser utilizado para:

- 🤖 Atendimento automatizado
- 💬 Respostas predefinidas baseadas em estágios da conversa
- 📞 Auxílio em processos interativos, como pedidos, suporte ou informações automatizadas

### 🛠️ Como foi construído?
O projeto foi desenvolvido utilizando:

- 📝 Linguagem: Go (Golang)
- 📚 Bibliotecas: [whatsmeow](https://pkg.go.dev/go.mau.fi/whatsmeow)
- 🏗️ Arquitetura: Modularizada para facilitar a manutenção e expansão
- 🔄 Baseado em sessões, conversas com pessoas distintas e em estágios de conversa 

### 📂 Estrutura do projeto

```
~/gongo-wabot/
├── internal/           # Implementação de lógica interna
│   ├── bot/            # Configuração e lógica do bot
│   │   ├── connect.go  # Conexão com o WhatsApp
│   │   └── handler.go  # Manipulação de mensagens recebidas
│   ├── session/        # Gerenciamento de estágios da conversa e sessões
│   │   ├── config.go   # Configuração da sessão
│   │   ├── methods.go  # Métodos de manipulação da sessão
│   │   └── stages.go   # Definição dos estágios da conversa
│   ├── types/          # Definição de tipos usados no projeto
│   ├── helpers/        # Funções auxiliares
│   ├── services/       # Configuração de serviços externos (banco de dados, storage)
│   ├── storage/        # Armazenamento da sessão
│   └── utils/          # Funções utilitárias 
├── main.go             # Arquivo principal do bot
└── README.md           # Este arquivo
```

### ⚙️  Funcionalidades

- 📥 Receber mensagens: O bot captura mensagens enviadas por usuários.
- 🔄 Gerenciamento de "sessions" : Uma conversa com um número de telefone representa um sessão
- ⚙️  Lógica de "stages": O bot responde com base no estágio atual da conversa.
- ✉️  Envio de respostas automatizadas: Respostas são definidas conforme o contexto da interação.

### 🚀 Como rodar o projeto?

#### 1 Pré-requisitos
- Go instalado 
- Conta no WhatsApp ou um número configurado com whatsmeow
- Banco de dados PostgreSQL e SQLite configurados

#### 2 Instalação 
Clone o repositório e entre na pasta do projeto:
```bash
git clone https://github.com/jhenriquem/gongo-wabot.git
cd gongo-wabot
``` 

Baixe as dependências:
```bash
go mod tidy
```

#### 3 Configuração
Antes de rodar o bot, configure as variáveis de ambiente no arquivo .env:
```sh
DATABASE_URL=postgres://user:password@host:port/dbname
STORAGE_KEY=
STORAGE_URL=
```

#### 4 Executando o bot
```bash
go run main.go
```

#### 📌 Como modificar o fluxo de mensagens?
Para personalizar as respostas, edite o arquivo stages.go e methods.go dentro da pasta internal/session/.
Cada estágio da conversa pode ser configurado conforme a necessidade do seu negócio.

