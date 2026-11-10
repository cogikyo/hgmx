import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M16 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const hands = el.querySelector('path[d^="M8 14"]');
  animate(hands, { rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
