/**
 * WASI Engine Implementation
 *
 * This implementation uses WASI to run the language server directly in VS Code.
 * The WASI binary provides a standard system interface that we can use.
 *
 * Architecture:
 * ```
 *  +----------------+         +-----------------+
 *  |   VS Code      |         |    gotmpls     |
 *  |  Extension     |<------->|     WASI       |
 *  +----------------+         +-----------------+
 *                   (stdio transport)
 * ```
 */

import * as vscode from "vscode";

import { ProcessOptions, Wasm } from "@vscode/wasm-wasi/v1";
import { startServer, createStdioOptions, createUriConverters  } from "@vscode/wasm-wasi-lsp";
import { LanguageClient, LanguageClientOptions, ServerOptions } from "vscode-languageclient/node";


// WASI module interface
declare global {
	interface WasiEnv {
		args: string[];
		env: { [key: string]: string };
		preopens: { [key: string]: string };
	}
}

export async function wasi_activate(
	context: vscode.ExtensionContext,
	outputChannel: vscode.OutputChannel,
	baseClientOptions: LanguageClientOptions,
	wasmPath: vscode.Uri,
	processName: string,
	env: { [key: string]: string },
) {
	const wasm: Wasm = await Wasm.load();

	const channel = outputChannel;
	// The server options to run the WebAssembly language server.
	const serverOptions: ServerOptions = async () => {
		const options: ProcessOptions = {
			trace: true,

			env: env,
			stdio: createStdioOptions(),
			mountPoints: [{ kind: "workspaceFolder" }],
		};

		// Load the WebAssembly code
		const filename = wasmPath;
		const bits = await vscode.workspace.fs.readFile(filename);
		const module = await WebAssembly.compile(bits);

		const memory: WebAssembly.MemoryDescriptor | WebAssembly.Memory = {
			initial: 10000,
			maximum: 10000,
			shared: true,
			buffer: new ArrayBuffer(10000),
		};

		// Create the wasm worker that runs the LSP server
		const process = await wasm.createProcess(processName, module, memory, options);

		// Hook stderr to the output channel
		const decoder = new TextDecoder("utf-8");
		process.stderr!.onData((data) => {
			channel.appendLine("[wasi-stderr] " + decoder.decode(data));
		});

		// process.stdout!.onData((data) => {
		// 	channel.appendLine("[wasi-stdout] " + decoder.decode(data));
		// });

		return startServer(process);
	};

	baseClientOptions.uriConverters = createUriConverters();

	return [serverOptions, baseClientOptions] as const;
}
