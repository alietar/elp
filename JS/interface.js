import inquirer from 'inquirer';
import sleep from 'sleep';

// Couleurs simples pour l'affichage console.
const BG_GREY = "\x1b[47m";
const FG_BLACK = "\x1b[30m";
const RESET = "\x1b[0m";  

class Interface {
  async askPlayerCount() {
    // Demande le nombre de joueurs (min 2).
    const config = await inquirer.prompt([
      {
        type: 'number',
        name: 'playerCount',
        message: 'How many players?',
        default: 2,
        validate: (value) => Number.isInteger(value) && value >= 2 ? true : 'Min 2 players',
      },
    ]);
    return config.playerCount;
  }

  async askPlayerName(player) {
    // Saisie du pseudo.
    const pseudo = await inquirer.prompt([
      {
        type: 'input',
        name: 'playerName',
        message: `Player ${player.player_nb}, what's your name?`,
      },
    ]);
    console.log('Hi', pseudo.playerName, ', Welcome !');
    return pseudo.playerName;
  }

  async chooseTarget(currentPlayer, players) {
    // IMPORTANT : la cible peut être n'importe quel joueur actif, y compris soi-même.
    const targets = players
      .filter(p => p.state)
      .map(p => ({ name: p.name || `Player ${p.player_nb}`, value: p }));

    if (targets.length === 0) return currentPlayer;

    const action = await inquirer.prompt([
      {
        type: 'rawlist',
        name: 'choice',
        message: 'You drew an action card, who do you want to target?',
        choices: targets,
      },
    ]);

    return action.choice;
  }

  async chooseSecondChanceTarget(currentPlayer, players) {
    // On ne peut donner Second Chance qu'à un autre joueur actif qui ne l'a pas déjà.
    const targets = players
      .filter(p => p.state && p !== currentPlayer && !p.hasSecondChance())
      .map(p => ({ name: p.name || `Player ${p.player_nb}`, value: p }));

    if (targets.length === 0) return null;

    const action = await inquirer.prompt([
      {
        type: 'rawlist',
        name: 'choice',
        message: 'You already have Second Chance. Who do you want to give it to?',
        choices: targets,
      },
    ]);

    return action.choice;
  }

  async askMove(player) {
    // Choix du tour : piocher, regarder sa main, ou s'arrêter.
    const play = await inquirer.prompt([
      {
        type: 'rawlist',
        name: 'choice',
        message: `${player.name || `Player ${player.player_nb}`} - Choose your move:`,
        choices: ['Flip a card', 'Watch my card', 'Stop'],
      },
    ]);

    return play.choice;
  }

  async showHand(player) {
    console.log('Numbers:', player.hand_number);
    console.log('Bonus:', player.hand_bonus);
    console.log('Actions:', player.hand_actions);
  }

  showDraw(player, card) {
    console.log(`${player.name || `Player ${player.player_nb}`} drew:`, card);
  }

  showSeparator() {
    // Séparateur visuel entre les tours.
    console.log('');
    console.log('---');
    console.log('');
  }

  async showRoundSummary(players, scores) {
    // Résumé de fin de manche.
    console.log(`${BG_GREY}${FG_BLACK}--- Résumé du tour ---${RESET}`);
    for (const player of players) {
      const total = scores.get(player.player_nb) || 0;
      console.log(`${BG_GREY}${FG_BLACK}${player.name || `Player ${player.player_nb}`} : +${player.score} (total ${total})${RESET}`);
      await sleep.sleep(1);
    }
  }

  showWinner(player, totalScore) {
    console.log(`🏆 ${player.name || `Player ${player.player_nb}`} wins with ${totalScore} points!`);
  }

  async pause(seconds) {
    await sleep.sleep(seconds);
  }
}

export { Interface };
