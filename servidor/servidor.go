package servidor

import (
	"crud/banco"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
    "strconv"

	"github.com/gorilla/mux"
)

type user struct {
    ID   uint32 `json:*id`
    Name string `json:*name`
    Email string `json:*email`

}

//POST
func CreateUser(w http.ResponseWriter, r *http.Request){

    // Le o corpo da requisição do tipo json
    bodyRequest, erro := io.ReadAll(r.Body)
    if erro != nil{
        w.Write([]byte("Fail to read body of requests"))
        return
    }


    var user user
    // transforma o json em bytes
    if erro = json.Unmarshal(bodyRequest, &user); erro != nil {
        w.Write([]byte("Error to convert user to struct"))
        return
    }


    db, erro := banco.Conect()
    if erro != nil {
        w.Write([]byte("Error to conect on database MySql"))
        return
    }
    defer db.Close()


    //Prepare statement -> para evitar sql injection
    statement, erro := db.Prepare("insert into ususarios (name, email) values (?, ?)")
    if erro != nil{
        w.Write([]byte("Error the create statement"))
        return
    }
    defer statement.Close()


    insert, erro := statement.Exec(user.Name, user.Email)
    if erro != nil{
        w.Write([]byte("Error to execute statement"))
        return
    }

    //Momento do insert no banco
    idInsert, erro := insert.LastInsertId()
    if erro != nil{
        w.Write([]byte("Error to get ID"))
        return
    }

    // Status_code
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(fmt.Sprintf("User insert with success: %d", idInsert)))
}

//GET
func GetAllUsers(w http.ResponseWriter, r *http.Request){

    db, erro := banco.Conect()
    if erro != nil{
        w.Write([]byte("Erro ao conectar no banco de dados "))
        return   
    }
    defer db.Close()


    linhas, erro := db.Query("SELECT * FROM ususarios")
    if erro != nil {
        w.Write([]byte("Erro ao retornar todos os usuarios !"))
        return
    }
    defer linhas.Close()


    var users []user

    for linhas.Next(){
        var user user

        if erro := linhas.Scan(&user.ID, &user.Name, &user.Email); erro != nil{
            w.Write([]byte("Erro ao escanear o usuário"))
            return
        }

        users = append(users, user)
    }

    w.WriteHeader((http.StatusOK))

    if erro := json.NewEncoder(w).Encode(users); erro != nil{
        w.Write([]byte("Erro ao converter usuario para JSON !"))
        return
    }
}

//GET
func GetUser(w http.ResponseWriter, r *http.Request){

    parameters := mux.Vars(r)

    // Identiica o parametro da requisição e converte no formato especifico
    ID, erro := strconv.ParseUint(parameters["id"], 10, 32)
    if erro != nil {
        w.Write([]byte("Erro ao converter parametro da requisição !\nDeve ser do tipo int"))
        return
    }

    
    db, erro := banco.Conect()
    if erro != nil{
        w.Write([]byte("Erro ao conectar no banco de dados "))
        return   
    }
    defer db.Close()


    linha, erro := db.Query("SELECT * FROM ususarios where id = ?", ID)
    if erro != nil {
        w.Write([]byte("Erro ao retornar o usuarios !"))
        return
    }
  

    var user user

    if linha.Next(){

        if erro := linha.Scan(&user.ID, &user.Name, &user.Email); erro != nil{
            w.Write([]byte("Erro ao escanear o usuário"))
            return
        }
    }

    if user.ID == 0{
        w.Write([]byte("Usuario não encontrado !"))
        w.WriteHeader(http.StatusNoContent)
        return
    }

    if erro := json.NewEncoder(w).Encode(user); erro != nil {
        w.Write([]byte("Erro ao converter resultado em json"))
        return
    }
}

//PUT
func UpdateUser(w http.ResponseWriter, r *http.Request){
    parameters := mux.Vars(r)

    // Identiica o parametro da requisição e converte no formato especifico
    ID, erro := strconv.ParseUint(parameters["id"], 10, 32)
    if erro != nil {
        w.Write([]byte("Erro ao converter parametro da requisição !\nDeve ser do tipo int"))
        return
    }

    
    // Le o corpo da requisição do tipo json
    bodyRequest, erro := io.ReadAll(r.Body)
    if erro != nil{
        w.Write([]byte("Fail to read body of requests"))
        return
    }


    var user user
    // transforma o json em bytes
    if erro = json.Unmarshal(bodyRequest, &user); erro != nil {
        w.Write([]byte("Error to convert user to struct"))
        return
    }

    db, erro := banco.Conect()
    if erro != nil{
        w.Write([]byte("Error to conect on database MySql"))
        return
    }
    defer db.Close()

    statement, erro := db.Prepare("Update ususarios set name = ?, email = ? where id = ?")
    if erro != nil {
        w.Write([]byte("Erro ao criar o statement"))
        return
    }
    defer statement.Close()

    if _, erro := statement.Exec(user.Name, user.Email, ID); erro != nil {
        w.Write([]byte("Erro ao atualizar usuario"))
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

//DEL
func DelUser(w http.ResponseWriter, r *http.Request){
    parameters := mux.Vars(r)

    // Identiica o parametro da requisição e converte no formato especifico
    ID, erro := strconv.ParseUint(parameters["id"], 10, 32)
    if erro != nil {
        w.Write([]byte("Erro ao converter parametro da requisição !\nDeve ser do tipo int"))
        return
    }

    db, erro := banco.Conect()
    if erro != nil{
        w.Write([]byte("Error to conect on database MySql"))
        return
    }
    defer db.Close()

    statement, erro := db.Prepare(" DELETE FROM ususarios WHERE id = ? ")
    if erro != nil {
        w.Write([]byte("Erro ao criar statement"))
        return
    }
    defer statement.Close()

    if _, erro := statement.Exec(ID); erro != nil {
        w.Write([]byte("Erro ao deletar usuario"))
        return
    }

    w.WriteHeader(http.StatusNoContent)

}