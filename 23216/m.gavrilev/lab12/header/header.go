package header

//header

const (
	MAXARGS = 256
	MAXCMDS = 50
)

type Command struct {
	Cmdargs [MAXARGS]string
	Cmdflag byte
}

/*  cmdflag's  */
const (
	OUTPIP = 01
	INPIP  = 02
)

var Cmds []Command
var Infile, Outfile, Appfile string
var Bkgrnd byte
