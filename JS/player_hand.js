class Hand {
    constructor(player_nb){
        this.player_nb = player_nb;
        this.hand_number = [];
        this.hand_bonus = [];
        this.hand_actions = [];
        this.state = true;
        this.score = 0;
    }

    addCard(card) {
    if (!this.state) return;

    if (this.isCardAction(card)) {
        this.hand_actions.push(card);
    } else if (this.isCardNumber(card)) {
        this.hand_number.push(card);
        this.checkForDuplicates(); // on verifie que la carte qu'on ajoute n'est pas un doublon
        if (this.state) this.checkWin(); // on verifie qu'on a gagné 
    } else {
        this.hand_bonus.push(card);
    }
    }


    isCardNumber(card) {
        return !isNaN(parseInt(card, 10));
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


    


    isCardAction(card) {
        const actionCards = ['Flip Three', 'Freeze', 'Second Chance']; // Liste des cartes d'action
        return actionCards.includes(card);
    }



    showHand() {
        console.log(`Main du joueur ${this.player_nb}:`, this.hand_number);
    }

    checkForDuplicates() {
    const unique = new Set(this.hand_number);

    if (unique.size !== this.hand_number.length) {
        this.state = false;
        this.score = 0;
        console.log(`❌ Doublon détecté : Joueur ${this.player_nb} éliminé (score = 0)`);
    }
    }

    endGame() {
        if (!this.state) return;

        this.state = false;
        this.score = this.pointInMyHand() + this.pointInBonus();
    }

    addToDeck(card) {
        console.log(`🃏 Carte ${card} remise dans la pioche`);
    } // a rajouté à remettre carte dans le jeu 


    playActionCard(card) {
        if (!this.state) return;

        if (!this.isCardAction(card)) return;

        const index = this.hand_actions.indexOf(card);
        if (index === -1) return;

        this.hand_actions.splice(index, 1);
        this.addToDeck(card);
    }

    checkWin() {
        if (this.hand_number.length === 7) {
            this.state = false;
            this.score = this.pointInMyHand() + this.pointInBonus();
            console.log(`🎉 FLIP 7 ! Joueur ${this.player_nb} gagne avec ${this.score} points`);
        }
}









    
}


