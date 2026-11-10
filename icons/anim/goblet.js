import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const goblet = el.querySelector("g");
  animate(goblet, { rotate: [0, -8, 6, -3, 0] }, { duration: τmid });

  const wine = el.querySelector('path[d^="M5 11"]');
  animate(wine, { rotate: [0, 12, -8, 0] }, { duration: τmid, delay: τmin });
}
