import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const piece = el.querySelectorAll('path[d^="M16.4"], path[d^="M2 9"]');
  animate(piece, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const topping = el.querySelector("g");
  animate(topping, { translateY: [-3, 0.5, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
