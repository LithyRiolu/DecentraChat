package peers

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"io/ioutil"
	//"log"

	//"io/ioutil"
	//"log"
	"net/http"
	"sync"
	//"time"
	"../httphelp"
	"../identity"
)

type AliveBeat struct {
	Identity   identity.Identity `json="identity"`
	AlivePeers Peers             `json="peers"`
	Hops       int8              `json="hops"`
	mux        sync.Mutex
}

type Peers struct {
	PeerMap map[string]Peer
	mux     sync.Mutex
}

type Peer struct {
	Id   string
	Addr string
}

func NewAliveBeat(Identity identity.Identity, peers Peers) AliveBeat {
	return AliveBeat{
		Identity:   Identity,
		AlivePeers: peers,
		Hops:       2,
	}
}

func NewPeers() Peers {
	return Peers{
		PeerMap: make(map[string]Peer, 0),
	}
}

func (peers *Peers) Add(peer Peer) {
	peers.mux.Lock()
	defer peers.mux.Unlock()

	peers.PeerMap[peer.Id] = peer
}

func (peers *Peers) Delete(id string) {
	peers.mux.Lock()
	defer peers.mux.Unlock()

	fmt.Println("deleting id - ", id)
	delete(peers.PeerMap, id)

}

func (peers *Peers) CopyPeers() map[string]Peer {
	peers.mux.Lock()
	defer peers.mux.Unlock()

	copyOfPeers := make(map[string]Peer)
	copyOfPeers = peers.PeerMap

	return copyOfPeers
}

func (peers *Peers) ConvertPeersToJSON() []byte {
	copyOfPeers := peers.CopyPeers()
	peersJSON, _ := json.Marshal(copyOfPeers)

	return peersJSON
}

func PrepareAliveBeatJSON(identity identity.Identity, peers Peers) string {
	aliveBeat := NewAliveBeat(identity, peers)
	aliveBeatJson, _ := json.Marshal(aliveBeat)
	return string(aliveBeatJson)
}

func SendAliveBeat(identity identity.Identity, peersDS *Peers, deletedPeersDS *Peers) {

	peerMapThisTime := peersDS.CopyPeers()
	peersThisTime := Peers{
		PeerMap: peerMapThisTime,
	}
	AliveBeatJSON := PrepareAliveBeatJSON(identity, peersThisTime)

	for i, p := range peersThisTime.PeerMap {
		//fmt.Println("Sending Alive Beat to : ", peersThisTime.PeerMap, "\n AliveBeatJSON :", AliveBeatJSON) //todo
		_, err := http.Post("http://"+p.Addr+"/peers", "json", bytes.NewBuffer([]byte(AliveBeatJSON)))
		if err != nil {
			fmt.Println("Error in SendAliveBeat - for : ", i, "  err : ", err)
			peersDS.Delete(i)
			deletedPeersDS.Add(p)
			continue
		}
	}

}

func RecvPeerAlive(w http.ResponseWriter, r *http.Request, idDs identity.Identity, peersDS *Peers, deletedPeersDS *Peers) {

	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	log.Println("Error in RecvPeerAlive : ", err)
	//}
	//defer r.Body.Close()

	body := httphelp.ReadHttpRequestBody(r)

	aliveBeat := JsonToAliveBeat(string(body))
	peersDS.InjectAliveBeatToPeers(aliveBeat, idDs, peersDS, deletedPeersDS)

}

func JsonToAliveBeat(aliveBeatJson string) AliveBeat {

	aliveBeat := AliveBeat{}
	jerr := json.Unmarshal([]byte(aliveBeatJson), &aliveBeat)
	if jerr != nil {
		fmt.Println("Err in unmarshall of AliveBeat : err", jerr)
	}

	return aliveBeat
}

func (peers *Peers) InjectAliveBeatToPeers(aliveBeat AliveBeat, idDS identity.Identity, alivePeers *Peers, deadPeers *Peers) {

	//fmt.Println("Injecting Alive beat")
	existingPeers := peers.CopyPeers()
	if _, ok := existingPeers[aliveBeat.Identity.Id]; !ok { //adding sender ids and sender
		thisPeer := Peer{
			Id:   aliveBeat.Identity.Id,
			Addr: aliveBeat.Identity.Addr,
		}
		peers.Add(thisPeer)

		if _, ok := deadPeers.PeerMap[thisPeer.Id]; ok { //does exist in deadPeers
			delete(deadPeers.PeerMap, aliveBeat.Identity.Id)
		}
	}

	for aliveid, alive := range aliveBeat.AlivePeers.PeerMap {
		//fmt.Println("aliveid : ", aliveid)

		if _, ok := existingPeers[aliveid]; !ok { //does not exist in alivePeers
			if _, ok := deadPeers.PeerMap[aliveid]; !ok { //does not exist in deadPeers
				if aliveid != idDS.Id {
					peers.Add(alive)

				}
			}
		}
	}

}

//
//func ConvertJSONToPeers(peersJSON string) map[string]Peer {
//
//	newPeers := make(map[string]Peer)
//	jerr := json.Unmarshal([]byte(peersJSON), &newPeers)
//	if jerr != nil {
//		log.Println("Err in unmarshalling ConvertPeersJSONToPeers - err :", jerr)
//	}
//
//	return newPeers
//}

func ShowPeerAlive(w http.ResponseWriter, r *http.Request, identitiy identity.Identity, peers Peers) {
	fmt.Println(peers)

	fmt.Fprint(w, "Peers : \n")
	AlivePeers := peers.CopyPeers()

	fmt.Fprint(w, "Peer List of : ", identitiy.Id+"\n")
	for peerId, peer := range AlivePeers {
		fmt.Fprint(w, peerId+" : "+peer.Addr+"\n")
	}

}
