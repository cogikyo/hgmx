import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const face = el.querySelector('path[d^="M21.92"]');
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const hands = el.querySelector('path[d^="M12 6v6"]');
  animate(hands, { rotate: [0, 360] }, { duration: τmax });
  const plus = el.querySelectorAll('path[d="M16 19h6"], path[d="M19 16v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
