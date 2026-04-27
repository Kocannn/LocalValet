import * as WailsApp from "../../wailsjs/go/main/App.js";

export async function openContextTerminal(projectDir: string = ""): Promise<void> {
  try {
    await (WailsApp as any).LaunchTerminal(projectDir);
  } catch (error) {
    console.error('Failed to open context terminal:', error);
    throw error;
  }
}
