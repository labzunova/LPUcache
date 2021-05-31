package main

import (
	"container/list"
	"errors"
	"fmt"
	"time"
	"unsafe"
)

type valueWithTTl struct {
	value string
	initTime time.Time
}

//var uintSize uint32 = 8
//var durationSize uint32 = 8
var initSize uintptr = 48 // каждый элемент изначальной структуры весит 8 байт
var mapElemSize uintptr = 16 // 8 ключ и 8 - указатель на слайс
var elemSize uintptr = 128 // 128 строка и 24 Time
//var timeSize uintptr = 24

type Cache struct {
	size uint32
	memory uintptr
	TTL time.Duration
	//values []valueWithTTl // для хранения значений в том порядке, в котором к ним обращались. Если индекс меньше => обращались раньше
	values *list.List // для хранения значений в том порядке, в котором к ним обращались. Если индекс меньше => обращались раньше
	valuesHashes map[uint32]*list.Element // для доступа к значению по ключу
	currentMemory uintptr
}

func NewCache(size, TTL, memory uint32) Cache {
	//return Cache{size, memory, time.Duration(TTL) * time.Second, make([]valueWithTTl, 0), make(map[uint32]*valueWithTTl),
	//	uintSize * 3 + durationSize +}

	return Cache{size, uintptr(memory), time.Duration(TTL) * time.Second, list.New(), make(map[uint32]*list.Element),
		initSize}
}

func (c *Cache) Put(i uint32, value string) error {
	newVal := &valueWithTTl{
		value: value,
		initTime: time.Now(),
	}
	//fmt.Println(unsafe.Sizeof(*newVal))
	//fmt.Println(unsafe.Sizeof(value))

	// если элемент с данным ключом уже есть в кэше
	elem, exists := c.valuesHashes[i]
	fmt.Println(unsafe.Sizeof(elem))
	if exists {
		c.values.MoveToFront(elem) // перемещаем в начало как самый "самый свежий"
		elem.Value = newVal // так как элементы List - интерфейс, необходимо преобразовать тип

		//c.values = c.values[:1] // убираем старый элемент // todo эффективно? GC? НЕ ПОСЛЕДНИЙ
		//c.values = append(c.values, []valueWithTTl{newVal}...) // вставляем в начало, так как элемент "свежее"
		//c.valuesHashes[i] = &c.values[0] // и обновляем ссылку в мапке
		return nil
	}
	
	// если элемента с данным ключом нет
	if c.currentMemory + elemSize + mapElemSize > c.memory {
		return errors.New("cache overflow")
	}

	memPlace := c.values.PushFront(newVal)
	c.valuesHashes[i] = memPlace
	c.currentMemory += elemSize + mapElemSize // размер элемента слайса + мапки
	fmt.Println("curr mem", c.currentMemory)
	//c.values = append(c.values, []valueWithTTl{newVal}...)
	//c.valuesHashes[i] = &c.values[0]

	// если кэш заполнен
	if c.values.Len() == int(c.size) {
		//удаляем элемент, который последний в очереди
		c.values.Back()

		//c.values = c.values[:1]
	}

	return nil
}

func (c *Cache) Get(i uint32) string {
	// если элемент с данным ключом уже есть в кэше
	existingValue, exist := c.valuesHashes[i]
	if exist {
		if time.Now().Sub(existingValue.Value.(*valueWithTTl).initTime) > c.TTL { // проверка на TTL
			return "not found"
		}

		// если элемент еще не "протух", то перемещаем в начало слайса и возвращаем его
		existingValue.Value.(*valueWithTTl).initTime = time.Now()
		c.values.MoveToFront(existingValue)

		//existingValue.initTime = time.Now()
		//c.values = append(c.values, []valueWithTTl{*existingValue}...)
		//// удаляем старый элемент:
		//c.values = append()
		//c.valuesHashes[i] = &c.values[0]

		return existingValue.Value.(*valueWithTTl).value
	}

	return "not found"
}
