CREATE TABLE lsp_cbt.mst_tipe_user (
  id VARCHAR(2) PRIMARY KEY,
  nama VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE lsp_cbt.mst_user (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NULL,
    id_tipe_user VARCHAR(2) NOT NULL,
    id_prodi VARCHAR(25) NULL,
    nama_prodi VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    CONSTRAINT fk_mst_user_tipe_user FOREIGN KEY (id_tipe_user) REFERENCES lsp_cbt.mst_tipe_user(id)
);