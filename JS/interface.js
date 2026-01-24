import inquirer from 'inquirer';
import { draw_card } from './draw.js';
import { packet, cardsTypes } from './game_init.js';

class Interface {
  constructor() {
    this.playersNames = [];
  }

  async namePlayer() {
    const pseudo = await inquirer.prompt([
      {
        type: 'input',
        name: 'playerName',
        message: "What's your name?",
      },
    ]);

    console.log('Hi', pseudo.playerName, ', welcome to the game!');
    this.playersNames.push(pseudo.playerName);
  }

  async play() {
    const play = await inquirer.prompt([
      {
        type: 'rawlist',
        name: 'choice',
        message: 'Chose your move :',
        choices: ['Flip a card', 'Stop'],
      },
    ]);

    if (play.choice === 'Flip a card') {
      const card = draw_card(packet, cardsTypes);
      console.log('You drew a', card);
      if (packet.get(card).type === 'action' && card !== 'Second Chance') {
        await this.action();
      }
    } else {
      console.log('You won', '', 'points');
    }
  }

  async action() {
    const action = await inquirer.prompt([
      {
        type: 'rawlist',
        name: 'choice',
        message: 'You drew an action card, who do you want to target?',
        choices: this.playersNames,
      },
    ]);

    if (action.choice === this.playersNames[0]) {
      console.log('You chose', this.playersNames[0]);
    } else if (action.choice === this.playersNames[1]) {
      console.log('You chose', this.playersNames[1]);
    } else {
      console.log('');
    }
  }

  async start() {
    await this.namePlayer();
    await this.play();
  }
}

const ui = new Interface();
await ui.start();

export { Interface };
