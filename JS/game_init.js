
global.packet_init = new Map();

global.packet_init.set('12', { type: 'number', quantity: 12 });
global.packet_init.set('11', { type: 'number', quantity: 11 });
global.packet_init.set('10', { type: 'number', quantity: 10 });
global.packet_init.set('9', { type: 'number', quantity: 9 });
global.packet_init.set('8', { type: 'number', quantity: 8 });
global.packet_init.set('7', { type: 'number', quantity: 7 });
global.packet_init.set('6', { type: 'number', quantity: 6 });
global.packet_init.set('5', { type: 'number', quantity: 5 });
global.packet_init.set('4', { type: 'number', quantity: 4 });
global.packet_init.set('3', { type: 'number', quantity: 3 });
global.packet_init.set('2', { type: 'number', quantity: 2 });
global.packet_init.set('1', { type: 'number', quantity: 1 });
global.packet_init.set('0', { type: 'number', quantity: 1 }); 


global.packet_init.set('+2', { type: 'modifier', quantity: 1 });
global.packet_init.set('+4', { type: 'modifier', quantity: 1 });
global.packet_init.set('+6', { type: 'modifier', quantity: 1 });
global.packet_init.set('+8', { type: 'modifier', quantity: 1 });
global.packet_init.set('+10', { type: 'modifier', quantity: 1 });

global.packet_init.set('Flip Three', { type: 'action', quantity: 3 });
global.packet_init.set('Freeze', { type: 'action', quantity: 3 });
global.packet_init.set('Second Chance', { type: 'action', quantity: 3 });

console.log(global.packet_init);

global.card = [
  '12', '11', '10', '9', '8', '7', '6', '5', '4', '3', '2', '1', '0',
  '+2', '+4', '+6', '+8', '+10',
  'Flip Three', 'Freeze', 'Second Chance'
];


