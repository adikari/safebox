const fs = require("fs");
const { join } = require("path");
const packageJson = require("./package.json");

// Mapping from Node's `process.arch` to Golang's `$GOARCH`
var ARCH_MAPPING = {
    "ia32": "386",
    "x64": "amd64",
    "arm": "arm",
    "arm64": "arm64"
};

// Mapping between Node's `process.platform` to Golang's
var PLATFORM_MAPPING = {
    "darwin": "darwin",
    "linux": "linux",
    "win32": "windows",
    "freebsd": "freebsd"
};

const error = (msg) => {
  console.error(msg);
  process.exit(1);
};

const name = packageJson.goBinary.name;
const binName = packageJson.goBinary.name;

if (!(process.arch in ARCH_MAPPING) || !(process.platform in PLATFORM_MAPPING)) {
    error(`Platform with type "${process.platform}" and architecture "${process.arch}" is not supported by ${name}.`);
}

if (process.platform === "win32") {
    binName += ".exe";
}

const binaryPath = `./dist/${name}_${process.platform}_${ARCH_MAPPING[process.arch]}/${binName}`;

const installDirectory = join(__dirname, "node_modules", ".bin");

if (!fs.existsSync(installDirectory)) {
  fs.mkdirSync(installDirectory, { recursive: true });
}

const installPath = join(installDirectory, name);

fs.copyFile(binaryPath, installPath, (err) => {
    if (err) error(err.message);
});
