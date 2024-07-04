package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"strings"
)

func (p Pool) GetFee() sdk.Dec {
	return p.Fee
}

func (p *Pool) SetFee(v sdk.Dec) {
	p.Fee = v
}

func (p *Pool) AddToToken1(token1 sdk.Coin) {
	p.Token1 = p.Token1.Add(token1)
}

func (p *Pool) AddToToken2(token1 sdk.Coin) {
	p.Token2 = p.Token2.Add(token1)
}

func (p *Pool) SubtractFromToken1(token1 sdk.Coin) (err error) {
	p.Token1, err = p.Token1.SafeSub(token1)
	if err != nil {
		return err
	}
	return nil
}

func (p *Pool) SubtractFromToken2(token2 sdk.Coin) (err error) {
	p.Token2, err = p.Token2.SafeSub(token2)
	if err != nil {
		return err
	}
	return nil
}

func (p Pool) GetCreatorAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(p.Creator)
}

func (p Pool) GetPoolAddress() sdk.AccAddress {
	key := append([]byte("pool"), sdk.Uint64ToBigEndian(p.Id)...)
	return address.Module(ModuleName, key)
}

func (p Pool) GetPoolShareDenom() string {
	return fmt.Sprintf("pool/%d", p.GetId())
}

func (p Pool) GetKey() []byte {
	return GetPoolKey(p.GetId())
}

func (p Pool) GetToken2by1Ratio() sdk.Dec {
	return p.GetToken2().Amount.ToLegacyDec().Quo(p.GetToken1().Amount.ToLegacyDec())
}

func (p Pool) GetPoolTokensFromName() (string, string) {
	list := strings.Split(p.Name, "/")
	return list[2], list[3]
}

func (p Pool) Validate() error {
	if p.Name == "" {
		return ErrEmptyPoolName
	}
	if err := p.Token1.Validate(); err != nil {
		return err
	}
	if err := p.Token2.Validate(); err != nil {
		return err
	}
	if p.Token1.Denom == p.Token2.Denom {
		return ErrInvalidToken.Wrapf("both tokens cannot be same")
	}
	token1, token2 := p.GetPoolTokensFromName()
	if token1 != p.Token1.Denom || token2 != p.Token2.Denom {
		return ErrInvalidToken.Wrapf("token denom does not match pool name tokens")
	}
	if (p.Token1.IsZero() && p.Token2.Amount.GT(sdk.ZeroInt())) || (p.Token2.IsZero() && p.Token1.Amount.GT(sdk.ZeroInt())) {
		return ErrPoolNotEmpty.Wrapf("one token is zero and other is not")
	}
	if p.Fee.LTE(sdk.ZeroDec()) || p.Fee.GTE(sdk.OneDec()) {
		return ErrInvalidFees.Wrapf("fees is <= 0 or >= 1")
	}
	if _, err := sdk.AccAddressFromBech32(p.Creator); err != nil {
		return err
	}
	return nil
}

func getPoolName(id uint64, token1, token2 sdk.Coin) string {
	return fmt.Sprintf("pool/%d/%s/%s", id, token1.Denom, token2.Denom)
}

func GetInitialPoolShares(token1, token2 sdk.Coin) (sdk.Int, error) {
	initialTotalShareDec, err := token1.Amount.ToLegacyDec().Mul(token2.Amount.ToLegacyDec()).ApproxSqrt()
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return initialTotalShareDec.Mul(sdk.OneDec().Quo(sdk.SmallestDec())).TruncateInt(), nil
}

func NewPool(id uint64, token1, token2 sdk.Coin, fee sdk.Dec, creator sdk.AccAddress) Pool {
	return Pool{
		Id:      id,
		Name:    getPoolName(id, token1, token2),
		Token1:  token1,
		Token2:  token2,
		Fee:     fee,
		Creator: creator.String(),
	}
}

func MustMarshalPool(cdc codec.BinaryCodec, pool *Pool) []byte {
	return cdc.MustMarshal(pool)
}

func MustUnmarshalPool(cdc codec.BinaryCodec, value []byte) (pool Pool) {
	err := cdc.Unmarshal(value, &pool)
	if err != nil {
		panic(err)
	}

	return pool
}
