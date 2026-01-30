// game_ui.js
import React, { useState, useEffect, useMemo } from 'react';
import { Box, Text, Newline } from 'ink';
import SelectInput from 'ink-select-input';
import { Match } from '../logic/match.js'; // Ta classe existante
import { CardVisual, CardText, Header } from './components.js';

export const GameController = ({ playerCount, playerNames, onGameOver }) => {
    // 1. Initialisation unique du Match
    const match = useMemo(() => {
        const m = new Match(playerCount, null); // Pas d'UI classique
        // On injecte les noms manuellement car on bypass match.ui
        m.players.forEach((p, i) => p.name = playerNames[i]);
        m.startRound();
        return m;
    }, [playerCount, playerNames]);

    // 2. États React pour suivre le jeu
    const [currentPlayerIndex, setCurrentPlayerIndex] = useState(match.dealerIndex);
    const [lastDrawnCard, setLastDrawnCard] = useState(null);
    const [message, setMessage] = useState("La manche commence !");
    const [viewState, setViewState] = useState('menu'); // 'menu', 'hand', 'target_selection', 'summary'
    const [forceUpdate, setForceUpdate] = useState(0); // Hack pour forcer le render quand l'objet match change

    // Helper pour récupérer le joueur actuel
    const currentPlayer = match.players[currentPlayerIndex];

    // --- LOGIQUE DE JEU ---

    const nextPlayer = () => {
        let nextIndex = (currentPlayerIndex + 1) % match.players.length;
        let loops = 0;
        
        // Trouver le prochain joueur actif (state = true)
        while (!match.players[nextIndex].state && loops < match.players.length) {
            nextIndex = (nextIndex + 1) % match.players.length;
            loops++;
        }

        // Si tout le monde a fini ou manche terminée
        if (loops === match.players.length || match.game.roundEnded || match.activePlayers().length === 0) {
            endRound();
        } else {
            setLastDrawnCard(null);
            setMessage(`C'est au tour de ${match.players[nextIndex].name}`);
            setCurrentPlayerIndex(nextIndex);
            setViewState('menu');
        }
    };

    const endRound = () => {
        match.endRound();
        setViewState('summary');
    };

    const handleAction = (item) => {
        if (item.value === 'view_hand') {
            setViewState('hand');
            return;
        }

        if (item.value === 'stop') {
            match.playTurn(currentPlayer, 'stop');
            setMessage(`${currentPlayer.name} s'arrête.`);
            setTimeout(nextPlayer, 1500);
            return;
        }

        if (item.value === 'draw') {
            // Utilisation de ta logique match.js
            const result = match.playTurn(currentPlayer, 'flip');
            
            if (result && result.card) {
                setLastDrawnCard(result.card);
                
                if (result.type === 'action' && result.needsTarget) {
                    setMessage(`Action ${result.card} ! Choisir une cible.`);
                    // TODO: Implémenter la sélection de cible (simplifié ici pour l'exemple)
                    // Pour l'instant, on applique l'action sans cible ou auto-cible pour ne pas bloquer
                    match.game.resolveAction(result.card, currentPlayer, null); 
                } 
                
                // Vérifier si le joueur a sauté (doublon)
                if (!currentPlayer.state) {
                    setMessage(`DOUBLON ! ${currentPlayer.name} est éliminé de la manche.`);
                    setTimeout(nextPlayer, 2000);
                } else if (match.game.roundEnded) {
                     setMessage(`Fin de manche !`);
                     setTimeout(endRound, 2000);
                } else {
                    setMessage(`${currentPlayer.name} a pioché.`);
                    setForceUpdate(n => n + 1); // Rafraichir l'affichage
                }
            }
        }
    };

    // --- VUES ---

    if (viewState === 'summary') {
        // Vérifier victoire globale
        const winner = match.getWinner();
        if (winner) {
            return (
                <Box flexDirection="column" alignItems="center" borderColor="green" borderStyle="double" padding={2}>
                    <Text bold color="green"> VICTOIRE FINALE </Text>
                    <Text>Le vainqueur est {winner.name} avec {match.scores.get(winner.player_nb)} points !</Text>
                </Box>
            );
        }

        return (
            <Box flexDirection="column">
                <Header title="Fin de la manche" />
                {match.players.map(p => (
                    <Box key={p.player_nb} justifyContent="space-between" width={40}>
                        <Text>{p.name}</Text>
                        <Text bold>Total: {match.scores.get(p.player_nb)} (+{p.score})</Text>
                    </Box>
                ))}
                <Box marginTop={2}>
                    <Text color="gray">Appuyez sur Entrée pour la manche suivante...</Text>
                    <SelectInput 
                        items={[{label: 'Manche suivante', value: 'next'}]} 
                        onSelect={() => {
                            match.startRound();
                            setCurrentPlayerIndex(match.dealerIndex);
                            setViewState('menu');
                            setMessage("Nouvelle manche !");
                        }} 
                    />
                </Box>
            </Box>
        );
    }

    if (viewState === 'hand') {
        return (
            <Box flexDirection="column">
                <Header title={`Main de ${currentPlayer.name}`} />
                <Text>Nombres: {currentPlayer.hand_number.map(c => c).join(', ')}</Text>
                <Text>Bonus: {currentPlayer.hand_bonus.join(', ')}</Text>
                <Text>Actions: {currentPlayer.hand_actions.join(', ')}</Text>
                <Box marginTop={1}>
                    <SelectInput items={[{label: 'Retour', value: 'back'}]} onSelect={() => setViewState('menu')} />
                </Box>
            </Box>
        );
    }

    // Vue Principale (Menu de jeu)
    const menuItems = [
        { label: 'Piocher une carte', value: 'draw' },
        { label: 'Voir ma main', value: 'view_hand' },
        { label: 'S\'arrêter', value: 'stop' }
    ];

    return (
        <Box flexDirection="column" width="100%">
            <Box justifyContent="space-between" width="100%">
                <Text>Joueur: <Text bold color="cyan">{currentPlayer.name}</Text></Text>
                <Text>Score manche: {currentPlayer.pointInMyHand()}</Text>
            </Box>
            
            <Box borderStyle="single" padding={1} marginY={1} flexDirection="column" alignItems="center">
                <Text italic color="gray">{message}</Text>
                {lastDrawnCard && (
                    <Box marginTop={1}>
                        <CardVisual card={lastDrawnCard} />
                    </Box>
                )}
            </Box>

            <Text>Que voulez-vous faire ?</Text>
            <SelectInput items={menuItems} onSelect={handleAction} />
        </Box>
    );
};