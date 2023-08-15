const { Binary } = require("binary-install");
const os = require("os");
const { join } = require("path");
const cTable = require("console.table");

const error = msg => {
  console.error(msg);
  process.exit(1);
};

const { version, repository } = require("./package.json");

const name = 'safebox';

const supportedPlatforms = [
  {
    TYPE: "Darwin",
    ARCHITECTURE: "x64",
    TARGET: "darwin_amd64",
    NAME: name
  },
  {
    TYPE: "Darwin",
    ARCHITECTURE: "arm64",
    TARGET: "darwin_arm64",
    NAME: name
  },
  {
    TYPE: "Windows_NT",
    ARCHITECTURE: "x64",
    TARGET: "windows_amd64.exe",
    NAME: `${name}.exe`
  },
  {
    TYPE: "Windows_NT",
    ARCHITECTURE: "arm64",
    TARGET: "windows_arm64.exe",
    NAME: `${name}.exe`
  },
  {
    TYPE: "Windows_NT",
    ARCHITECTURE: "ia32",
    TARGET: "windows_386.exe",
    NAME: `${name}.exe`
  },
  {
    TYPE: "Linux",
    ARCHITECTURE: "ia32",
    TARGET: "linux_386",
    NAME: name
  },
  {
    TYPE: "Linux",
    ARCHITECTURE: "x64",
    TARGET: "linux_amd64",
    NAME: name
  },
  {
    TYPE: "Linux",
    ARCHITECTURE: "arm64",
    TARGET: "linux_arm64",
    NAME: name
  },
];

const getPlatform = () => {
  const type = os.type();
  const architecture = os.arch();

  for (let supportedPlatform of supportedPlatforms) {
    if (
      type === supportedPlatform.TYPE &&
      architecture === supportedPlatform.ARCHITECTURE
    ) {
      return supportedPlatform;
    }
  }

  error(
    `Platform with type "${type}" and architecture "${architecture}" is not supported by ${name}.\nYour system must be one of the following:\n\n${cTable.getTable(
      supportedPlatforms
    )}`
  );
};

const getBinary = () => {
  const platform = getPlatform();
  
  const url = `${repository.url}/releases/download/v${version}/${name}_${version}_${platform.TARGET}.tar.gz`;
  
  return new Binary(platform.NAME, url);
};

const run = () => {
  const binary = getBinary();
  binary.run();
};

const install = (supressLogs = false) => {
  const binary = getBinary();
  return binary.install({}, supressLogs);
};

module.exports = {
  install,
  run
};
