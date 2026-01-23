import { Hand } from './player_hand.js';
import { Game } from './game.js';

class Match {
    constructor(playerCount) {
        this.playerCount = playerCount;
        this.players = [];
        this.scores = new Map();
        this.game = null;
        this.dealerIndex = 0;
        this.targetScore = 200;
        this.gameOver = false;
        this.initPlayers();
    }

    initPlayers() {
        for (let i = 1; i <= this.playerCount; i += 1) {
            const player = new Hand(i);
            this.players.push(player);
            this.scores.set(i, 0);
        }
        this.game = new Game(this.players);
    }

    startRound() {
        this.game.roundEnded = false;
        for (const player of this.players) {
            player.resetForNewRound();
        }
        this.dealInitialCards();
    }

    endRound() {
        for (const player of this.players) {
            if (player.state) {
                player.endGame(false);
            }
        }

        for (const player of this.players) {
            const total = this.scores.get(player.player_nb) || 0;
            this.scores.set(player.player_nb, total + player.score);
        }

        for (const player of this.players) {
            player.returnAllCardsToDeck();
        }

        this.game.roundEnded = false;
        this.advanceDealer();
        this.checkGameOver();
    }

    dealInitialCards() {
        const order = this.turnOrder();
        let safety = 0;
        while (this.players.some(p => p.totalCards() === 0) && !this.game.roundEnded) {
            for (const player of order) {
                if (player.totalCards() === 0) {
                    this.game.resolveDraw(player);
                    if (this.game.roundEnded) return;
                }
            }
            safety += 1;
            if (safety > 200) break;
        }
    }

    turnOrder() {
        const order = [];
        for (let i = 0; i < this.players.length; i += 1) {
            const index = (this.dealerIndex + i) % this.players.length;
            order.push(this.players[index]);
        }
        return order;
    }

    activePlayers() {
        return this.players.filter(p => p.state);
    }

    playerStay(player) {
        if (!player || !player.state) return;
        player.endGame(false);
    }

    advanceDealer() {
        this.dealerIndex = (this.dealerIndex + 1) % this.players.length;
    }

    checkGameOver() {
        for (const total of this.scores.values()) {
            if (total >= this.targetScore) {
                this.gameOver = true;
                return true;
            }
        }
        this.gameOver = false;
        return false;
    }
}

export { Match };
