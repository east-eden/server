package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

var battlePropertyEntries *BattlePropertyEntries //BattleProperty.json全局变量

// 属性表 角色战斗属性-拳皇/battleProperty
type BattlePropertyEntry struct {
	Id                  int32  `json:"id,omitempty" bson:"id"`                                   //##属性id
	Level               int32  `json:"level,omitempty" bson:"level"`                             //角色等级左侧区间
	MAXHP               string `json:"mAXHP,omitempty" bson:"mAXHP"`                             //血量上限
	ATTACK              string `json:"aTTACK,omitempty" bson:"aTTACK"`                           //攻击力
	DEFENCE             string `json:"dEFENCE,omitempty" bson:"dEFENCE"`                         //防御力
	SPINITRECOVER       string `json:"sPINITRECOVER,omitempty" bson:"sPINITRECOVER"`             //初始怒气回复
	MAXSP               string `json:"mAXSP,omitempty" bson:"mAXSP"`                             //最大怒气
	SPRESERVES          string `json:"sPRESERVES,omitempty" bson:"sPRESERVES"`                   //怒气储量
	ATTACKRANGE         string `json:"aTTACKRANGE,omitempty" bson:"aTTACKRANGE"`                 //攻击距离
	ATTACKSPEED         string `json:"aTTACKSPEED,omitempty" bson:"aTTACKSPEED"`                 //攻速
	MOVESPEED           string `json:"mOVESPEED,omitempty" bson:"mOVESPEED"`                     //移动速度
	DODGE               string `json:"dODGE,omitempty" bson:"dODGE"`                             //闪避
	IGNDODGE            string `json:"iGNDODGE,omitempty" bson:"iGNDODGE"`                       //闪避忽视
	BLASTS              string `json:"bLASTS,omitempty" bson:"bLASTS"`                           //暴击上限
	STIFFNESS           string `json:"sTIFFNESS,omitempty" bson:"sTIFFNESS"`                     //硬直
	HITRATE             string `json:"hITRATE,omitempty" bson:"hITRATE"`                         //命中率
	ROLEFIRSTRATE       string `json:"rOLEFIRSTRATE,omitempty" bson:"rOLEFIRSTRATE"`             //角色先手值
	PENETRATE           string `json:"pENETRATE,omitempty" bson:"pENETRATE"`                     //穿透率
	CRITAL              string `json:"cRITAL,omitempty" bson:"cRITAL"`                           //暴击率
	IGNCRITAL           string `json:"iGNCRITAL,omitempty" bson:"iGNCRITAL"`                     //免暴率
	IGNBLASTS           string `json:"iGNBLASTS,omitempty" bson:"iGNBLASTS"`                     //暴击抵抗
	PARRY               string `json:"pARRY,omitempty" bson:"pARRY"`                             //格挡
	BREAKPARRY          string `json:"bREAKPARRY,omitempty" bson:"bREAKPARRY"`                   //格挡抵抗
	STIFFNESSDEF        string `json:"sTIFFNESSDEF,omitempty" bson:"sTIFFNESSDEF"`               //硬直抵抗
	ENDURE              string `json:"eNDURE,omitempty" bson:"eNDURE"`                           //霸体免控率
	DEFSTUN             string `json:"dEFSTUN,omitempty" bson:"dEFSTUN"`                         //眩晕抵抗
	DEFFREEZE           string `json:"dEFFREEZE,omitempty" bson:"dEFFREEZE"`                     //冰冻抵抗
	NEARADD             string `json:"nEARADD,omitempty" bson:"nEARADD"`                         //近战增伤
	NEARSUB             string `json:"nEARSUB,omitempty" bson:"nEARSUB"`                         //近战减伤
	FARADD              string `json:"fARADD,omitempty" bson:"fARADD"`                           //远程增伤
	FARSUB              string `json:"fARSUB,omitempty" bson:"fARSUB"`                           //远程减伤
	TRICKADD            string `json:"tRICKADD,omitempty" bson:"tRICKADD"`                       //大招增伤
	DEFTRICKADD         string `json:"dEFTRICKADD,omitempty" bson:"dEFTRICKADD"`                 //大招减伤
	DAMAGEADD           string `json:"dAMAGEADD,omitempty" bson:"dAMAGEADD"`                     //最终增伤
	DAMAGESUB           string `json:"dAMAGESUB,omitempty" bson:"dAMAGESUB"`                     //最终减伤
	SPINCREASE          string `json:"sPINCREASE,omitempty" bson:"sPINCREASE"`                   //怒气增加
	SPRECOVERY          string `json:"sPRECOVERY,omitempty" bson:"sPRECOVERY"`                   //怒气回复
	SPINCRERESPOSTIVE   string `json:"sPINCRERESPOSTIVE,omitempty" bson:"sPINCRERESPOSTIVE"`     //攻击回怒
	SPINCRERESNEGIATIVE string `json:"sPINCRERESNEGIATIVE,omitempty" bson:"sPINCRERESNEGIATIVE"` //受击回怒
	SPINCRERESNEFAILED  string `json:"sPINCRERESNEFAILED,omitempty" bson:"sPINCRERESNEFAILED"`   //战败回怒
	SPRETENTION         string `json:"sPRETENTION,omitempty" bson:"sPRETENTION"`                 //怒气截留
	SPDESTRUCTION       string `json:"sPDESTRUCTION,omitempty" bson:"sPDESTRUCTION"`             //自爆减怒
	EFFECTHEAL          string `json:"eFFECTHEAL,omitempty" bson:"eFFECTHEAL"`                   //被治疗效果
	SPINITIAL           string `json:"sPINITIAL,omitempty" bson:"sPINITIAL"`                     //初始怒气
	SP                  string `json:"sP,omitempty" bson:"sP"`                                   //怒气
	SPINCRERESADD       string `json:"sPINCRERESADD,omitempty" bson:"sPINCRERESADD"`             //回怒增幅
	SPINCRERESSUB       string `json:"sPINCRERESSUB,omitempty" bson:"sPINCRERESSUB"`             //回怒降幅
	SPKILLRECOVER       string `json:"sPKILLRECOVER,omitempty" bson:"sPKILLRECOVER"`             //击杀回怒
	SILENTDEFENCE       string `json:"sILENTDEFENCE,omitempty" bson:"sILENTDEFENCE"`             //沉默抵抗
	SKILLTRIGGERADD     string `json:"sKILLTRIGGERADD,omitempty" bson:"sKILLTRIGGERADD"`         //技能触发率增长
	DEFSKILLTRIGGERADD  string `json:"dEFSKILLTRIGGERADD,omitempty" bson:"dEFSKILLTRIGGERADD"`   //技能触发率抵抗
	HEALADD             string `json:"hEALADD,omitempty" bson:"hEALADD"`                         //治疗增幅
	COUNTERING          string `json:"cOUNTERING,omitempty" bson:"cOUNTERING"`                   //反击率
	DEFCOUNTER          string `json:"dEFCOUNTER,omitempty" bson:"dEFCOUNTER"`                   //反击抵抗
	BLOODSUCK           string `json:"bLOODSUCK,omitempty" bson:"bLOODSUCK"`                     //吸血系数
	DEFBLOODSUCK        string `json:"dEFBLOODSUCK,omitempty" bson:"dEFBLOODSUCK"`               //吸血抗性
	COMBO               string `json:"cOMBO,omitempty" bson:"cOMBO"`                             //连击
	DEFCOMBO            string `json:"dEFCOMBO,omitempty" bson:"dEFCOMBO"`                       //连击抗性
	ENDUREBREAK         string `json:"eNDUREBREAK,omitempty" bson:"eNDUREBREAK"`                 //霸体免控抵抗
	PHYSICALADD         string `json:"pHYSICALADD,omitempty" bson:"pHYSICALADD"`                 //物理增伤
	PHYSICALSUB         string `json:"pHYSICALSUB,omitempty" bson:"pHYSICALSUB"`                 //物理减伤
	MAGICADD            string `json:"mAGICADD,omitempty" bson:"mAGICADD"`                       //魔法增伤
	MAGICSUB            string `json:"mAGICSUB,omitempty" bson:"mAGICSUB"`                       //魔法减伤
	COUNTERADD          string `json:"cOUNTERADD,omitempty" bson:"cOUNTERADD"`                   //反击增伤
	HEALSUB             string `json:"hEALSUB,omitempty" bson:"hEALSUB"`                         //治疗抑制
	ALLDEFENCEADD       string `json:"aLLDEFENCEADD,omitempty" bson:"aLLDEFENCEADD"`             //全状态抗性
	ALLSTATEADD         string `json:"aLLSTATEADD,omitempty" bson:"aLLSTATEADD"`                 //全状态属性
	MAXDEFENCEADD       string `json:"mAXDEFENCEADD,omitempty" bson:"mAXDEFENCEADD"`             //最高状态抗性
	MAXSTATEADD         string `json:"mAXSTATEADD,omitempty" bson:"mAXSTATEADD"`                 //最高状态属性
	STIFFNESSSTRENG     string `json:"sTIFFNESSSTRENG,omitempty" bson:"sTIFFNESSSTRENG"`         //硬直增强
	STIFFNESSWEAKEN     string `json:"sTIFFNESSWEAKEN,omitempty" bson:"sTIFFNESSWEAKEN"`         //硬直削弱
	SUCKSTRENGTH        string `json:"sUCKSTRENGTH,omitempty" bson:"sUCKSTRENGTH"`               //吸血增幅
	MAXSTATEADDEFFECT   string `json:"mAXSTATEADDEFFECT,omitempty" bson:"mAXSTATEADDEFFECT"`     //最高进阶效果增幅
	COMBOADD            string `json:"cOMBOADD,omitempty" bson:"cOMBOADD"`                       //连击增幅
	PARRYADD            string `json:"pARRYADD,omitempty" bson:"pARRYADD"`                       //格挡免伤增幅
	HEALCRITAL          string `json:"hEALCRITAL,omitempty" bson:"hEALCRITAL"`                   //治疗暴击
	HEALBLASTS          string `json:"hEALBLASTS,omitempty" bson:"hEALBLASTS"`                   //治疗暴击增幅
	DEFENCESTENGTH      string `json:"dEFENCESTENGTH,omitempty" bson:"dEFENCESTENGTH"`           //防御免伤增量77
	HIDEDAMAGEDD        string `json:"hIDEDAMAGEDD,omitempty" bson:"hIDEDAMAGEDD"`               //隐藏增伤78
	HEALDEFENCE         string `json:"hEALDEFENCE,omitempty" bson:"hEALDEFENCE"`                 //被治疗效果降低
	DEFENCEWEAKNESS     string `json:"dEFENCEWEAKNESS,omitempty" bson:"dEFENCEWEAKNESS"`         //防御减免
	ATTACKSUBSP         string `json:"aTTACKSUBSP,omitempty" bson:"aTTACKSUBSP"`                 //攻击减怒
	BEHURTSUBSP         string `json:"bEHURTSUBSP,omitempty" bson:"bEHURTSUBSP"`                 //受击减怒
	PVPATTACKFACTOR     string `json:"pVPATTACKFACTOR,omitempty" bson:"pVPATTACKFACTOR"`         //pvp攻击百分比
	PVPDEFENCEFACTOR    string `json:"pVPDEFENCEFACTOR,omitempty" bson:"pVPDEFENCEFACTOR"`       //pvp防御百分比
	PVPHPFACTOR         string `json:"pVPHPFACTOR,omitempty" bson:"pVPHPFACTOR"`                 //pvp血量百分比
	SKILLDMGADD         string `json:"sKILLDMGADD,omitempty" bson:"sKILLDMGADD"`                 //技能增伤
	SKILLDMGSUB         string `json:"sKILLDMGSUB,omitempty" bson:"sKILLDMGSUB"`                 //技能减伤
	MUILTIDAMAGEADD     string `json:"mUILTIDAMAGEADD,omitempty" bson:"mUILTIDAMAGEADD"`         //群体伤害加成
	SINGLEDAMAGEADD     string `json:"sINGLEDAMAGEADD,omitempty" bson:"sINGLEDAMAGEADD"`         //单体伤害加成
	MUILTIDAMAGESUB     string `json:"mUILTIDAMAGESUB,omitempty" bson:"mUILTIDAMAGESUB"`         //群体伤害减益
	SINGLEDAMAGESUB     string `json:"sINGLEDAMAGESUB,omitempty" bson:"sINGLEDAMAGESUB"`         //单体伤害减益
	MONSTERDMGADD       string `json:"mONSTERDMGADD,omitempty" bson:"mONSTERDMGADD"`             //非Boss伤害加成
	BOSSDMGADD          string `json:"bOSSDMGADD,omitempty" bson:"bOSSDMGADD"`                   //Boss伤害加成
	CLEANSESTRENGTH     string `json:"cLEANSESTRENGTH,omitempty" bson:"cLEANSESTRENGTH"`         //净化增强
	CLEANSEDEFENCE      string `json:"cLEANSEDEFENCE,omitempty" bson:"cLEANSEDEFENCE"`           //净化抑制
}

// BattleProperty.json属性表集合
type BattlePropertyEntries struct {
	Rows   map[string]*BattlePropertyEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*BattlePropertyEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("BattleProperty.json", (*BattlePropertyEntries)(nil))
}

func (e *BattlePropertyEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	battlePropertyEntries = &BattlePropertyEntries{
		Rows:   make(map[string]*BattlePropertyEntry),
		Arrays: make([]*BattlePropertyEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &BattlePropertyEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return battlePropertyEntries, err
		}

		key := fmt.Sprintf("%s_%s", cast.ToString(entry.Id), cast.ToString(entry.Level))
		battlePropertyEntries.Rows[key] = entry
		battlePropertyEntries.Arrays = append(battlePropertyEntries.Arrays, entry)
	}

	sort.Slice(battlePropertyEntries.Arrays, func(i, j int) bool {
		return battlePropertyEntries.Arrays[i].Id < battlePropertyEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("battleProperty.json", battlePropertyEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return battlePropertyEntries, nil

}

func (e *BattlePropertyEntries) GetArrays() any {
	return battlePropertyEntries.Arrays
}

func GetBattlePropertyEntry(keys ...string) (*BattlePropertyEntry, bool) {
	keyName := make([]string, 0, len(keys))
	keyName = append(keyName, keys...)

	finalKey := strings.Join(keyName, "_")
	entry, ok := battlePropertyEntries.Rows[finalKey]
	return entry, ok
}

func GetBattlePropertySize() int {
	return int(len(battlePropertyEntries.Rows))
}

func GetBattlePropertyRows() map[string]*BattlePropertyEntry {
	return battlePropertyEntries.Rows
}

func GetBattlePropertyArrays() []*BattlePropertyEntry {
	return battlePropertyEntries.Arrays
}
