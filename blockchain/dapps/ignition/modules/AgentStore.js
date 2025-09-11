const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

module.exports = buildModule("AgentStoreModule", (m) => {
  const store = m.contract("AgentStore");

  return { store };
});
