// Conteúdo institucional do SISMARQ — fonte única para a página inicial, o
// Sobre e o FAQ. Tudo aqui vem dos documentos de docs/SISMARQ/: a Lei
// Municipal nº 10.301/2019, as Tabelas de Temporalidade (TTDD) publicadas
// no DIORONDON-e e o artigo técnico do SISMARQ. Cada item cita a origem;
// nada entra sem base nesses documentos (ver docs/ROADMAP_SISMARQ.md).

export interface Fonte {
  /** Ex.: "Lei 10.301/2019, art. 18". */
  ref: string;
}

export const LEI = {
  numero: "Lei Municipal nº 10.301",
  data: "10 de junho de 2019",
  ementa:
    "Dispõe sobre a criação do Arquivo Público Municipal de Rondonópolis, define as diretrizes da Política Municipal de Arquivos Públicos e Privados e cria o Sistema Municipal de Arquivos - SISMARQ.",
  publicacao: "DIORONDON-e",
} as const;

export const MISSAO =
  "Elaborar, acompanhar e implementar a Política Municipal de Arquivos e Protocolos; promover a gestão e o recolhimento dos documentos permanentes, de valor histórico, produzidos pelo Poder Executivo Municipal e preservar e facilitar o acesso aos documentos sob sua guarda.";

/** Art. 7º — finalidades do SISMARQ. */
export const FINALIDADES: string[] = [
  "Garantir ao cidadão e aos órgãos da administração municipal, de forma ágil, transparente e segura, o acesso aos documentos de arquivo, resguardados o sigilo e as restrições legais.",
  "Harmonizar as atividades nas diversas fases da administração dos documentos, atendendo às especificidades de cada órgão.",
  "Integrar e coordenar a gestão de documentos dos órgãos setoriais e seccionais.",
  "Disseminar normas relativas à gestão de documentos de arquivo.",
  "Racionalizar a produção da documentação arquivística pública.",
  "Racionalizar e reduzir os custos operacionais e de armazenagem.",
  "Preservar o patrimônio documental arquivístico da administração municipal.",
  "Articular-se com os demais sistemas que atuam na gestão da informação pública municipal.",
];

export interface Idade {
  nome: string;
  descricao: string;
  onde: string;
  fonte: string;
}

/** Art. 18 — as três idades; art. 29 — onde cada uma fica. */
export const IDADES: Idade[] = [
  {
    nome: "Corrente",
    descricao: "Documentos em curso ou que, mesmo sem movimentação, são objeto de consultas frequentes.",
    onde: "Na unidade produtora (arquivo corrente).",
    fonte: "Lei 10.301/2019, arts. 18, § 1º, e 29, I",
  },
  {
    nome: "Intermediária",
    descricao:
      "Documentos que não são mais de uso corrente e, por interesse administrativo, aguardam a eliminação ou o recolhimento para a guarda permanente.",
    onde: "Na unidade com atribuições de arquivo intermediário do órgão.",
    fonte: "Lei 10.301/2019, arts. 18, § 2º, e 29, II",
  },
  {
    nome: "Permanente",
    descricao:
      "Conjuntos de valor histórico, probatório e informativo, que devem ser definitivamente preservados. São inalienáveis e imprescritíveis.",
    onde: "No Arquivo Público Municipal (ou no órgão produtor, enquanto necessários).",
    fonte: "Lei 10.301/2019, arts. 18, § 3º, 20 e 42",
  },
];

/** Art. 30 — prazos que definem o tempo de guarda. */
export const PRAZOS = [
  { nome: "Vigência", descricao: "Tempo em que o documento produz efeitos administrativos e legais plenos." },
  {
    nome: "Prescrição",
    descricao:
      "Tempo em que se pode recorrer ao Judiciário para fazer valer direitos. A guarda é dilatada se a prescrição for interrompida ou suspensa.",
  },
  { nome: "Precaução", descricao: "Tempo de guarda por precaução, antes de eliminar ou recolher para a guarda permanente." },
];

export interface Etapa {
  titulo: string;
  descricao: string;
  fonte: string;
}

/** O rito da eliminação (arts. 19, 33 a 40 e 65; Anexos I a IV). */
export const RITO_ELIMINACAO: Etapa[] = [
  {
    titulo: "Avaliação pela CSPAD",
    descricao:
      "A Comissão Setorial Permanente de Avaliação de Documentos do órgão identifica os documentos com prazo cumprido na Tabela de Temporalidade e confere se há ações judiciais em curso.",
    fonte: "Lei 10.301/2019, arts. 19, 32 e 44",
  },
  {
    titulo: "Relação de Eliminação (Anexo I)",
    descricao:
      "Os documentos são listados com código de classificação, assunto, datas-limite e quantificação. A relação é assinada e depende do “Autorizo” do titular do Arquivo Público.",
    fonte: "Lei 10.301/2019, art. 34 e Anexo I",
  },
  {
    titulo: "Ata da CSPAD (Anexo IV)",
    descricao: "A comissão aprova a listagem em reunião registrada em ata, com a lista de presença.",
    fonte: "Lei 10.301/2019, art. 33 e Anexo IV",
  },
  {
    titulo: "Edital de Ciência (Anexo III)",
    descricao:
      "O edital é publicado no Diário Oficial do Município e abre 30 dias para manifestações. Nesse prazo, os interessados podem pedir o desentranhamento de documentos ou cópias.",
    fonte: "Lei 10.301/2019, art. 37",
  },
  {
    titulo: "Amostragem",
    descricao:
      "Dos documentos destinados à eliminação, uma amostra representativa, escolhida por critérios qualitativos e quantitativos, passa à guarda permanente.",
    fonte: "Lei 10.301/2019, art. 39",
  },
  {
    titulo: "Eliminação e Termo (Anexo II)",
    descricao:
      "A eliminação é feita por fragmentação, e o papel é doado (Decreto Federal nº 5.940/2006). O Termo de Eliminação registra o ato, e uma cópia vai à Divisão de Arquivo Público e à Secretaria de Administração.",
    fonte: "Lei 10.301/2019, arts. 38 e 40 e Anexo II",
  },
];

export interface Pergunta {
  q: string;
  a: string;
  fonte: string;
}

export interface SecaoFAQ {
  categoria: string;
  itens: Pergunta[];
}

export const FAQ: SecaoFAQ[] = [
  {
    categoria: "O Arquivo Público e o SISMARQ",
    itens: [
      {
        q: "O que é o SISMARQ?",
        a: "É o Sistema Municipal de Arquivos de Rondonópolis, criado pela Lei nº 10.301/2019. Ele organiza, sob a forma de sistema, as atividades de gestão de documentos de todos os órgãos e entidades da administração pública municipal.",
        fonte: "Lei 10.301/2019, arts. 1º e 6º",
      },
      {
        q: "Quem faz parte do SISMARQ?",
        a: "O Arquivo Público Municipal, como órgão central, e, como órgãos setoriais, as unidades das secretarias que coordenam o protocolo e a gestão de documentos. Os servidores que chefiam atividades de arquivo são os agentes setoriais do sistema.",
        fonte: "Lei 10.301/2019, arts. 8º e 9º",
      },
      {
        q: "A que órgão o Arquivo Público Municipal está ligado?",
        a: "A lei o subordina ao Gabinete da Secretaria Municipal de Administração, sob a responsabilidade da Divisão de Arquivo Público. As Tabelas de Temporalidade publicadas em 2025 e 2026 são assinadas pela Superintendência de Arquivo Público da Secretaria Municipal de Administração, Gestão de Pessoas e Inovação.",
        fonte: "Lei 10.301/2019, art. 5º; TTDD (DIORONDON-e nº 6.017/2025)",
      },
      {
        q: "Qual a relação do SISMARQ com o SEI?",
        a: "São complementares. O SEI tramita os processos em curso (fase corrente). O SISMARQ cuida da guarda intermediária e permanente, do passivo físico e digitalizado e do rito de eliminação com os modelos da lei municipal.",
        fonte: "Artigo técnico do SISMARQ, Quadro 1; Lei 10.301/2019, art. 18",
      },
    ],
  },
  {
    categoria: "Acesso aos documentos",
    itens: [
      {
        q: "Posso consultar documentos públicos municipais?",
        a: "Sim. A lei assegura ao cidadão o acesso pleno aos documentos públicos municipais, de forma ágil e transparente. As exceções são o sigilo imprescindível à segurança da sociedade e do Estado e a proteção da intimidade, da vida privada, da honra e da imagem das pessoas.",
        fonte: "Lei 10.301/2019, arts. 2º e 7º, I",
      },
      {
        q: "Documentos de arquivos privados podem ser consultados?",
        a: "Arquivos privados declarados de interesse público e social podem ser consultados mediante autorização do proprietário ou possuidor.",
        fonte: "Lei 10.301/2019, art. 56, § 2º",
      },
    ],
  },
  {
    categoria: "Classificação e Tabela de Temporalidade",
    itens: [
      {
        q: "O que é a Tabela de Temporalidade e Destinação de Documentos (TTDD)?",
        a: "É o instrumento resultante da avaliação documental, aprovado pela autoridade competente, que define os prazos de guarda e a destinação de cada série: eliminação ou guarda permanente. As tabelas das secretarias são aprovadas pela Comissão Central e publicadas no Diário Oficial do Município.",
        fonte: "Lei 10.301/2019, arts. 27, 28 e 32, § 8º",
      },
      {
        q: "Onde consulto a TTDD em vigor?",
        a: "Na página Temporalidade deste portal, com as séries das tabelas publicadas no DIORONDON-e, os prazos nas fases corrente e intermediária, a destinação e as observações.",
        fonte: "TTDD publicada no DIORONDON-e (2025–2026)",
      },
      {
        q: "Como os documentos são classificados?",
        a: "Pelo método duplex, do geral para o particular. O código tem cinco unidades: órgão produtor, classe, subclasse, assunto e espécie documental. Nas tabelas publicadas, ele aparece como órgão, função, subfunção e série (ex.: 2.0.01.00.01).",
        fonte: "Lei 10.301/2019, arts. 24 e 25",
      },
      {
        q: "Como se conta o prazo de guarda?",
        a: "A TTDD indica quantos anos o documento fica no arquivo corrente e no intermediário. Esses prazos consideram a vigência, a prescrição e a precaução. Se a prescrição for interrompida ou suspensa, a guarda é dilatada.",
        fonte: "Lei 10.301/2019, arts. 29 e 30",
      },
    ],
  },
  {
    categoria: "Avaliação e eliminação",
    itens: [
      {
        q: "Quem pode autorizar a eliminação de documentos?",
        a: "Nenhuma eliminação ocorre sem autorização prévia do Arquivo Público Municipal. Só pode ser eliminado o que está previsto na Tabela de Temporalidade, depois de cumprido o prazo de guarda.",
        fonte: "Lei 10.301/2019, arts. 19 e 65",
      },
      {
        q: "O que são a CCPAD e as CSPADs?",
        a: "A Comissão Central Permanente de Avaliação de Documentos fica junto ao Arquivo Público. Cada órgão tem uma Comissão Setorial (CSPAD), multidisciplinar e preferencialmente em número ímpar, com servidor de Arquivologia ou História, da assessoria jurídica, de administração e finanças e das áreas produtoras.",
        fonte: "Lei 10.301/2019, art. 31",
      },
      {
        q: "Fui citado num Edital de Ciência de Eliminação. O que posso fazer?",
        a: "Em até 30 dias da publicação, você pode pedir, às suas expensas, o desentranhamento de documentos ou cópias de peças, com petição dirigida à Comissão Permanente de Avaliação do órgão, demonstrando a legitimidade do pedido.",
        fonte: "Lei 10.301/2019, art. 37, § 2º, e Anexo III",
      },
      {
        q: "Documentos com processo judicial podem ser eliminados?",
        a: "Não enquanto houver ação em curso. As comissões devem analisar a existência de ações judiciais nas quais o Município figure como autor ou réu antes de aplicar os prazos de guarda.",
        fonte: "Lei 10.301/2019, arts. 32, II, e 44",
      },
      {
        q: "Depois de digitalizado, o documento permanente pode ser descartado?",
        a: "Não. Documentos de guarda permanente não podem ser eliminados após microfilmagem, digitalização ou qualquer outra forma de reprodução.",
        fonte: "Lei 10.301/2019, art. 42",
      },
      {
        q: "O que acontece com quem destrói documentos de forma irregular?",
        a: "Fica sujeito à responsabilidade penal, civil e administrativa, nos termos do art. 25 da Lei Federal nº 8.159/1991 e da Lei nº 9.605/1998, quem desfigurar ou destruir documento de valor permanente ou de interesse público e social.",
        fonte: "Lei 10.301/2019, art. 66",
      },
    ],
  },
  {
    categoria: "Recolhimento e arquivos privados",
    itens: [
      {
        q: "Como os documentos permanentes chegam ao Arquivo Público?",
        a: "Pelo recolhimento. Os documentos devem estar classificados, avaliados, organizados, higienizados, acondicionados e acompanhados de instrumento descritivo. As despesas de preparo e transporte são do órgão produtor, e o Arquivo Público registra a entrada e envia cópia do registro à unidade de origem.",
        fonte: "Lei 10.301/2019, arts. 5º, VII, e 53",
      },
      {
        q: "Um arquivo particular pode ser protegido pelo Município?",
        a: "Sim. Arquivos privados relevantes para a história, a cultura e o desenvolvimento científico e tecnológico do município podem ser declarados de interesse público e social por decreto, após parecer de comissão especial. Eles podem ser doados ou depositados no Arquivo Público, e o Município tem preferência, em 60 dias, em caso de alienação.",
        fonte: "Lei 10.301/2019, arts. 56 a 58",
      },
    ],
  },
];

/** Base legal citada pelos documentos de docs/SISMARQ. */
export const BASE_LEGAL: { norma: string; tema: string }[] = [
  { norma: "Constituição Federal, art. 216, § 2º", tema: "Dever do Poder Público de gerir a documentação e franquear a consulta" },
  { norma: "Lei Federal nº 8.159/1991", tema: "Política Nacional de Arquivos Públicos e Privados" },
  { norma: "Lei Municipal nº 10.301/2019", tema: "Cria o Arquivo Público Municipal e o SISMARQ" },
  { norma: "Resolução CONARQ nº 7/1997", tema: "Procedimentos de eliminação de documentos no Poder Público" },
  { norma: "e-ARQ Brasil (CONARQ)", tema: "Requisitos para sistemas informatizados de gestão arquivística (art. 13 da lei municipal)" },
  { norma: "Decreto Federal nº 10.278/2020", tema: "Digitalização com o mesmo valor probatório do original" },
  { norma: "Lei Federal nº 9.605/1998", tema: "Responsabilização pela destruição de documentos de valor permanente" },
];

/** As TTDDs publicadas que compõem docs/SISMARQ/Tabela temporalidade.pdf. */
export const TTDD_PUBLICACOES = [
  "DIORONDON-e nº 6.017, de 25/08/2025",
  "DIORONDON-e nº 6.111, de 13/01/2026 (suplementar)",
  "DIORONDON-e nº 6.275, de 11/09/2026",
  "DIORONDON-e nº 6.281, de 21/09/2026",
];

/** Art. 5º — competências do Arquivo Público Municipal. */
export const COMPETENCIAS: string[] = [
  "Formular a política municipal de arquivos e exercer orientação normativa para a gestão e a proteção dos documentos, qualquer que seja o suporte.",
  "Implementar, acompanhar e supervisionar a gestão dos documentos produzidos, recebidos e acumulados pela administração municipal.",
  "Promover a organização, a preservação e o acesso aos documentos de valor permanente ou histórico recolhidos.",
  "Elaborar e divulgar normas para documentos físicos, digitais e híbridos, conforme o CONARQ e o SIARQ/MT.",
  "Coordenar a classificação e a avaliação dos documentos e aprovar os Planos de Classificação e as Tabelas de Temporalidade.",
  "Autorizar a eliminação dos documentos desprovidos de valor permanente (art. 9º da Lei Federal nº 8.159/1991).",
  "Acompanhar o recolhimento dos documentos permanentes, registrar a entrada e enviar cópia do registro à unidade de origem.",
  "Promover o treinamento e a orientação técnica dos profissionais de arquivo das unidades do SISMARQ.",
  "Incentivar a cooperação entre os órgãos para integrar as atividades arquivísticas.",
  "Difundir informações sobre o Arquivo Público e garantir o acesso aos documentos, observadas as restrições legais.",
  "Realizar ações educativas e culturais para divulgar e preservar o patrimônio documental do município.",
];

/** Artigo técnico, Quadro 1 — a divisão de responsabilidades. */
export const SEI_X_SISMARQ = {
  sei: {
    titulo: "SEI — processo eletrônico",
    escopo: "Fase corrente (art. 18, § 1º)",
    itens: [
      "Documentos nato-digitais da rotina diária",
      "Tramitação de processos entre as secretarias",
      "Despachos, assinaturas e instrução processual",
      "Autuação dos processos das comissões de avaliação",
    ],
  },
  sismarq: {
    titulo: "SISMARQ — repositório arquivístico",
    escopo: "Fases intermediária e permanente (art. 18, §§ 2º e 3º)",
    itens: [
      "Passivo físico e documentos digitalizados",
      "Digitalização com PDF/A, OCR, metadados e hash (Decreto 10.278/2020)",
      "Localização física: depósito, estante e caixa",
      "Temporalidade, bloqueio judicial e o rito de eliminação com os Anexos I a IV",
    ],
  },
} as const;
