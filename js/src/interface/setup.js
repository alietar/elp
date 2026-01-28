import {Text, Box} from 'ink';
import TextInput from 'ink-text-input';
import React, { useState } from 'react';

export const PlayerCountSetup = ({ onSubmit }) => {
	const [count, setCount] = useState('');

	return (
		<Box flexDirection="column">
			<Text>Combien de joueurs vont participer ?</Text>
			<Box borderStyle="round" borderColor="cyan">
				<Text>Nombre : </Text>
				<TextInput
					value={count}
					onChange={setCount}
					onSubmit={(val) => {
                        const num = parseInt(val, 10);
                        if (!isNaN(num) && num > 0) onSubmit(num);
                    }}
				/>
			</Box>
            <Text color="gray">(Appuyez sur Entrée pour valider)</Text>
		</Box>
	);
};

export const PlayerNameSetup = ({ playerCount, onFinished }) => {
	const [names, setNames] = useState([]);
	const [currentInput, setCurrentInput] = useState('');

	const currentIndex = names.length;

    // Si on a tous les noms, on ne rend rien (le useEffect ou le parent gérera la suite),
    // mais par sécurité on peut renvoyer null ici.
    if (names.length >= playerCount) {
        return <Text color="green">Configuration terminée !</Text>;
    }

	return (
		<Box flexDirection="column">
			<Text>Entrez le nom du Joueur {currentIndex + 1} / {playerCount}</Text>
			<Box borderStyle="round" borderColor="green">
				<Text>Nom : </Text>
				<TextInput
					value={currentInput}
					onChange={setCurrentInput}
					onSubmit={(val) => {
						if (val.trim() === '') return;
						const newNames = [...names, val];
						setNames(newNames);
						setCurrentInput(''); // Reset du champ
                        
                        // Si c'était le dernier joueur, on remonte l'info au parent
						if (newNames.length === playerCount) {
							onFinished(newNames);
						}
					}}
				/>
			</Box>
		</Box>
	);
};