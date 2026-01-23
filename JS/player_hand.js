class Hand {
    constructor(hand, state){
        this.hand = [] || [];
        this.state = true || true;
    }

    pointInMyHand() {
        let sum = 0;
        for (let card of this.hand) {
            sum += card.value; 
        }
        return sum;
    }


}