import { packet, cardsTypes, defausse } from './game_init.js';
import { draw_card } from './draw.js';

class Game {
    constructor(players = []) {
        this.players = players;
        this.roundEnded = false;
    }

    drawCard() {
        if (cardsTypes.length === 0) {
            this.refillDeckFromDefausse();
        }
        return draw_card(packet, cardsTypes);
    }

    isActionCard(card) {
        return card === 'Flip Three' || card === 'Freeze' || card === 'Second Chance';
    }

    resolveDraw(player) {
        const card = this.drawCard();
        if (!card) return null;

        if (this.isActionCard(card)) {
            return { card, type: 'action' };
        }

        const addResult = player.addCard(card);
        if (addResult === 'second_chance') {
            return this.resolveDraw(player);
        }
        if (player.flip7) this.roundEnded = true;
        return { card, type: 'normal' };
    }

    resolveAction(card, currentPlayer, targetPlayer) {
        if (!currentPlayer) return;

        if (card === 'Second Chance') {
            if (targetPlayer && targetPlayer !== currentPlayer) {
                if (!targetPlayer.hasSecondChance()) {
                    targetPlayer.receiveSecondChance([targetPlayer]);
                } else {
                    this.discardCard('Second Chance');
                }
            } else {
                const target = currentPlayer;
                target.receiveSecondChance(this.players);
            }
            return;
        }

        if (card === 'Freeze') {
            const target = targetPlayer || this.chooseTarget(currentPlayer);
            if (target) {
                target.addActionCard('Freeze');
                target.frozen();
            } else {
                this.discardCard('Freeze');
            }
            return;
        }

        if (card === 'Flip Three') {
            const target = targetPlayer || this.chooseTarget(currentPlayer);
            if (target) {
                target.addActionCard('Flip Three');
                return this.resolveFlipThree(target);
            } else {
                this.discardCard('Flip Three');
            }
        }
    }

    resolveFlipThree(currentPlayer) {
        const pendingActions = [];

        for (let i = 0; i < 3; i += 1) {
            const card = this.drawCard();
            if (!card) break;

            const playerLabel = currentPlayer.name || `Player ${currentPlayer.player_nb}`;

            console.log(`${playerLabel} draws (Flip Three):`, card);

            if (this.isActionCard(card)) {
                if (card === 'Second Chance') {
                    currentPlayer.receiveSecondChance(this.players);
                } else {
                    pendingActions.push(card);
                }
            } else {
                currentPlayer.addCard(card);
                if (currentPlayer.flip7) {
                    this.roundEnded = true;
                }
            }

            if (this.roundEnded) break;
        }

        if (this.roundEnded) return null;

        return { pendingActions };
    }

    chooseTarget(currentPlayer) { // reecrire quand on aura fait l'interface et ca retournera le joueur
        const activePlayers = this.players.filter(p => p.state);
        if (activePlayers.length === 0) return currentPlayer;
        if (activePlayers.length === 1) return activePlayers[0];
        const other = activePlayers.find(p => p !== currentPlayer);
        return other || activePlayers[0];
    }

    refillDeckFromDefausse() {
        for (const [key, value] of defausse.entries()) {
            if (value.quantity > 0) {
                packet.set(key, { type: value.type, quantity: value.quantity });
                if (!cardsTypes.includes(key)) cardsTypes.push(key);
                value.quantity = 0;
            }
        }
    }

    discardCard(card) {
        const cardData = defausse.get(card);
        if (cardData) {
            cardData.quantity += 1;
        }
    }
}

export { Game };
