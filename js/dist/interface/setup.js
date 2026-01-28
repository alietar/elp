import { Text, Box } from 'ink';
import TextInput from 'ink-text-input';
import React, { useState } from 'react';
export const PlayerCountSetup = ({
  onSubmit
}) => {
  const [count, setCount] = useState('');
  return /*#__PURE__*/React.createElement(Box, {
    flexDirection: "column"
  }, /*#__PURE__*/React.createElement(Text, null, "Combien de joueurs vont participer ?"), /*#__PURE__*/React.createElement(Box, {
    borderStyle: "round",
    borderColor: "cyan"
  }, /*#__PURE__*/React.createElement(Text, null, "Nombre : "), /*#__PURE__*/React.createElement(TextInput, {
    value: count,
    onChange: setCount,
    onSubmit: val => {
      const num = parseInt(val, 10);
      if (!isNaN(num) && num > 0) onSubmit(num);
    }
  })), /*#__PURE__*/React.createElement(Text, {
    color: "gray"
  }, "(Appuyez sur Entr\xE9e pour valider)"));
};
export const PlayerNameSetup = ({
  playerCount,
  onFinished
}) => {
  const [names, setNames] = useState([]);
  const [currentInput, setCurrentInput] = useState('');
  const currentIndex = names.length;

  // Si on a tous les noms, on ne rend rien (le useEffect ou le parent gérera la suite),
  // mais par sécurité on peut renvoyer null ici.
  if (names.length >= playerCount) {
    return /*#__PURE__*/React.createElement(Text, {
      color: "green"
    }, "Configuration termin\xE9e !");
  }
  return /*#__PURE__*/React.createElement(Box, {
    flexDirection: "column"
  }, /*#__PURE__*/React.createElement(Text, null, "Entrez le nom du Joueur ", currentIndex + 1, " / ", playerCount), /*#__PURE__*/React.createElement(Box, {
    borderStyle: "round",
    borderColor: "green"
  }, /*#__PURE__*/React.createElement(Text, null, "Nom : "), /*#__PURE__*/React.createElement(TextInput, {
    value: currentInput,
    onChange: setCurrentInput,
    onSubmit: val => {
      if (val.trim() === '') return;
      const newNames = [...names, val];
      setNames(newNames);
      setCurrentInput(''); // Reset du champ

      // Si c'était le dernier joueur, on remonte l'info au parent
      if (newNames.length === playerCount) {
        onFinished(newNames);
      }
    }
  })));
};