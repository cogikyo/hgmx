import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const hand = el.querySelectorAll('path:not([d^="m14.45"])');
  animate(hand, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const heart = el.querySelector('path[d^="m14.45"]');
  animate(heart, { scale: [1, 1.25, 0.9, 1.15, 1] }, { duration: τmid, delay: τmin });
}
