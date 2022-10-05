package main

import (
	"strings"
)

type Backpack []Item

type Item struct {
	name string
}

type Room struct {
	id         int
	def        bool
	name       string
	roomItem   []RoomItem
	associated map[int]string
}

type RoomItem struct {
	name  string
	items []Item
}

type Rooms []Room

type Door struct {
	isLocked bool
}

type Players []Player

type Player struct {
	id           int
	room         Room
	backpack     Backpack
	haveBackpack bool
	quests       Quests
}
type Quests []Quest

type Quest struct {
	isDone bool
	name   string
}

func DefaultString(prod, defaultString string) string {
	if len(prod) == 0 {
		return defaultString
	}
	return prod
}

func DefaultForEmptyString(s, defaultString string) string {
	if len(s) == 0 {
		return defaultString
	}
	return s
}

func (b Backpack) addProduct() string {
	var prod string
	for _, a := range b {
		prod += a.name

	}
	prod = DefaultString(prod, "нет такого")
	return prod
}

func (b Backpack) ItemsToString() string {
	var res string
	for _, v := range b {
		res += v.name
	}
	res = DefaultForEmptyString(res, "Ничего нет")
	return res
}

func (ri RoomItem) HaveItem(s string) (Item, bool) {
	for _, v := range ri.items {
		if v.name == s {
			return v, true
		}
	}
	return Item{}, false
}

func (ri RoomItem) ItemToString() string {
	s := make([]string, 0, len(ri.items))
	for _, v := range ri.items {
		s = append(s, v.name)
	}
	return strings.Join(s, ", ")
}

func (r Room) ItemsToString() string {
	var res string
	for _, v := range r.roomItem {
		res += v.name + " " + v.ItemToString()
	}
	res = DefaultForEmptyString(res, "ничего интересного")
	return res
}

func (r Rooms) GetByName(name string) Room {
	for _, v := range r {
		if v.name == name {
			return v
		}
	}
	panic("Комната не найдена")
}

func (r *Room) Lookup() string {
	if r.def {
		return "ты находишься на кухне, " + r.ItemsToString() + ", надо " + currentPlayer.quests.ToString() + ". можно пройти - " + r.AssociatedToString()
	}
	if len(room.ItemsToString()) > 0 {
		return r.ItemsToString() + ". можно пройти - " + r.AssociatedToString()
	}
	return "пустая " + r.name + ". можно пройти - " + r.AssociatedToString()
}

func (r Room) AssociatedToString() string {
	var res string
	roomsArray := make([]string, 0, 2)
	for _, v := range r.associated {
		roomsArray = append(roomsArray, v)
	}
	res = strings.Join(roomsArray, ", ")

	return res
}
func (r Room) CanIGoToRoom(s string) bool {
	room := rooms.GetByName(s)
	_, ok := r.associated[room.id]
	return ok
}

func (r Room) GetItemByName(s string) (Item, bool) {
	for _, v := range r.roomItem {
		if item, ok := v.HaveItem(s); ok {
			return item, true
		}
	}
	return Item{}, false
}

func (p *Players) NewPlayer(room Room) Player {
	pl := Player{id: len(*p) + 1, room: room}
	pl.quests.StartQuests()
	*p = append(*p, pl)
	return pl
}

func (p *Player) ClotheBackpack() string {
	currentPlayer.haveBackpack = true
	return "вы надели: рюкзак"
}

func (p *Player) TakeItem(s string) string {
	if !currentPlayer.haveBackpack {
		return "некуда класть"
	}
	if i, ok := p.room.GetItemByName(s); ok {
		currentPlayer.backpack = append(currentPlayer.backpack, i)
		return "предмет добавлен в инвентарь: " + i.name
	}
	return "нет такого"
}
func (p *Player) Go(nameRoom string) string {
	if ok := p.goToRoom(rooms.GetByName(nameRoom)); ok {
		currentPlayer.goToRoom(rooms.GetByName(nameRoom))
		return p.room.name
	}
	return "нет пути в " + nameRoom
}

func (p *Player) goToRoom(r Room) bool {
	if currentPlayer.room.CanIGoToRoom(r.name) {
		p.room = r
		return true
	}
	return false
}

func (q *Quests) StartQuests() {
	*q = Quests{
		{name: "собрать рюкзак"},
		{name: "идти в универ"},
	}
}

func (q Quests) ToString() string {
	var res string
	quests := make([]string, 0, 2)
	for _, v := range q {
		if v.isDone == false {
			quests = append(quests, v.name)
		}
	}
	res = strings.Join(quests, " и ")
	res = DefaultForEmptyString(res, "заданий нет")

	return res
}

func main() {
	/*
		в этой функции можно ничего не писать
		но тогда у вас не будет работать через go run main.go
		очень круто будет сделать построчный ввод команд тут, хотя это и не требуется по заданию
	*/
}

var players Players
var currentPlayer Player
var room Room
var rooms Rooms

func initGame() {
	rooms = Rooms{
		{id: 1, name: "кухня", associated: map[int]string{2: "коридор"}, def: true,
			roomItem: []RoomItem{
				{"на столе:", []Item{{"чай"}}},
			},
		},
		{id: 2, name: "коридор", associated: map[int]string{1: "кухня", 3: "комната", 4: "улица"}},
		{id: 3, name: "комната", associated: map[int]string{2: "коридор"},
			roomItem: []RoomItem{
				{name: "на столе:", items: []Item{{"ключи"}, {"конспекты"}}},
				{"на стуле -", []Item{{"рюкзак"}}},
			},
		},
	}

	currentPlayer = players.NewPlayer(rooms.GetByName("кухня"))
}

func handleCommand(command string) string {
	initGame()
	commands := strings.Split(command, " ")
	switch commands[0] {
	case "осмотреться":
		return currentPlayer.room.Lookup()
	case "идти":
		currentPlayer.Go(commands[1])
		return currentPlayer.room.Lookup()
	case "надеть":
		return currentPlayer.ClotheBackpack()
	case "взять":
		return currentPlayer.TakeItem(commands[1])
	case "применить":
		return currentPlayer.useItem(player.fromStringToItemBackpack(commands[1]), commands[2], commands[1])

	default:
		return "неизвестная команда"
	}
	return "not implemented"
}
