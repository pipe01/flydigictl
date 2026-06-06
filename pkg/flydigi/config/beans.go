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

func interpolateGradient(colors []LedUnit, t float32) LedUnit {
	if len(colors) == 0 {
		return LedUnit{}
	}
	if len(colors) == 1 {
		return colors[0]
	}

	pos := t * float32(len(colors)-1)
	i := int(pos)
	f := pos - float32(i)

	c1 := colors[i]
	c2 := colors[min(i+1, len(colors)-1)]

	return LedUnit{
		R: uint8(float32(c1.R)*(1-f) + float32(c2.R)*f),
		G: uint8(float32(c1.G)*(1-f) + float32(c2.G)*f),
		B: uint8(float32(c1.B)*(1-f) + float32(c2.B)*f),
	}
}

func (b *NewLedConfigBean) SetGradient(colors []LedUnit, speed float32) {
	const ledCount = 10
	const groupCount = 16

	b.LedMode = LedModeGradient
	b.Loop_time = 100 - byte(speed*100)
	b.Rgb_num = ledCount
	b.Loop_End = ledCount - 1
	b.Type = 0

	b.LedGroups = utils.RepeatFunc(func() *LedGroup {
		return &LedGroup{
			Units: utils.RepeatFunc(func() *LedUnit {
				return &LedUnit{}
			}, ledCount),
		}
	}, groupCount)

	for i := 0; i < ledCount; i++ {
		t := float32(i) / float32(ledCount-1)

		c := interpolateGradient(colors, t)

		for g := 0; g < groupCount; g++ {
			b.LedGroups[g].Units[i] = &c
		}
	}
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

func (t *Trigger) ApplyDefaultTrigger() {
	t.Type = 0
	t.AutoTrigger.Mode = 0
	t.AutoTrigger.VibrationBind.TriggerParams = []int32{1, 10, 1, 90, 0}

	for i := range t.AutoTrigger.MixedParams {
		t.AutoTrigger.MixedParams[i] = 0
	}

	lg := t.TriggerMotor.LineGear
	lg.Type = 1
	lg.Min = 30
	lg.Max = 80
	lg.Filter = 5
	lg.Vibrlimit = 1
	lg.Scale = 50
	lg.TimeLimit = 0
}

func (t *Trigger) ApplyRaceTrigger(initialPos int32, pressure int32) {
	t.Type = 0
	t.AutoTrigger.Mode = 1
	t.AutoTrigger.VibrationBind.TriggerParams = []int32{100, 1, 255, 70, 0}

	for i := range t.AutoTrigger.MixedParams {
		t.AutoTrigger.MixedParams[i] = 0
	}

	t.AutoTrigger.MixedParams[0] = int32(initialPos)
	t.AutoTrigger.MixedParams[1] = int32(pressure)

	lg := t.TriggerMotor.LineGear
	lg.Type = 1
	lg.Min = 90
	lg.Max = 100
	lg.Filter = 5
	lg.Vibrlimit = 1
	lg.Scale = 100
	lg.TimeLimit = 0
}

func (t *Trigger) ApplyRecoilTrigger(InitialPos int32, InitialStrength int32, Intensity int32, Frequency int32, OutputFromVibration int32) {
	t.Type = 0
	t.AutoTrigger.Mode = 2
	t.AutoTrigger.VibrationBind.Type = 0
	t.AutoTrigger.VibrationBind.MinFilter = 10
	t.AutoTrigger.VibrationBind.Scale = 50
	t.AutoTrigger.VibrationBind.TriggerParams =
		[]int32{100, 1, 255, 70, 0}

	t.AutoTrigger.MixedBorder = 0
	t.AutoTrigger.MixedParams = []int32{
		InitialPos, InitialStrength, Intensity, Frequency, OutputFromVibration,
		0, 0, 0, 0, 0,
	}

	lg := t.TriggerMotor.LineGear
	lg.Type = 1
	lg.Min = 30
	lg.Max = 80
	lg.Filter = 5
	lg.Vibrlimit = 1
	lg.Scale = 50
	lg.TimeLimit = 0
}

func (t *Trigger) ApplySniperTrigger(InitialPos int32, Length int32, Pressure int32, OutputFromVibration int32) {
	t.Type = 0
	t.AutoTrigger.Mode = 3
	t.AutoTrigger.VibrationBind.Type = 0
	t.AutoTrigger.VibrationBind.MinFilter = 10
	t.AutoTrigger.VibrationBind.Scale = 50
	t.AutoTrigger.VibrationBind.TriggerParams =
		[]int32{100, 1, 255, 70, 0}

	t.AutoTrigger.MixedBorder = 0
	t.AutoTrigger.MixedParams = []int32{
		InitialPos, Length, Pressure, 0, OutputFromVibration,
		0, 0, 0, 0, 0,
	}

	lg := t.TriggerMotor.LineGear
	lg.Type = 1
	lg.Min = 30
	lg.Max = 80
	lg.Filter = 5
	lg.Vibrlimit = 1
	lg.Scale = 50
	lg.TimeLimit = 0
}

func (t *Trigger) ApplyLockTrigger(InitialPos int32) {
	t.Type = 0
	t.AutoTrigger.Mode = 4
	t.AutoTrigger.VibrationBind.Type = 0
	t.AutoTrigger.VibrationBind.MinFilter = 10
	t.AutoTrigger.VibrationBind.Scale = 50
	t.AutoTrigger.VibrationBind.TriggerParams =
		[]int32{100, 1, 255, 70, 0}

	t.AutoTrigger.MixedBorder = 0
	t.AutoTrigger.MixedParams = []int32{
		InitialPos, 250, 1, 0, 0,
		0, 0, 0, 0, 0,
	}

	lg := t.TriggerMotor.LineGear
	lg.Type = 1
	lg.Min = 30
	lg.Max = 80
	lg.Filter = 5
	lg.Vibrlimit = 1
	lg.Scale = 50
	lg.TimeLimit = 0
}

func (t *Trigger) ApplyVibrationTrigger(Coefficient int32, Threshold int32, TravelRange int32, Frequency int32) {
	t.Type = 0
	t.AutoTrigger.Mode = 5
	t.AutoTrigger.VibrationBind.Type = 2
	t.AutoTrigger.VibrationBind.Scale = Coefficient
	t.AutoTrigger.VibrationBind.MinFilter = Threshold
	t.AutoTrigger.VibrationBind.TriggerParams =
		[]int32{TravelRange, 1, 1, Frequency, 0}

	t.AutoTrigger.MixedBorder = 0
	t.AutoTrigger.MixedParams = []int32{
		1, 1, 1, 90, 0,
		0, 0, 0, 0, 0,
	}

	lg := t.TriggerMotor.LineGear
	lg.Type = 1
	lg.Min = 30
	lg.Max = 80
	lg.Filter = 5
	lg.Vibrlimit = 1
	lg.Scale = 50
	lg.TimeLimit = 0
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
