package types

// WeightedOutputs is to make structs which is using weight to be based on
type WeightedOutputs struct {
	ResultEntries []int
	Weight        string
}

// WeightedOutputsList is a struct to keep items which can be generated by weight;
// ItemOutput and CoinOutput is possible in current stage
type WeightedOutputsList []WeightedOutputs
