const packageJson = require('./package.json');
const path = require('path');

// Mapping from Node's `process.arch` to Golang's `$GOARCH`
const ARCH_MAPPING = {
  ia32: '386',
  x64: 'amd64',
  arm: 'arm',
  arm64: 'arm64'
};

// Mapping between Node's `process.platform` to Golang's
const PLATFORM_MAPPING = {
  darwin: 'darwin',
  linux: 'linux',
  win32: 'windows',
  freebsd: 'freebsd'
};

const name = 'safebox';
const version = packageJson.version;
const platform = PLATFORM_MAPPING[process.platform];
const arch = ARCH_MAPPING[process.arch];
const binaryName = platform === 'win32' ? `${name}.ext` : name;
const tarUrl = `https://github.com/adikari/safebox/releases/download/v${version}/safebox_${version}_${platform}_${arch}.tar.gz`;
const bin = path.join(__dirname, "bin");

const constants = {
  name,
  version,
  platform,
  arch,
  binaryName,
  bin,
  tarUrl,
};

module.exports = constants;
