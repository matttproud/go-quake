package prog

import "io"

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
