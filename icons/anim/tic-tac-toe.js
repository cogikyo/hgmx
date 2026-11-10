import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const crosses = el.querySelectorAll('path[d^="m21 16"], path[d^="m21 21"], path[d^="M8 3"], path[d^="M8 8"]');
  animate(crosses, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid });
  const noughts = el.querySelectorAll("circle");
  animate(noughts, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
