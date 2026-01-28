function doIHaveToDraw(draw, myHand) {
  let count_same_card = 0;
  for (let card of myHand) {
    for (let [key, value] of draw) {
      if (key === card) {
        count_same_card += value["quantity"];
      }
    }
  }
  let tot = 0;
  for (let [key, value] of draw) {
    tot += value["quantity"];
  }
  let proba = count_same_card / tot;
  return proba;
}
export { doIHaveToDraw };