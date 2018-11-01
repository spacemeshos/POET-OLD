package poet

import "flag"

var (
	n        = *flag.Int("DAG_Size", 4, "Set DAG Size") // Low for testing. TODO: Set correct default
	m        = *flag.Int("DAG_Size_Store", 4, "Set DAG Size to Store")
	t        = *flag.Int("Security_Param_t", 150, "Set the Security Parameter t")
	w        = *flag.Int("Security_Param_w", 256, "Set the Security Parameter w")
	filepath = *flag.String("filepath", "./test.dat", "Set the prover file storage path")
)

func init() {
	flag.Parse()
}
