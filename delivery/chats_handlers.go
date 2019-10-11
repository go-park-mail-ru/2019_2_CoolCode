package delivery

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_2_CoolCode/models"
	"github.com/go-park-mail-ru/2019_2_CoolCode/repository"
	"github.com/go-park-mail-ru/2019_2_CoolCode/useCase"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ChatHandlers struct {
	Chats useCase.ChatsUseCase
}

func NewChatHandlers() ChatHandlers {
	return ChatHandlers{Chats: useCase.NewChatsUseCase(repository.NewChatArrayRepository())}
}

func (c *ChatHandlers)PostChat(w http.ResponseWriter, r *http.Request){
	//TODO: Check auth
	var newChatModel models.CreateChatModel
	decoder:=json.NewDecoder(r.Body)
	err:=decoder.Decode(&newChatModel)
	if err!=nil{

	}
	model:=models.Chat{
		Name:          "MEM",
	}
	err=c.Chats.PutChat(&model)
	w.WriteHeader(http.StatusOK)
}

func (c *ChatHandlers)GetChatsByUser(w http.ResponseWriter, r *http.Request){
	//TODO: Check auth
	requestedID, _ := strconv.Atoi(mux.Vars(r)["id"])
	chats,err:=c.Chats.GetChatByUserID(uint64(requestedID))
	if err!=nil{

	}
	jsonChat,err:=json.Marshal(chats)
	_,err=w.Write(jsonChat)
}


