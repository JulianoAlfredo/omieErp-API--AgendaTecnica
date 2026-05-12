CREATE TABLE amm_contas_omie_x_agenda
(
    id_conta_agenda VARCHAR(255) NULL,
    id_os FLOAT NULL,
    id_conta_omie FLOAT NULL,
    faturada INT NULL,
    boleto_gerado VARCHAR(100) NULL,
    id_nf INT NULL,
    numero_nf VARCHAR(255) NULL,
    numero_rps VARCHAR(255) NULL,
    numero_os VARCHAR(100) NULL,
    id_cliente FLOAT NULL,
    codigo_barras_boleto VARCHAR(255) NULL,
    boleto_numero VARCHAR(255) NULL
);

CREATE TABLE amm_omie_faturamento_log
(
    id INT IDENTITY(1,1) PRIMARY KEY,
    cod_int_os VARCHAR(255) NOT NULL,
    etapa VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL,
    mensagem NVARCHAR(MAX) NULL,
    dados NVARCHAR(MAX) NULL,
    criado_em DATETIME NOT NULL DEFAULT GETDATE()
);

CREATE TABLE amm_omie_relaciona_clientes
(
    cliente_agenda NVARCHAR(200),
    cliente_omie BIGINT,
    cnpj NVARCHAR(30)
);

ALTER TABLE amm_contas_omie_x_agenda
    ADD 
        data_baixa       DATETIME       NULL,
        data_cred        DATETIME       NULL,
        observacao_baixa NVARCHAR(MAX)  NULL,
        valor_baixa      DECIMAL(15, 2) NULL;

CREATE TABLE amm_omie_logs
(
    id INT IDENTITY(1,1) PRIMARY KEY,
    metodo VARCHAR(10) NOT NULL,
    rota VARCHAR(255) NOT NULL,
    status_code INT NOT NULL,
    ip_origem VARCHAR(50) NULL,
    corpo_requisicao NVARCHAR(MAX) NULL,
    corpo_resposta NVARCHAR(MAX) NULL,
    duracao_ms INT NOT NULL,
    criado_em DATETIME NOT NULL DEFAULT GETDATE()
);