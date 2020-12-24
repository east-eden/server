package excel

import (
)

var	heroEntries	*HeroEntries	//hero.xlsx全局变量

// 英雄属性表
type HeroEntry struct {
	AttID     	int       	`json:"AttID,omitempty"`	//属性id      
	AttList   	[]int     	`json:"AttList,omitempty"`	//属性列表      
	ID        	int       	`json:"Id"`	//id        
	Name      	string    	`json:"Name,omitempty"`	//名字        
	Quality   	int       	`json:"Quality,omitempty"`	//品质        
}

// 英雄属性表集合
type HeroEntries struct {
	Rows      	map[int]*HeroEntry	`json:"Rows"`	//          
}

func  init()  {
	AddEntries(heroEntries, "hero.xlsx")
}

func (e *HeroEntries) Load() error {
	return nil
}

