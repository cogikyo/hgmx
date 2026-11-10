import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const cards = el.querySelector("g");
  animate(cards, { rotate: [0, -6, 3, 0] }, { duration: τmax });
  const suit = el.querySelector('path[d^="M14.832"]');
  animate(suit, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
