package poet

import "errors"

// This type will provide the inteface to the Prover. It implements the
// io.ReadWriter interface, which will allow it to sit behind an RPC Server
// or be linked directly to a verifier.
// CurrentState is used to
type Prover struct {
	CreateChallenge bool
	CurrentState    State
	// other types based on implementation. Eg leveldb client & DAG
}


func NewProver(CreateChallenge bool) *Prover {
	p := new(Prover)
	p.CreateChallenge = CreateChallenge
	return p
}

// Satifying io.ReadWriter interface. In Start State it returns Proof from
// Commitment. In WaitingChalleng State, it returns Challenge Proof. Both
// commitment and challenge are encoded as a byte slice (b). To retrieve
// the proof, the verifier calls Read.
func (p *Prover) Write(b []byte) (n int, err error) {
	if p.CurrentState == Start {
		err = p.CalcCommitProof(b)
		p.CurrentState = Commited
	} else if p.CurrentState == WaitingChallenge {
		err = p.CalcChallengeProof(b)
		p.CurrentState = ProofDone
	} else {
		return 0, errors.New("Prover in Wrong State for Write")
	}
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (p *Prover) Read(b []byte) (n int, err error) {
	// TODO: Check size of b. Read only supposed to send len(b) bytes. If
	// not big enough, need to return error
	if p.CurrentState == Commited {
		b, err = p.SendCommitProof()
		p.CurrentState = WaitingChallenge
	} else if p.CurrentState == ProofDone {
		b, err = p.SendChallengeProof()
		p.CurrentState = Start // For now, this just loops back.
		// TODO: What is the logic to change state? May have other functions to
		// reset prover to Start state (clear DAG and ready for new statement)
	} else {
		return 0, errors.New("Prover in Wrong State for Read")
	}
	return 0, nil
}

func (p *Prover) ConstructDag(m int) *Graph {
	graph := NewGraph()
	
	binstrs := make([]BinaryID)

	for level := 0; level < (m+1); level++ {

		// add graph nodes
		for i := 0; i < (2^level); i++ {
			binstrs[i] = BinaryID(level, i)
			graph.AddNode(BinaryID(level, i))
		}

		// add graph edges
		if level > 0 {
			for node := 0; node < len(binstrs); node++ {
				bit_list := binstrs[node].GetBitList()
				graph.AddEdge(binstrs[node], BinaryID(level - 1, BitsToInt(bit_list[:level - 1])))
			}
		}
	}

	// add edge to the leaf
	for leaf := 0; leaf < len(binstrs); leaf++ {
		bit_list := binstrs[j].GetBitList()
		for i := 0; i < len(bit_list); i++ {
			if bit_list[i - 1] == 1 {
                G.add_edge(BinaryID(i, BitsToInt(bit_list[:i - 1] + [0])), binstrs[leaf])
			}
		}
	}

	return graph
}

func (p *Prover) GetParents(){
    parents = []
   
	bit_list = binary_str.get_bit_list()
    length = binary_str.length
    if length == n:
        for i in range(1, length + 1):
            if bit_list[i - 1] == 1:
                parents.append(BinaryString(i, bits_to_int(bit_list[:i - 1] + [0])))
    else:
        parents.append(BinaryString(length + 1, bits_to_int(bit_list + [0])))
        parents.append(BinaryString(length + 1, bits_to_int(bit_list + [1])))
    return sorted(parents)
}

func (p *Prover) CalcProof(n int, hash, ) error {
	G := p.ConstructDag(n)

	sorted_dag, err := G.TopologicalSort("")

	if err != nil {
		panic("an error occurred")
	}

	for i := 0; i < len(sorted_dag); i++ {
		
	}

    for elem in nx.topological_sort(G):
        hash_str = str(elem)
        for parent in get_parents(elem, n):
            hash_str += str(G.node[parent]['label'])
        G.node[elem]['label'] = H(chi, hash_str)
	return G
	
	return nil
}


func (p *Prover) CalcCommitProof([]byte) error {
	
	return nil
}

func (p *Prover) CalcChallengeProof([]byte) error {
	return nil
}

func (p *Prover) SendCommitProof() (b []byte, err error) {
	return b, nil
}

func (p *Prover) SendChallengeProof() (b []byte, err error) {
	return b, nil
}

type State int

const (
	Start State = iota
	Commited
	WaitingChallenge
	ProofDone
)
