import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const face = el.querySelector('path[d^="M13.227"]');
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const hands = el.querySelector('path[d^="M12 6v6"]');
  animate(hands, { rotate: [0, 360] }, { duration: τmax });
  const arrow = el.querySelectorAll('path[d^="m14 18 4-4"], path[d^="M18 22v-8"]');
  animate(arrow, { translateY: [3, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
