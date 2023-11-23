package collection

import (
	"encoding/xml"
)

// HEAD ...
type HEAD struct {
	COMPANYAttr string `xml:"COMPANY,attr,omitempty"`
	PROGRAMAttr string `xml:"PROGRAM,attr,omitempty"`
}

// LOCATION ...
type LOCATION struct {
	DIRAttr      string `xml:"DIR,attr,omitempty"`
	FILEAttr     string `xml:"FILE,attr,omitempty"`
	VOLUMEAttr   string `xml:"VOLUME,attr,omitempty"`
	VOLUMEIDAttr string `xml:"VOLUMEID,attr,omitempty"`
}

// ALBUM ...
type ALBUM struct {
	OFTRACKSAttr uint16 `xml:"OF_TRACKS,attr,omitempty"`
	TRACKAttr    uint32 `xml:"TRACK,attr,omitempty"`
	TITLEAttr    string `xml:"TITLE,attr,omitempty"`
}

// MODIFICATIONINFO ...
type MODIFICATIONINFO struct {
	XMLName        xml.Name `xml:"MODIFICATION_INFO"`
	AUTHORTYPEAttr string   `xml:"AUTHOR_TYPE,attr,omitempty"`
}

// INFO ...
type INFO struct {
	BITRATEAttr       uint32  `xml:"BITRATE,attr,omitempty"`
	GENREAttr         string  `xml:"GENRE,attr,omitempty"`
	LABELAttr         string  `xml:"LABEL,attr,omitempty"`
	COMMENTAttr       string  `xml:"COMMENT,attr,omitempty"`
	KEYAttr           string  `xml:"KEY,attr,omitempty"`
	PLAYCOUNTAttr     uint16  `xml:"PLAYCOUNT,attr,omitempty"`
	PLAYTIMEAttr      uint16  `xml:"PLAYTIME,attr,omitempty"`
	PLAYTIMEFLOATAttr float64 `xml:"PLAYTIME_FLOAT,attr,omitempty"`
	IMPORTDATEAttr    string  `xml:"IMPORT_DATE,attr,omitempty"`
	LASTPLAYEDAttr    string  `xml:"LAST_PLAYED,attr,omitempty"`
	RELEASEDATEAttr   string  `xml:"RELEASE_DATE,attr,omitempty"`
	FLAGSAttr         uint8   `xml:"FLAGS,attr,omitempty"`
	FILESIZEAttr      uint16  `xml:"FILESIZE,attr,omitempty"`
	COLORAttr         uint8   `xml:"COLOR,attr,omitempty"`
	COVERARTIDAttr    string  `xml:"COVERARTID,attr,omitempty"`
	RANKINGAttr       uint8   `xml:"RANKING,attr,omitempty"`
	PRODUCERAttr      string  `xml:"PRODUCER,attr,omitempty"`
	RATINGAttr        string  `xml:"RATING,attr,omitempty"`
	REMIXERAttr       string  `xml:"REMIXER,attr,omitempty"`
	KEYLYRICSAttr     string  `xml:"KEY_LYRICS,attr,omitempty"`
}

// TEMPO ...
type TEMPO struct {
	BPMAttr        float64 `xml:"BPM,attr,omitempty"`
	BPMQUALITYAttr float64 `xml:"BPM_QUALITY,attr,omitempty"`
}

// LOUDNESS ...
type LOUDNESS struct {
	PEAKDBAttr      float64 `xml:"PEAK_DB,attr,omitempty"`
	PERCEIVEDDBAttr float64 `xml:"PERCEIVED_DB,attr,omitempty"`
	ANALYZEDDBAttr  float64 `xml:"ANALYZED_DB,attr,omitempty"`
}

// LOOPINFO ...
type LOOPINFO struct {
	SAMPLETYPEINFOAttr uint8 `xml:"SAMPLE_TYPE_INFO,attr,omitempty"`
}

// MUSICALKEY ...
type MUSICALKEY struct {
	XMLName   xml.Name `xml:"MUSICAL_KEY"`
	VALUEAttr uint8    `xml:"VALUE,attr,omitempty"`
}

// CUEV2 ...
type CUEV2 struct {
	XMLName        xml.Name `xml:"CUE_V2"`
	NAMEAttr       string   `xml:"NAME,attr,omitempty"`
	DISPLORDERAttr uint8    `xml:"DISPL_ORDER,attr,omitempty"`
	TYPEAttr       uint8    `xml:"TYPE,attr,omitempty"`
	STARTAttr      float64  `xml:"START,attr,omitempty"`
	LENAttr        float64  `xml:"LEN,attr,omitempty"`
	REPEATSAttr    int8     `xml:"REPEATS,attr,omitempty"`
	HOTCUEAttr     int8     `xml:"HOTCUE,attr,omitempty"`
}

// STEMS ...
type STEMS struct {
	STEMSAttr string `xml:"STEMS,attr,omitempty"`
}

// ENTRY ...
type ENTRY struct {
	MODIFIEDDATEAttr         string             `xml:"MODIFIED_DATE,attr,omitempty"`
	MODIFIEDTIMEAttr         uint32             `xml:"MODIFIED_TIME,attr,omitempty"`
	AUDIOIDAttr              string             `xml:"AUDIO_ID,attr,omitempty"`
	TITLEAttr                string             `xml:"TITLE,attr,omitempty"`
	ARTISTAttr               string             `xml:"ARTIST,attr,omitempty"`
	LOCKAttr                 uint8              `xml:"LOCK,attr,omitempty"`
	LOCKMODIFICATIONTIMEAttr string             `xml:"LOCK_MODIFICATION_TIME,attr,omitempty"`
	LOCATION                 []LOCATION         `xml:"LOCATION"`
	ALBUM                    []ALBUM            `xml:"ALBUM"`
	MODIFICATIONINFO         []MODIFICATIONINFO `xml:"MODIFICATION_INFO"`
	INFO                     []INFO             `xml:"INFO"`
	TEMPO                    []TEMPO            `xml:"TEMPO"`
	LOUDNESS                 []LOUDNESS         `xml:"LOUDNESS"`
	LOOPINFO                 []LOOPINFO         `xml:"LOOPINFO"`
	MUSICALKEY               []MUSICALKEY       `xml:"MUSICAL_KEY"`
	CUEV2                    []CUEV2            `xml:"CUE_V2"`
	STEMS                    []STEMS            `xml:"STEMS"`
	PRIMARYKEY               *PRIMARYKEY        `xml:"PRIMARYKEY"`
}

// COLLECTION ...
type COLLECTION struct {
	ENTRIESAttr uint16  `xml:"ENTRIES,attr,omitempty"`
	ENTRY       []ENTRY `xml:"ENTRY"`
}

// CELL ...
type CELL struct {
	INDEXAttr       uint8   `xml:"INDEX,attr,omitempty"`
	CELLNAMEAttr    string  `xml:"CELLNAME,attr,omitempty"`
	COLORAttr       uint8   `xml:"COLOR,attr,omitempty"`
	SYNCAttr        uint8   `xml:"SYNC,attr,omitempty"`
	REVERSEAttr     uint8   `xml:"REVERSE,attr,omitempty"`
	MODEAttr        uint8   `xml:"MODE,attr,omitempty"`
	TYPEAttr        uint8   `xml:"TYPE,attr,omitempty"`
	SPEEDAttr       float64 `xml:"SPEED,attr,omitempty"`
	TRANSPOSEAttr   float64 `xml:"TRANSPOSE,attr,omitempty"`
	OFFSETAttr      float64 `xml:"OFFSET,attr,omitempty"`
	NUDGEAttr       float64 `xml:"NUDGE,attr,omitempty"`
	GAINAttr        float64 `xml:"GAIN,attr,omitempty"`
	STARTMARKERAttr float64 `xml:"START_MARKER,attr,omitempty"`
	ENDMARKERAttr   float64 `xml:"END_MARKER,attr,omitempty"`
	BPMAttr         float64 `xml:"BPM,attr,omitempty"`
	DIRAttr         string  `xml:"DIR,attr,omitempty"`
	FILEAttr        string  `xml:"FILE,attr,omitempty"`
	VOLUMEAttr      string  `xml:"VOLUME,attr,omitempty"`
}

// SLOT ...
type SLOT struct {
	KEYLOCKAttr         uint8  `xml:"KEYLOCK,attr,omitempty"`
	FXENABLEAttr        uint8  `xml:"FXENABLE,attr,omitempty"`
	PUNCHMODEAttr       uint8  `xml:"PUNCHMODE,attr,omitempty"`
	ACTIVECELLINDEXAttr uint8  `xml:"ACTIVE_CELL_INDEX,attr,omitempty"`
	CELL                []CELL `xml:"CELL"`
}

// SET ...
type SET struct {
	TITLEAttr        string           `xml:"TITLE,attr,omitempty"`
	ARTISTAttr       string           `xml:"ARTIST,attr,omitempty"`
	QUANTVAlUEAttr   uint8            `xml:"QUANT_VAlUE,attr,omitempty"`
	QUANTSTATEAttr   uint8            `xml:"QUANT_STATE,attr,omitempty"`
	LOCATION         LOCATION         `xml:"LOCATION"`
	ALBUM            ALBUM            `xml:"ALBUM"`
	MODIFICATIONINFO MODIFICATIONINFO `xml:"MODIFICATION_INFO"`
	INFO             INFO             `xml:"INFO"`
	TEMPO            TEMPO            `xml:"TEMPO"`
	SLOT             []SLOT           `xml:"SLOT"`
}

// SETS ...
type SETS struct {
	ENTRIESAttr uint8 `xml:"ENTRIES,attr,omitempty"`
	SET         []SET `xml:"SET"`
}

// PLAYLIST ...
type PLAYLIST struct {
	ENTRIESAttr string   `xml:"ENTRIES,attr,omitempty"`
	TYPEAttr    string   `xml:"TYPE,attr,omitempty"`
	UUIDAttr    string   `xml:"UUID,attr,omitempty"`
	ENTRY       []*ENTRY `xml:"ENTRY"`
}

// PRIMARYKEY ...
type PRIMARYKEY struct {
	TYPEAttr string `xml:"TYPE,attr,omitempty"`
	KEYAttr  string `xml:"KEY,attr,omitempty"`
}

// SEARCHEXPRESSION ...
type SEARCHEXPRESSION struct {
	XMLName     xml.Name `xml:"SEARCH_EXPRESSION"`
	VERSIONAttr uint8    `xml:"VERSION,attr,omitempty"`
	QUERYAttr   string   `xml:"QUERY,attr,omitempty"`
}

// SMARTLIST ...
type SMARTLIST struct {
	UUIDAttr         string           `xml:"UUID,attr,omitempty"`
	SEARCHEXPRESSION SEARCHEXPRESSION `xml:"SEARCH_EXPRESSION"`
}

// NODE ...
type NODE struct {
	TYPEAttr  string    `xml:"TYPE,attr,omitempty"`
	NAMEAttr  string    `xml:"NAME,attr,omitempty"`
	SMARTLIST SMARTLIST `xml:"SMARTLIST"`
	PLAYLIST  PLAYLIST  `xml:"PLAYLIST"`
	SUBNODES  SUBNODES  `xml:"SUBNODES"`
}

// SUBNODES ...
type SUBNODES struct {
	COUNTAttr uint8  `xml:"COUNT,attr,omitempty"`
	NODE      []NODE `xml:"NODE"`
}

// PLAYLISTS ...
type PLAYLISTS struct {
	NODE []NODE `xml:"NODE"`
}

// CRITERIA ...
type CRITERIA struct {
	ATTRIBUTEAttr uint8 `xml:"ATTRIBUTE,attr,omitempty"`
	DIRECTIONAttr uint8 `xml:"DIRECTION,attr,omitempty"`
}

// SORTINGINFO ...
type SORTINGINFO struct {
	XMLName  xml.Name `xml:"SORTING_INFO"`
	PATHAttr string   `xml:"PATH,attr,omitempty"`
	CRITERIA CRITERIA `xml:"CRITERIA"`
}

// INDEXING ...
type INDEXING struct {
	SORTINGINFO []SORTINGINFO `xml:"SORTING_INFO"`
}

// NML ...
type NML struct {
	VERSIONAttr uint8      `xml:"VERSION,attr,omitempty"`
	HEAD        HEAD       `xml:"HEAD"`
	COLLECTION  COLLECTION `xml:"COLLECTION"`
	SETS        SETS       `xml:"SETS"`
	PLAYLISTS   PLAYLISTS  `xml:"PLAYLISTS"`
	INDEXING    INDEXING   `xml:"INDEXING"`
}
