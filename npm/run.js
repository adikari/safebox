const { spawnSync } = require('child_process'),
  constants = require('./constants');

const binaryPath = `${constants.bin}/${constants.binaryName}`;

const [, , ...args] = process.argv;

const options = { cwd: process.cwd(), stdio: "inherit" };

const result = spawnSync(binaryPath, args, options);

if (result.error) {
  console.error(result.error);
  process.exit(1);
}

process.exit(result.status);
