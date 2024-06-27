package types

func NewGenesisState(params Params, pools []Pool, poolShares []PoolShare) *GenesisState {
	return &GenesisState{
		Params:         params,
		NextPoolNumber: 1,
		Pools:          pools,
		PoolShares:     poolShares,
	}
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:         DefaultParams(),
		NextPoolNumber: 1,
	}
}

func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	return nil
}
