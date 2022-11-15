package builder

// 供多种固定模式选择的生成器
type House struct {
	WindowType string
	DoorType   string
	Floor      int
}

// ---------------房屋建造器
type IBuilder interface {
	setWindowType()
	setDoorType()
	setNumFloor()
	getHouse() House
}

func GetBuilder(builderType string) IBuilder {
	if builderType == "normal" {
		return NewNormalBuilder()
	}

	if builderType == "igloo" {
		return NewIglooBuilder()
	}
	return nil
}

// ------------------NormalBuilder普通房屋生成器
type NormalBuilder struct {
	House
}

func NewNormalBuilder() *NormalBuilder {
	return &NormalBuilder{}
}

func (b *NormalBuilder) setWindowType() {
	b.WindowType = "Wooden Window"
}

func (b *NormalBuilder) setDoorType() {
	b.DoorType = "Wooden Door"
}

func (b *NormalBuilder) setNumFloor() {
	b.Floor = 2
}

func (b *NormalBuilder) getHouse() House {
	return House{
		DoorType:   b.DoorType,
		WindowType: b.WindowType,
		Floor:      b.Floor,
	}
}

// ------------------IglooBuilder冰屋生成器
type IglooBuilder struct {
	House
}

func NewIglooBuilder() *IglooBuilder {
	return &IglooBuilder{}
}

func (b *IglooBuilder) setWindowType() {
	b.WindowType = "Snow Window"
}

func (b *IglooBuilder) setDoorType() {
	b.DoorType = "Snow Door"
}

func (b *IglooBuilder) setNumFloor() {
	b.Floor = 1
}

func (b *IglooBuilder) getHouse() House {
	return House{
		DoorType:   b.DoorType,
		WindowType: b.WindowType,
		Floor:      b.Floor,
	}
}

// -----------------执行者对建造过程进行组织
type Director struct {
	builder IBuilder
}

func NewDirector(b IBuilder) *Director {
	return &Director{
		builder: b,
	}
}

func (d *Director) SetBuilder(b IBuilder) {
	d.builder = b
}

func (d *Director) BuildHouse() House {
	d.builder.setDoorType()
	d.builder.setWindowType()
	d.builder.setNumFloor()
	return d.builder.getHouse()
}
