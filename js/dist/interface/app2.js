// app.js
import React, { useState } from 'react';
import { Text, Box } from 'ink';
import BigText from 'ink-big-text';
import { PlayerCountSetup, PlayerNameSetup } from './setup.js';
import { GameController } from './game_ui.js';
export default function App() {
  const [step, setStep] = useState('count'); // 'count', 'names', 'game'
  const [playerCount, setPlayerCount] = useState(0);
  const [playerNames, setPlayerNames] = useState([]);
  return /*#__PURE__*/React.createElement(MainLayout, null, /*#__PURE__*/React.createElement(Box, {
    marginBottom: 2
  }, /*#__PURE__*/React.createElement(BigText, {
    text: "Flip 7",
    font: "tiny"
  })), step === 'count' && /*#__PURE__*/React.createElement(PlayerCountSetup, {
    onSubmit: num => {
      setPlayerCount(num);
      setStep('names');
    }
  }), step === 'names' && /*#__PURE__*/React.createElement(PlayerNameSetup, {
    playerCount: playerCount,
    onFinished: names => {
      setPlayerNames(names);
      setStep('game');
    }
  }), step === 'game' && /*#__PURE__*/React.createElement(GameController, {
    playerCount: playerCount,
    playerNames: playerNames
  }));
}
function MainLayout({
  children
}) {
  return /*#__PURE__*/React.createElement(Box, {
    borderStyle: "round",
    padding: 2,
    flexDirection: "column",
    width: "100%"
  }, children);
}