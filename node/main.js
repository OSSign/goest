// import "node:fs";
// import { readFileSync, readFile } from 'node:fs';
// import { fs } from 'node:fs/promises';
// import { global } from 'node:global';

const fs = import('node:fs');

import { readFileSync } from 'node:fs';
import './wasm_exec.js'; // This adds global.Go

var module = null;

const go = new global.Go();

go.exit = process.exit;

function insertText(text) {
   // Get the address of the writable memory.
   let addr = module.exports.MkBuffer()
   let buffer = module.exports.mem.buffer

   let mem = new Int8Array(buffer)
   let view = mem.subarray(addr, addr + text.length)

   for (let i = 0; i < text.length; i++) {
      view[i] = text.charCodeAt(i)
   }

   // Return the address we started at.
   return [addr, text.length];
}

function textArguments(...texts) {
    let args = [];

    for (let text of texts) {
        let [addr, length] = insertText(text);
        args.push(addr, length);
    }

    return args;
}

WebAssembly.instantiate(readFileSync("node/goast.wasm"), go.importObject).then((result) => {
	process.on("exit", (code) => { // Node.js exits if no event handler is pending
		if (code === 0 && !go.exited) {
			// deadlock, make Go print error and stack traces
			go._pendingEvent = { id: 0 };
			go._resume();
		}
	});

    module = result.instance;
	go.run(result.instance);

    result.instance.exports.Configure(...textArguments(
        "",
        "",
        "",
        "",
    ));

    // result.instance.exports.Sign(...textArguments("node/test1.ps1;node/test2.ps1"));
    result.instance.exports.Sign(...textArguments("test1.ps1"));
}).catch((err) => {
	console.error(err);
	process.exit(1);
});

// Runner();
