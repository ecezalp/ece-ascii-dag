package main

import (
	. "ece-ascii-dag/dag"
)

func main() {
	DAGtoText("random -> pool_urbg\nrandom -> nonsecure_base\nrandom -> seed_sequence\nrandom -> distribution\n\nnonsecure_base -> pool_urbg\nnonsecure_base -> salted_seed_seq\n\nseed_sequence -> pool_urbg\nseed_sequence -> salted_seed_seq\nseed_sequence -> seed_material\n\ndistribution -> strings\n\npool_urbg -> seed_material\n\nsalted_seed_seq -> seed_material\n\nseed_material -> strings")
}
