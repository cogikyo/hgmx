import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const person = el.querySelector('path[d^="M2 21"]');
  animate(person, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const shield = el.querySelector('path[d^="M22 17.5"]');
  animate(shield, { translateY: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
