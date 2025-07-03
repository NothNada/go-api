package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"api/db"
	"api/models"

	"github.com/gorilla/mux"
)

// CreateData cria um novo registro na tabela dados
func CreateData(w http.ResponseWriter, r *http.Request) {
	var data models.Data
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Define a data de registro como a data atual se não for fornecida
	if data.DataRegistrada.IsZero() {
		data.DataRegistrada = time.Now()
	}

	stmt, err := db.DB.Prepare("INSERT INTO dados(usuario_que_registrou, data_registrada, descricao_curta, descricao_longa, data_de_expiracao) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(data.UsuarioQueRegistrou, data.DataRegistrada, data.DescricaoCurta, data.DescricaoLonga, data.DataDeExpiracao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

// GetData recupera um registro por ID
func GetData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var data models.Data
	row := db.DB.QueryRow("SELECT id, usuario_que_registrou, data_registrada, descricao_curta, descricao_longa, data_de_expiracao FROM dados WHERE id = ?", id)
	err = row.Scan(&data.ID, &data.UsuarioQueRegistrou, &data.DataRegistrada, &data.DescricaoCurta, &data.DescricaoLonga, &data.DataDeExpiracao)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Registro não encontrado", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GetAllData recupera todos os registros
func GetAllData(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, usuario_que_registrou, data_registrada, descricao_curta, descricao_longa, data_de_expiracao FROM dados")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var allData []models.Data
	for rows.Next() {
		var data models.Data
		if err := rows.Scan(&data.ID, &data.UsuarioQueRegistrou, &data.DataRegistrada, &data.DescricaoCurta, &data.DescricaoLonga, &data.DataDeExpiracao); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		allData = append(allData, data)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allData)
}

// UpdateData atualiza um registro existente
func UpdateData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var data models.Data
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.DB.Prepare("UPDATE dados SET usuario_que_registrou = ?, data_registrada = ?, descricao_curta = ?, descricao_longa = ?, data_de_expiracao = ? WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(data.UsuarioQueRegistrou, data.DataRegistrada, data.DescricaoCurta, data.DescricaoLonga, data.DataDeExpiracao, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Registro atualizado com sucesso!"})
}

// DeleteData exclui um registro
func DeleteData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	stmt, err := db.DB.Prepare("DELETE FROM dados WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Registro não encontrado para exclusão", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Registro excluído com sucesso!"})
}

func CheckForExpiredData() {
	ticker := time.NewTicker(time.Hour) // Verifica a cada 24 horas (ajuste conforme necessário)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Verificando registros expirados...")
		rows, err := db.DB.Query("SELECT id, usuario_que_registrou, data_registrada, descricao_curta, descricao_longa, data_de_expiracao FROM dados WHERE data_de_expiracao < ?", time.Now())
		if err != nil {
			log.Printf("Erro ao consultar dados expirados: %v\n", err)
			continue
		}
		defer rows.Close()

		for rows.Next() {
			var data models.Data
			if err := rows.Scan(&data.ID, &data.UsuarioQueRegistrou, &data.DataRegistrada, &data.DescricaoCurta, &data.DescricaoLonga, &data.DataDeExpiracao); err != nil {
				log.Printf("Erro ao escanear linha de dados expirados: %v\n", err)
				continue
			}
			log.Printf("REGISTRO EXPIRADO ENCONTRADO: ID: %d, Descrição: %s, Expira em: %s\n", data.ID, data.DescricaoCurta, data.DataDeExpiracao.Format("2006-01-02"))

			deleteStmt, err := db.DB.Prepare("DELETE FROM dados WHERE id = ?")
			if err != nil {
				log.Printf("Erro ao preparar exclusão para registro %d: %v\n", data.ID, err)
				continue
			}
			_, err = deleteStmt.Exec(data.ID)
			if err != nil {
				log.Printf("Erro ao excluir registro expirado %d: %v\n", data.ID, err)
			} else {
				log.Printf("Registro expirado %d excluído com sucesso.\n", data.ID)
			}
			deleteStmt.Close()
		}

		if err = rows.Err(); err != nil {
			log.Printf("Erro após iteração de linhas: %v\n", err)
		}
		log.Println("Verificação de registros expirados concluída.")
	}
}
