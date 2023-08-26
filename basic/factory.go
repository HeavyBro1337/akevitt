package basic

func CreateItem[T Item](item T, ip *ItemParams) T {
	mapItemParams(item, ip)

	return item
}

func mapItemParams(item Item, itemParams *ItemParams) {
	item.SetName(itemParams.Name)
	item.SetDescription(itemParams.Description)
	item.SetCallback(itemParams.onUse)
	item.SetQuantity(itemParams.Quantity)
}

type ItemParams struct {
	Quantity    int
	Name        string
	Description string
	onUse       InteractFunc
}

func NewItemParams() *ItemParams {
	return &ItemParams{Quantity: 1}
}

func (ip *ItemParams) WithName(name string) *ItemParams {
	ip.Name = name
	return ip
}
func (ip *ItemParams) WithDescription(description string) *ItemParams {
	ip.Description = description
	return ip
}
func (ip *ItemParams) WithQuantity(quantity int) *ItemParams {
	ip.Quantity = quantity
	return ip
}
func (ip *ItemParams) WithCallback(f InteractFunc) *ItemParams {
	ip.onUse = f
	return ip
}

func CreateNpc(name, description string) *NPC {
	return &NPC{Name: name, Description: description}
}

var lastRoomKey uint64 = 0

func CreateRoom(name, description string) *Room {
	r := &Room{Name: name, DescriptionData: description, Key: lastRoomKey}
	lastRoomKey++

	return r
}
