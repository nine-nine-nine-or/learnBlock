// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

contract RCCStake is Initializable,
    UUPSUpgradeable,
    PausableUpgradeable,
    AccessControlUpgradeable{

    using SafeERC20 for ERC20;
    using Math for uint256;
    using Address for address;

    bytes32 public constant STAKER_ROLE = keccak256("STAKER_ROLE");
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");

    uint256 public constant nativeCurrency_PID = 0;

    struct tool {
        // 质押代币的地址
        address stTokenAddress;
        // 质押池的权重，影响奖励分配
        uint256 poolWeight;
        // 最后一次计算奖励的区块号
        uint256 lastRewardBlock;
        // 每个质押代币累积的 RCC 数量
        uint256 accRCCPerST;
        //池中的总质押代币量
        uint256 stTokenAmount;
        // 最小质押金额
        uint256 minDepositAmount;
        // 解除质押的锁定区块数
        uint256 unstakeLockedBlocks;
    }

    
    struct UnstakeRequest {
        // 请求提现金额
        uint256 amount;
        // 请求提现金额的可释放区块号
        uint256 unlockBlocks;
    }
    struct User {
        // 用户质押的代币数量
        uint256 stAmount;
        // 已分配的 RCC 数量
        uint256 finishedRCC;
        // 待领取的 RCC 数量
        uint256 pendingRCC;
        // 解质押请求列表，每个请求包含解质押数量和解锁区块
        UnstakeRequest[] requests;
    }

    //设置rcctoken的地址，以及设置项目部署时候的基本信息
    function initialize()  returns () {
        
    }
}