export var packet = new Map();
packet.set('12', {
  type: 'number',
  quantity: 12
});
packet.set('11', {
  type: 'number',
  quantity: 11
});
packet.set('10', {
  type: 'number',
  quantity: 10
});
packet.set('9', {
  type: 'number',
  quantity: 9
});
packet.set('8', {
  type: 'number',
  quantity: 8
});
packet.set('7', {
  type: 'number',
  quantity: 7
});
packet.set('6', {
  type: 'number',
  quantity: 6
});
packet.set('5', {
  type: 'number',
  quantity: 5
});
packet.set('4', {
  type: 'number',
  quantity: 4
});
packet.set('3', {
  type: 'number',
  quantity: 3
});
packet.set('2', {
  type: 'number',
  quantity: 2
});
packet.set('1', {
  type: 'number',
  quantity: 1
});
packet.set('0', {
  type: 'number',
  quantity: 1
});
packet.set('+2', {
  type: 'modifier',
  quantity: 1
});
packet.set('+4', {
  type: 'modifier',
  quantity: 1
});
packet.set('+6', {
  type: 'modifier',
  quantity: 1
});
packet.set('+8', {
  type: 'modifier',
  quantity: 1
});
packet.set('+10', {
  type: 'modifier',
  quantity: 1
});
packet.set('x2', {
  type: 'modifier',
  quantity: 1
});
packet.set('Flip Three', {
  type: 'action',
  quantity: 3
});
packet.set('Freeze', {
  type: 'action',
  quantity: 3
});
packet.set('Second Chance', {
  type: 'action',
  quantity: 3
});

// création de la défausse
export var defausse = new Map();
for (const [key, value] of packet.entries()) {
  defausse.set(key, {
    type: value.type,
    quantity: 0
  });
}
export var cardsTypes = ['12', '11', '10', '9', '8', '7', '6', '5', '4', '3', '2', '1', '0', '+2', '+4', '+6', '+8', '+10', 'x2', 'Flip Three', 'Freeze', 'Second Chance'];