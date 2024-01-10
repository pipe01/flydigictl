package config

import "github.com/pipe01/flydigictl/pkg/utils"

type AllConfigBean struct {
	Version        string
	PackageLength  int32
	Name           string
	Basic          *BasicBean
	KeyMapping     []*KeyMappingBean
	JoyMapping     *JoyMappingBean
	TriggerMapping *TriggerMappingBean
	MotionMapping  *MotionMappingBean
	Macro          *MacroGPBean
}

type BasicBean struct {
	LunPanMapping *LunPanMappingBean
	Motor         *MotorBean
	Led           *LedBean
	NewLedConfig  *NewLedConfigBean
}

type LunPanMappingBean struct {
	Rev  int32
	Type int32
}

type MotorBean struct {
	MainSwitch bool
	LeftMotor  *MotorNewBean
	RightMotor *MotorNewBean
}

type MotorNewBean struct {
	Type  int32
	Min   int32
	Max   int32
	Scale int32
}

type LedBean struct {
	Header    int32
	Mode      int32
	Peroid    int32
	Light     int32
	RgbColor0 []int32
	RgbColor1 []int32
}

type LedMode byte

const (
	LedModeOff LedMode = iota
	LedModeStreamlined
	LedModeBreathing
	LedModeGradient
	LedModeFeedback
	LedModeSteady
)

type NewLedConfigBean struct {
	Version     []byte
	Type        byte
	Loop_Start  byte
	Loop_End    byte
	Loop_time   byte // Speed, 0-100
	Light_scale byte // Brightness, 0-100
	Rgb_num     byte
	LedMode     LedMode
	Reserve     []byte
	LedGroups   []*LedGroup
}

func (b *NewLedConfigBean) SetSteady(color LedUnit) {
	b.LedMode = LedModeSteady
	b.Rgb_num = 5
	b.LedGroups = utils.RepeatFunc(func() *LedGroup { return &LedGroup{} }, 16)
	b.Loop_End = 0
	b.Type = 0

	for _, g := range b.LedGroups[:b.Rgb_num] {
		c := color

		g.Units = utils.RepeatFunc(func() *LedUnit { return &LedUnit{} }, 10)
		g.Units[0] = &c
	}
}

func (b *NewLedConfigBean) SetStreamlined(speed float32) {
	b.LedMode = LedModeStreamlined
	b.Loop_End = 5
	b.Loop_time = 100 - byte(speed*100)
	b.Rgb_num = 5
	b.LedGroups = getLedGroupList(0, 5)
}

type LedGroup struct {
	Units []*LedUnit
}

type LedUnit struct {
	R, G, B byte
}

type KeyMappingBean struct {
	MapKey            string
	MapMacroKeyList   []string
	OneMacroNewAction *OneMacroNewBean
	Turbo             int32
	TurboType         int32
	Priority          int32
	ListPriority      int32
	IsOwn             bool
	IsMap             bool
	IsFuncKey         bool
	ShowRes           string
	MapShowRes        string
	KeyPosType        KeyPosTypeEnum
	IsTriggerNative   bool
	TrigMode          int32
	MapData           string
	MapType           KeyMapType
	MenuFuncData      MenuFuncEnum
	Key               string
	KeyId             int32
	IsShowPic         bool
	IsShow            bool
}

type OneMacroNewBean struct {
	MacroId        int32
	LocalMacroId   int32
	GPID           int32
	Name           string
	FileName       string
	Btn            int32
	BtnNum         int32
	TotalTime      int32
	TotalActionNum int32
	TriggerType    int32
	ListStep       []*MacroNewActionBean
}

type MacroNewActionBean struct {
	ShowRes    string
	SelShowRes string
	KeyName    string
	IsFirst    bool
	Duration   int32
	Interval   int32
	UpDuration int32
	Btn        int32
	State      int32
	EventType  int32
}

type KeyPosTypeEnum int32

const (
	KeyPosTypeFront KeyPosTypeEnum = iota
	KeyPosTypeSide
	KeyPosTypeExt
	KeyPosTypeFn
	KeyPosTypeVirtual
)

type KeyMapType int32

const (
	KeyMapTypeGamePad KeyMapType = iota
	KeyMapTypeKeyboard
	KeyMapTypeMouse
	KeyMapTypeBurst
	KeyMapTypeVirtual
	KeyMapTypeMacro
	KeyMapTypeKeyboardBurst
	KeyMapTypeMouseBurst
)

type MenuFuncEnum int32

const (
	MenuFuncOpenApp MenuFuncEnum = iota
	MenuFuncADD_VOLUME
	MenuFuncSUB_VOLUME
	MenuFuncPLAY_NEXT
	MenuFuncPLAY_BEFORE
	MenuFuncADD_BRIGHT
	MenuFuncSUB_BRIGHT
	MenuFuncPLAY_PAUSE
	MenuFuncNULL
	MenuFuncXBOXGAMEBAR
	MenuFuncMOUSE_LEFTB_BUTTON
)

type JoyMappingBean struct {
	LeftJoystic  *JoyStickBean
	RightJoystic *JoyStickBean
}

type JoyStickBean struct {
	Curve  *Curve
	Config *JoyStickConfig // Unused
}

type Curve struct {
	Point0_X, Point0_Y int32
	Point1_X, Point1_Y int32

	End  int32
	Zero int32
	Type int32
}

type JoyStickConfig struct {
	OpenMethod any
	MapType    int32
	MapData    any
	Sens       float64
	SensX      float64
	SensY      float64
	DeadZone   int32
	DirectAlg  any
	DirectNum  int32
}

type TriggerMappingBean struct {
	MainSwitch   bool
	Type         int32
	LeftTrigger  *Trigger
	RightTrigger *Trigger
}

type Trigger struct {
	Type         int32
	Curve        *Curve
	AutoTrigger  *Autotrigger
	TriggerMotor *TriggerMotor
}

type Autotrigger struct {
	Mode          int32
	VibrationBind *Vibrationbind
	MixedBorder   int32
	MixedParams   []int32
}

type Vibrationbind struct {
	Type          int32
	MinFilter     int32
	Scale         int32
	TriggerParams []int32
}

type TriggerMotor struct {
	LineGear, MicrGear *TriggerMotorSet
}

type TriggerMotorSet struct {
	Type      int32
	Min       int32
	Max       int32
	Filter    int32
	Vibrlimit int32
	Scale     int32
	TimeLimit int32
}

type MotionMappingBean struct {
	MapType       int32
	OpenKeyId     int32
	OpenMapMethod int32
	OpenKey       string
	DeadZero      int32
	Zero          int32
	OpenKeyExtId  int32
	OpenKeyExt    string
	OpenKeyRes    string
	OpenKeyExtRes string
	Mode          int32
	Sensity       float32
	SensityX      float32
	SensityY      float32
	SensX         float32
	SensY         float32
	SimType       int32
}

type MacroGPBean struct {
	Nums       int32
	ListOffset []int32
	ListMacro  []*OneMacroBeanGP
}

type OneMacroBeanGP struct {
	Btn         int32
	Count_l     byte
	Count_h     byte
	TriggerType int32
	ListStep    []*MacroActionBeanGP
}

type MacroActionBeanGP struct {
	TriggerTime int32
	Btn         int32
	State       int32
}

type GamePadKeyConfig struct {
	KeyId      int32
	KeyName    string
	KeyPosType KeyPosTypeEnum
	ShowRes    string
	MapShowRes string
	IsShowPic  bool
	IsShow     bool
}
