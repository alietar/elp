import React from 'react';
import { Box, Text } from 'ink';
export const CardText = ({
  card
}) => {
  const cardStr = card.toString();
  let color = 'white';
  if (['0', '3', '6', '7'].includes(cardStr)) color = 'magenta';else if (['1', '12'].includes(cardStr)) color = 'gray';else if (['2', '5', '8'].includes(cardStr)) color = 'green';else if (['4', '11'].includes(cardStr)) color = 'blue';else if (['9'].includes(cardStr)) color = 'yellow';else if (['+2', '+4', '+6', '+8', '+10', 'x2', 'Flip Three', 'Freeze', 'Second Chance'].includes(cardStr) || parseInt(cardStr) < 0) color = 'red';
  return /*#__PURE__*/React.createElement(Text, {
    bold: true,
    color: color
  }, cardStr);
};
export const CardVisual = ({
  card
}) => {
  if (!card) return /*#__PURE__*/React.createElement(Box, {
    height: 5,
    width: 14,
    borderStyle: "single",
    justifyContent: "center",
    alignItems: "center"
  }, /*#__PURE__*/React.createElement(Text, null, "?"));
  const cardStr = card.toString();
  // Centrage approximatif

  return /*#__PURE__*/React.createElement(Box, {
    borderStyle: "double",
    borderColor: "white",
    flexDirection: "column",
    width: 15,
    paddingY: 5,
    alignItems: "center"
  }, /*#__PURE__*/React.createElement(CardText, {
    card: card
  }));
};
export const Header = ({
  title
}) => /*#__PURE__*/React.createElement(Box, {
  borderStyle: "round",
  borderColor: "cyan",
  paddingX: 1,
  marginBottom: 1
}, /*#__PURE__*/React.createElement(Text, {
  bold: true
}, title));