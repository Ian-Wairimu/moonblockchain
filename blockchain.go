package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const MINING_DIFFICULTY = 3

type Block struct {
	nonce       int
	prevHash    [32]byte
	time        int64
	transaction []*Transaction
}

func NewBlock(nonce int, prevHash [32]byte, transaction []*Transaction) *Block {
	b := new(Block)
	b.time = time.Now().UnixNano()
	b.nonce = nonce
	b.prevHash = prevHash
	b.transaction = transaction
	return b
}
func (b *Block) Print() {
	fmt.Printf("previous_hash  %x \n", b.prevHash)
	fmt.Printf("timestamp      %d \n", b.time)
	fmt.Printf("nonce          %d \n", b.nonce)
	for _, t := range b.transaction {
		t.PrintTransaction()
	}
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Transaction  []*Transaction `json:"transactions"`
	}{
		b.time,
		b.nonce,
		b.prevHash,
		b.transaction,
	})
}

type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
}

func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, prevHash [32]byte) *Block {
	b := NewBlock(nonce, prevHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}
func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}
func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s \n", strings.Repeat("=", 10), i, strings.Repeat("=", 20))
		block.Print()
	}
	fmt.Printf("%s \n", strings.Repeat("*", 40))
}
func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
}
func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, ts := range bc.transactionPool {
		transactions = append(transactions, NewTransaction(
			ts.senderBlockchainAddress,
			ts.recipientBlockchainAddress,
			ts.value))
	}
	return transactions
}
func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transaction []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{nonce, previousHash, 0, transaction}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}
func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{
		sender, recipient, value,
	}
}
func (t *Transaction) PrintTransaction() {
	fmt.Printf("%s\n", strings.Repeat("_", 40))
	fmt.Printf("sender_blockchain_address   %s\n", t.senderBlockchainAddress)
	fmt.Printf("recipeient_blockchain_address   %s\n", t.recipientBlockchainAddress)
	fmt.Printf("value   %.1f\n", t.value)
	fmt.Printf("%s\n", strings.Repeat("_", 40))
}
func (t *Transaction) MarshalTransactionJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		t.senderBlockchainAddress,
		t.recipientBlockchainAddress,
		t.value,
	})
}
func main() {
	blockchain := NewBlockchain()

	blockchain.AddTransaction("a", "b", 20000)
	previousHash := blockchain.LastBlock().Hash()
	nonce := blockchain.ProofOfWork()
	blockchain.CreateBlock(nonce, previousHash)
	blockchain.Print()

	blockchain.AddTransaction("x", "y", 200)
	blockchain.AddTransaction("c", "k", 3070)
	previousHash = blockchain.LastBlock().Hash()
	nonce = blockchain.ProofOfWork()
	blockchain.CreateBlock(nonce, previousHash)
	blockchain.Print()
}
