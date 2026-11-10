import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const hand = el.querySelectorAll('path[d^="m15.4"], path[d^="M2 14"], path[d^="M5 14"]');
  animate(hand, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const cloche = el.querySelectorAll('path[d^="M5 10"], path[d^="M12 3"]');
  animate(cloche, { translateY: [0, -3, 0] }, { duration: τmid, delay: τmin });
}
