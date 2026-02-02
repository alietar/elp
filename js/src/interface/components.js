import React from 'react';
import { Box, Text } from 'ink';

export const CardText = ({ card }) => {
    const cardStr = card.toString();
    let color = 'white';

    if (['0', '3', '6', '7'].includes(cardStr)) color = 'magenta';
    else if (['1', '12'].includes(cardStr)) color = 'gray';
    else if (['2', '5', '8'].includes(cardStr)) color = 'green';
    else if (['4', '11'].includes(cardStr)) color = 'blue';
    else if (['9'].includes(cardStr)) color = 'yellow';
    else if (['+2', '+4', '+6', '+8', '+10', 'x2', 'Flip Three', 'Freeze', 'Second Chance'].includes(cardStr) || parseInt(cardStr) < 0) color = 'red';

    return <Text bold color={color}>{cardStr}</Text>;
};

export const CardVisual = ({ card }) => {
    if (!card) return <Box height={10} width={14} borderStyle="single" justifyContent="center" alignItems="center"><Text>?</Text></Box>;

    const cardStr = card.toString();
    // Centrage approximatif

    return (
        <Box borderStyle="double" borderColor="white" flexDirection="column" flexBasis="10" paddingY="3" alignItems="center">
            <CardText card={card} />
        </Box>
    );
};

export const Header = ({ title }) => (
    <Box borderStyle="round" borderColor="cyan" paddingX={1} marginBottom={1}>
        <Text bold>{title}</Text>
    </Box>
);