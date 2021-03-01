package auto

import (
	"bitbucket.org/funplus/server/excel"
	"bitbucket.org/funplus/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var tokenEntries *TokenEntries //Token.xlsx全局变量

// Token.xlsx属性表
type TokenEntry struct {
	Id      int32 `json:"Id,omitempty"`      // 主键
	MaxHold int32 `json:"MaxHold,omitempty"` //品质
}

// Token.xlsx属性表集合
type TokenEntries struct {
	Rows map[int32]*TokenEntry `json:"Rows,omitempty"` //
}

func init() {
	excel.AddEntryLoader("Token.xlsx", (*TokenEntries)(nil))
}

func (e *TokenEntries) Load(excelFileRaw *excel.ExcelFileRaw) error {

	tokenEntries = &TokenEntries{
		Rows: make(map[int32]*TokenEntry, 100),
	}

	for _, v := range excelFileRaw.CellData {
		entry := &TokenEntry{}
		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return err
		}

		tokenEntries.Rows[entry.Id] = entry
	}

	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return nil

}

func GetTokenEntry(id int32) (*TokenEntry, bool) {
	entry, ok := tokenEntries.Rows[id]
	return entry, ok
}

func GetTokenSize() int32 {
	return int32(len(tokenEntries.Rows))
}

func GetTokenRows() map[int32]*TokenEntry {
	return tokenEntries.Rows
}
