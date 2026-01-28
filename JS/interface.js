import inquirer from 'inquirer';
import colors from 'colors';
import { packet } from './game_init.js';

// Couleurs simples pour l'affichage console.
const BG_GREY = "\x1b[47m";
const FG_BLACK = "\x1b[30m";
const RESET = "\x1b[0m";  


class Interface {
  
  getColoredCard(cardName) { // mettre les numéros en couleur comme dans le jeu

    const cardData = packet.get(cardName);
    if (!cardData) return cardName;

    switch (cardData.type) {
      case 'number':
        if (cardName == '0' || cardName === '3' || cardName === '6' || cardName === '7') return colors.magenta(cardName);
        else if (cardName === '1' || cardName === '12' ) return colors.grey(cardName);
        else if (cardName === '2' || cardName === '5' || cardName === '8') return colors.green(cardName);
        else if (cardName === '4' || cardName === '11') return colors.blue(cardName);
        else if (cardName === '9') return colors.yellow(cardName);
        else return colors.red(cardName);
      
      case 'modifier':
        return colors.red.bold(cardName);
      
      case 'action':
        return colors.red.bold(cardName);
      
      default:
        return cardName;
    }
  }


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
    const numbers = player.hand_number.map(c => this.getColoredCard(c)).join(', ');
    const bonus = player.hand_bonus.map(c => this.getColoredCard(c)).join(', ');
    const actions = player.hand_actions.map(c => this.getColoredCard(c)).join(', ');

    console.log('Numbers:', numbers);
    console.log('Bonus:', bonus);
    console.log('Actions:', actions);
  }

  showDraw(player, card) {

    // On calcule la longueur de la carte (ex: "10" = 2, "Freeze" = 6)
    // On utilise card.toString() pour ne pas compter les caractères invisibles des couleurs
    const cardLength = card.toString().length;
    
    // On définit la variable dashes ICI
    const dashes = '-'.repeat(cardLength + 2);

    console.log(`${player.name || `Player ${player.player_nb}`} drew:`);
    console.log(`  +${dashes}+`);
    console.log(`  | ${this.getColoredCard(card)} |`);
    console.log(`  +${dashes}+`);
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
      await this.pause(1);   
    }
  }

  showWinner(player, totalScore) {
    console.log(`🏆 ${player.name || `Player ${player.player_nb}`} wins with ${totalScore} points!`);
  }

  async pause(seconds) {
    return new Promise(resolve => setTimeout(resolve, seconds * 1000));
  }
}

export { Interface };
