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

	return (
		<MainLayout>
			<Box marginBottom={2}>
        <BigText text="Flip 7" font="tiny" />
      </Box>
            
			{step === 'count' && (
				<PlayerCountSetup
					onSubmit={(num) => {
						setPlayerCount(num);
						setStep('names');
					}}
				/>
			)}
            
			{step === 'names' && (
				<PlayerNameSetup
          playerCount={playerCount}
          onFinished={(names) => {
            setPlayerNames(names);
            setStep('game');
          }}
				/>
			)}
            
      {step === 'game' && (
        <GameController 
          playerCount={playerCount} 
          playerNames={playerNames}
        />
      )}
		</MainLayout>
	);
}

function MainLayout({children}) {
  return (
      <Box
        borderStyle="round"
        padding={2}
        flexDirection="column"
        width="100%"
      >
        {children}
      </Box>
  );
}