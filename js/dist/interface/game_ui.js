// game_ui.js
import React, { useState, useEffect, useMemo } from 'react';
import { Box, Text, Newline } from 'ink';
import SelectInput from 'ink-select-input';
import { Match } from '../logic/match.js'; // Ta classe existante
import { doIHaveToDraw } from '../logic/helper.js';
import { packet } from '../logic/game_init.js';
import { CardVisual, CardText, Header } from './components.js';
export const GameController = ({
  playerCount,
  playerNames,
  onGameOver
}) => {
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
  const [lastDrawnCards, setLastDrawnCards] = useState([]);
  const [message, setMessage] = useState("La manche commence !");
  const [viewState, setViewState] = useState('menu'); // 'menu', 'hand', 'target_selection', 'summary'
  const [forceUpdate, setForceUpdate] = useState(0); // Hack pour forcer le render quand l'objet match change
  const [pendingAction, setPendingAction] = useState(null); // { card, owner, targets }
  const [pendingQueue, setPendingQueue] = useState([]); // Actions en attente (Flip Three)

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
      setLastDrawnCards([]);
      setMessage(`C'est au tour de ${match.players[nextIndex].name}`);
      setCurrentPlayerIndex(nextIndex);
      setViewState('menu');
    }
  };
  const endRound = () => {
    match.endRound();
    setViewState('summary');
  };
  const finalizeAfterAction = () => {
    const activePlayer = match.players[currentPlayerIndex];
    if (!activePlayer.state) {
      setMessage(`💥 DOUBLON ! ${activePlayer.name} est éliminé de la manche.`);
      setTimeout(nextPlayer, 2000);
    } else if (match.game.roundEnded) {
      setMessage(`🎉 FLIP 7 ou Fin de manche !`);
      setTimeout(endRound, 2000);
    } else {
      setMessage(`${activePlayer.name} a pioché.`);
      setForceUpdate(n => n + 1);
      setTimeout(nextPlayer, 1500);
    }
  };
  const getTargetCandidates = (card, owner) => {
    const activeTargets = match.players.filter(p => p.state);
    if (card === 'Second Chance' && owner.hasSecondChance()) {
      const eligible = activeTargets.filter(p => p !== owner && !p.hasSecondChance());
      if (eligible.length <= 1) {
        return {
          needsSelection: false,
          targets: eligible.length === 1 ? eligible : [null]
        };
      }
      return {
        needsSelection: true,
        targets: eligible
      };
    }
    if (activeTargets.length <= 1) {
      return {
        needsSelection: false,
        targets: activeTargets.length === 1 ? activeTargets : [owner]
      };
    }
    return {
      needsSelection: true,
      targets: activeTargets
    };
  };
  const resolveActionChain = (card, owner, target) => {
    const enqueuePending = (result, actionOwner, queue) => {
      if (result && result.pendingActions && result.pendingActions.length > 0 && !match.game.roundEnded) {
        return [...result.pendingActions.map(pendingCard => ({
          card: pendingCard,
          owner: actionOwner
        })), ...queue];
      }
      return queue;
    };
    const step = queue => {
      if (queue.length === 0 || match.game.roundEnded) {
        setPendingQueue([]);
        finalizeAfterAction();
        return;
      }
      const [next, ...rest] = queue;
      const targetInfo = getTargetCandidates(next.card, next.owner);
      if (targetInfo.needsSelection) {
        setPendingQueue(rest);
        setPendingAction({
          card: next.card,
          owner: next.owner,
          targets: targetInfo.targets
        });
        setViewState('target_selection');
        return;
      }
      const autoTarget = targetInfo.targets[0] ?? null;
      const result = match.game.resolveAction(next.card, next.owner, autoTarget);
      if (result && result.drawnCards && result.drawnCards.length > 0) {
        setLastDrawnCards(result.drawnCards);
        const nextQueue = enqueuePending(result, result.actionOwner || next.owner, rest);
        setTimeout(() => step(nextQueue), 1000);
        return;
      }
      const nextQueue = enqueuePending(result, result && result.actionOwner ? result.actionOwner : next.owner, rest);
      step(nextQueue);
    };
    const initialResult = match.game.resolveAction(card, owner, target);
    if (initialResult && initialResult.drawnCards && initialResult.drawnCards.length > 0) {
      setLastDrawnCards(initialResult.drawnCards);
      const initialQueue = enqueuePending(initialResult, initialResult.actionOwner || owner, pendingQueue);
      setTimeout(() => step(initialQueue), 1000);
      return;
    }
    const startQueue = enqueuePending(initialResult, initialResult && initialResult.actionOwner ? initialResult.actionOwner : owner, pendingQueue);
    step(startQueue);
  };
  const startActionFlow = (card, owner) => {
    const targetInfo = getTargetCandidates(card, owner);
    if (targetInfo.needsSelection) {
      setPendingAction({
        card,
        owner,
        targets: targetInfo.targets
      });
      setTimeout(() => setViewState('target_selection'), 1200);
      return;
    }
    const autoTarget = targetInfo.targets[0] ?? null;
    resolveActionChain(card, owner, autoTarget);
  };
  const handleAction = item => {
    if (item.value === 'view_hand') {
      setViewState('hand');
      return;
    }
    if (item.value === 'helper') {
      const proba = doIHaveToDraw(packet, currentPlayer.hand_number);
      const percent = (proba * 100).toFixed(2);
      setMessage(`${currentPlayer.name} - Risque de doublon : ${percent}%`);
      setForceUpdate(n => n + 1);
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
        setLastDrawnCards([result.card]);
        if (result.type === 'action' && result.needsTarget) {
          setMessage(`Action ${result.card} ! Choisir une cible.`);
          startActionFlow(result.card, currentPlayer);
          return;
        }

        // Vérifier si le joueur a sauté (doublon)
        if (!currentPlayer.state) {
          setMessage(`💥 DOUBLON ! ${currentPlayer.name} est éliminé de la manche.`);
          setTimeout(nextPlayer, 2000);
        } else if (match.game.roundEnded) {
          setMessage(`🎉 FLIP 7 ou Fin de manche !`);
          setTimeout(endRound, 2000);
        } else {
          setMessage(`${currentPlayer.name} a pioché.`);
          setForceUpdate(n => n + 1); // Rafraichir l'affichage
          setTimeout(nextPlayer, 1500);
        }
      }
    }
  };

  // --- VUES ---

  if (viewState === 'summary') {
    // Vérifier victoire globale
    const winner = match.getWinner();
    if (winner) {
      return /*#__PURE__*/React.createElement(Box, {
        flexDirection: "column",
        alignItems: "center",
        borderColor: "green",
        borderStyle: "double",
        padding: 2
      }, /*#__PURE__*/React.createElement(Text, {
        bold: true,
        color: "green"
      }, "\uD83C\uDFC6 VICTOIRE FINALE \uD83C\uDFC6"), /*#__PURE__*/React.createElement(Text, null, "Le vainqueur est ", winner.name, " avec ", match.scores.get(winner.player_nb), " points !"));
    }
    return /*#__PURE__*/React.createElement(Box, {
      flexDirection: "column"
    }, /*#__PURE__*/React.createElement(Header, {
      title: "Fin de la manche"
    }), match.players.map(p => /*#__PURE__*/React.createElement(Box, {
      key: p.player_nb,
      justifyContent: "space-between",
      width: 40
    }, /*#__PURE__*/React.createElement(Text, null, p.name), /*#__PURE__*/React.createElement(Text, {
      bold: true
    }, "Total: ", match.scores.get(p.player_nb), " (+", p.score, ")"))), /*#__PURE__*/React.createElement(Box, {
      marginTop: 2
    }, /*#__PURE__*/React.createElement(Text, {
      color: "gray"
    }, "Appuyez sur Entr\xE9e pour la manche suivante..."), /*#__PURE__*/React.createElement(SelectInput, {
      items: [{
        label: 'Manche suivante',
        value: 'next'
      }],
      onSelect: () => {
        match.startRound();
        setCurrentPlayerIndex(match.dealerIndex);
        setLastDrawnCards([]);
        setViewState('menu');
        setMessage("Nouvelle manche !");
      }
    })));
  }
  if (viewState === 'hand') {
    return /*#__PURE__*/React.createElement(Box, {
      flexDirection: "column"
    }, /*#__PURE__*/React.createElement(Header, {
      title: `Main de ${currentPlayer.name}`
    }), /*#__PURE__*/React.createElement(Text, null, "Nombres: ", currentPlayer.hand_number.map(c => c).join(', ')), /*#__PURE__*/React.createElement(Text, null, "Bonus: ", currentPlayer.hand_bonus.join(', ')), /*#__PURE__*/React.createElement(Text, null, "Actions: ", currentPlayer.hand_actions.join(', ')), /*#__PURE__*/React.createElement(Box, {
      marginTop: 1
    }, /*#__PURE__*/React.createElement(SelectInput, {
      items: [{
        label: 'Retour',
        value: 'back'
      }],
      onSelect: () => setViewState('menu')
    })));
  }
  if (viewState === 'target_selection' && pendingAction) {
    return /*#__PURE__*/React.createElement(Box, {
      flexDirection: "column"
    }, /*#__PURE__*/React.createElement(Header, {
      title: "Choisir une cible"
    }), /*#__PURE__*/React.createElement(Text, null, pendingAction.owner.name || `Player ${pendingAction.owner.player_nb}`, " a tir\xE9 la carte."), /*#__PURE__*/React.createElement(Text, null, "Action: ", pendingAction.card), /*#__PURE__*/React.createElement(Box, {
      marginTop: 1
    }, /*#__PURE__*/React.createElement(SelectInput, {
      items: pendingAction.targets.map(p => ({
        label: p.name || `Player ${p.player_nb}`,
        value: p
      })),
      onSelect: item => {
        const action = pendingAction;
        setPendingAction(null);
        if (item?.value) {
          setMessage(`Action ${action.card} sur ${item.value.name}.`);
        }
        resolveActionChain(action.card, action.owner, item?.value || null);
      }
    })));
  }

  // Vue Principale (Menu de jeu)
  const menuItems = [{
    label: 'Piocher une carte',
    value: 'draw'
  }, {
    label: 'Voir ma main',
    value: 'view_hand'
  }, {
    label: "Besoin d'un coup de pouce ?",
    value: 'helper'
  }, {
    label: 'Stop (S\'arrêter)',
    value: 'stop'
  }];
  return /*#__PURE__*/React.createElement(Box, {
    flexDirection: "column",
    width: "100%"
  }, /*#__PURE__*/React.createElement(Box, {
    justifyContent: "space-between",
    width: "100%"
  }, /*#__PURE__*/React.createElement(Text, null, "Joueur: ", /*#__PURE__*/React.createElement(Text, {
    bold: true,
    color: "cyan"
  }, currentPlayer.name)), /*#__PURE__*/React.createElement(Text, null, "Score manche: ", currentPlayer.pointInMyHand())), /*#__PURE__*/React.createElement(Box, {
    borderStyle: "single",
    padding: 1,
    marginY: 1,
    flexDirection: "column",
    alignItems: "center"
  }, /*#__PURE__*/React.createElement(Text, {
    italic: true,
    color: "gray"
  }, message), lastDrawnCards.length > 0 && /*#__PURE__*/React.createElement(Box, {
    marginTop: 1,
    flexDirection: "row",
    flexWrap: "wrap",
    gap: 1
  }, lastDrawnCards.map((card, index) => /*#__PURE__*/React.createElement(CardVisual, {
    key: `${card}-${index}`,
    card: card
  })))), /*#__PURE__*/React.createElement(Text, null, "Que voulez-vous faire ?"), /*#__PURE__*/React.createElement(SelectInput, {
    items: menuItems,
    onSelect: handleAction
  }));
};