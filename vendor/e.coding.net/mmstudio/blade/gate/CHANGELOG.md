# Changelog


### [0.0.2-beta.9](https://e.coding.net/mmstudio/blade/gate/compare/v0.0.2-beta.9%0Dv0.0.2-beta.8) (2021-05-27)

### [0.0.2-beta.8](https://e.coding.net/mmstudio/blade/gate/compare/v0.0.2-beta.8%0Dv0.0.2-beta.7) (2021-05-19)


### Refactor

* **script:** modify release script ([bd4641b](https://e.coding.net/mmstudio/blade/gate/commits/bd4641b4ed8d4c94717fa4b79f713ec4d30abb3c))

### [0.0.2-beta.7](https://e.coding.net/mmstudio/blade/gate/compare/v0.0.2-beta.7%0Dv0.0.2-beta.6) (2021-05-19)


### Refactor

* release script ([44507d5](https://e.coding.net/mmstudio/blade/gate/commits/44507d59e519f34ff32636af9fc89757506c6f88))

### [0.0.2-beta.6](https://e.coding.net/mmstudio/blade/gate/compare/v0.0.2-beta.6%0Dv0.0.2-beta.5) (2021-05-19)


### Refactor

* release script ([1d143e9](https://e.coding.net/mmstudio/blade/gate/commits/1d143e98746f62b66cb898cfcb540b48d6402064))

### [0.0.2-beta.5](https://e.coding.net/mmstudio/blade/gate/compare/v0.0.2-beta.4%0Dv0.0.2-beta.5) (2021-05-13)

### Download URL

 * [Mac](https://zhongtai.s3.amazonaws.com/software/gate/gate_osx_0.0.2-beta.5.zip)  /
[Linux](https://zhongtai.s3.amazonaws.com/software/gate/gate_linux_0.0.2-beta.5.zip)

### Features

* **websocket:** add CheckOrigin ([13008f7](https://e.coding.net/mmstudio/blade/gate/commit/13008f7260d67f02c2c3c604b41544a3f56fd017))


### Refactor

* **golib:** upgrade golib to v0.3.11 ([bf34286](https://e.coding.net/mmstudio/blade/gate/commit/bf3428666a6ae850e6a117afbfa8516e0b1529b1))
* **websocket:** websocket use WsProtocol different from tcp ([c4d955c](https://e.coding.net/mmstudio/blade/gate/commit/c4d955c40b6c320e1e21c97965b397cc7c752286))

### [0.0.2-beta.4](https://e.coding.net/mmstudio/blade/gate/compare/v0.0.2-beta.3%0Dv0.0.2-beta.4) (2021-05-08)

### Download URL

 * [Mac](https://zhongtai.s3.amazonaws.com/software/gate/gate_osx_0.0.2-beta.4.zip)  /
[Linux](https://zhongtai.s3.amazonaws.com/software/gate/gate_linux_0.0.2-beta.4.zip)

### Features

* **secret:** add security check of handshake (sign and time) ([ccdc954](https://e.coding.net/mmstudio/blade/gate/commit/ccdc954fac42585c7d930004d92adc51faec4782))

### [0.0.2-beta.3](https://e.coding.net/mmstudio/blade/gate/compare/v0.0.2-beta.2%0Dv0.0.2-beta.3) (2021-05-06)

### Download URL

 * [Mac](https://zhongtai.s3.amazonaws.com/software/gate/gate_osx_0.0.2-beta.3.zip)  /
[Linux](https://zhongtai.s3.amazonaws.com/software/gate/gate_linux_0.0.2-beta.3.zip)

### Refactor

* **proto:** change meta from map to repeated struct for proto2 ([a80bb80](https://e.coding.net/mmstudio/blade/gate/commit/a80bb80dad0ab746dac450b9b0e63f2744077ddd))

### [0.0.2-beta.2](https://e.coding.net/mmstudio/blade/gate/compare/v0.0.2-beta%0Dv0.0.2-beta.2) (2021-05-06)

### Download URL

 * [Mac](https://zhongtai.s3.amazonaws.com/software/gate/gate_osx_0.0.2-beta.2.zip)  /
[Linux](https://zhongtai.s3.amazonaws.com/software/gate/gate_linux_0.0.2-beta.2.zip)

### Features

* **selector:** add etcd registry support ([ede44e1](https://e.coding.net/mmstudio/blade/gate/commit/ede44e151d649d4db15f0383b5c89d884a3f2c67))

### [0.0.2-beta](https://e.coding.net/mmstudio/blade/gate/compare/v0.0.1...v0.0.2-beta) (2021-04-19)

### Download URL

 * [Mac](https://zhongtai.s3.amazonaws.com/software/gate/gate_osx_0.0.2-beta.zip)  /
[Linux](https://zhongtai.s3.amazonaws.com/software/gate/gate_linux_0.0.2-beta.zip)

### Features

* add common gate server ([d464e5e](https://e.coding.net/mmstudio/blade/gate/commit/d464e5ea9aaddaf78265be4b49722f8a3f89309f))
* add hotfix test ([93c58fd](https://e.coding.net/mmstudio/blade/gate/commit/93c58fd4315a3139ea76a0dd472385c6371ae8b7))
* add plugin manager and example ([b066afb](https://e.coding.net/mmstudio/blade/gate/commit/b066afb4e5405cdb9ce14ba68ef94c0d976cecc9))
* add selector ([582d519](https://e.coding.net/mmstudio/blade/gate/commit/582d519689f9739225e387821cae99ba52a4f2a3))
* add test ([0bcce7e](https://e.coding.net/mmstudio/blade/gate/commit/0bcce7eee9d264420bbe0202cd8d7cbba87e6a30))


### Bug Fixes

* close conn if handshake fail ([58fabd6](https://e.coding.net/mmstudio/blade/gate/commit/58fabd69c539418eac5d78aa1e5a2bb1b19d1d3f))


### Refactor

* add logbus; add mock ([1c92b54](https://e.coding.net/mmstudio/blade/gate/commit/1c92b54d47441896794c117601e2946154509f1d))
* add release script ([289ba34](https://e.coding.net/mmstudio/blade/gate/commit/289ba3450e56f4332dfeffb2f67fad784da3722a))
* add Selector,BackEndHandshake option ([01449c7](https://e.coding.net/mmstudio/blade/gate/commit/01449c75ff2cfca65f9a5caa8efb7a619ed40d74))
* modify release ([c409146](https://e.coding.net/mmstudio/blade/gate/commit/c409146d515a28b79d5857cb08567ee53db2200e))
* use link.NewWebSocketConn ([f485174](https://e.coding.net/mmstudio/blade/gate/commit/f48517407b768d41a5957fbf3a76e1cc46808866))
