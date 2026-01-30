import { Hand } from './player_hand.js';
import { Game } from './game.js';

class Match {
    constructor(playerCount, ui = null) {
        this.playerCount = playerCount;
        this.players = [];
        this.scores = new Map();
        this.game = null;
        this.dealerIndex = 0;
        this.targetScore = 200;
        this.gameOver = false;
        this.ui = ui;
        this.initPlayers();
    }

    initPlayers() {
        // Initialise les joueurs et leurs scores.
        for (let i = 1; i <= this.playerCount; i += 1) {
            const player = new Hand(i);
            this.players.push(player);
            this.scores.set(i, 0);
        }
        this.game = new Game(this.players);
    }

    startRound() {
        // Nouvelle manche : reset des joueurs et du flag de fin de manche.
        this.game.roundEnded = false;
        for (const player of this.players) {
            player.resetForNewRound();
        }
    }

    endRound() {
        // Termine la manche et cumule les scores.
        for (const player of this.players) {
            if (player.state) {
                player.endGame(false);
            } else if (player.eliminatedByDuplicate) {
                player.score = 0;
                player.flip7 = false;
            }
        }

        for (const player of this.players) {
            const total = this.scores.get(player.player_nb) || 0;
            this.scores.set(player.player_nb, total + player.score);
        }

        for (const player of this.players) {
            player.returnAllCardsToDeck();
        }
        console.log('Cartes mises à la défausse.');

        this.game.roundEnded = false;
        this.advanceDealer();
        this.checkGameOver();
    }


    activePlayers() {
        return this.players.filter(p => p.state);
    }

    playTurn(player, action, targetPlayer = null) {
        if (!player || !player.state) return null;

        if (action === 'flip') {
            const result = this.game.resolveDraw(player);
            if (result && result.type === 'action') {
                // Second Chance : si le joueur en a déjà une, il doit choisir un autre joueur.
                if (result.card === 'Second Chance' && player.hasSecondChance()) {
                    return { ...result, needsTarget: true, secondChanceTarget: true };
                }
                // Autres actions : si aucune cible fournie, on la demandera à l'UI.
                if (result.card !== 'Second Chance' && !targetPlayer) {
                    return { ...result, needsTarget: true };
                }
                this.game.resolveAction(result.card, player, targetPlayer);
                return result;
            }
            return result;
        }

        if (action === 'stop') {
            this.playerStay(player);
            return { card: null, type: 'stop' };
        }

        return { card: null, type: 'watch' };
    }

    async resolveAction(card, currentPlayer, targetPlayer) {
        // Pour Second Chance, on demande la cible si nécessaire.
        if (card === 'Second Chance' && currentPlayer.hasSecondChance() && !targetPlayer) {
            targetPlayer = await this.ui.chooseSecondChanceTarget(currentPlayer, this.players);
        }
        const result = this.game.resolveAction(card, currentPlayer, targetPlayer);
        // Si un Flip Three a généré des actions, on les résout après les 3 cartes.
        if (!result || !result.pendingActions || result.pendingActions.length === 0) return result;
        if (this.game.roundEnded) return result;

        // actionOwner = joueur qui a tiré les actions via Flip Three.
        const actionOwner = result.actionOwner || currentPlayer;
        for (const action of result.pendingActions) {
            if (action === 'Second Chance') {
                const target = actionOwner.hasSecondChance()
                    ? await this.ui.chooseSecondChanceTarget(actionOwner, this.players)
                    : null;
                await this.resolveAction(action, actionOwner, target);
            } else {
                const target = await this.ui.chooseTarget(actionOwner, this.players);
                await this.resolveAction(action, actionOwner, target);
            }
        }
        return result;
    }

    playerStay(player) {
        // Le joueur choisit "Stop".
        if (!player || !player.state) return;
        player.endGame(false);
    }

    advanceDealer() {
        this.dealerIndex = (this.dealerIndex + 1) % this.players.length;
    }

    checkGameOver() {
        // La partie se termine quand un joueur atteint le score cible.
        for (const total of this.scores.values()) {
            if (total >= this.targetScore) {
                this.gameOver = true;
                return true;
            }
        }
        this.gameOver = false;
        return false;
    }

    getWinner() {
        for (const player of this.players) {
            const total = this.scores.get(player.player_nb) || 0;
            if (total >= this.targetScore) return player;
        }
        return null;
    }
}

export { Match };
