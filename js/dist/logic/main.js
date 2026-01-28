import { Interface } from './interface.js';
import { Match } from './match.js';
const ui = new Interface();
const playerCount = await ui.askPlayerCount();
const match = new Match(playerCount, ui);
await match.startGame();