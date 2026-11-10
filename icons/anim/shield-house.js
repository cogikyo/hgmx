import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const shield = el.querySelector('path[d^="M20 10.418"]');
  animate(shield, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const house = el.querySelector('path[d^="M15 22"]');
  animate(house, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmin });

  const door = el.querySelector('path[d="M18 22v-3"]');
  animate(door, { scaleY: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
