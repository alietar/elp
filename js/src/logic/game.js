import { packet, cardsTypes, defausse } from './game_init.js';
import { draw_card } from './draw.js';

class Game {
    constructor(players = []) {
        this.players = players;
        this.roundEnded = false;
    }

    drawCard() {
        // Si le paquet est vide, on le recharge depuis la défausse.
        if (cardsTypes.length === 0) {
            this.refillDeckFromDefausse();
        }
        return draw_card(packet, cardsTypes);
    }

    isActionCard(card) {
        // Cartes d'action "spéciales" qui ne vont pas dans la main.
        return card === 'Flip Three' || card === 'Freeze' || card === 'Second Chance';
    }

    resolveDraw(player) {
        const card = this.drawCard();
        if (!card) return null;

        // Une action ne va pas dans la main : on la gère plus tard.
        if (this.isActionCard(card)) {
            return { card, type: 'action' };
        }

        const addResult = player.addCard(card);
        if (addResult === 'second_chance') {
            // La carte est ignorée grâce à Second Chance, le tour s'arrête ici.
            return { card, type: 'second_chance' };
        }
        if (player.flip7) this.roundEnded = true;
        return { card, type: 'normal' };
    }

    resolveAction(card, currentPlayer, targetPlayer) {
        if (!currentPlayer) return;

        if (card === 'Second Chance') {
            if (targetPlayer && targetPlayer !== currentPlayer) {
                if (!targetPlayer.hasSecondChance()) {
                    // Donne Second Chance au joueur ciblé.
                    targetPlayer.receiveSecondChance(this.players);
                } else {
                    // Cible invalide (déjà Second Chance) => défausse.
                    this.discardCard('Second Chance');
                }
            } else {
                // Pas de cible : le joueur courant récupère Second Chance.
                const target = currentPlayer;
                target.receiveSecondChance(this.players);
            }
            return;
        }

        if (card === 'Freeze') {
            const target = targetPlayer || this.chooseTarget(currentPlayer);
            if (target) {
                // Freeze est immédiat et ne reste jamais en main.
                target.frozen();
            }
            this.discardCard('Freeze');
            return;
        }

        if (card === 'Flip Three') {
            const target = targetPlayer || this.chooseTarget(currentPlayer);
            if (target) {
                // On pioche 3 cartes pour la cible, puis on traite les actions trouvées.
                const result = this.resolveFlipThree(target);
                this.discardCard('Flip Three');
                if (result && result.pendingActions) {
                    // actionOwner = joueur qui a réellement tiré l'action pendant le Flip Three.
                    return { ...result, actionOwner: target };
                }
                return result;
            } else {
                // Pas de cible => on défausse l'action.
                this.discardCard('Flip Three');
            }
        }
    }

    resolveFlipThree(currentPlayer) {
        const pendingActions = [];
        const drawnCards = [];

        for (let i = 0; i < 3; i += 1) {
            const card = this.drawCard();
            if (!card) break;
            drawnCards.push(card);

            const playerLabel = currentPlayer.name || `Player ${currentPlayer.player_nb}`;

            console.log(`${playerLabel} draws (Flip Three):`, card);

            if (this.isActionCard(card)) {
                if (card === 'Second Chance') {
                    // Second Chance se résout immédiatement.
                    currentPlayer.receiveSecondChance(this.players);
                } else {
                    // Freeze / Flip Three seront résolus après la pioche des 3 cartes.
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

        // actionOwner permet à l'UI de savoir qui doit choisir la cible.
        return { pendingActions, actionOwner: currentPlayer, drawnCards };
    }

    chooseTarget(currentPlayer) {
        //Si l'UI ne choisit pas, on prend un joueur actif.
        const activePlayers = this.players.filter(p => p.state);
        if (activePlayers.length === 0) return currentPlayer;
        if (activePlayers.length === 1) return activePlayers[0];
        const other = activePlayers.find(p => p !== currentPlayer);
        return other || activePlayers[0];
    }

    refillDeckFromDefausse() {
        // Recharge toutes les cartes de la défausse dans le paquet.
        for (const [key, value] of defausse.entries()) {
            if (value.quantity > 0) {
                packet.set(key, { type: value.type, quantity: value.quantity });
                if (!cardsTypes.includes(key)) cardsTypes.push(key);
                value.quantity = 0;
            }
        }
    }

    discardCard(card) {
        // Ajoute une carte jouée dans la défausse.
        const cardData = defausse.get(card);
        if (cardData) {
            cardData.quantity += 1;
        }
    }
}

export { Game };
