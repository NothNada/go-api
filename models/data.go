package models

import "time"

type Data struct {
	ID                  int       `json:"id"`
	UsuarioQueRegistrou string    `json:"usuario_que_registrou"`
	DataRegistrada      time.Time `json:"data_registrada"`
	DescricaoCurta      string    `json:"descricao_curta"`
	DescricaoLonga      string    `json:"descricao_longa"`
	DataDeExpiracao     time.Time `json:"data_de_expiracao"`
}
