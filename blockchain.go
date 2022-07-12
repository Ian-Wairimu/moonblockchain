package main

import (
	"fmt"
	"strings"
	"time"
)

type Block struct {
	nonce       int
	prevHash    string
	time        int64
	transaction []string
}

func NewBlock(nonce int, prevHash string) *Block {
	b := new(Block)
	b.time = time.Now().UnixNano()
	b.nonce = nonce
	b.prevHash = prevHash
	return b
}
func (b *Block) Print() {
	fmt.Printf("previous_hash  %s \n", b.prevHash)
	fmt.Printf("timestamp      %d \n", b.time)
	fmt.Printf("nonce          %d \n", b.nonce)
	fmt.Printf("transaction    %s \n", b.transaction)
}

type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

func NewBlockchain() *Blockchain {
	bc := new(Blockchain)
	bc.CreateBlock(0, "int hash")
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, prevHash string) *Block {
	b := NewBlock(nonce, prevHash)
	bc.chain = append(bc.chain, b)
	return b
}
func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s \n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s \n", strings.Repeat("*", 50))
}

func main() {
	blockchain := NewBlockchain()
	blockchain.Print()
	blockchain.CreateBlock(0, "this is hash1")
	blockchain.Print()
}
