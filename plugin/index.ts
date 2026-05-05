/**
 * Pinchtab OpenClaw Plugin
 *
 * Two tools:
 * - `pinchtab`: Full-featured browser control with all actions
 * - `browser`: OpenClaw-compatible simplified interface
 */

import { definePluginEntry } from "openclaw/plugin-sdk/plugin-entry";
import type { PluginApi, PluginConfig, PluginTool } from "./types.js";
import { pinchtabToolSchema, pinchtabToolDescription, executePinchtabAction } from "./tools/pinchtab.js";
import { browserToolSchema, browserToolDescription, executeBrowserAction } from "./tools/browser.js";

function getConfig(api: PluginApi): PluginConfig {
  return (api.pluginConfig ?? api.config?.plugins?.entries?.pinchtab?.config ?? {}) as PluginConfig;
}

export default definePluginEntry({
  id: "pinchtab",
  name: "Pinchtab",
  description: "Browser control for AI agents via Pinchtab.",
  register(api) {
    const cfg = getConfig(api);

    const pinchtabTool = {
      name: "pinchtab",
      label: "PinchTab",
      description: pinchtabToolDescription,
      parameters: pinchtabToolSchema,
      async execute(_id: string, params: any) {
        return executePinchtabAction(getConfig(api), params);
      },
    } satisfies PluginTool;
    api.registerTool(pinchtabTool);

    if (cfg.registerBrowserTool !== false) {
      const browserTool = {
        name: "browser",
        label: "Browser",
        description: browserToolDescription,
        parameters: browserToolSchema,
        async execute(_id: string, params: any) {
          return executeBrowserAction(getConfig(api), params);
        },
      } satisfies PluginTool;
      api.registerTool(browserTool);
    }
  },
});
