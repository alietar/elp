import { defausse } from './game_init.js';

class Hand {
    constructor(player_nb){
        this.player_nb = player_nb;
        this.hand_number = [];
        this.hand_bonus = [];
        this.hand_actions = [];
        this.state = true;
        this.score = 0;
        this.flip7 = false;
        this.eliminated = false;
        this.frozen_state = false;
        this.stayed = false;
    }

    addCard(card) {
    if (!this.state) return;

    if (this.isCardNumber(card)) {
        if (this.hand_number.includes(card)) {
            if (this.useSecondChanceIfAvailable()) {
                this.addToDeck(card);
                return;
            }
        }
        this.hand_number.push(card);
        if (this.checkForDuplicates()) {
            this.hand_number.pop();
            this.addToDeck(card);
            return;
        }
        if (this.state) this.checkWin(); // on verifie qu'on a gagné 
    } else {
        this.hand_bonus.push(card);
    }
    }

    showHand() {
        console.log(`Main du joueur ${this.player_nb}:`, this.hand_number);
    }

    totalCards() {
        return this.hand_number.length + this.hand_bonus.length + this.hand_actions.length;
    }

    addActionCard(card) {
        this.hand_actions.push(card);
    }

    hasSecondChance() {
        return this.hand_actions.includes('Second Chance');
    }

    receiveSecondChance(players = []) {
        if (!this.hasSecondChance()) {
            this.hand_actions.push('Second Chance');
            return;
        }

        const target = players.find(p => p !== this && p.state && !p.hasSecondChance());
        if (target) {
            target.hand_actions.push('Second Chance');
        } else {
            this.addToDeck('Second Chance');
        }
    }

    isCardNumber(card) {
        return !isNaN(parseInt(card, 10));
    }

    checkForDuplicates() {
    const unique = new Set(this.hand_number);

    if (unique.size !== this.hand_number.length) {
        this.state = false;
        this.score = 0;
        this.eliminated = true;
        console.log(`Doublon détecté : Joueur ${this.player_nb} éliminé (score = 0)`);
        return true;
    }
    return false;
    }

    useSecondChanceIfAvailable() {
        const index = this.hand_actions.indexOf('Second Chance');
        if (index === -1) return false;
        this.hand_actions.splice(index, 1);
        this.addToDeck('Second Chance');
        console.log(`✨ Second Chance utilisée : doublon ignoré pour Joueur ${this.player_nb}`);
        return true;
    }

    pointInMyHand() {
        let sum = 0;

        for (let card of this.hand_number) {
            if (this.isCardNumber(card)) {
                sum += parseInt(card, 10);
            }
        }

        return sum;
    }

    pointInBonus() {
        let sum = 0;

        for (let card of this.hand_bonus) {
            const value = parseInt(card.replace("+", ""), 10);
            if (!isNaN(value)) {
                sum += value;
            }
        }

        return sum;
    }

    hasX2Bonus() {
        return this.hand_bonus.includes('x2');
    }

    checkWin() {
        if (this.hand_number.length === 7) {
            this.state = false;
            this.flip7 = true;
            const numberScore = this.pointInMyHand();
            const bonusScore = this.pointInBonus();
            this.score = (this.hasX2Bonus() ? numberScore * 2 : numberScore) + bonusScore + 15;
            console.log(` FLIP 7 ! Joueur ${this.player_nb} gagne avec ${this.score} points`);
        }
    }

    endGame(isFreeze = false) {
        if (!this.state) return;

        this.state = false;
        if (isFreeze) {
            this.score = 0;
            this.frozen_state = true;
        } else {
            const numberScore = this.pointInMyHand();
            const bonusScore = this.pointInBonus();
            this.score = (this.hasX2Bonus() ? numberScore * 2 : numberScore) + bonusScore;
            this.stayed = true;
        }
    }

    frozen() {
        this.endGame(true);
    }

    addToDeck(card) {
        const cardData = defausse.get(card);
        if (cardData) {
            cardData.quantity += 1; // ajout dans la défausse de +1 à la carte action joué
        }
        console.log(`🃏 Carte ${card} remise dans la pioche`);
    } 

    returnAllCardsToDeck() {
        for (const card of this.hand_number) {
            this.addToDeck(card);
        }
        for (const card of this.hand_bonus) {
            this.addToDeck(card);
        }
        for (const card of this.hand_actions) {
            this.addToDeck(card);
        }
        this.hand_number = [];
        this.hand_bonus = [];
        this.hand_actions = [];
    }

    resetForNewRound() {
        this.hand_number = [];
        this.hand_bonus = [];
        this.hand_actions = [];
        this.state = true;
        this.score = 0;
        this.flip7 = false;
        this.eliminated = false;
        this.frozen_state = false;
        this.stayed = false;
    }

}

export { Hand };
