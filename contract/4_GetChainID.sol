// SPDX-License-Identifier: MIT
pragma solidity 0.8.0;

contract ChainIDTest {

    function getChainID() external view returns (uint256) {
        uint256 id;
        assembly {
            id := chainid()
        }
        return id;
    }
}