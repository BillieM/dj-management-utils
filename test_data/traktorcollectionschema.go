package collection

import (
	"encoding/xml"
)

// HEAD ...
type HEAD struct {
	COMPANYAttr string `xml:"COMPANY,attr,omitempty"`
	PROGRAMAttr string `xml:"PROGRAM,attr,omitempty"`
	Value       string `xml:",chardata"`
}

// LOCATION ...
type LOCATION struct {
	DIRAttr      string `xml:"DIR,attr,omitempty"`
	FILEAttr     string `xml:"FILE,attr,omitempty"`
	VOLUMEAttr   string `xml:"VOLUME,attr,omitempty"`
	VOLUMEIDAttr string `xml:"VOLUMEID,attr,omitempty"`
	Value        string `xml:",chardata"`
}

// ALBUM ...
type ALBUM struct {
	OFTRACKSAttr int16  `xml:"OF_TRACKS,attr,omitempty"`
	TRACKAttr    int32  `xml:"TRACK,attr,omitempty"`
	TITLEAttr    string `xml:"TITLE,attr,omitempty"`
	Value        string `xml:",chardata"`
}

// MODIFICATIONINFO ...
type MODIFICATIONINFO struct {
	XMLName        xml.Name `xml:"MODIFICATION_INFO"`
	AUTHORTYPEAttr string   `xml:"AUTHOR_TYPE,attr,omitempty"`
	Value          string   `xml:",chardata"`
}

// INFO ...
type INFO struct {
	BITRATEAttr       int     `xml:"BITRATE,attr,omitempty"`
	GENREAttr         string  `xml:"GENRE,attr,omitempty"`
	COMMENTAttr       string  `xml:"COMMENT,attr,omitempty"`
	PLAYTIMEAttr      int16   `xml:"PLAYTIME,attr,omitempty"`
	PLAYTIMEFLOATAttr float32 `xml:"PLAYTIME_FLOAT,attr,omitempty"`
	IMPORTDATEAttr    string  `xml:"IMPORT_DATE,attr,omitempty"`
	RELEASEDATEAttr   string  `xml:"RELEASE_DATE,attr,omitempty"`
	FLAGSAttr         int8    `xml:"FLAGS,attr,omitempty"`
	FILESIZEAttr      int32   `xml:"FILESIZE,attr,omitempty"`
	LABELAttr         string  `xml:"LABEL,attr,omitempty"`
	COVERARTIDAttr    string  `xml:"COVERARTID,attr,omitempty"`
	KEYAttr           string  `xml:"KEY,attr,omitempty"`
	Value             string  `xml:",chardata"`
}

// TEMPO ...
type TEMPO struct {
	BPMAttr        float32 `xml:"BPM,attr,omitempty"`
	BPMQUALITYAttr float32 `xml:"BPM_QUALITY,attr,omitempty"`
	Value          string  `xml:",chardata"`
}

// LOUDNESS ...
type LOUDNESS struct {
	PEAKDBAttr      float32 `xml:"PEAK_DB,attr,omitempty"`
	PERCEIVEDDBAttr float32 `xml:"PERCEIVED_DB,attr,omitempty"`
	ANALYZEDDBAttr  float32 `xml:"ANALYZED_DB,attr,omitempty"`
	Value           string  `xml:",chardata"`
}

// MUSICALKEY ...
type MUSICALKEY struct {
	XMLName   xml.Name `xml:"MUSICAL_KEY"`
	VALUEAttr int8     `xml:"VALUE,attr,omitempty"`
	Value     string   `xml:",chardata"`
}

// CUEV2 ...
type CUEV2 struct {
	XMLName        xml.Name `xml:"CUE_V2"`
	NAMEAttr       string   `xml:"NAME,attr,omitempty"`
	DISPLORDERAttr int8     `xml:"DISPL_ORDER,attr,omitempty"`
	TYPEAttr       int8     `xml:"TYPE,attr,omitempty"`
	STARTAttr      float32  `xml:"START,attr,omitempty"`
	LENAttr        float32  `xml:"LEN,att,omitemptyr"`
	REPEATSAttr    int8     `xml:"REPEATS,attr,omitempty"`
	HOTCUEAttr     int8     `xml:"HOTCUE,attr,omitempty"`
	Value          string   `xml:",chardata"`
}

// STEMS ...
type STEMS struct {
	STEMSAttr string `xml:"STEMS,attr,omitempty"`
	Value     string `xml:",chardata"`
}

// ENTRY ...
type ENTRY struct {
	MODIFIEDDATEAttr string            `xml:"MODIFIED_DATE,attr,omitempty"`
	MODIFIEDTIMEAttr int               `xml:"MODIFIED_TIME,attr,omitempty"`
	AUDIOIDAttr      string            `xml:"AUDIO_ID,attr,omitempty"`
	TITLEAttr        string            `xml:"TITLE,attr,omitempty"`
	ARTISTAttr       string            `xml:"ARTIST,attr,omitempty"`
	LOCATION         *LOCATION         `xml:"LOCATION"`
	ALBUM            *ALBUM            `xml:"ALBUM"`
	MODIFICATIONINFO *MODIFICATIONINFO `xml:"MODIFICATION_INFO"`
	INFO             *INFO             `xml:"INFO"`
	TEMPO            *TEMPO            `xml:"TEMPO"`
	LOUDNESS         *LOUDNESS         `xml:"LOUDNESS"`
	MUSICALKEY       *MUSICALKEY       `xml:"MUSICAL_KEY"`
	CUEV2            *CUEV2            `xml:"CUE_V2"`
	STEMS            *STEMS            `xml:"STEMS"`
	PRIMARYKEY       *PRIMARYKEY       `xml:"PRIMARYKEY"`
}

// COLLECTION ...
type COLLECTION struct {
	ENTRIESAttr int32    `xml:"ENTRIES,attr,omitempty"`
	ENTRY       []*ENTRY `xml:"ENTRY"`
}

// CELL ...
type CELL struct {
	INDEXAttr       int8    `xml:"INDEX,attr,omitempty"`
	CELLNAMEAttr    string  `xml:"CELLNAME,attr,omitempty"`
	COLORAttr       int8    `xml:"COLOR,attr,omitempty"`
	SYNCAttr        int8    `xml:"SYNC,attr,omitempty"`
	REVERSEAttr     int8    `xml:"REVERSE,attr,omitempty"`
	MODEAttr        int8    `xml:"MODE,attr,omitempty"`
	TYPEAttr        int8    `xml:"TYPE,attr,omitempty"`
	SPEEDAttr       float32 `xml:"SPEED,attr,omitempty"`
	TRANSPOSEAttr   float32 `xml:"TRANSPOSE,attr,omitempty"`
	OFFSETAttr      float32 `xml:"OFFSET,attr,omitempty"`
	NUDGEAttr       float32 `xml:"NUDGE,attr,omitempty"`
	GAINAttr        float32 `xml:"GAIN,attr,omitempty"`
	STARTMARKERAttr float32 `xml:"START_MARKER,attr,omitempty"`
	ENDMARKERAttr   float32 `xml:"END_MARKER,attr,omitempty"`
	BPMAttr         float32 `xml:"BPM,attr,omitempty"`
	DIRAttr         string  `xml:"DIR,attr,omitempty"`
	FILEAttr        string  `xml:"FILE,attr,omitempty"`
	VOLUMEAttr      string  `xml:"VOLUME,attr,omitempty"`
	Value           string  `xml:",chardata"`
}

// SLOT ...
type SLOT struct {
	KEYLOCKAttr         int8    `xml:"KEYLOCK,attr"`
	FXENABLEAttr        int8    `xml:"FXENABLE,attr"`
	PUNCHMODEAttr       int8    `xml:"PUNCHMODE,attr"`
	ACTIVECELLINDEXAttr int8    `xml:"ACTIVE_CELL_INDEX,attr"`
	CELL                []*CELL `xml:"CELL"`
}

// SET ...
type SET struct {
	TITLEAttr        string            `xml:"TITLE,attr,omitempty"`
	ARTISTAttr       string            `xml:"ARTIST,attr,omitempty"`
	QUANTVAlUEAttr   int8              `xml:"QUANT_VAlUE,attr,omitempty"`
	QUANTSTATEAttr   int8              `xml:"QUANT_STATE,attr,omitempty"`
	LOCATION         *LOCATION         `xml:"LOCATION"`
	MODIFICATIONINFO *MODIFICATIONINFO `xml:"MODIFICATION_INFO"`
	INFO             *INFO             `xml:"INFO"`
	TEMPO            *TEMPO            `xml:"TEMPO"`
	SLOT             []*SLOT           `xml:"SLOT"`
}

// SETS ...
type SETS struct {
	ENTRIESAttr int16 `xml:"ENTRIES,attr,omitempty"`
	SET         *SET  `xml:"SET"`
}

// PRIMARYKEY ...
type PRIMARYKEY struct {
	TYPEAttr string `xml:"TYPE,attr,omitempty"`
	KEYAttr  string `xml:"KEY,attr,omitempty"`
	Value    string `xml:",chardata"`
}

// PLAYLIST ...
type PLAYLIST struct {
	ENTRIESAttr int16    `xml:"ENTRIES,attr,omitempty"`
	TYPEAttr    string   `xml:"TYPE,attr,omitempty"`
	UUIDAttr    string   `xml:"UUID,attr,omitempty"`
	ENTRY       []*ENTRY `xml:"ENTRY"`
}

type SMARTLIST struct {
	UUIDAttr          string             `xml:"UUID,attr,omitempty"`
	SEARCH_EXPRESSION *SEARCH_EXPRESSION `xml:"SEARCH_EXPRESSION"`
}

type SEARCH_EXPRESSION struct {
	VERSIONAttr int8   `xml:"VERSION,attr,omitempty"`
	QUERYAttr   string `xml:"QUERY,attr,omitempty"`
}

// NODE ...
type NODE struct {
	TYPEAttr  string     `xml:"TYPE,attr,omitempty"`
	NAMEAttr  string     `xml:"NAME,attr,omitempty"`
	PLAYLIST  *PLAYLIST  `xml:"PLAYLIST"`
	SMARTLIST *SMARTLIST `xml:"SMARTLIST"`
	SUBNODES  *SUBNODES  `xml:"SUBNODES"`
}

// SUBNODES ...
type SUBNODES struct {
	COUNTAttr int16   `xml:"COUNT,attr,omitempty"`
	NODE      []*NODE `xml:"NODE"`
}

// PLAYLISTS ...
type PLAYLISTS struct {
	NODE []*NODE `xml:"NODE"`
}

// CRITERIA ...
type CRITERIA struct {
	ATTRIBUTEAttr int8   `xml:"ATTRIBUTE,attr,omitempty"`
	DIRECTIONAttr int8   `xml:"DIRECTION,attr,omitempty"`
	Value         string `xml:",chardata"`
}

// SORTINGINFO ...
type SORTINGINFO struct {
	XMLName  xml.Name  `xml:"SORTING_INFO"`
	PATHAttr string    `xml:"PATH,attr,omitempty"`
	CRITERIA *CRITERIA `xml:"CRITERIA"`
}

// INDEXING ...
type INDEXING struct {
	SORTINGINFO []*SORTINGINFO `xml:"SORTING_INFO"`
}

// NML ...
type NML struct {
	VERSIONAttr int8        `xml:"VERSION,attr,omitempty"`
	HEAD        *HEAD       `xml:"HEAD"`
	COLLECTION  *COLLECTION `xml:"COLLECTION"`
	SETS        *SETS       `xml:"SETS"`
	PLAYLISTS   *PLAYLISTS  `xml:"PLAYLISTS"`
	INDEXING    *INDEXING   `xml:"INDEXING"`
}
