# goPaxos
goPaxos is an implementation of a _single instance_ of the Paxos consensus algorithm in Go. It is based on the outline given in [Paxos Made Simple](https://lamport.azurewebsites.net/pubs/paxos-simple.pdf). The main components of the algorithm - proposers, acceptors and learners - are implemented in their respective packages.

## Usage
1. Run `make` in the `src/paxos` directory.
```bash
> make
go install ./...
```
2. Add `bin/` from the repository to your `PATH`.
3. Create a configuration file with the _host:port_ information for proposers, acceptors and learners. 
```bash
> init-config -np 2 -na 3 -nl 2
{
  "Proposers": [
    "localhost:34379",
    "localhost:34380"
  ],
  "Acceptors": [
    "localhost:34381",
    "localhost:34382",
    "localhost:34383"
  ],
  "Learners": [
    "localhost:34384",
    "localhost:34385"
  ]
}
```
4. Start proposer, acceptors, learners.
```bash
> init-proposer &
> init-acceptor &
> init-learner &
```
## Tests

1. **Simple**: This is the normal case. A proposer proposes a value. All learners should get that value eventually.
```bash
> simple
PASS.
```  
2. **Concur**: Two proposers propose a value concurrently. All learners should get the same value eventually.
```bash
> concur
PASS with hello_world2
```

### TODO
- [ ] **Distinguished Proposer** - Use zookeeper for leader election.
- [ ] **Multi-Paxos** - Support for multiple instances of the algorithm to build a replicated state machine.
- [ ]  **Fault Tolerance** - Test single/multi paxos with random failures.