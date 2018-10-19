package poet

import "flag"

var (
	n = *flag.Int("DAG_Size", 4, "Set DAG Size") // Low for testing. TODO: Set correct default
	m = *flag.Int("DAF_Size_Store", 2, "Set DAG Size to Store")
	t = *flag.Int("Security_Param_t", 150, "Set the Security Parameter t")
	w = *flag.Int("Security_Param_w", 256, "Set the Security Parameter w")
)

func init() {
	flag.Parse()
}
