// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import { ExtensionContext } from 'vscode';
import {
    LanguageClient,
    LanguageClientOptions,
    ServerOptions,
    TransportKind
} from 'vscode-languageclient/node';

import * as vscode from "vscode";
import { wasi_activate } from "./wasi";

// this method is called when your extension is activated
// your extension is activated the very first time the command is executed
export async function activate(context: ExtensionContext) {
    // The server is implemented in node
    // const serverModule = context.asAbsolutePath(
    //     path.join('bin', 'protobuf-language-server')
    // );

    // const serverModule = 'protobuf-language-server';
    // let debugOptions = { execArgv: ["--nolazy", "--debug=6004"] };
    // If the extension is launched in debug mode then the debug server options are used
    // Otherwise the run options are used
    // const serverOptions: ServerOptions = {
    //     module: ,
    //     debug: { command: serverModule, transport: TransportKind.stdio },
    // };

    // // Options to control the language client
    const baseClientOptions: LanguageClientOptions = {
        // Register the server for plain text documents
        documentSelector: [{ scheme: 'file', language: 'proto' }],
    };

    // // Create the language client and start the client.
    // let client = new LanguageClient(
    //     'proto',
    //     'protobuf-language-server',
    //     serverOptions,
    //     clientOptions
    // );

    const outputChannel = vscode.window.createOutputChannel("protolsp");


    const wasmPath = vscode.Uri.joinPath(context.extensionUri, "out", "protolsp.wasi.wasm");


    const [serverOptions, clientOptions] = await wasi_activate(context, outputChannel, baseClientOptions, wasmPath, "protolsp", {
        PROTOLSP_DEBUG: "1",
        RUST_BACKTRACE: "1",
        RUST_LOG: "debug",
    });

    let client = new LanguageClient(
        'proto',
        'protolsp',
        serverOptions,
        clientOptions
    );
    

    // Start the client. This will also launch the server
    client.start();


}

// this method is called when your extension is deactivated
export function deactivate() { }