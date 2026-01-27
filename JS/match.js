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
                if (result.card === 'Second Chance' && player.hasSecondChance()) {
                    return { ...result, needsTarget: true, secondChanceTarget: true };
                }
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
        if (card === 'Second Chance' && currentPlayer.hasSecondChance() && !targetPlayer) {
            targetPlayer = await this.ui.chooseSecondChanceTarget(currentPlayer, this.players);
        }
        const result = this.game.resolveAction(card, currentPlayer, targetPlayer);
        if (result && result.pendingActions && result.pendingActions.length > 0 && !this.game.roundEnded) {
            for (const action of result.pendingActions) {
                if (action === 'Second Chance') {
                    const target = currentPlayer.hasSecondChance()
                        ? await this.ui.chooseSecondChanceTarget(currentPlayer, this.players)
                        : null;
                    await this.resolveAction(action, currentPlayer, target);
                } else {
                    const target = await this.ui.chooseTarget(currentPlayer, this.players);
                    await this.resolveAction(action, currentPlayer, target);
                }
            }
        }
        return result;
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

    getWinner() {
        for (const player of this.players) {
            const total = this.scores.get(player.player_nb) || 0;
            if (total >= this.targetScore) return player;
        }
        return null;
    }

    async startGame() {
        if (!this.ui) throw new Error('UI is required to start the game.');

        for (const player of this.players) {
            player.name = await this.ui.askPlayerName(player);
        }

        while (!this.gameOver) {
            this.startRound();
            await this.playRound();
            this.endRound();
            await this.ui.showRoundSummary(this.players, this.scores);
            await this.ui.pause(1);

            if (this.gameOver) {
                const winner = this.getWinner();
                if (winner) {
                    this.ui.showWinner(winner, this.scores.get(winner.player_nb));
                }
            }
        }
    }

    async playRound() {
        let index = this.dealerIndex;
        const totalPlayers = this.players.length;

        while (!this.game.roundEnded && this.activePlayers().length > 0) {
            const player = this.players[index % totalPlayers];
            if (!player.state) {
                index += 1;
                continue;
            }

            await this.ui.pause(1);
            this.ui.showSeparator();

            const choice = await this.ui.askMove(player);

            if (choice === 'Flip a card') {
                const result = this.playTurn(player, 'flip');
                if (result) {
                    this.ui.showDraw(player, result.card);
                    if (result.type === 'action' && result.needsTarget) {
                        if (result.secondChanceTarget) {
                            const target = await this.ui.chooseSecondChanceTarget(player, this.players);
                            await this.resolveAction(result.card, player, target);
                        } else {
                            const target = await this.ui.chooseTarget(player, this.players);
                            await this.resolveAction(result.card, player, target);
                        }
                    }
                }
                index += 1;
                continue;
            }

            if (choice === 'Watch my card') {
                await this.ui.showHand(player);
                await this.ui.pause(1);
                continue;
            }

            this.playTurn(player, 'stop');
            index += 1;
        }
    }
}

export { Match };
