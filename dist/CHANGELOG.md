## Changelog

73f4d2f add .goreleaser.yml
13fcd47 [make] release
ed0e065 [rpc] invoke method
cc20160 [cmd] sign tx util
89805e0 update dockerfile
c020049 go mod vendor
d87ff73 replace dep with go mod
cc9ec85 remove cached artifacts
a2314ec [fileutil] check file exists before removing
183206d removed bug that was being caused by directories being removed before files
719f6b4 [docker] log container output
f1bdb36 call ipfs api refs
403762c bugfix
2f30ff3 now removing temp dirs as well
b66385d [ethereumclient] init
7d704ba [cmd] error color
ac992b6 [registry] ipfs host option
e4c40ea bugfixes
573d00e changed registry.ipfsClient.Refs to recursive
e8a8a87 remade deps
14ceb05 remade deps
05fb6d4 fixed deps
8daa0a0 [node] nil check
46f6854 [node] eos client config
9d98ab9 [deps] update proto, libp2p crypto/netutil revisions
9f40a25 [eosclient] push action method
b19db92 Merge branch 'tmp' into develop
5c189a6 [eosclient] init
5e1c57a minor changes to readme
1ac76ab Merge branch 'tmp' of github.com:c3systems/c3-go into tmp
336d0c4 install jq supports more package managers
40ef784 install jq supports more package managers
7f3b13b use libp2pcrypto pr branch
d3aa75e [trie] fix size
a728eac pegged libp2p/go-libp2p-crypto to a specific rev
6786c22 updated go-libp2p-crypto
111ade3 [config] temp block difficulty
4bf5f81 [trie] init
ecb2e30 [miner] override latest state file
0c38ae9 [miner] write latest state temp file
7a9c88e [cli] image snapshot from state block
30bc978 [fileutil] init
376de39 [wip] image snapshot
825bcce [docker] commit container method
8f37920 [p2p] restore from latest block
3d72650 [rpc] get state chain block
e24d361 [rpc] cli flag
49e00a4 updated diffing lib to handle tar files
535d0c0 [rpc] getBlock error response
7c7b4ca [rpc] get block by number
1fdd92c [rpc] latest block number
e7527a3 [redisstore] return if empty list
b1025f7 [cli] mempool type flag
206a139 redisstore add pending block methods
94207ac rpc init
0f4509a coverage make command
d990e0f use libp2p-crypto pr branch
a92f78c cli display peer id + use patched libp2p-crypto
9a891e2 remove cached
7487978 cli node config filepath option
d66ee87 circleci update coverage file path
9147cf2 merge travis.yml with circleci.yml; closes #66
e1b4a5c rename to hashutil
7676345 c3 config as toml
7ed5ad7 fix circleci go pkg path
1170d0a dep ensure
1676241 lowercase hexutil output
9107d58 dep ensure
1ada3b6 update readme
2b9430f update readme
f5cfa7a add hello world tutorial
b896b68 added clib for config
7898fe0 Merge branch 'feature/clib' into develop
b560c6c added some clibs; closes #77
054caa3 got base58 tests to pass
331340c read ipfs gateway url; closes #71
36ddb33 WIP
b298b90 update circle-ci badge
b7d3ef0 updated readme
5241bae oops, pulled the wrong GNU AGPL v3
27d42d3 Merge branch 'feature/license' into develop
77f87a9 replaced Apache 2.0 with GNU AGPL v3; closes #76
7c330b9 c3 pull docker host flag
cbe7427 remove artifact
cc67550 add generate key command
219861a Merge tag 'v0.0.2' into develop
8a2929e Merge branch 'release/0.0.2'
d67f946 tmp
9349cd3 fix import
af7c8f1 process tx now adds to the pending pool
695acfe update makefile
2335adf add missing type to interface
3fca500 updated process tx
6c321cd sandbox use docker copy instead of mounting
4804a49 remove unused method in docker client test
03d9a64 Merge branch 'develop' of github.com:c3systems/c3 into develop
6c635ab docker copy to/from container methods
b81f41d made deps
021c941 renamed package from c3 to c3-go
029d08b tmp
82f1278 sandbox wait
3f05fe9 sandbox use cached ipfs image
f0d9162 clean up image pull tag name
6066403 docker image tag helper
e9db917 Merge branch 'develop' of github.com:c3systems/c3 into develop
ea02130 moved sdk to its own lib, updated go example
e948cf7 comment
5642b3b Merge branch 'develop' of github.com:c3systems/c3 into develop
3cfcff8 debug
a35b850 chmod and made deps
8131cdd chmod
868bc84 invalid transactions are now removed from the pending transactions pool; closes #61
997dc2f invalid transactions are now removed from the pending transactions pool; closes #61
bceb0d8 invalid transactions are now removed from the pending transactions pool; closes #61
dd17555 changed log.Print -> log.Error where approps
4bc9b22 color package is now returning the correct source
f108ef0 moved miner to core
f554e18 moved miner to core
b83b54f updated travis-ci coveralls token
53901f9 updated travis ci and added coveralls
dec4893 Merge branch 'feature/tests' into develop
c398335 most tests for utils are done
b40a63c added some unit tests for the miner
336defe Merge branch 'develop' of github.com:c3systems/c3 into develop
6dc4ab9 debug logs
4632869 another fix to patch
799cf29 debug logs
4d1095a fixed merge
f403a3f merged develop
b906c4d fixed bugs in diffs and state generation
5144b04 block difficulty cli flag
c25498c hashing pkg tests
5564437 resolve conflicts
0ebc80c stringutil json fix
98aed80 miner now requires sandbox interface
634a01b ensuring the services implement the interfaces
53953d8 debug logs
c9da003 debug logs
8fea60c fix hex encoding
a788ed6 fixed some linting issues
48f6a66 added some readmes
faedcac added some tests; added interfaces and mocks to p2p, docker and registry packages.
c786506 debug logs
a979b88 fixed bug that was causing the wrong file to close
9213133 updated .gitignore
a4cbcc8 Merge branch 'feature/refactor' into develop
2fc3975 Merge branch 'develop' into feature/refactor
28ecf03 mostly finished miner refactor, could be improved; closes #56
71ba774 ipfs cid use standard sha256
22ac5d1 Merge branch 'develop' of github.com:c3systems/c3 into develop
5d908c9 add to mainchain genesis block to ipfs
5852844 updated coder readme
076a3b7 refactored the miner; still need to refactor the verifier
1d58780 fixed logrus source; closes #59
410af0a Merge branch 'feature/routing' into develop
b351b6b added peer discovery; added dht routing
66e9319 added peer discovery; added dht routing
ab56d32 wip
e18464a Merge branch 'develop' of github.com:c3systems/c3 into develop
9622694 logs
041e8d0 on connection, peers are being added to the peer store
a2f9547 sending echo response on startup; on echo receipt, adding remote peer to peerstore
93df98e added IPFSTimeout to configs
1822603 coder now appends a type to bytes
12eba41 resolve merge
94a5681 comment out check for nil state blocks map
f1887b1 added gren to readme
e3168b5 added gren to readme
c3f89c0 made deps
0ec1908 Merge branch 'feature/refactorMiner' into develop
16dc615 fixed merge conflicts
f4578f2 now using context in miner; closes #57
e75aac3 method types pkg
35173de update makefile test
c2c9e65 Merge branch 'develop' of github.com:c3systems/c3 into develop
f12d8ce tmp work around for genesis state block
1b48ccf make temp file helper function
90ff826 fixed make deps
796fecd Merge branch 'develop' of github.com:c3systems/c3 into develop
2af30e7 lint gxundo
be48224 made deps
9ca58ee add debug logs
3b86c1b logrus line number
e650d81 use logrus
615b298 Merge branch 'feature/test' into develop
4707a12 Merge branch 'develop' of github.com:c3systems/c3 into feature/test
d261786 clean cache
6f74a63 fix c3crypto makefile test
fd2574d added mined block test for serialze/deserialize string
1b27de2 closes #55
67858e8 Merge branch 'develop' of github.com:c3systems/c3 into feature/test
1d427ec update gxundo
14d2a0a moved c3crypto to common; closes #51
ad6ac3a docker test use t.Error
853840a Merge branch 'feature/test' into develop
c748f0a make deps
62eeade hexutil test
ecf53b4 resolve merge + fix EncodeBigInt
d767712 update make test
255c5c2 Merge branch 'feature/fixes' into develop
faecdca fix hexutil DecodeString + DecodeBigInt
3a268a7 Merge branch 'feature/VSB' into develop
ec71fd0 can mine, broadcast and validate received blocks
eca329a wip
c74f34d rsa test key
fb6f6f7 wip
099ad8b fixed peer connection issue
f05961d wip
6c9655a Merge branch 'develop' of github.com:c3systems/c3 into feature/constructor
020ad1f wip
46c141a swapped flat disk store for leveldb
4e363d8 ipfs core package
8ce8244 perm fix for libp2p pubsub eff up
e34dc58 fix proto global namespace bug
9f3169b using rlp for serialize/deserialize; closes #50
0480a2a Merge branch 'develop' of github.com:c3systems/c3 into develop
b5928a8 launch ipfs if not running on registry calls
22f2108 Merge branch 'feature/rpc' into develop
ab423d2 Merge branch 'develop' into feature/rpc
3a15f0a added protocol buffers; now fetching headblock on start; closes #12
2a96103 use ipfs api lib
a538499 circleci cleanup
111e3e2 Merge branch 'develop' of github.com:c3systems/c3 into develop
85eb044 circleci move ipfs bin to home bin
f1b8762 Merge branch 'feature/custGOB' into develop
005bd0a gob -> JSON
41e676a circleci ipfs non root
ce0b9f1 now signing block before broadcasting
e0dcfa0 tests set parallel=1 bc of ipfs lock
4c25ecb circleci update working directory
0d6b330 circleci echo tee sudo
f2dd892 circleci fix path
756b4c1 circleci cd working directory
b339b16 circleci update gopath
9c0d86b circleci ls path
871b6f8 circleci change checkout path
9107e78 circleci disable restart
e5dcad8 circleci edit docker daemon.json
8f407d6 circleci checkout repo
1db21c7 circle ipfs sudo
3f9b291 circleci golang ppa
588f271 circleci install golang
61a9a9b circleci change workdir
8809f31 circleci classic image
d6690f9 circleci reference image name update
a8900f1 circleci reference image
4750106 circleci remove setup_remote_docker
2def747 circleci machine map
b1c7ab7 circleci fix workflow
f37004a circleci machine:true
9366217 circleci yml fix
0b3792c circlecli test workflow
abf3d20 circleci machine:false
c913988 Merge branch 'develop' of github.com:c3systems/c3 into develop
715cb06 circleci machine:true
b846cb6 Merge branch 'feature/verifyMinedBlock' into develop
d5b8a41 Merge branch 'develop' into feature/verifyMinedBlock
8a5cc85 added functions to verify a mined block struct
dfd3bdc circleci docker local registry host
565a5b7 circleci loopback alias
53b9a6a circlecli docker local registry host
73e0b2f circleci touch daemon.json
837dac5 circleci set insecure registry
b6be425 circleci test registry
85d75c2 circleci typo
7275ca1 circlecli copy to bin
8dea65e circleci typo fix
469b423 circleci wget rename output
1b48aec circleci wget
7d02e34 circleci fix path
33073ca circleci install IPFS
4ecce06 circleci docker config
6b64126 circleci install docker
45b8636 update makefile
b57e081 add circleci badge
6af3bb8 make deps
f42d7b2 add .circleci
66abc77 miner run container with state
1125e5e Merge branch 'feature/testCID' into develop
54749da fixed merge conflicts
102f714 finished block chain
2b3cbe6 docker sandbox container accept array input
8010b73 sandbox initial state prop
bf954a4 wip
5d84118 add readme file to packages
0ebb385 rename ditto -> registry
ad04bb4 rename ditto -> registry
a86df31 fix container read state
461bd03 fix makefile test
78dbd5b container exec cmd arg
43cfde2 resolve merge
53e8e4b return container tx result
14dca19 Merge branch 'feature/blockchain' into develop
e10ae62 merged with develop; closes #25; closes #26; closes #27; closes #28; closes #36; closes #37; closes #30
16cbe06 first blockchain build complete
6145e28 add log line
6d769c5 update makefile for example
be0fdbe forward payload to c3 registered methods
1377478 wip
bfe910b wip
8205bf3 wip
6df351a resize
11df982 add logo + badges
bcb75f2 sandbox readme
3d93f07 Merge branch 'feature/gob' into develop
8bca097 use gob instead of json marshaller to preserve interfaces
9233f1c resolve merge
7ed261a no longer sending tx on startup
c2a8a71 removed errant file
f5958d8 finished verify tx; closes #34
08df22d finished tx schema; closes #33
25d22fc fixed bug in statechain
cf687d6 added hashing library, built out hexutil lib
6345772 added isDir to diffing.Diff; closes #39
bcf3af6 Merge branch 'feature/miguel' into develop
cbe2a12 kill container after sending message
6fc2dd7 wip
4ce3440 remove underline
8ddc72c align center image
33e82dd Merge branch 'feature/ipfsregistry' into feature/miguel
77403b2 wip
7a36e1a makefile disable echo
ab5e42d use cobra RunE
b6052d6 Merge branch 'develop' into feature/miguel
2e9616e moved hexutil to common
40f831e updated readme
1d5f0d6 Merge branch 'feature/license' into develop
8010025 added apache, removed gnu; closes #14
c15fd53 added diffing
29ddd50 added apache 2.0
5f262cc added ecdsa and ecies crypto lib; closes #15
98a05b6 fixed bug in flat-fs; no longer sending messages on channel when received from self
3de1dd9 wip
c8d8b8d finished ipfs
d36a15a removed erroneous dir
99680f5 merged origin
9002944 finished ipfs
961d3f6 p2p service test init
799bf11 ipfs server up and running
2821dcc it compiles!
011edc3 Merge branch 'feature/core' into develop
5ad47b1 finally fixed deps
f6617d3 dep ensure
400af60 making progress on gx hell
db6180f gx hell
288086b make test cleanup
6f992f1 text
484f453 add GNU LGPLv3 license for now
36ccd9d make cmd test
ec3f517 fix import
6269e35 stop container method
c905c52 run container
8eabbde wip
c185251 Merge branch 'develop' into feature/adam
c8d1f52 wip
3c55e0e wip
7533e0e typo
a9fa458 c3 pull image from ipfs
0cb2d94 docker registry server ipfs proxy
9f67a34 wip
8e502f0 ipfs download image
a3a5f4c close ##11
1f0ab0b c3 deploy image cli
8844fc8 update readme
c8a3861 update readme
2b406df build ldflags
15f7daf makefile build
b96b125 init cobra cli
701f2a4 dep init
b6ef1df spacing
f6c4204 ditto image
5beb3e6 manifest corrections
c99654e base32 hash for docker
eccb4c9 log ipfs upload
a2435ba ipfs recursive upload
5b4af17 gzip image for ipfs
95b6270 image to ipfs wip
1eda6ed untar function
24bc542 update readme
dec64c0 fix import name
43d1563 closes #3
fb2b76a closes #4; closes #5
9fb1eb5 Merge branch 'develop' into feature/adam
da39c4e add method to push images - closes #7
b497b9c update desc
ea3ceb5 add method to pull images
bbc3d10 wip
