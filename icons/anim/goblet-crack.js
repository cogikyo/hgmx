import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const goblet = el.querySelector("g");
  animate(goblet, { rotate: [0, -8, 6, -3, 0] }, { duration: τmid });

  const crack = el.querySelector('path[d^="m13 11"]');
  animate(crack, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
