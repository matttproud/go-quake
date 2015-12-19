package prog

import (
	"io"

	. "github.com/matttproud/go-quake/qtype"
)

type GlobalVars struct {
	// progdefs.q1

	_                 [28]Int
	Self              Int
	Other             Int
	World             Int
	Time              Float
	FrameTime         Float
	ForceRetouch      Float
	MapName           Int
	Deathmatch        Float
	Coop              Float
	Teamplay          Float
	ServerFlags       Float
	TotalSecrets      Float
	TotalMonsters     Float
	FoundSecrets      Float
	KilledMonsters    Float
	Parm1             Float
	Parm2             Float
	Parm3             Float
	Parm4             Float
	Parm5             Float
	Parm6             Float
	Parm7             Float
	Parm8             Float
	Parm9             Float
	Parm10            Float
	Parm11            Float
	Parm12            Float
	Parm13            Float
	Parm14            Float
	Parm15            Float
	Parm16            Float
	VForward          Vec3
	VUp               Vec3
	VRight            Vec3
	TraceAllSolid     Float
	TraceStartSolid   Float
	TraceFraction     Float
	TraceEndPos       Vec3
	TracePlaneNormal  Vec3
	TracePlaneDist    Float
	TraceEnt          Int
	TraceInOpen       Float
	TraceInWater      Float
	MsgEntity         Int
	Main              Int
	StartFrame        Int
	PlayerPreThink    Int
	PlayerPostThink   Int
	ClientKill        Int
	ClientConnect     Int
	PutClientInServer Int
	SetNewParms       Int
	SetChangeParms    Int
}

func decodeGlobalVars(r io.Reader) (*GlobalVars, error) {
	var out GlobalVars
	if err := read(r, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type EntVars struct {
	ModelIndex   Float
	AbsMin       Vec3
	AbsMax       Vec3
	LTime        Float
	MoveType     Float
	Solid        Float
	Origin       Vec3
	OldOrigin    Vec3
	Velocity     Vec3
	Angles       Vec3
	AVelocity    Vec3
	PunchAngle   Vec3
	classname    String
	Model        String
	Frame        Float
	Skin         Float
	Effects      Float
	Mins         Vec3
	Maxs         Vec3
	Size         Vec3
	Touch        Func
	Use          Func
	Think        Func
	Blocked      Func
	NextThink    Float
	GroundEntity Int
	Health       Float
	Frags        Float
	Weapon       Float
	WeaponModel  String
	WeaponFrame  Float
	CurrentAmmo  Float
	AmmoShells   Float
	AmmoNails    Float
	AmmoRockets  Float
	AmmoCells    Float
	Items        Float
	TakeDamage   Float
	Chain        Int
	Deadflag     Float
	ViewOfs      Vec3
	Button0      Float
	Button1      Float
	Button2      Float
	Impulse      Float
	FixAngle     Float
	VAngle       Vec3
	IdealPitch   Float
	NetName      String
	Enemy        Int
	Flags        Float
	ColorMap     Float
	Team         Float
	MaxHealth    Float
	TeleportTime Float
	ArmorType    Float
	ArmorValue   Float
	WaterLevel   Float
	WaterType    Float
	IdealYaw     Float
	YawSpeed     Float
	AimEnt       Int
	GoalEntity   Int
	SpawnFlags   Float
	Target       String
	TargetName   String
	DmgTake      Float
	DmgSave      Float
	DmgInflictor Int
	Owner        Int
	MoveDir      Vec3
	Message      String
	Sounds       Float
	Noise        String
	Noise1       String
	Noise2       String
	Noise3       String
}
