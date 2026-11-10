import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const letter = el.querySelectorAll('path[d^="m4 16"], path[d="M6 12h8"]');
  animate(letter, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const tick = el.querySelector('path[d^="m20 15"]');
  animate(tick, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
