import inquirer from 'inquirer';
import {draw_card} from './draw.js';
import {packet, cardsTypes} from './game_init.js';

var playersNames = []

async function name_player() {
  const pseudo = await inquirer.prompt([
    {
      type: 'input',
      name: 'playerName',
      message: "What's your name?",
    }
  ]);

  console.log("Hi",pseudo.playerName,", welcome to the game!");
  playersNames.push(pseudo.playerName)
  
}


async function play() {
  const play = await inquirer.prompt([
    {
      type: 'rawlist',
      name: 'choice',
      message: 'Chose your move :',
      choices: ["Flip a card", "Stop"]
    }
  ]);

  if (play.choice === "Flip a card") {
    const card = draw_card(packet, cardsTypes)
    console.log("You drew a", card);
    if (packet.get(card).type === 'action' && card !== 'Second Chance'){
      await action()
    }
  } else {
    console.log("You won", "" ,"points");
    
  }
}

async function action() {
  const action = await inquirer.prompt([
    {
      type: 'rawlist',
      name: 'choice',
      message: 'You drew an action card, who do you want to target?',
      choices: playersNames
    }
  ]);

  if (action.choice === playersNames[0]) {
    console.log("You chose", playersNames[0]);
  } else if (action.choice === playersNames[1]) {
    console.log("You chose", playersNames[1]);
    
  } else {
    console.log("")
  }
}

await name_player()
await play()
