// SPDX-License-Identifier: MIT

pragma solidity 0.8.7;

/**
 * THIS IS AN EXAMPLE CONTRACT THAT USES UN-AUDITED CODE.
 * DO NOT USE THIS CODE IN PRODUCTION.
 */

// 接口定义函数而没有实现
interface numberComparison {
    function isSameNum(uint a, uint b) external view returns (bool);
}

// 继承接口
contract Test is numberComparison {
    constructor() {}

    // override is necessary, because it overrides the base function contained
    // pure instead of view because the isSameNum function in the Test contract does not return a storage variable.
    function isSameNum(uint a, uint b) external pure override returns (bool) {
        if (a == b) {
            return true;
        } else {
            return false;
        }
    }
}
