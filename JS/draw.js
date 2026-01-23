export function draw_card(packet, cardsTypes) {

    if (cardsTypes.length === 0) 
        return null;

    const index = Math.floor(Math.random() * cardsTypes.length);
    const card = cardsTypes[index];
    const cardData = packet.get(card);
    cardData.quantity -= 1;    
    if (cardData.quantity === 0){
        const i = cardsTypes.indexOf(card);
        cardsTypes.splice(i, 1);
    }
    return card;
}

