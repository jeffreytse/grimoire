/**
 * grimoire plugin for OpenCode
 *
 * Registers all grimoire skills and injects bootstrap context (AGENTS.md).
 */

import path from 'path';
import fs from 'fs';
import os from 'os';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

const normalizePath = (p, homeDir) => {
  if (!p || typeof p !== 'string') return null;
  let normalized = p.trim();
  if (!normalized) return null;
  if (normalized.startsWith('~/')) {
    normalized = path.join(homeDir, normalized.slice(2));
  } else if (normalized === '~') {
    normalized = homeDir;
  }
  return path.resolve(normalized);
};

let _bootstrapCache = undefined;

export const GrimoirePlugin = async ({ client, directory }) => {
  const homeDir = os.homedir();
  const skillsDir = path.resolve(__dirname, '../../skills');
  const envConfigDir = normalizePath(process.env.OPENCODE_CONFIG_DIR, homeDir);
  const configDir = envConfigDir || path.join(homeDir, '.config/opencode');

  const getBootstrapContent = () => {
    if (_bootstrapCache !== undefined) return _bootstrapCache;

    const agentsPath = path.resolve(__dirname, '../../AGENTS.md');
    if (!fs.existsSync(agentsPath)) {
      _bootstrapCache = null;
      return null;
    }

    const content = fs.readFileSync(agentsPath, 'utf8');

    const toolMapping = `**Tool Mapping for OpenCode:**
When skills reference tools you don't have, substitute OpenCode equivalents:
- \`Skill\` tool → OpenCode's native \`skill\` tool
- \`Read\`, \`Write\`, \`Edit\`, \`Bash\` → Your native tools
- \`TodoWrite\` → \`todowrite\`

Use OpenCode's native \`skill\` tool to list and load grimoire skills.`;

    _bootstrapCache = `<EXTREMELY_IMPORTANT>
You have grimoire — quality-gated best-practice skills across every domain.

${content}

${toolMapping}
</EXTREMELY_IMPORTANT>`;

    return _bootstrapCache;
  };

  return {
    config: async (config) => {
      config.skills = config.skills || {};
      config.skills.paths = config.skills.paths || [];
      if (!config.skills.paths.includes(skillsDir)) {
        config.skills.paths.push(skillsDir);
      }
    },

    'experimental.chat.messages.transform': async (_input, output) => {
      const bootstrap = getBootstrapContent();
      if (!bootstrap || !output.messages.length) return;
      const firstUser = output.messages.find(m => m.info.role === 'user');
      if (!firstUser || !firstUser.parts.length) return;

      if (firstUser.parts.some(p => p.type === 'text' && p.text.includes('EXTREMELY_IMPORTANT'))) return;

      const ref = firstUser.parts[0];
      firstUser.parts.unshift({ ...ref, type: 'text', text: bootstrap });
    }
  };
};

export default GrimoirePlugin;
