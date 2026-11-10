import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const button = el.querySelector('line[y1="2"]');
  animate(button, { translateY: [0, 2, 0] }, { duration: τmin });
  const hand = el.querySelector('line[y1="14"]');
  animate(hand, { rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
