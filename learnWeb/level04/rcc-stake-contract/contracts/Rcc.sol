// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract RccToken is ERC20{

  constructor() ERC20("RccToken", "RCC") {
    _mint(msg.sender, 1000000 * 10 ** decimals());
  }
}